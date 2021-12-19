// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package geoip

import (
	"net"

	"github.com/krossdev/iam-ms/mss/config"
	"github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
)

// geoip configuration
var geoipConfig *config.Geoip

func Setup(conf *config.Geoip) error {
	db, err := geoip2.Open(conf.Path)
	if err != nil {
		return err
	}
	geoipConfig = conf

	return db.Close()
}

// lookup city from ip address
func City(ip net.IP) (*geoip2.City, error) {
	db, err := geoip2.Open(geoipConfig.Path)
	if err != nil {
		return nil, errors.Wrap(err, "geoip")
	}
	defer db.Close()

	city, err := db.City(ip)
	if err != nil {
		return nil, errors.Wrap(err, "geoip")
	}
	return city, nil
}

// CityName return localed city name
func CityName(city *geoip2.City, lang string) string {
	if name, ok := city.City.Names[lang]; ok {
		return name
	}
	return city.City.Names["en"]
}

// CountryName return localed city name
func CountryName(city *geoip2.City, lang string) string {
	if name, ok := city.Country.Names[lang]; ok {
		return name
	}
	return city.Country.Names["en"]
}
