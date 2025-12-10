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
	avatar.ID = bson.NewObjectID()
	_, err := a.dal.InsertOne(ctx, *avatar)
	if err != nil {
		a.logger.Errorf("Unable to create avatar with error: %s", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AvatarStorage) Update(ctx context.Context, id string, avatar *model.Avatar) error {
	a.logger.Infof("[Update] updating avatar for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AvatarMapper(*avatar)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		a.logger.Errorf("[Update] failed to update avatar: %v", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[Update] avatar updated successfully")
	return nil
}

func (a *AvatarStorage) Delete(ctx context.Context, id string) error {
	a.logger.Infof("[Delete] deleting avatar for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = a.dal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("[Delete] failed to delete avatar: %v", err)
		return err
	}
	a.logger.Infof("[Delete] avatar deleted successfully")
	return nil
}

func (a *AvatarStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	a.logger.Infof("[EnableOrDisable] processing avatar enable/disable for id: %s, enabled: %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enable": enable}
	_, err = a.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[EnableOrDisable] avatar not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[EnableOrDisable] failed to enable/disable avatar: %v", err)
		return err
	}
	a.logger.Infof("[EnableOrDisable] avatar enable/disable completed successfully")
	return nil
}

func (a *AvatarStorage) FindByID(ctx context.Context, id string) (*model.Avatar, error) {
	a.logger.Infof("[FindByID] fetching avatar by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[FindByID] avatar not found")
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[FindByID] failed to find avatar: %v", err)
		return nil, err
	}
	a.logger.Infof("[FindByID] avatar retrieved successfully")
	return result, nil
}

func (a *AvatarStorage) FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]*model.Avatar, error) {
	result, err := a.dal.FindAll(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[FindAll] failed to fetch avatars: %v", err)
		return nil, err
	}
	a.logger.Infof("[FindAll] retrieved %d avatars", len(result))
	return result, nil
}

func (a *AvatarStorage) Find(ctx context.Context, filter bson.M, projection bson.M) (*model.Avatar, error) {
	a.logger.Infof("[Find] searching for avatar")
	filter["is_deleted"] = false
	avatar, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[Find] avatar not found")
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[Find] failed to find avatar: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[Find] avatar retrieved successfully")
	return avatar, nil
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
		s.logger.Errorf("[FindAllWithPagination] failed to fetch avatars: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count avatars: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	s.logger.Infof("[FindAllWithPagination] retrieved %d avatars", len(data))

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.Avatar]{
		Data: data,
		Meta: meta,
	}, nil
}
