package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	client Vape
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// client
	httpClient := NewHTTPClient()
	url, _ := url.Parse(server.URL)
	client = NewVape(httpClient, url)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestMatchingStatusCodePasses(t *testing.T) {
	setup()
	defer teardown()

	check := StatusCodeCheck{
		ExpectedStatusCode: 200,
		URI:                "/health",
	}

	result, err := client.performCheck(check)

	if err != nil {
		t.Error(err)
	}

	if result.Pass != true {
		t.Errorf("expected Pass: true, got: %v", result.Pass)
	}

	if result.ActualStatusCode != 200 {
		t.Errorf("expected ActualStatusCode: 200, got: %v", result.ActualStatusCode)
	}

	if result.Check != check {
		t.Errorf("expected check to match original, got: %v", result.Check)
	}
}

func TestDifferentStatusCodeFails(t *testing.T) {
	setup()
	defer teardown()

	check := StatusCodeCheck{
		ExpectedStatusCode: 404,
		URI:                "/health",
	}

	result, err := client.performCheck(check)

	if err != nil {
		t.Error(err)
	}

	if result.Pass != false {
		t.Errorf("expected Pass: false, got: %v", result.Pass)
	}

	if result.ActualStatusCode != 200 {
		t.Errorf("expected ActualStatusCode: 200, got: %v", result.ActualStatusCode)
	}

	if result.Check != check {
		t.Errorf("expected check to match original, got: %v", result.Check)
	}
}
