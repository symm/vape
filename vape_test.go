package main

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

type mockHTTPClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m mockHTTPClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}

var httpClient = new(mockHTTPClient)

var test = SmokeTest{
	URI:                "test",
	ExpectedStatusCode: 200,
}

func TestProcess(t *testing.T) {
	resCh, errCh := make(chan SmokeTestResult, 1), make(chan error, 1)
	baseURL, err := url.Parse("http://base.url")
	if err != nil {
		t.Fatal(err)
	}
	vape := NewVape(httpClient, baseURL, resCh, errCh)

	t.Run("TestHTTPErrorResult", func(t *testing.T) {
		httpClient.GetFunc = func(url string) (*http.Response, error) {
			return nil, errors.New("HTTP error")
		}
		vape.Process(SmokeTests{test})
		select {
		case <-resCh:
			t.Error("expected to recieve on error chan, not result chan")
		case <-errCh:
		}
	})

	t.Run("TestHTTPErrorResult", func(t *testing.T) {
		httpClient.GetFunc = func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
			}, nil
		}
		vape.Process(SmokeTests{test})
		select {
		case <-resCh:
		case <-errCh:
			t.Error("expected to recieve on result chan, not error chan")
		}
	})
}

func TestPerformTest(t *testing.T) {
	baseURL, err := url.Parse("http://base.url")
	if err != nil {
		t.Fatal(err)
	}
	vape := NewVape(httpClient, baseURL, nil, nil)

	t.Run("TestHTTPGetError", func(t *testing.T) {
		httpClient.GetFunc = func(url string) (*http.Response, error) {
			return nil, errors.New("HTTP error")
		}

		_, err := vape.performTest(test)
		if err == nil {
			t.Error("expected error: 'HTTP error', got: nil")
		}
	})

	t.Run("TestHTTPGetSuccess", func(t *testing.T) {
		httpClient.GetFunc = func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
			}, nil
		}

		result, err := vape.performTest(test)
		if err != nil {
			t.Errorf("expected error: nil, got: %v", err)
		}
		if result.ActualStatusCode != 200 {
			t.Errorf("expected status code: 200, got: %d", result.ActualStatusCode)
		}
	})
}
