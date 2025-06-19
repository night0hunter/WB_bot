package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
	keyboard "wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type WarehouseHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequences
}

func (h *WarehouseHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData prevCommandInfo) (prevCommandInfo, error) {
	var msg tgbotapi.MessageConfig

	data, err := Unmarshal[dto.WarehouseData](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "Unmarshal")
	}

	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
			data.FromDate.Format(dto.TimeFormat),
			data.ToDate.Format(dto.TimeFormat),
			"",
			"",
			"",
			BotCommands[enum.BotCommandNameTypeInputWarehouse],
		))
	}

	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
			"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
			data.FromDate.Format(dto.TimeFormat),
			data.ToDate.Format(dto.TimeFormat),
			"",
			"",
			"",
			BotCommands[enum.BotCommandNameTypeInputWarehouse],
		))
	}

	msg, err = keyboard.DrawWarehouseKeyboard(msg, dto.KeyboardData{})
	if err != nil {
		return tmpData, errors.Wrap(err, "keyboard.DrawWarehouseKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return tmpData, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID

	return tmpData, nil
}

func (h *WarehouseHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData prevCommandInfo) (prevCommandInfo, error) {
	if update.CallbackQuery == nil && update.Message == nil {
		return tmpData, nil
	}

	data, err := Unmarshal[dto.WarehouseData](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "Unmarshal")
	}

	if update.CallbackQuery == nil && update.Message != nil {
		data.Warehouse = 0
		return tmpData, &myError.MyError{
			ErrType: myError.WarehouseInputError,
			Message: "warehouse - user input error",
		}

		// return tmpData, nil
	}

	var buttonData dto.ButtonData

	err = json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
	if err != nil {
		return tmpData, errors.Wrap(err, "json.Unmarshal")
	}

	data.Warehouse = buttonData.Value

	json, err := Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *WarehouseHandler) GetCommandName() enum.CommandSequences {
	return h.commandName
}
