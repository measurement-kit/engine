// Package ndt7 implements the ndt7 nettest.
package ndt7

import (
	"context"
	"time"

	"github.com/measurement-kit/engine/model"
	"github.com/measurement-kit/engine/nettest"
	ndt7model "github.com/measurement-kit/engine/nettest/ndt7/runner/model"
	"github.com/measurement-kit/engine/nettest/ndt7/runner"
)

// Config contains the ndt7 nettest configuration.
type Config struct {
	// FQDNs contains the optional server FQDNs to use.
	FQDNs []string
}

// getservers returns you a list of suitable servers
func getservers(ctx context.Context, config Config) ([]string, error) {
	if len(config.FQDNs) > 0 {
		return config.FQDNs, nil
	}
	return runner.GetServers(ctx)
}

// getServersResults contains the results of getting servers.
type getServersResults struct {
	// Failure indicates whether there was a failure.
	Failure string `json:"failure"`

	// Servers contains the list of available servers.
	Servers []string `json:"fqdns"`
}

// subtestResults contains the results of subtests.
type subtestResults struct {
	// Server is the server we're using.
	Server string `json:"fqdn"`

	// Failure indicates whether we could not start the nettest.
	Failure string `json:"failure"`

	// Measurements contains the measurements.
	Measurements []ndt7model.Measurement `json:"measurements"`
}

// results contains the nettest results.
type results struct {
	// GetServersResults contains the results of getting servers.
	GetServersResults getServersResults `json:"get_servers_results"`

	// DownloadResults contains the download results.
	DownloadResults []subtestResults `json:"download_results"`

	// UploadResults contains the upload results.
	UploadResults []subtestResults `json:"upload_results"`
}

// run runs the ndt7 test.
func run(ctx context.Context, config Config, out chan<- model.Event) results {
	var results results
	FQDNs, err := getservers(ctx, config)
	if err != nil {
		results.GetServersResults.Failure = err.Error()
		return results
	}
	results.GetServersResults.Servers = FQDNs
	for _, FQDN := range FQDNs {
		var sr subtestResults
		sr.Server = FQDN
		in, err := runner.StartDownload(ctx, FQDN)
		if err != nil {
			sr.Failure = err.Error()
			results.DownloadResults = append(results.DownloadResults, sr)
			continue
		}
		for measurement := range in {
			out <- model.Event{Key: "performance.ndt7", Value: measurement}
			sr.Measurements = append(sr.Measurements, measurement)
		}
		results.DownloadResults = append(results.DownloadResults, sr)
		break
	}
	for _, FQDN := range FQDNs {
		in, err := runner.StartUpload(ctx, FQDN)
		var sr subtestResults
		sr.Server = FQDN
		if err != nil {
			sr.Failure = err.Error()
			results.UploadResults = append(results.UploadResults, sr)
			continue
		}
		for measurement := range in {
			out <- model.Event{Key: "performance.ndt7", Value: measurement}
			sr.Measurements = append(sr.Measurements, measurement)
		}
		results.UploadResults = append(results.UploadResults, sr)
		break
	}
	return results
}

// NewNettest creates a new ndt7 nettest. This function
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
		TestName:      "ndt7",
		TestVersion:   "0.0.1",
		TestStartTime: nettest.FormatTimeNowUTC(),
		Main: func(input string, m *model.Measurement, ch chan<- model.Event) {
			t0 := time.Now()
			m.TestKeys = run(ctx, config, ch)
			m.MeasurementRuntime = time.Now().Sub(t0).Seconds()
		},
	}
}
