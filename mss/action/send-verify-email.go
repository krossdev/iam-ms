// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package action

import (
	"github.com/krossdev/iam-ms/msc"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func init() {
	handlers[msc.ASendVerifyEmail] = sendVerifyEmail
}

// send-verify-email action
func sendVerifyEmail(payload interface{}, l *logrus.Entry) (interface{}, error) {
	var p msc.SendVerifyEmailPayload

	if err := mapstructure.Decode(payload, &p); err != nil {
		return nil, err
	}
	l.Infof("p=%+v", p)

	return nil, nil
}
