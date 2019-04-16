package mlabns

import (
	"context"
	"fmt"
	"testing"
)

func TestGeoOptionsIntegration(t *testing.T) {
	config := Config{Tool: "ndt_ssl"}
	servers, err := GeoOptions(context.Background(), config)
	for _, server := range(servers) {
		fmt.Println(server.FQDN)
	}
	if err != nil {
		t.Fatal(err)
	}
}
