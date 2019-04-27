package collector

import (
	"context"
	"errors"
	"testing"

	"github.com/measurement-kit/engine/internal/model"
)

type fakeTestKeys struct {
	ClientResolver string `json:"client_resolver"`
}

func makeMeasurement(rt ReportTemplate, ID string) model.Measurement {
	return model.Measurement{
		DataFormatVersion:    "0.2.0",
		ID:                   "bdd20d7a-bba5-40dd-a111-9863d7908572",
		MeasurementStartTime: "2018-11-01 15:33:20",
		ProbeASN:             rt.ProbeASN,
		ProbeCC:              rt.ProbeCC,
		ReportID:             ID,
		SoftwareName:         rt.SoftwareName,
		SoftwareVersion:      rt.SoftwareVersion,
		TestKeys: fakeTestKeys{
			ClientResolver: "91.80.37.104",
		},
		TestName:           rt.TestName,
		MeasurementRuntime: 5.0565230846405,
		TestStartTime:      "2018-11-01 15:33:17",
		TestVersion:        rt.TestVersion,
	}
}

// TestIntegration submits a measurement.
func TestIntegration(t *testing.T) {
	config := Config{
		BaseURL: "https://b.collector.ooni.io",
	}
	template := ReportTemplate{
		ProbeASN:        "AS0",
		ProbeCC:         "ZZ",
		SoftwareName:    "measurement-kit",
		SoftwareVersion: "0.0.1",
		TestName:        "dummy",
		TestVersion:     "0.0.1",
	}
	ctx := context.Background()
	report, err := Open(ctx, config, template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Report %s: OPEN", report.ID)
	defer func() {
		err := report.Close(ctx)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Report %s: CLOSED", report.ID)
	}()
	ID, err := report.Update(ctx, makeMeasurement(template, report.ID))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Report %s: UPDATED (Measurement ID: %s)", report.ID, ID)
}

// TestOpenJSONMarshalError verifies that we deal with
// JSON marshalling errors in Open.
func TestOpenJSONMarshalError(t *testing.T) {
	config := Config{}
	template := ReportTemplate{}
	savedJSONMarshal := jsonMarshal
	mockedError := errors.New("mocked error")
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, mockedError
	}
	ctx := context.Background()
	_, err := Open(ctx, config, template)
	if err == nil {
		t.Fatal("We did not expect a success here")
	}
	jsonMarshal = savedJSONMarshal
}

// TestOpenHTTPError verifies that Open deals with HTTP errors.
func TestOpenHTTPError(t *testing.T) {
	config := Config{
		BaseURL: "\t", // should be enough to fail httpx
	}
	template := ReportTemplate{}
	_, err := Open(context.Background(), config, template)
	if err == nil {
		t.Fatal("We expected an error here")
	}
}

// TestUpdateJSONMarshalError verifies that we deal with
// JSON marshalling errors in Update.
func TestUpdateJSONMarshalError(t *testing.T) {
	savedJSONMarshal := jsonMarshal
	mockedError := errors.New("mocked error")
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, mockedError
	}
	ctx := context.Background()
	var r Report
	_, err := r.Update(ctx, model.Measurement{})
	if err != mockedError {
		t.Fatal("Not the error we were expecting")
	}
	jsonMarshal = savedJSONMarshal
}

// TestUpdateHTTPError verifies that Update deals with HTTP errors.
func TestUpdateHTTPError(t *testing.T) {
	r := Report{
		Conf: Config{
			BaseURL: "\t", // should be enough to fail httpx
		},
	}
	_, err := r.Update(context.Background(), model.Measurement{})
	if err == nil {
		t.Fatal("We expected an error here")
	}
}

// TestUpdateJSONUnmarshalError verifies that we deal with
// JSON unmarshalling errors in Update.
func TestUpdateJSONUnmarshalError(t *testing.T) {
	savedFunc := httpxPOSTWithBaseURL
	httpxPOSTWithBaseURL = func(ctx context.Context, baseURL, path, contentType string, body []byte) ([]byte, error) {
		return []byte("{"), nil // this is not valid JSON
	}
	ctx := context.Background()
	var r Report
	_, err := r.Update(ctx, model.Measurement{})
	if err == nil {
		t.Fatal("We expected an error here")
	}
	httpxPOSTWithBaseURL = savedFunc
}
