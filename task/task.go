// Package task defines Measurement Kit tasks.
package task

import (
	"context"
	"time"

	"github.com/measurement-kit/engine/internal/nettest"
	"github.com/measurement-kit/engine/internal/nettest/ndt7"
	"github.com/measurement-kit/engine/model"
)

// Config contains the task settings
type Config struct{}

func runWithoutInput(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event,
) {
	defer close(out)
	out <- model.NewLogInfoEvent("discovering available collectors")
	err := nt.DiscoverAvailableCollectors(ctx)
	if err != nil {
		out <- model.NewLogWarningEvent(
			err, "cannot discover available collectors",
		)
		return
	}
	out <- model.NewLogInfoEvent("discovering available test helpers")
	err = nt.DiscoverAvailableTestHelpers(ctx)
	if err != nil {
		out <- model.NewLogWarningEvent(
			err, "cannot discover available test helpers",
		)
		return
	}
	out <- model.NewLogInfoEvent("opening report")
	for err := range nt.OpenReport(ctx) {
		out <- model.NewLogWarningEvent(
			err, "cannot open report; trying other collectors",
		)
	}
	if nt.Report.ID == "" {
		out <- model.NewLogWarningEvent(
			nil, "failed to open report with all collectors",
		)
		return
	}
	defer nt.CloseReport(ctx)
	out <- model.NewLogInfoEvent("starting measurement")
	measurement := nt.NewMeasurement()
	start := time.Now()
	for ev := range nt.StartMeasurement(ctx, "", &measurement) {
		out <- ev
	}
	measurement.MeasurementRuntime = time.Now().Sub(start).Seconds()
	out <- model.NewLogInfoEvent("measurement complete")
	measurementEvent, err := model.NewMeasurementEvent(measurement)
	if err != nil {
		out <- model.NewLogWarningEvent(err, "cannot serialize measurement")
		return
	}
	out <- measurementEvent
	out <- model.NewLogInfoEvent("submitting the measurement")
	err = nt.SubmitMeasurement(ctx, &measurement)
	if err != nil {
		out <- model.NewLogWarningEvent(
			err, "failed to submit the measurement",
		)
		return
	}
	out <- model.NewLogInfoEvent("measurement submitted")
}

// StartNdt7 starts a new ndt7 task.
func StartNdt7(ctx context.Context, config Config) <-chan model.Event {
	out := make(chan model.Event)
	nt := ndt7.NewNettest()
	go runWithoutInput(ctx, nt, config, out)
	return out
}
