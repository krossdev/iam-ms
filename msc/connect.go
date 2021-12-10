// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var (
	asyncErrorHandler        nats.ErrHandler
	disconnectedErrorHandler nats.ConnErrHandler
	reconnectedHandler       nats.ConnHandler
	closedHandler            nats.ConnHandler
)

var conn *nats.EncodedConn
var logger *logrus.Logger

// set logger
func SetLogger(l *logrus.Logger) {
	logger = l
}

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
	opts.InboxPrefix = "REPLY-TO"

	opts.AsyncErrorCB = asyncErrorHandler
	opts.DisconnectedErrCB = disconnectedErrorHandler
	opts.ReconnectedCB = reconnectedHandler
	opts.ClosedCB = closedHandler

	c, err := opts.Connect()
	if err != nil {
		return err
	}
	conn, err = nats.NewEncodedConn(c, nats.JSON_ENCODER)
	return err
}

// MaxPayloadSize return maximum allowed payload size.
// this size limit is set by nats server, client cannot modify it.
// before send large message(like email attachment), client may needs to
// check if the message size is exceed the payload size limit.
func MaxPayloadSize() int64 {
	return conn.Conn.MaxPayload()
}

// Disconnect from message broker
func Disconnect() {
	if conn != nil {
		conn.FlushTimeout(3 * time.Second)
		conn.Close()
	}
}
