package audit

import (
	"fmt"
	"net/mail"
)

type AccountAddEmailPayload struct {
	NewEmail string   `json:"newemail"` // new email address
	Userid   string   `json:"userid"`   // user id(name)
	HttpURL  string   `json:"httpurl"`  // KrossIAM host url
	Locale   string   `json:"locale"`   // i18n locale
	MailTo   []string `json:"mailto"`   // send mail to
	Subject  string   `json:"subject"`  // mail subject
}

// publish this audit message when account add new email
func AccountAddEmail(payload *AccountAddEmailPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is empty")
	}
	// validation new email address
	if _, err := mail.ParseAddress(payload.NewEmail); err != nil {
		return err
	}
	// send the request
	return publish(KAccountAddEmail, payload)
}
