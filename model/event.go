package model

// Event is an event emitted by a nettest.
type Event struct {
	// Key is the key that uniquely identifies the event.
	Key string `json:"key"`

	// Value contains event specific variables.
	Value interface{} `json:"value"`
}
