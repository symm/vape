package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func tmpFile(data string) (*os.File, func(), error) {
	tmpfile, err := ioutil.TempFile("", "Vapefile")
	if err != nil {
		return nil, nil, err
	}
	if _, err = tmpfile.Write([]byte(data)); err != nil {
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

func TestReadVapeFileFileDoesNotExist(t *testing.T) {
	_, err := parseVapefile("dummy.file")

	if err == nil {
		t.Error("expected error accessing JSON file, got: nil")
	}

}

func TestReadVapefileInvalidJSON(t *testing.T) {
	tmpfile, cleanup, err := tmpFile("invalid JSON")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_, err = parseVapefile(tmpfile.Name())
	if err == nil {
		t.Error("expected error reading JSON, got: nil")
	}
}

func TestReadVapefileSuccess(t *testing.T) {
	json := `[
  {
    "uri": "/status/200",
    "expected_status_code": 200
  },
  {
    "uri": "/status/500",
    "expected_status_code": 500
  }
]`
	tmpfile, cleanup, err := tmpFile(json)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_, err = parseVapefile(tmpfile.Name())
	if err != nil {
		t.Errorf("expected error: nil, got: %v", err)
	}
}

func TestReadVapefileMissingFields(t *testing.T) {
	t.Run("TestUriMissing", func(t *testing.T) {
		json := `[
	  {
	    "expected_status_code": 200
	  }
	]`
		tmpfile, cleanup, err := tmpFile(json)
		if err != nil {
			t.Fatal(err)
		}
		defer cleanup()
		_, err = parseVapefile(tmpfile.Name())
		if err == nil {
			t.Errorf("expected error: got: %v", err)
		}

		expectedError := "Each test should have at least a uri and expected_status_code"
		if strings.Contains(err.Error(), expectedError) == false {
			t.Errorf("expected message: %s got %s", expectedError, err.Error())
		}
	})

	t.Run("TestExpectedStatusCodeMissing", func(t *testing.T) {
		t.Run("TestUriMissing", func(t *testing.T) {
			json := `[
			{
				"uri": "/health"
			}
		]`
			tmpfile, cleanup, err := tmpFile(json)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()
			_, err = parseVapefile(tmpfile.Name())
			if err == nil {
				t.Errorf("expected error: got: %v", err)
			}

			expectedError := "Each test should have at least a uri and expected_status_code"
			if strings.Contains(err.Error(), expectedError) == false {
				t.Errorf("expected message: %s got %s", expectedError, err.Error())
			}
		})
	})
}

func TestParseBaseURL(t *testing.T) {
	t.Run("TestInvalidURL", func(t *testing.T) {
		url := ":"
		_, err := parseBaseURL(url)
		if err == nil {
			t.Error("expected error parsing invalid URL, got: nil")
		}
	})

	t.Run("TestInvalidURLScheme", func(t *testing.T) {
		url := "test.com"
		_, err := parseBaseURL(url)
		if err == nil {
			t.Error("expected error parsing invalid URL scheme, got: nil")
		}
	})

	t.Run("TestValidURL", func(t *testing.T) {
		url := "http://test.com"
		_, err := parseBaseURL(url)
		if err != nil {
			t.Errorf("expected error: nil, got: %v", err)
		}
	})
}

func TestFormatResult(t *testing.T) {
	result := SmokeTestResult{
		Test: SmokeTest{
			URI:                "/health",
			ExpectedStatusCode: 200,
		},
		ActualStatusCode: 200,
	}

	t.Run("TestSuccess", func(t *testing.T) {
		output := formatResult(result)
		expected := "\033[32m✓ [200:200] /health\033[0m"
		if output != expected {
			t.Errorf("expected output: %s, got: %s", expected, output)
		}
	})

	result = SmokeTestResult{
		Test: SmokeTest{
			URI:                "/health",
			ExpectedStatusCode: 200,
		},
		ActualStatusCode: 500,
	}

	t.Run("TestFail", func(t *testing.T) {
		output := formatResult(result)
		expected := "\033[31m✘ [200:500] /health\033[0m"
		if output != expected {
			t.Errorf("expected output: %s, got: %s", expected, output)
		}
	})
}
