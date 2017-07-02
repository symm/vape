package main

import (
	"io/ioutil"
	"net/url"
	"path"
	"strings"
)

// SmokeTest contains a URI and expected status code.
type SmokeTest struct {
	URI                string `json:"uri"`
	ExpectedStatusCode int    `json:"expected_status_code"`
	Content            string `json:"content"`
}

// SmokeTestResult is the result of a SmokeTest.
type SmokeTestResult struct {
	Test             SmokeTest
	ActualStatusCode int
	ActualContent    []byte
}

// Passed determines if the SmokeTest passed successfully
func (result SmokeTestResult) Passed() bool {
	if result.Test.Content != "" {
		return result.contentMatched()
	}

	return result.statusCodeMatched()
}

func (result SmokeTestResult) statusCodeMatched() bool {
	return (result.Test.ExpectedStatusCode == result.ActualStatusCode)
}

func (result SmokeTestResult) contentMatched() bool {
	return (strings.Contains(string(result.ActualContent), result.Test.Content) == false)
}

// SmokeTests is a slice of smoke tests to perform.
type SmokeTests []SmokeTest

// Vape contains dependencies used to run the application.
type Vape struct {
	client  HTTPClient
	baseURL *url.URL
	resCh   chan SmokeTestResult
	errCh   chan error
}

// NewVape builds a Vape from the given dependencies.
func NewVape(client HTTPClient, baseURL *url.URL, resCh chan SmokeTestResult, errCh chan error) Vape {
	return Vape{
		client:  client,
		baseURL: baseURL,
		resCh:   resCh,
		errCh:   errCh,
	}
}

// Process takes a list of URIs and concurrently performs a smoke test on each.
func (v Vape) Process(SmokeTests SmokeTests) {
	// TODO: limit the numer of concurrent requests so we don't DoS the server
	go func() {
		for _, test := range SmokeTests {
			go func(test SmokeTest) {
				result, err := v.performTest(test)
				if err != nil {
					v.errCh <- err
					return
				}
				v.resCh <- result
			}(test)
		}
	}()
}

// performTest tests the status code of a HTTP request of a given URI.
func (v Vape) performTest(test SmokeTest) (SmokeTestResult, error) {
	url := *v.baseURL
	url.Path = path.Join(url.Path, test.URI)

	resp, err := v.client.Get(url.String())
	if err != nil {
		return SmokeTestResult{}, err
	}

	result := SmokeTestResult{
		ActualStatusCode: resp.StatusCode,
		Test:             test,
	}

	if test.Content != "" {
		defer resp.Body.Close()
		actualContent, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return result, err
		}

		result.ActualContent = actualContent
	}

	return result, nil
}
