package main

import (
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	t.Run("TestHttpTimeoutIsSet", func(t *testing.T) {
		client := NewHTTPClient(true)

		expectedTimeout := time.Duration(5 * time.Second)
		actualTimeout := client.Timeout
		if expectedTimeout != actualTimeout {
			t.Errorf("expected timeout: %v, got: %v", expectedTimeout, actualTimeout)
		}
	})
}

func TestNewHTTPTransport(t *testing.T) {
	t.Run("TestSSLVerificationOn", func(t *testing.T) {
		transport := newHTTPTransport(false)

		if transport.TLSClientConfig.InsecureSkipVerify != false {
			t.Errorf("expected InsecureSkipVerify: false, got: %v", transport.TLSClientConfig.InsecureSkipVerify)
		}
	})

	t.Run("TestSSLVerificationOff", func(t *testing.T) {
		transport := newHTTPTransport(true)

		if transport.TLSClientConfig.InsecureSkipVerify != true {
			t.Errorf("expected InsecureSkipVerify: true, got: %v", transport.TLSClientConfig.InsecureSkipVerify)
		}
	})
}
