package audit

import (
	"fmt"
	"net/mail"
)

type AccountAddEmailPayload struct {
	Userid string   `json:"userid"` // user id(name)
	Email  string   `json:"email"`  // new email address
	To     []string `json:"to"`     // recipient addresses
	Locale string   `json:"locale"` // i18n locale
}

// publish this audit message when account add new email
func AccountAddEmail(payload *AccountAddEmailPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is empty")
	}
	// validation to address
	if _, err := mail.ParseAddress(payload.Email); err != nil {
		return err
	}
	// send the request
	return publish(KAccountAddEmail, payload)
}
