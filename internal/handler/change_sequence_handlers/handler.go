package changeHandler

import (
	"context"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service interface {
	BotSlashCommandTypeChange(ctx context.Context, chatID int64) ([]dto.WarehouseData, error)
	DeleteTrackingService(ctx context.Context, chatID int64, trackingID int64) error
	ChangeStatusService(ctx context.Context, chatID, trackingID int64) error
}

type ActionChoiceHandler struct {
	bot         *tgbotapi.BotAPI
	service     Service
	commandName enum.CommandSequence
}

func New(bot *tgbotapi.BotAPI, svc Service) map[enum.CommandSequence]struct {
	Prev    model.HandlerStruct
	Current model.HandlerStruct
	Next    model.HandlerStruct
} {
	return map[enum.CommandSequence]struct {
		Prev    model.HandlerStruct
		Current model.HandlerStruct
		Next    model.HandlerStruct
	}{
		enum.BotCommandNameTypeChange: {
			Prev:    nil,
			Current: nil,
			Next:    &TrackingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
		},
		enum.BotCommandNameTypeSaveStatus: {
			Prev:    nil,
			Current: &SaveStatusHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeSaveStatus},
			Next:    &TrackingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
		},
		enum.BotCommandNameTypeTracking: {
			Prev:    nil,
			Current: &TrackingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
			Next:    &ActionChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeAction},
		},
		enum.BotCommandNameTypeAction: {
			Prev:    &TrackingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
			Current: &ActionChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeAction},
			Next:    nil,
		},
	}
}
