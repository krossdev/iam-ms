// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	//t.Log(MustGetwd())
	//err := LoadConfig("../conf/conf.yaml")
	////fmt.Printf("stack trace:\n%+v\n", err)
	//assert.NoError(t, err, "LoadConfig should return no error")

	LoadTestConfig()

	t.Log(G.RootDir)
	t.Log(G.Mail.SchemaDir)
}
