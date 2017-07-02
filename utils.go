package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

// parseVapefile reads a given Vapefile and returns the contents.
func parseVapefile(file string) (SmokeTests, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var tests SmokeTests
	err = json.Unmarshal(raw, &tests)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

// parseBaseURL tests a given URL is valid.
func parseBaseURL(baseURL string) (*url.URL, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("invalid protocol scheme")
	}
	return u, nil
}

// formatResult returns a readable string summarizing the result
func formatResult(result SmokeTestResult) string {
	message := fmt.Sprintf("[%d:%d] %s", result.Test.ExpectedStatusCode, result.ActualStatusCode, result.Test.URI)

	if result.Test.Content != "" {
		message = fmt.Sprintf("%s %s", message, result.Test.Content)
	}

	colour := 32
	icon := '✓'

	if !result.Passed() {
		icon = '✘'
		colour = 31
	}

	return fmt.Sprintf("\033[%dm%c %s\033[0m", colour, icon, message)
}
