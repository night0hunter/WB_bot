package keyboard

import (
	"fmt"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DrawTrackingsKeyboard(msg tgbotapi.MessageConfig, warehouses []dto.WarehouseData) tgbotapi.MessageConfig {
	var buttons []dto.Button
	var button dto.Button

	for _, wh := range warehouses {
		button.Text = constmsg.WarehouseNames[int(wh.Warehouse)] + " " + wh.FromDate.Format(dto.TimeFormat) + "-" + wh.ToDate.Format(dto.TimeFormat)
		button.Data.Type = enum.ButtonTypeUserTrackings
		button.Data.Value = int(wh.Warehouse)

		buttons = append(buttons, button)
	}

	tmpMarkup, err := GenerateKeyboard(buttons...)
	if err != nil {
		fmt.Printf("generateKeyboard: %s\n", err.Error())
	}

	msg.ReplyMarkup = tmpMarkup

	return msg
}
