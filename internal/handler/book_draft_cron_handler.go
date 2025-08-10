package handler

import (
	"context"
	"fmt"
	constmsg "wb_bot/internal/const_message"
	myError "wb_bot/internal/error"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (h *handler) BookDraftCron(ctx context.Context) error {
	id, err := h.service.BookDraftCron(ctx)
	if err != nil {
		var myerr *myError.MyError
		if errors.As(err, &myerr) {
			fmt.Println(err)

			msg := tgbotapi.NewMessage(id, constmsg.MatchErrorType[myerr.GetErrorType()])
			_, err := h.bot.Send(msg)
			if err != nil {
				return errors.Wrap(err, "bot.Send")
			}

			return nil
		}

		return errors.Wrap(err, "service.BookDraftCron")
	}

	return nil
}
