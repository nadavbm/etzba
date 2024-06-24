package prompusher

import (
	"github.com/nadavbm/zlog"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Client interface {
}

func NewClient(logger *zlog.Logger, pushUrl, jobName string) Client {
	pusher := push.New(pushUrl, jobName)
	return &client{
		logger: logger,
		Pusher: pusher,
	}
}

type client struct {
	logger *zlog.Logger
	Pusher *push.Pusher
}

func (c *client) PushMetrics() error {
	return nil
}
