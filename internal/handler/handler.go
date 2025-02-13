package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	keyboard "wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

var prevCommands = map[int64]enum.BotCommandNameType{}

type Service interface {
	ButtonTypeWarehouseService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error
	ButtonTypeCoeffLimitService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error
	ButtonTypeSupplyTypeService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error
	// BotSlashCommandTypeAddService(ctx context.Context, chatID int64) error
	// BotSlashCommandTypeStopService(ctx context.Context, chatID int64) error // ADD
	BotSlashCommandTypeCheckService(ctx context.Context, chatID int64) ([]string, error)
	BotAnswerInputDateService(ctx context.Context, chatID int64, date string) error
	// BotAnswerInputWarehouseService(ctx context.Context, chatID int64) error
	BotAnswerInputCoeffLimitService(ctx context.Context, chatID int64, coeffLimit string) error
	// BotAnswerInputSupplyTypeService(ctx context.Context, chatID int64) error
}

type handler struct {
	bot     *tgbotapi.BotAPI
	service Service
}

func NewHandler(bot *tgbotapi.BotAPI, svc Service) *handler {
	return &handler{bot: bot, service: svc}
}

func (h *handler) Run(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := h.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		select {
		case <-ctx.Done():
			return errors.New("context cancelled")
		default:
		}

		if update.CallbackQuery == nil && update.Message == nil {
			// TODO: info log
			continue
		}

		if update.CallbackQuery != nil {
			var buttonData dto.ButtonData

			err := json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
			if err != nil {
				return errors.Wrap(err, "json.Unmarshal")
			}

			switch buttonData.Type {
			case enum.ButtonTypeWarehouse:
				err := h.ButtonTypeWarehouseHandler(ctx, update, buttonData)
				if err != nil {
					return errors.Wrap(err, "ButtonTypeWarehouseHandler")
				}
			case enum.ButtonTypeCoeffLimit:
				err := h.ButtonTypeCoeffLimitHandler(ctx, update, buttonData)
				if err != nil {
					return errors.Wrap(err, "ButtonTypeCoeffLimitHandler")
				}
			case enum.ButtonTypeSupplyType:
				err := h.ButtonTypeSupplyTypeHandler(ctx, update, buttonData)
				if err != nil {
					return errors.Wrap(err, "ButtonTypeSupplyTypeHandler")
				}
			}

			continue
		}

		if update.Message != nil {
			err := h.messageHandler(ctx, update)
			if err != nil {
				return errors.Wrap(err, "messageHandler")
			}
		}

	}

	return nil
}

func (h *handler) ButtonTypeWarehouseHandler(ctx context.Context, update tgbotapi.Update, buttonData dto.ButtonData) error {
	prevCommands[update.CallbackQuery.Message.Chat.ID] = enum.BotCommandNameTypeInputCoeffLimit

	err := h.service.ButtonTypeWarehouseService(ctx, update.CallbackQuery.Message.Chat.ID, buttonData)
	if err != nil {
		return errors.Wrap(err, "service.ButtonTypeWarehouseService")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Вы выбрали склад %s", enum.WarehouseNames[buttonData.Value]))
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputCoeffLimit])
	msg, err = keyboard.DrawCoeffKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawCoeffKeyboard")
	}

	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) ButtonTypeCoeffLimitHandler(ctx context.Context, update tgbotapi.Update, buttonData dto.ButtonData) error {
	prevCommands[update.CallbackQuery.Message.Chat.ID] = enum.BotCommandNameTypeInputSupplyType

	err := h.service.ButtonTypeCoeffLimitService(ctx, update.CallbackQuery.Message.Chat.ID, buttonData)
	if err != nil {
		return errors.Wrap(err, "service.ButtonTypeCoeffLimitService")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Вы выбрали %dx", buttonData.Value))
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputSupplyType])
	msg, err = keyboard.DrawSupplyKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawSupplyKeyboard")
	}

	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) ButtonTypeSupplyTypeHandler(ctx context.Context, update tgbotapi.Update, buttonData dto.ButtonData) error {
	prevCommands[update.CallbackQuery.Message.Chat.ID] = enum.BotCommandNameTypeInputDate

	tmpTracking := dto.Trackings[update.CallbackQuery.Message.Chat.ID]
	tmpTracking.SupplyType = fmt.Sprint(buttonData.Value)
	dto.Trackings[update.CallbackQuery.Message.Chat.ID] = tmpTracking

	err := h.service.ButtonTypeSupplyTypeService(ctx, update.CallbackQuery.Message.Chat.ID, buttonData)
	if err != nil {
		return errors.Wrap(err, "service.ButtonTypeSupplyTypeService")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Вы выбрали тип поставки %d", buttonData.Value))
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Склад успешно добавлен!")
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) messageHandler(ctx context.Context, update tgbotapi.Update) error {
	switch update.Message.Text {
	case BotSlashCommands[enum.BotSlashCommandTypeHelp]:
		err := h.BotSlashCommandTypeHelpHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeHelpHandler")
		}
	case BotSlashCommands[enum.BotSlashCommandTypeAdd]:
		err := h.BotSlashCommandTypeAddHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeAddHandler")
		}
	case BotSlashCommands[enum.BotSlashCommandTypeStop]:

	case BotSlashCommands[enum.BotSlashCommandTypeCheck]:
		err := h.BotSlashCommandTypeCheckHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeCheckHandler")
		}

	default:
		err := h.BotSlashCommandTypeDefaultHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeDefaultHandler")
		}
	}

	return nil
}

