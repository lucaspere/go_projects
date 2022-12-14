package server

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Table-Drive Testing
// See more https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
func TestServer(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "index",
			path:     "/api",
			expected: "Hello, world!",
		},
		{
			name:     "healthcheck",
			path:     "/healthz",
			expected: "ok",
		},
	}

	mux := http.NewServeMux()
	setupHandlers(mux)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Get(ts.URL + tc.path)
			if err != nil {
				log.Fatal(err)
			}
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			// Asserts
			if string(resBody) != tc.expected {
				t.Errorf(
					"Expected: %s, Got: %s",
					tc.expected, string(resBody),
				)
			}
		})
	}

}
