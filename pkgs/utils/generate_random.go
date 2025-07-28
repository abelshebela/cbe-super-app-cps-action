package utils

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	cRand "crypto/rand"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

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

func GeneratePrefixedName(prefix, value string, logger shared_utils.Logger) (string, error) {
	logger.Infof("Generating prefixed name", "prefix", prefix, "value", value)

	if prefix == "" || value == "" {
		logger.Errorf("Invalid input for GeneratePrefixedName", "prefix", prefix, "value", value)
		return "", fmt.Errorf("prefix and value must not be empty")
	}

	digits := "0123456789"
	max := big.NewInt(int64(len(digits)))
	code := make([]byte, 7)

	for i := range code {
		n, err := cRand.Int(cRand.Reader, max)
		if err != nil {
			logger.Errorf("Failed to generate random digit", "error", err)
			return "", fmt.Errorf("failed to generate random digit: %v", err)
		}
		code[i] = digits[n.Int64()]
	}

	value = strings.ReplaceAll(value, " ", "")
	prefix = strings.ReplaceAll(prefix, " ", "")
	result := strings.Join([]string{prefix, value, string(code)}, "-")
	logger.Infof("Successfully generated prefixed name", "result", result)
	return result, nil
}
