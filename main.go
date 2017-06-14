package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
)

// StatusCodeCheck is a check to perform on the server
type StatusCodeCheck struct {
	URL                string
	ExpectedStatusCode int
}

// CheckResult is the result of a StatusCodeCheck
type CheckResult struct {
	Check            StatusCodeCheck
	ActualStatusCode int
	Pass             bool
}

func getStatusCodeChecks() []StatusCodeCheck {
	// TODO: read these in from a config file
	return []StatusCodeCheck{
		StatusCodeCheck{URL: "http://localhost:8000/", ExpectedStatusCode: 200},
		StatusCodeCheck{URL: "http://localhost:8000/missing", ExpectedStatusCode: 404},
		//	StatusCodeCheck{URL: "http://localhost:8000/chicken", ExpectedStatusCode: 200},
		//StatusCodeCheck{URL: "http://localhost:4444/chicken", ExpectedStatusCode: 200},
	}
}

func main() {
	statusCodeChecks := getStatusCodeChecks()

	resc, errc := make(chan CheckResult), make(chan error)

	for _, check := range statusCodeChecks {
		go func(check StatusCodeCheck) {
			result, err := performCheck(check)

			if err != nil {
				errc <- err
				return
			}

			resc <- result
		}(check)
	}

	pass := color.New(color.FgGreen).PrintfFunc()
	fail := color.New(color.FgRed).PrintfFunc()

	failed := false
	for i := 0; i < len(statusCodeChecks); i++ {
		select {
		case res := <-resc:
			message := fmt.Sprintf("%s (expected: %d, actual: %d)\n", res.Check.URL, res.Check.ExpectedStatusCode, res.ActualStatusCode)

			if res.Pass == true {
				pass(message)
			} else {
				failed = true
				fail(message)
			}

		case err := <-errc:
			fail(fmt.Sprintf("%s\n", err))
			failed = true
		}
	}

	if failed {
		os.Exit(1)
	}

	os.Exit(0)
}

func performCheck(check StatusCodeCheck) (CheckResult, error) {
	resp, err := http.Get(check.URL)

	if err != nil {
		return CheckResult{}, err
	}

	wasSuccessful := check.ExpectedStatusCode == resp.StatusCode

	return CheckResult{ActualStatusCode: resp.StatusCode, Check: check, Pass: wasSuccessful}, nil
}
