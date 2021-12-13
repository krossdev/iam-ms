// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package action

import (
	"html/template"
	"net/mail"
	"os"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/email"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	handlers[msc.ASendVerifyEmail] = sendVerifyEmail
}

// send-verify-email action
func sendVerifyEmail(payload interface{}, l *logrus.Entry) (interface{}, error) {
	var p msc.SendVerifyEmailPayload

	// convert map to struct
	if err := mapstructure.Decode(payload, &p); err != nil {
		return nil, err
	}
	// parse recipient address
	to, err := mail.ParseAddress(p.To)
	if err != nil {
		return nil, errors.Wrap(err, "email address invalid")
	}
	// generate mail body
	t, err := template.New("email").Parse("")
	if err != nil {
		return nil, errors.Wrap(err, "parse mail template error")
	}
	if err = t.Execute(os.Stdout, &p); err != nil {
		return nil, errors.Wrap(err, "generate mail content error")
	}

	// contract a mail message to send
	m := email.HTMLMessage("Please verify your email address", "")
	m.AddTO(to)

	// send mail
	if err = m.Send(); err != nil {
		l.WithError(err).Errorf("failed to send verify email to %s", to)
		return nil, err
	}
	return nil, nil
}
