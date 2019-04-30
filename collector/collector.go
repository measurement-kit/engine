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

// ResubmitResults contains the results of resubmitting
// a measurement to the OONI collector.
type ResubmitResults struct {
	// Good indicates whether we succeded or not.
	Good bool

	// UpdatedSerializedMeasurement is the measurement with updated fields.
	UpdatedSerializedMeasurement string

	// UpdatedReportID is the updated report ID.
	UpdatedReportID string

	// Logs contains logs useful for debugging.
	Logs string
}

// ResubmitSettings contains settings indicating how to
// resubmit a specific OONI measurement.
type ResubmitSettings struct {
	// SerializedMeasurement is the measurement to resubmit.
	SerializedMeasurement string

	// Timeout is the number of seconds after which we abort resubmitting.
	Timeout int64
}

// ResubmitInto is like resubmit but takes the results as a pointer.
func ResubmitInto(settings *ResubmitSettings, out *ResubmitResults) {
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

// Resubmit resubmits a measurement and returns the results.
func Resubmit(settings *ResubmitSettings) *ResubmitResults {
	var out ResubmitResults
	ResubmitInto(settings, &out)
	return &out
}
