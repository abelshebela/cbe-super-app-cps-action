package utils

import (
	"strconv"
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
