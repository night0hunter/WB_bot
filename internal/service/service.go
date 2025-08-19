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
	SelectTrackingStatus(ctx context.Context, trackingID int64) (bool, error)
	SelectBookingStatus(ctx context.Context, bookingID int64) (bool, error)
	SelectBookings(ctx context.Context) ([]dto.BookingData, error)
	SelectUserBookings(ctx context.Context, chatID int64) ([]dto.BookingData, error)

	InsertTracking(ctx context.Context, params dto.WarehouseData) error
	InsertBooking(ctx context.Context, params dto.BookingData) error

	UpdateTrackingStatus(ctx context.Context, trackingID int64, isActive bool) error
	UpdateBookingStatus(ctx context.Context, bookingID int64, isActive bool) error
	DeleteTracking(ctx context.Context, trackingID int64) error
	DeleteBooking(ctx context.Context, trackingID int64) error
	JobSelect(ctx context.Context, dateTo time.Time) ([]dto.WarehouseData, error)
	UpdateSendingTime(ctx context.Context, date time.Time, id int64) error
	UpdatePreorderID(ctx context.Context, id int64, preorderID int) error

	SelectState(ctx context.Context, id int64) (dto.PrevCommandInfo, error)
	UpdateState(ctx context.Context, id int64, prevCommand dto.PrevCommandInfo) error
	InsertState(ctx context.Context, id int64, prevCommand dto.PrevCommandInfo) error
	DeleteState(ctx context.Context, id int64) error
}

type Adapter interface {
	GetTrackingsList(ctx context.Context, url string) ([]dto.Response, error)
	GetWarehouseGoodsV2(ctx context.Context, input dto.GetWarehouseGoodsV2Request, url string) (dto.GetWarehouseGoodsV2Response, error)
	Create(ctx context.Context, input dto.GetCreateRequest, url string) (dto.GetCreateResponse, error)
}

type Service struct {
	Repository Repository
	Adapter    Adapter
}

func New(rep Repository, adp Adapter) *Service {
	return &Service{
		Repository: rep,
		Adapter:    adp,
	}
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

func (s *Service) BotSlashCommandTypeCheckTrackingsService(ctx context.Context, chatID int64) ([]string, error) {
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
				*wh.CoeffLimit,
				constmsg.SupplyTypes[wh.SupplyType],
				utils.BoolToActiveRU(wh.IsActive),
			),
		)
	}

	return warehouseStrs, nil
}

func (s *Service) BotSlashCommandTypeCheckBookingsService(ctx context.Context, chatID int64) ([]string, error) {
	var warehouseStrs []string

	bookings, err := s.Repository.SelectUserBookings(ctx, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "Repository.SelectQuery")
	}

	for _, wh := range bookings {
		warehouseStrs = append(
			warehouseStrs,
			fmt.Sprintf(
				"Лимит даты бронирования: %s-%s\nID черновика: %s\nЗащита от бронирования: %d\nСклад: %s\nЛимит коэффициента: %d\nТип поставки: %s\nАктивно/Неактивно: %s",
				wh.FromDate.Format(dto.TimeFormat),
				wh.ToDate.Format(dto.TimeFormat),
				wh.DraftID,
				*wh.Protection,
				constmsg.WarehouseNames[wh.Warehouse],
				*wh.CoeffLimit,
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

func (s *Service) BotSlashCommandTypeChange(ctx context.Context, chatID int64) ([]dto.WarehouseData, error) {
	warehouses, err := s.Repository.SelectQuery(ctx, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "Repository.SelectQuery")
	}

	return warehouses, nil
}

func (s *Service) BotSlashCommandTypeChangeBooking(ctx context.Context, chatID int64) ([]dto.BookingData, error) {
	bookings, err := s.Repository.SelectUserBookings(ctx, chatID)
	if err != nil {
		return nil, errors.Wrap(err, "Repository.SelectUserBookings")
	}

	return bookings, nil
}

func (s *Service) ChangeStatusService(ctx context.Context, trackingID int64) error {
	status, err := s.Repository.SelectTrackingStatus(ctx, trackingID)
	if err != nil {
		return errors.Wrap(err, "Repository.SelectTrackingStatus")
	}

	err = s.Repository.UpdateTrackingStatus(ctx, trackingID, status)
	if err != nil {
		return errors.Wrap(err, "Repository.ChangeTrackingStatus")
	}

	return nil
}

func (s *Service) ChangeBookingStatusService(ctx context.Context, bookingID int64) error {
	status, err := s.Repository.SelectBookingStatus(ctx, bookingID)
	if err != nil {
		return errors.Wrap(err, "Repository.SelectTrackingStatus")
	}

	err = s.Repository.UpdateBookingStatus(ctx, bookingID, status)
	if err != nil {
		return errors.Wrap(err, "Repository.ChangeTrackingStatus")
	}

	return nil
}

func (s *Service) DeleteTrackingService(ctx context.Context, trackingID int64) error {
	err := s.Repository.DeleteTracking(ctx, trackingID)
	if err != nil {
		return errors.Wrap(err, "Repository.DeleteTracking")
	}

	return nil
}

func (s *Service) DeleteBookingService(ctx context.Context, bookingID int64) error {
	err := s.Repository.DeleteBooking(ctx, bookingID)
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

	err = s.Repository.InsertTracking(ctx, unmarshData)
	if err != nil {
		return errors.Wrap(err, "Repository.InsertQuery")
	}

	return nil
}

func (s *Service) BookSequenceEndService(ctx context.Context, chatID int64, data []byte) error {
	var unmarshData dto.BookingData
	err := json.Unmarshal(data, &unmarshData)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	err = s.Repository.InsertBooking(ctx, unmarshData)
	if err != nil {
		return errors.Wrap(err, "Repository.InsertBooking")
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
