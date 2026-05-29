# Redis Storage Implementation

This directory contains a comprehensive Redis storage implementation for the CBE Super App Member Auth service. The implementation provides Redis-based storage for temporary data, session management, OTP operations, and rate limiting.

## Architecture

The Redis storage implementation follows a layered architecture:

```
Redis Storage Layer
├── RedisRepository (Base operations)
├── OTPSessionRepository (OTP management)
├── UserSessionRepository (Session management)
├── TemporaryDataRepository (Temporary data & caching)
├── RateLimitRepository (Rate limiting)
└── RedisStorageFactory (Factory pattern)
```

## Features

### 1. Base Redis Operations (`RedisRepository`)
- **Basic Operations**: Set, Get, Delete, Exists, Expire
- **Hash Operations**: HSet, HGet, HGetAll, HDel
- **List Operations**: LPush, RPush, LPop, RPop, LRange
- **Set Operations**: SAdd, SRem, SMembers, SIsMember
- **Key Management**: Keys, TTL, FlushDB

### 2. OTP Session Management (`OTPSessionRepository`)
- **OTP Storage**: Save, retrieve, and delete OTP codes
- **Attempt Tracking**: Increment and track OTP attempts
- **Verification Status**: Store and retrieve OTP verification status
- **Automatic Expiration**: TTL-based OTP expiration

### 3. User Session Management (`UserSessionRepository`)
- **Session Storage**: Save and retrieve user session data
- **Device Mapping**: Map users to their devices
- **Online Status**: Track user online/offline status
- **Session Updates**: Update specific session fields

### 4. Temporary Data Storage (`TemporaryDataRepository`)
- **Temporary Data**: Store temporary data with expiration
- **Form Data**: Save form data during multi-step processes
- **Caching**: General-purpose caching operations
- **Pattern-based Cleanup**: Clear cache by patterns

### 5. Rate Limiting (`RateLimitRepository`)
- **Request Counting**: Track request counts within time windows
- **IP-based Limiting**: Rate limiting by IP address
- **IP Blocking**: Block/unblock IP addresses
- **Limit Checking**: Check if requests exceed limits

## Configuration

### Environment Variables

Add the following Redis configuration to your environment:

```env
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
REDIS_DATABASE=0
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5
REDIS_MAX_RETRIES=3
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s
```

### Vault Configuration

The Redis configuration is integrated with Vault for secure configuration management:

```go
type VaultConfig struct {
    // ... existing fields ...
    RedisConfig RedisConfig
}
```

## Usage Examples

### 1. Basic Redis Operations

```go
// Initialize Redis storage
redisFactory := redis.NewRedisStorageFactory(client, logger)
redisFactory.Initialize()

// Get repositories
redisRepo := redisFactory.GetRedisRepository()

// Basic operations
err := redisRepo.Set(ctx, "key", "value", 1*time.Hour)
value, err := redisRepo.Get(ctx, "key")
exists, err := redisRepo.Exists(ctx, "key")
```

### 2. OTP Management

```go
otpRepo := redisFactory.GetOTPSessionRepository()

// Save OTP
err := otpRepo.SaveOTP(ctx, "+1234567890", "123456", 5*time.Minute)

// Get OTP
otp, err := otpRepo.GetOTP(ctx, "+1234567890")

// Track attempts
attempts, err := otpRepo.IncrementOTPAttempts(ctx, "+1234567890")

// Save verification status
err := otpRepo.SaveOTPVerification(ctx, "+1234567890", true, 10*time.Minute)
```

### 3. User Session Management

```go
sessionRepo := redisFactory.GetUserSessionRepository()

// Save user session
userData := map[string]interface{}{
    "user_id": "user123",
    "username": "john_doe",
    "last_login": time.Now().Unix(),
}
err := sessionRepo.SaveUserSession(ctx, "session123", userData, 24*time.Hour)

// Get session data
sessionData, err := sessionRepo.GetUserSession(ctx, "session123")

// Set user online
err := sessionRepo.SetUserOnline(ctx, "user123", "device456", 30*time.Minute)

// Check online status
isOnline, err := sessionRepo.IsUserOnline(ctx, "user123")
```

### 4. Temporary Data Storage

