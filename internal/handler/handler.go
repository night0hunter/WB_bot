package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
	"wb_bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Service interface {
	DeleteTrackingService(ctx context.Context, chatID int64, trackingID int64) error
	ChangeStatusService(ctx context.Context, chatID, trackingID int64) error
	BotSlashCommandTypeHelpService(ctx context.Context, chatID int64) string
	BotSlashCommandTypeCheckService(ctx context.Context, chatID int64) ([]string, error)
	BotAnswerInputDateService(ctx context.Context, chatID int64, date string) (dto.TrackingDate, error)
	BotAnswerInputCoeffLimitService(ctx context.Context, chatID int64, coeffLimit string) (int, error)
	BotSlashCommandTypeChange(ctx context.Context, chatID int64) ([]dto.WarehouseData, error)
	AddSequenceEndService(ctx context.Context, chatID int64, data []byte) error
	GetTrackings(ctx context.Context) ([]dto.MergedResp, error)
	KeepSendingTime(ctx context.Context, tracking dto.MergedResp) error

	SelectState(ctx context.Context, chatID int64) (dto.PrevCommandInfo, error)
	InsertState(ctx context.Context, chatID int64, prevCommand dto.PrevCommandInfo) error
	UpdateState(ctx context.Context, chatID int64, prevCommand dto.PrevCommandInfo) error
	DeleteState(ctx context.Context, chatID int64) error
}

// var prevCommands = map[int64]dto.PrevCommandInfo{}

type HandlerStruct interface {
	Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error)
	Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error)
	GetCommandName() enum.CommandSequences
}

type handler2 struct {
	bot      *tgbotapi.BotAPI
	service  Service
	handlers map[enum.Sequences]map[enum.CommandSequences]struct {
		Prev    HandlerStruct
		Current HandlerStruct
		Next    HandlerStruct
	}
}

var SequenceToFirstCommand = map[enum.Sequences]enum.CommandSequences{
	enum.Add:    enum.BotCommandNameTypeAdd,
	enum.Change: enum.BotCommandNameTypeChange,
}

func New(bot *tgbotapi.BotAPI, svc Service) *handler2 {
	handlers := map[enum.Sequences]map[enum.CommandSequences]struct {
		Prev    HandlerStruct
		Current HandlerStruct
		Next    HandlerStruct
	}{
		enum.Add: {
			enum.BotCommandNameTypeAdd: {
				Prev:    nil,
				Current: nil,
				Next:    &InputDateHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputDate},
			},
			enum.BotCommandNameTypeSaveStatus: {
				Prev:    nil,
				Current: &SaveStatusHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeSaveStatus},
				Next:    &InputDateHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputDate},
			},
			enum.BotCommandNameTypeInputDate: {
				Prev:    nil,
				Current: &InputDateHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputDate},
				Next:    &WarehouseHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputWarehouse},
			},
			enum.BotCommandNameTypeInputWarehouse: {
				Prev:    &InputDateHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputDate},
				Current: &WarehouseHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputWarehouse},
				Next:    &CoeffLimitHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputCoeffLimit},
			},
			enum.BotCommandNameTypeInputCoeffLimit: {
				Prev:    &WarehouseHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputWarehouse},
				Current: &CoeffLimitHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputCoeffLimit},
				Next:    &SupplyTypeHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputSupplyType},
			},
			enum.BotCommandNameTypeInputSupplyType: {
				Prev:    &CoeffLimitHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputCoeffLimit},
				Current: &SupplyTypeHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputSupplyType},
				Next:    nil,
			},
		},
		enum.Change: {
			enum.BotCommandNameTypeChange: {
				Prev:    nil,
				Current: nil,
				Next:    &TrackingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
			},
			enum.BotCommandNameTypeSaveStatus: {
				Prev:    nil,
				Current: &SaveStatusHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeSaveStatus},
				Next:    &TrackingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
			},
			enum.BotCommandNameTypeTracking: {
				Prev:    nil,
				Current: &TrackingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
				Next:    &ActionChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeAction},
			},
			enum.BotCommandNameTypeAction: {
				Prev:    &TrackingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
				Current: &ActionChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeAction},
				Next:    nil,
			},
		},
	}

	return &handler2{bot: bot, service: svc, handlers: handlers}
}

