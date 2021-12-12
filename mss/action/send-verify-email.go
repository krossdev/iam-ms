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
	body := `
	<html>
		<head>
		</head>
		<body>
			<h1>REPLY</h1>
			<img width='64px' height='64px' src="cid:logo" alt='Logo' />
			<img src='https://miro.medium.com/max/1400/1*3LEMkgStgOhX1kSi08oZhQ.jpeg'
				width='32px' height='32px' />
			<p>Here is some content</p>
		</body>
	</html>
	`

	// contract a mail message to send
	m := email.HTMLMessage("Please verify your email address", body)
	m.AddTO(to)

	if err = m.Send(); err != nil {
		l.WithError(err).Errorf("failed to send verify email to %s", to)
		return nil, err
	}
	return nil, nil
}
