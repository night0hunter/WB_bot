package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	keyboard "wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type prevCommandInfo struct {
	CommandName enum.BotCommandNameType
	MessageID   int
	Info        dto.WarehouseData
}

// var prevCommands = map[int64]enum.BotCommandNameType{}
var prevCommands = map[int64]prevCommandInfo{}

type Service interface {
	ButtonTypeWarehouseService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error
	ButtonTypeCoeffLimitService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error
	ButtonTypeSupplyTypeService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error
	ButtonTypeChangeService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error
	BotSlashCommandTypeHelpService(ctx context.Context, chatID int64) string
	BotSlashCommandTypeCheckService(ctx context.Context, chatID int64) ([]string, error)
	BotAnswerInputDateService(ctx context.Context, chatID int64, date string) (dto.TrackingDate, error)
	BotAnswerInputCoeffLimitService(ctx context.Context, chatID int64, coeffLimit string) (int, error)
	BotSlashCommandTypeChange(ctx context.Context, chatID int64) ([]dto.WarehouseData, error)
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

			err = h.ButtonHandler(ctx, update, buttonData)
			if err != nil {
				return errors.Wrap(err, "ButtonHandler")
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

func (h *handler) ButtonHandler(ctx context.Context, update tgbotapi.Update, buttonData dto.ButtonData) error {
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
	case enum.ButtonTypeUserTrackings:
		err := h.ButtonTypeUserTrackingsHandler(ctx, update, buttonData)
		if err != nil {
			return errors.Wrap(err, "ButtonTypeUserTrackingsHandler")
		}
	}

	return nil
}

func (h *handler) ButtonTypeUserTrackingsHandler(ctx context.Context, update tgbotapi.Update, buttonData dto.ButtonData) error {
	err := h.service.ButtonTypeChangeService(ctx, update.CallbackQuery.Message.Chat.ID, buttonData)
	if err != nil {
		return errors.Wrap(err, "service.ButtonTypeChangeService")
	}

	prevCommand, ok := prevCommands[update.CallbackQuery.Message.Chat.ID]
	if ok {
		deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, prevCommand.MessageID)
		if _, err := h.bot.Send(deleteMsg); err != nil {
			fmt.Printf("bot.Send(deleteMsg): %s", err.Error())
		}
	}

	prevCommands[update.CallbackQuery.Message.Chat.ID] = prevCommandInfo{}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Статус отслеживания успешно изменён")
	_, err = h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) ButtonTypeWarehouseHandler(ctx context.Context, update tgbotapi.Update, buttonData dto.ButtonData) error {
	err := h.service.ButtonTypeWarehouseService(ctx, update.CallbackQuery.Message.Chat.ID, buttonData)
	if err != nil {
		return errors.Wrap(err, "service.ButtonTypeWarehouseService")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
		"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.FromDate.Format(dto.TimeFormat),
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.ToDate.Format(dto.TimeFormat),
		constmsg.WarehouseNames[buttonData.Value],
		"",
		"",
		BotCommands[enum.BotCommandNameTypeInputCoeffLimit],
	))

	msg, err = keyboard.DrawCoeffKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawCoeffKeyboard")
	}

	prevCommand, ok := prevCommands[update.CallbackQuery.Message.Chat.ID]
	if ok {
		deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, prevCommand.MessageID)
		if _, err := h.bot.Send(deleteMsg); err != nil {
			fmt.Printf("bot.Send(deleteMsg): %s", err.Error())
		}
	}

	if buttonData.Value == -1 {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputDate])

		message, err := h.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "bot.Send")
		}

		prevCommands[update.CallbackQuery.Message.Chat.ID] = prevCommandInfo{
			CommandName: enum.BotCommandNameTypeInputDate,
			MessageID:   message.MessageID,
		}

		return nil
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	prevCommands[update.CallbackQuery.Message.Chat.ID] = prevCommandInfo{
		CommandName: enum.BotCommandNameTypeInputCoeffLimit,
		MessageID:   message.MessageID,
		Info: dto.WarehouseData{
			FromDate:      prevCommand.Info.FromDate,
			ToDate:        prevCommand.Info.ToDate,
			WarehouseName: constmsg.WarehouseNames[buttonData.Value],
		},
	}

	return nil
}

