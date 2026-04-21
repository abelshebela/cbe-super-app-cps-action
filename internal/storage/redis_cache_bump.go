package storage

import (
	"context"
	"strconv"
	"time"
)

// BumpRedisCacheKey sets key to a new monotonic version token (Unix nanoseconds).
// Use for cache invalidation / coordination; failures are ignored (same pattern as services repo).
func BumpRedisCacheKey(ctx context.Context, r RedisRepository, key string) {
	if r == nil || key == "" {
		return
	}
	_ = r.Set(ctx, key, strconv.FormatInt(time.Now().UnixNano(), 10), 0)
}
