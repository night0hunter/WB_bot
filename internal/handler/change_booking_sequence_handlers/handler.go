package changeBookingHandler

import (
	"context"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service interface {
	BotSlashCommandTypeChangeBooking(ctx context.Context, chatID int64) ([]dto.BookingData, error)
	DeleteBookingService(ctx context.Context, bookingID int64) error
	ChangeBookingStatusService(ctx context.Context, bookingID int64) error
}

type ActionChoiceBookingHandler struct {
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
		enum.BotCommandNameTypeChangeBooking: {
			Prev:    nil,
			Current: nil,
			Next:    &BookingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
		},
		enum.BotCommandNameTypeSaveStatus: {
			Prev:    nil,
			Current: &SaveStatusHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeSaveStatus},
			Next:    &BookingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
		},
		enum.BotCommandNameTypeTracking: {
			Prev:    nil,
			Current: &BookingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
			Next:    &ActionChoiceBookingHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeAction},
		},
		enum.BotCommandNameTypeAction: {
			Prev:    &BookingChoiceHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeTracking},
			Current: &ActionChoiceBookingHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeAction},
			Next:    nil,
		},
	}
}
