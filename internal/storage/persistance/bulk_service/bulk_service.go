package bulk_service

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BulkServicePersistence struct {
	cfg                    *config.VaultConfig
	mongoDalCpsAction      dal.MongoDal[model.CPSAction, model.CPSAction]
	mongoDalbulkService    dal.MongoDal[model.APPAccessList, model.APPAccessList]
	kafkaProducer          kafka.ClientOrchestrationProducer
	kafkaProducerForClient kafka.NotificationProducer
	logger                 utils.Logger
}

func InitBulkServicePersistence(client *mongo.Client, cfg *config.VaultConfig, dbName string, collections []string, clientOrchestrationProducer kafka.ClientOrchestrationProducer, kafkaProducerForClient kafka.NotificationProducer, logger utils.Logger) storage.BulkServiceRepository {
	mongoDalCpsAction := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, cfg, dbName, collections[0])
	mongoDalbulkService := dal.NewMongoDal[model.APPAccessList, model.APPAccessList](client, cfg, dbName, collections[1])
	return &BulkServicePersistence{
		cfg:                    cfg,
		mongoDalCpsAction:      mongoDalCpsAction,
		mongoDalbulkService:    mongoDalbulkService,
		kafkaProducer:          clientOrchestrationProducer,
		kafkaProducerForClient: kafkaProducerForClient,
		logger:                 logger,
	}
}

func (b BulkServicePersistence) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.APPAccessList], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"ussd_enabled", "enabled"}
	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"access_list_name": searchRegex},
			{"key": searchRegex},
		}
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := b.mongoDalbulkService.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return empty paginated response
			meta := local_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage)
			return &types.PaginatedResponse[[]model.APPAccessList]{
				Data: []model.APPAccessList{},
				Meta: meta,
			}, nil
		}
		b.logger.Errorf("[BulkServicePersistence][FindAllWithPagination] failed to fetch bulk services: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	b.logger.Infof("[BulkServicePersistence][FindAllWithPagination] fetched %d bulk services with filter %v result %v", len(data), filter, data)

	// 6. Count total
	total, err := b.mongoDalbulkService.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[BulkServicePersistence][FindAllWithPagination] failed to count bulk services: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]model.APPAccessList]{
		Data: data,
		Meta: meta,
	}, nil
}

func (b BulkServicePersistence) FindAll(ctx context.Context) ([]model.APPAccessList, error) {
	b.logger.Infof("[BulkServicePersistence][FindAll] fetching all bulk services")
	filter := bson.M{}

	projection := bson.M{}
	bulkServices, err := b.mongoDalbulkService.FindAll(ctx, filter, projection)
	if err != nil {
		b.logger.Errorf("[BulkServicePersistence][FindAll] failed to fetch bulk services: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	b.logger.Infof("[BulkServicePersistence][FindAll] retrieved %d bulk services", len(bulkServices))
	return bulkServices, nil
}

func (b BulkServicePersistence) Update(ctx context.Context, keys []string, state bool) error {
	b.logger.Infof("[Update] updating bulk services, enabled: %v keys: %v", state, keys)
	// parentKeys := []string{}

	// for _, key := range keys {
	// 	// Try updating parent
	// 	parentFilter := bson.M{"key": key}
	// 	parentUpdate := bson.M{"enabled": state}

	// 	result, err := b.mongoDalbulkService.UpdateOne(ctx, parentFilter, parentUpdate)
	// 	if err != nil {
	// 		if err == mongo.ErrNoDocuments {
	// 			// Parent not found → update child
	// 			childFilter := bson.M{"sub_access_list.key": key}
	// 			childUpdate := bson.M{"sub_access_list.$.enabled": state}
	// 			child, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
	// 			if err != nil {
	// 				b.logger.Errorf("[Update] failed to update child: %v", err)
	// 				return errors.New(localization.ErrorFailToUpdateChild.Code)
	// 			}
	// 			b.kafkaProducer.PublishMessage(ctx, child, string(constants.ClientOrchestrationAccessControlTopic), string(constants.ClientOrchestrationAccessControlTopic), "bulk enable/disable access-list child")
	// 			parentKeys = append(parentKeys, child.Key)

	// 			continue
	// 		}
	// 		// Other parent update errors
	// 		b.logger.Errorf("[Update] failed to update parent: %v", err)
	// 		return errors.New(localization.ErrorFailToUpdateParent.Code)
	// 	}

	// 	b.kafkaProducer.PublishMessage(ctx, result, string(constants.ClientOrchestrationAccessControlTopic), string(constants.ClientOrchestrationAccessControlTopic), "bulk enable/disable access-list parent")

	// 	// Parent exists → manually loop over children and update each one
	// 	for _, sub := range result.SubAccessList {
	// 		childFilter := bson.M{"sub_access_list.key": sub.Key}
	// 		childUpdate := bson.M{"sub_access_list.$.enabled": state}

	// 		_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
	// 		if err != nil {
	// 			b.logger.Errorf("[Update] failed to update child: %v", err)
	// 			return errors.New(localization.ErrorFailToUpdateChild.Code)
	// 		}
	// 	}

	// }

	publishBody := []model.APPAccessList{}
	if state {
		for _, key := range keys {
			parentFilter := bson.M{"key": key}
			parentUpdate := bson.M{"enabled": state}

			updatedParent, err := b.mongoDalbulkService.UpdateOne(ctx, parentFilter, parentUpdate)
			if err != nil {
				b.logger.Errorf("[UpdateBulkService] failed to update access list: %v", err)
				return errors.New(localization.ErrorFailToUpdateBulkService.Code)
			}
			publishBody = append(publishBody, updatedParent)
		}
	}

	b.kafkaProducer.PublishMessage(ctx, publishBody, string(b.cfg.KafkaBulkServiceUpdateTopic), string(b.cfg.KafkaBulkServiceUpdateTopic), "bulk enable/disable access-list parent")

	return nil
}
