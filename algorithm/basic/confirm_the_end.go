package basic

import (
	"strings"
)

// Check if a string (first argument, str) ends with the given target string (second argument, target).
func ConfirmTheEnd(str string, target string) (isEnding bool) {
	strArr := strings.Split(str, "")[len(str)-len(target):]
	isEnding = strings.EqualFold(strings.Join(strArr, ""), target)

	return
}
