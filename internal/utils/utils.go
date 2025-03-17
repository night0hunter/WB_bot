package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"wb_bot/internal/dto"

	"github.com/pkg/errors"
)

const (
	TimeFormat         = "02.01.2006"
	MoscowLocationName = "Europe/Moscow"

	active   = "Активно"
	inactive = "Неактивно"
)

var MoscowLocation *time.Location

func init() {
	var err error
	MoscowLocation, err = time.LoadLocation(MoscowLocationName)
	if err != nil {
		log.Fatal(err)
	}
}

// returns begin time and end time
func ParseTimeRange(dateString string) (time.Time, time.Time, error) {
	datesRaw := strings.Split(dateString, "-")
	if len(datesRaw) != 2 {
		return time.Time{}, time.Time{}, errors.New("There must be 2 dates")
	}

	dateFrom, err := time.ParseInLocation(TimeFormat, datesRaw[0], MoscowLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "dateFrom: time.ParseInLocation")
	}

	dateTo, err := time.ParseInLocation(TimeFormat, datesRaw[1], MoscowLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "dateTo: time.ParseInLocation")
	}

	return dateFrom, dateTo, nil
}

func ParseCoeffLimit(coeff string) (int, error) {
	parsed, err := strconv.Atoi(coeff)
	if err != nil {
		return 0, fmt.Errorf("coeffLimit: strconv.Atoi: %w", err)
	}

	return parsed, nil
}

func BoolToActiveRU(input bool) string {
	if input {
		return active
	}

	return inactive
}

func SortResponse(response []dto.Response) map[int][]dto.Response {
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
