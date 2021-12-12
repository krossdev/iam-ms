// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"fmt"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/xlog"
	"github.com/sirupsen/logrus"
)

// subscribe to action subject
func subscribeAction() error {
	_, err := conn.Subscribe(msc.SubjectAction, actionHandler)
	return err
}

func actionHandler(subject, reply string, qt *msc.Request) {
	l := xlog.X.WithFields(logrus.Fields{
		"version": qt.Version,
		"time":    qt.Time,
		"reqid":   qt.ReqId,
		"action":  qt.Action,
	})
	l.Infof("request %s on %s, reply to = %s", qt.Action, subject, reply)

	// reply subject(inbox) cannot be empty
	if len(reply) == 0 {
		err := fmt.Errorf("missing reply subject")
		l.WithError(err).Error("bad request")
		return
	}
	// function to publish reply
	replyTo := func(code int32, message string, payload interface{}) {
		rp := msc.MakeReply(code, message, payload, qt.ReqId)
		if err := conn.Publish(reply, &rp); err != nil {
			l.WithError(err).Errorf("failed to reply to %s", qt.Action)
		}
	}
	if perr := checkRequest(qt); perr != nil {
		l.WithError(perr.err).Error("bad request")
		replyTo(perr.code, perr.Error(), nil)
		return
	}
	if len(qt.Action) == 0 {
		err := fmt.Errorf("request missing action")
		l.WithError(err).Error("bad request")
		replyTo(msc.ReplyNoAction, err.Error(), nil)
		return
	}
	l.Infof("request %s done!", qt.Action)

	replyTo(msc.ReplyOk, "", nil)
}
