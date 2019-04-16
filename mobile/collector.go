// Package mobile contains the mobile API
package mobile

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/measurement-kit/engine/mobile/internal"
	"github.com/measurement-kit/engine/model"
	"github.com/measurement-kit/engine/nettest"
)

// MKECollectorResubmitResults contains the results of resubmitting
// a measurement to the OONI collector.
type MKECollectorResubmitResults struct {
	// Good indicates whether we succeded or not.
	Good bool

	// UpdatedSerializedMeasurement is the measurement with updated fields.
	UpdatedSerializedMeasurement string

	// Logs contains logs useful for debugging.
	Logs string
}

// MKECollectorResubmitSettings contains settings indicating how to
// resubmit a specific OONI measurement.
type MKECollectorResubmitSettings struct {
	// SerializedMeasurement is the measurement to resubmit.
	SerializedMeasurement string

	// Timeout is the number of seconds after which we abort resubmitting.
	Timeout int64
}

// Perform resubmits a measurement and returns the results.
func (x *MKECollectorResubmitSettings) Perform() *MKECollectorResubmitResults {
	// Implementation note: here we basically run the normal nettest workflow
	// except that the measurement result is known ahead of time.
	var out MKECollectorResubmitResults
	var measurement model.Measurement
	err := json.Unmarshal([]byte(x.SerializedMeasurement), &measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot unmarshal measurement: %s\n", err.Error())
		return &out
	}
	var nettest nettest.Nettest
	duration, err := internal.MakeTimeout(x.Timeout)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot make duration: %s\n", err.Error())
		return &out
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
		return &out
	}
	err = nettest.OpenReport()
	if err != nil {
		out.Logs = fmt.Sprintf("cannot open report: %s\n", err.Error())
		return &out
	}
	defer nettest.CloseReport()
	measurement.ReportID = nettest.Report.ID
	err = nettest.SubmitMeasurement(&measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot submit measurement: %s\n", err.Error())
		return &out
	}
	data, err := json.Marshal(measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot marshal measurement: %s\n", err.Error())
		return &out
	}
	out.UpdatedSerializedMeasurement = string(data)
	out.Good = true
	return &out
}