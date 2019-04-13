// Package assets contains code to manage assets.
package assets

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type asset struct {
	URLPath  string
	SHA256   string
	Filename string
}

var allAssets = []asset{
	asset{
		URLPath:  "download/20190327/asn.mmdb.gz",
		SHA256:   "6fcae12937b383e1f067e14d1eb728a75a360279df8240517ac70ef6d401c2be",
		Filename: "asn.mmdb",
	},
	asset{
		URLPath:  "download/20190327/country.mmdb.gz",
		SHA256:   "d0a499d15506c54111217f30af9dfd11476ded076c55a3e28a73715c890b5d66",
		Filename: "country.mmdb",
	},
}

func get(ctx context.Context, URL string) ([]byte, error) {
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("The request failed")
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func save(ctx context.Context, destdir string, asset asset) error {
	const baseURL = `https://github.com/measurement-kit/generic-assets/releases/`
	data, err := get(ctx, baseURL+asset.URLPath)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%x", sha256.Sum256(data)) != asset.SHA256 {
		return errors.New("SHA256 does not match expected SHA256")
	}
	gunzipper, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gunzipper.Close()
	data, err = ioutil.ReadAll(gunzipper)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(destdir, asset.Filename), data, 0600)
}

// Download downloads assets in destdir.
func Download(ctx context.Context, destdir string) error {
	if err := os.MkdirAll(destdir, 0700); err != nil {
		return err
	}
	for _, asset := range allAssets {
		if err := save(ctx, destdir, asset); err != nil {
			return err
		}
	}
	return nil
}
