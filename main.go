package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
)

// StatusCodeCheck is a check to perform on the server
type StatusCodeCheck struct {
	URI                string `json:"uri"`
	ExpectedStatusCode int    `json:"expectedStatusCode"`
}

// CheckResult is the result of a StatusCodeCheck
type CheckResult struct {
	Check            StatusCodeCheck
	ActualStatusCode int
	Pass             bool
}

func getStatusCodeChecks(vapeFile string) []StatusCodeCheck {
	raw, err := ioutil.ReadFile(vapeFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var checks []StatusCodeCheck
	//TODO: make sure we've got a valid URI and ExpectedStatusCode
	json.Unmarshal(raw, &checks)

	return checks
}

func usage() {
	fmt.Printf("Usage: ./vape {BaseURL}\n\n")
	fmt.Println("Where {BaseURL} e.g. https://example.com")
	os.Exit(0)
}

func main() {

	if len(os.Args) != 2 {
		usage()
	}

	// TODO: verify it's a proper good baseURL
	baseURL := os.Args[1]

	statusCodeChecks := getStatusCodeChecks("Vapefile")

	resc, errc := make(chan CheckResult), make(chan error)

	// TODO: limit the numer of concurrent requests so we don't DoS the server
	for _, check := range statusCodeChecks {
		go func(check StatusCodeCheck) {
			result, err := performCheck(baseURL, check)

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
			message := fmt.Sprintf("%s (expected: %d, actual: %d)\n", res.Check.URI, res.Check.ExpectedStatusCode, res.ActualStatusCode)

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

func performCheck(baseURL string, check StatusCodeCheck) (CheckResult, error) {
	resp, err := http.Get(baseURL + check.URI)

	if err != nil {
		return CheckResult{}, err
	}

	wasSuccessful := check.ExpectedStatusCode == resp.StatusCode

	return CheckResult{ActualStatusCode: resp.StatusCode, Check: check, Pass: wasSuccessful}, nil
}
