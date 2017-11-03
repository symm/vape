package main

import (
	"io/ioutil"
	"net/url"
	"path"
	"strings"
	"sync"
	"net/http"
)

// SmokeTest contains a URI and expected status code.
type SmokeTest struct {
	URI                string `json:"uri"`
	ExpectedStatusCode int    `json:"status_code"`
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
	if !result.statusCodeMatched() {
		return false
	}

	if result.Test.Content != "" && !result.contentMatched() {
		return false
	}

	return true
}

func (result SmokeTestResult) statusCodeMatched() bool {
	return (result.Test.ExpectedStatusCode == result.ActualStatusCode)
}

func (result SmokeTestResult) contentMatched() bool {
	return (strings.Contains(string(result.ActualContent), result.Test.Content) == true)
}

// SmokeTestResults is a slice of smoke test results
type SmokeTestResults []SmokeTestResult

// PassedCount is the number of smoke tests that passed
func (results SmokeTestResults) PassedCount() int {
	var passedCount int
	for _, result := range results {
		if result.Passed() {
			passedCount++
		}
	}
	return passedCount
}

// SmokeTests is a slice of smoke tests to perform.
type SmokeTests []SmokeTest

// Vape contains dependencies used to run the application.
type Vape struct {
	client        HTTPClient
	baseURL       *url.URL
	concurrency   int
	authorization string
}

// NewVape builds a Vape from the given dependencies.
func NewVape(client HTTPClient, baseURL *url.URL, concurrency int, authorization string) Vape {
	return Vape{
		client:      client,
		baseURL:     baseURL,
		concurrency: concurrency,
		authorization: authorization,
	}
}

func (v Vape) worker(wg *sync.WaitGroup, tests <-chan SmokeTest, resultCh chan<- SmokeTestResult, errorCh chan<- error) {
	for test := range tests {

		result, err := v.performTest(test)

		if err != nil {
			errorCh <- err
		} else {
			resultCh <- result
		}

		wg.Done()
	}
}

// Process takes a list of URIs and concurrently performs a smoke test on each.
func (v Vape) Process(tests SmokeTests) (results SmokeTestResults, errors []error) {
	testCount := len(tests)

	jobCh := make(chan SmokeTest, testCount)
	resultCh := make(chan SmokeTestResult, testCount)
	errorCh := make(chan error, testCount)

	var wg sync.WaitGroup
	for w := 1; w <= v.concurrency; w++ {
		go v.worker(&wg, jobCh, resultCh, errorCh)
	}

	for _, job := range tests {
		jobCh <- job
		wg.Add(1)
	}
	close(jobCh)

	wg.Wait()

	for i := 0; i < testCount; i++ {
		select {
		case err := <-errorCh:
			errors = append(errors, err)
		case result := <-resultCh:
			results = append(results, result)
		}
	}

	return results, errors
}

// performTest tests the status code of a HTTP request of a given URI.
func (v Vape) performTest(test SmokeTest) (SmokeTestResult, error) {
	url := *v.baseURL
	u, err := url.Parse(path.Join(url.Path + test.URI))

	if err != nil {
		return SmokeTestResult{}, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return SmokeTestResult{}, err
	}

	if v.authorization != "" {
		req.Header.Add("Authorization", v.authorization)
	}

	resp, err := v.client.Do(req)

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