func (h *handler) BotSlashCommandTypeHelpHandler(ctx context.Context, update tgbotapi.Update) error {
	var msg tgbotapi.MessageConfig

	// err := h.service.BotSlashCommandTypeHelpService(ctx, update.Message.Chat.ID)
	// if err != nil {
	// 	return errors.Wrap(err, "service.BotSlashCommandTypeHelpService")
	// }

	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) BotSlashCommandTypeAddHandler(ctx context.Context, update tgbotapi.Update) error {
	prevCommands[update.Message.Chat.ID] = enum.BotCommandNameTypeInputDate

	// err := h.service.BotSlashCommandTypeAddService(ctx, update.Message.Chat.ID)
	// if err != nil {
	// 	return errors.Wrap(err, "BotSlashCommandTypeAddService")
	// }

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputDate])
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) BotSlashCommandTypeCheckHandler(ctx context.Context, update tgbotapi.Update) error {
	whs, err := h.service.BotSlashCommandTypeCheckService(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.BotSlashCommandTypeCheck")
	}

	if whs == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("На данный момент вы не отслеживаете ни одного склада, чтобы добавить, используйте %s:", BotSlashCommands[enum.BotSlashCommandTypeCheck]))
		if _, err := h.bot.Send(msg); err != nil {
			return errors.Wrap(err, "bot.Send")
		}

		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список отслеживаемых складов:")
	if _, err := h.bot.Send(msg); err != nil {
		errors.Wrap(err, "bot.Send")
	}

	for _, wh := range whs {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, wh)
		if _, err := h.bot.Send(msg); err != nil {
			errors.Wrap(err, "bot.Send")
		}
	}

	return nil
}

// TODO: add /stop function
func (h *handler) BotSlashCommandTypeStopHandler(ctx context.Context, update tgbotapi.Update) error {
	// err := h.service.BotSlashCommandTypeStopService(ctx, update.Message.Chat.ID)
	// if err != nil {
	// 	return errors.Wrap(err, "service.BotSlashCommandTypeStopService")
	// }

	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите отслеживание из списка ниже, чтобы приостановить/удалить")
	// if _, err := h.bot.Send(keyboard.DrawTrackingsKeyboard(msg, warehouses)); err != nil {
	// 	errors.Wrap(err, "bot.Send")
	// }

	return nil
}

func (h *handler) BotSlashCommandTypeDefaultHandler(ctx context.Context, update tgbotapi.Update) error {
	prevCommand, ok := prevCommands[update.Message.Chat.ID]

	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know this command")
		if _, err := h.bot.Send(msg); err != nil {
			errors.Wrap(err, "bot.Send")
		}

		return errors.New("Unknown command")
	}

	switch prevCommand {
	case enum.BotCommandNameTypeInputDate:
		err := h.BotAnswerInputDateHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotAnswerInptuDateHandler")
		}
	case enum.BotCommandNameTypeInputWarehouse:
		err := h.BotAnswerInputWarehouseHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotAnswerInputWarehouseHandler")
		}
	case enum.BotCommandNameTypeInputCoeffLimit:
		err := h.BotAnswerInputCoeffLimitHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotAnswerInputCoeffLimitHandler")
		}
	case enum.BotCommandNameTypeInputSupplyType:
		err := h.BotAnswerInputSupplyType(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotAnswerInputSupplyType")
		}
	default:
		// never reaches
	}

	return nil
}

func (h *handler) BotAnswerInputDateHandler(ctx context.Context, update tgbotapi.Update) error {
	var msg tgbotapi.MessageConfig

	prevCommands[update.Message.Chat.ID] = enum.BotCommandNameTypeInputWarehouse

	err := h.service.BotAnswerInputDateService(ctx, update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		return errors.Wrap(err, "service.BotAnswerInputDateService")
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputWarehouse])
	if _, err := h.bot.Send((msg)); err != nil {
		errors.Wrap(err, "bot.Send")
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputWarehouse])
	msg, err = keyboard.DrawWarehouseKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawWarehouseKeyboard")
	}

	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) BotAnswerInputWarehouseHandler(ctx context.Context, update tgbotapi.Update) error {
	// err := h.service.BotAnswerInputWarehouseService(ctx, update.Message.Chat.ID)
	// if err != nil {
	// 	return errors.Wrap(err, "service.BotAnswerInputWarehouseService")
	// }

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, выберите нужный склад из предложенного списка")
	if _, err := h.bot.Send((msg)); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) BotAnswerInputCoeffLimitHandler(ctx context.Context, update tgbotapi.Update) error {
	prevCommands[update.Message.Chat.ID] = enum.BotCommandNameTypeInputSupplyType

	err := h.service.BotAnswerInputCoeffLimitService(ctx, update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		return errors.Wrap(err, "service.BotAnswerInputCoeffLimitService")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputSupplyType])
	msg, err = keyboard.DrawSupplyKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawSupplyKeyboard")
	}

	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) BotAnswerInputSupplyType(ctx context.Context, update tgbotapi.Update) error {
	// err := h.service.BotAnswerInputSupplyTypeService(ctx, update.Message.Chat.ID)
	// if err != nil {
	// 	return errors.Wrap(err, "service.BotAnswerInputSupplyTypeService")
	// }

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, выберите нужный тип поставки из предложенного списка")
	msg, err := keyboard.DrawSupplyKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawSupplyKeyboard")
	}

	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}
