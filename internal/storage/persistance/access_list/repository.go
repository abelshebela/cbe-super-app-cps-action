package access_list

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

type AccessListStorage struct {
	dal    dal.MongoDal[model.APPAccessList, model.APPAccessList]
	client *mongo.Client
	logger utils.Logger
}

func NewAccessListRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.AppAccessListRepository {
	return &AccessListStorage{
		dal:    dal.NewMongoDal[model.APPAccessList, model.APPAccessList](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (a *AccessListStorage) Create(ctx context.Context, accessList *model.APPAccessList) error {

	_, err := a.dal.InsertOne(ctx, *accessList)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *AccessListStorage) Update(ctx context.Context, id string, accessList *model.APPAccessList) error {
	a.logger.Infof("[Update] updating access list for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AccessListMapper(*accessList)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[Update] access list not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[Update] failed to update access list: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[Update] access list updated successfully")
	return nil
}

func (a *AccessListStorage) Delete(ctx context.Context, id string) error {
	a.logger.Infof("[Delete] deleting access list for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = a.dal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("[Delete] failed to delete access list: %v", err)
		return err
	}
	a.logger.Infof("[Delete] access list deleted successfully")
	return nil
}

func (a *AccessListStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	a.logger.Infof("[EnableOrDisable] processing access list enable/disable for id: %s, enabled: %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	_, err = a.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[EnableOrDisable] access list not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[EnableOrDisable] failed to enable/disable access list: %v", err)
		return err
	}
	a.logger.Infof("[EnableOrDisable] access list enable/disable completed successfully")
	return nil
}

func (a *AccessListStorage) FindByID(ctx context.Context, id string) (*model.APPAccessList, error) {
	a.logger.Infof("[FindByID] fetching access list by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[FindByID] failed to find access list: %v", err)
		return nil, err
	}
	a.logger.Infof("[FindByID] access list retrieved successfully")
	return result, nil
}

func (a *AccessListStorage) FindAllWithPagination(ctx context.Context, department string, filterParam types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"access_list_name", "ussd_enabled", "enabled"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["access_list_name"] = searchRegex
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := a.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		a.logger.Errorf("[FindAllWithPagination] failed to fetch access lists: %v", err)
		return nil, err
	}

	total, err := a.dal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("[FindAllWithPagination] failed to count access lists: %v", err)
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	a.logger.Infof("[FindAllWithPagination] retrieved %d access lists", len(data))

	return &types.PaginatedResponse[[]*model.APPAccessList]{
		Data: data,
		Meta: meta,
	}, nil
}
