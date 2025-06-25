package utils

import (
	"fmt"
)

func UpdatePhoneNumber(phoneNumber interface{}) string {
	var strPhone string
	switch v := phoneNumber.(type) {
	case int, int64, float64:
		strPhone = fmt.Sprintf("%v", v)
	case string:
		strPhone = v
	default:
		return ""
	}

	if len(strPhone) == 0 {
		return ""
	}

	firstDigit := strPhone[0]
	switch firstDigit {
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
