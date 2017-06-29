package main

import (
	"net/http"
	"time"
)

// HTTPClient is a custom interface that wraps the basic Get method.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultClient returns a HTTP client with configured timeouts.
var DefaultClient = &http.Client{
	Timeout: time.Duration(5 * time.Second),
}