```go
tempRepo := redisFactory.GetTemporaryDataRepository()

// Save temporary data
data := map[string]interface{}{
    "step": "2",
    "form_data": "partial_form_data",
}
err := tempRepo.SaveTemporaryData(ctx, "temp_key", data, 1*time.Hour)

// Cache operations
err := tempRepo.SetCache(ctx, "user_profile:123", userProfile, 30*time.Minute)
cachedData, err := tempRepo.GetCache(ctx, "user_profile:123")
```

### 5. Rate Limiting

```go
rateLimitRepo := redisFactory.GetRateLimitRepository()

// Increment request count
count, err := rateLimitRepo.IncrementRequestCount(ctx, "api_key:123", 1*time.Minute)

// Check if rate limited
isLimited, err := rateLimitRepo.IsRateLimited(ctx, "api_key:123", 100)

// IP-based rate limiting
ipCount, err := rateLimitRepo.IncrementIPRequestCount(ctx, "192.168.1.1", 1*time.Minute)

// Block IP
err := rateLimitRepo.BlockIP(ctx, "192.168.1.1", 1*time.Hour)
```

## Redis Logger

The implementation includes a specialized Redis logger for monitoring and debugging:

```go
redisLogger := logger.NewRedisLogger(baseLogger)

// Log Redis operations
redisLogger.LogRedisOperation(ctx, "SET", "user:123", 5*time.Millisecond, nil)

// Log performance metrics
redisLogger.LogRedisPerformance(ctx, "GET", "cache:key", 10*time.Millisecond, 1024)

// Log cache hits/misses
redisLogger.LogRedisCacheHit(ctx, "user_profile:123", true, 2*time.Millisecond)

// Log rate limiting events
redisLogger.LogRedisRateLimit(ctx, "api_key:123", 100, 95, false)
```

## Key Naming Conventions

The implementation uses consistent key naming conventions:

- **OTP**: `otp:{phone_number}`
- **OTP Attempts**: `otp_attempts:{phone_number}`
- **OTP Verification**: `otp_verified:{phone_number}`
- **User Sessions**: `user_session:{session_id}`
- **User Devices**: `user_device:{user_id}`
- **Online Users**: `online_users` (set)
- **User Online Status**: `user_online:{user_id}`
- **Temporary Data**: `temp_data:{key}`
- **Form Data**: `form_data:{session_id}`
- **Cache**: `cache:{key}`
- **Rate Limits**: `rate_limit:{key}`
- **IP Rate Limits**: `ip_rate_limit:{ip}`
- **Blocked IPs**: `blocked_ip:{ip}`

## Error Handling

All Redis operations include comprehensive error handling:

- **Connection Errors**: Logged and returned
- **Key Not Found**: Handled gracefully with appropriate error messages
- **Serialization Errors**: JSON marshaling/unmarshaling errors are logged
- **Timeout Errors**: Redis operation timeouts are handled
- **Rate Limiting**: Exceeded limits are logged as warnings

## Performance Considerations

- **Connection Pooling**: Configurable connection pool size
- **Timeouts**: Configurable dial, read, and write timeouts
- **Retry Logic**: Automatic retry for failed operations
- **Expiration**: Automatic TTL for temporary data
- **Batch Operations**: Support for batch operations where applicable

## Monitoring and Health Checks

```go
// Health check
err := redisFactory.HealthCheck()

// Memory usage monitoring
redisLogger.LogRedisMemory(ctx, usedMemory, maxMemory, memoryUsage)

// Performance monitoring
redisLogger.LogRedisPerformance(ctx, operation, key, duration, size)
```

## Security Features

- **Password Authentication**: Redis password support
- **Database Isolation**: Configurable database selection
- **Key Expiration**: Automatic cleanup of temporary data
- **Rate Limiting**: Protection against abuse
- **IP Blocking**: Security against malicious IPs

## Testing

The Redis implementation can be tested using:

1. **Unit Tests**: Mock Redis client for testing
2. **Integration Tests**: Real Redis instance
3. **Performance Tests**: Load testing with Redis
4. **Health Checks**: Automated health monitoring

## Dependencies

- `github.com/redis/go-redis/v9`: Redis client
- `gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils`: Shared utilities
- `go.uber.org/zap`: Logging framework

## Future Enhancements

- **Redis Cluster Support**: Multi-node Redis cluster
- **Redis Sentinel**: High availability with Redis Sentinel
- **Redis Streams**: Real-time messaging capabilities
- **Redis Modules**: Custom Redis modules integration
- **Metrics Integration**: Prometheus metrics collection
- **Distributed Locking**: Redis-based distributed locks 