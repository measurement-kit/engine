// Package task defines Measurement Kit tasks.
package task

import (
	"context"
	"errors"
	"time"

	"github.com/measurement-kit/engine/internal/nettest"
	"github.com/measurement-kit/engine/internal/nettest/ndt7"
	"github.com/measurement-kit/engine/model"
)

// Config contains the task settings
type Config struct {
	// Inputs is the list of inputs for the measurement task
	Inputs []string
}

func discoverAvailableCollectors(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event,
) error {
	out <- model.NewLogInfoEvent("discovering available collectors")
	err := nt.DiscoverAvailableCollectors(ctx)
	if err != nil {
		out <- model.NewLogWarningEvent(
			err, "cannot discover available collectors",
		)
		return err
	}
	return nil
}

func discoverAvailableTestHelpers(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event,
) error {
	out <- model.NewLogInfoEvent("discovering available test helpers")
	err := nt.DiscoverAvailableTestHelpers(ctx)
	if err != nil {
		out <- model.NewLogWarningEvent(
			err, "cannot discover available test helpers",
		)
		return err
	}
	return nil
}

func openReport(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event,
) error {
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
		return errors.New("cannot open report")
	}
	return nil
}

func performMeasurement(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event, input string,
) (model.Measurement, error) {
	out <- model.NewLogInfoEvent("starting measurement")
	measurement := nt.NewMeasurement()
	start := time.Now()
	for ev := range nt.StartMeasurement(ctx, input, &measurement) {
		out <- ev
	}
	measurement.Input = input
	measurement.MeasurementRuntime = time.Now().Sub(start).Seconds()
	out <- model.NewLogInfoEvent("measurement complete")
	return measurement, nil
}

func emitMeasurement(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event, measurement model.Measurement,
) error {
	ev, err := model.NewMeasurementEvent(measurement)
	if err != nil {
		out <- model.NewLogWarningEvent(err, "cannot serialize measurement")
		return err
	}
	out <- ev
	return nil
}

func submitMeasurement(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event, measurement model.Measurement,
) error {
	out <- model.NewLogInfoEvent("submitting the measurement")
	err := nt.SubmitMeasurement(ctx, &measurement)
	if err != nil {
		out <- model.NewLogWarningEvent(
			err, "failed to submit the measurement",
		)
		return err
	}
	out <- model.NewLogInfoEvent("measurement submitted")
	return nil
}

func performEmitAndSubmitMeasurement(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event, input string,
) {
	measurement, err := performMeasurement(ctx, nt, config, out, input)
	if err != nil {
		return
	}
	err = emitMeasurement(ctx, nt, config, out, measurement)
	if err != nil {
		return // if we cannot emit it's not serializable, so stop here
	}
	err = submitMeasurement(ctx, nt, config, out, measurement)
	if err != nil {
		return
	}
}

func performTask(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event,
) {
	defer close(out) // tell the reader we're done
	err := discoverAvailableCollectors(ctx, nt, config, out)
	if err != nil {
		return
	}
	err = discoverAvailableTestHelpers(ctx, nt, config, out)
	if err != nil {
		return
	}
	// TODO(bassosimone): discover probe IP
	// TODO(bassosimone): discover probe ASN
	// TODO(bassosimone): discover probe CC
	// TODO(bassosimone): discover probe network name
	// TODO(bassosimone): discover probe resolver IP
	err = openReport(ctx, nt, config, out)
	if err != nil {
		return
	}
	defer nt.CloseReport(ctx)
	// TODO(bassosimone): implement parallelism
	for _, input := range config.Inputs {
		performEmitAndSubmitMeasurement(ctx, nt, config, out, input)
	}
}

func startTaskAndFilterEvents(
	ctx context.Context, nt *nettest.Nettest,
	config Config, out chan<- model.Event,
) {
	// Implementation note: this is the right place where to implement
	// logging on file and saving measurements on file. We just need to
	// intercept the proper events and write on the file system.
	//
	// Therefore we create a cancellable ctx for this function and we
	// use a child channel so we can filter events.
	//
	// TODO(bassosimone): open the required file descriptors.
	defer close(out)
	innerctx, cancel := context.WithCancel(ctx)
	defer cancel()
	in := make(chan model.Event)
	go performTask(innerctx, nt, config, in)
	for ev := range in {
		// TODO(bassosimone): filter and write on file system here.
		// TODO(bassosimone): also fiter logs by verbosity.
		out <- ev
	}
}

// StartNdt7 starts a new ndt7 task.
func StartNdt7(ctx context.Context, config Config) <-chan model.Event {
	out := make(chan model.Event)
	nt := ndt7.NewNettest()
	config.Inputs = []string{""} // force running just once
	go startTaskAndFilterEvents(ctx, nt, config, out)
	return out
}
