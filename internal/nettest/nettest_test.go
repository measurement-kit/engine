package nettest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/measurement-kit/engine/internal/collector"
	"github.com/measurement-kit/engine/internal/model"
)

// TestGetAvailableBouncersDefault ensures that we can set available
// default bouncers and we'll get them back.
func TestGetAvailableBouncersDefault(t *testing.T) {
	svc := model.Service{
		Address: "a",
		Type:    "b",
	}
	nettest := &Nettest{
		Ctx:               context.Background(),
		AvailableBouncers: []model.Service{svc},
	}
	svcs := nettest.getAvailableBouncers()
	if len(svcs) != 1 {
		t.Fatal("Unexpected returned value length")
	}
	if svcs[0] != svc {
		t.Fatal("The returned value is wrong")
	}
}

// TestDiscoverAvailableCollectorsIntegration discovers available collectors.
func TestDiscoverAvailableCollectorsIntegration(t *testing.T) {
	nettest := &Nettest{
		Ctx: context.Background(),
	}
	err := nettest.DiscoverAvailableCollectors()
	if err != nil {
		t.Fatal(err)
	}
}

// TestDiscoverAvailableCollectorsFailure deals with the case
// where we cannot discover any available collector.
func TestDiscoverAvailableCollectorsFailure(t *testing.T) {
	nettest := &Nettest{
		Ctx: context.Background(),
		AvailableBouncers: []model.Service{
			{
				Address: "httpo://42q7ug46dspcsvkw.onion",
				Type:    "onion",
			},
		},
	}
	err := nettest.DiscoverAvailableCollectors()
	if err == nil {
		t.Fatal("We expected a failure here")
	}
}

// TestDiscoverAvailableCollectorsQueryFailure deals with the case
// where we fail in querying the bouncer.
func TestDiscoverAvailableCollectorsQueryFailure(t *testing.T) {
	nettest := &Nettest{
		Ctx: context.Background(),
		AvailableBouncers: []model.Service{
			{
				Address: "\t", // fail b/c URL is invalid
				Type:    "https",
			},
		},
	}
	err := nettest.DiscoverAvailableCollectors()
	if err == nil {
		t.Fatal("We expected a failure here")
	}
}

// TestDiscoverAvailableTestHelpersIntegration discovers available test helpers.
func TestDiscoverAvailableTestHelpersIntegration(t *testing.T) {
	nettest := &Nettest{
		Ctx: context.Background(),
	}
	err := nettest.DiscoverAvailableTestHelpers()
	if err != nil {
		t.Fatal(err)
	}
}

// TestDiscoverAvailableTestHelpersFailure deals with the case
// where we cannot discover any available test helper.
func TestDiscoverAvailableTestHelpersFailure(t *testing.T) {
	nettest := &Nettest{
		Ctx: context.Background(),
		AvailableBouncers: []model.Service{
			{
				Address: "httpo://42q7ug46dspcsvkw.onion",
				Type:    "onion",
			},
		},
	}
	err := nettest.DiscoverAvailableTestHelpers()
	if err == nil {
		t.Fatal("We expected a failure here")
	}
}

// TestDiscoverAvailableTestHelpersQueryFailure deals with the case
// where we fail in querying the bouncer.
func TestDiscoverAvailableTestHelpersQueryFailure(t *testing.T) {
	nettest := &Nettest{
		Ctx: context.Background(),
		AvailableBouncers: []model.Service{
			{
				Address: "\t", // fail b/c URL is invalid
				Type:    "https",
			},
		},
	}
	err := nettest.DiscoverAvailableTestHelpers()
	if err == nil {
		t.Fatal("We expected a failure here")
	}
}

// TestOpenReportIntegration opens a report.
func TestOpenReportIntegration(t *testing.T) {
	nettest := &Nettest{
		Ctx:             context.Background(),
		ProbeASN:        "AS0",
		ProbeCC:         "ZZ",
		SoftwareName:    "MKEngine",
		SoftwareVersion: "0.0.1",
		TestName:        "dummy",
		TestVersion:     "0.0.1",
		AvailableCollectors: []model.Service{
			{
				Address: "https://b.collector.ooni.io",
				Type:    "https",
			},
		},
	}
	for err := range nettest.OpenReport() {
		if err != nil {
			t.Log(err)
		}
	}
	if nettest.Report.ID == "" {
		t.Fatal("OpenReport: failed")
	}
}

// TestOpenReportMultipleTimes ensures that we cannot
// open an already openned report.
func TestOpenReportMultipleTimes(t *testing.T) {
	nettest := &Nettest{
		Ctx:             context.Background(),
		ProbeASN:        "AS0",
		ProbeCC:         "ZZ",
		SoftwareName:    "MKEngine",
		SoftwareVersion: "0.0.1",
		TestName:        "dummy",
		TestVersion:     "0.0.1",
		AvailableCollectors: []model.Service{
			{
				Address: "https://b.collector.ooni.io",
				Type:    "https",
			},
		},
	}
	for err := range nettest.OpenReport() {
		if err != nil {
			t.Log(err)
		}
	}
	if nettest.Report.ID == "" {
		t.Fatal("OpenReport: failed")
	}
	reportID := nettest.Report.ID
	for err := range nettest.OpenReport() {
		if err != nil {
			t.Log(err)
		}
	}
	if nettest.Report.ID != reportID {
		t.Fatal("OpenReport: changed the report ID")
	}
}

