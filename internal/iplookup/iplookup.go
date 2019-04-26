// Package iplookup discovers the probe IP (aka probeIP).
package iplookup

import (
	"context"
	"encoding/xml"
	"errors"
	"net"

	"github.com/measurement-kit/engine/internal/httpx"
)

type response struct {
	XMLName xml.Name `xml:"Response"`
	IP      string   `xml:"Ip"`
}

// httpxGET allows mocking httpx.GET
var httpxGET = httpx.GET

// Perform lookups the probeIP. On failure, probeIP is set to "127.0.0.1".
func Perform(ctx context.Context) (string, error) {
	data, err := httpxGET(ctx, "https://geoip.ubuntu.com/lookup")
	if err != nil {
		return "127.0.0.1", err
	}
	v := response{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		return "127.0.0.1", err
	}
	if net.ParseIP(v.IP) == nil {
		err = errors.New("Invalid IP address")
		return "127.0.0.1", err
	}
	return v.IP, nil
}
