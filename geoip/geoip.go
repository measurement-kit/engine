// Package geoip implements mkall's GeoIP API.
package geoip

import (
	"context"
	"fmt"

	"github.com/measurement-kit/engine/internal"
	"github.com/measurement-kit/engine/internal/assets"
	"github.com/measurement-kit/engine/internal/geolookup"
	"github.com/measurement-kit/engine/internal/iplookup"
)

// LookupResults contains the results of a GeoIP lookup.
type LookupResults struct {
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

// LookupSettings contains the GeoIP lookup settings.
type LookupSettings struct {
	// Timeout is the number of seconds after which we abort.
	Timeout int64

	// WorkDirPath is the path to the working directory.
	WorkDirPath string
}

// LookupInto is like Lookup but the results are passed as a pointer.
func LookupInto(settings *LookupSettings, out *LookupResults) {
	if settings.WorkDirPath == "" {
		out.Logs = fmt.Sprintf("WorkDirPath is not set\n")
		return
	}
	duration, err := internal.MakeTimeout(settings.Timeout)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot make duration: %s\n", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	err = assets.Download(ctx, settings.WorkDirPath)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot download assets: %s\n", err.Error())
		return
	}
	probeIP, err := iplookup.Perform(ctx)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot discover probe IP: %s\n", err.Error())
		return
	}
	out.ProbeIP = probeIP
	probeASN, probeOrg, err := geolookup.GetASN(
		assets.ASNDatabasePath(settings.WorkDirPath), out.ProbeIP,
	)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot discover probe ASN: %s\n", err.Error())
		return
	}
	out.ProbeASN, out.ProbeOrg = probeASN, probeOrg
	probeCC, err := geolookup.GetCC(
		assets.CountryDatabasePath(settings.WorkDirPath), out.ProbeIP,
	)
	if err != nil {
		out.Logs = fmt.Sprintf("cannot discover probe CC: %s\n", err.Error())
		return
	}
	out.ProbeCC = probeCC
	out.Good = true
}

// Lookup performs a GeoIP lookup.
func Lookup(settings *LookupSettings) *LookupResults {
	var out LookupResults
	LookupInto(settings, &out)
	return &out
}
