package main

import (
	"crypto/tls"
	"net/http"
	"time"
)

// HTTPClient is a custom interface that wraps the basic Get method.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

// NewHTTPClient returns a configrued HTTP client.
func NewHTTPClient(insecureSSL bool) *http.Client {
	return &http.Client{
		Timeout:   time.Duration(5 * time.Second),
		Transport: newHTTPTransport(insecureSSL),
	}
}

func newHTTPTransport(insecureSSL bool) *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSSL,
		},
	}
}
