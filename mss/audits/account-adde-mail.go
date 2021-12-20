package audits

import (
	"html/template"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/krossdev/iam-ms/msc/audit"
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/email"
	"github.com/pkg/errors"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func init() {
	handlers[audit.KAccountAddEmail] = accountAddEmail
}

type accountAddEmailTemplateData struct {
	audit.AccountAddEmailPayload
	TemplateData email.TemplateData
}

func accountAddEmail(p interface{}, c *config.ServiceAudits, l *logrus.Entry) error {
	var payload audit.AccountAddEmailPayload

	// convert map to struct
	if err := mapstructure.Decode(p, &payload); err != nil {
		return err
	}

	l.Debugf("payload %v", payload)

	// inline logo cid identifer
	logo_cid := strings.ReplaceAll(uuid.NewString(), "-", "")

	templateData := email.TemplateData{
		Logo:  template.URL("cid:" + logo_cid),
		Title: payload.Subject,
	}
	// data to execute template
	data := accountAddEmailTemplateData{payload, templateData}

	// generate mail content
	content, err := email.ExecTemplate(email.TAccountAddEmail, data.Locale, &data)
	if err != nil {
		return err
	}
	// contract a mail message
	m := email.HTMLMessage(data.Subject, content)

	for _, t := range payload.MailTo {
		if to, err := mail.ParseAddress(t); err != nil {
			l.WithError(err).Errorf("parse mail address(%s) error", t)
		} else {
			m.AddTO(to)
		}
	}
	// inline logo
	m.Inline(email.TemplateLogoPath(), logo_cid)

	// send mail
	if err = m.Send(); err != nil {
		return errors.Wrapf(err,
			"failed to send email to %s", strings.Join(payload.MailTo, ","),
		)
	}
	return nil
}
