// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE

package email

import (
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/xlog"
)

// configuration
var mailConfig *config.Mail

// propagation configuration from config file
func Setup(conf *config.Mail) {
	// sort the mta, take preferred mta to the first
	if len(conf.PreferredMta) > 0 {
		found := false

		for i, m := range conf.Mtas {
			if m.Name == conf.PreferredMta {
				m2 := conf.Mtas[i]

				// keep the order in the configuration file
				for j := i; j > 0; j -= 1 {
					conf.Mtas[j] = conf.Mtas[j-1]
				}
				conf.Mtas[0] = m2
				found = true
				break
			}
		}
		if !found {
			xlog.X.Warnf("preferred mta '%s' not found", conf.PreferredMta)
		}
	}
	mailConfig = conf
}
