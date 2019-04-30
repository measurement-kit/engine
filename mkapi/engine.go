// Package mkapi implements the Measurement Kit API.
//
// This API is designed for generating mobile libraries. As such, we do not
// follow some Go best practices such as returning structures conforming
// to specific interfaces rather than interfaces. In this case, we have to
// return interfaces for gomobile to do its job properly. If you are a Go
// user, you should probably use the more specific APIs that are part of this
// repository rather than using this API.
//
// See https://github.com/measurement-kit/api.
package mkapi

import (
	"github.com/measurement-kit/engine/collector"
)

// Engine is a Measurement Kit engine.
type Engine interface {
	// NewCollectorSubmitTask creates a new CollectorSubmitTask initialized with
	// the specified softwareName, softwareVersion and measurement.
	NewCollectorSubmitTask(
		softwareName, softwareVersion, measurement string) CollectorSubmitTask
}

// GoEngine is the Measurement Kit engine implemented in Go.
type GoEngine struct {
}

// NewCollectorSubmitTask creates a new concrete implementation
// of the CollectorSubmitTask interface.
func (engine *GoEngine) NewCollectorSubmitTask(
	softwareName, softwareVersion, measurement string) CollectorSubmitTask {
	return &goCollectorSubmitTask{
		t: collector.NewSubmitTask(
			softwareName, softwareVersion, measurement,
		),
	}
}
