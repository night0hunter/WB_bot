package dto

import (
	"time"
	"wb_bot/internal/enum"

	uuid "github.com/google/uuid"
)

const TimeFormat = "02.01.2006"

type PrevCommandInfo struct {
	SequenceName enum.Sequences
	CommandName  enum.CommandSequence
	MessageID    int
	Info         []byte
	KeyboardInfo []byte
}

type ButtonData struct {
	Type  enum.ButtonType
	Value int
}

type Button struct {
	Data ButtonData
	Text string
}

type WarehouseData struct {
	TrackingID  int64
	ChatID      int64
	SendingDate time.Time
	FromDate    time.Time
	ToDate      time.Time
	Warehouse   int
	CoeffLimit  *int
	SupplyType  string
	IsActive    bool
}

type ChangeStatusInfo struct {
	TrackingID int64
	Choice     int
}

type BookingData struct {
	BookingID  int64
	ChatID     int64
	DraftID    uuid.UUID
	FromDate   time.Time
	ToDate     time.Time
	Protection *int
	Warehouse  int
	CoeffLimit *int
	SupplyType string
}

type TrackingDate struct {
	DateFrom time.Time
	DateTo   time.Time
}

type TrackingStatus struct {
	UserID int64
	Status int
}

type KeyboardData struct {
	Warehouses []WarehouseData
}

type MergedResp struct {
	TrackingID      int64
	UserID          int64
	SendingDate     time.Time
	Date            time.Time
	Coefficient     int
	WarehouseID     int
	WarehouseName   string
	BoxTypeName     string
	BoxTypeID       int
	IsSortingCenter bool
	IsActive        bool
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