func (h *handler) ButtonTypeCoeffLimitHandler(ctx context.Context, update tgbotapi.Update, buttonData dto.ButtonData) error {
	err := h.service.ButtonTypeCoeffLimitService(ctx, update.CallbackQuery.Message.Chat.ID, buttonData)
	if err != nil {
		return errors.Wrap(err, "service.ButtonTypeCoeffLimitService")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
		"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %dx и меньше\nТип поставки: %s\n---------------\n%s",
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.FromDate.Format(dto.TimeFormat),
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.ToDate.Format(dto.TimeFormat),
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.WarehouseName,
		buttonData.Value,
		"",
		BotCommands[enum.BotCommandNameTypeInputSupplyType],
	))

	// msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputSupplyType])
	msg, err = keyboard.DrawSupplyKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawSupplyKeyboard")
	}

	prevCommand, ok := prevCommands[update.CallbackQuery.Message.Chat.ID]
	if ok {
		deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, prevCommand.MessageID)
		if _, err := h.bot.Send(deleteMsg); err != nil {
			fmt.Printf("bot.Send(deleteMsg): %s", err.Error())
		}
	}

	if buttonData.Value == -1 {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
			"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
			prevCommands[update.CallbackQuery.Message.Chat.ID].Info.FromDate.Format(dto.TimeFormat),
			prevCommands[update.CallbackQuery.Message.Chat.ID].Info.ToDate.Format(dto.TimeFormat),
			"",
			"",
			"",
			BotCommands[enum.BotCommandNameTypeInputWarehouse],
		))
		msg, err = keyboard.DrawWarehouseKeyboard(msg)
		if err != nil {
			return errors.Wrap(err, "keyboard.DrawWarehouseKeyboard")
		}

		message, err := h.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "bot.Send")
		}

		prevCommands[update.CallbackQuery.Message.Chat.ID] = prevCommandInfo{
			CommandName: enum.BotCommandNameTypeInputWarehouse,
			MessageID:   message.MessageID,
			Info: dto.WarehouseData{
				FromDate: prevCommand.Info.FromDate,
				ToDate:   prevCommand.Info.ToDate,
			},
		}

		return nil
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	prevCommands[update.CallbackQuery.Message.Chat.ID] = prevCommandInfo{
		CommandName: enum.BotCommandNameTypeInputSupplyType,
		MessageID:   message.MessageID,
		Info: dto.WarehouseData{
			FromDate:      prevCommand.Info.FromDate,
			ToDate:        prevCommand.Info.ToDate,
			WarehouseName: prevCommand.Info.WarehouseName,
			CoeffLimit:    buttonData.Value,
		},
	}

	return nil
}

func (h *handler) ButtonTypeSupplyTypeHandler(ctx context.Context, update tgbotapi.Update, buttonData dto.ButtonData) error {
	prevCommand, ok := prevCommands[update.CallbackQuery.Message.Chat.ID]
	if ok {
		deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, prevCommand.MessageID)
		if _, err := h.bot.Send(deleteMsg); err != nil {
			fmt.Printf("bot.Send(deleteMsg): %s", err.Error())
		}
	}

	if buttonData.Value == -1 {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
			"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
			prevCommands[update.CallbackQuery.Message.Chat.ID].Info.FromDate.Format(dto.TimeFormat),
			prevCommands[update.CallbackQuery.Message.Chat.ID].Info.ToDate.Format(dto.TimeFormat),
			prevCommands[update.CallbackQuery.Message.Chat.ID].Info.WarehouseName,
			"",
			"",
			BotCommands[enum.BotCommandNameTypeInputCoeffLimit],
		))
		msg, err := keyboard.DrawCoeffKeyboard(msg)
		if err != nil {
			return errors.Wrap(err, "keyboard.DrawCoeffKeyboard")
		}

		message, err := h.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "bot.Send")
		}

		prevCommands[update.CallbackQuery.Message.Chat.ID] = prevCommandInfo{
			CommandName: enum.BotCommandNameTypeInputCoeffLimit,
			MessageID:   message.MessageID,
			Info: dto.WarehouseData{
				FromDate:      prevCommand.Info.FromDate,
				ToDate:        prevCommand.Info.ToDate,
				WarehouseName: prevCommand.Info.WarehouseName,
			},
		}

		return nil
	}

	err := h.service.ButtonTypeSupplyTypeService(ctx, update.CallbackQuery.Message.Chat.ID, buttonData)
	if err != nil {
		return errors.Wrap(err, "service.ButtonTypeSupplyTypeService")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
		"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %dx и меньше\nТип поставки: %s\n---------------\n%s",
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.FromDate.Format(dto.TimeFormat),
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.ToDate.Format(dto.TimeFormat),
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.WarehouseName,
		prevCommands[update.CallbackQuery.Message.Chat.ID].Info.CoeffLimit,
		constmsg.SupplyTypes[strconv.Itoa(buttonData.Value)],
		"Склад успешно добавлен",
	))

	// msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Склад успешно добавлен!")

	_, err = h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	prevCommands[update.CallbackQuery.Message.Chat.ID] = prevCommandInfo{}

	return nil
}

