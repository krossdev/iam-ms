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
	l := xlog.X.WithFields(logrus.Fields{
		"version": a.Version,
		"time":    a.Time,
		"reqid":   a.ReqId,
		"name":    a.Name,
	})
	l.Infof("receive action %s on %s, reply to = %s", a.Name, subject, reply)

	// reply subject(inbox) cannot be empty
	if len(reply) == 0 {
		err := fmt.Errorf("missing reply subject")
		l.WithError(err).Error("request invalid")
		replyActionWithError(a.ReqId, reply, err, msc.ReplyCodeNoReply)
	}
	if err := checkRequestAction(a); err != nil {
		l.WithError(err).Error("request invalid")
		replyActionWithError(a.ReqId, reply, err, msc.ReplyCodeNoReply)
	}
	replyActionWithOk(a.ReqId, reply)

	l.Infof("action %s done!", a.Name)
}

// check request action
func checkRequestAction(a *msc.RequestAction) error {
	// check request version
	if err := checkVersion(a.Version); err != nil {
		return err
	}
	// check request time
	td := time.Since(time.UnixMicro(a.Time))
	if td < 0 || td > 10*time.Second {
		return fmt.Errorf("request time %d invalid", a.Time)
	}
	return nil
}

// response with error
func replyActionWithError(reqid string, reply string, err error, code int32) {
	conn.Publish(reply, &msc.ReplyAction{
		Version: msc.Version,
		Time:    time.Now().UnixMicro(),
		ReqId:   reqid,
		Code:    code,
		Message: err.Error(),
	})
}

// response with success
func replyActionWithOk(reqid string, reply string) {
	conn.Publish(reply, &msc.ReplyAction{
		Version: msc.Version,
		Time:    time.Now().UnixMicro(),
		ReqId:   reqid,
		Code:    msc.ReplyCodeOk,
	})
}
