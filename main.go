package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("vape: usage: no base URL specified")
		os.Exit(0)
	}

	baseURL, err := parseBaseURL(os.Args[1])
	if err != nil {
		fmt.Printf("vape: invalid base URL: %v\n", err)
		os.Exit(0)
	}

	statusCodeChecks, err := parseVapefile(vapefile)
	if err != nil {
		log.Fatal(err)
	}

	checksLen := len(statusCodeChecks)
	resCh, errCh := make(chan CheckResult, checksLen), make(chan error, checksLen)
	vape := NewVape(DefaultClient, baseURL, resCh, errCh)
	vape.Process(statusCodeChecks)

	for i := 0; i < checksLen; i++ {
		select {
		case res := <-resCh:
			output := fmt.Sprintf("%s (expected: %d, actual: %d)", res.Check.URI, res.Check.ExpectedStatusCode, res.ActualStatusCode)
			fmt.Println(parseOutput(output, res.Pass))
		case err := <-errCh:
			fmt.Println(err)
		}
	}
}
