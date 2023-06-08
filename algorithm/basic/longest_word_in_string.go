package basic

import "strings"

// Return the length of the longest word in the provided sentence.
func FindLongestWordLength(str string) (num int) {
	strArray := strings.Split(str, " ")
	num = len(strArray[0])

	for i := 1; i < len(strArray); i++ {
		if num < len(strArray[i]) {
			num = len(strArray[i])
		}
	}

	return
}
