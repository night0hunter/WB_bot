package db

import (
	"context"
	"wb_bot/internal/dto"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (pg *Postgres) SelectState(ctx context.Context, id int64) (dto.PrevCommandInfo, error) {
	query := `SELECT seq_name, com_name, mes_id, info, keyboard FROM state WHERE chat_id = (@ID)`
	args := pgx.NamedArgs{
		"ID": id,
	}

	row := pg.db.QueryRow(ctx, query, args)

	prevCommand := dto.PrevCommandInfo{}
	err := row.Scan(
		&prevCommand.SequenceName,
		&prevCommand.CommandName,
		&prevCommand.MessageID,
		&prevCommand.Info,
		&prevCommand.KeyboardInfo,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return dto.PrevCommandInfo{}, errors.Wrap(err, "row.Scan: unable to scan row")
	}
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return dto.PrevCommandInfo{}, nil
	}

	return prevCommand, nil
}

func (pg *Postgres) UpdateState(ctx context.Context, id int64, prevCommand dto.PrevCommandInfo) error {
	query := `UPDATE state SET seq_name = (@seqName), com_name = (@comName), mes_id = (@mesID), info = (@info), keyboard = (@keyboard) WHERE chat_id = (@ID)`
	args := pgx.NamedArgs{
		"ID":       id,
		"seqName":  prevCommand.SequenceName,
		"comName":  prevCommand.CommandName,
		"mesID":    prevCommand.MessageID,
		"info":     prevCommand.Info,
		"keyboard": prevCommand.KeyboardInfo,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "pg.db.Exec")
	}

	return nil
}

func (pg *Postgres) InsertState(ctx context.Context, id int64, prevCommand dto.PrevCommandInfo) error {
	query := `INSERT INTO state (chat_id, seq_name, com_name, mes_id, info, keyboard) VALUES (@ChatID, @seqName, @comName, @mesID, @info, @keyboard)`
	args := pgx.NamedArgs{
		"ChatID":   id,
		"seqName":  prevCommand.SequenceName,
		"comName":  prevCommand.CommandName,
		"mesID":    prevCommand.MessageID,
		"info":     prevCommand.Info,
		"keyboard": prevCommand.KeyboardInfo,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "unable to insert row")
	}

	return nil
}

func (pg *Postgres) DeleteState(ctx context.Context, id int64) error {
	query := `DELETE FROM state WHERE chat_id=(@ID)`
	args := pgx.NamedArgs{
		"ID": id,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "pg.db.Exec")
	}

	return nil
}
