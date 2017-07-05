package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	vapeFile := flag.String("config", "Vapefile", "The full path to the Vape configuration file")
	insecureSSL := flag.Bool("skip-ssl-verification", false, "Ignore bad SSL certs")
	flag.Parse()

	start := time.Now()

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

	httpClient := NewHTTPClient(*insecureSSL)
	vape := NewVape(httpClient, baseURL, resCh, errCh)
	vape.Process(smokeTests)

	var passedCount int
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

	fmt.Printf("\n✨  [%d/%d] tests passed in %s\n", passedCount, testsLen, time.Since(start))
	if passedCount < testsLen {
		fmt.Println("🔥  Some tests failed. You may have a bad deployment")
		os.Exit(2)
	}
}
