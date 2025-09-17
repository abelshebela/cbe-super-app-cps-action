package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string `mapstructure:"REDIS_HOST"`
	Port         string `mapstructure:"REDIS_PORT"`
	Password     string `mapstructure:"REDIS_PASSWORD"`
	Database     int    `mapstructure:"REDIS_DATABASE"`
	PoolSize     int    `mapstructure:"REDIS_POOL_SIZE"`
	MinIdleConns int    `mapstructure:"REDIS_MIN_IDLE_CONNS"`
	MaxRetries   int    `mapstructure:"REDIS_MAX_RETRIES"`
	DialTimeout  string `mapstructure:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout  string `mapstructure:"REDIS_READ_TIMEOUT"`
	WriteTimeout string `mapstructure:"REDIS_WRITE_TIMEOUT"`
}

// ConnectRedis establishes connection to Redis
func ConnectRedis(logger interface{}, env *VaultConfig) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Parse timeouts
	dialTimeout, err := time.ParseDuration(env.RedisConfig.DialTimeout)
	if err != nil {
		dialTimeout = 5 * time.Second
	}

	readTimeout, err := time.ParseDuration(env.RedisConfig.ReadTimeout)
	if err != nil {
		readTimeout = 3 * time.Second
	}

	writeTimeout, err := time.ParseDuration(env.RedisConfig.WriteTimeout)
	if err != nil {
		writeTimeout = 3 * time.Second
	}

	// Set default values if not provided
	poolSize := env.RedisConfig.PoolSize
	if poolSize == 0 {
		poolSize = 10
	}

	minIdleConns := env.RedisConfig.MinIdleConns
	if minIdleConns == 0 {
		minIdleConns = 5
	}

	maxRetries := env.RedisConfig.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3
	}

	redisAddr := env.RedisConfig.Host + ":" + env.RedisConfig.Port
	if env.RedisConfig.Port == "" {
		redisAddr = env.RedisConfig.Host + ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     env.RedisConfig.Password,
		DB:           env.RedisConfig.Database,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		MaxRetries:   maxRetries,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	// Type assertion for logger
	if log, ok := logger.(interface {
		Info(context.Context, string, ...zap.Field)
	}); ok {
		log.Info(ctx, "Successfully Connected to Redis!")
	}

	return client, nil
}

// DisconnectRedis gracefully disconnects from Redis
func DisconnectRedis(ctx context.Context, client *redis.Client, logger interface{}) {
	if err := client.Close(); err != nil {
		if log, ok := logger.(interface {
			Error(context.Context, string, ...zap.Field)
		}); ok {
			log.Error(ctx, "Error disconnecting Redis", zap.Error(err))
		}
	} else {
		if log, ok := logger.(interface {
			Info(context.Context, string, ...zap.Field)
		}); ok {
			log.Info(ctx, "Disconnected Redis successfully")
		}
	}
}
