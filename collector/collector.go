// Package collector implements Measurement Kit's collector API.
//
// See https://github.com/measurement-kit/api#collector.
package collector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/measurement-kit/engine/internal"
	"github.com/measurement-kit/engine/internal/model"
	"github.com/measurement-kit/engine/internal/nettest"
)

// SubmitResults contains the results of submitting or resubmitting
// a measurement to the OONI collector.
type SubmitResults struct {
	// Good indicates whether we succeded or not.
	Good bool

	// UpdatedSerializedMeasurement is the measurement with updated fields.
	UpdatedSerializedMeasurement string

	// UpdatedReportID is the updated report ID.
	UpdatedReportID string

	// Logs contains logs useful for debugging.
	Logs string
}

// SubmitTask is a synchronous task for submitting or resubmitting a
// specific OONI measurement to the OONI collector.
type SubmitTask struct {
	// SerializedMeasurement is the measurement to submit.
	SerializedMeasurement string

	// SoftwareName is the name of the software submitting the measurement.
	SoftwareName string

	// SoftwareVersion is the name of the software submitting the measurement.
	SoftwareVersion string

	// Timeout is the number of seconds after which we abort submitting.
	Timeout int64
}

// defaultTimeout is the default timeout in seconds.
var defaultTimeout int64 = 30

// NewSubmitTask creates a new SubmitTask instance.
func NewSubmitTask(softwareName, softwareVersion, serializedMeasurement string) *SubmitTask {
	return &SubmitTask{
		SerializedMeasurement: serializedMeasurement,
		SoftwareName: softwareName,
		SoftwareVersion: softwareVersion,
		Timeout: defaultTimeout,
	}
}

// SubmitInto is like Submit but takes the results pointers.
func (task *SubmitTask) SubmitInto(out *SubmitResults) {
	// Implementation note: here we basically run the normal nettest workflow
	// except that the measurement result is known ahead of time.
	var measurement model.Measurement
	err := json.Unmarshal([]byte(task.SerializedMeasurement), &measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot unmarshal measurement: %s\n", err.Error())
		return
	}
	var nettest nettest.Nettest
	duration, err := internal.MakeTimeout(task.Timeout)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot make duration: %s\n", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	nettest.Ctx = ctx
	nettest.TestName = measurement.TestName
	nettest.TestVersion = measurement.TestVersion
	nettest.SoftwareName = task.SoftwareName
	nettest.SoftwareVersion = task.SoftwareVersion
	nettest.TestStartTime = measurement.TestStartTime
	err = nettest.DiscoverAvailableCollectors()
	if err != nil {
		out.Logs = fmt.Sprintf("cannot discover collectors: %s\n", err.Error())
		return
	}
	for err := range nettest.OpenReport() {
		out.Logs += fmt.Sprintf("cannot open report: %s\n", err.Error())
	}
	if nettest.Report.ID == "" {
		out.Logs += fmt.Sprintf("empty report ID, assuming failure\n")
		return
	}
	defer nettest.CloseReport()
	measurement.ReportID = nettest.Report.ID
	err = nettest.SubmitMeasurement(&measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot submit measurement: %s\n", err.Error())
		return
	}
	data, err := json.Marshal(measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot marshal measurement: %s\n", err.Error())
		return
	}
	out.UpdatedSerializedMeasurement = string(data)
	out.UpdatedReportID = measurement.ReportID
	out.Good = true
}

// Submit submits (or resubmits) a measurement and returns the results.
func (task *SubmitTask) Submit() *SubmitResults {
	var out SubmitResults
	task.SubmitInto(&out)
	return &out
}
