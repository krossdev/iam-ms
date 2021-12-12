// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Broker struct {
	servers []string
	conn    *nats.EncodedConn
}

var broker *Broker

// create unique broker
func createBroker(servers []string) {
	if broker != nil {
		broker.disconnect()
	}
	broker = &Broker{
		servers: servers,
	}
}

// connect to message server
func (b *Broker) connect() error {
	opts := nats.GetDefaultOptions()

	// connection options
	opts.Servers = b.servers
	opts.Name = "kiam"
	opts.Timeout = 10 * time.Second
	opts.InboxPrefix = "kiam._inbox"

	// never give up reconnect once connected
	opts.MaxReconnect = -1

	// listen connection events
	opts.AsyncErrorCB = asyncErrorHandler
	opts.DisconnectedErrCB = disconnectedErrorHandler
	opts.ReconnectedCB = reconnectedHandler
	opts.DiscoveredServersCB = discoveredServersHandler
	opts.ClosedCB = closedHandler

	// connect to message broker
	c, err := opts.Connect()
	if err != nil {
		return err
	}
	// make a json encoded connection to send/receive data
	b.conn, err = nats.NewEncodedConn(c, nats.JSON_ENCODER)
	return err
}

// disconnect from message server
func (b *Broker) disconnect() {
	if b.conn != nil && b.conn.Conn.IsConnected() {
		if err := b.conn.FlushTimeout(3 * time.Second); err != nil {
			logger.WithError(err).Warn("flush error")
		}
		b.conn.Close()
	}
}

// maxPayloadSize return maximum allowed payload size.
// this size limit is set by nats server, client cannot modify it.
// before send large message(like email attachment), client may needs to
// check if the message size is exceed the payload size limit.
func (b *Broker) maxPayloadSize() int64 {
	return b.conn.Conn.MaxPayload()
}

// request is a wrapper to nats.Conn.Reqeust
func (b *Broker) request(subject, action string, payload interface{}) (interface{}, error) {
	if b.conn == nil || !b.conn.Conn.IsConnected() {
		if err := b.connect(); err != nil {
			return nil, err
		}
	}
	qt := makeRequest(action, payload)

	l := logger.WithFields(logrus.Fields{
		"version": qt.Version,
		"time":    qt.Time,
		"reqid":   qt.ReqId,
		"action":  action,
	})
	l.Tracef("request %s on %s", action, subject)

	var rp Reply

	// send request and wait for first reply
	if err := b.conn.Request(subject, &qt, &rp, 10*time.Second); err != nil {
		return nil, errors.Wrap(err, "communication error")
	}
	if err := checkReply(&rp, qt.ReqId); err != nil {
		return nil, errors.Wrap(err, "proto error")
	}
	dt := rp.Time - qt.Time

	// if reply code is not ok, return err
	if rp.Code != ReplyOk {
		l.Errorf("request %s reply %d: %s, latency %.2fms",
			action, rp.Code, rp.Message, float64(dt/1000),
		)
		return nil, fmt.Errorf("%d-%s", rp.Code, rp.Message)
	}
	l.Tracef("request %s reply ok, latency %.2fms", action, float64(dt/1000))

	return rp.Payload, nil
}