func (h *handler2) Run(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := h.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// select {
		// case <-ctx.Done():
		// 	return errors.New("context cancelled")
		// default:
		// }

		if update.CallbackQuery == nil && update.Message == nil {
			// TODO: info log
			continue
		}

		if update.CallbackQuery != nil {
			err := h.ButtonHandler(ctx, update)
			if err != nil {
				// return errors.Wrap(err, "ButtonHandler")
				fmt.Println(errors.Wrap(err, "h.messageHandler2"))
			}

			continue
		}

		if update.Message != nil {
			err := h.messageHandler2(ctx, update)
			if err != nil {
				// return errors.Wrap(err, "h.messageHandler2")
				fmt.Println(errors.Wrap(err, "h.messageHandler2"))
			}
		}

	}

	return nil
}

func (h *handler2) messageHandler2(ctx context.Context, update tgbotapi.Update) error {
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

func (h *handler2) BotSlashCommandTypeHelpHandler(ctx context.Context, update tgbotapi.Update) error {
	text := h.service.BotSlashCommandTypeHelpService(ctx, update.Message.Chat.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler2) BotSlashCommandTypeAddHandler(ctx context.Context, update tgbotapi.Update) error {
	deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	_, err := h.bot.Send(deleteMsg)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
		fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
	}

	state, err := h.service.SelectState(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.SelectState")
	}

	if state.Info != nil {
		deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, state.MessageID)
		_, err := h.bot.Send(deleteMsg)
		if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
			fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
		}
	}

	if state.Info != nil && state.SequenceName == enum.Add {
		if state.CommandName != enum.BotCommandNameTypeSaveStatus {
			// state.CommandName = SequenceToFirstCommand[state.SequenceName]
			if h.handlers[state.SequenceName][state.CommandName].Prev != nil {
				state.CommandName = h.handlers[state.SequenceName][state.CommandName].Prev.GetCommandName()
			}
		}

		prevCommand, err := h.handlers[state.SequenceName][enum.BotCommandNameTypeSaveStatus].Current.Question(ctx, update, state)
		if err != nil {
			return errors.Wrap(err, "handlers[state.SequenceName][enum.BotCommandNameTypeSaveStatus].Current.Question")
		}

		err = h.service.UpdateState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
			SequenceName: state.SequenceName,
			CommandName:  enum.BotCommandNameTypeSaveStatus,
			MessageID:    prevCommand.MessageID,
			Info:         prevCommand.Info,
		})
		if err != nil {
			return errors.Wrap(err, "service.InsertState")
		}

		return nil
	}

	err = h.service.DeleteState(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.DeleteState")
	}

	prevCommand, err := h.handlers[enum.Add][SequenceToFirstCommand[enum.Add]].Next.Question(ctx, update, dto.PrevCommandInfo{})
	if err != nil {
		return errors.Wrap(err, "handlers[enum.Add][SequenceToFirstCommand[enum.Add]].Next.Question")
	}

	err = h.service.InsertState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
		SequenceName: enum.Add,
		CommandName:  enum.BotCommandNameTypeInputDate,
		MessageID:    prevCommand.MessageID,
		Info:         prevCommand.Info,
	})
	if err != nil {
		return errors.Wrap(err, "service.InsertState")
	}

	return nil
}

