package utils

import (
	"context"
	"time"
)

const defaultSessionMinute = 15

type User struct {
	ID               string
	SessionExpiresOn time.Time
	LastModified     time.Time
}

type UserDAL interface {
	UpdateUserSession(ctx context.Context, userID string, sessionExpiresOn, lastModified time.Time) error
}

func Session(ctx context.Context, routePath string, user *User, dal UserDAL) error {
	if !contains(routePath, "extendsession") {
		if time.Now().UTC().Add(1 * time.Minute).After(user.SessionExpiresOn) {
			updateSessionExpiresOn := time.Now().Add(defaultSessionMinute * time.Minute)
			lastModified := time.Now()
			err := dal.UpdateUserSession(ctx, user.ID, updateSessionExpiresOn, lastModified)
			if err != nil {
				return err
			}
			user.SessionExpiresOn = updateSessionExpiresOn
			user.LastModified = lastModified
		}
	}
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
