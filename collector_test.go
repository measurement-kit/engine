package engine

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/measurement-kit/engine/internal/model"
	"github.com/measurement-kit/engine/internal/nettest"
)

const origMeasurement = `{
	"data_format_version": "0.2.0",
	"input": "torproject.org",
	"measurement_start_time": "2016-06-04 17:53:13",
	"probe_asn": "AS0",
	"probe_cc": "ZZ",
	"probe_ip": "127.0.0.1",
	"software_name": "ooniprobe-android",
	"software_version": "2.0.0",
	"test_keys": {
		"failure": null,
		"received": [],
		"sent": []
	},
	"test_name": "tcp_connect",
	"test_runtime": 0.253494024276733,
	"test_start_time": "2016-06-04 17:53:13",
	"test_version": "0.0.1"
}`

// TestCollectorSubmitIntegration covers the common case of submitting
// a measurement to the OONI collector.
func TestCollectorSubmitIntegration(t *testing.T) {
	task := NewCollectorSubmitTask("ooniprobe-android", "2.1.0", origMeasurement)
	results := task.Run()
	fmt.Println(results.Logs)
	fmt.Println(results.UpdatedSerializedMeasurement)
	fmt.Println(results.UpdatedReportID)
	if !results.Good {
		t.Fatal("resubmission failed")
	}
}

// TestCollectorSubmitExpectFailureV200 covers the case where we want
// a measurement to fail because the client is ooniprobe-android v2.0.0.
func TestCollectorSubmitExpectFailureV200(t *testing.T) {
	task := NewCollectorSubmitTask("ooniprobe-android", "2.0.0", origMeasurement)
	results := task.Run()
	fmt.Println(results.Logs)
	fmt.Println(results.UpdatedSerializedMeasurement)
	fmt.Println(results.UpdatedReportID)
	if results.Good {
		t.Fatal("we expected a failure here")
	}
}

// TestCollectorSubmitConstructorAndSetters ensures that we can use
// either the constructor or the setters to configure the task.
func TestCollectorSubmitConstructorAndSetters(t *testing.T) {
	task := NewCollectorSubmitTask("ooniprobe-android", "2.0.0", origMeasurement)
	if task.SoftwareName != "ooniprobe-android" {
		t.Fatal("the constructor cannot set the softwareName")
	}
	if task.SoftwareVersion != "2.0.0" {
		t.Fatal("the constructor cannot set the softwareVersion")
	}
	if task.SerializedMeasurement != origMeasurement {
		t.Fatal("the constructor cannot set the serializedMeasurement")
	}
	if task.Timeout != defaultTimeout {
		t.Fatal("the constructor cannot set the timeout")
	}
}

// TestCollectorSubmitUnmarshalError covers the case where we're
// passed an invalid serialized JSON.
func TestCollectorSubmitUnmarshalError(t *testing.T) {
	task := NewCollectorSubmitTask("ooniprobe-android", "2.1.0", "{")
	results := task.Run()
	if results.Good {
		t.Fatal("We expected a failure here")
	}
}

// TestCollectorSubmitInvalidTimeout covers the case where we're
// passed an invalid timeout value.
func TestCollectorSubmitInvalidTimeout(t *testing.T) {
	task := NewCollectorSubmitTask("ooniprobe-android", "2.1.0", origMeasurement)
	task.Timeout = -1
	results := task.Run()
	if results.Good {
		t.Fatal("We expected a failure here")
	}
}

// TestCollectorSubmitDiscoverFailure covers the case where there
// is a failure when discovering available collectors.
func TestCollectorSubmitDiscoverFailure(t *testing.T) {
	savedFunc := discoverAvailableCollectors
	discoverAvailableCollectors = func(ctx context.Context, nt *nettest.Nettest) error {
		return errors.New("mocked error")
	}
	task := NewCollectorSubmitTask("ooniprobe-android", "2.1.0", origMeasurement)
	results := task.Run()
	if results.Good {
		t.Fatal("We expected a failure here")
	}
	discoverAvailableCollectors = savedFunc
}

// TestCollectorSubmitSubmitFailure covers the case where there
// is a failure when submitting the actual measurement.
func TestCollectorSubmitSubmitFailure(t *testing.T) {
	savedFunc := submitMeasurement
	submitMeasurement = func(ctx context.Context, nt *nettest.Nettest, m *model.Measurement) error {
		return errors.New("mocked error")
	}
	task := NewCollectorSubmitTask("ooniprobe-android", "2.1.0", origMeasurement)
	results := task.Run()
	if results.Good {
		t.Fatal("We expected a failure here")
	}
	submitMeasurement = savedFunc
}

// TestCollectorSubmitMarshalFailure covers the case where there
// is a failure when marshalling the updated measurement.
func TestCollectorSubmitMarshalFailure(t *testing.T) {
	savedFunc := jsonMarshal
	jsonMarshal = func(m *model.Measurement) ([]byte, error) {
		return nil, errors.New("mocked error")
	}
	task := NewCollectorSubmitTask("ooniprobe-android", "2.1.0", origMeasurement)
	results := task.Run()
	if results.Good {
		t.Fatal("We expected a failure here")
	}
	jsonMarshal = savedFunc
}
