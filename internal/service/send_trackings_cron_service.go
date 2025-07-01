package service

import (
	"context"
	"net/http"
	"time"
	"wb_bot/internal/api"
	"wb_bot/internal/dto"
	"wb_bot/internal/utils"

	"github.com/pkg/errors"
)

func (s *Service) GetTrackings(ctx context.Context) ([]dto.MergedResp, error) {
	var result []dto.MergedResp
	var userTrackings []dto.WarehouseData

	response, err := api.GetTrackingsList(ctx, http.Client{Timeout: time.Second * 2})
	if err != nil {
		return []dto.MergedResp{}, errors.Wrap(err, "api.GetTrackingsList")
	}

	sortedResponse := utils.SortResponse(response)

	for _, val := range sortedResponse {
		userTrackings, err = s.Repository.JobSelect(ctx, val[len(val)-1].Date)
		if err != nil {
			return []dto.MergedResp{}, errors.Wrap(err, "Repository.JobSelect")
		}

		break
	}

	for _, tr := range userTrackings {
		for j := 0; j < len(sortedResponse[tr.Warehouse]); j++ {
			if sortedResponse[tr.Warehouse][j].Coefficient == -1 {
				continue
			}

			if *tr.CoeffLimit < sortedResponse[tr.Warehouse][j].Coefficient {
				continue
			}

			if tr.SupplyType != sortedResponse[tr.Warehouse][j].BoxTypeName {
				continue
			}

			if sortedResponse[tr.Warehouse][j].Date.Before(tr.FromDate) || sortedResponse[tr.Warehouse][j].Date.After(tr.ToDate) {
				continue
			}

			if tr.SendingDate.Add(time.Minute * 5).After(time.Now()) {
				continue
			}

			tmp := dto.MergedResp{
				TrackingID:      tr.TrackingID,
				UserID:          tr.ChatID,
				Date:            sortedResponse[tr.Warehouse][j].Date,
				Coefficient:     sortedResponse[tr.Warehouse][j].Coefficient,
				WarehouseID:     sortedResponse[tr.Warehouse][j].WarehouseID,
				WarehouseName:   sortedResponse[tr.Warehouse][j].WarehouseName,
				BoxTypeName:     sortedResponse[tr.Warehouse][j].BoxTypeName,
				BoxTypeID:       sortedResponse[tr.Warehouse][j].BoxTypeID,
				IsSortingCenter: sortedResponse[tr.Warehouse][j].IsSortingCenter,
				IsActive:        tr.IsActive,
			}

			result = append(result, tmp)
		}
	}

	return result, nil
}

func (s *Service) KeepSendingTime(ctx context.Context, tracking dto.MergedResp) error {
	err := s.Repository.UpdateSendingTime(ctx, time.Now(), tracking.TrackingID)
	if err != nil {
		return errors.Wrap(err, "Repository.UpdateSendingDate")
	}

	return nil
}
