// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package email

import (
	"github.com/stretchr/testify/assert"
	"net/mail"
	"testing"
)

func TestSendHtml(t *testing.T) {
	//html, err := os.ReadFile(util.Config.SchemaDir + "/mail.html")
	//if err != nil {
	//	t.Errorf("Failed to read mail.html: %v", err)
	//	return
	//}
	//m2 := HTMLMessage("测试邮件1", string(html))
	m := TextMessage("hello测试", "world你好")

	// Add to/cc recipients
	m.AddTo(&mail.Address{Name: "张三", Address: "11560648@qq.com"})
	//m.AddCc(&mail.Address{Name: "李四", Address: "11560648@qq.com"})

	// Add attachments
	//m.Attach("/tmp/dumb.sqlite")
	//m.Attach("/tmp/kiam.log")
	m.AttachBuffer("buffer.txt", []byte("Hello"), false)

	e := &Mail{
		Message: m,

		// file schema
		Schema: "hello.yaml",

		// AfterHook is called after the email is actually sent
		//AfterHook: func(mi *Mail, err error) error {
		//	return err
		//},
	}
	err := Send(e)
	assert.NoError(t, err, "Send() should not return error")
	if err != nil {
		t.Errorf("Send mail failed: %v", err)
		panic(err)
	}
}
