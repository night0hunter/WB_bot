package cronjob

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type BookDraftCronHandler interface {
	BookDraftCron(ctx context.Context) error
}

type BookDraftCron struct {
	handler BookDraftCronHandler
}

func NewBookDraftCron(handler BookDraftCronHandler) *BookDraftCron {
	return &BookDraftCron{
		handler: handler,
	}
}

func (c *BookDraftCron) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.handler.BookDraftCron(ctx)
	if err != nil {
		fmt.Println(errors.Wrap(err, "handler.BookDraftCron"))
	}
}
