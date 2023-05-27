package basic

import (
	"testing"
)

func TestFactorializeNumber(t *testing.T) {
	tests := []struct {
		input  int
		expect int64
	}{{
		input:  5,
		expect: 120,
	}, {
		input:  10,
		expect: 3628800,
	}, {
		input:  20,
		expect: 2432902008176640000,
	}, {
		input:  0,
		expect: 1,
	}}

	for _, test := range tests {
		fah := FactorializeNumber(test.input)
		if test.expect != fah {
			t.Errorf("Expect input to be %v, Got %v", test.expect, fah)
		}
	}
}
