package bouncer

import (
	"context"
	"testing"
)

// TestGetCollectorsIntegration just fetches the collectors.
func TestGetCollectorsIntegration(t *testing.T) {
	entries, err := GetCollectors(context.Background(), Config{
		BaseURL: "https://bouncer.ooni.io/",
	})
	t.Logf("%+v", entries)
	if err != nil {
		t.Fatal(err)
	}
}

// TestGetCollectorsFailure checks we deal with a failure.
func TestGetCollectorsFailure(t *testing.T) {
	_, err := GetCollectors(context.Background(), Config{
		BaseURL: "\t", // this should be an invalid URL
	})
	if err == nil {
		t.Fatal("We expected an error here")
	}
}

// TestGetTestHelpersIntegration just fetches the test helpers.
func TestGetTestHelpersIntegration(t *testing.T) {
	entries, err := GetTestHelpers(context.Background(), Config{
		BaseURL: "https://bouncer.ooni.io/",
	})
	t.Logf("%+v", entries)
	if err != nil {
		t.Fatal(err)
	}
}

// TestGetTestHelpersFailure checks we deal with a failure.
func TestGetTestHelpersFailure(t *testing.T) {
	_, err := GetTestHelpers(context.Background(), Config{
		BaseURL: "\t", // this should be an invalid URL
	})
	if err == nil {
		t.Fatal("We expected an error here")
	}
}
