// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"fmt"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/audit"
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/xlog"

	"github.com/sirupsen/logrus"
)

// subscribe to audit messages
func subscribeAudit(conf *config.ServiceAudits) error {
	auditHandler := func(subject string, reply string, qt *msc.Request) {
		logger := xlog.X.WithFields(logrus.Fields{
			"version": qt.Version,
			"time":    qt.Time,
			"reqid":   qt.ReqId,
			"subject": subject,
		})
		logger.Infof("receive audit message on %s", subject)

		// check request
		if perr := checkRequest(qt); perr != nil {
			logger.WithError(perr.err).Error("validate audit request error")
			return
		}
		// dispatch
		if err := audit.Handler(subject, qt.Payload, conf, logger); err != nil {
			logger.WithError(err).Errorf("failed to process audit message %s", subject)
		}
		logger.Infof("%s is done!", subject)
	}
	// subscribe all audit messages with wildcard
	subject := fmt.Sprintf("%s.>", msc.SubjectAudit)

	xlog.X.Tracef("subscribed audit messages on %s ...", subject)

	_, err := conn.Subscribe(subject, auditHandler)
	return err
}

// subscribe audit messages
func subscribeAuditsWithConfig(c *config.ServiceAudits) error {
	if !c.Subscribe {
		return nil
	}
	return subscribeAudit(c)
}
