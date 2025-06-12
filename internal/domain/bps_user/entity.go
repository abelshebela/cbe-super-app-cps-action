package bpsuser

import (
	"time"
)

type BPSUser struct {
	ID                string
	UserCode          string
	FullName          string
	Username          string
	PhoneNumber       string
	BranchCode        []string
	BranchName        string
	HomeBranch        string
	Role              string
	Realm             string
	LoginAttemptCount uint8
	FirstPasswordSet  bool
	OTPVerifyCount    uint8
	OTPLastTriedAt    time.Time
	OTPLastVerifiedAt time.Time
	Password          Password
	Enabled           bool
	IsDeleted         bool
	LastLoginAttempt  time.Time
	NextLoginAttempt  time.Time
	LastLogin         time.Time
	CreatedAt         time.Time
	LastModifiedAt    time.Time
}

type Password struct {
	Salt             string
	CurrentPassword  string
	OldPassword      [4]string
	PasswordChangeAt time.Time
}
