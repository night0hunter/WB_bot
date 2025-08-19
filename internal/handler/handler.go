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
	addHandler "wb_bot/internal/handler/add_sequence_handlers"
	bookHandler "wb_bot/internal/handler/book_sequence_handlers"
	changeBookingHandler "wb_bot/internal/handler/change_booking_sequence_handlers"
	changeTrackingHandler "wb_bot/internal/handler/change_tracking_sequence_handlers"
	"wb_bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/pkg/errors"
)

type Service interface {
	DeleteTrackingService(ctx context.Context, trackingID int64) error
	DeleteBookingService(ctx context.Context, trackingID int64) error
	ChangeStatusService(ctx context.Context, trackingID int64) error
	ChangeBookingStatusService(ctx context.Context, trackingID int64) error
	BotSlashCommandTypeHelpService(ctx context.Context, chatID int64) string
	BotSlashCommandTypeCheckTrackingsService(ctx context.Context, chatID int64) ([]string, error)
	BotSlashCommandTypeCheckBookingsService(ctx context.Context, chatID int64) ([]string, error)
	BotSlashCommandTypeChange(ctx context.Context, chatID int64) ([]dto.WarehouseData, error)
	BotSlashCommandTypeChangeBooking(ctx context.Context, chatID int64) ([]dto.BookingData, error)
	BotAnswerInputDateService(ctx context.Context, chatID int64, date string) (dto.TrackingDate, error)
	BotAnswerInputCoeffLimitService(ctx context.Context, chatID int64, coeffLimit string) (int, error)

	AddSequenceEndService(ctx context.Context, chatID int64, data []byte) error
	BookSequenceEndService(ctx context.Context, chatID int64, data []byte) error

	SelectState(ctx context.Context, chatID int64) (dto.PrevCommandInfo, error)
	InsertState(ctx context.Context, chatID int64, prevCommand dto.PrevCommandInfo) error
	UpdateState(ctx context.Context, chatID int64, prevCommand dto.PrevCommandInfo) error
	DeleteState(ctx context.Context, chatID int64) error

	BookDraftCron(ctx context.Context) (int64, error)
	GetTrackings(ctx context.Context) ([]dto.MergedResp, error)
	KeepSendingTime(ctx context.Context, tracking dto.MergedResp) error
}

// var prevCommands = map[int64]dto.PrevCommandInfo{}

type handler struct {
	bot      *tgbotapi.BotAPI
	service  Service
	handlers map[enum.Sequences]map[enum.CommandSequence]struct {
		Prev    model.HandlerStruct
		Current model.HandlerStruct
		Next    model.HandlerStruct
	}
}

func New(bot *tgbotapi.BotAPI, svc Service) *handler {
	handlers := map[enum.Sequences]map[enum.CommandSequence]struct {
		Prev    model.HandlerStruct
		Current model.HandlerStruct
		Next    model.HandlerStruct
	}{}

	handlers[enum.Add] = addHandler.New(bot, svc)
	handlers[enum.ChangeTracking] = changeTrackingHandler.New(bot, svc)
	handlers[enum.Booking] = bookHandler.New(bot, svc)
	handlers[enum.ChangeBooking] = changeBookingHandler.New(bot, svc)

	return &handler{bot: bot, service: svc, handlers: handlers}
}

func (h *handler) Run(ctx context.Context) error {
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
				fmt.Println(errors.Wrap(err, "h.messageHandler"))
			}

			continue
		}

		if update.Message != nil {
			err := h.messageHandler(ctx, update)
			if err != nil {
				// return errors.Wrap(err, "h.messageHandler")
				fmt.Println(errors.Wrap(err, "h.messageHandler"))
			}
		}

	}

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
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeChangeTracking]:
		err := h.BotSlashCommandTypeChangeTrackingHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeChangeTrackingHandler")
		}
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeCheckTrackings]:
		err := h.BotSlashCommandTypeCheckTrackingsHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeCheckHandler")
		}
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeBook]:
		err := h.BotSlashCommandTypeBookHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeBookHandler")
		}
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeChangeBooking]:
		err := h.BotSlashCommandTypeChangeBookingHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeChangeBookingHandler")
		}
	case constmsg.BotSlashCommands[enum.BotSlashCommandTypeCheckBookings]:
		err := h.BotSlashCommandTypeCheckBookingsHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeChangeBookingHandler")
		}
	default:
		err := h.BotSlashCommandTypeDefaultHandler(ctx, update)
		if err != nil {
			return errors.Wrap(err, "BotSlashCommandTypeDefaultHandler")
		}
	}

	return nil
}

