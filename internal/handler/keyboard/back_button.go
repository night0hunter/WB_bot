package keyboard

import (
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DrawBackKeyboard(msg tgbotapi.MessageConfig, data dto.KeyboardData) (tgbotapi.MessageConfig, error) {
	tmpMarkup, err := GenerateKeyboard([]dto.Button{
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
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
