package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func getStatusCodeChecks(vapeFile string) []StatusCodeCheck {
	file, err := os.Open(vapeFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	checks := []StatusCodeCheck{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		url := words[0]
		code, err := strconv.Atoi(words[1])
		if err != nil {
			log.Fatal(err)
		}
		checks = append(checks, StatusCodeCheck{URL: url, ExpectedStatusCode: code})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return checks
}

func main() {
	statusCodeChecks := getStatusCodeChecks(".smoke")

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
