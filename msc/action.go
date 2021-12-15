// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
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
	Subject string `json:"subject"` // mail subject
	Name    string `json:"name"`    // recipient name
	To      string `json:"to"`      // recipient address
	Href    string `json:"href"`    // verify url
	Locale  string `json:"locale"`  // i18n locale
}

func SendVerifyEmail(payload *SendVerifyEmailPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is empty")
	}
	if _, err := mail.ParseAddress(payload.To); err != nil {
		return err
	}
	_, err := requestAction(ASendVerifyEmail, &payload)
	return err
}
