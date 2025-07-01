package bookHandler

import (
	"context"
	"encoding/json"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	keyboard "wb_bot/internal/handler/keyboard"
	"wb_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type CoeffLimitHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequence
}

func (h *CoeffLimitHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	if update.CallbackQuery == nil {
		return tmpData, nil
	}

	text, err := CraftMessage(tmpData, h.GetCommandName())
	if err != nil {
		return tmpData, errors.Wrap(err, "SequenceController")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)

	msg, err = keyboard.DrawCoeffKeyboard(msg, dto.KeyboardData{})
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "keyboard.DrawCoeffKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID

	return tmpData, nil
}

func (h *CoeffLimitHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	data, err := utils.Unmarshal[dto.BookingData](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "utils.Unmarshal")
	}

	if update.Message != nil {
		coeff, err := h.service.BotAnswerInputCoeffLimitService(ctx, update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			return tmpData, errors.Wrap(err, "service.BotAnswerInputCoeffLimitService")
		}

		data.CoeffLimit = &coeff
	}

	if update.CallbackQuery != nil {
		var buttonData dto.ButtonData

		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
		if err != nil {
			return tmpData, errors.Wrap(err, "json.Unmarshal")
		}

		data.CoeffLimit = &buttonData.Value
	}

	json, err := utils.Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "utils.Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *CoeffLimitHandler) GetCommandName() enum.CommandSequence {
	return h.commandName
}
