package array

import "strings"

// Contains sequentially searches haystack for needle
func Contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

// DeepContains sequentially searches haystack for needle
func DeepContains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if strings.Contains(item, needle) {
			return true
		}
	}

	return false
}
