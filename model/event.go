package model

import (
	"encoding/json"
	"fmt"
)

// Event is an event emitted by a nettest.
type Event struct {
	// Key is the key that uniquely identifies the event.
	Key string `json:"key"`

	// Value contains event specific variables.
	Value interface{} `json:"value"`
}

// failureMeasurementEvent is a measurement failure
type failureMeasurementEvent struct {
	// Failure is the error that occurred
	Failure string `json:"failure"`

	// Idx is the measurement index
	Idx int64 `json:"idx"`
}

// NewFailureMeasurementEvent creates a new failure measurement event
func NewFailureMeasurementEvent(idx int64, err error) Event {
	return Event{
		Key: "failure.measurement",
		Value: failureMeasurementEvent{
			Failure: err.Error(),
			Idx:     idx,
		},
	}
}

// logEvent is a log event
type logEvent struct {
	// LogLevel is the log level
	LogLevel string `json:"log_level"`

	// Message is the log message
	Message string `json:"message"`
}

// NewLogInfoEvent generates a new log info event
func NewLogInfoEvent(message string) Event {
	return Event{
		Key: "log",
		Value: logEvent{
			LogLevel: "INFO",
			Message:  message,
		},
	}
}

// NewLogWarningEvent generates a new log warning event
func NewLogWarningEvent(err error, message string) Event {
	if err != nil {
		message = fmt.Sprintf("%s: %s", message, err.Error())
	}
	return Event{
		Key: "log",
		Value: logEvent{
			LogLevel: "WARNING",
			Message:  message,
		},
	}
}

// measurementEvent is a measurement event
type measurementEvent struct {
	// JSONStr is the serialized measurement
	JSONStr string `json:"json_str"`
}

// NewMeasurementEvent creates a new measurement event.
func NewMeasurementEvent(measurement Measurement) (Event, error) {
	data, err := json.Marshal(measurement)
	if err != nil {
		return Event{}, err
	}
	return Event{
		Key: "measurement",
		Value: measurementEvent{
			JSONStr: string(data),
		},
	}, nil
}
