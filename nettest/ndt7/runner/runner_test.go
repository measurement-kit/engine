package runner

import (
	"context"
	"fmt"
	"testing"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	FQDNs, err := GetServers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, FQDN := range FQDNs {
		fmt.Printf("try: %s\n", FQDN)
		ch, err := StartDownload(ctx, FQDN)
		if err != nil {
			t.Error(err)
			continue
		}
		for ev := range ch {
			fmt.Printf("%+v\n", ev)
		}
		break
	}
	for _, FQDN := range FQDNs {
		fmt.Printf("try: %s\n", FQDN)
		ch, err := StartUpload(ctx, FQDN)
		if err != nil {
			t.Error(err)
			continue
		}
		for ev := range ch {
			fmt.Printf("%+v\n", ev)
		}
		break
	}
}
