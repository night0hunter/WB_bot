package keyboard

import (
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func DrawActionChoiceKeyboard(msg tgbotapi.MessageConfig, data dto.KeyboardData) (tgbotapi.MessageConfig, error) {
	tmpMarkup, err := GenerateKeyboard([]dto.Button{
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeActionChoice,
				Value: 1,
			},
			Text: "Изменить статус отслеживания",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeActionChoice,
				Value: 2,
			},
			Text: "Удалить отслеживание",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeActionChoice,
				Value: -1,
			},
			Text: "Назад",
		},
	}...)
	if err != nil {
		return msg, errors.Wrap(err, "GenerateKeyboard")
	}

	msg.ReplyMarkup = tmpMarkup

	return msg, nil
}
