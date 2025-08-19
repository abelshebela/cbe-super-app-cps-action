package archived_user

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type archivedUserStorage struct {
	dal    dal.MongoDal[model.ArchivedUser, model.ArchivedUser]
	client *mongo.Client
	logger utils.Logger
}

func NewArchivedUserRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ArchivedUserRepository {
	return &archivedUserStorage{
		dal:    dal.NewMongoDal[model.ArchivedUser, model.ArchivedUser](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (a *archivedUserStorage) Create(ctx context.Context, user *model.ArchivedUser) error {
	_, err := a.dal.InsertOne(ctx, *user)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *archivedUserStorage) FindByID(ctx context.Context, id string) (*model.ArchivedUser, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	result, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *archivedUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ArchivedUser], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["key"] = searchRegex
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

	return &types.PaginatedResponse[[]*model.ArchivedUser]{
		Data: data,
		Meta: meta,
	}, nil
}
