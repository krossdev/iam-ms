package audit

import (
	"fmt"

	"github.com/krossdev/iam-ms/msc"
)

const (
	KAccountAddEmail = "accountaddemail"
)

// publish audit message
func publish(key string, payload interface{}) error {
	b := msc.Borker()
	if b == nil {
		msc.Logger.Panicf("broker is nil")
	}
	subject := fmt.Sprintf("%s.%s", msc.SubjectAudit, key)

	return b.Publish(subject, payload)
}
