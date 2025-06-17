package persistence

import (
	"context"
	"github.com/rs/zerolog/log"
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

func (repo AmountBasedAuthRepo) ApproveAmountBasedAuth(ctx context.Context, id string) (string, error) {
	objectID, idParser := bson.ObjectIDFromHex(id)
	if idParser != nil {
		return "Failed to parse Object id", nil
	}
	filter := bson.M{"_id": objectID}
	projection := bson.M{}

	tierById, err := repo.mongoDal.FindOne(ctx, filter, projection)
	if err != nil {
		return "Error while fetching tier", err
	}

	tier := bson.M{
		"minAmount": tierById.MinAmount,
		"maxAmount": tierById.MaxAmount,
		"enabled":   true,
	}

	_, updateErr := repo.mongoDal.UpdateOne(ctx, filter, tier)
	if updateErr != nil {
		return "Error while updating tier", updateErr
	}

	return "Update Success", nil
}

func (repo AmountBasedAuthRepo) UpdateAuthTier(request amount_based_auth_domain.AmountBasedAuthRequest, ctx context.Context, logger utils.Logger) (amount_based_auth_domain.AuthTier, error) {
	// 1. Use filter by ID (or unique identifier)
	objectID, idParser := bson.ObjectIDFromHex(request.Id)
	if idParser != nil {
		return amount_based_auth_domain.AuthTier{}, nil
	}
	filter := bson.M{"_id": objectID}
	projection := bson.M{}

	// 2. Fetch existing document
	tierById, err := repo.mongoDal.FindOne(ctx, filter, projection)
	if err != nil {
		log.Printf(err.Error())
		//logger.Errorf("find error: %v", err)
		return amount_based_auth_domain.AuthTier{}, nil
	}
	if tierById == nil {
		//logger.Errorf("auth tier not found")
		return amount_based_auth_domain.AuthTier{}, nil
	}

	tier := bson.M{
		"minAmount": request.MinAmount,
		"maxAmount": request.MaxAmount,
		"enabled":   false,
	}

	// Perform proper update with $set
	_, updateErr := repo.mongoDal.UpdateOne(ctx, filter, tier)

	if updateErr != nil {
		log.Printf(updateErr.Error())
		return amount_based_auth_domain.AuthTier{}, updateErr
	}

	// re-fetch updated doc
	updatedDoc, err := repo.mongoDal.FindOne(ctx, filter, projection)
	if err != nil {
		return amount_based_auth_domain.AuthTier{}, err
	}
	return *updatedDoc, nil

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
