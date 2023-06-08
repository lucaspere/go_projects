package basic

import (
	"testing"
)

func TestFindLongestWordLength(t *testing.T) {
	tests := []struct {
		input  string
		expect int
	}{{
		input:  "The quick brown fox jumped over the lazy dog",
		expect: 6,
	}, {
		input:  "May the force be with you",
		expect: 5,
	}, {
		input:  "Google do a barrel roll",
		expect: 6,
	}, {
		input:  "What is the average airspeed velocity of an unladen swallow",
		expect: 8,
	}, {
		input:  "What if we try a super-long word such as otorhinolaryngology",
		expect: 19,
	}}

	for _, test := range tests {
		fah := FindLongestWordLength(test.input)
		if test.expect != fah {
			t.Errorf("Expect the number to be %v, Got %v", test.expect, fah)
		}
	}
}
