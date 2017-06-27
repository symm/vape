package main

import "net/http"

// HTTPClient is a custom interface that wraps the basic Get method.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}
