// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE

package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/xlog"
)

// Attachment represents an email attachment.
type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

// Header represents an additional email header.
type Header struct {
	Key   string
	Value string
}

// Message represents a smtp message.
type Message struct {
	To          []*mail.Address
	Cc          []*mail.Address
	Bcc         []*mail.Address
	ReplyTo     string
	Subject     string
	Body        string
	BodyType    string
	Headers     []Header
	Attachments map[string]*Attachment
}

func (m *Message) attach(file string, inline bool) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	_, filename := filepath.Split(file)

	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}
	return nil
}

// AddTO add a recipient
func (m *Message) AddTO(address *mail.Address) []*mail.Address {
	m.To = append(m.To, address)
	return m.To
}

// AddCC add a cc recipient
func (m *Message) AddCC(address *mail.Address) []*mail.Address {
	m.Cc = append(m.Cc, address)
	return m.Cc
}

// AddBCC add a bcc recipient
func (m *Message) AddBCC(address *mail.Address) []*mail.Address {
	m.Bcc = append(m.Bcc, address)
	return m.Bcc
}

// AttachBuffer attaches a binary attachment.
func (m *Message) AttachBuffer(filename string, buf []byte, inline bool) error {
	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     buf,
		Inline:   inline,
	}
	return nil
}

// Attach attaches a file.
func (m *Message) Attach(file string) error {
	return m.attach(file, false)
}

// Inline includes a file as an inline attachment.
func (m *Message) Inline(file string) error {
	return m.attach(file, true)
}

// AddHeader a Header to message
func (m *Message) AddHeader(key string, value string) Header {
	newHeader := Header{Key: key, Value: value}
	m.Headers = append(m.Headers, newHeader)
	return newHeader
}

// bytes returns the mail data
func (m *Message) bytes(sender *mail.Address) []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + sender.String() + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC1123Z) + "\r\n")

	buf.WriteString("To: " + strings.Join(recipients(true, m.To), ",") + "\r\n")
	if len(m.Cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(recipients(true, m.Cc), ",") + "\r\n")
	}

	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(m.Subject)) + "?="
	buf.WriteString("Subject: " + subject + "\r\n")

	if len(m.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + m.ReplyTo + "\r\n")
	}
	buf.WriteString("MIME-Version: 1.0\r\n")

	// add custom headers
	if len(m.Headers) > 0 {
		for _, header := range m.Headers {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", header.Key, header.Value))
		}
	}
	boundary := "f46d043c813270fc6b04c2d223da"

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.BodyType))
	buf.WriteString(m.Body)
	buf.WriteString("\r\n")

	if len(m.Attachments) == 0 {
		return buf.Bytes()
	}
	for _, attachment := range m.Attachments {
		buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

		if attachment.Inline {
			buf.WriteString("Content-Type: message/rfc822\r\n")
			buf.WriteString("Content-Disposition: inline; filename=\"" +
				attachment.Filename + "\"\r\n\r\n")

			buf.Write(attachment.Data)
		} else {
			ext := filepath.Ext(attachment.Filename)
			mimetype := mime.TypeByExtension(ext)
			if mimetype != "" {
				mime := fmt.Sprintf("Content-Type: %s\r\n", mimetype)
				buf.WriteString(mime)
			} else {
				buf.WriteString("Content-Type: application/octet-stream\r\n")
			}
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")

			buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
			buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
			buf.WriteString("?=\"\r\n\r\n")

			b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
			base64.StdEncoding.Encode(b, attachment.Data)

			// write base64 content in lines of up to 76 chars
			for i, l := 0, len(b); i < l; i++ {
				buf.WriteByte(b[i])
				if (i+1)%76 == 0 {
					buf.WriteString("\r\n")
				}
			}
		}
		buf.WriteString("\r\n--" + boundary)
	}
	buf.WriteString("--")
	return buf.Bytes()
}

// Send email with smtp protocol.
func (m *Message) Send() error {
	if mailConfig == nil {
		xlog.X.Panic("email config is empty, please Setup()")
	}
	if len(mailConfig.Mtas) == 0 {
		return fmt.Errorf("no smtp mta found")
	}
	if len(mailConfig.SubjectPrefix) > 0 {
		m.Subject = fmt.Sprintf("%s %s", mailConfig.SubjectPrefix, m.Subject)
	}
	sent := false

	for _, mta := range mailConfig.Mtas {
		if err := m.send(&mta); err != nil {
			xlog.X.Warnf("send mail with '%s' error %v", mta.Name, err)
			continue
		}
		sent = true
		break
	}
	if !sent {
		return fmt.Errorf("try all mta but no one can send this mail")
	}
	return nil
}

// send with smtp
func (m *Message) send(mta *config.Mta) error {
	if len(mta.Host) == 0 {
		return fmt.Errorf("mta '%s' host is empty", mta.Name)
	}
	if len(mta.Sender) == 0 {
		return fmt.Errorf("mta '%s' sender address is empty", mta.Name)
	}
	from, err := mail.ParseAddress(mta.Sender)
	if err != nil {
		return fmt.Errorf("parse sender address '%s' error: %v", mta.Sender, err)
	}
	to := recipients(false, m.To, m.Cc, m.Bcc)

	if mta.SSLMode {
		return m.sendWithSSL(mta, from, to)
	}
	dest := fmt.Sprintf("%s:%d", mta.Host, mta.Port)
	auth := smtp.PlainAuth("", from.Address, mta.Passwd, mta.Host)

	return smtp.SendMail(dest, auth, from.Address, to, m.bytes(from))
}

// here need to call tls.Dial instead of smtp.Dial
// for smtp servers running on 465 that require an ssl connection
// from the very beginning (no starttls)
func (m *Message) sendWithSSL(mta *config.Mta, from *mail.Address, to []string) error {
	dest := fmt.Sprintf("%s:%d", mta.Host, mta.Port)

	conn, err := tls.Dial("tcp", dest, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         mta.Host,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, mta.Host)
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth("", from.Address, mta.Passwd, mta.Host)
	if err = c.Auth(auth); err != nil {
		return err
	}
	// add sender and recipients
	if err = c.Mail(from.Address); err != nil {
		return err
	}
	for _, r := range to {
		if err = c.Rcpt(r); err != nil {
			return err
		}
	}
	// write message body
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(m.bytes(from)); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

func newMessage(subject string, body string, bodyType string) *Message {
	return &Message{
		Subject:     subject,
		Body:        body,
		BodyType:    bodyType,
		Attachments: make(map[string]*Attachment),
	}
}

// TextMessage returns a new Message that can compose an email with attachments
func TextMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/plain")
}

// HTMLMessage returns a new Message that can compose an HTML email with attachments
func HTMLMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/html")
}

// recipients convert mail.Address array to string array
//   with name: Name <mail@addr.com>
//   with no name: mail@addr.com
func recipients(withName bool, args ...[]*mail.Address) []string {
	var list []string

	for _, arg := range args {
		for _, a := range arg {
			if withName {
				list = append(list, a.String())
			} else {
				list = append(list, a.Address)
			}
		}
	}
	return list
}
