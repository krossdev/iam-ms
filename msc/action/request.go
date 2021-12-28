package action

import (
	"fmt"

	"github.com/krossdev/iam-ms/msc"
)

const (
	KIPLocation      = "iplocation"
	KSendVerifyEmail = "sendverifyemail"
)

// send action request
func request(key string, payload interface{}) (interface{}, error) {
	b := msc.Borker()
	if b == nil {
		msc.Logger.Panicf("broker is nil")
	}
	subject := fmt.Sprintf("%s.%s", msc.SubjectAction, key)

	return b.Request(subject, payload)
}
