package bookHandler

import (
	"context"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
	keyboard "wb_bot/internal/handler/keyboard"
	"wb_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type DraftIdHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequence
}

func (h *DraftIdHandler) Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	var msg tgbotapi.MessageConfig
	var data dto.BookingData
	var err error

	if tmpData.Info != nil {
		data, err = utils.Unmarshal[dto.BookingData](tmpData.Info)
		if err != nil {
			return dto.PrevCommandInfo{}, errors.Wrap(err, "utils.Unmarshal")
		}
	}

	text, err := CraftMessage(tmpData, h.GetCommandName())
	if err != nil {
		return tmpData, errors.Wrap(err, "SequenceController")
	}

	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
	}

	if update.CallbackQuery != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	}

	msg, err = keyboard.DrawBackKeyboard(msg, dto.KeyboardData{})
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

func (h *DraftIdHandler) Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error) {
	data, err := utils.Unmarshal[dto.BookingData](tmpData.Info)
	if err != nil {
		return tmpData, errors.Wrap(err, "utils.Unmarshal")
	}

	if update.Message == nil {
		data.DraftID = uuid.UUID{}
		json, err := utils.Marshal(data)
		if err != nil {
			return tmpData, errors.Wrap(err, "utils.Marshal")
		}

		tmpData.Info = json

		return tmpData, nil
	}

	id, err := uuid.Parse(update.Message.Text)
	if err != nil {
		return tmpData, &myError.MyError{
			ErrType: myError.BookingIdError,
			Message: "uuid.Parse: bookingID - user input error",
		}
	}

	data.DraftID = id

	json, err := utils.Marshal(data)
	if err != nil {
		return tmpData, errors.Wrap(err, "utils.Marshal")
	}

	tmpData.Info = json

	return tmpData, nil
}

func (h *DraftIdHandler) GetCommandName() enum.CommandSequence {
	return h.commandName
}