func (h *handler) BotSlashCommandTypeChangeBookingHandler(ctx context.Context, update tgbotapi.Update) error {
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

	if state.Info != nil && state.SequenceName == enum.ChangeBooking {
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

	prevCommand, err := h.handlers[enum.ChangeBooking][model.SequenceToFirstCommand[enum.ChangeBooking]].Next.Question(ctx, update, dto.PrevCommandInfo{})
	if err != nil {
		return errors.Wrap(err, "handlers[enum.BotCommandNameTypeInputDate].Value.Question")
	}

	jsonData, err := json.Marshal(prevCommand.Info)
	if err != nil {
		return errors.Wrap(err, "json.Marshall")
	}

	err = h.service.InsertState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
		SequenceName: enum.ChangeBooking,
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

func (h *handler) BotSlashCommandTypeHelpHandler(ctx context.Context, update tgbotapi.Update) error {
	text := h.service.BotSlashCommandTypeHelpService(ctx, update.Message.Chat.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "bot.Send")
	}

	return nil
}

func (h *handler) BotSlashCommandTypeAddHandler(ctx context.Context, update tgbotapi.Update) error {
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
			if h.handlers[state.SequenceName][state.CommandName].Prev != nil {
				state.CommandName = h.handlers[state.SequenceName][state.CommandName].Prev.GetCommandName()
			} else {
				state.CommandName = model.SequenceToFirstCommand[state.SequenceName]
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

	prevCommand, err := h.handlers[enum.Add][model.SequenceToFirstCommand[enum.Add]].Next.Question(ctx, update, dto.PrevCommandInfo{})
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

func (h *handler) BotSlashCommandTypeCheckTrackingsHandler(ctx context.Context, update tgbotapi.Update) error {
	whs, err := h.service.BotSlashCommandTypeCheckTrackingsService(ctx, update.Message.Chat.ID)
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

func (h *handler) BotSlashCommandTypeCheckBookingsHandler(ctx context.Context, update tgbotapi.Update) error {
	whs, err := h.service.BotSlashCommandTypeCheckBookingsService(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.BotSlashCommandTypeCheck")
	}

	if whs == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"На данный момент У вас нет автоброни, чтобы добавить, используйте %s",
			constmsg.BotSlashCommands[enum.BotSlashCommandTypeBook]))
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

func (h *handler) BotSlashCommandTypeChangeTrackingHandler(ctx context.Context, update tgbotapi.Update) error {
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

	if state.Info != nil && state.SequenceName == enum.ChangeTracking {
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

	prevCommand, err := h.handlers[enum.ChangeTracking][model.SequenceToFirstCommand[enum.ChangeTracking]].Next.Question(ctx, update, dto.PrevCommandInfo{})
	if err != nil {
		return errors.Wrap(err, "handlers[enum.BotCommandNameTypeInputDate].Value.Question")
	}

	jsonData, err := json.Marshal(prevCommand.Info)
	if err != nil {
		return errors.Wrap(err, "json.Marshall")
	}

	err = h.service.InsertState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
		SequenceName: enum.ChangeTracking,
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

func (h *handler) BotSlashCommandTypeBookHandler(ctx context.Context, update tgbotapi.Update) error {
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

	if state.Info != nil && state.SequenceName == enum.Booking {
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

	prevCommand, err := h.handlers[enum.Booking][model.SequenceToFirstCommand[enum.Booking]].Next.Question(ctx, update, dto.PrevCommandInfo{})
	if err != nil {
		return errors.Wrap(err, "handlers[enum.Add][SequenceToFirstCommand[enum.Add]].Next.Question")
	}

	err = h.service.InsertState(ctx, update.Message.Chat.ID, dto.PrevCommandInfo{
		SequenceName: enum.Booking,
		CommandName:  enum.BotCommandNameTypeInputDate,
		MessageID:    prevCommand.MessageID,
		Info:         prevCommand.Info,
	})
	if err != nil {
		return errors.Wrap(err, "service.InsertState")
	}

	return nil
}

func (h *handler) BotSlashCommandTypeDefaultHandler(ctx context.Context, update tgbotapi.Update) error {
	prevCommand, err := h.service.SelectState(ctx, update.Message.Chat.ID)
	if err != nil {
		return errors.Wrap(err, "service.SelectState")
	}

	if prevCommand.Info == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Такой команды не существует")
		if _, msgErr := h.bot.Send(msg); err != nil {
			return errors.Wrap(msgErr, "bot.Send")
		}

		return nil
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

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, constmsg.MatchErrorType[myerr.GetErrorType()])
			msg, err = model.CommandToKeyboard[prevCommand.CommandName](msg, prevCommand.KeyboardInfo)
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

func (h *handler) ButtonHandler(ctx context.Context, update tgbotapi.Update) error {
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

		copy := update
		copy.CallbackQuery.Data = "{}"
		prevCommand, err = h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Prev.Answer(ctx, copy, prevCommand)
		if err != nil {
			return errors.Wrap(err, "handlers[prevCommand.SequenceName][prevCommand.CommandName].Prev.Answer")
		}

		prevCommand, err = h.handlers[prevCommand.SequenceName][prevCommand.CommandName].Prev.Question(ctx, update, prevCommand)
		if err != nil {
			return errors.Wrap(err, "handlers[prevCommand.SequenceName][prevCommand.CommandName].Prev.Question")
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

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Отслеживание успешно добавлено")
			if _, err = h.bot.Send(msg); err != nil {
				return errors.Wrap(err, "bot.Send")
			}
		case enum.Booking:
			err = h.service.BookSequenceEndService(ctx, update.CallbackQuery.Message.Chat.ID, prevCommand.Info)
			if err != nil {
				return errors.Wrap(err, "service.BookSequenceEndService")
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Автобронирование успешно добавлено")
			if _, err = h.bot.Send(msg); err != nil {
				return errors.Wrap(err, "bot.Send")
			}
		case enum.ChangeTracking:
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
