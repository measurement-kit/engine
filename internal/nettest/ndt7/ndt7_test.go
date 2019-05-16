package ndt7

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

// TestIntegration runs a ndt7 nettest.
func TestIntegration(t *testing.T) {
	ctx := context.Background()
	nettest := NewNettest()
	err := nettest.DiscoverAvailableCollectors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = nettest.OpenReport(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer nettest.CloseReport(ctx)
	measurement := nettest.NewMeasurement()
	for ev := range nettest.StartMeasurement(ctx, "", &measurement) {
		t.Logf("%+v", ev)
	}
	data, err := json.MarshalIndent(measurement, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	t.Logf("%s", string(data))
	err = nettest.SubmitMeasurement(ctx, &measurement)
	if err != nil {
		t.Fatal(err)
	}
}
