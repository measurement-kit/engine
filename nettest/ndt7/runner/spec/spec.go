// Package spec contains ndt7 constants and data structures.
package spec

import (
	"time"
)

// SecWebSocketProtocol is the value of the Sec-WebSocket-Protocol header.
const SecWebSocketProtocol = "net.measurementlab.ndt.v7"

// MaxMessageSize is the maximum accepted message size.
const MaxMessageSize = 1 << 20

// DownloadTimeout is the time after which the download must stop.
const DownloadTimeout = 15 * time.Second

// IOTimeout is the timeout for I/O operations.
const IOTimeout = 7 * time.Second

// DownloadURLPath is the URL path used for the download.
const DownloadURLPath = "/ndt/v7/download"

// UploadURLPath is the URL path used for the download.
const UploadURLPath = "/ndt/v7/upload"

// UploadTimeout is the time after which the upload must stop.
const UploadTimeout = 10 * time.Second

// BulkMessageSize is the size of uploader messages
const BulkMessageSize = 1 << 13

// UpdateInterval is the interval between client side upload measurements.
const UpdateInterval = 250 * time.Millisecond

// EventValue is the value of a ndt7 event. This is the model.Event.Value value
// that will be emitted by ndt7-generated nettest events.
type EventValue struct {
	// JSONStr is a serialized JSON message. See the documentation of the
	// github.com/measurement-kit/engine/nettest/ndt7/runner package for a
	// description of the possible values of this field.
	JSONStr string `json:"json_str"`
}
