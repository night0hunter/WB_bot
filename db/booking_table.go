package db

import (
	"context"
	"wb_bot/internal/dto"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (pg *Postgres) InsertBooking(ctx context.Context, params dto.BookingData) error {
	query := `INSERT INTO bookings (
					chat_id,
					from_date,
					to_date,
					draft_id,
					protection,
					warehouse,
					coeff_limit,
					supply_type
				)
  			  VALUES (@ChatID, @FromDate, @ToDate, @DraftID, @Protection, @Warehouse, @CoeffLimit, @SupplyType)`
	args := pgx.NamedArgs{
		"ChatID":     params.ChatID,
		"FromDate":   params.FromDate,
		"ToDate":     params.ToDate,
		"DraftID":    params.DraftID,
		"Protection": params.Protection,
		"Warehouse":  params.Warehouse,
		"CoeffLimit": params.CoeffLimit,
		"SupplyType": int(params.SupplyType),
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "unable to insert row")
	}

	return nil
}

func (pg *Postgres) SelectBookings(ctx context.Context) ([]dto.BookingData, error) {
	query := `
		SELECT id, 
			   chat_id,
			   from_date,
			   to_date,
			   draft_id,
			   protection,
			   warehouse,
			   coeff_limit,
			   supply_type,
			   is_active
		FROM bookings
		WHERE preorder_id IS NULL AND is_active = TRUE
		ORDER BY from_date ASC
	`

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "unable to scan row")
	}

	bookings := []dto.BookingData{}
	for rows.Next() {
		booking := dto.BookingData{}
		err := rows.Scan(
			&booking.BookingID,
			&booking.ChatID,
			&booking.FromDate,
			&booking.ToDate,
			&booking.DraftID,
			&booking.Protection,
			&booking.Warehouse,
			&booking.CoeffLimit,
			&booking.SupplyType,
			&booking.IsActive)
		if err != nil {
			return nil, errors.Wrap(err, "unable to scan row")
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (pg *Postgres) SelectUserBookings(ctx context.Context, chatID int64) ([]dto.BookingData, error) {
	query := `
		SELECT id,
			   chat_id,
			   from_date,
			   to_date,
			   draft_id,
			   protection,
			   warehouse,
			   coeff_limit,
			   supply_type,
			   is_active
		FROM bookings
		WHERE chat_id = (@ChatID)
		ORDER BY from_date ASC
	`
	args := pgx.NamedArgs{
		"ChatID": chatID,
	}

	rows, err := pg.db.Query(ctx, query, args)
	if err != nil {
		return nil, errors.Wrap(err, "unable to scan row")
	}

	bookings := []dto.BookingData{}
	for rows.Next() {
		booking := dto.BookingData{}
		err := rows.Scan(
			&booking.BookingID,
			&booking.ChatID,
			&booking.FromDate,
			&booking.ToDate,
			&booking.DraftID,
			&booking.Protection,
			&booking.Warehouse,
			&booking.CoeffLimit,
			&booking.SupplyType,
			&booking.IsActive)
		if err != nil {
			return nil, errors.Wrap(err, "unable to scan row")
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (pg *Postgres) SelectBookingStatus(ctx context.Context, bookingID int64) (bool, error) {
	query := `SELECT is_active FROM bookings WHERE id=(@BookingID)`
	args := pgx.NamedArgs{
		"BookingID": bookingID,
	}

	row := pg.db.QueryRow(ctx, query, args)

	var status bool

	err := row.Scan(&status)
	if err != nil {
		return true, errors.Wrap(err, "unable to scan row")
	}

	return status, nil
}

func (pg *Postgres) UpdatePreorderID(ctx context.Context, id int64, preorderID int) error {
	query := `UPDATE bookings SET preorder_id=(@PreorderID) WHERE id=(@TrackingID)`
	args := pgx.NamedArgs{
		"PreorderID": preorderID,
		"TrackingID": id,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "UpdatePreorderID")
	}

	return nil
}

func (pg *Postgres) UpdateBookingStatus(ctx context.Context, bookingID int64, isActive bool) error {
	query := `UPDATE bookings SET is_active=(@IsActive) WHERE id=(@BookingID)`
	args := pgx.NamedArgs{
		"IsActive":  !isActive,
		"BookingID": bookingID,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "UpdateBookingStatus")
	}

	return nil
}

func (pg *Postgres) DeleteBooking(ctx context.Context, bookingID int64) error {
	query := `DELETE FROM bookings WHERE id=(@BookingID)`
	args := pgx.NamedArgs{
		"BookingID": bookingID,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "DeleteBooking")
	}

	return nil
}
