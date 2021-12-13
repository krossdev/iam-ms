// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package email

import (
	"github.com/krossdev/iam-ms/xlog"
	"github.com/pkg/errors"
	"net/mail"
)

// Mail contains email information that needs to be sent
type Mail struct {
	Message *Message
	Schema  string
	Reply   bool
}

// Send an email with smtp protocol.
func Send(m *Mail) (err error) {
	var ps []*provider

	if len(m.Schema) == 0 {
		return errors.New("EMail.Schema can not be empty")
	}
	if m.Message == nil {
		return errors.New("EMail.Message can not be empty")
	}
	if ps, err = providersForSchema(m.Schema); err != nil {
		return
	}

	sent := false
	for _, p := range ps {
		var sender *mail.Address
		sender, err = mail.ParseAddress(p.Sender)
		if err != nil {
			xlog.F("schema", m.Schema, "provider", p.Name, "error", err).
				Errorf("Sender '%s' address incorrect", p.Sender)
			continue
		}
		if err = sendMail(p, sender, m.Message); err != nil {
			xlog.F("schema", m.Schema, "provider", p.Name, "error", err).
				Error("Failed to send mail with provider")
			continue
		}
		sent = true
		break
	}

	if !sent {
		err = errors.New("tried all the providers but no one can send the mail")
	}
	return err
}
