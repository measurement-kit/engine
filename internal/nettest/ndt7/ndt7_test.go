package ndt7

import (
	"context"
	"errors"
	"testing"

	"github.com/apex/log"
	"github.com/measurement-kit/engine/internal/model"
	ndt7model "github.com/measurement-kit/engine/internal/nettest/ndt7/runner/model"
)

// TestGetserversCustomFQDNs checks whether we can configure
// custom FQDNs for running a ndt7 nettest.
func TestGetserversCustomFQDNs(t *testing.T) {
	config := Config{
		FQDNs: []string{
			"a", "b", "c", "d",
		},
	}
	FQDNs, err := getservers(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	if len(FQDNs) != len(config.FQDNs) {
		t.Fatal("The two slices length do not match")
	}
	for i := range FQDNs {
		if FQDNs[i] != config.FQDNs[i] {
			t.Fatal("The two elements do not match")
		}
	}
}

// TestNewNettestIntegration runs a ndt7 nettest
func TestNewNettestIntegration(t *testing.T) {
	config := Config{}
	nettest := NewNettest(context.Background(), config)
	nettest.ASNDatabasePath = "../../asn.mmdb.gz"
	nettest.CountryDatabasePath = "../../country.mmdb.gz"
	nettest.SoftwareName = "measurement-kit"
	nettest.SoftwareVersion = "0.1.0"
	err := nettest.DiscoverAvailableCollectors()
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("AvailableCollectors: %+v", nettest.AvailableCollectors)
	err = nettest.GeoLookup()
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("ProbeIP: %+v", nettest.ProbeIP)
	log.Infof("ProbeASN: %+v", nettest.ProbeASN)
	log.Infof("ProbeCC: %+v", nettest.ProbeCC)
	log.Infof("ProbeNetworkName: %+v", nettest.ProbeNetworkName)
	for err := range nettest.OpenReport() {
		log.Warnf("OpenReport: %+v", err)
	}
	if nettest.Report.ID == "" {
		t.Fatal("OpenReport: failed")
	}
	defer nettest.CloseReport()
	log.Infof("Report: %+v", nettest.Report)
	measurement := nettest.NewMeasurement()
	for ev := range nettest.StartMeasurement("", &measurement) {
		log.Infof("ev: %+v => %+v", ev.Key, ev.Value)
	}
	log.Infof("measurement: %+v", measurement)
	err = nettest.SubmitMeasurement(&measurement)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("measurementID: %+v", measurement.OOID)
}

// TestRunGetserversFailure checks whether Run deals with
// a failure when trying to discover servers.
func TestRunGetServersFailure(t *testing.T) {
	savedFunc := mockableGetservers
	mockedError := errors.New("mocked error")
	mockableGetservers = func(ctx context.Context, config Config) ([]string, error) {
		return nil, mockedError
	}
	ctx := context.Background()
	config := Config{}
	out := make(chan<- model.Event)
	results := run(ctx, config, out)
	if results.GetServersResults.Failure != mockedError.Error() {
		t.Fatal("Unexpected result of getservers")
	}
	mockableGetservers = savedFunc
}

// TestRunRunnerStartDownloadFailure checks whether Run deals with
// a failure when trying to start the download.
func TestRunRunnerStartDownloadFailure(t *testing.T) {
	savedFunc := runnerStartDownload
	mockedError := errors.New("mocked error")
	runnerStartDownload = func(ctx context.Context, FQDN string) (<-chan ndt7model.Measurement, error) {
		return nil, mockedError
	}
	ctx := context.Background()
	config := Config{
		FQDNs: []string{
			"ndd-iupui-mlab4-mil03.measurement-lab.org",
		},
	}
	out := make(chan<- model.Event)
	results := run(ctx, config, out)
	if len(results.DownloadResults) != 1 {
		t.Fatal("Unexpected number of download results")
	}
	if results.DownloadResults[0].Failure != mockedError.Error() {
		t.Fatal("Unexpected error")
	}
	runnerStartDownload = savedFunc
}

// TestRunRunnerStartUploadFailure checks whether Run deals with
// a failure when trying to start the upload.
func TestRunRunnerStartUploadFailure(t *testing.T) {
	savedDloadFunc := runnerStartDownload
	runnerStartDownload = func(ctx context.Context, FQDN string) (<-chan ndt7model.Measurement, error) {
		out := make(chan ndt7model.Measurement)
		defer close(out)
		return out, nil
	}
	savedUploadFunc := runnerStartUpload
	mockedError := errors.New("mocked error")
	runnerStartUpload = func(ctx context.Context, FQDN string) (<-chan ndt7model.Measurement, error) {
		return nil, mockedError
	}
	ctx := context.Background()
	config := Config{
		FQDNs: []string{
			"ndd-iupui-mlab4-mil03.measurement-lab.org",
		},
	}
	out := make(chan<- model.Event)
	results := run(ctx, config, out)
	if len(results.UploadResults) != 1 {
		t.Fatal("Unexpected number of upload results")
	}
	if results.UploadResults[0].Failure != mockedError.Error() {
		t.Fatal("Unexpected error")
	}
	runnerStartUpload = savedUploadFunc
	runnerStartDownload = savedDloadFunc
}
