// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
	"net"
	"net/mail"
	"net/url"

	"github.com/mitchellh/mapstructure"
)

const (
	SubjectAction = "kiam.console.action"
	SubjectEvent  = "kiam.console.event"
	SubjectAudit  = "kiam.console.audit"
)

// send request to action subject
func requestAction(action string, payload interface{}) (interface{}, error) {
	subject := fmt.Sprintf("%s.%s", SubjectAction, action)
	return broker.request(subject, payload)
}

const (
	ActionSendVerifyEmail = "sendverifyemail"
	ActionIpLocation      = "iplocation"
)

type SendVerifyEmailPayload struct {
	Subject string `json:"subject"` // mail subject
	Name    string `json:"name"`    // recipient name
	To      string `json:"to"`      // recipient address
	Href    string `json:"href"`    // verify url
	Locale  string `json:"locale"`  // i18n locale
	Expire  string `json:"expire"`  // expire
}

type IPLocationPayload struct {
	IpAddr string `json:"ipaddr"` // ip address to lookup
	Locale string `json:"locale"` // i18n locale
}

type IPLocationReply struct {
	Country string `json:"country"` // country name
	City    string `json:"city"`    // city name
}

// Ask message services to send a verify email
func SendVerifyEmail(payload *SendVerifyEmailPayload) error {
	if payload == nil {
		return fmt.Errorf("payload is empty")
	}
	// validation to address
	if _, err := mail.ParseAddress(payload.To); err != nil {
		return err
	}
	// validation href
	u, err := url.Parse(payload.Href)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return fmt.Errorf("verify url invalid")
	}
	// send the request
	_, err = requestAction(ActionSendVerifyEmail, payload)
	return err
}

// Ask message services to lookup ip location
func IPLocation(payload *IPLocationPayload) (*IPLocationReply, error) {
	if payload == nil {
		return nil, fmt.Errorf("payload is empty")
	}
	ip := net.ParseIP(payload.IpAddr)
	if ip == nil {
		return nil, fmt.Errorf("parse ip(%s) error", payload.IpAddr)
	}
	// if ip.IsLoopback() {
	// 	return i18n.P(c, "Localhost"), nil
	// }
	// if ip.IsPrivate() {
	// 	return i18n.P(c, "Local area network"), nil
	// }
	// city, err := geoip.City(ip)
	// if err != nil {
	// 	return "", err
	// }
	// arr := []string{
	// 	geoip.CountryName(city, i18n.Language(c)),
	// 	geoip.CityName(city, i18n.Language(c)),
	// }
	// return strings.TrimRight(strings.Join(arr, ","), ","), nil

	// send the request and wait for reply
	rp, err := requestAction(ActionIpLocation, payload)
	if err != nil {
		return nil, err
	}
	var reply IPLocationReply

	if err := mapstructure.Decode(rp, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
