package avatar

import (
	"cbe-super-app-cps-action/internal/constants/lib"
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

type AvatarStorage struct {
	dal    dal.MongoDal[model.Avatar, model.Avatar]
	client *mongo.Client
	logger utils.Logger
}

func NewAvatarRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.AvatarRepository {
	return &AvatarStorage{
		dal:    dal.NewMongoDal[model.Avatar, model.Avatar](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (a *AvatarStorage) Create(ctx context.Context, avatar *model.Avatar) error {
	_, err := a.dal.InsertOne(ctx, *avatar)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *AvatarStorage) Update(ctx context.Context, id string, avatar *model.Avatar) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AvatarMapper(*avatar)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *AvatarStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return a.dal.DeleteOne(ctx, filter)
}

func (a *AvatarStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	_, err = a.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
	}
	return nil
}

func (a *AvatarStorage) FindByID(ctx context.Context, id string) (*model.Avatar, error) {
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

func (a *AvatarStorage) FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]*model.Avatar, error) {
	return a.dal.FindAll(ctx, filter, projection)
}

func (s *AvatarStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Avatar], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"enable"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["label"] = searchRegex // choose your searchable field(s)
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.Avatar]{
		Data: data,
		Meta: meta,
	}, nil
}
