// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package geoip

import (
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/xlog"
	"github.com/oschwald/geoip2-golang"
)

// geoip configuration
var geoipConfig *config.Geoip

func Setup(conf *config.Geoip) error {
	db, err := geoip2.Open(conf.Path)
	if err != nil {
		return err
	}
	geoipConfig = conf

	xlog.X.Debugf("use geoip database %s", conf.Path)
	return db.Close()
}
