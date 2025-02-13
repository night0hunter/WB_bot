package dto

import (
	"time"
	"wb_bot/internal/enum"
)

const TimeFormat = "02.01.2006"

type ButtonData struct {
	Type  enum.ButtonType
	Value int
}

type Button struct {
	Data ButtonData
	Text string
}

type WarehouseData struct {
	TrackingID int64
	ChatID     int64
	FromDate   time.Time
	ToDate     time.Time
	Warehouse  int
	CoeffLimit int
	SupplyType string
	IsActive   bool
}

var Trackings = map[int64]WarehouseData{}

// type CheckWarehouse struct {
// 	Text     string
// 	IsActive bool
// }
