package audit

import (
	"fmt"
	"strings"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/sirupsen/logrus"
)

// handler
type handlerFunc func(payload interface{}, conf *config.ServiceAudits, l *logrus.Entry) error

// all registered handlers
var handlers = map[string]handlerFunc{}

// find a matched handler to process the message
func Handler(subject string, payload interface{}, conf *config.ServiceAudits, l *logrus.Entry) error {
	if !strings.HasPrefix(subject, msc.SubjectAudit) {
		return fmt.Errorf("subject '%s' is not a audit message", subject)
	}
	if len(subject) <= len(msc.SubjectAudit)+1 {
		return fmt.Errorf("subject '%s' is not a audit message", subject)
	}
	suffix := subject[len(msc.SubjectAudit)+1:]

	// match handler with subject suffix
	fn, ok := handlers[suffix]
	if !ok {
		l.Warnf("no handler registered for audit message %s", suffix)
		return nil
	}
	return fn(payload, conf, l)
}
