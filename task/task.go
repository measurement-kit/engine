// Package task implements mkall's task API.
package task

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/measurement-kit/engine/internal/assets"
	"github.com/measurement-kit/engine/internal/model"
	"github.com/measurement-kit/engine/internal/nettest"
	"github.com/measurement-kit/engine/internal/nettest/ndt7"
	"github.com/measurement-kit/engine/internal/nettest/psiphontunnel"
)

type taskAccounting struct {
	DownloadedKB float64 `json:"downloaded_kb"`
	Failure      string  `json:"failure"`
	UploadedKB   float64 `json:"uploaded_kb"`
}

// Handle contains the task internals.
type Handle struct {
	accounting taskAccounting
	cancel     context.CancelFunc
	ctx        context.Context
	done       int64
	marshal    func(interface{}) ([]byte, error)
	out        chan interface{}
}

var semaphore sync.Mutex

// Start starts a new task.
func Start(settings string) *Handle {
	handle := &Handle{}
	handle.ctx, handle.cancel = context.WithCancel(context.Background())
	handle.out = make(chan interface{})
	if os.Getenv("MK_EVENT_PRETTY") == "1" {
		handle.marshal = func(d interface{}) ([]byte, error) {
			return json.MarshalIndent(d, "", "  ")
		}
	} else {
		handle.marshal = json.Marshal
	}
	go func() {
		handle.out <- model.Event{Key: "status.queued", Value: struct{}{}}
		semaphore.Lock()
		defer semaphore.Unlock()
		handle.run(settings)
		handle.out <- model.Event{
			Key:   "status.end",
			Value: handle.accounting,
		}
		close(handle.out)
	}()
	return handle
}

// IsDone indicates whether the task is done.
func (handle *Handle) IsDone() bool {
	return atomic.LoadInt64(&handle.done) != 0
}

// WaitForNextEvent blocks until the task emits its next event.
func (handle *Handle) WaitForNextEvent() string {
	something, ok := <-handle.out
	if !ok {
		atomic.StoreInt64(&handle.done, 1)
		return `{"key":"status.terminated","value":{}}`
	}
	data, err := handle.marshal(something)
	if err == nil {
		return string(data)
	}
	failure := model.Event{
		Key: "bug.json_dump",
		Value: taskFailure{
			Failure: err.Error(),
		},
	}
	data, err = json.Marshal(failure)
	if err == nil {
		return string(data)
	}
	return `{"key":"bug.json_dump","value":{Failure:"generic_error"}}`
}

// Interrupt interrupts a running task.
func (handle *Handle) Interrupt() {
	handle.cancel()
}

type taskOptions struct {
	ConfigFilePath   string `json:"config_file_path"`
	NoBouncer        bool   `json:"no_bouncer"`
	NoCollector      bool   `json:"no_collector"`
	NoGeoLookup      bool   `json:"no_geolookup"`
	NoResolverLookup bool   `json:"no_resolver_lookup"`
	SoftwareName     string `json:"software_name"`
	SoftwareVersion  string `json:"software_version"`
	WorkDirPath      string `json:"work_dir_path"`
}

type taskSettings struct {
	Inputs  []string    `json:"inputs"`
	Name    string      `json:"name"`
	Options taskOptions `json:"options"`
}

type taskFailure struct {
	Failure string `json:"failure"`
}

type taskLog struct {
	LogLevel   string  `json:"log_level,omitempty"`
	Message    string  `json:"message"`
	Percentage float64 `json:"percentage,omitempty"`
}

type taskMeasurement struct {
	Failure string `json:"failure,omitempty"`
	Idx     int    `json:"idx"`
	Input   string `json:"input"`
	JSONStr string `json:"json_str,omitempty"`
}

