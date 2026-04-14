package common

import "strings"

// Clamp ensures a value stays within a min/max range.
func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// IsHiddenFile returns true if the file name starts with a dot.
func IsHiddenFile(name string) bool {
	return strings.HasPrefix(name, ".")
}
