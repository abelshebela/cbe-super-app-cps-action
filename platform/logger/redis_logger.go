package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// RedisLogger provides specialized logging for Redis operations
type RedisLogger struct {
	logger Logger
}

// NewRedisLogger creates a new Redis logger instance
func NewRedisLogger(logger Logger) *RedisLogger {
	return &RedisLogger{
		logger: logger,
	}
}

// LogRedisOperation logs Redis operations with detailed information
func (r *RedisLogger) LogRedisOperation(ctx context.Context, operation string, key string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("key", key),
		zap.Duration("duration", duration),
		zap.Time("timestamp", time.Now()),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		r.logger.Error(ctx, fmt.Sprintf("Redis operation failed: %s", operation), fields...)
	} else {
		r.logger.Info(ctx, fmt.Sprintf("Redis operation completed: %s", operation), fields...)
	}
}

// LogRedisConnection logs Redis connection events
func (r *RedisLogger) LogRedisConnection(ctx context.Context, event string, details map[string]interface{}) {
	fields := []zap.Field{
		zap.String("event", event),
		zap.Time("timestamp", time.Now()),
	}

	for key, value := range details {
		fields = append(fields, zap.Any(key, value))
	}

	r.logger.Info(ctx, fmt.Sprintf("Redis connection event: %s", event), fields...)
}

// LogRedisPerformance logs Redis performance metrics
func (r *RedisLogger) LogRedisPerformance(ctx context.Context, operation string, key string, duration time.Duration, size int) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("key", key),
		zap.Duration("duration", duration),
		zap.Int("size_bytes", size),
		zap.Time("timestamp", time.Now()),
	}

	// Log as warning if operation takes too long
	if duration > 100*time.Millisecond {
		r.logger.Warn(ctx, fmt.Sprintf("Slow Redis operation: %s", operation), fields...)
	} else {
		r.logger.Info(ctx, fmt.Sprintf("Redis performance: %s", operation), fields...)
	}
}

// LogRedisError logs Redis-specific errors with context
func (r *RedisLogger) LogRedisError(ctx context.Context, operation string, key string, err error, context map[string]interface{}) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("key", key),
		zap.Error(err),
		zap.Time("timestamp", time.Now()),
	}

	for key, value := range context {
		fields = append(fields, zap.Any(key, value))
	}

	r.logger.Error(ctx, fmt.Sprintf("Redis error in operation: %s", operation), fields...)
}

// LogRedisCacheHit logs cache hit/miss events
func (r *RedisLogger) LogRedisCacheHit(ctx context.Context, key string, hit bool, duration time.Duration) {
	fields := []zap.Field{
		zap.String("key", key),
		zap.Bool("cache_hit", hit),
		zap.Duration("duration", duration),
		zap.Time("timestamp", time.Now()),
	}

	if hit {
		r.logger.Info(ctx, "Redis cache hit", fields...)
	} else {
		r.logger.Info(ctx, "Redis cache miss", fields...)
	}
}

// LogRedisRateLimit logs rate limiting events
func (r *RedisLogger) LogRedisRateLimit(ctx context.Context, key string, limit int, current int, blocked bool) {
	fields := []zap.Field{
		zap.String("key", key),
		zap.Int("limit", limit),
		zap.Int("current", current),
		zap.Bool("blocked", blocked),
		zap.Time("timestamp", time.Now()),
	}

	if blocked {
		r.logger.Warn(ctx, "Redis rate limit exceeded", fields...)
	} else {
		r.logger.Info(ctx, "Redis rate limit check", fields...)
	}
}

// LogRedisSession logs session-related Redis operations
func (r *RedisLogger) LogRedisSession(ctx context.Context, operation string, sessionID string, userID string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("session_id", sessionID),
		zap.String("user_id", userID),
		zap.Duration("duration", duration),
		zap.Time("timestamp", time.Now()),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		r.logger.Error(ctx, fmt.Sprintf("Redis session operation failed: %s", operation), fields...)
	} else {
		r.logger.Info(ctx, fmt.Sprintf("Redis session operation completed: %s", operation), fields...)
	}
}

// LogRedisOTP logs OTP-related Redis operations
func (r *RedisLogger) LogRedisOTP(ctx context.Context, operation string, phoneNumber string, attempts int, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("phone_number", phoneNumber),
		zap.Int("attempts", attempts),
		zap.Duration("duration", duration),
		zap.Time("timestamp", time.Now()),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		r.logger.Error(ctx, fmt.Sprintf("Redis OTP operation failed: %s", operation), fields...)
	} else {
		r.logger.Info(ctx, fmt.Sprintf("Redis OTP operation completed: %s", operation), fields...)
	}
}

// LogRedisHealth logs Redis health check results
func (r *RedisLogger) LogRedisHealth(ctx context.Context, healthy bool, latency time.Duration, details map[string]interface{}) {
	fields := []zap.Field{
		zap.Bool("healthy", healthy),
		zap.Duration("latency", latency),
		zap.Time("timestamp", time.Now()),
	}

	for key, value := range details {
		fields = append(fields, zap.Any(key, value))
	}

	if healthy {
		r.logger.Info(ctx, "Redis health check passed", fields...)
	} else {
		r.logger.Error(ctx, "Redis health check failed", fields...)
	}
}

// LogRedisMemory logs Redis memory usage
func (r *RedisLogger) LogRedisMemory(ctx context.Context, usedMemory int64, maxMemory int64, memoryUsage float64) {
	fields := []zap.Field{
		zap.Int64("used_memory", usedMemory),
		zap.Int64("max_memory", maxMemory),
		zap.Float64("memory_usage_percent", memoryUsage),
		zap.Time("timestamp", time.Now()),
	}

	if memoryUsage > 80 {
		r.logger.Warn(ctx, "Redis memory usage high", fields...)
	} else {
		r.logger.Info(ctx, "Redis memory usage normal", fields...)
	}
}