func (handle *Handle) run(s string) {
	var settings taskSettings
	err := json.Unmarshal([]byte(s), &settings)
	if err != nil {
		handle.out <- model.Event{
			Key: "failure.startup",
			Value: taskFailure{
				Failure: err.Error(),
			},
		}
		return
	}

	if settings.Options.SoftwareName == "" {
		handle.out <- model.Event{
			Key: "failure.startup",
			Value: taskFailure{
				Failure: "empty_variable: software_name",
			},
		}
		return
	}
	if settings.Options.SoftwareVersion == "" {
		handle.out <- model.Event{
			Key: "failure.startup",
			Value: taskFailure{
				Failure: "empty_variable: software_version",
			},
		}
		return
	}
	if settings.Options.WorkDirPath == "" {
		handle.out <- model.Event{
			Key: "failure.startup",
			Value: taskFailure{
				Failure: "empty_variable: work_dir_path",
			},
		}
		return
	}

	var nt *nettest.Nettest
	if settings.Name == "PsiphonTunnel" {
		nt = psiphontunnel.NewNettest(handle.ctx, psiphontunnel.Config{
			ConfigFilePath: settings.Options.ConfigFilePath,
			WorkDirPath:    settings.Options.WorkDirPath,
		})
		settings.Inputs = []string{""} // run exactly once

	} else if settings.Name == "Ndt7" {
		nt = ndt7.NewNettest(handle.ctx, ndt7.Config{})
		settings.Inputs = []string{""} // run exactly once
	}

	if nt == nil {
		handle.out <- model.Event{
			Key: "failure.startup",
			Value: taskFailure{
				Failure: "unknown_nettest_error",
			},
		}
		return
	}

	nt.SoftwareName = settings.Options.SoftwareName
	nt.SoftwareVersion = settings.Options.SoftwareVersion

	if !settings.Options.NoBouncer {
		err = nt.DiscoverAvailableCollectors()
		if err != nil {
			handle.out <- model.Event{
				Key: "log",
				Value: taskLog{
					Message:  fmt.Sprintf("discover_collector_error: %s", err.Error()),
					LogLevel: "WARNING",
				},
			}
			// FALLTHROUGH
		}
		handle.out <- model.Event{
			Key:   "status.available_collectors",
			Value: nt.AvailableCollectors,
		}
		err = nt.DiscoverAvailableTestHelpers()
		if err != nil {
			handle.out <- model.Event{
				Key: "log",
				Value: taskLog{
					Message:  fmt.Sprintf("discover_test_helpers_error: %s", err.Error()),
					LogLevel: "WARNING",
				},
			}
			// FALLTHROUGH
		}
		handle.out <- model.Event{
			Key:   "status.available_test_helpers",
			Value: nt.AvailableTestHelpers,
		}
	}
	handle.out <- model.Event{
		Key: "status.progress",
		Value: taskLog{
			Percentage: 0.1,
			Message:    "contacted bouncer",
		},
	}

	if !settings.Options.NoGeoLookup {
		err = assets.Download(handle.ctx, settings.Options.WorkDirPath)
		if err != nil {
			handle.out <- model.Event{
				Key: "failure.startup",
				Value: taskFailure{
					Failure: fmt.Sprintf("download_assets_error: %s", err.Error()),
				},
			}
			return
		}
		nt.ASNDatabasePath = assets.ASNDatabasePath(settings.Options.WorkDirPath)
		nt.CountryDatabasePath = assets.CountryDatabasePath(settings.Options.WorkDirPath)
		err = nt.GeoLookup()
		if err != nil {
			handle.out <- model.Event{
				Key: "log",
				Value: taskLog{
					Message:  fmt.Sprintf("geolookup_error: %s", err.Error()),
					LogLevel: "WARNING",
				},
			}
			// FALLTHROUGH
		}
		handle.out <- model.Event{
			Key: "status.geoip_lookup",
			Value: struct {
				ProbeIP          string `json:"probe_ip"`
				ProbeASN         string `json:"probe_asn"`
				ProbeCC          string `json:"probe_cc"`
				ProbeNetworkName string `json:"probe_network_name"`
			}{
				nt.ProbeIP,
				nt.ProbeASN,
				nt.ProbeCC,
				nt.ProbeNetworkName,
			},
		}
	}
	handle.out <- model.Event{
		Key: "status.progress",
		Value: taskLog{
			Percentage: 0.2,
			Message:    "geoip lookup",
		},
	}

	if !settings.Options.NoResolverLookup {
		err = nt.ResolverLookup()
		if err != nil {
			handle.out <- model.Event{
				Key: "log",
				Value: taskLog{
					Message:  fmt.Sprintf("resolver_lookup_error: %s", err.Error()),
					LogLevel: "WARNING",
				},
			}
			// FALLTHROUGH
		}
		handle.out <- model.Event{
			Key: "status.resolver_lookup",
			Value: struct {
				ResolverIP string `json:"resolver_ip"`
			}{
				nt.ResolverIP,
			},
		}
	}
	handle.out <- model.Event{
		Key: "status.progress",
		Value: taskLog{
			Percentage: 0.3,
			Message:    "resolver lookup",
		},
	}

	if !settings.Options.NoCollector {
		for err := range nt.OpenReport() {
			handle.out <- model.Event{
				Key: "log",
				Value: taskLog{
					Message:  fmt.Sprintf("open_report_error: %s", err.Error()),
					LogLevel: "WARNING",
				},
			}
			// FALLTHROUGH
		}
		if nt.Report.ID != "" {
			defer nt.CloseReport()
			handle.out <- model.Event{
				Key: "status.report_create",
				Value: struct {
					ReportID string `json:"report_id"`
				}{
					nt.Report.ID,
				},
			}
		} else {
			handle.out <- model.Event{
				Key: "failure.report_create",
				Value: taskFailure{
					Failure: "sequential_operation_error",
				},
			}
		}
	}
	handle.out <- model.Event{
		Key: "status.progress",
		Value: taskLog{
			Percentage: 0.4,
			Message:    "open report",
		},
	}

	for idx, input := range settings.Inputs {
		handle.out <- model.Event{
			Key: "status.measurement_start",
			Value: taskMeasurement{
				Idx:   idx,
				Input: input,
			},
		}

		measurement := nt.NewMeasurement()
		for ev := range nt.StartMeasurement(input, &measurement) {
			handle.out <- ev
		}
		jsonbytes, err := json.Marshal(measurement)
		if err != nil {
			handle.out <- model.Event{
				Key: "bug.json_dump",
				Value: taskFailure{
					Failure: err.Error(),
				},
			}
			continue
		}
		jsonstr := string(jsonbytes)
		handle.out <- model.Event{
			Key: "measurement",
			Value: taskMeasurement{
				Idx:     idx,
				Input:   input,
				JSONStr: jsonstr,
			},
		}

		if nt.Report.ID != "" {
			err = nt.SubmitMeasurement(&measurement)
			if err != nil {
				handle.out <- model.Event{
					Key: "failure.measurement_submission",
					Value: taskMeasurement{
						Failure: err.Error(),
						Idx:     idx,
						Input:   input,
						JSONStr: jsonstr,
					},
				}
			} else {
				handle.out <- model.Event{
					Key: "status.measurement_submission",
					Value: taskMeasurement{
						Idx:   idx,
						Input: input,
					},
				}
			}
		}

		handle.out <- model.Event{
			Key: "status.measurement_done",
			Value: taskMeasurement{
				Idx:   idx,
				Input: input,
			},
		}
		handle.out <- model.Event{
			Key: "status.progress",
			Value: taskLog{
				Percentage: 0.4 + float64(idx)/float64(len(settings.Inputs))/0.6,
				Message:    fmt.Sprintf("measured input: '%s'", input),
			},
		}
	}
}
