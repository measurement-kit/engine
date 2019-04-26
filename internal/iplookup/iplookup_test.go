package iplookup

import (
	"context"
	"errors"
	"testing"
)

// TestIntegration is a normal IP lookup
func TestIntegration(t *testing.T) {
	IP, err := Perform(context.Background())
	t.Logf("IP: %s", IP)
	if err != nil {
		t.Fatal(err)
	}
}

// TestHTTPFailure deals with the case where HTTP fails
func TestHTTPFailure(t *testing.T) {
	savedFunc := httpxGET
	mockedError := errors.New("mocked error")
	httpxGET = func(ctx context.Context, URL string) ([]byte, error) {
		return nil, mockedError
	}
	IP, err := Perform(context.Background())
	if err != mockedError {
		t.Fatal("Not the error we were expecting")
	}
	if IP != "127.0.0.1" {
		t.Fatal("Not the IP we were expecting")
	}
	httpxGET = savedFunc
}

// TestInvalidXML deals with the case where XML is invalid
func TestInvalidXML(t *testing.T) {
	savedFunc := httpxGET
	httpxGET = func(ctx context.Context, URL string) ([]byte, error) {
		return []byte("<Result"), nil
	}
	IP, err := Perform(context.Background())
	if err == nil {
		t.Fatal("We were expecting a error here")
	}
	if IP != "127.0.0.1" {
		t.Fatal("Not the IP we were expecting")
	}
	httpxGET = savedFunc
}

// TestNotIP deals with the case where we get an invalid IP.
func TestNotIP(t *testing.T) {
	savedFunc := httpxGET
	httpxGET = func(ctx context.Context, URL string) ([]byte, error) {
		return []byte("<Response><Ip>1.2.3</Ip></Response>"), nil
	}
	IP, err := Perform(context.Background())
	if err == nil {
		t.Fatal("We were expecting a error here")
	}
	if IP != "127.0.0.1" {
		t.Fatal("Not the IP we were expecting")
	}
	httpxGET = savedFunc
}
