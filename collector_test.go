package engine

import (
	"fmt"
	"testing"
)

const origMeasurement = `{
	"data_format_version": "0.2.0",
	"input": "torproject.org",
	"measurement_start_time": "2016-06-04 17:53:13",
	"probe_asn": "AS0",
	"probe_cc": "ZZ",
	"probe_ip": "127.0.0.1",
	"software_name": "measurement_kit",
	"software_version": "0.2.0-alpha.1",
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

func TestCollectorResubmitIntegration(t *testing.T) {
	settings := &CollectorResubmitSettings{}
	settings.SerializedMeasurement = origMeasurement
	settings.Timeout = 14
	results := settings.Perform()
	fmt.Println(results.Logs)
	fmt.Println(results.UpdatedSerializedMeasurement)
	if !results.Good {
		t.Fatal("resubmission failed")
	}
}
