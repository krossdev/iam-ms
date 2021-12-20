// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package actions

import (
	"html/template"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/krossdev/iam-ms/msc/action"
	"github.com/krossdev/iam-ms/mss/email"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// data pass to template execute
type sendVerifyEmailTemplateData struct {
	action.SendVerifyEmailPayload
	TemplateData email.TemplateData
}

// send-verify-email action handler
func SendVerifyEmailHandler(p interface{}, c interface{}, l *logrus.Entry) (interface{}, error) {
	var payload action.SendVerifyEmailPayload

	// convert map to struct
	if err := mapstructure.Decode(p, &payload); err != nil {
		return nil, err
	}
	// parse recipient address
	to, err := mail.ParseAddress(payload.To)
	if err != nil {
		return nil, errors.Wrap(err, "email address invalid")
	}
	// inline logo cid identifer
	logo_cid := strings.ReplaceAll(uuid.NewString(), "-", "")

	templateData := email.TemplateData{
		Logo:  template.URL("cid:" + logo_cid),
		Title: payload.Subject,
	}
	// data to execute template
	data := sendVerifyEmailTemplateData{payload, templateData}

	// generate mail content
	content, err := email.ExecTemplate(email.TVerifyEmail, data.Locale, &data)
	if err != nil {
		return nil, err
	}
	// contract a mail message
	m := email.HTMLMessage(data.Subject, content)
	m.AddTO(to)

	// inline logo
	m.Inline(email.TemplateLogoPath(), logo_cid)

	// send mail
	if err = m.Send(); err != nil {
		return nil, errors.Wrapf(err, "failed to send verify email to %s", to)
	}
	return nil, nil
}
