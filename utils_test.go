package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func tmpFile(data string) (*os.File, func(), error) {
	tmpfile, err := ioutil.TempFile("", "Vapefile")
	if err != nil {
		return nil, nil, err
	}
	if _, err := tmpfile.Write([]byte(data)); err != nil {
		return nil, nil, err
	}
	if err = tmpfile.Close(); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		os.Remove(tmpfile.Name())
	}

	return tmpfile, cleanup, nil
}

func TestReadVapefileInvalidJSON(t *testing.T) {
	tmpfile, cleanup, err := tmpFile("invalid JSON")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_, err = readVapefile(tmpfile.Name())
	if err == nil {
		t.Error("expected error reading JSON, got: nil")
	}
}

func TestReadVapefilieSuccess(t *testing.T) {
	json := `[
  {
    "uri": "/status/200",
    "expectedStatusCode": 200
  },
  {
    "uri": "/status/500",
    "expectedStatusCode": 500
  }
]`
	tmpfile, cleanup, err := tmpFile(json)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_, err = readVapefile(tmpfile.Name())
	if err != nil {
		t.Errorf("expected error: nil, got: %v", err)
	}
}
