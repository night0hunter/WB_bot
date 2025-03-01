package keyboard

import (
	constmsg "wb_bot/internal/const_message"
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
			Text: constmsg.SupplyTypes[enum.Box],
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSupplyType,
				Value: 5,
			},
			Text: constmsg.SupplyTypes[enum.Monopallet],
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSupplyType,
				Value: 6,
			},
			Text: constmsg.SupplyTypes[enum.SuperSafe],
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeSupplyType,
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
