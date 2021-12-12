// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import "github.com/sirupsen/logrus"

// the default logger has no any configuration(like log to file...),
// app should replace it with SetLogger()
var logger = logrus.StandardLogger().WithField("realm", "ms")

// set a comfortable logger
func SetLogger(e *logrus.Entry) {
	logger = e
}
