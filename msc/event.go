package msc

import "fmt"

// publish event
func publishEvent(event string, payload interface{}) error {
	subject := fmt.Sprintf("%s.%s", SubjectEvent, event)
	return broker.publish(subject, payload)
}

const (
	EventAccountAddEmail = "accountaddemail"
)
