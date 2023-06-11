package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLargestNumbersInArray(t *testing.T) {
	tests := []struct {
		input  [][]int
		expect []int
	}{{
		input:  [][]int{{4, 5, 1, 8}, {10, 3, 6, 2}, {9, 7, 11, 5}, {12, 15, 13, 14}},
		expect: []int{8, 10, 11, 15},
	}, {
		input:  [][]int{{-4, -5, -1, -8}, {-10, -3, -6, -2}, {-9, -7, -11, -5}, {-12, -15, -13, -14}},
		expect: []int{-1, -2, -5, -12},
	}, {
		input:  [][]int{{-4, 5, -1, 8}, {-10, 3, -6, 2}, {-9, -7, 11, -5}, {12, -15, 13, -14}},
		expect: []int{8, 3, 11, 13},
	}, {
		input:  [][]int{{5, 5, 5, 5}, {10, 10, 10, 10}, {3, 3, 3, 3}, {15, 15, 15, 15}},
		expect: []int{5, 10, 3, 15},
	}, {
		input:  [][]int{{}, {}, {}, {}},
		expect: []int{},
	}, {
		input:  [][]int{{4, 5, 1, 8}},
		expect: []int{8},
	}}

	for _, test := range tests {
		fah := FindLargestNumbersInArray(test.input)
		assert.Equal(t, test.expect, fah)
	}
}
