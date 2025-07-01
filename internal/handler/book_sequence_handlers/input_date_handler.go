package bookHandler

import (
	"context"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/handler/keyboard"
	"wb_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type InputDateHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequence
}

func (h *InputDateHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	var msg tgbotapi.MessageConfig
	var data dto.WarehouseData
	var err error

	if tmpData.Info != nil {
		data, err = utils.Unmarshal[dto.WarehouseData](tmpData.Info)
		if err != nil {
			return dto.PrevCommandInfo{}, errors.Wrap(err, "utils.Unmarshal")
		}
	}

	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, constmsg.BotCommands[enum.BotCommandNameTypeInputDate])
		data.ChatID = update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, constmsg.BotCommands[enum.BotCommandNameTypeInputDate])
		data.ChatID = update.CallbackQuery.Message.Chat.ID
	}

	msg, err = keyboard.DrawCancelKeyboard(msg, dto.KeyboardData{})
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "keyboard.DrawCancelKeyboard")
	}

	message, err := h.bot.Send(msg)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "bot.Send")
	}

	tmpData.MessageID = message.MessageID
	json, err := utils.Marshal(data)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "utils.Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *InputDateHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	data, err := utils.Unmarshal[dto.BookingData](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "utils.Unmarshal")
	}

	if update.Message == nil {
		return tmpData, nil
	}

	timeRange, err := h.service.BotAnswerInputDateService(ctx, update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		return tmpData, errors.Wrap(err, "service.BotAnswerInputDateService")
	}

	data.FromDate = timeRange.DateFrom
	data.ToDate = timeRange.DateTo

	json, err := utils.Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "utils.Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *InputDateHandler) GetCommandName() enum.CommandSequence {
	return h.commandName
}
