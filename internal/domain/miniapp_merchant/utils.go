package miniappmerchant


import "time"

// nonEmptyString returns s if non-empty, otherwise fallback
func nonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

// nonZeroTime returns t if non-zero, otherwise fallback
func nonZeroTime(t, fallback time.Time) time.Time {
	if !t.IsZero() {
		return t
	}
	return fallback
}

// nonZeroUint64 returns n if non-zero, otherwise fallback
func nonZeroUint64(n, fallback uint64) uint64 {
	if n != 0 {
		return n
	}
	return fallback
}
