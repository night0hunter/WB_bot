package keyboard

import (
	"encoding/json"
	"wb_bot/internal/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

// func GenerateKeyboard(buttons ...dto.Button) (tgbotapi.InlineKeyboardMarkup, error) {
// 	keyboardButtons := make([]tgbotapi.InlineKeyboardButton, 0, len(buttons))

// 	for _, button := range buttons {
// 		jsonData, err := json.Marshal(button.Data)
// 		if err != nil {
// 			return tgbotapi.InlineKeyboardMarkup{}, errors.Wrap(err, "json.Marshal")
// 		}

// 		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(button.Text, string(jsonData)))
// 	}

// 	return tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
// 	), nil
// }

func GenerateKeyboard(buttons ...dto.Button) (tgbotapi.InlineKeyboardMarkup, error) {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(buttons))
	for _, button := range buttons {
		jsonData, err := json.Marshal(button.Data)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, errors.Wrap(err, "json.Marshal")
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(button.Text, string(jsonData))))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}
