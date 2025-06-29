package handler

import (
	"context"
	"encoding/json"
	"strconv"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
	keyboard "wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type SupplyTypeHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequence
}

func (h *SupplyTypeHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	var msg tgbotapi.MessageConfig

	text, err := SequenceController(tmpData, h.GetCommandName())
	if err != nil {
		return tmpData, errors.Wrap(err, "SequenceController")
	}

	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
	}

	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	}

	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, BotCommands[enum.BotCommandNameTypeInputWarehouse])
	msg, err = keyboard.DrawSupplyKeyboard(msg, dto.KeyboardData{})
	if err != nil {
		return tmpData, errors.Wrap(err, "keyboard.DrawSupplyKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return tmpData, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID

	return tmpData, nil
}

func (h *SupplyTypeHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	data, err := Unmarshal[dto.WarehouseData](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "Unmarshal")
	}

	if update.Message != nil {
		data.SupplyType = ""

		return tmpData, &myError.MyError{
			ErrType: myError.SupplyTypeError,
			Message: "supplyType - user input error",
		}
	}

	var buttonData dto.ButtonData
	err = json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
	if err != nil {
		return tmpData, errors.Wrap(err, "json.Unmarshal")
	}

	data.SupplyType = constmsg.SupplyTypes[strconv.Itoa(buttonData.Value)]

	json, err := Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *SupplyTypeHandler) GetCommandName() enum.CommandSequence {
	return h.commandName
}
