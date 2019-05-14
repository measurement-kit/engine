// Package bouncer contains a OONI bouncer client implementation.
//
// Specifically we implement v2.0.0 of the OONI bouncer specification defined
// in https://github.com/ooni/spec/blob/master/backends/bk-004-bouncer.md.
package bouncer

import (
	"context"
	"encoding/json"

	"github.com/measurement-kit/engine/internal/httpx"
	"github.com/measurement-kit/engine/model"
)

// Config contains the bouncer configuration.
type Config struct {
	// BaseURL is the optional bouncer base URL to use.
	BaseURL string
}

// GetCollectors queries the bouncer for collectors. Returns a list of
// entries on success; an error on failure.
func GetCollectors(ctx context.Context, config Config) ([]model.Service, error) {
	data, err := httpx.GETWithBaseURL(ctx, config.BaseURL, "/api/v1/collectors")
	if err != nil {
		return nil, err
	}
	var result []model.Service
	err = json.Unmarshal(data, &result)
	return result, err
}

// GetTestHelpers is like GetCollectors but for test helpers.
func GetTestHelpers(ctx context.Context, config Config) (map[string][]model.Service, error) {
	data, err := httpx.GETWithBaseURL(ctx, config.BaseURL, "/api/v1/test-helpers")
	if err != nil {
		return nil, err
	}
	var result map[string][]model.Service
	err = json.Unmarshal(data, &result)
	return result, err
}
