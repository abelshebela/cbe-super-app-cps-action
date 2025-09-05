package color

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

type ColorStorage struct {
	dal    dal.MongoDal[model.Color, model.Color]
	client *mongo.Client
	logger utils.Logger
}

func NewColorRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ColorRepository {
	return &ColorStorage{
		dal:    dal.NewMongoDal[model.Color, model.Color](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (c *ColorStorage) Create(ctx context.Context, color *model.Color) error {
	_, err := c.dal.InsertOne(ctx, *color)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (c *ColorStorage) Update(ctx context.Context, id string, color *model.Color) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := ColorMapper(*color)

	_, err = c.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (c *ColorStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return c.dal.DeleteOne(ctx, filter)
}

func (c *ColorStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	_, err = c.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
	}
	return nil
}

func (c *ColorStorage) Find(ctx context.Context, filter bson.M) (*model.Color, error) {

	data, err := c.dal.FindOne(ctx, filter, nil)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		return nil, errors.New(code)
	}

	return data, nil
}

func (c *ColorStorage) FindByID(ctx context.Context, id string) (*model.Color, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := c.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (s *ColorStorage) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.Color], error) {
    searchKeys := bson.M{}
    allowedKeys := []string{"color"} 

    fbFilter, skip, limit := lib.FilterBuilder(*filterParam, searchKeys, allowedKeys)

    fbFilter["is_deleted"] = false

    data, err := s.dal.FindAllWithPagination(ctx, fbFilter, bson.M{}, skip, limit)
    if err != nil {
        return nil, errors.New(localization.ErrorUnexpectedError.Message)
    }

    total, err := s.dal.TotalCount(ctx, fbFilter)
    if err != nil {
        return nil, errors.New(localization.ErrorUnexpectedError.Message)
    }

    meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

    return &types.PaginatedResponse[[]*model.Color]{
        Data: data,
        Meta: meta,
    }, nil
}

