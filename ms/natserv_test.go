// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package ms

import (
	"github.com/krossdev/iam-ms/config"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	config.LoadTestConfig()
	log.Printf("TestMain conf")
	m.Run()
}

func Test1(t *testing.T) {
	c := config.G.Nats
	log.Println(c.Url)
	err := Setup()
	assert.NoErrorf(t, err, "nats setup")

	select {}
}
