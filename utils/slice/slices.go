package slice_utils

// Contains checks if a slice contains a specified element.
func Contains[T comparable](slice []T, elem T) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate elements from a slice.
func RemoveDuplicates[T comparable](slice []T) []T {
	unique := make(map[T]bool)
	var result []T
	for _, v := range slice {
		if !unique[v] {
			unique[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Sum calculates the sum of integers in a slice.
func Sum(nums []int) int {
	total := 0
	for _, v := range nums {
		total += v
	}
	return total
}
