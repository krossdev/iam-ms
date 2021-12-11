// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"fmt"
	"time"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/xlog"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func subscribeActions(subject string, handler nats.Handler) error {
	_, err := conn.Subscribe(subject, actionHandler)
	return err
}

func actionHandler(subject, reply string, a *msc.RequestAction) {
	logger := xlog.X.WithFields(logrus.Fields{
		"version": a.Version,
		"time":    a.Time,
		"reqid":   a.ReqId,
		"name":    a.Name,
	})
	logger.Infof("receive action %s on %s, reply to %s", a.Name, subject, reply)

	// check request version
	if err := checkVersion(a.Version); err != nil {
		logger.WithError(err).Error("version %s incompatable", a.Version)
		replyActionWithError(a, reply, err, msc.ReplyCodeBadVersion)
		return
	}
	// check request time
	td := time.Since(time.UnixMicro(a.Time))
	if td < 0 || td > 10*time.Second {
		err := fmt.Errorf("request time %d invalid", a.Time)
		logger.WithError(err).Errorf("please check system clock")
		replyActionWithError(a, reply, err, msc.ReplyCodeBadTime)
	}
	logger.Infof("response ok to action %s", a.Name)
	replyActionOk(a, reply)
}

// response with error
func replyActionWithError(a *msc.RequestAction, reply string, err error, code int32) {
	conn.Publish(reply, &msc.ReplyAction{
		Version: msc.Version,
		Time:    time.Now().UnixMicro(),
		ReqId:   a.ReqId,
		Code:    code,
		Message: err.Error(),
	})
}

// response with success
func replyActionOk(a *msc.RequestAction, reply string) {
	conn.Publish(reply, &msc.ReplyAction{
		Version: msc.Version,
		Time:    time.Now().UnixMicro(),
		ReqId:   a.ReqId,
		Code:    msc.ReplyCodeOk,
	})
}
