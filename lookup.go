package main

import (
	"errors"
	"net"
	"strings"
)

type MaxmindRecord struct {
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Continent struct {
		Code  string            `maxminddb:"code"`
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Country struct {
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		IsoCode           string            `maxminddb:"iso_code"`
		Names             map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
		TimeZone  string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
	Subdivisions []struct {
		IsoCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	Traits struct {
		IsAnonymousProxy bool `maxminddb:"is_anonymous_proxy"`
	} `maxminddb:"traits"`
}

type IpInfo struct {
	IP            net.IP   `json:"ip"`
	Name          *string  `json:"name"`
	GpsLatitude   *float64 `json:"gps_lat"`
	GpsLongitude  *float64 `json:"gps_lng"`
	Timezone      *string  `json:"timezone"`
	CityName      *string  `json:"city_name"`
	ContinentCode *string  `json:"continent_code"`
	ContinentName *string  `json:"continent_name"`
	CountryCode   *string  `json:"country_code"`
	CountryName   *string  `json:"country_name"`
	StateCode     *string  `json:"state_code"`
	StateName     *string  `json:"state_name"`
	IsProxy       bool     `json:"is_proxy"`
}

func lookupIp(inputIp string) (*IpInfo, error) {
	var info IpInfo
	if inputIp == "" {
		return &info, errors.New("An IP address must be provided.")
	}
	ip := net.ParseIP(inputIp)
	if ip == nil {
		return &info, errors.New("IP address is not a valid IP.")
	}
	info.IP = ip
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return &info, errors.New("IP address is not allowed.")
	}
	var record MaxmindRecord
	err := ipDatabase.Lookup(ip, &record)
	if err != nil {
		return &info, err
	}

	info.GpsLatitude = &record.Location.Latitude
	info.GpsLongitude = &record.Location.Longitude
	info.Timezone = &record.Location.TimeZone
	info.CityName = getNameFromIpNames(record.City.Names)
	info.ContinentCode = &record.Continent.Code
	info.ContinentName = getNameFromIpNames(record.Continent.Names)
	info.CountryCode = &record.Country.IsoCode
	info.CountryName = getNameFromIpNames(record.Country.Names)
	if len(record.Subdivisions) != 0 {
		info.StateCode = &record.Subdivisions[0].IsoCode
		info.StateName = getNameFromIpNames(record.Subdivisions[0].Names)
	}
	info.IsProxy = record.Traits.IsAnonymousProxy
	info.Name = buildCompositeName(info)

	return &info, nil
}

func buildCompositeName(info IpInfo) *string {
	var names []string
	if info.CityName != nil && *info.CityName != "" {
		names = append(names, *info.CityName)
	}
	if info.StateName != nil && *info.StateName != "" {
		names = append(names, *info.StateName)
	}
	if info.CountryName != nil && *info.CountryName != "" {
		names = append(names, *info.CountryName)
	}
	if len(names) != 0 {
		name := strings.Join(names, ", ")
		return &name
	}
	return nil
}

func getNameFromIpNames(names map[string]string) *string {
	val, ok := names["en"]
	if ok {
		return &val
	} else {
		return nil
	}
}
