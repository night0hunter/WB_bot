package handler

import (
	"context"
	"encoding/json"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	keyboard "wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type BookProtectionHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequence
}

func (h *BookProtectionHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	if update.CallbackQuery == nil {
		return tmpData, nil
	}

	var buttonData dto.ButtonData

	err := json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "json.Unmarshal")
	}

	text, err := SequenceController(tmpData, h.GetCommandName())
	if err != nil {
		return tmpData, errors.Wrap(err, "SequenceController")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)

	msg, err = keyboard.DrawBookProtectKeyboard(msg, dto.KeyboardData{})
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, " keyboard.DrawBookProtectKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID

	return tmpData, nil
}

func (h *BookProtectionHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	data, err := Unmarshal[dto.WarehouseData](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "Unmarshal")
	}

	if update.Message != nil {
		coeff, err := h.service.BotAnswerInputCoeffLimitService(ctx, update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			return tmpData, errors.Wrap(err, "service.BotAnswerInputCoeffLimitService")
		}

		data.CoeffLimit = coeff
	}

	if update.CallbackQuery != nil {
		var buttonData dto.ButtonData

		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
		if err != nil {
			return tmpData, errors.Wrap(err, "json.Unmarshal")
		}

		data.CoeffLimit = buttonData.Value
	}

	json, err := Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *BookProtectionHandler) GetCommandName() enum.CommandSequence {
	return h.commandName
}
