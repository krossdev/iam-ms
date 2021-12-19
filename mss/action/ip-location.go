// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package action

import (
	"fmt"
	"net"

	"github.com/krossdev/iam-ms/msc"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func IPLocationHandler(p interface{}, params interface{}, l *logrus.Entry) (interface{}, error) {
	var payload msc.IPLocationPayload

	// convert map to struct
	if err := mapstructure.Decode(p, &payload); err != nil {
		return nil, err
	}
	ip := net.ParseIP(payload.IpAddr)
	if ip == nil {
		return nil, fmt.Errorf("parse ip(%s) error", payload.IpAddr)
	}
	return nil, nil
}
