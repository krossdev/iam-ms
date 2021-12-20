package audit

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type handlerFunc func(payload interface{}, params interface{}, l *logrus.Entry) error

var handlers map[string]handlerFunc

func Handler(audit string, payload interface{}, params interface{}, l *logrus.Entry) error {
	f, ok := handlers[audit]
	if !ok {
		return fmt.Errorf("audit %s not found", audit)
	}
	return f(payload, params, l)
}
