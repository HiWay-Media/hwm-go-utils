package int_utils

// Min returns the smallest integer from a list of integers.
func Min(nums ...int) int {
	if len(nums) == 0 {
		return 0 // or some error value
	}
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}

// Max returns the largest integer from a list of integers.
func Max(nums ...int) int {
	if len(nums) == 0 {
		return 0 // or some error value
	}
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

// IsEven returns true if the integer is even.
func IsEven(n int) bool {
	return n%2 == 0
}

// IsOdd returns true if the integer is odd.
func IsOdd(n int) bool {
	return n%2 != 0
}
