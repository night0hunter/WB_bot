package handler

import (
	"context"
	"fmt"
	"wb_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (h *handler) TrackingCron(ctx context.Context) error {
	var msg tgbotapi.MessageConfig

	trackings, err := h.service.GetTrackings(ctx)
	if err != nil {
		return errors.Wrap(err, "service.TrackingCron")
	}

	if len(trackings) == 0 {
		return nil
	}

	for _, tr := range trackings {
		if !tr.IsActive {
			continue
		}

		msg = tgbotapi.NewMessage(tr.UserID, fmt.Sprintf("Найден таймслот\nСклад: %s\nКоэффициент: %dх\nТип поставки: %s\nДата: %s",
			tr.WarehouseName,
			tr.Coefficient,
			tr.BoxTypeName,
			tr.Date.Format(utils.TimeFormat),
		))

		if _, err := h.bot.Send(msg); err != nil {
			return errors.Wrap(err, "bot.Send")
		}

		err = h.service.KeepSendingTime(ctx, tr)
		if err != nil {
			return errors.Wrap(err, "service.KeepSendingTime")
		}

		return nil
	}

	return nil
}
