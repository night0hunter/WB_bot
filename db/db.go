package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Postgres struct {
	db *pgxpool.Pool
}

var (
	PgInstance *Postgres
	pgOnce     sync.Once
)

type WarehouseData struct {
	ChatID     int64
	FromDate   time.Time
	toDate     time.Time
	Warehouse  string
	CoeffLimit uint8
	SupplyType string
	IsActive   bool
}

func NewPG(ctx context.Context, connString string) (*Postgres, error) {
	var err error
	pgOnce.Do(func() {
		db, err1 := pgxpool.New(ctx, connString)
		if err1 != nil {
			err = errors.Wrap(err1, "pgpool.New")
			return
		}

		err = db.Ping(ctx)
		if err != nil {
			return
		}

		PgInstance = &Postgres{db}
	})
	if err != nil {
		return nil, errors.Wrap(err, "pgOnce.Do")
	}

	return PgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.db.Close()
}

func (pg *Postgres) SelectQuery(ctx context.Context, ChatID int64) ([]WarehouseData, error) {
	query := `SELECT from_date, to_date, warehouse, coeff_limit, supply_type, is_active FROM supplies WHERE chat_id = (@ChatID)`
	args := pgx.NamedArgs{
		"ChatID": ChatID,
	}

	rows, err := pg.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to scan row: %w", err)
	}

	warehouses := []WarehouseData{}
	for rows.Next() {
		warehouse := WarehouseData{}
		err := rows.Scan(&warehouse.FromDate, &warehouse.toDate, &warehouse.Warehouse, &warehouse.CoeffLimit, &warehouse.SupplyType, &warehouse.IsActive)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}

		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil
}

func (pg *Postgres) InsertQuery(ctx context.Context, params WarehouseData) error {
	query := `INSERT INTO supplies (chat_id, from_date, to_date, Warehouse, coeff_limit, supply_type, is_active) VALUES (@ChatID, @FromDate, @toDate, @Warehouse, @CoeffLimit, @SupplyType, @IsActive)`
	args := pgx.NamedArgs{
		"ChatID":     params.ChatID,
		"FromDate":   params.FromDate,
		"toDate":     params.toDate,
		"Warehouse":  params.Warehouse,
		"CoeffLimit": params.CoeffLimit,
		"SupplyType": params.SupplyType,
		"IsActive":   params.IsActive,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}
