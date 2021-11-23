// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package test

import (
	"github.com/krossdev/iam-ms/config"
	"github.com/krossdev/iam-ms/ms"
	"github.com/krossdev/iam-ms/mscli"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	config.LoadTestConfig()
	log.Printf("load config")

	c := config.G.Nats

	err := ms.Setup()
	if err != nil {
		log.Fatalf("ms setup error %+v", err)
	} else {
		log.Println("ms setup ok")
	}

	params := mscli.Params{}
	params.Url = c.Url
	err = mscli.Setup(&params)
	if err != nil {
		log.Fatalf("mscli nats setup error %+v", err)
	} else {
		log.Println("mscli setup ok")
	}

	m.Run()
}
