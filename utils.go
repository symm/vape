package main

import (
	"encoding/json"
	"io/ioutil"
)

const vapefile = "Vapefile"

func readVapefile(file string) (StatusCodeChecks, error) {
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
