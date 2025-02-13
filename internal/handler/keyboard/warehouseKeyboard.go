package keyboard

import (
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DrawWarehouseKeyboard(msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	tmpMarkup, err := GenerateKeyboard([]dto.Button{
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 507, // wareID
			},
			Text: "Коледино", // wareName
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 316646,
			},
			Text: "Шушары СГТ",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 301229,
			},
			Text: "Подольск 4",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 120762,
			},
			Text: "Электросталь",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 206348,
			},
			Text: "Тула",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 130744,
			},
			Text: "Краснодар",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 208277,
			},
			Text: "Невинномысск",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 117986,
			},
			Text: "Казань",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 117501,
			},
			Text: "Подольск",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 1733,
			},
			Text: "Екатеринбург - Испытателей 14г",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 218644,
			},
			Text: "СЦ Хабаровск", // add "Хабаровск"
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 218644, // find id
			},
			Text: "Санкт-Петербург Уткина Заводь",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 206236,
			},
			Text: "Белые Столбы",
		},
		{
			Data: dto.ButtonData{
				Type:  enum.ButtonTypeWarehouse,
				Value: 686,
			},
			Text: "Новосибирск",
		},
	}...)
	if err != nil {
		return msg, err
	}

	msg.ReplyMarkup = tmpMarkup

	return msg, nil

}
