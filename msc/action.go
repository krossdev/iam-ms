// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"net/mail"
)

const (
	SubjectAction = "kiam.console.action"
	SubjectEvent  = "kiam.console.event"
	SubjectAudit  = "kiam.console.audit"
)

// send request to action subject
func requestAction(action string, payload interface{}) (interface{}, error) {
	return broker.request(SubjectAction, action, payload)
}

const (
	ASendVerifyEmail = "send-verify-email"
)

type SendVerifyEmailPayload struct {
	To string `json:"to"`
}

func SendVerifyEmail(to string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return err
	}
	payload := SendVerifyEmailPayload{
		To: to,
	}
	_, err := requestAction(ASendVerifyEmail, &payload)
	return err
}
