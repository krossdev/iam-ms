// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
	"net/mail"
	"net/url"
)

const (
	SubjectAction = "kiam.console.action"
	SubjectEvent  = "kiam.console.event"
	SubjectAudit  = "kiam.console.audit"
)

// send request to action subject
func requestAction(action string, payload interface{}) (interface{}, error) {
	subject := fmt.Sprintf("%s.%s", SubjectAction, action)
	return broker.request(subject, payload)
}

const (
	ActionSendVerifyEmail = "send-verify-email"
)

type SendVerifyEmailPayload struct {
	Subject string `json:"subject"` // mail subject
	Name    string `json:"name"`    // recipient name
	To      string `json:"to"`      // recipient address
	Href    string `json:"href"`    // verify url
	Locale  string `json:"locale"`  // i18n locale
	Expire  string `json:"expire"`  // expire
}

// Ask message services to send a verify email
func SendVerifyEmail(payload *SendVerifyEmailPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is empty")
	}
	// validation to address
	if _, err := mail.ParseAddress(payload.To); err != nil {
		return err
	}
	// validation href
	u, err := url.Parse(payload.Href)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return fmt.Errorf("verify url invalid")
	}
	// send the request
	_, err = requestAction(ActionSendVerifyEmail, payload)
	return err
}
