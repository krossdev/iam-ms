package msc

import "fmt"

// publish audit message
func publishAudit(audit string, payload interface{}) error {
	subject := fmt.Sprintf("%s.%s", SubjectAudit, audit)
	return broker.publish(subject, payload)
}

const (
	AuditAccountAddEmail = "accountaddemail"
)
