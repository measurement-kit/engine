// Package psiphontunnel implements the psiphontunnel nettest.
package psiphontunnel

import (
	"context"

	"github.com/measurement-kit/engine/internal/nettest"
	"github.com/measurement-kit/engine/internal/nettest/psiphontunnel/runner"
	"github.com/measurement-kit/engine/internal/version"
	"github.com/measurement-kit/engine/model"
)

// Config contains the psiphontunnel nettest configuration.
type Config = runner.Config

// NewNettest creates a new psiphontunnel nettest.
func NewNettest(config Config) *nettest.Nettest {
	return &nettest.Nettest{
		TestName:        "psiphontunnel",
		TestVersion:     "0.0.1",
		SoftwareName:    "MKEngine",
		SoftwareVersion: version.Version,
		TestStartTime:   nettest.FormatTimeNowUTC(),
		Main: func(
			ctx context.Context,
			input string,
			measurement *model.Measurement,
			out chan<- model.Event,
		) {
			measurement.TestKeys = runner.Run(ctx, config)
		},
	}
}
