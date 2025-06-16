package persistence

import (
	"context"
	amount_based_auth_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/amount_based_auth"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type AmountBasedAuthRepo struct {
	client   *mongo.Client
	mongoDal dal.MongoDal[amount_based_auth_domain.AuthTier, amount_based_auth_domain.AuthTier]
	timeOut  time.Duration
	logger   utils.Logger
}

func (repo AmountBasedAuthRepo) UpdateAuthTier(request amount_based_auth_domain.AmountBasedAuthRequest, ctx context.Context, logger utils.Logger) (amount_based_auth_domain.AuthTier, error) {
	// 1. Use filter by ID (or unique identifier)
	filter := bson.M{"_id": request.Id}
	projection := bson.M{}

	// 2. Fetch existing document
	tierById, err := repo.mongoDal.FindOne(ctx, filter, projection)
	if err != nil {
		logger.Errorf("find error: %w", err)
		return amount_based_auth_domain.AuthTier{}, nil
	}
	if tierById == nil {
		logger.Errorf("auth tier not found")
		return amount_based_auth_domain.AuthTier{}, nil
	}

	// 3. Update fields
	tier := *tierById
	tier.MinAmount = uint64(request.MinAmount)
	tier.MaxAmount = uint64(request.MaxAmount)

	// 4. Apply update
	_, err = repo.mongoDal.UpdateOne(ctx, filter, bson.M{"$set": tier})
	if err != nil {
		logger.Errorf("failed to update auth tier: %w", err)
		return amount_based_auth_domain.AuthTier{}, nil
	}

	return tier, nil
}

func InitAmountBasedAuth(client *mongo.Client, database string, collection string) *AmountBasedAuthRepo {
	mongoDal := dal.NewMongoDal[amount_based_auth_domain.AuthTier, amount_based_auth_domain.AuthTier](client, database, collection)

	return &AmountBasedAuthRepo{
		client:   client,
		mongoDal: mongoDal,
		timeOut:  5 * time.Second,
		logger:   nil,
	}
}
