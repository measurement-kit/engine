package mobile

import (
	"fmt"
	"testing"
)

func TestGeoIPLookupIntegration(t *testing.T) {
	settings := &MKEGeoIPLookupSettings{}
	settings.ASNDatabasePath = "../asn.mmdb"
	settings.CountryDatabasePath = "../country.mmdb"
	settings.Timeout = 14
	results := settings.Perform()
	fmt.Println(results.ProbeIP)
	fmt.Println(results.ProbeASN)
	fmt.Println(results.ProbeCC)
	fmt.Println(results.ProbeOrg)
	fmt.Println(results.Logs)
	if !results.Good {
		t.Fatal("geoiplookup failed")
	}
}
