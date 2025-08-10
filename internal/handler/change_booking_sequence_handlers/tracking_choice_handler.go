package changeBookingHandler

import (
	"context"
	"encoding/json"
	"fmt"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
	keyboard "wb_bot/internal/handler/keyboard"
	"wb_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type BookingChoiceHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequence
}

func (h *BookingChoiceHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	var warehouses []dto.BookingData
	var err error

	if update.Message != nil {
		warehouses, err = h.service.BotSlashCommandTypeChangeBooking(ctx, update.Message.Chat.ID)
		if err != nil {
			return tmpData, errors.Wrap(err, "service.BotSlashCommandTypeChangeBooking")
		}
	}

	if update.CallbackQuery != nil {
		warehouses, err = h.service.BotSlashCommandTypeChangeBooking(ctx, update.CallbackQuery.Message.Chat.ID)
		if err != nil {
			return tmpData, errors.Wrap(err, "service.BotSlashCommandTypeChangeBooking")
		}
	}

	if len(warehouses) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("На данный момент У вас нет автоброни, чтобы добавить, используйте %s", constmsg.BotSlashCommands[enum.BotSlashCommandTypeBook]))
		if _, err := h.bot.Send(msg); err != nil {
			return tmpData, errors.Wrap(err, "bot.Send")
		}

		return tmpData, nil
	}

	var msg tgbotapi.MessageConfig
	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeTracking])
	}

	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите бронирование из списка ниже, чтобы изменить его статус/удалить")
	}

	j, err := utils.Marshal(warehouses)
	if err != nil {
		return tmpData, errors.Wrap(err, "utils.Marshal")
	}
	msg, err = keyboard.DrawBookingsKeyboard(msg, j)
	if err != nil {
		return tmpData, errors.Wrap(err, "keyboard.DrawBookingsKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return tmpData, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID
	tmpData.KeyboardInfo = j

	return tmpData, nil
}

func (h *BookingChoiceHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	if update.CallbackQuery == nil && update.Message == nil {
		return tmpData, nil
	}

	if update.Message != nil {
		return tmpData, &myError.MyError{
			ErrType: myError.TrackingChoiceError,
			Message: "trackingChoice - user input error",
		}
	}

	var buttonData dto.ButtonData

	err := json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
	if err != nil {
		return tmpData, errors.Wrap(err, "json.Unmarshal")
	}

	var data dto.ChangeBookingStatusInfo

	if tmpData.Info != nil {
		data, err = utils.Unmarshal[dto.ChangeBookingStatusInfo](tmpData.Info)
		if err != nil {
			return tmpData, errors.Wrap(err, "Unmarshal")
		}
	}

	data.BookingID = int64(buttonData.Value)

	json, err := utils.Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *BookingChoiceHandler) GetCommandName() enum.CommandSequence {
	return h.commandName
}