// TestOpenReportNoCollector deals with the case where we
// fail when opening the report because there is no collector
// that we know how to handle.
func TestOpenReportNoCollector(t *testing.T) {
	nettest := &Nettest{
		Ctx: context.Background(),
		AvailableCollectors: []model.Service{
			{
				Address: "httpo://42q7ug46dspcsvkw.onion",
				Type:    "onion",
			},
		},
	}
	for err := range nettest.OpenReport() {
		if err != nil {
			t.Log(err)
		}
	}
	if nettest.Report.ID != "" {
		t.Fatal("OpenReport: we expected a failure here")
	}
}

// TestOpenReportCollectorOpenError deals with the case
// where we fail when opening the report because there are
// errors when attempting to open a report.
func TestOpenReportCollectorOpenError(t *testing.T) {
	nettest := &Nettest{
		Ctx: context.Background(),
		AvailableCollectors: []model.Service{
			{
				Address: "\t", // fail b/c URL is invalid
				Type:    "https",
			},
		},
	}
	for err := range nettest.OpenReport() {
		if err != nil {
			t.Log(err)
		}
	}
	if nettest.Report.ID != "" {
		t.Fatal("OpenReport: we expected a failure here")
	}
}

// TestNewMeasurementWorks helps us to gain some confidence
// that NewMeasurement is working as intended.
func TestNewMeasurementWorks(t *testing.T) {
	nettest := &Nettest{
		ProbeIP:  "130.192.91.211",
		ProbeASN: "AS30722",
		ProbeCC:  "IT",
		Report: collector.Report{
			ID: "1234567",
		},
		SoftwareName:    "ooniprobe-android",
		SoftwareVersion: "1.2.4",
		TestName:        "antani",
		TestStartTime:   FormatTimeNowUTC(),
		TestVersion:     "4.2.1",
	}
	m := nettest.NewMeasurement()
	if m.DataFormatVersion != "0.2.0" {
		t.Fatal("invalid DateFormatVersion")
	}
	_, err := time.Parse(DateFormat, m.MeasurementStartTime)
	if err != nil {
		t.Fatal("the serialized MeasurementStartTime is invalid")
	}
	if m.ProbeIP != "127.0.0.1" {
		t.Fatal("invalid ProbeIP")
	}
	if m.ProbeASN != nettest.ProbeASN {
		t.Fatal("invalid ProbeASN")
	}
	if m.ProbeCC != nettest.ProbeCC {
		t.Fatal("invalid ProbeCC")
	}
	if m.SoftwareName != nettest.SoftwareName {
		t.Fatal("invalid SoftwareName")
	}
	if m.SoftwareVersion != nettest.SoftwareVersion {
		t.Fatal("invalid SoftwareVersion")
	}
	if m.TestName != nettest.TestName {
		t.Fatal("invalid TestName")
	}
	if m.TestStartTime != nettest.TestStartTime {
		t.Fatal("invalid TestStartTime")
	}
	if m.TestVersion != nettest.TestVersion {
		t.Fatal("invalid TestVersion")
	}
}

func measurementLifecycle(t *testing.T, expectedErr error) {
	nettest := &Nettest{
		Ctx:             context.Background(),
		SoftwareName:    "ooniprobe-mocked",
		SoftwareVersion: "1.2.4",
		TestName:        "antani",
		TestStartTime:   FormatTimeNowUTC(),
		TestVersion:     "4.2.1",
	}
	err := nettest.DiscoverAvailableCollectors()
	if err != nil {
		t.Fatal(err)
	}
	for err := range nettest.OpenReport() {
		if err != nil {
			t.Fatal(err)
		}
	}
	if nettest.Report.ID == "" {
		t.Fatal("cannot open report")
	}
	m := nettest.NewMeasurement()
	m.TestKeys = struct{}{}
	m.MeasurementRuntime = 1.17
	m.Input = ""
	err = nettest.SubmitMeasurement(&m)
	if err != expectedErr {
		t.Fatalf("SubmitMeasurement did not return: %+v", expectedErr)
	}
	err = nettest.CloseReport()
	if err != nil {
		t.Fatal(err)
	}
}

// TestMeasurementLifecycle runs the lifecycle of a measurement.
func TestMeasurementLifecycle(t *testing.T) {
	measurementLifecycle(t, nil)
}

// TestMeasurementSubmitMeasurementError runs the measurement
// lifecycle and makes sure we handle a measurement submission error.
func TestMeasurementSubmitMeasurementError(t *testing.T) {
	savedFunc := updateReport
	mockedError := errors.New("mocked error")
	updateReport = func(ctx context.Context, r *collector.Report, m *model.Measurement) (string, error) {
		return "", mockedError
	}
	measurementLifecycle(t, mockedError)
	updateReport = savedFunc
}
