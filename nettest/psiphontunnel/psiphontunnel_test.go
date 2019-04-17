package psiphontunnel

import (
	"context"
	"testing"

	"github.com/apex/log"
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
	log.Infof("AvailableCollectors: %+v", nettest.AvailableCollectors)
	err = nettest.GeoLookup()
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("ProbeIP: %+v", nettest.ProbeIP)
	log.Infof("ProbeASN: %+v", nettest.ProbeASN)
	log.Infof("ProbeCC: %+v", nettest.ProbeCC)
	log.Infof("ProbeNetworkName: %+v", nettest.ProbeNetworkName)
	err = nettest.OpenReport()
	if err != nil {
		t.Fatal(err)
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
