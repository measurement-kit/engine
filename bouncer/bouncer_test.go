package bouncer

import (
	"context"
	"testing"
)

func TestGetCollectorsIntegration(t *testing.T) {
	entries, err := GetCollectors(context.Background(), Config{
		BaseURL: "https://bouncer.ooni.io/",
	})
	t.Logf("%+v", entries)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTestHelpersIntegration(t *testing.T) {
	entries, err := GetTestHelpers(context.Background(), Config{
		BaseURL: "https://bouncer.ooni.io/",
	})
	t.Logf("%+v", entries)
	if err != nil {
		t.Fatal(err)
	}
}
