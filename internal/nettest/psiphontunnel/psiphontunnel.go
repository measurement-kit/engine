// Package psiphontunnel implements the psiphontunnel nettest.
package psiphontunnel

import (
	"context"
	"time"

	"github.com/measurement-kit/engine/internal/model"
	"github.com/measurement-kit/engine/internal/nettest"
	"github.com/measurement-kit/engine/internal/nettest/psiphontunnel/runner"
)

// Config contains the psiphontunnel nettest configuration.
type Config = runner.Config

// NewNettest creates a new psiphontunnel nettest. This function
// initializes the following nettest fields:
//
// - Ctx
// - TestName
// - TestVersion
// - TestStartTime
// - Measure
//
// Call nettest.StartMeasurement("", &measurement) to perform a measurement.
func NewNettest(ctx context.Context, config Config) *nettest.Nettest {
	return &nettest.Nettest{
		Ctx:           ctx,
		TestName:      "psiphontunnel",
		TestVersion:   "0.0.1",
		TestStartTime: nettest.FormatTimeNowUTC(),
		Main: func(input string, m *model.Measurement, ch chan<- model.Event) {
			t0 := time.Now()
			m.TestKeys = runner.Run(ctx, config)
			m.MeasurementRuntime = time.Now().Sub(t0).Seconds()
		},
	}
}
