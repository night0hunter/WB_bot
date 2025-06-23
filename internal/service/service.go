package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/utils"

	"github.com/pkg/errors"
)

type Repository interface {
	SelectQuery(ctx context.Context, ChatID int64) ([]dto.WarehouseData, error)
	InsertQuery(ctx context.Context, params dto.WarehouseData) error
	SelectTrackingStatus(ctx context.Context, chatID int64, trackingID int64) (bool, error)
	ChangeTrackingStatus(ctx context.Context, chatID int64, isActive bool) error
	DeleteTracking(ctx context.Context, trackingID int64) error
	JobSelect(ctx context.Context, date time.Time) ([]dto.WarehouseData, error)
	UpdateSendingTime(ctx context.Context, date time.Time, id int64) error

	SelectState(ctx context.Context, id int64) (dto.PrevCommandInfo, error)
	UpdateState(ctx context.Context, id int64, prevCommand dto.PrevCommandInfo) error
	InsertState(ctx context.Context, id int64, prevCommand dto.PrevCommandInfo) error
	DeleteState(ctx context.Context, id int64) error
}

type Service struct {
	Repository Repository
}

func NewService(rep Repository) *Service {
	return &Service{Repository: rep}
}

func (s *Service) BotAnswerInputDateService(ctx context.Context, chatID int64, date string) (dto.TrackingDate, error) {
	dateFrom, dateTo, err := utils.ParseTimeRange(date)
	if err != nil {
		return dto.TrackingDate{}, errors.Wrap(err, "utils.ParseTimeRange")
	}

	return dto.TrackingDate{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}, nil
}

func (s *Service) BotAnswerInputCoeffLimitService(ctx context.Context, chatID int64, coeffLimit string) (int, error) {
	parsedCoeff, err := utils.ParseCoeffLimit(coeffLimit)
	if err != nil {
		return 0, errors.Wrap(err, "utils.ParseCoeffLimit")
	}

	return parsedCoeff, nil
}

func (s *Service) BotSlashCommandTypeCheckService(ctx context.Context, chatID int64) ([]string, error) {
	var warehouseStrs []string

	warehouses, err := s.Repository.SelectQuery(ctx, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "Repository.SelectQuery")
	}

	for _, wh := range warehouses {
		warehouseStrs = append(
			warehouseStrs,
			fmt.Sprintf(
				"Склад: %s\nДата отслеживания: %s-%s\nЛимит коэффициента: x%d и меньше\nТип поставки: %s\nАктивно/Неактивно: %s",
				constmsg.WarehouseNames[int(wh.Warehouse)],
				wh.FromDate.Format(dto.TimeFormat),
				wh.ToDate.Format(dto.TimeFormat),
				wh.CoeffLimit,
				wh.SupplyType,
				utils.BoolToActiveRU(wh.IsActive),
			),
		)
	}

	return warehouseStrs, nil
}

func (s *Service) BotSlashCommandTypeHelpService(ctx context.Context, chatID int64) string {
	var text string

	for cmd, desc := range constmsg.BotSlashCommandsHelp {
		text += cmd + " - " + desc + "\n"
	}

	return text
}

func (s *Service) BotSlashCommandTypeChange(ctx context.Context, chatID int64) ([]dto.WarehouseData, error) {
	warehouses, err := s.Repository.SelectQuery(ctx, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "Repository.SelectQuery")
	}

	return warehouses, nil
}

func (s *Service) ChangeStatusService(ctx context.Context, chatID, trackingID int64) error {
	status, err := s.Repository.SelectTrackingStatus(ctx, chatID, trackingID)
	if err != nil {
		return errors.Wrap(err, "Repository.SelectTrackingStatus")
	}

	err = s.Repository.ChangeTrackingStatus(ctx, trackingID, status)
	if err != nil {
		return errors.Wrap(err, "Repository.ChangeTrackingStatus")
	}

	return nil
}

func (s *Service) DeleteTrackingService(ctx context.Context, chatID int64, trackingID int64) error {
	err := s.Repository.DeleteTracking(ctx, trackingID)
	if err != nil {
		return errors.Wrap(err, "Repository.DeleteTracking")
	}

	return nil
}

func (s *Service) AddSequenceEndService(ctx context.Context, chatID int64, data []byte) error {
	var unmarshData dto.WarehouseData
	err := json.Unmarshal(data, &unmarshData)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	err = s.Repository.InsertQuery(ctx, unmarshData)
	if err != nil {
		return errors.Wrap(err, "Repository.InsertQuery")
	}

	return nil
}

func (s *Service) SelectState(ctx context.Context, chatID int64) (dto.PrevCommandInfo, error) {
	prevCommand, err := s.Repository.SelectState(ctx, chatID)
	if err != nil {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "Repository.SelectState")
	}

	return prevCommand, nil
}

func (s *Service) InsertState(ctx context.Context, chatID int64, prevCommand dto.PrevCommandInfo) error {
	err := s.Repository.InsertState(ctx, chatID, prevCommand)
	if err != nil {
		return errors.Wrap(err, "Repository.InsertState")
	}

	return nil
}

func (s *Service) UpdateState(ctx context.Context, chatID int64, prevCommand dto.PrevCommandInfo) error {
	err := s.Repository.UpdateState(ctx, chatID, prevCommand)
	if err != nil {
		return errors.Wrap(err, "Repository.UpdateState")
	}

	return nil
}

func (s *Service) DeleteState(ctx context.Context, chatID int64) error {
	err := s.Repository.DeleteState(ctx, chatID)
	if err != nil {
		return errors.Wrap(err, "Repository.DeleteState")
	}

	return nil
}
