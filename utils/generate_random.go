package utils

import (
	"math/rand"
	"time"
)

func GenerateRandom(digit int) string {
	if digit <= 0 {
		return "0"
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, digit)
	for i := 0; i < digit; i++ {
		if i == 0 {
			result[i] = byte(r.Intn(9)+1) + '0' // avoid leading 0
		} else {
			result[i] = byte(r.Intn(10)) + '0'
		}
	}
	return string(result)
}
