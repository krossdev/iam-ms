// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package action

import (
	"net/mail"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/email"
	"github.com/pkg/errors"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func init() {
	handlers[msc.ASendVerifyEmail] = sendVerifyEmail
}

// send-verify-email action
func sendVerifyEmail(payload interface{}, l *logrus.Entry) (interface{}, error) {
	var p msc.SendVerifyEmailPayload

	if err := mapstructure.Decode(payload, &p); err != nil {
		return nil, err
	}

	to, err := mail.ParseAddress(p.To)
	if err != nil {
		return nil, errors.Wrap(err, "email address invalid")
	}
	m := email.HTMLMessage("Please verify your email address", "please verify")
	m.AddTO(to)
	err = m.Send()
	return nil, err
}
