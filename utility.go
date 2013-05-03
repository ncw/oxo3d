package main

// Returns the minimum of a and b
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// Returns the maximum of a and b
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
