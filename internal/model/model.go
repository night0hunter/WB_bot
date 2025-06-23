package model

import (
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var CommandToKeyboard = map[enum.CommandSequences]func(msg tgbotapi.MessageConfig, data dto.KeyboardData) (tgbotapi.MessageConfig, error){
	enum.BotCommandNameTypeInputDate:       keyboard.DrawCancelKeyboard,
	enum.BotCommandNameTypeInputWarehouse:  keyboard.DrawWarehouseKeyboard,
	enum.BotCommandNameTypeInputCoeffLimit: keyboard.DrawCoeffKeyboard,
	enum.BotCommandNameTypeInputSupplyType: keyboard.DrawSupplyKeyboard,
	enum.BotCommandNameTypeChange:          keyboard.DrawCancelKeyboard,
	enum.BotCommandNameTypeTracking:        keyboard.DrawTrackingsKeyboard,
	enum.BotCommandNameTypeAction:          keyboard.DrawActionChoiceKeyboard,

	enum.BotCommandNameTypeSaveStatus: keyboard.DrawSaveStatusKeyboard,
}
