package basic

import (
	"testing"
)

func TestReverseString(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{{
		input:  "hello",
		expect: "olleh",
	}, {
		input:  "Howdy",
		expect: "ydwoH",
	}, {
		input:  "Greetings from Earth",
		expect: "htraE morf sgniteerG",
	}}

	for _, test := range tests {
		fah := ReverseString(test.input)
		if test.expect != fah {
			t.Errorf("Expect input to be %v, Got %v", test.expect, fah)
		}
	}
}
