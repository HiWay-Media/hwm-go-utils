package conv_utils

import (
	"strconv"
)

// StrToInt converts a string to an integer.
func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// IntToStr converts an integer to a string.
func IntToStr(i int) string {
	return strconv.Itoa(i)
}

// IntToFloat converts an integer to a float64.
func IntToFloat(i int) float64 {
	return float64(i)
}

// FloatToInt converts a float64 to an integer.
func FloatToInt(f float64) int {
	return int(f)
}