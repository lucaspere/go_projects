package main

import (
	"errors"
	"testing"
)

type testConfig struct {
	args           []string
	expectedErr    error
	expectedConfig config
}

func TestParseArgs(t *testing.T) {
	var c config
	tests := []testConfig{
		{
			args:           []string{},
			expectedErr:    errors.New("must specify one argument"),
			expectedConfig: config{},
		},
		{
			args:           []string{"-h"},
			expectedErr:    errors.New("invalid -h argument. Must specify a valid argument -n"),
			expectedConfig: config{},
		},
		{
			args:           []string{"-n"},
			expectedErr:    errors.New("must insert some int value to -n argument"),
			expectedConfig: config{},
		},
		{
			args:           []string{"-n", "test"},
			expectedErr:    errors.New("test is not a int type"),
			expectedConfig: config{},
		},
		{
			args:           []string{"-n", "2"},
			expectedErr:    nil,
			expectedConfig: config{numTimes: 2},
		},
	}
	for _, tc := range tests {
		err := c.parseArgs(tc.args)

		if tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n", tc.expectedErr, err)
		}
		if tc.expectedErr == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if tc.expectedConfig.numTimes != c.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.expectedConfig.numTimes, c.numTimes)
		}
	}

}
