// Package model defines several structures.
package model

// Measurement is a OONI measurement.
//
// This structure is compatible with the definition of the base data format in
// https://github.com/ooni/spec/blob/master/data-formats/df-000-base.md.
type Measurement struct {
	// Annotations contains results annotations
	Annotations map[string]string `json:"annotations,omitempty"`

	// DataFormatVersion is the version of the data format
	DataFormatVersion string `json:"data_format_version"`

	// ID is the locally generated measurement ID
	ID string `json:"id,omitempty"`

	// Input is the measurement input
	Input string `json:"input,omitempty"`

	// InputHashes contains input hashes
	InputHashes []string `json:"input_hashes,omitempty"`

	// MeasurementStartTime is the time when the measurement started
	MeasurementStartTime string `json:"measurement_start_time"`

	// MeasurementRuntime contains the measurement runtime. The JSON name
	// is test_runtime because this is the name expected by the OONI backend
	// even though that name is clearly a misleading one.
	MeasurementRuntime float64 `json:"test_runtime"`

	// OOID is the measurement ID stamped by the OONI collector.
	OOID string `json:"ooid,omitempty"`

	// Options contains command line options
	Options []string `json:"options,omitempty"`

	// ProbeASN contains the probe autonomous system number
	ProbeASN string `json:"probe_asn"`

	// ProbeCC contains the probe country code
	ProbeCC string `json:"probe_cc"`

	// ProbeCity contains the probe city
	ProbeCity string `json:"probe_city,omitempty"`

	// ProbeIP contains the probe IP
	ProbeIP string `json:"probe_ip,omitempty"`

	// ReportID contains the report ID
	ReportID string `json:"report_id"`

	// SoftwareName contains the software name
	SoftwareName string `json:"software_name"`

	// SoftwareVersion contains the software version
	SoftwareVersion string `json:"software_version"`

	// TestHelpers contains the test helpers
	TestHelpers map[string]string `json:"test_helpers,omitempty"`

	// TestKeys contains the real test result. This field is opaque because
	// each nettest will insert here a different structure.
	TestKeys interface{} `json:"test_keys"`

	// TestName contains the test name
	TestName string `json:"test_name"`

	// TestStartTime contains the test start time
	TestStartTime string `json:"test_start_time"`

	// TestVersion contains the test version
	TestVersion string `json:"test_version"`
}
