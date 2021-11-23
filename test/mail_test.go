// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package test

import (
	"fmt"
	"github.com/krossdev/iam-ms/mscli"
	"github.com/stretchr/testify/assert"
	"net/mail"
	"sync"
	"testing"
	"time"
)

func TestSendMail(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	m := mscli.TextMessage(fmt.Sprintf("测试 - %d", time.Now().Unix()), "world你好")

	// Add to/cc recipients
	m.AddTo(&mail.Address{Name: "张三", Address: "11560648@qq.com"})

	// Add attachments
	m.AttachBuffer("buffer.txt", []byte("Hello World"), false)

	e := &mscli.Mail{
		Reply:   true,
		Message: m,

		// file schema
		Schema: "hello.yaml",

		// AfterHook is called after the email is actually sent
		AfterHook: func(_ *mscli.Mail, err error) error {
			if err != nil {
				t.Logf("mscli SendMail return : %s", err)
			} else {
				t.Logf("mscli SendMail return OK")
			}
			wg.Done()
			return err
		},
	}
	err := mscli.SendMail(e)
	assert.NoError(t, err, "SendMail() should not return error")

	wg.Wait()

	//time.Sleep(10 * time.Second)
}

func TestSendMailFail(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	m := mscli.TextMessage(fmt.Sprintf("测试 - %d", time.Now().Unix()), "world你好")
	m.AddTo(&mail.Address{Name: "No exist address", Address: "xxx-11560648@qq.com"})

	e := &mscli.Mail{
		Reply:   true,
		Message: m,
		Schema:  "hello-error.yaml", // schema不存在
		AfterHook: func(_ *mscli.Mail, err error) error {
			assert.Error(t, err, "schema not exist")
			t.Logf("Hook: mscli SendMail return : %s", err)
			wg.Done()
			return err
		},
	}
	err := mscli.SendMail(e)
	assert.NoError(t, err, "SendMail() should not return error")

	wg.Wait()

	//time.Sleep(10 * time.Second)
}
