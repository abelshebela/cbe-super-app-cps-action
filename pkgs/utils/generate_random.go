package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandom(digit int) string {
	fmt.Println("generating number...")

	numberDigit := "1"
	multiplier := "9"

	for i := 1; i < digit; i++ {
		numberDigit += "0"
		multiplier += "9"
	}

	fmt.Println(numberDigit, multiplier)

	min, _ := strconv.Atoi(numberDigit)
	max, _ := strconv.Atoi(multiplier)
	rand.Seed(time.Now().UnixNano())
	generatedNumber := min + rand.Intn(max+1)

	result := strconv.Itoa(generatedNumber)
	if len(result) > digit {
		result = result[:digit]
	}
	return result
}
