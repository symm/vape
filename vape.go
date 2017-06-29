package main

import (
	"net/url"
	"path"
)

// StatusCodeCheck contains a URI and expected status code.
type StatusCodeCheck struct {
	URI                string `json:"uri"`
	ExpectedStatusCode int    `json:"expected_status_code"`
}

// CheckResult is the result of a StatusCodeCheck.
type CheckResult struct {
	Check            StatusCodeCheck
	ActualStatusCode int
	Pass             bool
}

// StatusCodeChecks is a slice of checks to perform.
type StatusCodeChecks []StatusCodeCheck

// Vape contains dependencies used to run the application.
type Vape struct {
	client  HTTPClient
	baseURL *url.URL
	resCh   chan CheckResult
	errCh   chan error
}

// NewVape builds a Vape from the given dependencies.
func NewVape(client HTTPClient, baseURL *url.URL, resCh chan CheckResult, errCh chan error) Vape {
	return Vape{
		client:  client,
		baseURL: baseURL,
		resCh:   resCh,
		errCh:   errCh,
	}
}

// Process takes a list of URIs and concurrently performs a smoke test on each.
func (v Vape) Process(statusCodeChecks StatusCodeChecks) {
	// TODO: limit the numer of concurrent requests so we don't DoS the server
	for _, check := range statusCodeChecks {
		go func(check StatusCodeCheck) {
			result, err := v.performCheck(check)
			if err != nil {
				v.errCh <- err
				return
			}
			v.resCh <- result
		}(check)
	}
}

// performCheck checks the status code of a HTTP request of a given URI.
func (v Vape) performCheck(check StatusCodeCheck) (CheckResult, error) {
	url := *v.baseURL
	url.Path = path.Join(url.Path, check.URI)

	resp, err := v.client.Get(url.String())
	if err != nil {
		return CheckResult{}, err
	}

	return CheckResult{
		ActualStatusCode: resp.StatusCode,
		Check:            check,
		Pass:             check.ExpectedStatusCode == resp.StatusCode,
	}, nil
}