func (h *handler2) BotSlashCommandTypeCheckHandler(ctx context.Context, update tgbotapi.Update) error {
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

func (h *handler2) BotSlashCommandTypeChangeHandler(ctx context.Context, update tgbotapi.Update) error {
	deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	_, err := h.bot.Send(deleteMsg)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
		fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
	}

	state, err := h.service.SelectState(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.SelectState")
	}

	if state.Info != nil {
		deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, state.MessageID)
		_, err := h.bot.Send(deleteMsg)
		if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
			fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
		}
	}

	if state.Info != nil && state.SequenceName == enum.Change {
		if state.CommandName != enum.BotCommandNameTypeSaveStatus {
			// state.CommandName = SequenceToFirstCommand[state.SequenceName]
			if h.handlers[state.SequenceName][state.CommandName].Prev != nil {
				state.CommandName = h.handlers[state.SequenceName][state.CommandName].Prev.GetCommandName()
			}
		}

		prevCommand, err := h.handlers[state.SequenceName][enum.BotCommandNameTypeSaveStatus].Current.Question(ctx, update, state)
		if err != nil {
			return errors.Wrap(err, "handlers[state.SequenceName][enum.BotCommandNameTypeSaveStatus].Current.Question")
		}

		err = h.service.UpdateState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
			SequenceName: state.SequenceName,
			CommandName:  enum.BotCommandNameTypeSaveStatus,
			MessageID:    prevCommand.MessageID,
			Info:         prevCommand.Info,
		})
		if err != nil {
			return errors.Wrap(err, "service.InsertState")
		}

		return nil
	}

	err = h.service.DeleteState(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.DeleteState")
	}

	prevCommand, err := h.handlers[enum.Change][SequenceToFirstCommand[enum.Change]].Next.Question(ctx, update, dto.PrevCommandInfo{})
	if err != nil {
		return errors.Wrap(err, "handlers[enum.BotCommandNameTypeInputDate].Value.Question")
	}

	jsonData, err := json.Marshal(prevCommand.Info)
	if err != nil {
		return errors.Wrap(err, "json.Marshall")
	}

	err = h.service.InsertState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
		SequenceName: enum.Change,
		CommandName:  enum.BotCommandNameTypeTracking,
		MessageID:    prevCommand.MessageID,
		Info:         jsonData,
		KeyboardInfo: prevCommand.KeyboardInfo,
	})
	if err != nil {
		return errors.Wrap(err, "service.InsertState")
	}

	return nil
}

func (h *handler2) BotSlashCommandTypeDefaultHandler(ctx context.Context, update tgbotapi.Update) error {
	prevCommand, err := h.service.SelectState(ctx, update.Message.Chat.ID)
	if err != nil {
		errors.Wrap(err, "service.SelectState")
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Такой команды не существует")
		if _, msgErr := h.bot.Send(msg); err != nil {
			return errors.Wrap(msgErr, "bot.Send")
		}

		// return errors.New("Unknown command")
		return err
	}

	var data dto.WarehouseData
	err = json.Unmarshal(prevCommand.Info, &data)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	prevCommand, err = h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Current.Answer(ctx, update, prevCommand)
	if err != nil {
		var myerr *myError.MyError
		if errors.As(err, &myerr) {
			fmt.Println(err)

			deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
			_, err = h.bot.Send(deleteMsg)
			if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
				fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
			}

			deleteMsg = tgbotapi.NewDeleteMessage(update.Message.Chat.ID, prevCommand.MessageID)
			_, err = h.bot.Send(deleteMsg)
			if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
				fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
			}

			var whs []dto.WarehouseData
			if prevCommand.KeyboardInfo != nil {
				whs, err = Unmarshal[[]dto.WarehouseData](prevCommand.KeyboardInfo)
				if err != nil {
					return errors.Wrap(err, "Unmarshal")
				}
			}

			keyboardInfo := dto.KeyboardData{
				Warehouses: whs,
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, constmsg.MatchErrorType[myerr.GetErrorType()])
			msg, err = model.CommandToKeyboard[prevCommand.CommandName](msg, keyboardInfo)
			if err != nil {
				return errors.Wrap(err, "model.CommandToKeyboard")
			}

			message, err := h.bot.Send(msg)
			if err != nil {
				return errors.Wrap(err, "bot.Send")
			}

			prevCommand.MessageID = message.MessageID
			// prevCommands[update.Message.Chat.ID] = prevCommand
			err = h.service.UpdateState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
				SequenceName: prevCommand.SequenceName,
				CommandName:  prevCommand.CommandName,
				MessageID:    prevCommand.MessageID,
				Info:         prevCommand.Info,
			})
			if err != nil {
				return errors.Wrap(err, "service.UpdateState")
			}

			return nil
		}

		return errors.Wrap(err, "handlers[prevCommand.CommandName].Current.Answer")
	}

	deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	_, err = h.bot.Send(deleteMsg)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
		fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
	}

	deleteMsg = tgbotapi.NewDeleteMessage(update.Message.Chat.ID, prevCommand.MessageID)
	_, err = h.bot.Send(deleteMsg)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
		fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
	}

	prevCommand, err = h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Next.Question(ctx, update, prevCommand)
	if err != nil {
		return errors.Wrap(err, "handlers[prevCommand.CommandName].Next.Question")
	}

	err = h.service.UpdateState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
		SequenceName: prevCommand.SequenceName,
		CommandName:  h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Next.GetCommandName(),
		MessageID:    prevCommand.MessageID,
		Info:         prevCommand.Info,
	})
	if err != nil {
		return errors.Wrap(err, "service.UpdateState")
	}

	return nil
}

