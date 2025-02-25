package keyboard

import (
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DrawSupplyKeyboard(msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	tmpMarkup, err := GenerateKeyboard([]dto.Button{
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSupplyType,
				Value: 2,
			},
			Text: "Короб",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSupplyType,
				Value: 5,
			},
			Text: "Монопалеты",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSupplyType,
				Value: 6,
			},
			Text: "Суперсейф",
		},
	}...)
	if err != nil {
		return msg, err
	}

	msg.ReplyMarkup = tmpMarkup

	return msg, nil
}
