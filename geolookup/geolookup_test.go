package geolookup

import (
	"testing"
)

func TestGetCCIntegration(t *testing.T) {
	probeCC, err := GetCC("../country.mmdb.gz", "8.8.8.8")
	t.Logf("CC: %s", probeCC)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetASNIntegration(t *testing.T) {
	probeASN, probeOrg, err := GetASN("../asn.mmdb.gz", "8.8.8.8")
	t.Logf("ASN: %s; Org: %s", probeASN, probeOrg)
	if err != nil {
		t.Fatal(err)
	}
}
