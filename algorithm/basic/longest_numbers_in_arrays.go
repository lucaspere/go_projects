package basic

// Return an array consisting of the largest number from each provided sub-array.
func FindLargestNumbersInArray(arrNumbers [][]int) (largests []int) {
	largests = make([]int, 0, len(arrNumbers))
	for _, arr := range arrNumbers {
		if len(arr) > 0 {
			largest := arr[0]

			for i := 1; i < len(arr); i++ {
				if largest < arr[i] {
					largest = arr[i]
				}
			}

			largests = append(largests, largest)
		}
	}

	return
}
