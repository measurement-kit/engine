// Package mlabns contains an mlabns implementation.
package mlabns

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/measurement-kit/engine/internal/httpx"
)

// Config contains mlabns settings.
type Config struct {
	// BaseURL is the optional base URL for contacting mlabns.
	BaseURL string

	// Tool is the mandatory tool to use.
	Tool string
}

// Server describes a mlab server.
type Server struct {
	// FQDN is the the FQDN of the server.
	FQDN string `json:"fqdn"`
}

// httpGET allows mocking httpx.GET
var httpxGET = httpx.GET

// GeoOptions returns some nearby mlab servers.
func GeoOptions(ctx context.Context, config Config) ([]Server, error) {
	if config.BaseURL == "" {
		config.BaseURL = "https://mlab-ns.appspot.com/"
	}
	URL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}
	URL.Path = config.Tool
	query := URL.Query()
	query.Add("policy", "geo_options")
	URL.RawQuery = query.Encode()
	data, err := httpxGET(ctx, URL.String())
	if err != nil {
		return nil, err
	}
	var servers []Server
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}