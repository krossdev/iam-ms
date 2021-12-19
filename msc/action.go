// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
)

// send action request
func requestAction(action string, payload interface{}) (interface{}, error) {
	subject := fmt.Sprintf("%s.%s", SubjectAction, action)
	return broker.request(subject, payload)
}

const (
	ActionSendVerifyEmail = "sendverifyemail"
	ActionIpLocation      = "iplocation"
)
