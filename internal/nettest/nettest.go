// Package nettest contains code for running nettests.
//
// This API is such that every small operation that a test must perform
// is a separate operation. This allows you to handle errors and results
// of each separate operation in the way you find most convenient.
//
// This is an internal API that you should use to implement a nettest.
//
// Creating a nettest
//
// Instantiate a nettest variable using:
//
//     var nettest nettest.Nettest
//
// In such case, you MUST fill in the following fields:
//
// - nettest.TestName with the name of the nettest
//
// - nettest.TestVersion with the nettest version
//
// - nettest.SoftwareName with the app name
//
// - nettest.SoftwareVersion with the app version
//
// - nettest.TestStartTime with the UTC test start time formatted according
// to the format expected by OONI (you can use nettest.FormatTimeNowUTC to
// initialize this field with the current UTC time, or nettest.DateFormat to
// format another time according to the proper format -- just remember
// that you must use the UTC time here)
//
// - nettest.Main if the nettest is written in Go, otherwise, if the nettest
// is in a foreign language, you'll need to call corresponding code to get back
// the test keys, as described below.
//
// For example
//
//     nettest.TestName = "nettest"
//     nettest.TestVersion = "0.0.1"
//     nettest.SoftwareName = "example"
//     nettest.SoftwareVersion = "0.0.1"
//     nettest.TestStartTime = nettest.FormatTimeNowUTC()
//     nettest.Main = func(
//       ctx context.Context, input string, m *model.Measurement,
//       ch chan<- model.Event,
//     ) {
//       // perform measurement and initialize m with results
//     }
//
// Configuring specific bouncers
//
// The bouncer is used to discover collectors and test helpers. If you
// don't have a specific bouncer in mind, just skip this step and we'll
// use the default bouncer. Otherwise, do something like:
//
//     nettest.AvailableBouncers = []model.Service{
//       model.Service{
//         Type: "https",
//         Address: "https://bouncer.example.com",
//       },
//     }
//
// Add as many bouncers as you wish. Currently, only HTTPS bouncers
// are supported. We'll try them in order and use the first one
// that successfully returns us a valid response.
//
// Discovering collectors
//
// We recommend you to automatically discover collectors. Otherwise
// just initialize nettest.AvailableCollectors.
//
// To automatically discover collectors do the following:
//
//     err := nettest.DiscoverAvailableCollectors(ctx)
//     if err != nil {
//       return
//     }
//
// This will populate the nettest.AvailableCollectors field.
//
// Collectors will be automatically discovered using the OONI
// bouncer. You can set the nettest.AvailableBouncers field to
// force the code to use one, or few, specific bouncers.
//
// Alternatively, just populate nettest.AvailableCollectors manually.
//
// Discovering test helpers
//
// If your test needs test helpers, you should discover the available
// test helpers using:
//
//     err = nettest.DiscoverAvailableTestHelpers(ctx)
//     if err != nil {
//       return
//     }
//
// This will populate the nettest.AvailableTestHelpers field.
//
// Test helpers will be automatically discovered using the OONI
// bouncer. You can set the nettest.AvailableBouncers field to
// force the code to use one, or few, specific bouncers.
//
// Alternatively, just populate nettest.AvailableTestHelpers manually.
//
// Geolocation
//
// Geolocating a probe means discover its IP, CC (country code),
// ASN (autonomous system number), and network name (i.e. the
// commercial name bound to the ASN).
//
// If you already know this values, just initialize them; e.g.:
//
//     nettest.ProbeIP = "93.147.252.33"
//     nettest.ProbeCC = "IT"
//     nettest.ProbeASN = "AS30722"
//     nettest.ProbeNetworkName = "Vodafone Italia"
//
// Otherwise, you need to initialize the CountryDatabasePath and
// the ASNDatabasePath fields to point to valid and current MaxMind
// MMDB databases; e.g.,
//
//     nettest.CountryDatabasePath = "country.mmdb"
//     nettest.ASNDatabasePath = "asn.mmdb"
//
// Then run:
//
//     err = nettest.GeoLookup(ctx)
//     if err != nil {
//       return
//     }
//
// This will fill the nettest.Probe{IP,ASN,CC,NetworkName} fields. On
// error they will be initialized, respectively, to "127.0.0.1", "AS0",
// "ZZ", and "". Not setting the country and/or the ASN database path
// will cause GeoLookup to fail and return an error.
//
// Resolver lookup
//
// The resolver lookup step discovers the resolver IP. Run:
//
//     err = nettest.ResolverLookup(ctx)
//     if err != nil {
//       return
//     }
//
// The result will be saved in nettest.ResolverIP. In case of
// error this function will set the ResolverIP to 127.0.0.1.
//
// Opening a report
//
// This is required to submit measurements to a collector. Run:
//
//     err := range nettest.OpenReport(ctx) {
//     if err != nil {
//       return
//     }
//     defer nettest.CloseReport(ctx)
//
// This will attempt to open a report with all the available collectors
// and returns whether we succeeded or not. Note that, if the report has
// already been openned, this function will be a no-op.
//
// Creating a new measurement
//
// You are now ready to perform measurements. Ask the nettest to
// create for you a measurement with:
//
//     measurement := nettest.NewMeasurement()
//
// This will initialize all measurement fields except:
//
// - measurement.TestKeys, which should contains a JSON serializable
// interface{} containing the nettest specific results
//
// - measurement.MeasurementRuntime, which should contain the measurement
// runtime in seconds as a floating point
//
// - measurement.Input, which should only be initialized if your
// nettest requires input
//
// If nettest.Main is initialized, as it should be the case for all
// nettests created using a factory function, you can perform a
// measurement for a specific input and fill the above measurement
// fields by using:
//
//     for ev := range nettest.StartMeasurement(ctx, input, &measurement) {
//       // handle nettest generated events
//     }
//     // nettest done; you most likely want to inspect measurement
//
// where input is an empty string if the nettest does not take any
// input. Otherwise, you'll need to call the (possibly foreign)
// nettest specific code to get the test keys and initialize Runtime
// and Input yourself. Either way, when you're done, you can submit
// the measurement to the configured collector.
//
// Note that, by default, the ProbeIP in the measurement will be set
// to "127.0.0.1". If you want to submit the real probe IP, you'll
// need to override measurement.ProbeIP with nettest.ProbeIP manually
// before submitting the measurement.
//
// Submitting a measurement
//
// To submit a measurement, run:
//
//     err := nettest.SubmitMeasurement(ctx, &measurement)
//     if err != nil {
//       return
//     }
//
// Note that you need to have opened a report before, otherwise
// we will not know where to submit the measurement.
//
// If successful, this will set the measurement.OOID field, which
// may be empty if the collector does not support if. If this field
// isn't empty, later you can use this OOID to get the (possibly
// post processed) measurement from the OONI API.
package nettest

