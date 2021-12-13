// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package action

import (
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/email"
	"github.com/krossdev/iam-ms/mss/xlog"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	handlers[msc.ASendVerifyEmail] = sendVerifyEmail
}

type SendVerifyEmailTemplateData struct {
	msc.SendVerifyEmailPayload
	Logo string
}

// send-verify-email action
func sendVerifyEmail(p interface{}, l *logrus.Entry) (interface{}, error) {
	var payload msc.SendVerifyEmailPayload

	xlog.X.Infof("payload: %v", p)
	// convert map to struct
	if err := mapstructure.Decode(p, &payload); err != nil {
		return nil, err
	}
	xlog.X.Infof("data: %v", payload)

	// parse recipient address
	to, err := mail.ParseAddress(payload.To)
	if err != nil {
		return nil, errors.Wrap(err, "email address invalid")
	}
	// template data add a logo field
	tdata := SendVerifyEmailTemplateData{
		payload,
		strings.ReplaceAll(uuid.NewString(), "-", ""),
	}
	// generate mail content
	content, err := email.ExecTemplate("verify-email", &tdata)
	if err != nil {
		return nil, err
	}
	// contract a mail message
	m := email.HTMLMessage("Please verify your email address", content)
	m.AddTO(to)

	// inline logo
	m.Inline(email.LogoPath(), tdata.Logo)

	// send mail
	if err = m.Send(); err != nil {
		l.WithError(err).Errorf("failed to send verify email to %s", to)
		return nil, err
	}
	return nil, nil
}
