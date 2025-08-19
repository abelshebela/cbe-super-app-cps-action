package bps_user

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BPSUserStorage struct {
	dal    dal.MongoDal[model.BPSUser, model.BPSUser]
	client *mongo.Client
	logger utils.Logger
}

func NewBPSUserRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.BPSUserRepository {
	return &BPSUserStorage{
		dal:    dal.NewMongoDal[model.BPSUser, model.BPSUser](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (b *BPSUserStorage) GetByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error) {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	result, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *BPSUserStorage) GetAll(ctx context.Context, filter bson.M, projection bson.M) ([]*model.BPSUser, error) {
	return b.dal.FindAll(ctx, filter, projection)
}

func (b *BPSUserStorage) GetAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error) {
	filter := bson.M{"is_deleted": false}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"user_code": searchRegex},
			{"full_name": searchRegex},
		}
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := b.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := b.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.BPSUser]{
		Data: data,
		Meta: meta,
	}, nil
}

func (b *BPSUserStorage) DisableUser(ctx context.Context, userCode string) error {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"$set": bson.M{"enabled": false}}
	_, err := b.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *BPSUserStorage) EnableUser(ctx context.Context, userCode string) error {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"$set": bson.M{"enabled": true}}
	_, err := b.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
