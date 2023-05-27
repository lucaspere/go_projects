package basic

import (
	"testing"
)

func TestConvertCtoF(t *testing.T) {
	tests := []struct {
		input  float64
		expect float64
	}{{
		input:  0,
		expect: 32,
	}, {
		input:  -30,
		expect: -22,
	}, {
		input:  -10,
		expect: 14,
	}, {
		input:  20,
		expect: 68,
	}, {
		input:  30,
		expect: 86,
	}}

	for _, test := range tests {
		fah := ConvertCtoF(test.input)
		if test.expect != fah {
			t.Errorf("Expect the Fahrenheit to be %f, Got %f", test.expect, fah)
		}
	}
}
