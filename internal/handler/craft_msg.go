package handler

import (
	"fmt"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func SequenceController(prevCommand dto.PrevCommandInfo, comName enum.CommandSequence) (string, error) {
	switch prevCommand.SequenceName {
	case enum.Add:
		msg, err := CraftMsgAdd(prevCommand, comName)
		if err != nil {
			return "", errors.Wrap(err, "CraftMsgAdd")
		}

		return msg, nil
	case enum.Booking:
		msg, err := CraftMsgBooking(prevCommand, comName)
		if err != nil {
			return "", errors.Wrap(err, "CraftMsgBooking")
		}

		return msg, nil
	case enum.Change:
		return "", nil
	default:
		return "", errors.New("unsupported sequence")
	}
}

func CraftMsgAdd(prevCommand dto.PrevCommandInfo, comName enum.CommandSequence) (string, error) {
	data, err := Unmarshal[dto.WarehouseData](prevCommand.Info)
	if err != nil {
		return "", errors.Wrap(err, "Unmarshal")
	}

	return fmt.Sprintf(
		"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
		data.FromDate.Format(dto.TimeFormat),
		data.ToDate.Format(dto.TimeFormat),
		constmsg.WarehouseNames[data.Warehouse],
		constmsg.Coefficients[data.CoeffLimit],
		data.SupplyType,
		BotCommands[comName],
	), nil
}

func CraftMsgBooking(prevCommand dto.PrevCommandInfo, comName enum.CommandSequence) (string, error) {
	data, err := Unmarshal[dto.BookingData](prevCommand.Info)
	if err != nil {
		return "", errors.Wrap(err, "Unmarshal")
	}

	if data.DraftID == uuid.Nil {
		return fmt.Sprintf(
			"Лимит даты бронирования: %s-%s\nID черновика:\nЗащита от бронирования: %s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
			data.FromDate.Format(dto.TimeFormat),
			data.ToDate.Format(dto.TimeFormat),
			constmsg.Coefficients[data.Protection],
			constmsg.WarehouseNames[data.Warehouse],
			constmsg.Coefficients[data.CoeffLimit],
			data.SupplyType,
			BotCommands[comName],
		), nil
	}

	return fmt.Sprintf(
		"Лимит даты бронирования: %s-%s\nID черновика: %s\nЗащита от бронирования: %s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
		data.FromDate.Format(dto.TimeFormat),
		data.ToDate.Format(dto.TimeFormat),
		data.DraftID,
		constmsg.Coefficients[data.Protection],
		constmsg.WarehouseNames[data.Warehouse],
		constmsg.Coefficients[data.CoeffLimit],
		data.SupplyType,
		BotCommands[comName],
	), nil
}
