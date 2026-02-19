package icon

import (
	"cbe-super-app-cps-action/internal/constants/lib"
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

type IconStorage struct {
	dal    dal.MongoDal[model.Icon, model.Icon]
	client *mongo.Client
	logger utils.Logger
}

func NewIconRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.IconRepository {
	return &IconStorage{
		dal:    dal.NewMongoDal[model.Icon, model.Icon](client, cfg, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (i *IconStorage) Create(ctx context.Context, icon *model.Icon) error {
	_, err := i.dal.InsertOne(ctx, *icon)
	if err != nil {
		i.logger.Errorf("[Create] failed to create icon: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (i *IconStorage) Update(ctx context.Context, id string, icon *model.Icon) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		i.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := IconMapper(*icon)

	_, err = i.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		i.logger.Errorf("[Update] failed to update icon: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (i *IconStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		i.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	if err := i.dal.DeleteOne(ctx, filter); err != nil {
		i.logger.Errorf("[Delete] failed to delete icon: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (i *IconStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		i.logger.Errorf("[EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	_, err = i.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		i.logger.Errorf("[EnableOrDisable] failed to enable/disable icon: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (i *IconStorage) FindByID(ctx context.Context, id string) (*model.Icon, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		i.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := i.dal.FindOne(ctx, filter, nil)

	if err != nil {
		i.logger.Errorf("[FindByID] failed to find icon: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (s *IconStorage) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]model.Icon], error) {
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	allowedKeys := []string{"enabled"}

	filter, skip, limit := lib.FilterBuilder(*filterParam, searchKeys, allowedKeys)

	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to fetch icons: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count icons: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.Icon]{
		Data: data,
		Meta: meta,
	}, nil
}
