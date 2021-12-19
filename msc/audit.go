package msc

import "fmt"

// publish audit
func publishAudit(audit string, payload interface{}) error {
	subject := fmt.Sprintf("%s.%s", SubjectAudit, audit)
	return broker.publish(subject, payload)
}

const (
	AuditAccountAddEmail = "accountaddemail"
)
