package geoip

import (
	"fmt"
	"testing"
)

func TestLookupIntegration(t *testing.T) {
	settings := &LookupSettings{}
	settings.WorkDirPath = "/tmp"
	settings.Timeout = 60
	results := Lookup(settings)
	fmt.Println(results.ProbeIP)
	fmt.Println(results.ProbeASN)
	fmt.Println(results.ProbeCC)
	fmt.Println(results.ProbeOrg)
	fmt.Println(results.Logs)
	if !results.Good {
		t.Fatal("geoiplookup failed")
	}
}
