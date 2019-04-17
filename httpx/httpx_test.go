package httpx

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/proxy"

	"github.com/mccutchen/go-httpbin/httpbin"
)

// maxBodySize is the max body size we can request to httpbin
const maxBodySize = int64(10 * 1024 * 1024)

// withHTTPBin allows running tests with a custom HTTP server
func withHTTPBin(t *testing.T, fn func(baseURL string)) {
	httpbin := httpbin.New()
	httpbin.MaxBodySize = maxBodySize
	srv := httptest.NewServer(httpbin.Handler())
	defer srv.Close()
	fn(srv.URL)
}

// TestPerformSimple performs a simple httpx.Perform test
func TestPerformSimple(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		var r Request
		r.Ctx = context.Background()
		r.Method = "GET"
		r.URL = baseURL + "/status/200"
		response, err := r.Perform()
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != 200 {
			t.Fatal("Invalid response status code")
		}
	})
}

// TestPerformNewRequestError verifies that httpx.Perform
// handles an error occurring in http.NewRequest
func TestPerformNewRequestError(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		var r Request
		r.Ctx = context.Background()
		r.Method = "\t"
		r.URL = baseURL + "/status/200"
		_, err := r.Perform()
		if err == nil {
			t.Fatal("An error was expected")
		}
	})
}

// TestPerform400 verifies that httpx.Perform handles a 400 error
func TestPerform400(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		var r Request
		r.Ctx = context.Background()
		r.Method = "GET"
		r.URL = baseURL + "/status/400"
		_, err := r.Perform()
		if err == nil {
			t.Fatal("An error was expected")
		}
	})
}

// TestPerformReadBodyError verifies that httpx.Perform handles an
// error while reading the body
func TestPerformReadBodyError(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		mockedError := errors.New("mocked error")
		savedReadAll := ioutilReadAll
		ioutilReadAll = func(r io.Reader) ([]byte, error) {
			return nil, mockedError
		}
		var r Request
		r.Ctx = context.Background()
		r.Method = "GET"
		r.URL = baseURL + "/status/200"
		_, err := r.Perform()
		if err != mockedError {
			t.Fatal("Not the error we were expecting")
		}
		ioutilReadAll = savedReadAll
	})
}

// TestPerformProxySOCKS5Error verifies that httpx.Perform handles an
// error while settings up a proxy
func TestPerformProxySOCKS5Error(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		mockedError := errors.New("mocked error")
		savedProxySOCKS5 := proxySOCKS5
		proxySOCKS5 = func(network, address string, auth *proxy.Auth, dialer proxy.Dialer) (proxy.Dialer, error) {
			return nil, mockedError
		}
		var r Request
		r.Ctx = context.Background()
		r.Method = "GET"
		r.URL = baseURL + "/status/200"
		r.SOCKS5ProxyPort = 9999
		_, err := r.Perform()
		if err != mockedError {
			t.Fatal("Not the error we were expecting")
		}
		proxySOCKS5 = savedProxySOCKS5
	})
}

// TestGETSimple performs a simple httpx.GET test
func TestGETSimple(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		ctx := context.Background()
		_, err := GET(ctx, baseURL+"/status/200")
		if err != nil {
			t.Fatal(err)
		}
	})
}

// TestGETError checks whether httpx.GET correctly handles errors
func TestGETError(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		ctx := context.Background()
		const emptyURL = "" // The code fails because the schema is unknown
		_, err := GET(ctx, emptyURL)
		if err == nil {
			t.Fatal("An error was expected")
		}
	})
}

// TestGETWithBaseURLSimple performs a simple httpx.GETWithBaseURL test
func TestGETWithBaseURLSimple(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		ctx := context.Background()
		_, err := GETWithBaseURL(ctx, baseURL, "/status/200")
		if err != nil {
			t.Fatal(err)
		}
	})
}

// TestGETWithBaseURLError checks whether httpx.GETWithBaseURL correctly handles errors
func TestGETWithBaseURLError(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		ctx := context.Background()
		const invalidURL = "\t" // URL parsing fails because there are control characters
		_, err := GETWithBaseURL(ctx, invalidURL, "/foobar")
		if err == nil {
			t.Fatal("An error was expected")
		}
	})
}

// TestPOSTSimple performs a simple httpx.POST test
func TestPOSTSimple(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		ctx := context.Background()
		_, err := POST(ctx, baseURL+"/status/200", "text/plain", []byte("1234"))
		if err != nil {
			t.Fatal(err)
		}
	})
}

// TestPOSTError checks whether httpx.POST correctly handles errors
func TestPOSTError(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		ctx := context.Background()
		const emptyURL = "" // The code fails because the schema is unknown
		_, err := POST(ctx, emptyURL, "text/plain", []byte("1234"))
		if err == nil {
			t.Fatal("An error was expected")
		}
	})
}

// TestPOSTWithBaseURLSimple performs a simple httpx.POSTWithBaseURL test
func TestPOSTWithBaseURLSimple(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		ctx := context.Background()
		_, err := POSTWithBaseURL(ctx, baseURL, "/status/200", "text/plain", []byte("1234"))
		if err != nil {
			t.Fatal(err)
		}
	})
}

// TestPOSTWithBaseURLError checks whether httpx.POSTWithBaseURL correctly handles errors
func TestPOSTWithBaseURLError(t *testing.T) {
	withHTTPBin(t, func(baseURL string) {
		ctx := context.Background()
		const invalidURL = "\t" // URL parsing fails because there are control characters
		_, err := POSTWithBaseURL(ctx, invalidURL, "/foobar", "text/plain", []byte("1234"))
		if err == nil {
			t.Fatal("An error was expected")
		}
	})
}
