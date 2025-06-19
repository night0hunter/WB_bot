package handler

import (
	"encoding/json"
	"wb_bot/internal/dto"

	"github.com/pkg/errors"
)

type arg interface {
	dto.WarehouseData | dto.ChangeStatusInfo | []dto.WarehouseData
}

func Marshal[T arg](data T) ([]byte, error) {
	j, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}

	return j, nil
}

func Unmarshal[T arg](j []byte) (T, error) {
	var data T
	err := json.Unmarshal(j, &data)
	if err != nil {
		return data, errors.Wrap(err, "json.Unmarshal")
	}

	return data, nil
}
