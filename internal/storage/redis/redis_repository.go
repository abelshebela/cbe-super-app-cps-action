package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/storage"

	"github.com/redis/go-redis/v9"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// RedisRepository implements the RedisRepository interface
type RedisRepository struct {
	client *redis.Client
	logger utils.Logger
}

// NewRedisRepository creates a new Redis repository instance
func NewRedisRepository(client *redis.Client, logger utils.Logger) storage.RedisRepository {
	return &RedisRepository{
		client: client,
		logger: logger,
	}
}

// Set sets a key-value pair with optional expiration
func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var stringValue string

	switch v := value.(type) {
	case string:
		stringValue = v
	default:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			r.logger.Errorf("Failed to marshal value for key %s: %v", key, err)
			return err
		}
		stringValue = string(jsonValue)
	}

	err := r.client.Set(ctx, key, stringValue, expiration).Err()
	if err != nil {
		r.logger.Errorf("Failed to set key %s: %v", key, err)
		return err
	}

	r.logger.Infof("Successfully set key: %s", key)
	return nil
}

// Get retrieves a value by key
func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Warnf("Key not found: %s", key)
			return "", fmt.Errorf("key not found: %s", key)
		}
		r.logger.Errorf("Failed to get key %s: %v", key, err)
		return "", err
	}

	r.logger.Infof("Successfully retrieved key: %s", key)
	return value, nil
}

// Delete removes a key
func (r *RedisRepository) Delete(ctx context.Context, key string) error {
	result, err := r.client.Del(ctx, key).Result()
	if err != nil {
		r.logger.Errorf("Failed to delete key %s: %v", key, err)
		return err
	}

	if result == 0 {
		r.logger.Warnf("Key not found for deletion: %s", key)
		return fmt.Errorf("key not found: %s", key)
	}

	r.logger.Infof("Successfully deleted key: %s", key)
	return nil
}

// Exists checks if a key exists
func (r *RedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Errorf("Failed to check existence of key %s: %v", key, err)
		return false, err
	}

	exists := result > 0
	r.logger.Infof("Key %s exists: %t", key, exists)
	return exists, nil
}

// Expire sets expiration time for a key
func (r *RedisRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	result, err := r.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		r.logger.Errorf("Failed to set expiration for key %s: %v", key, err)
		return err
	}

	if !result {
		r.logger.Warnf("Key not found for expiration setting: %s", key)
		return fmt.Errorf("key not found: %s", key)
	}

	r.logger.Infof("Successfully set expiration for key: %s", key)
	return nil
}

// HSet sets a field in a hash
func (r *RedisRepository) HSet(ctx context.Context, key string, field string, value interface{}) error {
	var stringValue string

	switch v := value.(type) {
	case string:
		stringValue = v
	default:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			r.logger.Errorf("Failed to marshal value for hash field %s:%s: %v", key, field, err)
			return err
		}
		stringValue = string(jsonValue)
	}

	err := r.client.HSet(ctx, key, field, stringValue).Err()
	if err != nil {
		r.logger.Errorf("Failed to set hash field %s:%s: %v", key, field, err)
		return err
	}

	r.logger.Infof("Successfully set hash field: %s:%s", key, field)
	return nil
}

// HGet retrieves a field from a hash
func (r *RedisRepository) HGet(ctx context.Context, key string, field string) (string, error) {
	value, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Warnf("Hash field not found: %s:%s", key, field)
			return "", fmt.Errorf("hash field not found: %s:%s", key, field)
		}
		r.logger.Errorf("Failed to get hash field %s:%s: %v", key, field, err)
		return "", err
	}

	r.logger.Infof("Successfully retrieved hash field: %s:%s", key, field)
	return value, nil
}

// HGetAll retrieves all fields from a hash
func (r *RedisRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	values, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		r.logger.Errorf("Failed to get all hash fields for key %s: %v", key, err)
		return nil, err
	}

	if len(values) == 0 {
		r.logger.Warnf("Hash not found or empty: %s", key)
		return nil, fmt.Errorf("hash not found: %s", key)
	}

	r.logger.Infof("Successfully retrieved all hash fields for key: %s", key)
	return values, nil
}

// HDel removes fields from a hash
func (r *RedisRepository) HDel(ctx context.Context, key string, fields ...string) error {
	result, err := r.client.HDel(ctx, key, fields...).Result()
	if err != nil {
		r.logger.Errorf("Failed to delete hash fields %s:%v: %v", key, fields, err)
		return err
	}

	if result == 0 {
		r.logger.Warnf("Hash fields not found for deletion: %s:%v", key, fields)
		return fmt.Errorf("hash fields not found: %s:%v", key, fields)
	}

	r.logger.Infof("Successfully deleted hash fields: %s:%v", key, fields)
	return nil
}

// LPush pushes values to the left of a list
func (r *RedisRepository) LPush(ctx context.Context, key string, values ...interface{}) error {
	var stringValues []interface{}

	for _, value := range values {
		switch v := value.(type) {
		case string:
			stringValues = append(stringValues, v)
		default:
			jsonValue, err := json.Marshal(value)
			if err != nil {
				r.logger.Errorf("Failed to marshal value for list push %s: %v", key, err)
				return err
			}
			stringValues = append(stringValues, string(jsonValue))
		}
	}

	err := r.client.LPush(ctx, key, stringValues...).Err()
	if err != nil {
		r.logger.Errorf("Failed to push to list %s: %v", key, err)
		return err
	}

	r.logger.Infof("Successfully pushed to list: %s", key)
	return nil
}

