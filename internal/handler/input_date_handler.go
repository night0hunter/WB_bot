package handler

import (
	"context"
	"time"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type InputDateHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequences
}

func (h *InputDateHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData prevCommandInfo) (prevCommandInfo, error) {
	var msg tgbotapi.MessageConfig
	var data dto.WarehouseData
	var err error

	if tmpData.Info != nil {
		data, err = Unmarshal[dto.WarehouseData](tmpData.Info)
		if err != nil {
			return prevCommandInfo{}, errors.Wrap(err, "Unmarshal")
		}
	}

	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputDate])
		data.ChatID = update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputDate])
		data.ChatID = update.CallbackQuery.Message.Chat.ID
	}

	msg, err = keyboard.DrawCancelKeyboard(msg, dto.KeyboardData{})
	if err != nil {
		return prevCommandInfo{}, errors.Wrap(err, "keyboard.DrawCancelKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return prevCommandInfo{}, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID
	json, err := Marshal(data)
	if err != nil {
		return prevCommandInfo{}, errors.Wrap(err, "Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *InputDateHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData prevCommandInfo) (prevCommandInfo, error) {
	data, err := Unmarshal[dto.WarehouseData](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "Unmarshal")
	}

	if update.Message == nil {
		data.FromDate = time.Time{}
		data.ToDate = time.Time{}

		return tmpData, nil
	}

	timeRange, err := h.service.BotAnswerInputDateService(ctx, update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		return tmpData, errors.Wrap(err, "service.BotAnswerInputDateService")
	}

	data.FromDate = timeRange.DateFrom
	data.ToDate = timeRange.DateTo

	json, err := Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *InputDateHandler) GetCommandName() enum.CommandSequences {
	return h.commandName
}
