package mobile

import (
	"context"
	"fmt"

	"github.com/measurement-kit/engine/geolookup"
	"github.com/measurement-kit/engine/iplookup"
	"github.com/measurement-kit/engine/mobile/internal"
)

// MKEGeoIPLookupResults contains the results of a GeoIP lookup.
type MKEGeoIPLookupResults struct {
	// Good indicates whether we succeded.
	Good bool

	// ProbeIP is the probe IP.
	ProbeIP string

	// ProbeASN is the probe ASN.
	ProbeASN string

	// ProbeCC is the probe CC.
	ProbeCC string

	// ProbeOrg is the organization owning the ASN.
	ProbeOrg string

	// Logs contains logs useful to debug errors.
	Logs string
}

// MKEGeoIPLookupSettings contains the GeoIP lookup settings.
type MKEGeoIPLookupSettings struct {
	// Timeout is the number of seconds after which we abort.
	Timeout int64

	// ASNDatabasePath is the path to the ASN database.
	ASNDatabasePath string

	// CountryDatabasePath is the path to the country database.
	CountryDatabasePath string
}

// Perform performs a GeoIP lookup.
func (x *MKEGeoIPLookupSettings) Perform() *MKEGeoIPLookupResults {
	var out MKEGeoIPLookupResults
	duration, err := internal.MakeDuration(x.Timeout)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot make duration: %s\n", err.Error())
		return &out
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	probeIP, err := iplookup.Perform(ctx)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot discover probe IP: %s\n", err.Error())
		return &out
	}
	out.ProbeIP = probeIP
	probeASN, probeOrg, err := geolookup.GetASN(x.ASNDatabasePath, out.ProbeIP)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot discover probe ASN: %s\n", err.Error())
		return &out
	}
	out.ProbeASN, out.ProbeOrg = probeASN, probeOrg
	probeCC, err := geolookup.GetCC(x.CountryDatabasePath, out.ProbeIP)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot discover probe CC: %s\n", err.Error())
		return &out
	}
	out.ProbeCC = probeCC
	out.Good = true
	return &out
}
