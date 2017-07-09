package main

import (
	"bytes"
	"errors"
	"io/ioutil"
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

func TestSmokeTestResult(t *testing.T) {
	t.Run("PassedWhenStatusCodeMatches", func(t *testing.T) {
		result := SmokeTestResult{
			Test: SmokeTest{
				ExpectedStatusCode: 200,
			},
			ActualStatusCode: 200,
		}

		if result.Passed() != true {
			t.Errorf("expected pass: true, got: %v", result.Passed())
		}
	})

	t.Run("FailedWhenStatusCodeDifferent", func(t *testing.T) {
		result := SmokeTestResult{
			Test: SmokeTest{
				ExpectedStatusCode: 200,
			},
			ActualStatusCode: 404,
		}

		if result.Passed() != false {
			t.Errorf("expected pass: false, got: %v", result.Passed())
		}
	})

	t.Run("PassedWhenStatusCodeAndContentMatches", func(t *testing.T) {
		result := SmokeTestResult{
			Test: SmokeTest{
				ExpectedStatusCode: 200,
				Content:            "Hello",
			},
			ActualStatusCode: 200,
			ActualContent:    []byte("Hello"),
		}

		if result.Passed() != true {
			t.Errorf("expected pass: true, got: %v", result.Passed())
		}
	})

	t.Run("FailedWhenStatusCodeMatchesButContentDoesnt", func(t *testing.T) {
		result := SmokeTestResult{
			Test: SmokeTest{
				ExpectedStatusCode: 200,
				Content:            "Hello",
			},
			ActualStatusCode: 200,
			ActualContent:    []byte("Goodbye"),
		}

		if result.Passed() != false {
			t.Errorf("expected pass: false, got: %v", result.Passed())
		}
	})

	t.Run("FailedWhenStatusCodeDifferentButContentMatches", func(t *testing.T) {
		result := SmokeTestResult{
			Test: SmokeTest{
				ExpectedStatusCode: 200,
				Content:            "Hello",
			},
			ActualStatusCode: 404,
			ActualContent:    []byte("Hello"),
		}

		if result.Passed() != false {
			t.Errorf("expected pass: false, got: %v", result.Passed())
		}
	})
}

func TestPerformTest(t *testing.T) {
	baseURL, err := url.Parse("http://base.url")
	if err != nil {
		t.Fatal(err)
	}
	vape := NewVape(httpClient, baseURL, 3)

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

	t.Run("TestHttpGetSuccessWithMatchingContent", func(t *testing.T) {
		httpClient.GetFunc = func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("Hello World")),
			}, nil
		}

		var test = SmokeTest{
			URI:                "test",
			ExpectedStatusCode: 200,
			Content:            "Hello World",
		}

		result, err := vape.performTest(test)
		if err != nil {
			t.Errorf("expected error: nil, got: %v", err)
		}
		if result.ActualStatusCode != 200 {
			t.Errorf("expected status code: 200, got: %d", result.ActualStatusCode)
		}
		if string(result.ActualContent) != "Hello World" {
			t.Errorf("expected content: Hello World, got: %s", result.ActualContent)
		}
		if result.Passed() != true {
			t.Errorf("expected pass: true, got: %v", result.Passed())
		}

	})

	t.Run("TestHttpGetSuccessWithNonMatchingContent", func(t *testing.T) {
		httpClient.GetFunc = func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("Not the message you are looking for")),
			}, nil
		}

		var test = SmokeTest{
			URI:                "test",
			ExpectedStatusCode: 200,
			Content:            "Hello World",
		}

		result, err := vape.performTest(test)
		if err != nil {
			t.Errorf("expected error: nil, got: %v", err)
		}
		if result.ActualStatusCode != 200 {
			t.Errorf("expected status code: 200, got: %d", result.ActualStatusCode)
		}
		if string(result.ActualContent) != "Not the message you are looking for" {
			t.Errorf("expected content: Not the message you are looking for, got: %s", result.ActualContent)
		}

		if result.Passed() != false {
			t.Errorf("expected pass: false, got: %v", result.Passed())
		}
	})
}
