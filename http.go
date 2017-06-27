package main

import (
	"net/http"
	"time"
)

// HTTPClient is a custom interface that wraps the basic Get method.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// NewHTTPClient returns a HTTP client with configured timeouts.
func NewHTTPClient() HTTPClient {
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	return client
}
