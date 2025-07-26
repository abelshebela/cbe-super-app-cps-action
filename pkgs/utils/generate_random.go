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
	generatedNumber := min + rand.Intn(max-min+1)
	result := strconv.Itoa(generatedNumber)
	return result
}

func GenerateUniqueActionCode(length int) string {
	if length <= 0 {
		return ""
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
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
