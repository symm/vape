package main

import (
	"crypto/tls"
	"net/http"
	"time"
)

// HTTPClient is a custom interface that wraps the basic Get method.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// NewHTTPClient returns a configrued HTTP client.
func NewHTTPClient(sslSkip bool) *http.Client {
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	if sslSkip {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return client
}
