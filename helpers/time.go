package helpers

import "time"

// Helper function to find the maximum of two time.Time values
func Max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// Helper function to find the minimum of two time.Time values
func Min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
