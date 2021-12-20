// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"fmt"
	"strings"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/audit"
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/xlog"

	"github.com/sirupsen/logrus"
)

type AuditHandlerFunc func(payload interface{}, params interface{}, l *logrus.Entry) error

// subscribe to action subject
func subscribeAudit(params interface{}) error {
	auditHandler := func(subject string, reply string, qt *msc.Request) {
		logger := xlog.X.WithFields(logrus.Fields{
			"version": qt.Version,
			"time":    qt.Time,
			"reqid":   qt.ReqId,
			"subject": subject,
		})
		logger.Infof("receive audit on %s", subject)

		// check request
		if perr := checkRequest(qt); perr != nil {
			logger.WithError(perr.err).Error("audit request validation error")
			return
		}
		// call handler
		arr := strings.Split(subject, ".")
		err := audit.Handler(arr[len(arr)-1], qt.Payload, params, logger)
		// err := handler(qt.Payload, params, logger)
		if err != nil {
			logger.WithError(err).Errorf("audit %s execute error", subject)
			return
		}
		logger.Infof("%s is done!")
	}
	subject := fmt.Sprintf("%s.>", msc.SubjectAudit)

	xlog.X.Tracef("subscribed on %s ...", subject)

	_, err := conn.Subscribe(subject, auditHandler)
	return err
}

// subscribe audit
func subscribeAuditWithConfig(c *config.ServiceAudit) error {
	if c.Subscribe {
		return subscribeAudit(c)
	}
	return nil
}
