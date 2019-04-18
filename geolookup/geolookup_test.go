package geolookup

import (
	"errors"
	"io"
	"testing"
)

// TestOpenNoGZ checks whether we correctly deal with a input
// file that is not compressed using gzip.
func TestOpenNoGZ(t *testing.T) {
	_, err := open("../version/version.go")
	if err == nil {
		t.Fatal("We expected an error here")
	}
}

// TestOpenIoutilReadAllError checks whether we correctly deal with an
// error when we're reading the input file.
func TestOpenIoutilReadAllError(t *testing.T) {
	savedFunc := ioutilReadAll
	mockedError := errors.New("mocked error")
	ioutilReadAll = func(r io.Reader) ([]byte, error) {
		return nil, mockedError
	}
	_, err := open("../asn.mmdb.gz")
	if err != mockedError {
		t.Fatal("Not the error we were expecting")
	}
	ioutilReadAll = savedFunc
}

// TestOpenGzCloseError checks whether we correctly deal with an
// error when we're closing the input gzip stream.
func TestOpenGzCloseError(t *testing.T) {
	savedFunc := gzclose
	mockedError := errors.New("mocked error")
	gzclose = func(r io.Closer) error {
		return mockedError
	}
	_, err := open("../asn.mmdb.gz")
	if err != mockedError {
		t.Fatal("Not the error we were expecting")
	}
	gzclose = savedFunc
}

// TestGetCCIntegration tests the common CC-lookup case.
func TestGetCCIntegration(t *testing.T) {
	probeCC, err := GetCC("../country.mmdb.gz", "8.8.8.8")
	t.Logf("CC: %s", probeCC)
	if err != nil {
		t.Fatal(err)
	}
}

// TestGetCCNonExistentFile checks whether GetCC will deal
// correctly with the case where the file is missing.
func TestGetCCNonExistentFile(t *testing.T) {
	_, err := GetCC("../nonexistent.mmdb.gz", "8.8.8.8")
	if err == nil {
		t.Fatal("We expected an error here")
	}
}

// TestGetCCBadIP checks whether GetCC will deal correctly with the
// case where the input is not actually a real IP.
func TestGetCCNonExistentIP(t *testing.T) {
	_, err := GetCC("../country.mmdb.gz", "127.0.0")
	if err == nil {
		t.Fatal("We expected an error here")
	}
}

// TestGetASNIntegration tests the common ASN-lookup case.
func TestGetASNIntegration(t *testing.T) {
	probeASN, probeOrg, err := GetASN("../asn.mmdb.gz", "8.8.8.8")
	t.Logf("ASN: %s; Org: %s", probeASN, probeOrg)
	if err != nil {
		t.Fatal(err)
	}
}

// TestGetASNNonExistentFile checks whether GetASN will deal
// correctly with the case where the file is missing.
func TestGetASNNonExistentFile(t *testing.T) {
	_, _, err := GetASN("../nonexistent.mmdb.gz", "8.8.8.8")
	if err == nil {
		t.Fatal("We expected an error here")
	}
}

// TestGetASNBadIP checks whether GetASN will deal correctly with the
// case where the input is not actually a real IP.
func TestGetASNNonExistentIP(t *testing.T) {
	_, _, err := GetASN("../asn.mmdb.gz", "127.0.0")
	if err == nil {
		t.Fatal("We expected an error here")
	}
}
