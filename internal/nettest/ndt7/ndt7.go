// Package ndt7 contains the ndt7 client.
package ndt7

import (
	"context"

	upstream "github.com/m-lab/ndt7-client-go"
	upstreamSpec "github.com/m-lab/ndt7-client-go/spec"
	"github.com/measurement-kit/engine/internal/nettest"
	"github.com/measurement-kit/engine/internal/version"
	"github.com/measurement-kit/engine/model"
)

// Client is a ndt7 client
type Client struct {
	nettest *nettest.Nettest
}

// testKeys contains the test keys
type testKeys struct {
	// Failure is the failure string
	Failure string `json:"failure"`

	// Download contains download results
	Download []upstreamSpec.Measurement `json:"download"`

	// Upload contains upload results
	Upload []upstreamSpec.Measurement `json:"upload"`
}

// run runs a ndt7 test
func run(
	ctx context.Context,
	input string,
	measurement *model.Measurement,
	out chan<- model.Event,
) {
	testkeys := &testKeys{}
	measurement.TestKeys = testkeys
	client := upstream.NewClient("MKengine/" + version.Version)
	ch, err := client.StartDownload(ctx)
	if err != nil {
		testkeys.Failure = err.Error()
		out <- model.NewFailureMeasurementEvent(0, err)
		return
	}
	for ev := range ch {
		testkeys.Download = append(testkeys.Download, ev)
		out <- model.Event{
			Key:   "ndt7.download",
			Value: ev,
		}
	}
	ch, err = client.StartUpload(ctx)
	if err != nil {
		testkeys.Failure = err.Error()
		out <- model.NewFailureMeasurementEvent(0, err)
		return
	}
	for ev := range ch {
		testkeys.Upload = append(testkeys.Upload, ev)
		out <- model.Event{
			Key:   "ndt7.upload",
			Value: ev,
		}
	}
}

// NewNettest creates a new ndt7 client nettest
func NewNettest() *nettest.Nettest {
	return &nettest.Nettest{
		TestName:        "ndt7",
		TestVersion:     "0.1.0",
		SoftwareName:    "MKEngine",
		SoftwareVersion: version.Version,
		TestStartTime:   nettest.FormatTimeNowUTC(),
		Main:            run,
	}
}
