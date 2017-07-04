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

var defaultTransport = &http.Transport{
	// TODO: allow this to be toggled
	// So we can still smoke test sites with bad ssl cert
	// e.g. ./vape https://self-signed.badssl.com/
	TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
}

// DefaultClient returns a HTTP client with configured timeouts.
var DefaultClient = &http.Client{
	Timeout:   time.Duration(5 * time.Second),
	Transport: defaultTransport,
}
