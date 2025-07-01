package bookHandler

import (
	"context"
	"encoding/json"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
	keyboard "wb_bot/internal/handler/keyboard"
	"wb_bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type SaveStatusHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequence
}

func (h *SaveStatusHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	var msg tgbotapi.MessageConfig
	var err error

	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже начинали заполнение, выберите действие")
	}

	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы уже начинали заполнение, выберите действие")
	}

	msg, err = keyboard.DrawSaveStatusKeyboard(msg, dto.KeyboardData{})
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "keyboard.DrawCancelKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID

	stateJson, err := json.Marshal(tmpData)
	if err != nil {
		return tmpData, errors.Wrap(err, "json.Marshal")
	}

	tmpData.Info = stateJson

	return tmpData, nil
}

func (h *SaveStatusHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	if update.CallbackQuery == nil && update.Message != nil {
		return tmpData, &myError.MyError{
			ErrType: myError.SaveStatusChoiceError,
			Message: "saveStatusChoice - user input error",
		}
	}

	var buttonData dto.ButtonData

	err := json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
	if err != nil {
		return tmpData, errors.Wrap(err, "json.Unmarshal")
	}

	switch buttonData.Value {
	case 1: // Continue sequence
		var data dto.PrevCommandInfo
		err := json.Unmarshal(tmpData.Info, &data)
		if err != nil {
			return dto.PrevCommandInfo{}, errors.Wrap(err, "json.Unmarshal")
		}

		return data, nil
	case 2: // Start sequence over
		tmpData.CommandName = model.SequenceToFirstCommand[tmpData.SequenceName]
	default:
	}

	return tmpData, nil
}

func (h *SaveStatusHandler) GetCommandName() enum.CommandSequence {
	return h.commandName
}
