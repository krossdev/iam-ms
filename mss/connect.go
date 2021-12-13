// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"time"

	"github.com/krossdev/iam-ms/mss/xlog"
	"github.com/nats-io/nats.go"
)

// json encoded connection
var conn *nats.EncodedConn

// reconnect to the message broker
func reconnect(servers []string) error {
	opts := nats.GetDefaultOptions()

	// connection options
	opts.Servers = servers
	opts.Name = "kiam-ms"
	opts.Timeout = 10 * time.Second

	// never give up reconnect once connected
	opts.MaxReconnect = -1

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
	// close the exists connection(if exist)
	disconnect()

	conn, err = nats.NewEncodedConn(c, nats.JSON_ENCODER)
	return err
}

// disconnect from message broker
func disconnect() {
	if conn != nil {
		if err := conn.Drain(); err != nil {
			xlog.X.WithError(err).Warn("drain error")
		}
		conn.Close()
	}
}
