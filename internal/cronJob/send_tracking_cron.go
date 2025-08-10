package cronjob

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type SendTrackingsCronHandler interface {
	TrackingCron(ctx context.Context) error
}

type SendTrackingsCron struct {
	handler SendTrackingsCronHandler
}

func NewSendTrackingCron(handler SendTrackingsCronHandler) *SendTrackingsCron {
	return &SendTrackingsCron{
		handler: handler,
	}
}

func (c *SendTrackingsCron) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.handler.TrackingCron(ctx)
	if err != nil {
		fmt.Println(errors.Wrap(err, "svc.TrackingCron"))
	}
}
