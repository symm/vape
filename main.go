package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var vapeFile = flag.String("config", "Vapefile", "The full path to the Vape configuration file")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("vape: usage: no base URL specified")
		os.Exit(0)
	}

	baseURL, err := parseBaseURL(args[0])
	if err != nil {
		fmt.Printf("vape: invalid base URL: %v\n", err)
		os.Exit(0)
	}

	statusCodeChecks, err := parseVapefile(*vapeFile)
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
