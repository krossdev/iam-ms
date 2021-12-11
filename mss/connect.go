// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"time"

	"github.com/krossdev/iam-ms/mss/xlog"

	"github.com/nats-io/nats.go"
)

var conn *nats.EncodedConn

// reconnect to message broker
func reconnect(servers []string) error {
	xlog.X.Info("reconnect....")
	opts := nats.GetDefaultOptions()

	// connection options
	opts.Servers = servers
	opts.Name = "kiam-ms"
	opts.Timeout = 10 * time.Second
	opts.MaxReconnect = -1 // never give up reconnect

	// connection events
	opts.AsyncErrorCB = asyncErrorHandler
	opts.DisconnectedErrCB = disconnectedErrorHandler
	opts.ReconnectedCB = reconnectedHandler
	opts.DiscoveredServersCB = discoveredServersHandler
	opts.ClosedCB = closedHandler

	c, err := opts.Connect()
	if err != nil {
		return err
	}
	// disconnect brfore reconnect
	disconnect()

	conn, err = nats.NewEncodedConn(c, nats.JSON_ENCODER)
	return err
}

// disconnect from message broker
func disconnect() {
	if conn != nil {
		conn.Close()
	}
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

func discoveredServersHandler(c *nats.Conn) {
	xlog.X.Info("Message broker discover new server")
}

func closedHandler(c *nats.Conn) {
	xlog.X.Info("Message broker connection closed")
}
