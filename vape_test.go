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

func TestPerformCheck(t *testing.T) {
	var httpClient = new(mockHTTPClient)
	baseURL, err := url.Parse("http://base.url")
	if err != nil {
		t.Fatal(err)
	}
	check := StatusCodeCheck{
		URI:                "test",
		ExpectedStatusCode: 200,
	}

	t.Run("TestHTTPGetError", func(t *testing.T) {
		httpClient.GetFunc = func(url string) (*http.Response, error) {
			return nil, errors.New("HTTP error")
		}

		vape := NewVape(httpClient, baseURL, nil, nil)
		_, err := vape.performCheck(check)
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

		vape := NewVape(httpClient, baseURL, nil, nil)
		result, err := vape.performCheck(check)
		if err != nil {
			t.Errorf("expected error: nil, got: %v", err)
		}
		if result.ActualStatusCode != 200 {
			t.Errorf("expected status code: 200, got: %d", result.ActualStatusCode)
		}
	})
}
