package psiphontunnel

import (
	"context"
	"testing"
)

func TestNewNettestIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	nettest := NewNettest(context.Background(), config)
	nettest.ASNDatabasePath = "../../asn.mmdb.gz"
	nettest.CountryDatabasePath = "../../country.mmdb.gz"
	nettest.SoftwareName = "measurement-kit"
	nettest.SoftwareVersion = "0.1.0"
	err := nettest.DiscoverAvailableCollectors()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("AvailableCollectors: %+v", nettest.AvailableCollectors)
	err = nettest.GeoLookup()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ProbeIP: %+v", nettest.ProbeIP)
	t.Logf("ProbeASN: %+v", nettest.ProbeASN)
	t.Logf("ProbeCC: %+v", nettest.ProbeCC)
	t.Logf("ProbeNetworkName: %+v", nettest.ProbeNetworkName)
	err = nettest.OpenReport()
	if err != nil {
		t.Fatal(err)
	}
	defer nettest.CloseReport()
	t.Logf("Report: %+v", nettest.Report)
	measurement := nettest.NewMeasurement()
	for range nettest.StartMeasurement("", &measurement) {
		// nothing
	}
	t.Logf("measurement: %+v", measurement)
	err = nettest.SubmitMeasurement(&measurement)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("measurementID: %+v", measurement.OOID)
}
