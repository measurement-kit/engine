package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/measurement-kit/engine/internal"
	"github.com/measurement-kit/engine/internal/model"
	"github.com/measurement-kit/engine/internal/nettest"
)

// CollectorSubmitResults contains the results of submitting or resubmitting
// a measurement to the OONI collector.
type CollectorSubmitResults struct {
	good                         bool
	updatedSerializedMeasurement string
	updatedReportID              string
	logs                         string
}

// Good returns whether we succeded or not.
func (r *CollectorSubmitResults) Good() bool {
	return r.good
}

// Logs returns logs useful for debugging.
func (r *CollectorSubmitResults) Logs() string {
	return r.logs
}

// UpdatedReportID returns the updated report ID.
func (r *CollectorSubmitResults) UpdatedReportID() string {
	return r.updatedReportID
}

// UpdatedSerializedMeasurement returns the measurement with updated fields.
func (r *CollectorSubmitResults) UpdatedSerializedMeasurement() string {
	return r.updatedSerializedMeasurement
}

// CollectorSubmitTask is a synchronous task for submitting or resubmitting a
// specific OONI measurement to the OONI collector.
type CollectorSubmitTask struct {
	serializedMeasurement string
	softwareName          string
	softwareVersion       string
	timeout               int64
}

// defaultTimeout is the default timeout in seconds.
var defaultTimeout int64 = 30

// NewCollectorSubmitTask creates a new CollectorSubmitTask with the specified
// software name, software version, and serialized measurement fields.
func NewCollectorSubmitTask(swName, swVersion, measurement string) *CollectorSubmitTask {
	return &CollectorSubmitTask{
		serializedMeasurement: measurement,
		softwareName:          swName,
		softwareVersion:       swVersion,
		timeout:               defaultTimeout,
	}
}

// SetSerializedMeasurement sets the measurement to submit.
func (t *CollectorSubmitTask) SetSerializedMeasurement(measurement string) {
	t.serializedMeasurement = measurement
}

// SetSoftwareName sets the name of the software submitting the measurement.
func (t *CollectorSubmitTask) SetSoftwareName(softwareName string) {
	t.softwareName = softwareName
}

// SetSoftwareVersion sets the name of the software submitting the measurement.
func (t *CollectorSubmitTask) SetSoftwareVersion(softwareVersion string) {
	t.softwareVersion = softwareVersion
}

// SetTimeout sets the number of seconds after which we abort submitting.
func (t *CollectorSubmitTask) SetTimeout(timeout int64) {
	t.timeout = timeout
}

// discoverAvailableCollectors allows to simulate errors in unit tests.
var discoverAvailableCollectors = func(nt *nettest.Nettest) error {
	return nt.DiscoverAvailableCollectors()
}

// submitMeasurement allows to simulate errors in unit tests.
var submitMeasurement = func(nt *nettest.Nettest, m *model.Measurement) error {
	return nt.SubmitMeasurement(m)
}

// jsonMarshal allows to simulate errors in unit tests.
var jsonMarshal = func(m *model.Measurement) ([]byte, error) {
	return json.Marshal(*m)
}

func (t *CollectorSubmitTask) runWithResults(out *CollectorSubmitResults) {
	// Implementation note: here we basically run the normal nettest workflow
	// except that the measurement result is known ahead of time.
	var measurement model.Measurement
	err := json.Unmarshal([]byte(t.serializedMeasurement), &measurement)
	if err != nil {
		out.logs = fmt.Sprintf("cannot unmarshal measurement: %s\n", err.Error())
		return
	}
	var nettest nettest.Nettest
	duration, err := internal.MakeTimeout(t.timeout)
	if err != nil {
		out.logs = fmt.Sprintf("cannot make duration: %s\n", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	nettest.Ctx = ctx
	nettest.TestName = measurement.TestName
	nettest.TestVersion = measurement.TestVersion
	nettest.SoftwareName = t.softwareName
	nettest.SoftwareVersion = t.softwareVersion
	nettest.TestStartTime = measurement.TestStartTime
	err = discoverAvailableCollectors(&nettest)
	if err != nil {
		out.logs = fmt.Sprintf("cannot discover collectors: %s\n", err.Error())
		return
	}
	for err := range nettest.OpenReport() {
		out.logs += fmt.Sprintf("cannot open report: %s\n", err.Error())
	}
	if nettest.Report.ID == "" {
		out.logs += fmt.Sprintf("empty report ID, assuming failure\n")
		return
	}
	defer nettest.CloseReport()
	measurement.ReportID = nettest.Report.ID
	err = submitMeasurement(&nettest, &measurement)
	if err != nil {
		out.logs = fmt.Sprintf("cannot submit measurement: %s\n", err.Error())
		return
	}
	data, err := jsonMarshal(&measurement)
	if err != nil {
		out.logs = fmt.Sprintf("cannot marshal measurement: %s\n", err.Error())
		return
	}
	out.updatedSerializedMeasurement = string(data)
	out.updatedReportID = measurement.ReportID
	out.good = true
}

// Run submits (or resubmits) a measurement and returns the results.
func (t *CollectorSubmitTask) Run() *CollectorSubmitResults {
	var out CollectorSubmitResults
	t.runWithResults(&out)
	return &out
}