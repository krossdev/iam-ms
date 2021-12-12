// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package action

import (
	"github.com/sirupsen/logrus"
)

type Handler func(payload interface{}, l *logrus.Entry) (interface{}, error)

var handlers = map[string]Handler{}

func Find(action string) Handler {
	if f, ok := handlers[action]; ok {
		return f
	}
	return nil
}
