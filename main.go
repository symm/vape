package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var vapeFile = flag.String("config", "Vapefile", "The full path to the Vape configuration file")
var insecureSSL = flag.Bool("skip-ssl-verification", false, "")

func main() {
	start := time.Now()

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

	smokeTests, err := parseVapefile(*vapeFile)
	if err != nil {
		fmt.Printf("Error reading config: %s", err.Error())
		os.Exit(1)
	}

	testsLen := len(smokeTests)
	resCh, errCh := make(chan SmokeTestResult, testsLen), make(chan error, testsLen)

	client := DefaultClient

	if *insecureSSL == true {
		client = InsecureClient
	}

	vape := NewVape(client, baseURL, resCh, errCh)
	vape.Process(smokeTests)

	passedCount := 0

	for i := 0; i < testsLen; i++ {
		select {
		case res := <-resCh:
			if res.Passed() {
				passedCount++
			}

			fmt.Println(formatResult(res))
		case err := <-errCh:
			fmt.Println(err)
		}
	}

	if passedCount < testsLen {
		fmt.Println("\n🔥  Some tests failed. You may have a bad deployment")
		os.Exit(2)
	}

	fmt.Printf("\n✨  [%d/%d] tests passed in %s\n", passedCount, testsLen, time.Since(start))
}
