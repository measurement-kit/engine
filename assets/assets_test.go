package assets

import (
	"context"
	"testing"
)

func TestDownloadIntegration(t *testing.T) {
	if err := Download(context.Background(), ".."); err != nil {
		t.Fatal(err)
	}
}