// RPush pushes values to the right of a list
func (r *RedisRepository) RPush(ctx context.Context, key string, values ...interface{}) error {
	var stringValues []interface{}

	for _, value := range values {
		switch v := value.(type) {
		case string:
			stringValues = append(stringValues, v)
		default:
			jsonValue, err := json.Marshal(value)
			if err != nil {
				r.logger.Errorf("Failed to marshal value for list push %s: %v", key, err)
				return err
			}
			stringValues = append(stringValues, string(jsonValue))
		}
	}

	err := r.client.RPush(ctx, key, stringValues...).Err()
	if err != nil {
		r.logger.Errorf("Failed to push to list %s: %v", key, err)
		return err
	}

	r.logger.Infof("Successfully pushed to list: %s", key)
	return nil
}

// LPop pops a value from the left of a list
func (r *RedisRepository) LPop(ctx context.Context, key string) (string, error) {
	value, err := r.client.LPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Warnf("List is empty: %s", key)
			return "", fmt.Errorf("list is empty: %s", key)
		}
		r.logger.Errorf("Failed to pop from list %s: %v", key, err)
		return "", err
	}

	r.logger.Infof("Successfully popped from list: %s", key)
	return value, nil
}

// RPop pops a value from the right of a list
func (r *RedisRepository) RPop(ctx context.Context, key string) (string, error) {
	value, err := r.client.RPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Warnf("List is empty: %s", key)
			return "", fmt.Errorf("list is empty: %s", key)
		}
		r.logger.Errorf("Failed to pop from list %s: %v", key, err)
		return "", err
	}

	r.logger.Infof("Successfully popped from list: %s", key)
	return value, nil
}

// LRange gets a range of elements from a list
func (r *RedisRepository) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	values, err := r.client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		r.logger.Errorf("Failed to get range from list %s: %v", key, err)
		return nil, err
	}

	r.logger.Infof("Successfully retrieved range from list: %s", key)
	return values, nil
}

// SAdd adds members to a set
func (r *RedisRepository) SAdd(ctx context.Context, key string, members ...interface{}) error {
	var stringMembers []interface{}

	for _, member := range members {
		switch v := member.(type) {
		case string:
			stringMembers = append(stringMembers, v)
		default:
			jsonValue, err := json.Marshal(member)
			if err != nil {
				r.logger.Errorf("Failed to marshal member for set %s: %v", key, err)
				return err
			}
			stringMembers = append(stringMembers, string(jsonValue))
		}
	}

	err := r.client.SAdd(ctx, key, stringMembers...).Err()
	if err != nil {
		r.logger.Errorf("Failed to add members to set %s: %v", key, err)
		return err
	}

	r.logger.Infof("Successfully added members to set: %s", key)
	return nil
}

// SRem removes members from a set
func (r *RedisRepository) SRem(ctx context.Context, key string, members ...interface{}) error {
	var stringMembers []interface{}

	for _, member := range members {
		switch v := member.(type) {
		case string:
			stringMembers = append(stringMembers, v)
		default:
			jsonValue, err := json.Marshal(member)
			if err != nil {
				r.logger.Errorf("Failed to marshal member for set removal %s: %v", key, err)
				return err
			}
			stringMembers = append(stringMembers, string(jsonValue))
		}
	}

	err := r.client.SRem(ctx, key, stringMembers...).Err()
	if err != nil {
		r.logger.Errorf("Failed to remove members from set %s: %v", key, err)
		return err
	}

	r.logger.Infof("Successfully removed members from set: %s", key)
	return nil
}

// SMembers gets all members of a set
func (r *RedisRepository) SMembers(ctx context.Context, key string) ([]string, error) {
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		r.logger.Errorf("Failed to get members from set %s: %v", key, err)
		return nil, err
	}

	r.logger.Infof("Successfully retrieved members from set: %s", key)
	return members, nil
}

// SIsMember checks if a member exists in a set
func (r *RedisRepository) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	var stringMember string

	switch v := member.(type) {
	case string:
		stringMember = v
	default:
		jsonValue, err := json.Marshal(member)
		if err != nil {
			r.logger.Errorf("Failed to marshal member for set membership check %s: %v", key, err)
			return false, err
		}
		stringMember = string(jsonValue)
	}

	result, err := r.client.SIsMember(ctx, key, stringMember).Result()
	if err != nil {
		r.logger.Errorf("Failed to check membership in set %s: %v", key, err)
		return false, err
	}

	r.logger.Infof("Successfully checked membership in set: %s", key)
	return result, nil
}

// Keys gets keys matching a pattern
func (r *RedisRepository) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		r.logger.Errorf("Failed to get keys with pattern %s: %v", pattern, err)
		return nil, err
	}

	r.logger.Infof("Successfully retrieved keys with pattern: %s", pattern)
	return keys, nil
}

// TTL gets the time to live of a key
func (r *RedisRepository) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		r.logger.Errorf("Failed to get TTL for key %s: %v", key, err)
		return 0, err
	}

	r.logger.Infof("Successfully retrieved TTL for key: %s", key)
	return ttl, nil
}

// FlushDB clears all keys from the current database
func (r *RedisRepository) FlushDB(ctx context.Context) error {
	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		r.logger.Errorf("Failed to flush database: %v", err)
		return err
	}

	r.logger.Infof("Successfully flushed database")
	return nil
}
