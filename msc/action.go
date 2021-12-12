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
	Name string `json:"username"` // user name
	To   string `json:"to"`       // recipient address
	Href string `json:"href"`     // verify url
}

func SendVerifyEmail(name, to, href string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return err
	}
	payload := SendVerifyEmailPayload{
		Name: name,
		To:   to,
		Href: href,
	}
	_, err := requestAction(ASendVerifyEmail, &payload)
	return err
}
