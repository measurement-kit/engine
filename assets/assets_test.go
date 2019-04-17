package assets

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

// TestSaveIntegration invokes save with common arguments
func TestSaveIntegration(t *testing.T) {
	if testing.Short() {
		return
	}
	err := save(context.Background(), "../temp.mmdb.gz", allAssets[0])
	if err != nil {
		t.Fatal(err)
	}
}

// TestSaveHTTPXGetFailure checks whether save is able
// to deal with a httpx.GET failure.
func TestSaveHTTPXGetFailure(t *testing.T) {
	mockedError := errors.New("mocked error")
	savedHTTPXGET := httpxGET
	httpxGET = func(ctx context.Context, URL string) ([]byte, error) {
		return nil, mockedError
	}
	err := save(context.Background(), "../temp.mmdb.gz", allAssets[0])
	if err != mockedError {
		t.Fatal("Not the error that we were expecting")
	}
	httpxGET = savedHTTPXGET
}

// TestSaveSHA256Mismatch checks whether save is able
// to deal with a SHA256 mismatch.
func TestSaveSHA256Mismatch(t *testing.T) {
	savedHTTPXGET := httpxGET
	httpxGET = func(ctx context.Context, URL string) ([]byte, error) {
		// This body should be such that the SHA256 never matches
		return []byte("deadbeeef"), nil
	}
	err := save(context.Background(), "../temp.mmdb.gz", allAssets[0])
	if err == nil {
		t.Fatal("We were expecting an error")
	}
	httpxGET = savedHTTPXGET
}

// TestCacheOpenOSOpenError checks whether cacheOpen is able
// to deal with a os.Open failure
func TestCacheOpenOSOpenError(t *testing.T) {
	mockedError := errors.New("mocked error")
	savedOSOpen := osOpen
	osOpen = func(name string) (*os.File, error) {
		return nil, mockedError
	}
	err := cacheOpen(context.Background(), "..", allAssets[0])
	if err != mockedError {
		t.Fatal("Not the error that we were expecting")
	}
	osOpen = savedOSOpen
}

// TestCacheOpenIOCopyError checks whether cacheOpen is able
// to deal with a io.Copy failure
func TestCacheOpenIOCopyError(t *testing.T) {
	savedOSOpen := osOpen
	osOpen = func(name string) (*os.File, error) {
		return ioutil.TempFile("", "measurement-kit-engine")
	}
	mockedError := errors.New("mocked error")
	savedIOCopy := ioCopy
	ioCopy = func(dst io.Writer, src io.Reader) (int64, error) {
		return 0, mockedError
	}
	err := cacheOpen(context.Background(), "..", allAssets[0])
	if err != mockedError {
		t.Fatal("Not the error that we were expecting")
	}
	osOpen = savedOSOpen
	ioCopy = savedIOCopy
}

// TestCacheOpenSHA256Mismatch checks whether cacheOpen is able
// to deal with a SHA256 mismatch failure.
func TestCacheSHA256Mismatch(t *testing.T) {
	savedOSOpen := osOpen
	osOpen = func(name string) (*os.File, error) {
		return ioutil.TempFile("", "measurement-kit-engine")
	}
	err := cacheOpen(context.Background(), "..", allAssets[0])
	if err == nil || err.Error() != "SHA256 mismatch" {
		t.Fatal("No error or not the error that we were expecting")
	}
	osOpen = savedOSOpen
}

// TestSaveIdempotentOpenCacheError checks whether saveIdempotent
// is able to deal with a cacheOpen error
func TestSaveIdempotentOpenCacheError(t *testing.T) {
	savedCacheOpen := mockableCacheOpen
	mockableCacheOpen = func(ctx context.Context, filename string, asset asset) error {
		return errors.New("cannot access file from cache")
	}
	err := saveIdempotent(context.Background(), "..", allAssets[0])
	if err != nil {
		t.Fatal(err)
	}
	mockableCacheOpen = savedCacheOpen
}

// TestDownloadIntegration invokes Download with common arguments
func TestDownloadIntegration(t *testing.T) {
	if err := Download(context.Background(), ".."); err != nil {
		t.Fatal(err)
	}
}

// TestDownloadMkdirAllError checks whether Download
// is able to deal with a os.MkdirAll error
func TestDownloadMkdirAllError(t *testing.T) {
	mockedError := errors.New("mocked error")
	savedMkdirAll := osMkdirAll
	osMkdirAll = func(path string, perm os.FileMode) error {
		return mockedError
	}
	err := Download(context.Background(), "..")
	if err != mockedError {
		t.Fatal("Not the error we expected to see")
	}
	osMkdirAll = savedMkdirAll
}

// TestDownloadSaveIdempotentError checks whether Download
// is able to deal with a saveIdempotent error
func TestDownloadSaveIdempotent(t *testing.T) {
	mockedError := errors.New("mocked error")
	savedFunc := mockableSaveIdempotent
	mockableSaveIdempotent = func(ctx context.Context, destdir string, asset asset) error {
		return mockedError
	}
	err := Download(context.Background(), "..")
	if err != mockedError {
		t.Fatal("Not the error we expected to see")
	}
	mockableSaveIdempotent = savedFunc
}