import (
	"context"
	"errors"
	"time"

	"github.com/measurement-kit/engine/internal/bouncer"
	"github.com/measurement-kit/engine/internal/collector"
	"github.com/measurement-kit/engine/model"
)

// DateFormat is the format used by OONI for dates inside reports.
const DateFormat = "2006-01-02 15:04:05"

// FormatTimeNowUTC formats the current time in UTC using the OONI format.
func FormatTimeNowUTC() string {
	return time.Now().UTC().Format(DateFormat)
}

// MainFunc is the measurement main. The first argument (ctx) is the context
// that you want to use for measuring. The second argument (input) is the
// input of the measurement. A nettest that does not take any input expects you
// to pass an empty string here. The third argument (mstub) is the stub
// measurement, partially initialized by the nettest. The runner implementation
// MUST fill all the fields that are not initialized by NewMeasurement (see
// above for a complete list of such fields). The fourth argument (ch) is a
// channel where the nettest should post asynchronous events. The runner
// MUST NOT close the channel as it is managed by nettest.StartMeasurement.
type MainFunc = func(
	ctx context.Context,
	input string,
	mstub *model.Measurement,
	ch chan<- model.Event,
)

// Nettest is a nettest.
type Nettest struct {
	// TestName is the test name.
	TestName string

	// TestVersion is the test version.
	TestVersion string

	// SoftwareName contains the software name.
	SoftwareName string

	// SoftwareVersion contains the software version.
	SoftwareVersion string

	// TestStartTime is the UTC time when the test started.
	TestStartTime string

	// Main is the main function of the measurement
	Main MainFunc

	// AvailableBouncers contains all the available bouncers.
	AvailableBouncers []model.Service

	// AvailableCollectors contains all the available collectors.
	AvailableCollectors []model.Service

	// AvailableTestHelpers contains all the available test helpers.
	AvailableTestHelpers map[string][]model.Service

	// CountryDatabasePath contains the country MMDB database path.
	CountryDatabasePath string

	// ASNDatabasePath contains the ASN MMDB database path.
	ASNDatabasePath string

	// ProbeIP contains the probe IP.
	ProbeIP string

	// ProbeASN contains the probe ASN.
	ProbeASN string

	// ProbeCC contains the probe CC.
	ProbeCC string

	// ProbeNetworkName contains the probe network name.
	ProbeNetworkName string

	// ResolverIP is the resolver's IP.
	ResolverIP string

	// Report is the report bound to this nettest.
	Report collector.Report
}

// getAvailableBouncers always returns one or more available bouncers. If the
// user configured AvailableBouncers, this is what we return. Otherwise, we
// create one or more entry represting out default choices.
func (nettest *Nettest) getAvailableBouncers() []model.Service {
	if len(nettest.AvailableBouncers) > 0 {
		return nettest.AvailableBouncers
	}
	return []model.Service{
		{
			Address: "https://bouncer.ooni.io",
			Type:    "https",
		},
	}
}

