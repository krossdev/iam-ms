// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package mscli

import "github.com/pkg/errors"

// Mail contains email information that needs to be sent
type Mail struct {
	Message *Message
	Schema  string
	Reply   bool // 是否需要确认送达
	//BeforeHook func(*Mail) error `json:"-"`
	AfterHook func(*Mail, error) error `json:"-"`
}

// mailC is mail channel
var mailC = make(chan *Mail, 100)

// SendMail puts the mail in the sending queue and waits for dispatch.
func SendMail(m *Mail) error {
	if m == nil {
		return errors.Errorf("no mail to send, parameter is nil")
	}
	mailC <- m
	return nil
}
