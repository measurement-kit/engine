package engine

import (
	"fmt"
	"testing"
)

// TestPsiphonTunnelIntegration runs a psiphontunnel nettest.
func TestPsiphonTunnelIntegration(t *testing.T) {
	task := StartTask(`{
		"name": "PsiphonTunnel",
		"options": {
			"config_file_path": "/tmp/psiphon.json",
			"software_name": "mke-test",
			"software_version": "0.0.1",
			"work_dir_path": "/tmp"
		}
	}`)
	for !task.IsDone() {
		ev := task.WaitForNextEvent()
		fmt.Printf("%s\n\n", ev)
	}
}
