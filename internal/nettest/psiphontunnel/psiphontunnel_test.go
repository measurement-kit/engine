package psiphontunnel

import (
	"context"
	"testing"

	"github.com/apex/log"
)

func TestNewNettestIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "../../../testdata/psiphon_config.json",
		WorkDirPath:    "/tmp/",
	}
	ctx := context.Background()
	nettest := NewNettest(config)
	nettest.ASNDatabasePath = "../../asn.mmdb"
	nettest.CountryDatabasePath = "../../country.mmdb"
	nettest.SoftwareName = "measurement-kit"
	nettest.SoftwareVersion = "0.1.0"
	err := nettest.DiscoverAvailableCollectors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("AvailableCollectors: %+v", nettest.AvailableCollectors)
	err = nettest.OpenReport(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer nettest.CloseReport(ctx)
	log.Infof("Report: %+v", nettest.Report)
	measurement := nettest.NewMeasurement()
	for ev := range nettest.StartMeasurement(ctx, "", &measurement) {
		log.Infof("ev: %+v => %+v", ev.Key, ev.Value)
	}
	log.Infof("measurement: %+v", measurement)
	err = nettest.SubmitMeasurement(ctx, &measurement)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("measurementID: %+v", measurement.OOID)
}
