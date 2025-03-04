package db

import (
	"context"
	"sync"
	"wb_bot/internal/dto"

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

func (pg *Postgres) SelectQuery(ctx context.Context, ChatID int64) ([]dto.WarehouseData, error) {
	query := `
		SELECT chat_id,
			   from_date,
			   to_date,
			   warehouse,
			   coeff_limit,
			   supply_type,
			   is_active,
			   id
		FROM supplies
		WHERE chat_id = (@ChatID)`

	args := pgx.NamedArgs{
		"ChatID": ChatID,
	}

	rows, err := pg.db.Query(ctx, query, args)
	if err != nil {
		return nil, errors.Wrap(err, "unable to scan row")
	}

	warehouses := []dto.WarehouseData{}
	for rows.Next() {
		warehouse := dto.WarehouseData{}
		err := rows.Scan(
			&warehouse.ChatID,
			&warehouse.FromDate,
			&warehouse.ToDate,
			&warehouse.Warehouse,
			&warehouse.CoeffLimit,
			&warehouse.SupplyType,
			&warehouse.IsActive,
			&warehouse.TrackingID)
		if err != nil {
			return nil, errors.Wrap(err, "unable to scan row")
		}

		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil
}

func (pg *Postgres) InsertQuery(ctx context.Context, params dto.WarehouseData) error {
	query := `INSERT INTO supplies (chat_id, from_date, to_date, Warehouse, coeff_limit, supply_type) VALUES (@ChatID, @FromDate, @ToDate, @Warehouse, @CoeffLimit, @SupplyType)`
	args := pgx.NamedArgs{
		"ChatID":     params.ChatID,
		"FromDate":   params.FromDate,
		"ToDate":     params.ToDate,
		"Warehouse":  params.Warehouse,
		"CoeffLimit": params.CoeffLimit,
		"SupplyType": params.SupplyType,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "unable to insert row")
	}

	return nil
}

func (pg *Postgres) InsertTrackingStatus(ctx context.Context, params dto.TrackingStatus) error {
	query := `INSERT INTO tracking_status (user_id, status) VALUES (@ChatID, @Status)`
	args := pgx.NamedArgs{
		"ChatID": params.UserID,
		"Status": params.Status,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "InsertTrackingStatus: unable to insert row")
	}

	return nil
}

// func (pg *Postgres) SelectTrackingStatus(ctx context.Context, chatID int64) error {
// 	query := `SELECT status FROM tracking_status WHERE user_id=(@ChatID)`
// 	args := pgx.NamedArgs{
// 		"ChatID": chatID,
// 	}

// 	_, err := pg.db.Exec(ctx, query, args)
// 	if err != nil {
// 		return errors.Wrap(err, "SelectTrackingStatus: unable to insert row")
// 	}

// 	return nil
// }

func (pg *Postgres) SelectTrackingStatus(ctx context.Context, chatID int64, trackingID int64) (bool, error) {
	query := `SELECT is_active FROM supplies WHERE chat_id=(@ChatID) AND id=(@TrackingID)`
	args := pgx.NamedArgs{
		"ChatID":     chatID,
		"TrackingID": trackingID,
	}

	row := pg.db.QueryRow(ctx, query, args)

	var status bool

	err := row.Scan(&status)
	if err != nil {
		return true, errors.Wrap(err, "unable to scan row")
	}

	return status, nil
}

func (pg *Postgres) ChangeTrackingStatus(ctx context.Context, trackingID int64, isActive bool) error {
	query := `UPDATE supplies SET is_active=(@IsActive) WHERE id=(@TrackingID)`
	args := pgx.NamedArgs{
		"IsActive":   !isActive,
		"TrackingID": trackingID,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "ChangeTrackingStatus")
	}

	return nil
}
