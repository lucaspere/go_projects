package basic

import (
	"testing"
)

func TestRepeatStringNumTimes(t *testing.T) {
	tests := []struct {
		str    string
		num    int
		expect string
	}{{
		str:    "*",
		num:    3,
		expect: "***",
	}, {
		str:    "abc",
		num:    3,
		expect: "abcabcabc",
	}, {
		str:    "abc",
		num:    4,
		expect: "abcabcabcabc",
	}, {
		str:    "abc",
		num:    1,
		expect: "abc",
	}, {
		str:    "*",
		num:    8,
		expect: "********",
	}, {
		str:    "abc",
		num:    -2,
		expect: "",
	}, {
		str:    "abc",
		num:    0,
		expect: "",
	}}

	for _, test := range tests {
		fah := RepeatStringNumTimes(test.str, test.num)
		if test.expect != fah {
			t.Errorf("Expect the result of %s with num %d to be %v, Got %v", test.str, test.num, test.expect, fah)
		}
	}
}
