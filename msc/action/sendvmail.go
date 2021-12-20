package action

import (
	"fmt"
	"net/mail"
	"net/url"

	"github.com/pkg/errors"
)

type SendVerifyEmailPayload struct {
	Subject string `json:"subject"` // mail subject
	Userid  string `json:"userid"`  // recipient name
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
		return errors.Wrapf(err, "parse ip(%s) error", payload.To)
	}
	// validation href
	u, err := url.Parse(payload.Href)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return fmt.Errorf("verify url(%s) invalid", payload.Href)
	}
	// send the request
	_, err = request(KSendVerifyEmail, payload)
	return err
}
