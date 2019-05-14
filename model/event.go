package model

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
