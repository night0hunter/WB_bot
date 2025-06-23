package keyboard

import (
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func DrawSaveStatusKeyboard(msg tgbotapi.MessageConfig, data dto.KeyboardData) (tgbotapi.MessageConfig, error) {
	tmpMarkup, err := GenerateKeyboard([]dto.Button{
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSaveStatus,
				Value: 1,
			},
			Text: "Продолжить",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSaveStatus,
				Value: 2,
			},
			Text: "Начать заново",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSaveStatus,
				Value: -1,
			},
			Text: "Отмена",
		},
	}...)
	if err != nil {
		return msg, errors.Wrap(err, "GenerateKeyboard")
	}

	msg.ReplyMarkup = tmpMarkup

	return msg, nil
}
