package notification

// nonEmptyString returns s if non-empty, otherwise fallback
func nonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func nonEmptyNotificationFor(newFor string, oldFor NotificationFor) NotificationFor {
	if newFor != "" {
		return NotificationFor(newFor)
	}
	return oldFor
}
