package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"strconv"
// 	"wb_bot/internal/dto"
// 	"wb_bot/internal/enum"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func main() {
// 	buttons := []dto.Button{
// 		{
// 			Data: dto.ButtonData{
// 				Type:  enum.ButtonTypeCoeffLimit,
// 				Value: 0,
// 			},
// 			Text: "Бесплатно",
// 		},
// 		{
// 			Data: dto.ButtonData{
// 				Type:  enum.ButtonTypeCoeffLimit,
// 				Value: 1,
// 			},
// 			Text: "1x",
// 		},
// 		{
// 			Data: dto.ButtonData{
// 				Type:  enum.ButtonTypeCoeffLimit,
// 				Value: 2,
// 			},
// 			Text: "2x",
// 		},
// 		{
// 			Data: dto.ButtonData{
// 				Type:  enum.ButtonTypeCoeffLimit,
// 				Value: 3,
// 			},
// 			Text: "3x",
// 		}}

// 	bot, err := tgbotapi.NewBotAPI("7232888178:AAFPNqmEmgP-QA2hkHIeOxdv9Zfx5W2Laz0")
// 	if err != nil {
// 		log.Fatalf("tgbotapi.NewBotAPI: %s", err)
// 	}

// 	fmt.Printf("Bot has been started")

// 	updateConfig := tgbotapi.NewUpdate(0)

// 	updateConfig.Timeout = 30

// 	updates := bot.GetUpdatesChan(updateConfig)

// 	for update := range updates {
// 		if update.Message == nil && update.CallbackQuery == nil {
// 			continue
// 		}

// 		if update.Message != nil {
// 			var msg tgbotapi.MessageConfig

// 			tmp, err := generateKeyboard(buttons...)
// 			if err != nil {
// 				fmt.Errorf("generateKeyboard: %w", err)
// 			}

// 			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "hello")
// 			msg.ReplyMarkup = tmp
// 			if _, err := bot.Send(msg); err != nil {
// 				fmt.Errorf("bot.Send: %w", err)
// 			}
// 		}

// 		if update.CallbackQuery != nil {
// 			var buttonData dto.ButtonData

// 			err := json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
// 			if err != nil {
// 				fmt.Errorf("json.Unmarshal", err)
// 			}

// 			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, strconv.Itoa(buttonData.Value))
// 			if _, err := bot.Send(msg); err != nil {
// 				fmt.Errorf("bot.Send: %w", err)
// 			}
// 		}
// 	}

// }

// func generateKeyboard(buttons ...dto.Button) (tgbotapi.InlineKeyboardMarkup, error) {
// 	rows := make([][]tgbotapi.InlineKeyboardButton, len(buttons))
// 	for index, button := range buttons {
// 		// jsonData, err := json.Marshal(button.Data)
// 		// if err != nil {
// 		// 	return tgbotapi.InlineKeyboardMarkup{}, errors.Wrap(err, "json.Marshal")
// 		// }

// 		rows[index] = tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{
// 			Text:         button.Text,
// 			CallbackData: &button.Text,
// 		})
// 	}
// 	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
// }
