// Package runner implements the psiphontunnel runner.
package runner

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Psiphon-Labs/psiphon-tunnel-core/ClientLibrary/clientlib"
	"github.com/measurement-kit/engine/httpx"
)

// Config contains the nettest configuration.
type Config struct {
	// ConfigFilePath is the path where Psiphon config file is located.
	ConfigFilePath string `json:"config_file_path"`

	// WorkDirPath is the directory where Psiphon should store
	// its configuration database.
	WorkDirPath string `json:"work_dir_path"`
}

// Result contains the nettest result.
//
// This is what will end up into the Measurement.TestKeys field
// when you run this nettest.
type Result struct {
	// Failure contains the failure that occurred.
	Failure string `json:"failure"`

	// BootstrapTime is the time it took to bootstrap Psiphon.
	BootstrapTime float64 `json:"bootstrap_time"`
}

// osRemoveAll is a mockable os.RemoveAll
var osRemoveAll = os.RemoveAll

// osMkdirAll is a mockable os.MkdirAll
var osMkdirAll = os.MkdirAll

// ioutilReadFile is a mockable ioutil.ReadFile
var ioutilReadFile = ioutil.ReadFile

func processconfig(config Config) ([]byte, clientlib.Parameters, error) {
	if config.WorkDirPath == "" {
		return nil, clientlib.Parameters{}, errors.New("WorkDirPath is empty")
	}
	const testdirname = "oonipsiphontunnelcore"
	workdir := filepath.Join(config.WorkDirPath, testdirname)
	err := osRemoveAll(workdir)
	if err != nil {
		return nil, clientlib.Parameters{}, err
	}
	err = osMkdirAll(workdir, 0700)
	if err != nil {
		return nil, clientlib.Parameters{}, err
	}
	params := clientlib.Parameters{
		DataRootDirectory: &workdir,
	}
	configJSON, err := ioutilReadFile(config.ConfigFilePath)
	if err != nil {
		return nil, clientlib.Parameters{}, err
	}
	return configJSON, params, nil
}

func usetunnel(ctx context.Context, t *clientlib.PsiphonTunnel) error {
	_, err := httpx.Request{
		Ctx:             ctx,
		Method:          "GET",
		URL:             "https://www.google.com/humans.txt",
		SOCKS5ProxyPort: t.SOCKSProxyPort,
	}.Perform()
	return err
}

// clientlibStartTunnel is a mockable clientlib.StartTunnel
var clientlibStartTunnel = clientlib.StartTunnel

// mockableUsetunnel is mockable usetunnel
var mockableUsetunnel = usetunnel

// Run runs the nettest and returns the result.
func Run(ctx context.Context, config Config) Result {
	var result Result
	configJSON, params, err := processconfig(config)
	if err != nil {
		result.Failure = err.Error()
		return result
	}
	t0 := time.Now()
	tunnel, err := clientlibStartTunnel(ctx, configJSON, "", params, nil, nil)
	if err != nil {
		result.Failure = err.Error()
		return result
	}
	result.BootstrapTime = time.Now().Sub(t0).Seconds()
	defer tunnel.Stop()
	err = mockableUsetunnel(ctx, tunnel)
	if err != nil {
		result.Failure = err.Error()
		return result
	}
	return result
}
