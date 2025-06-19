package handler

import (
	"context"
	"encoding/json"
	"fmt"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
	keyboard "wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type TrackingChoiceHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequences
}

func (h *TrackingChoiceHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData prevCommandInfo) (prevCommandInfo, error) {
	var warehouses []dto.WarehouseData
	var err error

	if update.Message != nil {
		warehouses, err = h.service.BotSlashCommandTypeChange(ctx, update.Message.Chat.ID)
		if err != nil {
			return tmpData, errors.Wrap(err, "service.BotSlashCommandTypeChange")
		}
	}

	if update.CallbackQuery != nil {
		warehouses, err = h.service.BotSlashCommandTypeChange(ctx, update.CallbackQuery.Message.Chat.ID)
		if err != nil {
			return tmpData, errors.Wrap(err, "service.BotSlashCommandTypeChange")
		}
	}

	if len(warehouses) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("На данный момент вы не отслеживаете ни одного склада, чтобы добавить, используйте %s", constmsg.BotSlashCommands[enum.BotSlashCommandTypeAdd]))
		if _, err := h.bot.Send(msg); err != nil {
			return tmpData, errors.Wrap(err, "bot.Send")
		}

		return tmpData, nil
	}

	var msg tgbotapi.MessageConfig
	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите отслеживание из списка ниже, чтобы изменить его статус/удалить")
	}

	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите отслеживание из списка ниже, чтобы изменить его статус/удалить")
	}

	data := dto.KeyboardData{
		Warehouses: warehouses,
	}
	msg, err = keyboard.DrawTrackingsKeyboard(msg, data)
	if err != nil {
		return tmpData, errors.Wrap(err, "keyboard.DrawTrackingsKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return tmpData, errors.Wrap(err, "bot.Send")
	}

	j, err := json.Marshal(warehouses)
	if err != nil {
		return tmpData, errors.Wrap(err, "json.Marshal")
	}

	tmpData.MessageID = message.MessageID
	tmpData.KeyboardInfo = j

	return tmpData, nil
}

func (h *TrackingChoiceHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData prevCommandInfo) (prevCommandInfo, error) {
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

	var data dto.ChangeStatusInfo

	if tmpData.Info != nil {
		data, err = Unmarshal[dto.ChangeStatusInfo](tmpData.Info)
		if err != nil {
			return tmpData, errors.Wrap(err, "Unmarshal")
		}
	}

	data.TrackingID = int64(buttonData.Value)

	json, err := Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *TrackingChoiceHandler) GetCommandName() enum.CommandSequences {
	return h.commandName
}
