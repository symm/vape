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

	resCh, errCh := make(chan CheckResult), make(chan error)
	vape := NewVape(DefaultClient, baseURL, resCh, errCh)
	vape.Run(statusCodeChecks)
}
