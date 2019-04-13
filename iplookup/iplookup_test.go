package iplookup

import (
	"context"
	"testing"
)

func TestIntegration(t *testing.T) {
	IP, err := Perform(context.Background())
	t.Logf("IP: %s", IP)
	if err != nil {
		t.Fatal(err)
	}
}
