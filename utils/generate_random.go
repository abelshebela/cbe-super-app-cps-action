package utils

import (
	"math/rand"
	"time"
)

func GenerateRandom(digits int) string {
	if digits <= 0 {
		return "0"
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, digits)
	for i := range digits {
		n := r.Intn(10) 
		if i == 0 {
			n = r.Intn(9) + 1
		}
		result[i] = '0' + byte(n)
	}

	return string(result)
}

