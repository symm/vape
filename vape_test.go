package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

var httpClient = new(mockHTTPClient)

var test = SmokeTest{
	URI:                "/test",
	ExpectedStatusCode: 200,
}

func getPassedStatusCodeResult() SmokeTestResult {
	return SmokeTestResult{
		Test: SmokeTest{
			ExpectedStatusCode: 200,
		},
		ActualStatusCode: 200,
	}
}

func getFailedStatusCodeResult() SmokeTestResult {
	return SmokeTestResult{
		Test: SmokeTest{
			ExpectedStatusCode: 200,
		},
		ActualStatusCode: 404,
	}
}

func TestSmokeTestResult(t *testing.T) {
	t.Run("PassedWhenStatusCodeMatches", func(t *testing.T) {
		result := getPassedStatusCodeResult()

		if result.Passed() != true {
			t.Errorf("expected pass: true, got: %v", result.Passed())
		}
	})

	t.Run("FailedWhenStatusCodeDifferent", func(t *testing.T) {
		result := getFailedStatusCodeResult()

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

func TestSmokeTestResults(t *testing.T) {
	t.Run("TestPassedCount", func(t *testing.T) {
		results := SmokeTestResults{
			getPassedStatusCodeResult(),
			getPassedStatusCodeResult(),
		}

		if results.PassedCount() != 2 {
			t.Errorf("expected passed: 0, got: %d", results.PassedCount())
		}

		results = SmokeTestResults{
			getPassedStatusCodeResult(),
			getFailedStatusCodeResult(),
			getFailedStatusCodeResult(),
			getFailedStatusCodeResult(),
		}

		if results.PassedCount() != 1 {
			t.Errorf("expected passed: 1, got: %d", results.PassedCount())
		}

		results = SmokeTestResults{
			getFailedStatusCodeResult(),
			getFailedStatusCodeResult(),
		}

		if results.PassedCount() != 0 {
			t.Errorf("expected passed: 0, got: %d", results.PassedCount())
		}
	})
}

func getVapeClient() Vape {
	baseURL, err := url.Parse("http://base.url")
	if err != nil {
		log.Fatal(err)
	}
	return NewVape(httpClient, baseURL, 3, "")
}

func TestPerformTest(t *testing.T) {
	vape := getVapeClient()

	t.Run("TestHTTPGetError", func(t *testing.T) {
		httpClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("HTTP error")
		}

		_, err := vape.performTest(test)
		if err == nil {
			t.Error("expected error: 'HTTP error', got: nil")
		}
	})

	t.Run("TestHTTPGetSuccess", func(t *testing.T) {
		httpClient.DoFunc = func(req *http.Request) (*http.Response, error) {
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
		httpClient.DoFunc = func(req *http.Request) (*http.Response, error) {
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
		httpClient.DoFunc = func(req *http.Request) (*http.Response, error) {
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

func TestWorker(t *testing.T) {
	t.Run("TestProcessOkay", func(t *testing.T) {
		vape := getVapeClient()

		tests := SmokeTests{
			test,
			test,
			test,
		}
		results, errors := vape.Process(tests)

		if len(results) != 3 {
			t.Errorf("expected result count: 3, got: %v", len(results))
		}

		if len(errors) != 0 {
			t.Errorf("expected error count: 0, got: %v", len(errors))
		}
	})
	t.Run("TestProcessErrors", func(t *testing.T) {
		vape := getVapeClient()
		httpClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("HTTP error")
		}

		tests := SmokeTests{
			test,
			test,
			test,
		}
		results, errors := vape.Process(tests)

		if len(results) != 0 {
			t.Errorf("expected result count: 0, got: %v", len(results))
		}

		if len(errors) != 3 {
			t.Errorf("expected error count: 3, got: %v", len(errors))
		}
	})

}
