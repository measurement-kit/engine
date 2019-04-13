// Package iplookup discovers the probe IP (aka probeIP).
package iplookup

import (
	"context"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
)

func get(ctx context.Context, URL string) ([]byte, error) {
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("The request failed")
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

type response struct {
	XMLName xml.Name `xml:"Response"`
	IP      string   `xml:"Ip"`
}

// Perform lookups the probeIP. On failure, probeIP is set to "127.0.0.1".
func Perform(ctx context.Context) (string, error) {
	data, err := get(ctx, "https://geoip.ubuntu.com/lookup")
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
