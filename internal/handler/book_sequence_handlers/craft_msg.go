package bookHandler

import (
	"fmt"
	"strconv"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/utils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func CraftMessage(prevCommand dto.PrevCommandInfo, comName enum.CommandSequence) (string, error) {
	data, err := utils.Unmarshal[dto.BookingData](prevCommand.Info)
	if err != nil {
		return "", errors.Wrap(err, "Unmarshal")
	}

	coefLimit := ""
	if data.CoeffLimit != nil {
		coefLimit = strconv.Itoa(*data.CoeffLimit)
	}

	protection := ""
	if data.Protection != nil {
		protection = strconv.Itoa(*data.Protection)
	}

	if data.DraftID == uuid.Nil {
		return fmt.Sprintf(
			"Лимит даты бронирования: %s-%s\nID черновика:\nЗащита от бронирования: %s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
			data.FromDate.Format(dto.TimeFormat),
			data.ToDate.Format(dto.TimeFormat),
			protection,
			constmsg.WarehouseNames[data.Warehouse],
			coefLimit,
			data.SupplyType,
			constmsg.BotCommands[comName],
		), nil
	}

	return fmt.Sprintf(
		"Лимит даты бронирования: %s-%s\nID черновика: %s\nЗащита от бронирования: %s\nСклад: %s\nЛимит коэффициента: %s\nТип поставки: %s\n---------------\n%s",
		data.FromDate.Format(dto.TimeFormat),
		data.ToDate.Format(dto.TimeFormat),
		data.DraftID,
		protection,
		constmsg.WarehouseNames[data.Warehouse],
		coefLimit,
		data.SupplyType,
		constmsg.BotCommands[comName],
	), nil
}
