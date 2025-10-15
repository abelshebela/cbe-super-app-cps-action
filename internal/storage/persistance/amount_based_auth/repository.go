package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

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

func NewAmountBasedAuthRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.AmountBasedAuthRepository {
	return &AmountBasedAuthStorage{
		dal:    dal.NewMongoDal[model.AuthTier, model.AuthTier](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (a *AmountBasedAuthStorage) Update(ctx context.Context, id string, authTier *model.AuthTier) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("Failed to convert id to ObjectID: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AuthTierMapper(*authTier)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("UpdateOne failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AmountBasedAuthStorage) FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]*model.AuthTier, error) {
	result, err := a.dal.FindAll(ctx, filter, projection)
	if err != nil {
		a.logger.Errorf("FindAll failed with filter %v: %v", filter, err)
		return nil, err
	}

	return result, nil
}

func (a *AmountBasedAuthStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AuthTier], error) {
	filter := bson.M{"is_deleted": false}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["method"] = searchRegex
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := a.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := a.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.AuthTier]{
		Data: data,
		Meta: meta,
	}, nil
}

func (a *AmountBasedAuthStorage) FindByID(ctx context.Context, id string) (*model.AuthTier, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := a.dal.FindOne(ctx, filter, nil)

	if err != nil {
		return nil, err
	}
	return result, nil
}
