package msc

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type natsBroker struct {
	servers []string
	conn    *nats.EncodedConn
}

var broker *natsBroker

// Return broker
func Borker() *natsBroker {
	return broker
}

// createBroker create unique broker
func createBroker(servers []string) *natsBroker {
	if broker != nil {
		broker.disconnect()
	}
	broker = &natsBroker{
		servers: servers,
	}
	return broker
}

// connect to message server
func (b *natsBroker) connect() error {
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
func (b *natsBroker) disconnect() {
	if b.conn != nil && b.conn.Conn.IsConnected() {
		if err := b.conn.Drain(); err != nil {
			Logger.WithError(err).Warn("drain error")
		}
		b.conn.Close()
	}
}

// maxPayloadSize return maximum allowed payload size.
// this size limit is set by nats server, client cannot modify it.
// before send large message(like email attachment), client may needs to
// check if the message size is exceed the payload size limit.
func (b *natsBroker) maxPayloadSize() int64 {
	return b.conn.Conn.MaxPayload()
}

// Request is a wrapper to nats.Conn.Reqeust
func (b *natsBroker) Request(subject string, payload interface{}) (interface{}, error) {
	if b.conn == nil || !b.conn.Conn.IsConnected() {
		if err := b.connect(); err != nil {
			return nil, err
		}
	}
	qt := makeRequest(payload)

	l := Logger.WithFields(logrus.Fields{
		"version": qt.Version,
		"time":    qt.Time,
		"reqid":   qt.ReqId,
		"subject": subject,
	})
	l.Tracef("send request to %s", subject)

	var rp Reply

	// send request and wait for first reply
	if err := b.conn.Request(subject, &qt, &rp, 60*time.Second); err != nil {
		return nil, errors.Wrap(err, "communication error")
	}
	if err := checkReply(&rp, qt.ReqId); err != nil {
		return nil, errors.Wrap(err, "proto error")
	}
	dt := rp.Time - qt.Time

	// if reply code is not ok, return err
	if rp.Code != ReplyOk {
		l.Errorf("%s reply %d: %s, latency %.2fms",
			subject, rp.Code, rp.Message, float64(dt/1000),
		)
		return nil, fmt.Errorf("%s", rp.Message)
	}
	l.Tracef("%s reply ok, latency %.2fms", subject, float64(dt/1000))

	return rp.Payload, nil
}

// Publish is a wrapper to nats.Conn.Publish
func (b *natsBroker) Publish(subject string, payload interface{}) error {
	if b.conn == nil || !b.conn.Conn.IsConnected() {
		if err := b.connect(); err != nil {
			return err
		}
	}
	qt := makeRequest(payload)

	l := Logger.WithFields(logrus.Fields{
		"version": qt.Version,
		"time":    qt.Time,
		"reqid":   qt.ReqId,
		"subject": subject,
	})
	l.Tracef("publish message to %s", subject)

	return b.conn.Publish(subject, &qt)
}
