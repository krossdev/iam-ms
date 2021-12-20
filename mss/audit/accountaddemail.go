package audit

import (
	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/config"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func init() {
	handlers[msc.AuditAccountAddEmail] = accountAddEmail
}

func accountAddEmail(p interface{}, c *config.ServiceAudits, l *logrus.Entry) error {
	var payload msc.AccountAddEmailAuditPayload

	// convert map to struct
	if err := mapstructure.Decode(p, &payload); err != nil {
		return err
	}

	l.Debugf("payload %v", payload)
	return nil
}
