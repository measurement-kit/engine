package task_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/measurement-kit/engine/task"
)

// TestNdt7Integration runs a ndt7 nettest.
func TestNdt7Integration(t *testing.T) {
	ctx := context.Background()
	for ev := range task.StartNdt7(ctx, task.Config{}) {
		data, err := json.Marshal(ev)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
	}
}

// TestPsiphonTunnelIntegration runs a psiphontunnel nettest
func TestPsiphonTunnelIntegration(t *testing.T) {
	ctx := context.Background()
	config := task.Config{
		ConfigFilePath: "../testdata/psiphon_config.json",
		WorkDirPath:    "/tmp",
	}
	for ev := range task.StartPsiphonTunnel(ctx, config) {
		data, err := json.Marshal(ev)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(data))
	}
}
