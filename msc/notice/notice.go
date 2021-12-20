package notice

import (
	"fmt"

	"github.com/krossdev/iam-ms/msc"
)

const (
	KAccountAddEmail = "accountaddemail"
)

// publish notice message
func publish(key string, payload interface{}) error {
	b := msc.Borker()
	if b == nil {
		msc.Logger.Panicf("broker is nil")
	}
	subject := fmt.Sprintf("%s.%s", msc.SubjectNotice, key)
	return b.Publish(subject, payload)
}
