package access_list

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"

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
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AccessListMapper(*accessList)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *AccessListStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return a.dal.DeleteOne(ctx, filter)
}

func (a *AccessListStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
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

func (a *AccessListStorage) FindByID(ctx context.Context, id string) (*model.APPAccessList, error) {
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

func (a *AccessListStorage) FindAllWithPagination(ctx context.Context, department string, filterParam types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"access_list_name", "ussd_enabled", "enabled"}
	fmt.Println("*********************FindAllWithPagination**********************")
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["access_list_name"] = searchRegex
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := a.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := a.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.APPAccessList]{
		Data: data,
		Meta: meta,
	}, nil
}
