// Package httpx contains HTTP extensions
package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"

	"github.com/measurement-kit/engine/internal/version"
)

// Request is an HTTP request.
type Request struct {
	// Ctx is the mandatory request context.
	Ctx context.Context

	// Method is the mandatory request method.
	Method string

	// URL is the mandatory request URL.
	URL string

	// ContentType is the optional content type.
	ContentType string

	// UserAgent is the optional user agent.
	UserAgent string

	// Body is the optional request body.
	Body []byte

	// NoFailOnError controls whether an HTTP failure causes
	// the Perform function to fail or not.
	NoFailOnError bool

	// SOCKS5ProxyPort is the optional SOCKS5 proxy port to use. The
	// default value (zero) means no proxy is used.
	SOCKS5ProxyPort int
}

// Response is an HTTP response
type Response struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// ContentType is the optional content type.
	ContentType string

	// Body is the optional response body.
	Body []byte
}

// ioutilReadAll allows to mock ioutil.ReadAll when testing the code.
var ioutilReadAll = ioutil.ReadAll

// proxySOCKS5 allows to mock proxy.SOCKS5 when testing the code.
var proxySOCKS5 = proxy.SOCKS5

func (r Request) perform() (*Response, error) {
	request, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}
	if r.ContentType != "" {
		request.Header.Set("Content-Type", r.ContentType)
	}
	if r.UserAgent != "" {
		request.Header.Set("User-Agent", r.UserAgent)
	}
	request = request.WithContext(r.Ctx)
	var client *http.Client
	if r.SOCKS5ProxyPort != 0 {
		// TODO(bassosimone): for correctness here we MUST make sure that
		// this proxy implementation does not leak the DNS.
		endpoint := fmt.Sprintf("127.0.0.1:%d", r.SOCKS5ProxyPort)
		dialer, err := proxySOCKS5("tcp", endpoint, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		client = &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	} else {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 && !r.NoFailOnError {
		return nil, fmt.Errorf(
			"Request failed with status %d", response.StatusCode,
		)
	}
	data, err := ioutilReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode:  response.StatusCode,
		ContentType: response.Header.Get("Content-Type"),
		Body:        data,
	}, nil
}

// Perform performs an HTTP request and returns the response.
func (r Request) Perform() (*Response, error) {
	response, err := r.perform()
	if err != nil {
		return nil, fmt.Errorf(
			"%s %s failed: %s", r.Method, r.URL, err.Error(),
		)
	}
	return response, nil
}

// userAgent creates the user agent string
func userAgent() string {
	return fmt.Sprintf("MKEngine/%s", version.Version)
}

// GET performs a GET request and returns the body.
func GET(ctx context.Context, URL string) ([]byte, error) {
	response, err := Request{
		Ctx:       ctx,
		Method:    "GET",
		URL:       URL,
		UserAgent: userAgent(),
	}.Perform()
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

// GETWithBaseURL is like GET but with baseURL and path.
func GETWithBaseURL(ctx context.Context, baseURL, path string) ([]byte, error) {
	URL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	URL.Path = path
	return GET(ctx, URL.String())
}

// POST performs a POST request and returns the body.
func POST(ctx context.Context, URL, contentType string, body []byte) ([]byte, error) {
	response, err := Request{
		Ctx:         ctx,
		Method:      "POST",
		URL:         URL,
		Body:        body,
		ContentType: contentType,
		UserAgent:   userAgent(),
	}.Perform()
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

// POSTWithBaseURL performs a POST with a baseURL.
func POSTWithBaseURL(ctx context.Context, baseURL, path, contentType string, body []byte) ([]byte, error) {
	URL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	URL.Path = path
	return POST(ctx, URL.String(), contentType, body)
}
