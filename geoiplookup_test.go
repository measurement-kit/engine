package engine

import (
	"fmt"
	"testing"
)

func TestGeoIPLookupIntegration(t *testing.T) {
	settings := &GeoIPLookupSettings{}
	settings.ASNDatabasePath = "asn.mmdb.gz"
	settings.CountryDatabasePath = "country.mmdb.gz"
	settings.Timeout = 14
	results := GeoIPLookup(settings)
	fmt.Println(results.ProbeIP)
	fmt.Println(results.ProbeASN)
	fmt.Println(results.ProbeCC)
	fmt.Println(results.ProbeOrg)
	fmt.Println(results.Logs)
	if !results.Good {
		t.Fatal("geoiplookup failed")
	}
}
