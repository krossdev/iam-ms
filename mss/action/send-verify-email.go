// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package action

import (
	"html/template"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/email"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	handlers[msc.ASendVerifyEmail] = sendVerifyEmail
}

type SendVerifyEmailTemplateData struct {
	msc.SendVerifyEmailPayload
	TemplateData email.TemplateData
}

// send-verify-email action
func sendVerifyEmail(p interface{}, l *logrus.Entry) (interface{}, error) {
	var payload msc.SendVerifyEmailPayload

	// convert map to struct
	if err := mapstructure.Decode(p, &payload); err != nil {
		return nil, err
	}
	// parse recipient address
	to, err := mail.ParseAddress(payload.To)
	if err != nil {
		return nil, errors.Wrap(err, "email address invalid")
	}
	subject := "Please verify your email address"

	logo_cid := strings.ReplaceAll(uuid.NewString(), "-", "")

	templateData := email.TemplateData{
		Logo:  template.URL("cid:" + logo_cid),
		Title: subject,
	}
	// data to execute template
	execData := SendVerifyEmailTemplateData{payload, templateData}

	// generate mail content
	content, err := email.ExecTemplate("verify-email", &execData)
	if err != nil {
		return nil, err
	}
	// contract a mail message
	m := email.HTMLMessage(subject, content)
	m.AddTO(to)

	// inline logo
	m.Inline(email.LogoPath(), logo_cid)

	// send mail
	if err = m.Send(); err != nil {
		l.WithError(err).Errorf("failed to send verify email to %s", to)
		return nil, err
	}
	return nil, nil
}
