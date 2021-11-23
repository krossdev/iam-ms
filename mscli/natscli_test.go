// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package mscli

import (
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	params := Params{
		Url: "nats://3secret@localhost",
	}
	err := Setup(&params)
	if err != nil {
		log.Fatalf("mscli nats setup error %+v", err)
	} else {
		log.Println("mscli setup ok")
	}

	m.Run()
}
