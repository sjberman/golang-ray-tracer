package utils

import "math"

// Mod properly handles modding with negative numbers (returns positive).
func Mod(x, y float64) float64 {
	res := math.Mod(x, y)
	if (res < 0 && y > 0) || (res > 0 && y < 0) {
		return res + y
	}

	return res
}

// Max returns the maximum value in a list of values.
func Max(vals ...float64) float64 {
	max := vals[0]
	for _, v := range vals {
		if v > max {
			max = v
		}
	}

	return max
}

// Min returns the minimum value in a list of values.
func Min(vals ...float64) float64 {
	min := vals[0]
	for _, v := range vals {
		if v < min {
			min = v
		}
	}

	return min
}
