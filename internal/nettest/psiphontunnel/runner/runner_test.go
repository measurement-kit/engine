package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Psiphon-Labs/psiphon-tunnel-core/ClientLibrary/clientlib"
)

// TestProcessconfigEmptyWorkDir checks whether processconfig deals
// with an empty WorkDir in a sane way.
func TestProcessconfigEmptyWorkDir(t *testing.T) {
	_, _, err := processconfig(Config{})
	if err == nil || err.Error() != "WorkDirPath is empty" {
		t.Fatal("No error or unexpected error")
	}
}

// TestProcessconfigOSRemoveAllError checks whether processconfig deals
// with an error when running os.RemoveAll
func TestProcessconfigOSRemoveAllError(t *testing.T) {
	savedFunc := osRemoveAll
	mockedError := errors.New("mocked error")
	osRemoveAll = func(string) error {
		return mockedError
	}
	_, _, err := processconfig(Config{
		WorkDirPath: "/tmp",
	})
	if err != mockedError {
		t.Fatal("Not the error we expected")
	}
	osRemoveAll = savedFunc
}

// TestProcessconfigOSMkdirAllError checks whether processconfig deals
// with an error when running os.MkdirAll
func TestProcessconfigOSMkdirAllError(t *testing.T) {
	savedFunc := osMkdirAll
	mockedError := errors.New("mocked error")
	osMkdirAll = func(string, os.FileMode) error {
		return mockedError
	}
	_, _, err := processconfig(Config{
		WorkDirPath: "/tmp",
	})
	if err != mockedError {
		t.Fatal("Not the error we expected")
	}
	osMkdirAll = savedFunc
}

// TestProcessconfigIoutilReadFileError checks whether processconfig deals
// with an error when running ioutil.ReadFile
func TestProcessconfigIoutilReadFileError(t *testing.T) {
	savedFunc := ioutilReadFile
	mockedError := errors.New("mocked error")
	ioutilReadFile = func(string) ([]byte, error) {
		return nil, mockedError
	}
	_, _, err := processconfig(Config{
		WorkDirPath:    "/tmp",
		ConfigFilePath: "../../../../testdata/psiphon_config.json",
	})
	if err != mockedError {
		t.Fatal("Not the error we expected")
	}
	ioutilReadFile = savedFunc
}

// TestRunIntegration just runs Run in the common case.
func TestRunIntegration(t *testing.T) {
	config := Config{
		ConfigFilePath: "../../../../testdata/psiphon_config.json",
		WorkDirPath:    "/tmp/",
	}
	result := Run(context.Background(), config)
	fmt.Printf("%+v\n", result)
	if result.Failure != "" {
		t.Fatal("Failure is not empty")
	}
	if result.BootstrapTime <= 0.0 {
		t.Fatal("BootstrapTime is not positive")
	}
}

// TestRunProcessconfigFailure checks whether Run deals with a
// failure in processing the configuration.
func TestRunProcessconfigFailure(t *testing.T) {
	config := Config{
		ConfigFilePath: "/nonexistent/psiphon.json",
		WorkDirPath:    "/tmp/",
	}
	result := Run(context.Background(), config)
	if result.Failure == "" {
		t.Fatal("We expected a failure, found no error")
	}
}

// TestRunClientlibStartTunnelFailure checks whether Run deals
// with a failure in starting the psiphon tunnel.
func TestRunClientlibStartTunnelFailure(t *testing.T) {
	savedFunc := clientlibStartTunnel
	mockedError := errors.New("mocked error")
	clientlibStartTunnel = func(
		ctx context.Context,
		configJSON []byte,
		embeddedServerEntryList string,
		params clientlib.Parameters,
		paramsDelta clientlib.ClientParametersDelta,
		noticeReceiver func(clientlib.NoticeEvent)) (tunnel *clientlib.PsiphonTunnel, err error) {
		return nil, mockedError
	}
	config := Config{
		ConfigFilePath: "../../../../testdata/psiphon_config.json",
		WorkDirPath:    "/tmp/",
	}
	result := Run(context.Background(), config)
	if result.Failure != mockedError.Error() {
		t.Fatal("Not the error that we were expecting")
	}
	clientlibStartTunnel = savedFunc
}

// TestRunUsetunnelFailure checks whether Run deals
// with a failure in using the psiphon tunnel.
func TestRunUsetunnelFailure(t *testing.T) {
	savedFunc := mockableUsetunnel
	mockedError := errors.New("mocked error")
	mockableUsetunnel = func(ctx context.Context, t *clientlib.PsiphonTunnel) error {
		return mockedError
	}
	config := Config{
		ConfigFilePath: "../../../../testdata/psiphon_config.json",
		WorkDirPath:    "/tmp/",
	}
	result := Run(context.Background(), config)
	if result.Failure != mockedError.Error() {
		t.Fatal("Not the error that we were expecting")
	}
	mockableUsetunnel = savedFunc
}
