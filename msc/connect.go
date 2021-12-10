// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"time"

	"github.com/nats-io/nats.go"
)

var (
	asyncErrorHandler        nats.ErrHandler
	disconnectedErrorHandler nats.ConnErrHandler
	reconnectedHandler       nats.ConnHandler
	closedHandler            nats.ConnHandler
)

var conn *nats.Conn

// Must be call before Connect
func SetAsyncErrorHandler(handler nats.ErrHandler) {
	asyncErrorHandler = handler
}

func SetDisconnectErrHandler(handler nats.ConnErrHandler) {
	disconnectedErrorHandler = handler
}

func SetReconnectHandler(handler nats.ConnHandler) {
	reconnectedHandler = handler
}

func SetClosedHandler(handler nats.ConnHandler) {
	closedHandler = handler
}

// Connect to message broker
func Connect(servers []string) error {
	opts := nats.GetDefaultOptions()

	opts.Servers = servers
	opts.Name = "kiam"
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

func MaxPayloadSize() int64 {
	return conn.MaxPayload()
}

// Disconnect from message broker
func Disconnect() {
	if conn != nil {
		conn.Flush()
		conn.Close()
	}
}
