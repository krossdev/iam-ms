// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// connection events listener
var (
	AsyncErrorHandler        nats.ErrHandler
	DisconnectedErrorHandler nats.ConnErrHandler
	ReconnectedHandler       nats.ConnHandler
	DiscoveredServersHandler nats.ConnHandler
	ClosedHandler            nats.ConnHandler
)

// encoded connection
var conn *nats.EncodedConn

var logger = logrus.StandardLogger().WithField("realm", "ms")

// set logger
func SetLogger(e *logrus.Entry) {
	logger = e
}

// Connect to message broker
func Connect(servers []string) error {
	if len(servers) == 0 {
		logger.Panic("no servers give to connect")
	}
	opts := nats.GetDefaultOptions()

	// connection options
	opts.Servers = servers
	opts.Name = "kiam"
	opts.Timeout = 10 * time.Second
	opts.MaxReconnect = -1 // never give up reconnect
	opts.InboxPrefix = "REPLY-TO"

	// listen connection events
	if AsyncErrorHandler != nil {
		opts.AsyncErrorCB = AsyncErrorHandler
	}
	if DisconnectedErrorHandler != nil {
		opts.DisconnectedErrCB = DisconnectedErrorHandler
	}
	if ReconnectedHandler != nil {
		opts.ReconnectedCB = ReconnectedHandler
	}
	if DiscoveredServersHandler != nil {
		opts.DiscoveredServersCB = DiscoveredServersHandler
	}
	if ClosedHandler != nil {
		opts.ClosedCB = ClosedHandler
	}
	// connect to broker
	c, err := opts.Connect()
	if err != nil {
		return err
	}
	// make a json encoded connection to send/receive data
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
