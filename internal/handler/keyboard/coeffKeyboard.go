package keyboard

import (
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DrawCoeffKeyboard(msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	tmpMarkup, err := GenerateKeyboard([]dto.Button{
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeCoeffLimit,
				Value: 0,
			},
			Text: "Бесплатно",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeCoeffLimit,
				Value: 1,
			},
			Text: "1x",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeCoeffLimit,
				Value: 2,
			},
			Text: "2x",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeCoeffLimit,
				Value: 3,
			},
			Text: "3x",
		},
	}...)
	if err != nil {
		return msg, err
	}

	msg.ReplyMarkup = tmpMarkup

	return msg, nil
}
