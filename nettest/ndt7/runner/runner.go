// Package runner contains the ndt7 nettest runner.
//
// You can use this code through the nettest abstraction, or you can
// instead just use it directly. This documentaion explains how to use
// it directly and describes the emitted events.
//
// Discovering servers
//
// To run a ndt7 nettest you need to discover suitable servers first. To this
// end, use the GetServers function as follows:
//
//     servers, err := runner.GetServers()
//     if err != nil {
//       return
//     }
//
// Note that discovering servers may fail if your IP address is consuming
// too much bandwidth. In such case, the ErrNoAvailableServers error will be
// returned by the GetServers function. Background clients should handle
// this error by retrying after an exponential delay.
//
// Download subtest
//
// To perform a download subtest with a specific server, use
//
//     ch, err := runner.StartDownload(ctx, FQDN)
//     if err != nil {
//       return
//     }
//     for measurement := range ch {
//       // process measurement
//     }
//
// where FQDN is the FQDN of a mlab server. If StartDownload fails, it
// means that we could not connect to the specified server and/or upgrade
// the connection to WebSockets using the ndt7 subprotocol.
//
// Upload subtest
//
// To perform an upload subtest with a specific server, use
//
//     ch, err := runner.StartUpload(ctx, FQDN)
//     if err != nil {
//       return
//     }
//     for measurement := range ch {
//       // process measurement
//     }
//
// where FQDN is the FQDN of a mlab server. If StartUpload fails, it
// means that we could not connect to the specified server and/or upgrade
// the connection to WebSockets using the ndt7 subprotocol.
//
// On success, you must process measurements emitted by the upload.
package runner

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/measurement-kit/engine/mlabns"
	"github.com/measurement-kit/engine/nettest/ndt7/runner/download"
	"github.com/measurement-kit/engine/nettest/ndt7/runner/model"
	"github.com/measurement-kit/engine/nettest/ndt7/runner/spec"
	"github.com/measurement-kit/engine/nettest/ndt7/runner/upload"
)

// ErrNoAvailableServers is returned when there are no available servers. A
// background client should treat this error specially and schedule retrying
// after an exponentially distributed number of seconds.
var ErrNoAvailableServers = errors.New("No available M-Lab servers")

// GetServers gets ndt7 mlab servers using mlabns.
func GetServers(ctx context.Context) ([]string, error) {
	servers, err := mlabns.GeoOptions(ctx, mlabns.Config{
		// TODO(bassosimone): when ndt7 is deployed on the whole platform, we can
		// stop using the staging mlabns service and use the production one.
		BaseURL: "https://locate-dot-mlab-staging.appspot.com/",
		Tool:    "ndt_ssl",
	})
	if err != nil {
		return nil, err
	}
	var FQDNs []string
	for _, server := range servers {
		// TODO(bassosimone): we need to use mlab4 servers only until ndt7
		// is deployed on the whole M-Lab platform.
		if strings.Index(server.FQDN, "-mlab4-") == -1 {
			continue
		}
		FQDNs = append(FQDNs, server.FQDN)
	}
	if len(FQDNs) <= 0 {
		return nil, ErrNoAvailableServers
	}
	return FQDNs, nil
}

// connect establishes a websocket connection.
func connect(ctx context.Context, FQDN, URLPath string) (*websocket.Conn, error) {
	URL := url.URL{}
	URL.Scheme = "wss"
	URL.Host = FQDN
	URL.Path = URLPath
	dialer := websocket.Dialer{}
	headers := http.Header{}
	headers.Add("Sec-WebSocket-Protocol", spec.SecWebSocketProtocol)
	conn, _, err := dialer.DialContext(ctx, URL.String(), headers)
	return conn, err
}

// StartDownload starts the ndt7 download. On success returns a channel where
// measurements are emitted. This channel is closed when the download ends. On
// failure, the error is non nil and you should not attempt using the channel.
func StartDownload(ctx context.Context, FQDN string) (<-chan model.Measurement, error) {
	conn, err := connect(ctx, FQDN, spec.DownloadURLPath)
	if err != nil {
		return nil, err
	}
	ch := make(chan model.Measurement)
	go download.Run(ctx, conn, ch)
	return ch, nil
}

// StartUpload is like StartDownload but for the upload.
func StartUpload(ctx context.Context, FQDN string) (<-chan model.Measurement, error) {
	conn, err := connect(ctx, FQDN, spec.UploadURLPath)
	if err != nil {
		return nil, err
	}
	ch := make(chan model.Measurement)
	go upload.Run(ctx, conn, ch)
	return ch, nil
}