func (h *handler2) ButtonHandler(ctx context.Context, update tgbotapi.Update) error {
	prevCommand, err := h.service.SelectState(ctx, update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Такой команды не существует")
		if _, err := h.bot.Send(msg); err != nil {
			return errors.Wrap(err, "bot.Send")
		}

		return nil
	}

	var buttonData dto.ButtonData
	err = json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	if buttonData.Value == -1 {
		deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
		_, err = h.bot.Send(deleteMsg)
		if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
			fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
		}

		if h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Prev == nil {
			// delete(prevCommands, update.CallbackQuery.Message.Chat.ID)
			err = h.service.DeleteState(ctx, update.CallbackQuery.Message.Chat.ID)
			if err != nil {
				return errors.Wrap(err, "service.DeleteState")
			}

			return nil
		}

		prevCommand, err = h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Prev.Question(ctx, update, prevCommand)
		if err != nil {
			return errors.Wrap(err, "handlers[prevCommand.CommandName].Prev.Question")
		}

		err = h.service.UpdateState(ctx, update.CallbackQuery.Message.Chat.ID, dto.PrevCommandInfo{
			SequenceName: prevCommand.SequenceName,
			CommandName:  h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Prev.GetCommandName(),
			MessageID:    prevCommand.MessageID,
			Info:         prevCommand.Info,
			KeyboardInfo: prevCommand.KeyboardInfo,
		})
		if err != nil {
			return errors.Wrap(err, "service.UpdateState")
		}

		return nil
	}

	prevCommand, err = h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Current.Answer(ctx, update, prevCommand)
	if err != nil {
		return errors.Wrap(err, "handlers[prevCommand.CommandName].Current.Answer")
	}

	deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	_, err = h.bot.Send(deleteMsg)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool") {
		fmt.Printf("bot.Send(deleteMsg): %s\n", err.Error())
	}

	if h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Next == nil {
		switch prevCommand.SequenceName {
		case enum.Add:
			err = h.service.AddSequenceEndService(ctx, update.CallbackQuery.Message.Chat.ID, prevCommand.Info)
			if err != nil {
				return errors.Wrap(err, "service.AddSequenceEndService")
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Склад успешно добавлен")
			if _, err = h.bot.Send(msg); err != nil {
				return errors.Wrap(err, "bot.Send")
			}
		case enum.Change:
		default:
		}

		// delete(prevCommands, update.CallbackQuery.Message.Chat.ID)
		err = h.service.DeleteState(ctx, update.CallbackQuery.Message.Chat.ID)
		if err != nil {
			return errors.Wrap(err, "service.DeleteState")
		}

		return nil
	}

	prevCommand, err = h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Next.Question(ctx, update, prevCommand)
	if err != nil {
		return errors.Wrap(err, "handlers[prevCommand.CommandName].Next.Question")
	}

	err = h.service.UpdateState(ctx, update.CallbackQuery.Message.Chat.ID, dto.PrevCommandInfo{
		SequenceName: prevCommand.SequenceName,
		CommandName:  h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Next.GetCommandName(),
		MessageID:    prevCommand.MessageID,
		Info:         prevCommand.Info,
	})
	if err != nil {
		return errors.Wrap(err, "service.UpdateState")
	}

	return nil
}
