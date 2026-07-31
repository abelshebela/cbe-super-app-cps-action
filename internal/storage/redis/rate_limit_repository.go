package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"

	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// RateLimitRepository implements the RateLimitRepository interface
type RateLimitRepository struct {
	redisRepo storage.RedisRepository
	logger    utils.Logger
}

// NewRateLimitRepository creates a new rate limit repository instance
func NewRateLimitRepository(redisRepo storage.RedisRepository, logger utils.Logger) storage.RateLimitRepository {
	return &RateLimitRepository{
		redisRepo: redisRepo,
		logger:    logger,
	}
}

// IncrementRequestCount increments the request count for a key within a time window
func (r *RateLimitRepository) IncrementRequestCount(ctx context.Context, key string, window time.Duration) (int, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// Get current count
	currentCountStr, err := r.redisRepo.Get(ctx, redisKey)
	if err != nil {
		// If key doesn't exist, start with 1
		err = r.redisRepo.Set(ctx, redisKey, "1", window)
		if err != nil {
			log.Errorf("Failed to set initial request count for key %s: %v", key, err)
			return 0, err
		}
		log.Infof("Set initial request count for key: %s", key)
		return 1, nil
	}

	currentCount, err := strconv.Atoi(currentCountStr)
	if err != nil {
		log.Errorf("Failed to parse request count for key %s: %v", key, err)
		return 0, err
	}

	newCount := currentCount + 1
	err = r.redisRepo.Set(ctx, redisKey, strconv.Itoa(newCount), window)
	if err != nil {
		log.Errorf("Failed to increment request count for key %s: %v", key, err)
		return 0, err
	}

	log.Infof("Incremented request count for key %s: %d", key, newCount)
	return newCount, nil
}

// GetRequestCount gets the current request count for a key
func (r *RateLimitRepository) GetRequestCount(ctx context.Context, key string) (int, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	redisKey := fmt.Sprintf("rate_limit:%s", key)

	countStr, err := r.redisRepo.Get(ctx, redisKey)
	if err != nil {
		log.Warnf("No request count found for key %s", key)
		return 0, nil
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Errorf("Failed to parse request count for key %s: %v", key, err)
		return 0, err
	}

	log.Infof("Retrieved request count for key %s: %d", key, count)
	return count, nil
}

// ResetRequestCount resets the request count for a key
func (r *RateLimitRepository) ResetRequestCount(ctx context.Context, key string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	redisKey := fmt.Sprintf("rate_limit:%s", key)

	err := r.redisRepo.Delete(ctx, redisKey)
	if err != nil {
		log.Errorf("Failed to reset request count for key %s: %v", key, err)
		return err
	}

	log.Infof("Successfully reset request count for key: %s", key)
	return nil
}

// IsRateLimited checks if a key is rate limited based on a limit
func (r *RateLimitRepository) IsRateLimited(ctx context.Context, key string, limit int) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	count, err := r.GetRequestCount(ctx, key)
	if err != nil {
		log.Errorf("Failed to get request count for rate limit check %s: %v", key, err)
		return false, err
	}

	isLimited := count >= limit
	log.Infof("Rate limit check for key %s: count=%d, limit=%d, limited=%t", key, count, limit, isLimited)
	return isLimited, nil
}

// IncrementIPRequestCount increments the request count for an IP address
func (r *RateLimitRepository) IncrementIPRequestCount(ctx context.Context, ip string, window time.Duration) (int, error) {
	key := fmt.Sprintf("ip_rate_limit:%s", ip)
	return r.IncrementRequestCount(ctx, key, window)
}

// GetIPRequestCount gets the current request count for an IP address
func (r *RateLimitRepository) GetIPRequestCount(ctx context.Context, ip string) (int, error) {
	key := fmt.Sprintf("ip_rate_limit:%s", ip)
	return r.GetRequestCount(ctx, key)
}

// BlockIP blocks an IP address for a specified duration
func (r *RateLimitRepository) BlockIP(ctx context.Context, ip string, duration time.Duration) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	key := fmt.Sprintf("blocked_ip:%s", ip)

	err := r.redisRepo.Set(ctx, key, "blocked", duration)
	if err != nil {
		log.Errorf("Failed to block IP %s: %v", ip, err)
		return err
	}

	log.Infof("Successfully blocked IP: %s for %v", ip, duration)
	return nil
}

// IsIPBlocked checks if an IP address is blocked
func (r *RateLimitRepository) IsIPBlocked(ctx context.Context, ip string) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	key := fmt.Sprintf("blocked_ip:%s", ip)

	exists, err := r.redisRepo.Exists(ctx, key)
	if err != nil {
		log.Errorf("Failed to check if IP %s is blocked: %v", ip, err)
		return false, err
	}

	log.Infof("IP %s blocked status: %t", ip, exists)
	return exists, nil
}

// UnblockIP unblocks an IP address
func (r *RateLimitRepository) UnblockIP(ctx context.Context, ip string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	key := fmt.Sprintf("blocked_ip:%s", ip)

	err := r.redisRepo.Delete(ctx, key)
	if err != nil {
		log.Errorf("Failed to unblock IP %s: %v", ip, err)
		return err
	}

	log.Infof("Successfully unblocked IP: %s", ip)
	return nil
}
