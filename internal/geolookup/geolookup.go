// Package geolookup allows to geolookup a OONI probe.
//
// Specifically, the objective of the geolookup is to discover:
//
// 1. the autonomous system number (ASN) associated to such IP (aka probeASN);
//
// 2. the code of the country in which the IP is (aka probeCC);
//
// 3. the name associated to the probe' ASN (aka probeOrg).
//
// To this end, we use MaxMind databases using the MMDB data format.
package geolookup

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	// oschwald is a maxmind developer, therefore I expect this package
	// to work reasonably well even though it's not official.
	"github.com/oschwald/maxminddb-golang"
)

// ioutilReadAll is a mockable ioutil.ReadAll
var ioutilReadAll = ioutil.ReadAll

// gzclose allows to verify we deal with gzip.Close errors
var gzclose = func(c io.Closer) error {
	return c.Close()
}

// open opens a compressed database.
func open(dbpath string) (*maxminddb.Reader, error) {
	filep, err := os.Open(dbpath)
	if err != nil {
		return nil, err
	}
	defer filep.Close()
	gzfilep, err := gzip.NewReader(filep)
	if err != nil {
		return nil, err
	}
	// Implementation note: don't discard gzip.Close return value since
	// it may actually indicate that the file is corrupted.
	data, err := ioutilReadAll(gzfilep)
	if err != nil {
		gzclose(gzfilep)
		return nil, err
	}
	err = gzclose(gzfilep)
	if err != nil {
		return nil, err
	}
	return maxminddb.FromBytes(data)
}

// GetCC returns the probeCC. In case of failure, probeCC is "ZZ".
func GetCC(dbpath, IP string) (string, error) {
	db, err := open(dbpath)
	if err != nil {
		return "ZZ", err
	}
	defer db.Close()
	dataIP := net.ParseIP(IP)
	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}
	err = db.Lookup(dataIP, &record)
	if err != nil {
		return "ZZ", err
	}
	return record.Country.ISOCode, nil
}

// GetASN lookups the probeASN and the probeOrg. In case of failure, probeASN
// is "AS0", otherwise it will be "AS<number>". In case of failure, probeOrg is
// empty, otherwise it's the commercial name of the ASN.
func GetASN(dbpath, IP string) (string, string, error) {
	db, err := open(dbpath)
	if err != nil {
		return "AS0", "", err
	}
	defer db.Close()
	dataIP := net.ParseIP(IP)
	var record struct {
		ASN int    `maxminddb:"autonomous_system_number"`
		Org string `maxminddb:"autonomous_system_organization"`
	}
	err = db.Lookup(dataIP, &record)
	if err != nil {
		return "AS0", "", err
	}
	return fmt.Sprintf("AS%d", record.ASN), record.Org, nil
}
