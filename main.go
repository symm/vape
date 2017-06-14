package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// Check is a check to perform on the server
type StatusCodeCheck struct {
	URL        string
	StatusCode int
}

type CheckResult struct {
	Check      StatusCodeCheck
	StatusCode int
	Pass       bool
}

func main() {
	urls := []StatusCodeCheck{
		StatusCodeCheck{URL: "http://mysite.com/health", StatusCode: 200},
		StatusCodeCheck{URL: "http://mysite.com/fake-page-chicken", StatusCode: 404},
	}

	fmt.Println("Smoke tester")
	fmt.Println("------------")

	results := []CheckResult{}
	// TODO: goroutines send this all at once
	for _, check := range urls {
		results = append(results, performCheck(check))
	}

	fmt.Println("All pages checked, here are the results:")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"URL", "Expected Status Code", "Actual Status Code", "PASS?"})
	for _, result := range results {
		table.Append([]string{
			result.Check.URL,
			strconv.Itoa(result.Check.StatusCode),
			strconv.Itoa(result.StatusCode),
			strconv.FormatBool(result.Pass)})
	}

	table.Render()
}

func performCheck(check StatusCodeCheck) CheckResult {
	fmt.Printf("[+] Checking %s\n", check.URL)
	wasSuccessful := check.StatusCode == 500

	return CheckResult{StatusCode: 500, Check: check, Pass: wasSuccessful}
}
