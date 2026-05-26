package core

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

func GenerateUsername(firstName, lastName, middleName string) string {
	first := strings.ToLower(strings.TrimSpace(firstName))
	last := strings.ToLower(strings.TrimSpace(lastName))
	middle := strings.ToLower(strings.TrimSpace(middleName))

	base := first
	if len(middle) > 0 {
		base += string(middle[0])
	}
	base += last

	reg := regexp.MustCompile(`[^a-z0-9]`)
	base = reg.ReplaceAllString(base, "")

	if base == "" {
		base = "user"
	}

	// 6-digit suffix via crypto/rand — 900000 possible values reduces collision risk
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	var suffix int64
	if err != nil {
		suffix = 100000
	} else {
		suffix = n.Int64() + 100000
	}
	return fmt.Sprintf("%s%d", base, suffix)
}
