package basic

// Repeat a given string str for num times.
//
// Return an empty string if num is not a positive number
func RepeatStringNumTimes(str string, num int) (rs string) {
	if num > 0 {
		if num == 1 {
			rs = str
			return
		}

		for num != 0 {
			rs += str
			num--
		}
	}
	return
}
