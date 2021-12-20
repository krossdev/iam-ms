// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package action

import (
	"fmt"
	"net"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/geoip"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func isLookupable(ipaddr string) error {
	ip := net.ParseIP(ipaddr)
	if ip == nil {
		return fmt.Errorf("parse ip(%s) error", ipaddr)
	}
	if ip.IsLoopback() {
		return fmt.Errorf("ip(%s) is a loopback address", ipaddr)
	}
	if ip.IsPrivate() {
		return fmt.Errorf("ip(%s) is a private address", ipaddr)
	}
	if ip.IsMulticast() {
		return fmt.Errorf("ip(%s) is a multicast address", ipaddr)
	}
	return nil
}

func IPLocationHandler(p interface{}, c interface{}, l *logrus.Entry) (interface{}, error) {
	var payload msc.IPLocationActionPayload
	var conf config.ActionIPLocation

	// convert map to struct
	if err := mapstructure.Decode(p, &payload); err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(c, &conf); err != nil {
		return nil, err
	}
	if err := isLookupable(payload.IpAddr); err != nil {
		return nil, err
	}

	// lookup location with geoip database
	if conf.Engine == config.IPLocationEngineGeoip {
		return ipLocationWithGeoip(payload.IpAddr, payload.Locale)
	}

	return nil, nil
}

// lookup ip location by geoip engine
func ipLocationWithGeoip(ipaddr string, locale string) (interface{}, error) {
	city, err := geoip.CityLookup(ipaddr, locale)
	if err != nil {
		return nil, err
	}
	reply := msc.IPLocationActionReply{
		Continent: city.Continent(),
		Country:   city.Country(),
		City:      city.Name(),
		Longitude: city.Longitude(),
		Latitude:  city.Latitude(),
		TimeZone:  city.TimeZone(),
	}
	return &reply, nil
}
