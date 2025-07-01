package addHandler

import (
	"fmt"
	constmsg "wb_bot/internal/const_message"
	"wb_bot/internal/dto"
	"wb_bot/internal/enum"
	"wb_bot/internal/utils"

	"github.com/pkg/errors"
)

func CraftMessage(prevCommand dto.PrevCommandInfo, comName enum.CommandSequence) (string, error) {
	data, err := utils.Unmarshal[dto.WarehouseData](prevCommand.Info)
	if err != nil {
		return "", errors.Wrap(err, "utils.Unmarshal")
	}

	if data.CoeffLimit == nil {
		return fmt.Sprintf(
			"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента:\nТип поставки: %s\n---------------\n%s",
			data.FromDate.Format(dto.TimeFormat),
			data.ToDate.Format(dto.TimeFormat),
			constmsg.WarehouseNames[data.Warehouse],
			data.SupplyType,
			constmsg.BotCommands[comName],
		), nil
	}

	return fmt.Sprintf(
		"Дата отслеживания: %s-%s\nСклад: %s\nЛимит коэффициента: %d\nТип поставки: %s\n---------------\n%s",
		data.FromDate.Format(dto.TimeFormat),
		data.ToDate.Format(dto.TimeFormat),
		constmsg.WarehouseNames[data.Warehouse],
		data.CoeffLimit,
		data.SupplyType,
		constmsg.BotCommands[comName],
	), nil
}
