package model

import (
	"context"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/handler/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var CommandToKeyboard = map[enum.CommandSequence]func(msg tgbotapi.MessageConfig, data dto.KeyboardData) (tgbotapi.MessageConfig, error){
	// add seq
	enum.BotCommandNameTypeInputDate:       keyboard.DrawCancelKeyboard,
	enum.BotCommandNameTypeInputWarehouse:  keyboard.DrawWarehouseKeyboard,
	enum.BotCommandNameTypeInputCoeffLimit: keyboard.DrawCoeffKeyboard,
	enum.BotCommandNameTypeInputSupplyType: keyboard.DrawSupplyKeyboard,

	// change seq
	enum.BotCommandNameTypeChange:   keyboard.DrawCancelKeyboard,
	enum.BotCommandNameTypeTracking: keyboard.DrawTrackingsKeyboard,
	enum.BotCommandNameTypeAction:   keyboard.DrawActionChoiceKeyboard,

	// booking seq
	enum.BotCommandNameTypeBook:    nil,
	enum.BotCommandNameTypeDraftID: keyboard.DrawCancelKeyboard,

	enum.BotCommandNameTypeSaveStatus: keyboard.DrawSaveStatusKeyboard,
}

var SequenceToFirstCommand = map[enum.Sequences]enum.CommandSequence{
	enum.Add:     enum.BotCommandNameTypeAdd,
	enum.Change:  enum.BotCommandNameTypeChange,
	enum.Booking: enum.BotCommandNameTypeBook,
}

type HandlerStruct interface {
	Question(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error)
	Answer(ctx context.Context, update tgbotapi.Update, tmpData dto.PrevCommandInfo) (dto.PrevCommandInfo, error)
	GetCommandName() enum.CommandSequence
}

type Handler struct {
	Prev    HandlerStruct
	Current HandlerStruct
	Next    HandlerStruct
}
