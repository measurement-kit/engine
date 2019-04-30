package mkapi

import (
	"github.com/measurement-kit/engine/collector"
)

// CollectorSubmitResults contains the results of CollectorSubmitTask.
type CollectorSubmitResults interface {
	// SyncTaskResults is the base interface.
	SyncTaskResults

	// UpdatedReportID returns the updated report ID.
	UpdatedReportID() string

	// UpdatedSerializedMeasurement returns the updated measurement.
	UpdatedSerializedMeasurement() string
}

type goCollectorSubmitResults struct {
	r *collector.SubmitResults
}

func (r *goCollectorSubmitResults) Good() bool {
	return r.r.Good
}

func (r *goCollectorSubmitResults) Logs() string {
	return r.r.Logs
}

func (r *goCollectorSubmitResults) UpdatedReportID() string {
	return r.r.UpdatedReportID
}

func (r *goCollectorSubmitResults) UpdatedSerializedMeasurement() string {
	return r.r.UpdatedSerializedMeasurement
}

// CollectorSubmitTask is a sync task for submitting a measurement.
type CollectorSubmitTask interface {
	// SyncTask is the base interface.
	SyncTask

	// SetSerializedMeasurement sets the measurement to submit.
	SetSerializedMeasurement(measurement string)

	// SetSoftwareName sets the name of the software submitting the measurement.
	SetSoftwareName(softwareName string)

	// SetSoftwareVersion sets the version of the software submitting the measurement.
	SetSoftwareVersion(softwareVersion string)

	// Run runs the task until completion and returns the results.
	Run() CollectorSubmitResults
}

type goCollectorSubmitTask struct {
	t *collector.SubmitTask
}

func (t *goCollectorSubmitTask) SetSerializedMeasurement(measurement string) {
	t.t.SerializedMeasurement = measurement
}

func (t *goCollectorSubmitTask) SetSoftwareName(softwareName string) {
	t.t.SoftwareName = softwareName
}

func (t *goCollectorSubmitTask) SetSoftwareVersion(softwareVersion string) {
	t.t.SoftwareVersion = softwareVersion
}

func (t *goCollectorSubmitTask) SetTimeout(timeout int64) {
	t.t.Timeout = timeout
}

func (t *goCollectorSubmitTask) Run() CollectorSubmitResults {
	return &goCollectorSubmitResults{t.t.Submit()}
}
