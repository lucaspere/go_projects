package main

import (
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	var c config
	tests := []struct {
		args           []string
		expectedErr    error
		expectedConfig config
	}{
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

func TestValidateNumberArgs(t *testing.T) {
	tests := []struct {
		args     []string
		expected error
	}{
		{
			args:     []string{},
			expected: errors.New("must specify one argument"),
		},
		{
			args:     []string{"-n"},
			expected: nil,
		},
	}

	for _, tc := range tests {
		err := validateNumberArgs(tc.args)
		if err != nil && tc.expected.Error() != err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.expected, err)
		}
		if tc.expected == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
	}
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		args     []string
		expected error
	}{
		{
			args:     []string{"-h"},
			expected: errors.New("invalid -h argument. Must specify a valid argument -n"),
		},
		{
			args:     []string{"-n"},
			expected: errors.New("must insert some int value to -n argument"),
		},
		{
			args:     []string{"-n", "2"},
			expected: nil,
		},
	}

	for _, tc := range tests {
		err := validateArgs(tc.args)
		if err != nil && tc.expected.Error() != err.Error() {
			t.Errorf("expected error to be: %v, got: %v\n", tc.expected, err)
		}
		if len(tc.args) == 1 && tc.expected.Error() != err.Error() {
			t.Errorf("expected error to be: %v, got: %v\n", tc.expected, err)
		}
		if tc.expected == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
	}
}
