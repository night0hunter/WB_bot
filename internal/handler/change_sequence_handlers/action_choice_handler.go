package changeHandler

import (
	"context"
	"encoding/json"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
	keyboard "wb_bot/internal/handler/keyboard"
	"wb_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (h *ActionChoiceHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	var err error
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите действие")
	msg, err = keyboard.DrawActionChoiceKeyboard(msg, dto.KeyboardData{})
	if err != nil {
		return tmpData, errors.Wrap(err, "keyboard.DrawActionChoiceKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return tmpData, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID

	return tmpData, nil
}

func (h *ActionChoiceHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	if update.Message != nil {
		return tmpData, &myError.MyError{
			ErrType: myError.ActionChoiceError,
			Message: "actionChoice - user input error",
		}
	}

	data, err := utils.Unmarshal[dto.ChangeStatusInfo](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "Unmarshal")
	}

	var buttonData dto.ButtonData
	err = json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
	if err != nil {
		return tmpData, errors.Wrap(err, "json.Unmarshal")
	}

	switch buttonData.Value {
	case 1:
		err = h.service.ChangeStatusService(ctx, update.CallbackQuery.Message.Chat.ID, data.TrackingID)
		if err != nil {
			return tmpData, errors.Wrap(err, "service.ChangeStatusService")
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Статус отслеживания успешно изменён")
		if _, err = h.bot.Send(msg); err != nil {
			return tmpData, errors.Wrap(err, "bot.Send")
		}
	case 2:
		err = h.service.DeleteTrackingService(ctx, update.CallbackQuery.Message.Chat.ID, data.TrackingID)
		if err != nil {
			return tmpData, errors.Wrap(err, "service.DeleteTrackingService")
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Отслеживание успешно удалено")
		if _, err = h.bot.Send(msg); err != nil {
			return tmpData, errors.Wrap(err, "bot.Send")
		}
	default:
		return tmpData, errors.Wrap(err, "undefined choice")
	}

	return tmpData, nil
}

func (h *ActionChoiceHandler) GetCommandName() enum.CommandSequence {
	return h.commandName
}
