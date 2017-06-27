package main

import (
	"fmt"
	"net/url"
	"path"
)

// Vape ...
type Vape struct {
	client  HTTPClient
	baseURL *url.URL
}

// NewVape ...
func NewVape(client HTTPClient, baseURL *url.URL) Vape {
	return Vape{
		client:  client,
		baseURL: baseURL,
	}
}

// Run ...
func (v Vape) Run(statusCodeChecks StatusCodeChecks) {
	resCh, errCh := make(chan CheckResult), make(chan error)

	// TODO: limit the numer of concurrent requests so we don't DoS the server
	for _, check := range statusCodeChecks {
		go func(check StatusCodeCheck) {
			result, err := v.performCheck(check)
			if err != nil {
				errCh <- err
				return
			}
			resCh <- result
		}(check)
	}

	for i := 0; i < len(statusCodeChecks); i++ {
		select {
		case res := <-resCh:
			fmt.Printf("%s (expected: %d, actual: %d)\n", res.Check.URI, res.Check.ExpectedStatusCode, res.ActualStatusCode)
		case err := <-errCh:
			fmt.Println(err)
		}
	}
}

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
