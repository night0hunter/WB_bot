package keyboard

import (
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func DrawBookingsKeyboard(msg tgbotapi.MessageConfig, data []byte) (tgbotapi.MessageConfig, error) {
	var buttons []dto.Button
	var button dto.Button

	keyboardInfo, err := utils.Unmarshal[[]dto.BookingData](data)
	if err != nil {
		return msg, errors.Wrap(err, "utils.Unmarshal")
	}

	for _, wh := range keyboardInfo {
		button.Text = constmsg.WarehouseNames[wh.Warehouse] + " " + wh.FromDate.Format(dto.TimeFormat) + "-" + wh.ToDate.Format(dto.TimeFormat)
		button.Data.Type = enum.ButtonTypeUserBookingsChoice
		button.Data.Value = int(wh.BookingID)

		buttons = append(buttons, button)
	}

	button = dto.Button{
		Data: dto.ButtonData{
			Type:  enum.ButtonTypeUserBookingsChoice,
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
