package msc

import "fmt"

// publish notice message
func publishNotice(event string, payload interface{}) error {
	subject := fmt.Sprintf("%s.%s", SubjectNotice, event)
	return broker.publish(subject, payload)
}
