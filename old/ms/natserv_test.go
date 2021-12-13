// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package ms

import (
	"fmt"
	"github.com/krossdev/iam-ms/config"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	config.LoadTestConfig()
	log.Printf("TestMain conf")
	m.Run()
}

func Test1(t *testing.T) {
	err := connect()
	if err != nil {
		log.Fatal("connect to nats error")
	} else {
		log.Printf("nats connected")
	}

	testSubWithReplyN()

	testRequest()

	select {}
}

func check(err error) {
	if err != nil {
		log.Fatal(errors.WithStack(err))
	}
}

func testRequest() {
	msg, err := nc.Request("test", []byte("hello"), time.Second)
	check(err)
	log.Println("request-reply", string(msg.Data))
	sub, err := nc.Subscribe("test.ack", func(msg *nats.Msg) {
		log.Println("ack", string(msg.Data))
	})
	err = sub.AutoUnsubscribe(2)
	check(err)

	err = nc.PublishRequest("test", "test.ack", []byte("hello"))
	check(err)
	//log.Println("request-reply", string(msg.Data))
}
func testSubWithReplyN() {
	_, err := nc.Subscribe("test", func(msg *nats.Msg) {
		for i := 0; i < 5; i++ {
			// 多次答
			err := msg.Respond([]byte(fmt.Sprintf("%d", i)))
			check(err)
		}
	})
	check(err)

}
