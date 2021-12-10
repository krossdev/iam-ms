// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"time"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/xlog"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func subscribeAction(subject string, handler nats.Handler) {
	conn.Subscribe(subject, actionHandler)
}

func actionHandler(subject, reply string, a *msc.RequestAction) {
	logger := xlog.X.WithFields(logrus.Fields{
		"subject": a.Subject,
		"version": a.Version,
		"reqid":   a.ReqId,
		"name":    a.Name,
	})
	logger.Infof("subject: %s, reply: %s, a: %+v", subject, reply, a)

	if err := checkVersion(a.Version); err != nil {
		logger.WithError(err).Error("version incompatable")
		replyActionWithError(reply, err, msc.ReplyCodeVersionIncompatible)
		return
	}
	logger.Info("ok")
	replyActionOk(reply)
}

// response with error
func replyActionWithError(reply string, err error, code int32) {
	conn.Publish(reply, &msc.ReplyAction{
		Version:   msc.Version,
		Timestamp: time.Now().UnixMicro(),
		Code:      code,
		Message:   err.Error(),
	})
}

// response with success
func replyActionOk(reply string) {
	conn.Publish(reply, &msc.ReplyAction{
		Version:   msc.Version,
		Timestamp: time.Now().UnixMicro(),
		Code:      msc.ReplyCodeOk,
	})
}
