package basic

func ReverseString(str string) (rs string) {
	st := make([]rune, len(str))
	for idx, c := range str {
		st[len(str)-1-idx] = c
	}

	rs = string(st)
	return
}
