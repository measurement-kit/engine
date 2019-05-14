package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/measurement-kit/engine/internal"
	"github.com/measurement-kit/engine/internal/nettest"
	"github.com/measurement-kit/engine/model"
)

// CollectorSubmitResults contains the results of submitting or resubmitting
// a measurement to the OONI collector.
type CollectorSubmitResults struct {
	// Good indicates whether we succeeded or not.
	Good bool

	// UpdatedSerializedMeasurement returns the measurement with updated fields.
	UpdatedSerializedMeasurement string

	// UpdatedReportID returns the updated report ID.
	UpdatedReportID string

	// Logs returns logs useful for debugging.
	Logs string
}

// CollectorSubmitTask is a synchronous task for submitting or resubmitting a
// specific OONI measurement to the OONI collector.
type CollectorSubmitTask struct {
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

// NewCollectorSubmitTask creates a new CollectorSubmitTask with the specified
// software name, software version, and serialized measurement fields.
func NewCollectorSubmitTask(swName, swVersion, measurement string) *CollectorSubmitTask {
	return &CollectorSubmitTask{
		SerializedMeasurement: measurement,
		SoftwareName:          swName,
		SoftwareVersion:       swVersion,
		Timeout:               defaultTimeout,
	}
}

// discoverAvailableCollectors allows to simulate errors in unit tests.
var discoverAvailableCollectors = func(ctx context.Context, nt *nettest.Nettest) error {
	return nt.DiscoverAvailableCollectors(ctx)
}

// submitMeasurement allows to simulate errors in unit tests.
var submitMeasurement = func(ctx context.Context, nt *nettest.Nettest, m *model.Measurement) error {
	return nt.SubmitMeasurement(ctx, m)
}

// jsonMarshal allows to simulate errors in unit tests.
var jsonMarshal = func(m *model.Measurement) ([]byte, error) {
	return json.Marshal(*m)
}

func (t *CollectorSubmitTask) runWithResults(out *CollectorSubmitResults) {
	// Implementation note: here we basically run the normal nettest workflow
	// except that the measurement result is known ahead of time.
	var measurement model.Measurement
	err := json.Unmarshal([]byte(t.SerializedMeasurement), &measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot unmarshal measurement: %s\n", err.Error())
		return
	}
	var nettest nettest.Nettest
	duration, err := internal.MakeTimeout(t.Timeout)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot make duration: %s\n", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	nettest.TestName = measurement.TestName
	nettest.TestVersion = measurement.TestVersion
	nettest.SoftwareName = t.SoftwareName
	nettest.SoftwareVersion = t.SoftwareVersion
	nettest.TestStartTime = measurement.TestStartTime
	err = discoverAvailableCollectors(ctx, &nettest)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot discover collectors: %s\n", err.Error())
		return
	}
	for err := range nettest.OpenReport(ctx) {
		out.Logs += fmt.Sprintf("cannot open report: %s\n", err.Error())
	}
	if nettest.Report.ID == "" {
		out.Logs += fmt.Sprintf("empty report ID, assuming failure\n")
		return
	}
	defer nettest.CloseReport(ctx)
	measurement.ReportID = nettest.Report.ID
	err = submitMeasurement(ctx, &nettest, &measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot submit measurement: %s\n", err.Error())
		return
	}
	data, err := jsonMarshal(&measurement)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot marshal measurement: %s\n", err.Error())
		return
	}
	out.UpdatedSerializedMeasurement = string(data)
	out.UpdatedReportID = measurement.ReportID
	out.Good = true
}

// Run submits (or resubmits) a measurement and returns the results.
func (t *CollectorSubmitTask) Run() *CollectorSubmitResults {
	var out CollectorSubmitResults
	t.runWithResults(&out)
	return &out
}
