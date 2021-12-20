// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
)

type AccountAddEmailPayload struct {
	Subject string `json:"subject"` // mail subject
	Name    string `json:"name"`    // recipient name
	To      string `json:"to"`      // recipient address
	Href    string `json:"href"`    // verify url
	Locale  string `json:"locale"`  // i18n locale
	Expire  string `json:"expire"`  // expire
}

// Publish account add email audit
func AccountAddEmail(payload *AccountAddEmailPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is empty")
	}
	// // validation to address
	// if _, err := mail.ParseAddress(payload.To); err != nil {
	// 	return err
	// }
	// // validation href
	// u, err := url.Parse(payload.Href)
	// if err != nil {
	// 	return err
	// }
	// if !u.IsAbs() {
	// 	return fmt.Errorf("verify url invalid")
	// }
	// send the request
	return publishAudit(EventAccountAddEmail, payload)
}
