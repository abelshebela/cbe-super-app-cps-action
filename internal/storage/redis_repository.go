package storage

import (
	"context"
	"time"
)

// RedisRepository defines the main Redis operations interface
type RedisRepository interface {
	// Basic operations
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// Hash operations
	HSet(ctx context.Context, key string, field string, value interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// List operations
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	// Set operations
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)

	// Key management
	Keys(ctx context.Context, pattern string) ([]string, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	FlushDB(ctx context.Context) error
}

// OTPSessionRepository defines Redis operations for OTP session management
type OTPSessionRepository interface {
	// OTP operations
	SaveOTP(ctx context.Context, phoneNumber string, otp string, expiration time.Duration) error
	GetOTP(ctx context.Context, phoneNumber string) (string, error)
	DeleteOTP(ctx context.Context, phoneNumber string) error
	IncrementOTPAttempts(ctx context.Context, phoneNumber string) (int, error)
	GetOTPAttempts(ctx context.Context, phoneNumber string) (int, error)
	ResetOTPAttempts(ctx context.Context, phoneNumber string) error

	// OTP verification
	SaveOTPVerification(ctx context.Context, phoneNumber string, verified bool, expiration time.Duration) error
	GetOTPVerification(ctx context.Context, phoneNumber string) (bool, error)
	DeleteOTPVerification(ctx context.Context, phoneNumber string) error
}

// UserSessionRepository defines Redis operations for user session management
type UserSessionRepository interface {
	// User session operations
	SaveUserSession(ctx context.Context, sessionID string, userData map[string]interface{}, expiration time.Duration) error
	GetUserSession(ctx context.Context, sessionID string) (map[string]string, error)
	DeleteUserSession(ctx context.Context, sessionID string) error
	UpdateUserSession(ctx context.Context, sessionID string, updates map[string]interface{}) error

	// User device mapping
	SaveUserDeviceMapping(ctx context.Context, userID string, deviceUUID string, sessionData map[string]interface{}, expiration time.Duration) error
	GetUserDeviceMapping(ctx context.Context, userID string) (map[string]string, error)
	DeleteUserDeviceMapping(ctx context.Context, userID string) error

	// User online status
	SetUserOnline(ctx context.Context, userID string, deviceUUID string, expiration time.Duration) error
	SetUserOffline(ctx context.Context, userID string, deviceUUID string) error
	IsUserOnline(ctx context.Context, userID string) (bool, error)
	GetOnlineUsers(ctx context.Context) ([]string, error)
}

// TemporaryDataRepository defines Redis operations for temporary data storage
type TemporaryDataRepository interface {
	// Temporary data operations
	SaveTemporaryData(ctx context.Context, key string, data map[string]interface{}, expiration time.Duration) error
	GetTemporaryData(ctx context.Context, key string) (map[string]string, error)
	DeleteTemporaryData(ctx context.Context, key string) error
	UpdateTemporaryData(ctx context.Context, key string, updates map[string]interface{}) error

	// Form data storage
	SaveFormData(ctx context.Context, sessionID string, formData map[string]interface{}, expiration time.Duration) error
	GetFormData(ctx context.Context, sessionID string) (map[string]string, error)
	DeleteFormData(ctx context.Context, sessionID string) error

	// Cache operations
	SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetCache(ctx context.Context, key string) (string, error)
	DeleteCache(ctx context.Context, key string) error
	ClearCacheByPattern(ctx context.Context, pattern string) error
}

// RateLimitRepository defines Redis operations for rate limiting
type RateLimitRepository interface {
	// Rate limiting operations
	IncrementRequestCount(ctx context.Context, key string, window time.Duration) (int, error)
	GetRequestCount(ctx context.Context, key string) (int, error)
	ResetRequestCount(ctx context.Context, key string) error
	IsRateLimited(ctx context.Context, key string, limit int) (bool, error)

	// IP-based rate limiting
	IncrementIPRequestCount(ctx context.Context, ip string, window time.Duration) (int, error)
	GetIPRequestCount(ctx context.Context, ip string) (int, error)
	BlockIP(ctx context.Context, ip string, duration time.Duration) error
	IsIPBlocked(ctx context.Context, ip string) (bool, error)
	UnblockIP(ctx context.Context, ip string) error
}
