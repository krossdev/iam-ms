// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
	"net/mail"
)

type AccountAddEmailAuditPayload struct {
	Userid string   `json:"userid"` // user id(name)
	Email  string   `json:"email"`  // new email address
	To     []string `json:"to"`     // recipient addresses
	Locale string   `json:"locale"` // i18n locale
}

// publish this audit message when account add new email
func AccountAddEmailAudit(payload *AccountAddEmailAuditPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is empty")
	}
	// validation to address
	if _, err := mail.ParseAddress(payload.Email); err != nil {
		return err
	}
	// send the request
	return publishAudit(AuditAccountAddEmail, payload)
}
