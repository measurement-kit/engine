// Package bouncer contains a OONI bouncer client implementation.
//
// Specifically we implement v2.0.0 of the OONI bouncer specification defined
// in https://github.com/ooni/spec/blob/master/backends/bk-004-bouncer.md.
package bouncer

import (
	"context"
	"encoding/json"

	"github.com/measurement-kit/engine/httpx"
	"github.com/measurement-kit/engine/model"
)

// Config contains the bouncer configuration.
type Config struct {
	// BaseURL is the base URL to use.
	BaseURL string
}

// === BEGIN PRE spec v2.0.0 CODE ===

// TODO(bassosimone): if the v2.0.0 spec is approved then we should
// change the code to remove the result indirection.

type result struct {
	Results []model.Service `json:"results"`
}

func kludge(orig []model.Service) (edited []model.Service) {
	for _, e := range orig {
		if e.Type == "https" {
			e.Address = "https://" + e.Address
		}
		edited = append(edited, e)
	}
	return
}

// === END PRE spec v2.0.0 CODE ===

func get(ctx context.Context, config Config, path string) ([]model.Service, error) {
	data, err := httpx.GETWithBaseURL(ctx, config.BaseURL, path)
	if err != nil {
		return nil, err
	}
	var result result
	err = json.Unmarshal(data, &result)
	return kludge(result.Results), err
}

// GetCollectors queries the bouncer for collectors. Returns a list of
// entries on success; an error on failure.
func GetCollectors(ctx context.Context, config Config) ([]model.Service, error) {
	return get(ctx, config, "/api/v1/collectors")
}

// GetTestHelpers is like GetCollectors but for test helpers.
func GetTestHelpers(ctx context.Context, config Config) ([]model.Service, error) {
	return get(ctx, config, "/api/v1/test-helpers")
}
