// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

var nc *nats.Conn

// Connect to message broker
func Connect(url string) error {
	conn, err := nats.Connect(url, nil)
	if err != nil {
		return errors.Wrap(err, "nats connect")
	}
	nc = conn
	return nil
}
