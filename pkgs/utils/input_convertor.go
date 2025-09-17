package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// UpdatePhoneNumber normalizes a phone number to the +251 format if possible.
// Accepts int, int64, float64, or string. Returns an empty string for unsupported types or empty input.
func UpdatePhoneNumber(phoneNumber interface{}) string {
	var strPhone string
	switch v := phoneNumber.(type) {
	case int:
		strPhone = strconv.Itoa(v)
	case int64:
		strPhone = strconv.FormatInt(v, 10)
	case float64:
		strPhone = strconv.FormatInt(int64(v), 10)
	case string:
		strPhone = v
	default:
		return ""
	}

	strPhone = trimSpaces(strPhone)
	if len(strPhone) == 0 {
		return ""
	}

	switch strPhone[0] {
	case '0':
		return "+251" + strPhone[1:]
	case '+':
		return strPhone
	case '9', '7':
		return "+251" + strPhone
	default:
		return strPhone
	}
}

func trimSpaces(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func CleanAmountString(amountStr string) string {
	// Remove leading/trailing whitespace
	amountStr = strings.TrimSpace(amountStr)

	// Check if this looks like European format (comma as decimal separator)
	// European format: "1.000,50" -> 1000.50
	// US format: "1,000.50" -> 1000.50
	isEuropeanFormat := false
	if strings.Contains(amountStr, ",") && strings.Contains(amountStr, ".") {
		// If both comma and dot exist, check the pattern
		// European: "1.000,50" (dot before comma)
		// US: "1,000.50" (comma before dot)
		commaIndex := strings.Index(amountStr, ",")
		dotIndex := strings.Index(amountStr, ".")
		if dotIndex < commaIndex {
			isEuropeanFormat = true
		}
	} else if strings.Contains(amountStr, ",") && !strings.Contains(amountStr, ".") {
		// If only comma exists, it might be European format
		// Check if there are digits after the comma
		parts := strings.Split(amountStr, ",")
		if len(parts) == 2 && len(parts[1]) <= 2 {
			// Likely European format (e.g., "1000,50")
			isEuropeanFormat = true
		}
	}

	// Remove currency symbols and other non-numeric characters except digits, dots, and minus
	// This handles cases like "$1,000.50", "€1.000,50", "1,000 USD", etc.
	var cleaned strings.Builder
	decimalFound := false
	minusFound := false

	for _, char := range amountStr {
		switch {
		case char >= '0' && char <= '9':
			cleaned.WriteRune(char)
		case char == '.' && !decimalFound:
			if isEuropeanFormat {
				// In European format, dot is thousands separator, skip it
				continue
			}
			cleaned.WriteRune(char)
			decimalFound = true
		case char == ',' && !decimalFound:
			if isEuropeanFormat {
				// In European format, comma is decimal separator
				cleaned.WriteRune('.')
				decimalFound = true
			} else {
				// Skip thousands separators before decimal point
				continue
			}
		case char == ',' && decimalFound:
			// Skip thousands separators after decimal point
			continue
		case char == '-' && !minusFound && cleaned.Len() == 0:
			// Allow minus sign only at the beginning
			cleaned.WriteRune(char)
			minusFound = true
		}
	}

	result := cleaned.String()

	// Handle edge cases
	if result == "" || result == "-" {
		return "0"
	}

	// Remove trailing decimal point
	if strings.HasSuffix(result, ".") {
		result = strings.TrimSuffix(result, ".")
	}

	return result
}

func ParseDateString(dateStr string) (time.Time, error) {
	formats := []string{
		"02/01/2006",                // DD/MM/YYYY
		"01/02/2006",                // MM/DD/YYYY
		"2006-01-02",                // YYYY-MM-DD
		"2006-01-02T15:04:05Z07:00", // ISO format
		"2006-01-02T15:04:05",       // ISO format without timezone
		"02-01-2006",                // DD-MM-YYYY
		"01-02-2006",                // MM-DD-YYYY
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
