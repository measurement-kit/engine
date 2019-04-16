// Package assets contains code to manage assets.
package assets

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/measurement-kit/engine/httpx"
)

// asset is an asset that we can download. Assets are published at the
// github.com/measurement-kit/generic-assets repository. They are always
// gzipped, to save bandwidth.
type asset struct {
	// URLPath is the relative URL to download the asset from. The URL
	// is relative to github.com/measurement-kit/generic-assets/releases.
	URLPath string

	// SHA256 is the checksum of the compressed resource.
	SHA256 string
}

// allAssets describes all assets that we could download.
var allAssets = []asset{
	asset{
		URLPath: "download/20190327/asn.mmdb.gz",
		SHA256:  "6fcae12937b383e1f067e14d1eb728a75a360279df8240517ac70ef6d401c2be",
	},
	asset{
		URLPath: "download/20190327/country.mmdb.gz",
		SHA256:  "d0a499d15506c54111217f30af9dfd11476ded076c55a3e28a73715c890b5d66",
	},
}

// baseURL is the base URL used to download assets from.
const baseURL = `https://github.com/measurement-kit/generic-assets/releases/`

// errSHA256Mismatch is the error returned on SHA256 mismatch.
var errSHA256Mismatch = errors.New("SHA256 does not match expected SHA256")

// save saves the specified, compressed asset as filename, which is the absolute
// file path of the asset in the destination directory.
func save(ctx context.Context, filename string, asset asset) error {
	data, err := httpx.GET(ctx, baseURL+asset.URLPath)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%x", sha256.Sum256(data)) != asset.SHA256 {
		return errSHA256Mismatch
	}
	return ioutil.WriteFile(filename, data, 0600)
}

// saveIdempotent saves the specified compressed asset in destdir only
// if we have not downladed the same file already.
func saveIdempotent(ctx context.Context, destdir string, asset asset) error {
	filename := filepath.Join(destdir, filepath.Base(asset.URLPath))
	filep, err := os.Open(filename)
	if err != nil {
		return save(ctx, filename, asset)
	}
	defer filep.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, filep); err != nil {
		return save(ctx, filename, asset)
	}
	if fmt.Sprintf("%x", hash.Sum(nil)) != asset.SHA256 {
		return save(ctx, filename, asset)
	}
	return nil
}

// Download downloads assets in destdir.
func Download(ctx context.Context, destdir string) error {
	if err := os.MkdirAll(destdir, 0700); err != nil {
		return err
	}
	for _, asset := range allAssets {
		if err := saveIdempotent(ctx, destdir, asset); err != nil {
			return err
		}
	}
	return nil
}
