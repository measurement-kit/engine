// Package upload contains ndt7 upload code
package upload

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/measurement-kit/engine/model"
	"github.com/measurement-kit/engine/nettest/ndt7/runner/spec"
)

// makePreparedMessage generates a prepared message that should be sent
// over the network for generating network load.
func makePreparedMessage(size int) (*websocket.PreparedMessage, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	data := make([]byte, size)
	// This is not the fastest algorithm to generate a random string, yet it
	// is most likely good enough for our purposes. See [1] for a comprehensive
	// discussion regarding how to generate a random string in Golang.
	//
	// .. [1] https://stackoverflow.com/a/31832326/4354461
	//
	// Also, the ndt7 specification does not require us to use this algorithm
	// and we could send purely random data as well. We're sending textual data
	// here just because in some debugging cases it's easier to read.
	for i := range data {
		data[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return websocket.NewPreparedMessage(websocket.BinaryMessage, data)
}

// ignoreIncoming ignores any incoming message.
func ignoreIncoming(conn *websocket.Conn) {
	conn.SetReadLimit(spec.MaxMessageSize)
	for {
		// Implementation note: this guarantees that the websocket engine
		// is processing messages. Here we're using as timeout the timeout
		// for the whole upload, so that we know that this goroutine is
		// active for most of the time we care about, even in the case in
		// which the server is not sending us any messages.
		conn.SetReadDeadline(time.Now().Add(spec.UploadTimeout))
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// appInfo contains an application level measurement. This message is
// consistent with v0.7.0 of the ndt7 spec.
type appInfo struct {
	// NumBytes is the number of bytes transferred so far.
	NumBytes int64 `json:"num_bytes"`
}

// The measurement struct contains measurement results. This message is
// consistent with v0.7.0 of the ndt7 spec.
type measurement struct {
	// Elapsed is the number of seconds elapsed since the beginning.
	Elapsed float64 `json:"elapsed"`

	// AppInfo contains application level measurements.
	AppInfo *appInfo `json:"app_info,omitempty"`
}

// emit emits an event during the upload.
func emit(ch chan<- model.Event, elapsed float64, numBytes int64) {
	measurement := measurement{
		Elapsed: elapsed,
		AppInfo: &appInfo{
			NumBytes: numBytes,
		},
	}
	data, err := json.Marshal(measurement)
	if err != nil {
		return // should not happen in practice as all fields are serializable
	}
	ch <- model.Event{
		Key:   "ndt7.client_upload_measurement",
		Value: spec.EventValue{JSONStr: string(data)},
	}
}

// upload runs the upload and emits progress on ch.
func upload(ctx context.Context, conn *websocket.Conn, out chan<- int64) {
	defer close(out)
	wholectx, cancel := context.WithTimeout(ctx, spec.UploadTimeout)
	defer cancel()
	preparedMessage, err := makePreparedMessage(spec.BulkMessageSize)
	if err != nil {
		return // I believe this should not happen in practice
	}
	var total int64
	for {
		select {
		case <-wholectx.Done():
			return // time to stop uploading
		default:
			// nothing
		}
		conn.SetWriteDeadline(time.Now().Add(spec.IOTimeout))
		if err := conn.WritePreparedMessage(preparedMessage); err != nil {
			return // just bail if we cannot write
		}
		total += spec.BulkMessageSize
		out <- total
	}
}

// uploadAsync runs the upload and returns a channel where progress is emitted.
func uploadAsync(ctx context.Context, conn *websocket.Conn) <-chan int64 {
	out := make(chan int64)
	go upload(ctx, conn, out)
	return out
}

// Run runs the upload subtest.
func Run(ctx context.Context, conn *websocket.Conn, ch chan<- model.Event) {
	defer close(ch)
	defer conn.Close()
	go ignoreIncoming(conn)
	t0 := time.Now()
	prev := t0
	for tot := range uploadAsync(ctx, conn) {
		now := time.Now()
		if now.Sub(prev) > spec.UpdateInterval {
			emit(ch, now.Sub(t0).Seconds(), tot)
			prev = now
		}
	}
}
