// Package runner implements the psiphontunnel runner.
package runner

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/proxy"

	"github.com/Psiphon-Labs/psiphon-tunnel-core/ClientLibrary/clientlib"
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

var osRemoveAll = os.RemoveAll
var osMkdirAll = os.MkdirAll
var ioutilReadFile = ioutil.ReadFile
var clientlibStartTunnel = clientlib.StartTunnel
var urlParse = url.Parse
var clntGet = func(clnt *http.Client, URL string) (*http.Response, error) {
	return clnt.Get(URL)
}
var proxySOCKS5 = proxy.SOCKS5

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

func usetunnel(t *clientlib.PsiphonTunnel) error {
	// TODO(bassosimone): for correctness here we MUST make sure that
	// this proxy implementation does not leak the DNS.
	endpoint := fmt.Sprintf("127.0.0.1:%d", t.SOCKSProxyPort)
	dialer, err := proxySOCKS5("tcp", endpoint, nil, proxy.Direct)
	if err != nil {
		return err
	}
	clnt := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	const URL = "https://www.google.com/humans.txt"
	response, err := clntGet(clnt, URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("HTTP status code is not 200")
	}
	return nil
}

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
	err = usetunnel(tunnel)
	if err != nil {
		result.Failure = err.Error()
		return result
	}
	return result
}
