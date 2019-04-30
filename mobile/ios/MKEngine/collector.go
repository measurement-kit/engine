// +build ios

package MKEngine

import "github.com/measurement-kit/engine/collector"

// Contains the results of resubmitting a report.
type CollectorResubmitResults struct {
	rr collector.ResubmitResults
}

// Returns whether the resubmission succeded.
func (rr *CollectorResubmitResults) Good() bool {
	return rr.rr.Good
}

// Returns the updated measurement.
func (rr *CollectorResubmitResults) UpdatedSerializedMeasurement() string {
	return rr.rr.UpdatedSerializedMeasurement
}

// Returns the updated report ID.
func (rr *CollectorResubmitResults) UpdatedReportID() string {
	return rr.rr.UpdatedReportID
}

// Returns logs generated by resubmitting a report.
func (rr *CollectorResubmitResults) Logs() string {
	return rr.rr.Logs
}

// Contains the settings for resubmitting a report.
type CollectorResubmitSettings struct {
	rs collector.ResubmitSettings
}

// Sets the measurement to resubmit.
func (rs *CollectorResubmitSettings) SetSerializedMeasurement(value string) {
	rs.rs.SerializedMeasurement = value
}

// Sets the timeout (in seconds) for the resubmission.
func (rs *CollectorResubmitSettings) SetTimeout(value int64) {
	rs.rs.Timeout = value
}

// The task for resubmitting measurements.
type CollectorResubmitTask struct {
}

// Runs the resubmission task with specific settings and returns the results.
func (rt *CollectorResubmitTask) Run(rs *CollectorResubmitSettings) *CollectorResubmitResults {
	var rr CollectorResubmitResults
	collector.ResubmitInto(&rs.rs, &rr.rr)
	return &rr
}