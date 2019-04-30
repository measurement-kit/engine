// Package collector implements mkall's collector API.
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

// SubmitTask contains settings indicating how to submit or
// resubmit a specific OONI measurement.
type SubmitTask struct {
	// SerializedMeasurement is the measurement to submit.
	SerializedMeasurement string

	// Timeout is the number of seconds after which we abort submitting.
	Timeout int64
}

// SubmitInto is like Submit but takes the results pointers.
func SubmitInto(settings *SubmitTask, out *SubmitResults) {
	// Implementation note: here we basically run the normal nettest workflow
	// except that the measurement result is known ahead of time.
	var measurement model.Measurement
	err := json.Unmarshal([]byte(settings.SerializedMeasurement), &measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot unmarshal measurement: %s\n", err.Error())
		return
	}
	var nettest nettest.Nettest
	duration, err := internal.MakeTimeout(settings.Timeout)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot make duration: %s\n", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	nettest.Ctx = ctx
	nettest.TestName = measurement.TestName
	nettest.TestVersion = measurement.TestVersion
	nettest.SoftwareName = measurement.SoftwareName
	nettest.SoftwareVersion = measurement.SoftwareVersion
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
func Submit(settings *SubmitTask) *SubmitResults {
	var out SubmitResults
	SubmitInto(settings, &out)
	return &out
}
