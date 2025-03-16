package service

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"wb_bot/internal/api"
	"wb_bot/internal/dto"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func (s *Service) GetTrackings(ctx context.Context) ([]dto.MergedResp, error) {
	var result []dto.MergedResp
	var userTrackings []dto.WarehouseData

	response, err := api.GetTrackingsList(ctx, http.Client{Timeout: time.Second * 2})
	if err != nil {
		return []dto.MergedResp{}, errors.Wrap(err, "api.GetTrackingsList")
	}

	sortedResponse := sortResponse(response)

	for _, val := range sortedResponse {
		userTrackings, err = s.Repository.JobSelect(ctx, val[len(val)-1].Date)
		if err != nil {
			return []dto.MergedResp{}, errors.Wrap(err, "Repository.JobSelect")
		}

		break
	}
	// spew.Dump(sortedResponse["Коледино"])
	// userTrackings, err := s.Repository.JobSelect(ctx, sortedResponse["Коледино"][len(sortedResponse["Коледино"])-1].Date)

	for _, tr := range userTrackings {
		for j := 0; j < len(sortedResponse[tr.Warehouse]); j++ {
			if sortedResponse[tr.Warehouse][j].Coefficient == -1 {
				continue
			}

			if tr.CoeffLimit < sortedResponse[tr.Warehouse][j].Coefficient {
				continue
			}

			if tr.SupplyType != strconv.Itoa(sortedResponse[tr.Warehouse][j].BoxTypeID) {
				continue
			}

			if sortedResponse[tr.Warehouse][j].Date.Before(tr.FromDate) || sortedResponse[tr.Warehouse][j].Date.After(tr.ToDate) {
				continue
			}

			tmp := dto.MergedResp{
				UserID:          tr.ChatID,
				Date:            sortedResponse[tr.Warehouse][j].Date,
				Coefficient:     sortedResponse[tr.Warehouse][j].Coefficient,
				WarehouseID:     sortedResponse[tr.Warehouse][j].WarehouseID,
				WarehouseName:   sortedResponse[tr.Warehouse][j].WarehouseName,
				BoxTypeName:     sortedResponse[tr.Warehouse][j].BoxTypeName,
				BoxTypeID:       sortedResponse[tr.Warehouse][j].BoxTypeID,
				IsSortingCenter: sortedResponse[tr.Warehouse][j].IsSortingCenter,
				IsAvtive:        tr.IsActive,
			}

			result = append(result, tmp)
		}
	}

	spew.Dump(result)

	return result, nil
}

func sortResponse(response []dto.Response) map[int][]dto.Response {
	var result = map[int][]dto.Response{}

	for _, rp := range response {
		if len(result[rp.WarehouseID]) == 0 {
			result[rp.WarehouseID] = append(result[rp.WarehouseID], rp)

			continue
		}

		for j := 0; j < len(result[rp.WarehouseID]); j++ {
			if rp.Date.Before(result[rp.WarehouseID][j].Date) {
				result[rp.WarehouseID] = append(result[rp.WarehouseID][:j+1], result[rp.WarehouseID][j:]...)
				result[rp.WarehouseID][j] = rp

				break
			}
		}
	}

	return result
}
