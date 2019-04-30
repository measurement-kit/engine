package collector

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

func TestSubmitIntegration(t *testing.T) {
	task := NewSubmitTask("ooniprobe-android", "2.1.0", origMeasurement)
	results := task.Submit()
	fmt.Println(results.Logs)
	fmt.Println(results.UpdatedSerializedMeasurement)
	fmt.Println(results.UpdatedReportID)
	if !results.Good {
		t.Fatal("resubmission failed")
	}
}
