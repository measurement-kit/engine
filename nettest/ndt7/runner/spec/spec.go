// Package spec contains ndt7 constants. See also the ndt7 spec:
// https://github.com/m-lab/ndt-server/blob/master/spec/ndt7-protocol.md
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
