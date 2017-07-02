package main

import (
	"net/url"
	"path"
)

// SmokeTest contains a URI and expected status code.
type SmokeTest struct {
	URI                string `json:"uri"`
	ExpectedStatusCode int    `json:"expected_status_code"`
}

// SmokeTestResult is the result of a SmokeTest.
type SmokeTestResult struct {
	Test             SmokeTest
	ActualStatusCode int
	Pass             bool
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

	return SmokeTestResult{
		ActualStatusCode: resp.StatusCode,
		Test:             test,
		Pass:             test.ExpectedStatusCode == resp.StatusCode,
	}, nil
}
