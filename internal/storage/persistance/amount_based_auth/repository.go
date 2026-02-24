package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AmountBasedAuthStorage struct {
	dal    dal.MongoDal[model.AuthTier, model.AuthTier]
	client *mongo.Client
	logger utils.Logger
}

func NewAmountBasedAuthRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.AmountBasedAuthRepository {
	return &AmountBasedAuthStorage{
		dal:    dal.NewMongoDal[model.AuthTier, model.AuthTier](client, cfg, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (a *AmountBasedAuthStorage) Update(ctx context.Context, id string, authTier *model.AuthTier) error {
	a.logger.Infof("[AmountBasedAuthStorage][Update] updating amount-based auth tier for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[AmountBasedAuthStorage][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AuthTierMapper(*authTier)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		a.logger.Errorf("[AmountBasedAuthStorage][Update] failed to update amount-based auth tier: %v", err)
		return local_util.HandleDBError(err)
	}
	a.logger.Infof("[AmountBasedAuthStorage][Update] amount-based auth tier updated successfully")
	return nil
}

func (a *AmountBasedAuthStorage) FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]model.AuthTier, error) {
	a.logger.Infof("[AmountBasedAuthStorage][FindAll] fetching all amount-based auth tiers")
	result, err := a.dal.FindAll(ctx, filter, projection)
	if err != nil {
		a.logger.Errorf("[AmountBasedAuthStorage][FindAll] failed to fetch amount-based auth tiers: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	a.logger.Infof("[AmountBasedAuthStorage][FindAll] retrieved %d amount-based auth tiers", len(result))
	return result, nil
}

func (a *AmountBasedAuthStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.AuthTier], error) {
	filter := bson.M{"is_deleted": false}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["method"] = searchRegex
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := a.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		a.logger.Errorf("[AmountBasedAuthStorage][FindAllWithPagination] failed to fetch amount-based auth tiers: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := a.dal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("[AmountBasedAuthStorage][FindAllWithPagination] failed to count amount-based auth tiers: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	a.logger.Infof("[AmountBasedAuthStorage][FindAllWithPagination] retrieved %d amount-based auth tiers", len(data))

	return &types.PaginatedResponse[[]model.AuthTier]{
		Data: data,
		Meta: meta,
	}, nil
}

func (a *AmountBasedAuthStorage) FindByID(ctx context.Context, id string) (*model.AuthTier, error) {
	a.logger.Infof("[AmountBasedAuthStorage][FindByID] fetching amount-based auth tier by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[AmountBasedAuthStorage][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := a.dal.FindOne(ctx, filter, nil)

	if err != nil {
		a.logger.Errorf("[AmountBasedAuthStorage][FindByID] failed to find amount-based auth tier: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	a.logger.Infof("[AmountBasedAuthStorage][FindByID] amount-based auth tier retrieved successfully")
	return result, nil
}
