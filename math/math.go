package math

// Min get the minimum of two values
func Min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// Max get the maximum of two values
func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}
