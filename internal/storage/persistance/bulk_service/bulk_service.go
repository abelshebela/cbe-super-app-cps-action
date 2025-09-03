package bulk_service

import (
	// "cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"

	// "strings"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BulkServicePersistence struct {
	mongoDalCpsAction   dal.MongoDal[model.CPSAction, model.CPSAction]
	mongoDalbulkService dal.MongoDal[model.APPAccessList, model.APPAccessList]
	logger              utils.Logger
}

func InitBulkServicePersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) storage.BulkServiceRepository {
	mongoDalCpsAction := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collections[0])
	mongoDalbulkService := dal.NewMongoDal[model.APPAccessList, model.APPAccessList](client, dbName, collections[1])
	return &BulkServicePersistence{
		mongoDalCpsAction:   mongoDalCpsAction,
		mongoDalbulkService: mongoDalbulkService,
		logger:              logger,
	}
}

func (b BulkServicePersistence) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"ussd_enabled", "enabled"}
	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["access_list_name"] = searchRegex
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := b.mongoDalbulkService.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)

	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := b.mongoDalbulkService.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.APPAccessList]{
		Data: data,
		Meta: meta,
	}, nil
}

func (b BulkServicePersistence) FindAll(ctx context.Context) ([]*model.APPAccessList, error) {
	filter := bson.M{}

	projection := bson.M{}
	bulkServices, err := b.mongoDalbulkService.FindAll(ctx, filter, projection)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	return bulkServices, nil
}
func (b BulkServicePersistence) Update(ctx context.Context, keys []string, state bool) error {
	parentKeys := []string{}

	for _, key := range keys {
		// Try updating parent
		parentFilter := bson.M{"key": key}
		parentUpdate := bson.M{"$set": bson.M{"enabled": state}}

		result, err := b.mongoDalbulkService.UpdateOne(ctx, parentFilter, parentUpdate)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Parent not found → update child
				childFilter := bson.M{"subAccessList.key": key}
				childUpdate := bson.M{"$set": bson.M{"subAccessList.$.enabled": state}}
				_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
				if err != nil {
					b.logger.Errorf("failed to update child: %v", err)
					return fmt.Errorf("FAILED_TO_UPDATE_CHILD")
				}
				continue
			}

			// Other parent update errors
			b.logger.Errorf("failed to update parent: %v", err)
			return fmt.Errorf("FAILED_TO_UPDATE_PARENT")
		}

		// Parent exists → manually loop over children and update each one
		for _, sub := range result.SubAccessList {
			childFilter := bson.M{"subAccessList.key": sub.Key}
			childUpdate := bson.M{"$set": bson.M{"subAccessList.$.enabled": state}}

			_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
			if err != nil {
				b.logger.Errorf("failed to update child %s: %v", sub.Key, err)
				return fmt.Errorf("FAILED_TO_UPDATE_CHILD")
			}
		}

		// Collect parent keys for final consistency check
		parentKeys = append(parentKeys, key)
	}

	// Final step: enable parent only if all children are enabled
	for _, parentKey := range parentKeys {
		filter := bson.M{
			"key": parentKey,
			"subAccessList": bson.M{
				"$not": bson.M{"$elemMatch": bson.M{"enabled": false}},
			},
		}
		update := bson.M{"$set": bson.M{"enabled": true}}

		_, err := b.mongoDalbulkService.UpdateOne(ctx, filter, update)
		if err != nil && err != mongo.ErrNoDocuments {
			b.logger.Errorf("failed to enable parent %s: %v", parentKey, err)
			return fmt.Errorf("FAILED_TO_UPDATE_PARENT")
		}
	}
	return nil
}
