package token_provider_service

import (
	"context"
	"time"

	tp_client "cbe-super-app-cps-action/internal/storage/external_call/token_provider"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"golang.org/x/sync/singleflight"
)

const redisKey = "account_creation_token"

type tokenProviderService struct {
	client *tp_client.TokenProviderClient
	redis  storage.RedisRepository
	logger utils.Logger
	group  singleflight.Group
}

func NewTokenProviderService(client *tp_client.TokenProviderClient, redis storage.RedisRepository, logger utils.Logger) *tokenProviderService {
	return &tokenProviderService{
		client: client,
		redis:  redis,
		logger: logger,
	}
}

func (s *tokenProviderService) GetToken(ctx context.Context) (string, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	token, err := s.redis.Get(ctx, redisKey)
	if err == nil && token != "" {
		return token, nil
	}

	if err != nil {
		log.Warnf("[TokenProviderSvc] Redis get: %v — fetching fresh token", err)
	}

	result, fetchErr, _ := s.group.Do(redisKey, func() (interface{}, error) {
		// Re-check inside singleflight in case a concurrent call already refreshed it
		if t, rerr := s.redis.Get(ctx, redisKey); rerr == nil && t != "" {
			return t, nil
		}

		resp, err := s.client.FetchToken(ctx)
		if err != nil {
			return nil, err
		}

		ttl := time.Duration(resp.ExpiresIn) * time.Second
		if setErr := s.redis.Set(ctx, redisKey, resp.AccessToken, ttl); setErr != nil {
			log.Errorf("[TokenProviderSvc] Redis set: %v", setErr)
		} else {
			log.Infof("[TokenProviderSvc] token cached with TTL=%s", ttl)
		}

		return resp.AccessToken, nil
	})

	if fetchErr != nil {
		log.Errorf("[TokenProviderSvc] fetch: %v", fetchErr)
		return "", fetchErr
	}

	return result.(string), nil
}
