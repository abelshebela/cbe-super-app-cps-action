package term_and_condition_oracle_core

import (
	"fmt"
	"time"
)

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ParseActivationTime parses an activation_time string (RFC3339 or YYYY-MM-DD) into time.Time.
func ParseActivationTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("activation_time is required")
	}
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported activation_time format: %q", s)
}
