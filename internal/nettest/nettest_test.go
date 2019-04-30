package nettest

import (
	"context"
	"testing"

	"github.com/measurement-kit/engine/internal/model"
)

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
			model.Service{
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
