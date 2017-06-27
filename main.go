package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// StatusCodeCheck is a check to perform on the server
type StatusCodeCheck struct {
	URI                string `json:"uri"`
	ExpectedStatusCode int    `json:"expected_status_code"`
}

// CheckResult is the result of a StatusCodeCheck
type CheckResult struct {
	Check            StatusCodeCheck
	ActualStatusCode int
	Pass             bool
}

// StatusCodeChecks ...
type StatusCodeChecks []StatusCodeCheck

func main() {
	if len(os.Args) != 2 {
		fmt.Println("vape: usage: no base URL specified")
		os.Exit(0)
	}

	baseURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println("vape: invalid base URL")
		os.Exit(0)
	}

	statusCodeChecks, err := readVapefile(vapefile)
	if err != nil {
		log.Fatal(err)
	}

	vape := NewVape(http.DefaultClient, baseURL)
	vape.Run(statusCodeChecks)
}
