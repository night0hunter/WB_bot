package service

import (
	"context"
	"fmt"
	"strconv"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/utils"

	"github.com/pkg/errors"
)

type Repository interface {
	SelectQuery(ctx context.Context, ChatID int64) ([]dto.WarehouseData, error)
	InsertQuery(ctx context.Context, params dto.WarehouseData) error
	DeleteQuery(ctx context.Context, params dto.WarehouseData) error
}

type Service struct {
	Repository Repository
}

func NewService(rep Repository) *Service {
	return &Service{Repository: rep}
}

func (s *Service) ButtonTypeWarehouseService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error {
	tmpTracking := dto.Trackings[chatID]
	tmpTracking.Warehouse = buttonData.Value
	dto.Trackings[chatID] = tmpTracking

	return nil
}

func (s *Service) ButtonTypeCoeffLimitService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error {
	tmpTracking := dto.Trackings[chatID]
	tmpTracking.CoeffLimit = buttonData.Value
	dto.Trackings[chatID] = tmpTracking

	return nil
}

func (s *Service) ButtonTypeSupplyTypeService(ctx context.Context, chatID int64, buttonData dto.ButtonData) error {
	tmpTracking := dto.Trackings[chatID]
	tmpTracking.SupplyType = strconv.Itoa(buttonData.Value)
	dto.Trackings[chatID] = tmpTracking

	err := s.Repository.InsertQuery(ctx, dto.Trackings[chatID])
	if err != nil {
		return errors.Wrap(err, "Repository.InsertQuery")
	}

	return nil
}

func (s *Service) BotAnswerInputDateService(ctx context.Context, chatID int64, date string) (dto.TrackingDate, error) {
	dateFrom, dateTo, err := utils.ParseTimeRange(date)
	if err != nil {
		return dto.TrackingDate{}, errors.Wrap(err, "utils.ParseTimeRange")
	}

	// usersMutex.Lock()
	dto.Trackings[chatID] = dto.WarehouseData{ChatID: chatID}
	// usersMutex.Unlock()

	tmpTracking := dto.Trackings[chatID]

	tmpTracking.FromDate = dateFrom
	tmpTracking.ToDate = dateTo
	dto.Trackings[chatID] = tmpTracking

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

	tmpTracking := dto.Trackings[chatID]
	tmpTracking.CoeffLimit = parsedCoeff
	dto.Trackings[chatID] = tmpTracking

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
				constmsg.WarehouseNames[wh.Warehouse],
				wh.FromDate.Format(dto.TimeFormat),
				wh.ToDate.Format(dto.TimeFormat),
				wh.CoeffLimit,
				constmsg.SupplyTypes[wh.SupplyType],
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
