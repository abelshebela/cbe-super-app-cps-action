package access_list

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	config "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AccessListStorage struct {
	dal           dal.MongoDal[model.APPAccessList, model.APPAccessList]
	client        *mongo.Client
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewAccessListRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, clientOrchestrationProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.AppAccessListRepository {
	return &AccessListStorage{
		dal:           dal.NewMongoDal[model.APPAccessList, model.APPAccessList](client, cfg, dbName, collection),
		client:        client,
		kafkaProducer: clientOrchestrationProducer,
		logger:        logger,
	}
}

func (a *AccessListStorage) Create(ctx context.Context, accessList *model.APPAccessList) error {

	newAccessControl, err := a.dal.InsertOne(ctx, *accessList)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	a.kafkaProducer.PublishMessage(ctx, newAccessControl, string(constants.ClientOrchestrationAccessControlTopic), string(constants.ClientOrchestrationAccessControlTopic), "create access-list")

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

	updateAccessList, err := a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[Update] access list not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[Update] failed to update access list: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	a.kafkaProducer.PublishMessage(ctx, updateAccessList, string(constants.ClientOrchestrationAccessControlTopic), string(constants.ClientOrchestrationAccessControlTopic), "update access-list")

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
	updateAccessList, err := a.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[EnableOrDisable] access list not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[EnableOrDisable] failed to enable/disable access list: %v", err)
		return err
	}

	a.kafkaProducer.PublishMessage(ctx, updateAccessList, string(constants.ClientOrchestrationAccessControlTopic), string(constants.ClientOrchestrationAccessControlTopic), "enable/disable access-list")

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

func (a *AccessListStorage) FindAllWithPagination(ctx context.Context, department string, filterParam types.Filter) (*types.PaginatedResponse[[]model.APPAccessList], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search", "access_list_name", "ussd_enabled", "enabled"}
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

	return &types.PaginatedResponse[[]model.APPAccessList]{
		Data: data,
		Meta: meta,
	}, nil
}

func (a *AccessListStorage) FindByKeys(ctx context.Context, keys []string) (map[string]string, error) {
	a.logger.Infof("[FindByKeys] checking access list for keys")
	// Match if any key in keys is present in either the parent key or any sub_access_list.key
	filter := bson.M{
		"$or": []bson.M{
			{"key": bson.M{"$in": keys}},
			{"sub_access_list.key": bson.M{"$in": keys}},
		},
	}

	als, err := a.dal.FindAll(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[FindByKeys] failed to count access lists: %v", err)
		return map[string]string{}, err
	}

	// Build a set of found keys from both key and sub_access_list.key
	foundKeys := make(map[string]string)
	for _, al := range als {
		foundKeys[al.Key] = al.AccessListName
		// Sub access list keys
		if al.SubAccessList != nil {
			for _, sub := range al.SubAccessList {
				foundKeys[sub.Key] = sub.AccessListName
			}
		}
	}

	var missing []string
	for _, key := range keys {
		if _, ok := foundKeys[key]; !ok {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		msg := "keys: " + strings.Join(missing, ", ") + " not found in access list"
		if len(missing) == 1 {
			msg = "key: " + strings.Join(missing, ", ") + " not found in access list"
		}
		a.logger.Errorf("[FindByKeys] %s", msg)
		return map[string]string{}, errors.New(msg)
	}
	return foundKeys, nil
}

func (a *AccessListStorage) FindAllByKeys(ctx context.Context, keys []string) ([]model.APPAccessList, error) {
	a.logger.Infof("[FindByKeys] checking access list for keys")
	filter := bson.M{"enabled": true}
	// Match if any key in keys is present in either the parent key or any sub_access_list.key
	// filter := bson.M{
	// 	"$or": []bson.M{
	// 		{"key": bson.M{"$nin": keys}},
	// 		{"sub_access_list.key": bson.M{"$nin": keys}},
	// 	},
	// }
	// filter := bson.M{
	// 	"key":                 bson.M{"$nin": keys},
	// 	"sub_access_list.key": bson.M{"$nin": keys},
	// }
	if len(keys) == 0 {
		filter = bson.M{}
	}

	als, err := a.dal.FindAll(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[FindByKeys] failed to count access lists: %v", err)
		return []model.APPAccessList{}, err
	}
	return als, nil
}
