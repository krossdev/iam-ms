// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package mscli

import (
	"net/mail"
	"os"
	"path/filepath"
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

// AddTo a recipient
func (m *Message) AddTo(address *mail.Address) []*mail.Address {
	m.To = append(m.To, address)
	return m.To
}

// AddCc a cc recipient
func (m *Message) AddCc(address *mail.Address) []*mail.Address {
	m.Cc = append(m.Cc, address)
	return m.Cc
}

// AddBcc a bcc recipient
func (m *Message) AddBcc(address *mail.Address) []*mail.Address {
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
