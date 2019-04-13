package psiphontunnel

import (
	"context"
	"fmt"
	"testing"
)

func TestRunIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "/tmp/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	result := Run(context.Background(), config)
	fmt.Printf("%+v\n", result)
	if result.Failure != "" {
		t.Fatal("Failure is not empty")
	}
	if result.BootstrapTime <= 0.0 {
		t.Fatal("BootstrapTime is not positive")
	}
}
