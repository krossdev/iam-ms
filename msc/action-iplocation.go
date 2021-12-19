// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
	"net"

	"github.com/mitchellh/mapstructure"
)

type IPLocationPayload struct {
	IpAddr string `json:"ipaddr"` // ip address to lookup
	Locale string `json:"locale"` // i18n locale
}

type IPLocationReply struct {
	Country string `json:"country"` // country name
	City    string `json:"city"`    // city name
	Address string `json:"address"` // full address
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
	// loopback and private address cannot location
	if ip.IsLoopback() {
		return &IPLocationReply{City: "Localhost"}, nil
	}
	if ip.IsPrivate() {
		return &IPLocationReply{City: "Local area network"}, nil
	}

	// send the request and wait for reply
	rp, err := requestAction(ActionIpLocation, payload)
	if err != nil {
		return nil, err
	}
	var reply IPLocationReply

	// decode the response map to struct
	if err := mapstructure.Decode(rp, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
