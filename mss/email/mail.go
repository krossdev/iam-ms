// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE

package email

import (
	"github.com/krossdev/iam-ms/mss/config"
)

// configuration
var mailConfig *config.Mail

// propagation configuration from config file
func Setup(conf *config.Mail) {
	mailConfig = conf
}
