package addHandler

import (
	"context"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service interface {
	SelectState(ctx context.Context, chatID int64) (dto.PrevCommandInfo, error)
	InsertState(ctx context.Context, chatID int64, prevCommand dto.PrevCommandInfo) error
	UpdateState(ctx context.Context, chatID int64, prevCommand dto.PrevCommandInfo) error
	DeleteState(ctx context.Context, chatID int64) error
	BotAnswerInputDateService(ctx context.Context, chatID int64, date string) (dto.TrackingDate, error)
	BotAnswerInputCoeffLimitService(ctx context.Context, chatID int64, coeffLimit string) (int, error)
	AddSequenceEndService(ctx context.Context, chatID int64, data []byte) error
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
		enum.BotCommandNameTypeAdd: {
			Prev:    nil,
			Current: nil,
			Next:    &InputDateHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputDate},
		},
		enum.BotCommandNameTypeSaveStatus: {
			Prev:    nil,
			Current: &SaveStatusHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeSaveStatus},
			Next:    &InputDateHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputDate},
		},
		enum.BotCommandNameTypeInputDate: {
			Prev:    nil,
			Current: &InputDateHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputDate},
			Next:    &WarehouseHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputWarehouse},
		},
		enum.BotCommandNameTypeInputWarehouse: {
			Prev:    &InputDateHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputDate},
			Current: &WarehouseHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputWarehouse},
			Next:    &CoeffLimitHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputCoeffLimit},
		},
		enum.BotCommandNameTypeInputCoeffLimit: {
			Prev:    &WarehouseHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputWarehouse},
			Current: &CoeffLimitHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputCoeffLimit},
			Next:    &SupplyTypeHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputSupplyType},
		},
		enum.BotCommandNameTypeInputSupplyType: {
			Prev:    &CoeffLimitHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputCoeffLimit},
			Current: &SupplyTypeHandler{bot: bot, service: svc, commandName: enum.BotCommandNameTypeInputSupplyType},
			Next:    nil,
		},
	}
}
