package task

import (
	"fmt"
	"testing"
)

// TestPsiphonTunnelIntegration runs a psiphontunnel task.
func TestPsiphonTunnelIntegration(t *testing.T) {
	handle := Start(`{
		"name": "PsiphonTunnel",
		"options": {
			"config_file_path": "/tmp/psiphon.json",
			"software_name": "mke-test",
			"software_version": "0.0.1",
			"work_dir_path": "/tmp"
		}
	}`)
	for !handle.IsDone() {
		ev := handle.WaitForNextEvent()
		fmt.Printf("%s\n\n", ev)
	}
}
