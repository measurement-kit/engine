// Package download contains ndt7 download code.
package download

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/measurement-kit/engine/model"
	"github.com/measurement-kit/engine/nettest/ndt7/runner/spec"
)

// Run runs the download subtest.
func Run(ctx context.Context, conn *websocket.Conn, ch chan<- model.Event) {
	defer close(ch)
	defer conn.Close()
	wholectx, cancel := context.WithTimeout(ctx, spec.DownloadTimeout)
	defer cancel()
	conn.SetReadLimit(spec.MaxMessageSize)
	for {
		select {
		case <-wholectx.Done():
			return // don't fail the test if we're running for too much time
		default:
			// nothing
		}
		conn.SetReadDeadline(time.Now().Add(spec.IOTimeout))
		mtype, mdata, err := conn.ReadMessage()
		if err != nil {
			return // don't fail the test because of an I/O error
		}
		if mtype != websocket.TextMessage {
			continue
		}
		ch <- model.Event{
			Key:   "ndt7.server_download_measurement",
			Value: spec.EventValue{JSONStr: string(mdata)},
		}
	}
}
