// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package email

import (
	"github.com/krossdev/iam-ms/config"
	"log"
	"testing"
)

func TestSmtp1(t *testing.T) {
	t.Log(config.G.Mail.SchemaDir)
}

func TestMain(m *testing.M) {
	config.LoadTestConfig()
	log.Printf("TestMain conf")
	m.Run()
}
