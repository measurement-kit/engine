// +build !small

package task

import (
	"context"

	"github.com/measurement-kit/engine/internal/nettest"
	"github.com/measurement-kit/engine/internal/nettest/psiphontunnel"
)

// newPsiphonTunnel creates a new PsiphonTunnel nettest.
func newPsiphonTunnel(ctx context.Context, settings taskSettings) *nettest.Nettest {
	return psiphontunnel.NewNettest(ctx, psiphontunnel.Config{
		ConfigFilePath: settings.Options.ConfigFilePath,
		WorkDirPath:    settings.Options.WorkDirPath,
	})
}
