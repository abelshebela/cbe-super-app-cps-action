package utils

import (
	"fmt"
	"time"
)

func ParseTime(date string) time.Time {
	parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}
	}
	return parsedTime
}
