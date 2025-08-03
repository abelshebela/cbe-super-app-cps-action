package redis

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"github.com/redis/go-redis/v9"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// RedisStorageFactory manages all Redis repositories
type RedisStorageFactory struct {
	client            *redis.Client
	logger            utils.Logger
	redisRepo         storage.RedisRepository
	otpSessionRepo    storage.OTPSessionRepository
	userSessionRepo   storage.UserSessionRepository
	temporaryDataRepo storage.TemporaryDataRepository
	rateLimitRepo     storage.RateLimitRepository
}

// NewRedisStorageFactory creates a new Redis storage factory
func NewRedisStorageFactory(client *redis.Client, logger utils.Logger) *RedisStorageFactory {
	return &RedisStorageFactory{
		client: client,
		logger: logger,
	}
}

// Initialize initializes all Redis repositories
func (f *RedisStorageFactory) Initialize() {
	// Initialize base Redis repository
	f.redisRepo = NewRedisRepository(f.client, f.logger)

	// Initialize specialized repositories
	f.otpSessionRepo = NewOTPSessionRepository(f.redisRepo, f.logger)
	f.userSessionRepo = NewUserSessionRepository(f.redisRepo, f.logger)
	f.temporaryDataRepo = NewTemporaryDataRepository(f.redisRepo, f.logger)
	f.rateLimitRepo = NewRateLimitRepository(f.redisRepo, f.logger)

	f.logger.Infof("Redis storage factory initialized successfully")
}

// GetRedisRepository returns the base Redis repository
func (f *RedisStorageFactory) GetRedisRepository() storage.RedisRepository {
	return f.redisRepo
}

// GetOTPSessionRepository returns the OTP session repository
func (f *RedisStorageFactory) GetOTPSessionRepository() storage.OTPSessionRepository {
	return f.otpSessionRepo
}

// GetUserSessionRepository returns the user session repository
func (f *RedisStorageFactory) GetUserSessionRepository() storage.UserSessionRepository {
	return f.userSessionRepo
}

// GetTemporaryDataRepository returns the temporary data repository
func (f *RedisStorageFactory) GetTemporaryDataRepository() storage.TemporaryDataRepository {
	return f.temporaryDataRepo
}

// GetRateLimitRepository returns the rate limit repository
func (f *RedisStorageFactory) GetRateLimitRepository() storage.RateLimitRepository {
	return f.rateLimitRepo
}

// GetAllRepositories returns all Redis repositories
func (f *RedisStorageFactory) GetAllRepositories() (
	storage.RedisRepository,
	storage.OTPSessionRepository,
	storage.UserSessionRepository,
	storage.TemporaryDataRepository,
	storage.RateLimitRepository,
) {
	return f.redisRepo, f.otpSessionRepo, f.userSessionRepo, f.temporaryDataRepo, f.rateLimitRepo
}

// Close closes the Redis client connection
func (f *RedisStorageFactory) Close() error {
	if f.client != nil {
		err := f.client.Close()
		if err != nil {
			f.logger.Errorf("Failed to close Redis client: %v", err)
			return err
		}
		f.logger.Infof("Redis client closed successfully")
	}
	return nil
}

// HealthCheck performs a health check on Redis
func (f *RedisStorageFactory) HealthCheck() error {
	ctx := f.client.Context()
	err := f.client.Ping(ctx).Err()
	if err != nil {
		f.logger.Errorf("Redis health check failed: %v", err)
		return err
	}
	f.logger.Infof("Redis health check passed")
	return nil
}
