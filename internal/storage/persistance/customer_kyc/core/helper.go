package core

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func GenerateUsername(firstName, lastName, middleName string) string {
	// normalize: lowercase, trim spaces
	first := strings.ToLower(strings.TrimSpace(firstName))
	last := strings.ToLower(strings.TrimSpace(lastName))
	middle := strings.ToLower(strings.TrimSpace(middleName))

	// build base: first + middle initial (if any) + last
	base := first
	if len(middle) > 0 {
		base += string(middle[0]) // middle initial
	}
	base += last

	// strip non-alphanumeric characters (accents, symbols, spaces)
	reg := regexp.MustCompile(`[^a-z0-9]`)
	base = reg.ReplaceAllString(base, "")

	if base == "" {
		base = "user"
	}

	// append short random suffix to avoid collisions
	suffix := rand.Intn(9000) + 1000 // 1000–9999
	return fmt.Sprintf("%s%d", base, suffix)
}
