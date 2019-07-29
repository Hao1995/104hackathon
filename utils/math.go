package utils

func CountDigits(i int64) int {
	var count int
	for i != 0 {
		i /= 10
		count++
	}
	return count
}
