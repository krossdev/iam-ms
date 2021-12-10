// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"time"

	"github.com/krossdev/iam-ms/mss/xlog"
	"github.com/nats-io/nats.go"
)

var conn *nats.Conn

// connect to message broker
func connect(servers []string) error {
	opts := nats.GetDefaultOptions()

	opts.Servers = servers
	opts.Name = "kiam-ms"
	opts.Timeout = 10 * time.Second
	opts.PingInterval = 10 * time.Second
	opts.MaxPingsOut = 3

	opts.AsyncErrorCB = asyncErrorHandler
	opts.DisconnectedErrCB = disconnectedErrorHandler
	opts.ReconnectedCB = reconnectedHandler
	opts.ClosedCB = closedHandler

	nc, err := opts.Connect()
	if err != nil {
		return err
	}
	conn = nc
	return nil
}

// nats async error callback
func asyncErrorHandler(c *nats.Conn, s *nats.Subscription, e error) {
	xlog.X.WithError(e).Error("Message broker async error")
}

func disconnectedErrorHandler(c *nats.Conn, e error) {
	xlog.X.WithError(e).Warn("Message broker disconnected")
}

func reconnectedHandler(c *nats.Conn) {
	xlog.X.Info("Message broker reconnected")
}

func closedHandler(c *nats.Conn) {
	xlog.X.Info("Message broker connection closed")
}
