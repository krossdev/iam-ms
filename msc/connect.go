// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

// last connnection
var nc *nats.Conn

// Connect to message broker
func Connect(url string) (err error) {
	nc, err = nats.Connect(url, nil)
	if err != nil {
		return errors.Wrap(err, "nats connect")
	}
	return nil
}

// Disconnect from message broker
func Disconnect() {
	if nc != nil {
		nc.Flush()
		nc.Close()
	}
}

func Subscript() {
	nc.Subscribe("", func(msg *nats.Msg) {})
}
