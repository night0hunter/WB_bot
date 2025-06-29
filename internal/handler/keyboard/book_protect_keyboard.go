package keyboard

import (
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DrawBookProtectKeyboard(msg tgbotapi.MessageConfig, data dto.KeyboardData) (tgbotapi.MessageConfig, error) {
	tmpMarkup, err := GenerateKeyboard([]dto.Button{
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeBookProtection,
				Value: -2, // shit for constmsg.Coefficients map
			},
			Text: "0",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeBookProtection,
				Value: 1,
			},
			Text: "1",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeBookProtection,
				Value: 2,
			},
			Text: "2",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeBookProtection,
				Value: 3,
			},
			Text: "3",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeBookProtection,
				Value: -1,
			},
			Text: "Назад",
		},
	}...)
	if err != nil {
		return msg, err
	}

	msg.ReplyMarkup = tmpMarkup

	return msg, nil
}
