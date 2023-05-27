package basic

// Return the factorial of the provided integer.
//
// If the integer is represented with the letter _n_, a factorial is the product of all positive  integers less than or equal to _n_.
//
// Factorials are often represented with the shorthand notation _!n_
//
// For example: 5! = 1 * 2 * 3 * 4 * 5 = 120
//
// only integers greater than or equal to zero will be supplie to the function.
func FactorializeNumber(num int) (fac int64) {
	if num == 1 || num == 0 {
		fac = int64(1)
	} else {
		fac = int64(num) * FactorializeNumber(num-1)
	}

	return
}
