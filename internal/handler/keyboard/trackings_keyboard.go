package keyboard

import (
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func DrawTrackingsKeyboard(msg tgbotapi.MessageConfig, data dto.KeyboardData) (tgbotapi.MessageConfig, error) {
	var buttons []dto.Button
	var button dto.Button

	for _, wh := range data.Warehouses {
		button.Text = constmsg.WarehouseNames[int(wh.Warehouse)] + " " + wh.FromDate.Format(dto.TimeFormat) + "-" + wh.ToDate.Format(dto.TimeFormat)
		button.Data.Type = enum.ButtonTypeUserTrackingChoice
		button.Data.Value = int(wh.TrackingID)

		buttons = append(buttons, button)
	}

	button = dto.Button{
		Data: dto.ButtonData{
			Type:  enum.ButtonTypeUserTrackingChoice,
			Value: -1,
		},
		Text: "Отмена",
	}
	buttons = append(buttons, button)

	tmpMarkup, err := GenerateKeyboard(buttons...)
	if err != nil {
		return msg, errors.Wrap(err, "GenerateKeyboard")
	}

	msg.ReplyMarkup = tmpMarkup

	return msg, nil
}
