package handler

import (
	"context"
	"encoding/json"
	"fmt"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	keyboard "wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type CoeffLimitHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequences
}

func (h *CoeffLimitHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	if update.CallbackQuery == nil {
		return tmpData, nil
	}

	data, err := Unmarshal[dto.WarehouseData](tmpData.Info)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "Unmarshal")
	}

	var buttonData dto.ButtonData

	err = json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "json.Unmarshal")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
		"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
		data.FromDate.Format(dto.TimeFormat),
		data.ToDate.Format(dto.TimeFormat),
		constmsg.WarehouseNames[data.Warehouse],
		"",
		"",
		BotCommands[enum.BotCommandNameTypeInputCoeffLimit],
	))

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

func (h *CoeffLimitHandler) GetCommandName() enum.CommandSequences {
	return h.commandName
}
