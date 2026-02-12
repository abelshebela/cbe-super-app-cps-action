package redis

import (
	"context"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// TemporaryDataRepository implements the TemporaryDataRepository interface
type TemporaryDataRepository struct {
	redisRepo storage.RedisRepository
	logger    utils.Logger
}

// NewTemporaryDataRepository creates a new temporary data repository instance
func NewTemporaryDataRepository(redisRepo storage.RedisRepository, logger utils.Logger) storage.TemporaryDataRepository {
	return &TemporaryDataRepository{
		redisRepo: redisRepo,
		logger:    logger,
	}
}

// SaveTemporaryData saves temporary data with expiration
func (t *TemporaryDataRepository) SaveTemporaryData(ctx context.Context, key string, data map[string]interface{}, expiration time.Duration) error {
	redisKey := fmt.Sprintf("temp_data:%s", key)

	for field, value := range data {
		err := t.redisRepo.HSet(ctx, redisKey, field, value)
		if err != nil {
			t.logger.Errorf("Failed to save temporary data field %s for key %s: %v", field, key, err)
			return err
		}
	}

	// Set expiration for the entire hash
	err := t.redisRepo.Expire(ctx, redisKey, expiration)
	if err != nil {
		t.logger.Errorf("Failed to set expiration for temporary data %s: %v", key, err)
		return err
	}

	t.logger.Infof("Successfully saved temporary data: %s", key)
	return nil
}

// GetTemporaryData retrieves temporary data
func (t *TemporaryDataRepository) GetTemporaryData(ctx context.Context, key string) (map[string]string, error) {
	redisKey := fmt.Sprintf("temp_data:%s", key)

	data, err := t.redisRepo.HGetAll(ctx, redisKey)
	if err != nil {
		t.logger.Errorf("Failed to get temporary data %s: %v", key, err)
		return nil, err
	}

	t.logger.Infof("Successfully retrieved temporary data: %s", key)
	return data, nil
}

// DeleteTemporaryData deletes temporary data
func (t *TemporaryDataRepository) DeleteTemporaryData(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("temp_data:%s", key)

	err := t.redisRepo.Delete(ctx, redisKey)
	if err != nil {
		t.logger.Errorf("Failed to delete temporary data %s: %v", key, err)
		return err
	}

	t.logger.Infof("Successfully deleted temporary data: %s", key)
	return nil
}

// UpdateTemporaryData updates specific fields in temporary data
func (t *TemporaryDataRepository) UpdateTemporaryData(ctx context.Context, key string, updates map[string]interface{}) error {
	redisKey := fmt.Sprintf("temp_data:%s", key)

	for field, value := range updates {
		err := t.redisRepo.HSet(ctx, redisKey, field, value)
		if err != nil {
			t.logger.Errorf("Failed to update temporary data field %s for key %s: %v", field, key, err)
			return err
		}
	}

	t.logger.Infof("Successfully updated temporary data: %s", key)
	return nil
}

// SaveFormData saves form data for a session
func (t *TemporaryDataRepository) SaveFormData(ctx context.Context, sessionID string, formData map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("form_data:%s", sessionID)

	for field, value := range formData {
		err := t.redisRepo.HSet(ctx, key, field, value)
		if err != nil {
			t.logger.Errorf("Failed to save form data field %s for session %s: %v", field, sessionID, err)
			return err
		}
	}

	// Set expiration for the entire hash
	err := t.redisRepo.Expire(ctx, key, expiration)
	if err != nil {
		t.logger.Errorf("Failed to set expiration for form data %s: %v", sessionID, err)
		return err
	}

	t.logger.Infof("Successfully saved form data for session: %s", sessionID)
	return nil
}

// GetFormData retrieves form data for a session
func (t *TemporaryDataRepository) GetFormData(ctx context.Context, sessionID string) (map[string]string, error) {
	key := fmt.Sprintf("form_data:%s", sessionID)

	formData, err := t.redisRepo.HGetAll(ctx, key)
	if err != nil {
		t.logger.Errorf("Failed to get form data for session %s: %v", sessionID, err)
		return nil, err
	}

	t.logger.Infof("Successfully retrieved form data for session: %s", sessionID)
	return formData, nil
}

// DeleteFormData deletes form data for a session
func (t *TemporaryDataRepository) DeleteFormData(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("form_data:%s", sessionID)

	err := t.redisRepo.Delete(ctx, key)
	if err != nil {
		t.logger.Errorf("Failed to delete form data for session %s: %v", sessionID, err)
		return err
	}

	t.logger.Infof("Successfully deleted form data for session: %s", sessionID)
	return nil
}

// SetCache sets a cache value with expiration
func (t *TemporaryDataRepository) SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("cache:%s", key)

	err := t.redisRepo.Set(ctx, cacheKey, value, expiration)
	if err != nil {
		t.logger.Errorf("Failed to set cache for key %s: %v", key, err)
		return err
	}

	t.logger.Infof("Successfully set cache: %s", key)
	return nil
}

// GetCache retrieves a cache value
func (t *TemporaryDataRepository) GetCache(ctx context.Context, key string) (string, error) {
	cacheKey := fmt.Sprintf("cache:%s", key)

	value, err := t.redisRepo.Get(ctx, cacheKey)
	if err != nil {
		t.logger.Errorf("Failed to get cache for key %s: %v", key, err)
		return "", err
	}

	t.logger.Infof("Successfully retrieved cache: %s", key)
	return value, nil
}

// DeleteCache deletes a cache value
func (t *TemporaryDataRepository) DeleteCache(ctx context.Context, key string) error {
	cacheKey := fmt.Sprintf("cache:%s", key)

	err := t.redisRepo.Delete(ctx, cacheKey)
	if err != nil {
		t.logger.Errorf("Failed to delete cache for key %s: %v", key, err)
		return err
	}

	t.logger.Infof("Successfully deleted cache: %s", key)
	return nil
}

// ClearCacheByPattern clears cache keys matching a pattern
func (t *TemporaryDataRepository) ClearCacheByPattern(ctx context.Context, pattern string) error {
	cachePattern := fmt.Sprintf("cache:%s", pattern)

	keys, err := t.redisRepo.Keys(ctx, cachePattern)
	if err != nil {
		t.logger.Errorf("Failed to get cache keys for pattern %s: %v", pattern, err)
		return err
	}

	if len(keys) == 0 {
		t.logger.Infof("No cache keys found for pattern: %s", pattern)
		return nil
	}

	// Delete all matching keys
	for _, key := range keys {
		err := t.redisRepo.Delete(ctx, key)
		if err != nil {
			t.logger.Errorf("Failed to delete cache key %s: %v", key, err)
			return err
		}
	}

	t.logger.Infof("Successfully cleared %d cache keys for pattern: %s", len(keys), pattern)
	return nil
}
