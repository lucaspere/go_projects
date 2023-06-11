package basic

import (
	"testing"
)

func TestConfirmTheEnd(t *testing.T) {
	tests := []struct {
		str    string
		target string
		expect bool
	}{{
		str:    "Bastian",
		target: "n",
		expect: true,
	}, {
		str:    "lucas",
		target: "as",
		expect: true,
	}, {
		str:    "Connor",
		target: "n",
		expect: false,
	}, {
		str:    "Walking on water and developing software from a specification are easy if both are frozen",
		target: "specification",
		expect: false,
	}, {
		str:    "He has to give me a new name",
		target: "name",
		expect: true,
	}, {
		str:    "Open sesame",
		target: "same",
		expect: true,
	}, {
		str:    "Open sesame",
		target: "sage",
		expect: false,
	}, {
		str:    "Open sesame",
		target: "game",
		expect: false,
	}, {
		str:    "If you want to save our world, you must hurry. We dont know how much longer we can withstand the nothing",
		target: "mountain",
		expect: false,
	}, {
		str:    "Abstraction",
		target: "action",
		expect: true,
	}}

	for _, test := range tests {
		fah := ConfirmTheEnd(test.str, test.target)
		if test.expect != fah {
			t.Errorf("Expect the result of %s with target %s to be %v, Got %v", test.str, test.target, test.expect, fah)
		}
	}
}
