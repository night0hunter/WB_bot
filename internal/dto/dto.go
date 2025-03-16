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
	TrackingID    int64
	ChatID        int64
	FromDate      time.Time
	ToDate        time.Time
	Warehouse     int
	WarehouseName string
	CoeffLimit    int
	SupplyType    string
	IsActive      bool
}

// type WarehouseDataExt struct {
// 	WarehouseData
// 	WarehouseName string
// }

var Trackings = map[int64]WarehouseData{}

type TrackingDate struct {
	DateFrom time.Time
	DateTo   time.Time
}

type TrackingStatus struct {
	UserID int64
	Status int
}

type MergedResp struct {
	UserID          int64
	Date            time.Time
	Coefficient     int
	WarehouseID     int
	WarehouseName   string
	BoxTypeName     string
	BoxTypeID       int
	IsSortingCenter bool
	IsAvtive        bool
}

type Response struct {
	Date            time.Time `json:"date"`
	Coefficient     int       `json:"coefficient"`
	WarehouseID     int       `json:"warehouseID"`
	WarehouseName   string    `json:"warehouseName"`
	BoxTypeName     string    `json:"boxTypeName"`
	BoxTypeID       int       `json:"boxTypeId"`
	IsSortingCenter bool      `json:"isSortingCenter"`
}
