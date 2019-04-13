package bouncer

import (
	"context"
	"testing"

	"github.com/measurement-kit/engine/model"
)

func doit(t *testing.T, f func(context.Context, Config) ([]model.Service, error)) {
	entries, err := f(context.Background(), Config{
		BaseURL: "https://events.proteus.test.ooni.io",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		t.Logf("%+v", entry)
	}
}

func TestGetCollectorsIntegration(t *testing.T) {
	doit(t, GetCollectors)
}

func TestGetTestHelpersIntegration(t *testing.T) {
	doit(t, GetTestHelpers)
}
