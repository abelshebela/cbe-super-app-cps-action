package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateRandom returns a random numeric string of the specified digit length.
// If digit <= 0, it returns an empty string.
func GenerateRandom(digit int) string {
	if digit <= 0 {
		return ""
	}

	min := intPow(10, digit-1)
	max := intPow(10, digit) - 1
	if digit == 1 {
		min = 0
	}

	// Ensure the range is valid and non-negative for rand.Intn
	rangeSize := max - min + 1
	if rangeSize <= 0 {
		return ""
	}

	generatedNumber := min + rand.Intn(rangeSize)
	result := strconv.Itoa(generatedNumber)
	return result
}

func intPow(a, b int) int {
	result := 1
	for i := 0; i < b; i++ {
		result *= a
	}
	return result
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
