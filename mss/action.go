// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"fmt"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/action"
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/xlog"

	"github.com/sirupsen/logrus"
)

type ActionHandlerFunc func(payload interface{}, conf interface{}, l *logrus.Entry) (interface{}, error)

// subscribe to action subject
func subscribeAction(action string, handler ActionHandlerFunc, conf interface{}) error {
	actionHandler := func(subject string, reply string, qt *msc.Request) {
		logger := xlog.X.WithFields(logrus.Fields{
			"version": qt.Version,
			"time":    qt.Time,
			"reqid":   qt.ReqId,
			"subject": subject,
		})
		logger.Infof("receive action request on %s", subject)

		// reply subject(inbox) cannot be empty
		if len(reply) == 0 {
			err := fmt.Errorf("missing reply subject")
			logger.WithError(err).Error("action request missing reply subject")
			return
		}
		// function to reply message
		replyTo := func(code int32, message string, payload interface{}) {
			rp := msc.MakeReply(code, message, payload, qt.ReqId)
			if err := conn.Publish(reply, &rp); err != nil {
				logger.WithError(err).Errorf("failed to reply to %s", subject)
			}
		}
		// check request
		if perr := checkRequest(qt); perr != nil {
			logger.WithError(perr.err).Error("action request validation error")
			replyTo(perr.code, perr.Error(), nil)
			return
		}
		// call handler
		payload, err := handler(qt.Payload, conf, logger)
		if err != nil {
			logger.WithError(err).Errorf("action %s execute error", subject)
			replyTo(msc.ReplyError, err.Error(), nil)
			return
		}
		logger.Infof("%s is done! reply to %s", subject, reply)

		// must reply a message
		replyTo(msc.ReplyOk, "ok", payload)
	}
	subject := fmt.Sprintf("%s.%s", msc.SubjectAction, action)

	xlog.X.Tracef("subscribed action message on %s ...", subject)

	_, err := conn.Subscribe(subject, actionHandler)
	return err
}

// subscribe actions
func subscribeActionsWithConfig(c *config.ServiceActions) error {
	handlers := map[string]ActionHandlerFunc{
		msc.ActionIpLocation:      action.IPLocationHandler,
		msc.ActionSendVerifyEmail: action.SendVerifyEmailHandler,
	}
	subscribe := func(action string, conf interface{}) error {
		return subscribeAction(action, handlers[action], conf)
	}
	// scbscribe action which enabled
	if c.IPLocation.Subscribe {
		if err := subscribe(msc.ActionIpLocation, &c.IPLocation); err != nil {
			return err
		}
	}
	if c.SendVerifyEmail.Subscribe {
		if err := subscribe(msc.ActionSendVerifyEmail, &c.SendVerifyEmail); err != nil {
			return err
		}
	}
	return nil
}
