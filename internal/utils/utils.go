package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	TimeFormat         = "02.01.2006"
	MoscowLocationName = "Europe/Moscow"

	active   = "Активно"
	inactive = "Неактивно"
)

var moscowLocation *time.Location

func init() {
	var err error
	moscowLocation, err = time.LoadLocation(MoscowLocationName)
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

	dateFrom, err := time.ParseInLocation(TimeFormat, datesRaw[0], moscowLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "dateFrom: time.ParseInLocation")
	}

	dateTo, err := time.ParseInLocation(TimeFormat, datesRaw[1], moscowLocation)
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
