package main

import (
	"net/http"
	"testing"
)

type mockHTTPClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m mockHTTPClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}

func TestVapeRun(t *testing.T) {

}
