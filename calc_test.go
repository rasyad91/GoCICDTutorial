package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIncHandler(t *testing.T) {
	testSuite := []struct {
		name     string
		action   string
		value1   string
		value2   string
		expected string
		status   int
		err      string
	}{
		{name: "Addition Test Handler", value1: "2", value2: "3", action: "Add", expected: "5"},
		{name: "Addition Test Handler", value1: "2", value2: "-1", action: "Add", expected: "1"},
	}

	for _, testCase := range testSuite {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("localhost:8080/?action=%s&val1=%s&val2=%s",
				testCase.action, testCase.value1, testCase.value2), nil)
			if err != nil {
				t.Fatalf("Cannot create request: %v", err)
			}
			w := httptest.NewRecorder()
			incHandler(w, req)

			result := w.Result()
			defer result.Body.Close()

			body, err := ioutil.ReadAll(result.Body)
			if err != nil {
				t.Fatalf("response body cannot read: %v", err)
			}

			if testCase.err != "" {
				if result.StatusCode != http.StatusBadRequest {
					t.Errorf("expected bad request; got %v", result.StatusCode)
				}
				if msg := string(bytes.TrimSpace(body)); msg != testCase.err {
					t.Errorf("expected message %q but got %q", testCase.err, msg)
				}
				return
			}

			if result.StatusCode != http.StatusOK {
				t.Errorf("expected status OK but got %v", result.Status)
			}
			data := string(bytes.TrimSpace(body))

			if data != testCase.expected {
				t.Fatalf("expected result of %v but got %v", testCase.expected, data)
			}

		})
	}
}

func TestServerRouting(t *testing.T) {
	testSuite := []struct {
		name     string
		action   string
		value1   string
		value2   string
		expected string
		status   int
		err      string
	}{
		{name: "Addition Routing", value1: "2", value2: "3", action: "Add", expected: "5"},
	}

	server := httptest.NewServer(handlers())
	defer server.Close()
	for _, testCase := range testSuite {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := http.Get(fmt.Sprintf( /*"%s/?%s%s%s"*/ "%s/?action=%s&val1=%s&val2=%s",
				server.URL, testCase.action, testCase.value1, testCase.value2))
			if err != nil {
				t.Fatalf("Get request error: %v", err)
			}
			defer result.Body.Close()

			body, err := ioutil.ReadAll(result.Body)

			if err != nil {
				t.Fatalf("response body cannot read: %v", err)
			}

			if testCase.err != "" {
				if result.StatusCode != http.StatusBadRequest {
					t.Errorf("expected bad request; got %v", result.StatusCode)
				}
				if msg := string(bytes.TrimSpace(body)); msg != testCase.err {
					t.Errorf("expected message %q but got %q", testCase.err, msg)
				}
				return
			}

			if result.StatusCode != http.StatusOK {
				t.Errorf("expected status OK but got %v", result.Status)
			}
		})
	}
}
