// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/xlog"
	"github.com/nats-io/nats.go"
)

var consoleHanders = map[string]nats.Handler{
	msc.ActionNameResendVerifyEmail: nil,
}

func init() {

}

func consoleHandler(a *msc.RequestAction) {
	xlog.X.Infof("recv action: %v", a)
}

func subscribeConsoleActions() {
	subscribeAction(msc.SubjectConsoleActions, consoleHandler)
}
