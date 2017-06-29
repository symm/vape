package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

// parseVapefile reads a given Vapefile and returns the contents.
func parseVapefile(file string) (StatusCodeChecks, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var checks StatusCodeChecks
	err = json.Unmarshal(raw, &checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

// parseBaseURL checks a given URL is valid.
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

// parseOutput colours a given message.
func parseOutput(message string, pass bool) string {
	colour := 32
	if !pass {
		colour = 31
	}
	return fmt.Sprintf("\033[%dm%s\033[0m", colour, message)
}
