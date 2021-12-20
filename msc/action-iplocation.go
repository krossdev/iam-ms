package msc

import (
	"fmt"
	"net"

	"github.com/mitchellh/mapstructure"
)

type IPLocationActionPayload struct {
	IpAddr string `json:"ipaddr"` // ip address to lookup
	Locale string `json:"locale"` // i18n locale
}

type IPLocationActionReply struct {
	Continent string  `json:"continent"` // continent name
	Country   string  `json:"country"`   // country name
	City      string  `json:"city"`      // city name
	Longitude float64 `json:"longitude"` // longitude
	Latitude  float64 `json:"latitude"`  // latitude
	TimeZone  string  `json:"timezone"`  // time zone
}

// Ask message services to lookup ip location
func IPLocationAction(payload *IPLocationActionPayload) (*IPLocationActionReply, error) {
	if payload == nil {
		return nil, fmt.Errorf("payload is empty")
	}
	ip := net.ParseIP(payload.IpAddr)
	if ip == nil {
		return nil, fmt.Errorf("parse ip(%s) error", payload.IpAddr)
	}
	// loopback and private address cannot location
	if ip.IsLoopback() || ip.IsPrivate() {
		return nil, fmt.Errorf("ip(%s) is a private address", payload.IpAddr)
	}

	// send the request and wait for reply
	rp, err := requestAction(ActionIpLocation, payload)
	if err != nil {
		return nil, err
	}
	var reply IPLocationActionReply

	// decode the response map to struct
	if err := mapstructure.Decode(rp, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
