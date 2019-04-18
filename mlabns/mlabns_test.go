package mlabns

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// TestGeoOptionsIntegrations uses the GeoOptions policy.
func TestGeoOptionsIntegration(t *testing.T) {
	config := Config{Tool: "ndt_ssl"}
	servers, err := GeoOptions(context.Background(), config)
	for _, server := range servers {
		fmt.Println(server.FQDN)
	}
	if err != nil {
		t.Fatal(err)
	}
}

// TestGeoOptionsURLParseError ensures we deal with URL
// parsing errors.
func TestGeoOptionsURLParseError(t *testing.T) {
	config := Config{
		Tool:    "ndt_ssl",
		BaseURL: "\t", // enough to break URL parsing
	}
	_, err := GeoOptions(context.Background(), config)
	if err == nil {
		t.Fatal("We were expecting an error here")
	}
}

// TestGeoOptionsHTTPError ensures we deal a HTTP error.
func TestGeoOptionsHTTPError(t *testing.T) {
	savedFunc := httpxGET
	mockedError := errors.New("mocked error")
	httpxGET = func(ctx context.Context, URL string) ([]byte, error) {
		return nil, mockedError
	}
	config := Config{Tool: "ndt_ssl"}
	_, err := GeoOptions(context.Background(), config)
	if err != mockedError {
		t.Fatal("Not the error we were expecting")
	}
	httpxGET = savedFunc
}

// TestGeoOptionsJSONError ensures we deal an invalid JSON.
func TestGeoOptionsJSONError(t *testing.T) {
	savedFunc := httpxGET
	httpxGET = func(ctx context.Context, URL string) ([]byte, error) {
		return []byte("{"), nil
	}
	config := Config{Tool: "ndt_ssl"}
	_, err := GeoOptions(context.Background(), config)
	if err == nil {
		t.Fatal("We were expecting an error here")
	}
	httpxGET = savedFunc
}
