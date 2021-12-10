// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import "net/mail"

// SendConsoleAction send action to console subject
func sendConsoleAction(name string, payload interface{}) (interface{}, error) {
	return sendActionRequest(&RequestAction{
		Subject: SubjectConsoleActions,
		Name:    name,
		Payload: payload,
	})
}

const (
	ActionNameResendVerifyEmail = "resend-verify-email"
)

type ResendVerifyEmailPayload struct {
	To string `json:"to"`
}

// ResendVerifyEmail ask to resend a verify email to account email address
func ResendVerifyEmail(to string) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return err
	}
	payload := &ResendVerifyEmailPayload{
		To: to,
	}
	_, err := sendConsoleAction(ActionNameResendVerifyEmail, payload)
	if err != nil {
		return err
	}
	return nil
}
