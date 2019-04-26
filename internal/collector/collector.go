// Package collector contains a OONI collector client implementation.
//
// Specifically we implement v2.0.0 of the OONI collector specification defined
// in https://github.com/ooni/spec/blob/master/backends/bk-003-collector.md.
package collector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/measurement-kit/engine/internal/httpx"
	"github.com/measurement-kit/engine/internal/model"
)

// Config contains the collector configuration
type Config struct {
	// BaseURL is the collector base URL
	BaseURL string
}

// ReportTemplate is the template for opening a report
type ReportTemplate struct {
	// ProbeASN is the probe's autonomous system number (e.g. `AS1234`)
	ProbeASN string `json:"probe_asn"`

	// ProbeCC is the probe's country code (e.g. `IT`)
	ProbeCC string `json:"probe_cc"`

	// SoftwareName is the app name (e.g. `measurement-kit`)
	SoftwareName string `json:"software_name"`

	// SoftwareVersion is the app version (e.g. `0.9.1`)
	SoftwareVersion string `json:"software_version"`

	// TestName is the test name (e.g. `ndt`)
	TestName string `json:"test_name"`

	// TestVersion is the test version (e.g. `1.0.1`)
	TestVersion string `json:"test_version"`
}

// Report is an open report
type Report struct {
	// ID is the report ID
	ID string `json:"report_id"`

	// Conf is the configuration being used
	Conf Config
}

// jsonMarshal allows to mock json.Marshal in tests
var jsonMarshal = json.Marshal

func open(ctx context.Context, conf Config, rt ReportTemplate) (Report, error) {
	report := Report{Conf: conf}
	requestData, err := jsonMarshal(rt)
	if err != nil {
		return report, err
	}
	responseData, err := httpx.POSTWithBaseURL(
		ctx, conf.BaseURL, "/report", "application/json", requestData,
	)
	if err != nil {
		return report, fmt.Errorf("request with body '%s' has failed: %s",
			string(requestData), err.Error(),
		)
	}
	err = json.Unmarshal(responseData, &report)
	if err != nil {
		return report, fmt.Errorf(
			"cannot parse JSON returned by server: %s",
			err.Error(),
		)
	}
	return report, err
}

// Open opens a new report. Returns the report on success; an error on failure.
func Open(ctx context.Context, conf Config, rt ReportTemplate) (Report, error) {
	report, err := open(ctx, conf, rt)
	if err != nil {
		return report, fmt.Errorf(
			"Opening report failed: %s", err.Error(),
		)
	}
	return report, nil
}

type updateRequest struct {
	// Format is the data format
	Format string `json:"format"`

	// Content is the actual report
	Content interface{} `json:"content"`
}

type updateResponse struct {
	// ID is the measurement ID
	ID string `json:"measurement_id"`
}

// httpxPOSTWithBaseURL simplifies life in unit tests
var httpxPOSTWithBaseURL = httpx.POSTWithBaseURL

// Update updates a report by appending a new measurement to it.
//
// Returns the measurement ID on success; an error on failure.
func (r Report) Update(ctx context.Context, m model.Measurement) (string, error) {
	ureq := updateRequest{
		Format:  "json",
		Content: m,
	}
	data, err := jsonMarshal(ureq)
	if err != nil {
		return "", err
	}
	data, err = httpxPOSTWithBaseURL(
		ctx, r.Conf.BaseURL, fmt.Sprintf("/report/%s", r.ID),
		"application/json", data,
	)
	if err != nil {
		return "", err
	}
	var ures updateResponse
	err = json.Unmarshal(data, &ures)
	if err != nil {
		return "", err
	}
	return ures.ID, nil
}

// Close closes the report. Returns nil on success; an error on failure.
func (r Report) Close(ctx context.Context) error {
	_, err := httpx.POSTWithBaseURL(
		ctx, r.Conf.BaseURL, fmt.Sprintf("/report/%s/close", r.ID), "", nil,
	)
	return err
}
