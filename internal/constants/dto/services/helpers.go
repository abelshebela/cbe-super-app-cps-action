package services

func StringPointer(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return *s
}

func BoolPointer(b *bool, fallback bool) bool {
	if b == nil {
		return fallback
	}
	return *b
}

func Float64Pointer(f *float64, fallback float64) float64 {
	if f == nil {
		return fallback
	}
	return *f
}