// DiscoverAvailableCollectors discovers the available collectors.
func (nettest *Nettest) DiscoverAvailableCollectors(ctx context.Context) error {
	for _, e := range nettest.getAvailableBouncers() {
		if e.Type != "https" {
			continue
		}
		collectors, err := bouncer.GetCollectors(ctx, bouncer.Config{
			BaseURL: e.Address,
		})
		if err != nil {
			continue
		}
		nettest.AvailableCollectors = collectors
		return nil
	}
	return errors.New("Cannot discover available collectors")
}

// DiscoverAvailableTestHelpers discovers the available test helpers.
func (nettest *Nettest) DiscoverAvailableTestHelpers(ctx context.Context) error {
	for _, e := range nettest.getAvailableBouncers() {
		if e.Type != "https" {
			continue
		}
		testHelpers, err := bouncer.GetTestHelpers(ctx, bouncer.Config{
			BaseURL: e.Address,
		})
		if err != nil {
			continue
		}
		nettest.AvailableTestHelpers = testHelpers
		return nil
	}
	return errors.New("Cannot discover available test helpers")
}

// ErrNoDatabasesPath indicates that the MMDB databases path are not specified.
var ErrNoDatabasesPath = errors.New("unspecified ASN and/or country path")

// GeoLookup performs the geolookup (probe_ip, probe_asn, etc.)
func (nettest *Nettest) GeoLookup(ctx context.Context) error {
	nettest.ProbeIP = "127.0.0.1"
	nettest.ProbeASN = "AS0"
	nettest.ProbeCC = "ZZ"
	nettest.ProbeNetworkName = ""
	return errors.New("Not implemented")
}

// ResolverLookup discovers the resolver's IP address.
func (nettest *Nettest) ResolverLookup(ctx context.Context) error {
	return errors.New("Not implemented")
}

// OpenReport opens a new report for the nettest.
func (nettest *Nettest) OpenReport(ctx context.Context) error {
	if nettest.Report.ID != "" {
		return nil
	}
	for _, e := range nettest.AvailableCollectors {
		if e.Type != "https" {
			continue
		}
		report, err := collector.Open(ctx, collector.Config{
			BaseURL: e.Address,
		}, collector.ReportTemplate{
			ProbeASN:        nettest.ProbeASN,
			ProbeCC:         nettest.ProbeCC,
			SoftwareName:    nettest.SoftwareName,
			SoftwareVersion: nettest.SoftwareVersion,
			TestName:        nettest.TestName,
			TestVersion:     nettest.TestVersion,
		})
		if err == nil {
			nettest.Report = report
			return nil
		}
	}
	return errors.New("Cannot open report: all collectors failed")
}

// NewMeasurement returns a new measurement for this nettest. You should
// fill fields that are not initialized; see above for a description
// of what fields WILL NOT be initialized.
func (nettest *Nettest) NewMeasurement() model.Measurement {
	return model.Measurement{
		DataFormatVersion:    "0.2.0",
		MeasurementStartTime: time.Now().UTC().Format(DateFormat),
		ProbeIP:              "127.0.0.1", // override if you want to submit it
		ProbeASN:             nettest.ProbeASN,
		ProbeCC:              nettest.ProbeCC,
		ReportID:             nettest.Report.ID,
		SoftwareName:         nettest.SoftwareName,
		SoftwareVersion:      nettest.SoftwareVersion,
		TestName:             nettest.TestName,
		TestStartTime:        nettest.TestStartTime,
		TestVersion:          nettest.TestVersion,
	}
}

// StartMeasurement starts the measurement in a background goroutine. The
// input argument is the input required by the nettest. If the nettest does
// not take any input, use an empty string. The measurement argument is a
// measurement returned by NewMeasurement. The code running the measurement
// will initialize all the fields that NewMeasurement didn't initialize. The
// code will post asynchronous events on the returned channel. The channel
// will be closed when the measurement is complete. This function will cause
// a panic if the nettest.Main field is not initialized.
func (nettest *Nettest) StartMeasurement(
	ctx context.Context,
	input string,
	measurement *model.Measurement,
) <-chan model.Event {
	outch := make(chan model.Event)
	go func() {
		defer close(outch)
		nettest.Main(ctx, input, measurement, outch)
	}()
	return outch
}

// updateReport allows to inject errors in tests
var updateReport = func(
	ctx context.Context, r *collector.Report, m *model.Measurement,
) (string, error) {
	return r.Update(ctx, *m)
}

// SubmitMeasurement submits a measurement to the selected collector. It is
// safe to call this function from different goroutines concurrently as long
// as the measurement is not shared by the goroutines.
func (nettest *Nettest) SubmitMeasurement(
	ctx context.Context, measurement *model.Measurement,
) error {
	measurementID, err := updateReport(ctx, &nettest.Report, measurement)
	if err != nil {
		return err
	}
	measurement.OOID = measurementID
	return nil
}

// CloseReport closes an open report.
func (nettest *Nettest) CloseReport(ctx context.Context) error {
	return nettest.Report.Close(ctx)
}
