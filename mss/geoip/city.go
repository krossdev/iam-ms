// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package geoip

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type City struct {
	city   *geoip2.City
	locale string
}

// Name return localed city name
func (c *City) Name() string {
	if name, ok := c.city.City.Names[c.locale]; ok {
		return name
	}
	return c.city.City.Names["en"]
}

// Continent return localed continent name
func (c *City) Continent() string {
	if name, ok := c.city.Continent.Names[c.locale]; ok {
		return name
	}
	return c.city.Continent.Names["en"]
}

// Country return localed country name
func (c *City) Country() string {
	if name, ok := c.city.Country.Names[c.locale]; ok {
		return name
	}
	return c.city.Country.Names["en"]
}

// Longitude
func (c *City) Longitude() float64 {
	return c.city.Location.Longitude
}

// Latitude
func (c *City) Latitude() float64 {
	return c.city.Location.Latitude
}

// Time zone
func (c *City) TimeZone() string {
	return c.city.Location.TimeZone
}

// lookup city database
func CityLookup(ipaddr string, locale string) (*City, error) {
	ip := net.ParseIP(ipaddr)
	if ip == nil {
		return nil, fmt.Errorf("parse ip(%s} error", ipaddr)
	}
	// open geoip city database
	db, err := geoip2.Open(geoipConfig.Path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// lookup
	result, err := db.City(ip)
	if err != nil {
		return nil, err
	}
	return &City{city: result, locale: locale}, nil
}