func (h *handler) messageHandler(ctx context.Context, update tgbotapi.Update) error {
	switch update.Message.Text {
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeHelp]:
		err := h.BotSlashCommandTypeHelpHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeHelpHandler")
		}
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeAdd]:
		err := h.BotSlashCommandTypeAddHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeAddHandler")
		}
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeChange]:
		err := h.BotSlashCommandTypeChangeHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeChangeHandler")
		}
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeCheck]:
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
	text := h.service.BotSlashCommandTypeHelpService(ctx, update.Message.Chat.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) BotSlashCommandTypeAddHandler(ctx context.Context, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputDate])

	message, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	prevCommands[update.Message.Chat.ID] = prevCommandInfo{
		CommandName: enum.BotCommandNameTypeInputDate,
		MessageID:   message.MessageID,
	}

	return nil
}

func (h *handler) BotSlashCommandTypeCheckHandler(ctx context.Context, update tgbotapi.Update) error {
	whs, err := h.service.BotSlashCommandTypeCheckService(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.BotSlashCommandTypeCheck")
	}

	if whs == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("На данный момент вы не отслеживаете ни одного склада, чтобы добавить, используйте %s", constmsg.BotSlashCommands[enum.BotSlashCommandTypeAdd]))
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
func (h *handler) BotSlashCommandTypeChangeHandler(ctx context.Context, update tgbotapi.Update) error {
	warehouses, err := h.service.BotSlashCommandTypeChange(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.BotSlashCommandTypeChange")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите отслеживание из списка ниже, чтобы изменить его статус")
	msg = keyboard.DrawTrackingsKeyboard(msg, warehouses)

	message, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	prevCommands[update.Message.Chat.ID] = prevCommandInfo{
		CommandName: enum.BotCommandNameTypeChangeStatus,
		MessageID:   message.MessageID,
	}

	return nil
}

func (h *handler) BotSlashCommandTypeDefaultHandler(ctx context.Context, update tgbotapi.Update) error {
	prevCommand, ok := prevCommands[update.Message.Chat.ID]

	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know this command")
		if _, err := h.bot.Send(msg); err != nil {
			errors.Wrap(err, "bot.Send")
		}

		// return errors.New("Unknown command")
		return nil
	}

	switch prevCommand.CommandName {
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
	timeRange, err := h.service.BotAnswerInputDateService(ctx, update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		return errors.Wrap(err, "service.BotAnswerInputDateService")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
		timeRange.DateFrom.Format(dto.TimeFormat),
		timeRange.DateTo.Format(dto.TimeFormat),
		"",
		"",
		"",
		BotCommands[enum.BotCommandNameTypeInputWarehouse],
	))

	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputWarehouse])
	msg, err = keyboard.DrawWarehouseKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawWarehouseKeyboard")
	}

	prevCommand, ok := prevCommands[update.Message.Chat.ID]
	if ok {
		deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, prevCommand.MessageID)
		if _, err := h.bot.Send(deleteMsg); err != nil {
			fmt.Printf("bot.Send(deleteMsg): %s", err.Error())
		}
	}

	deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	if _, err := h.bot.Send(deleteMsg); err != nil {
		fmt.Printf("bot.Send(deleteMsg): %s", err.Error())
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	prevCommands[update.Message.Chat.ID] = prevCommandInfo{
		CommandName: enum.BotCommandNameTypeInputWarehouse,
		MessageID:   message.MessageID,
		Info:        dto.WarehouseData{FromDate: timeRange.DateFrom, ToDate: timeRange.DateTo},
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
	coeff, err := h.service.BotAnswerInputCoeffLimitService(ctx, update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		return errors.Wrap(err, "service.BotAnswerInputCoeffLimitService")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %dx и меньше\nТип поставки: %s\n%s",
		prevCommands[update.Message.Chat.ID].Info.FromDate.Format(dto.TimeFormat),
		prevCommands[update.Message.Chat.ID].Info.ToDate.Format(dto.TimeFormat),
		prevCommands[update.Message.Chat.ID].Info.WarehouseName,
		coeff,
		"",
		BotCommands[enum.BotCommandNameTypeInputSupplyType],
	))

	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputSupplyType])
	msg, err = keyboard.DrawSupplyKeyboard(msg)
	if err != nil {
		return errors.Wrap(err, "keyboard.DrawSupplyKeyboard")
	}

	prevCommand, ok := prevCommands[update.Message.Chat.ID]
	if ok {
		deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, prevCommand.MessageID)
		if _, err := h.bot.Send(deleteMsg); err != nil {
			fmt.Printf("bot.Send(deleteMsg): %s", err.Error())
		}
	}

	deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	if _, err := h.bot.Send(deleteMsg); err != nil {
		fmt.Printf("bot.Send(deleteMsg): %s", err.Error())
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	prevCommands[update.Message.Chat.ID] = prevCommandInfo{
		CommandName: enum.BotCommandNameTypeInputSupplyType,
		MessageID:   message.MessageID,
		Info: dto.WarehouseData{
			FromDate:   prevCommand.Info.FromDate,
			ToDate:     prevCommand.Info.ToDate,
			CoeffLimit: coeff,
		},
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
