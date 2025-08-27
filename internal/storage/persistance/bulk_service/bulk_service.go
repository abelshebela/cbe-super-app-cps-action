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

func (b BulkServicePersistence) FindAllWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error) {
	filter := bson.M{}
	projection := bson.M{}

	// Add search functionality
	if filterParams.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"key": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"accessListName": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			},
		}
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	bulkServices, err := b.mongoDalbulkService.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	totalDocs, err := b.mongoDalbulkService.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	meta := local_util.BuildPaginationMeta(totalDocs, page, limit)

	return &types.PaginatedResponse[[]*model.APPAccessList]{
		Data: bulkServices,
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

// func (b BulkServicePersistence) AuthorizeBulkServiceEnable(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
// 	// Extract keys from CurrentAction
// 	m, ok := action.CurrentAction.(map[string]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("INVALID_CURRENT_ACTION")
// 	}

// 	rawKeys, ok := m["keys"].([]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("INVALID_KEY_FORMAT")
// 	}

// 	var keys []string
// 	for _, v := range rawKeys {
// 		if str, ok := v.(string); ok {
// 			keys = append(keys, str)
// 		}
// 	}

// 	parentKeys := []string{}
// 	// Process each key
// 	for _, key := range keys {
// 		// Try updating parent
// 		parentFilter := bson.M{"key": key}
// 		parentUpdate := bson.M{"enabled": true}

// 		result, err := b.mongoDalbulkService.UpdateOne(ctx, parentFilter, parentUpdate)
// 		if err != nil {
// 			// If parent doesn't exist, try updating child
// 			if err == mongo.ErrNoDocuments {
// 				childFilter := bson.M{"subAccessList.key": key}
// 				childUpdate := bson.M{"subAccessList.$.enabled": true}

// 				parentDoc, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
// 				if err != nil {
// 					b.logger.Errorf("failed to update child: %v\n", err)
// 					return nil, fmt.Errorf("FAILED_TO_UPDATE_CHILD")
// 				}

// 				// Collect the parent keys
// 				exists := false
// 				for _, pk := range parentKeys {
// 					if pk == parentDoc.Key {
// 						exists = true
// 						break
// 					}
// 				}

// 				if parentDoc.Key != "" && !exists {
// 					parentKeys = append(parentKeys, parentDoc.Key)
// 				}
// 				continue
// 			}

// 			// Other errors
// 			b.logger.Errorf("failed to update parent: %v\n", err)
// 			return nil, fmt.Errorf("FAILED_TO_UPDATE_PARENT")
// 		}

// 		// Parent exists — now enable all sub-keys
// 		for _, sub := range result.SubAccessList {
// 			childFilter := bson.M{"subAccessList.key": sub.Key}
// 			childUpdate := bson.M{"subAccessList.$.enabled": true}

// 			_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
// 			if err != nil {
// 				b.logger.Errorf("failed to update child: %v\n", err)
// 				return nil, fmt.Errorf("FAILED_TO_UPDATE_CHILD")
// 			}
// 		}
// 	}

// 	for _, parentKey := range parentKeys {
// 		filter := bson.M{
// 			"key": parentKey,
// 			"subAccessList": bson.M{
// 				"$not": bson.M{
// 					"$elemMatch": bson.M{"enabled": false},
// 				},
// 			},
// 		}

// 		update := bson.M{"enabled": true}
// 		_, err := b.mongoDalbulkService.UpdateOne(ctx, filter, update)
// 		if err != nil {
// 			if err == mongo.ErrNoDocuments {
// 				continue
// 			}
// 			b.logger.Errorf("failed to enable parent %s: %v", parentKey, err)
// 			return nil, fmt.Errorf("FAILED_TO_UPDATE_PARENT")
// 		}
// 	}

// 	return action, nil

// }

// func (b BulkServicePersistence) AuthorizeBulkServiceDisable(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
// 	// Extract keys from CurrentAction
// 	m, ok := action.CurrentAction.(map[string]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("INVALID_CURRENT_ACTION")
// 	}

// 	rawKeys, ok := m["keys"].([]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("INVALID_KEY_FORMAT")
// 	}

// 	var keys []string
// 	for _, v := range rawKeys {
// 		if str, ok := v.(string); ok {
// 			keys = append(keys, str)
// 		}
// 	}

// 	// Process each key
// 	parentKeys := []string{}
// 	for _, key := range keys {
// 		// Try updating parent
// 		parentFilter := bson.M{"key": key}
// 		parentUpdate := bson.M{"enabled": false}

// 		result, err := b.mongoDalbulkService.UpdateOne(ctx, parentFilter, parentUpdate)
// 		if err != nil {
// 			// If parent doesn't exist, try updating child
// 			if err == mongo.ErrNoDocuments {
// 				childFilter := bson.M{"subAccessList.key": key}
// 				childUpdate := bson.M{"subAccessList.$.enabled": false}

// 				parentDoc, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
// 				if err != nil {
// 					b.logger.Errorf("failed to update child: %v\n", err)
// 					return nil, fmt.Errorf("FAILED_TO_UPDATE_CHILD")
// 				}

// 				exists := false
// 				for _, pk := range parentKeys {
// 					if pk == parentDoc.Key {
// 						exists = true
// 						break
// 					}
// 				}

// 				if parentDoc.Key != "" && !exists {
// 					parentKeys = append(parentKeys, parentDoc.Key)
// 				}

// 				continue
// 			}

// 			// Other errors
// 			b.logger.Errorf("failed to update parent: %v\n", err)
// 			return nil, fmt.Errorf("FAILED_TO_UPDATE_PARENT")
// 		}

// 		// Parent exists — now disable all sub-keys
// 		for _, sub := range result.SubAccessList {
// 			childFilter := bson.M{"subAccessList.key": sub.Key}
// 			childUpdate := bson.M{"subAccessList.$.enabled": false}

// 			_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
// 			if err != nil {
// 				b.logger.Errorf("failed to update child: %v\n", err)
// 				return nil, fmt.Errorf("FAILED_TO_UPDATE_CHILD")
// 			}
// 		}
// 	}

// 	for _, parentKey := range parentKeys {
// 		filter := bson.M{
// 			"key": parentKey,
// 			"subAccessList": bson.M{
// 				"$not": bson.M{
// 					"$elemMatch": bson.M{"enabled": true},
// 				},
// 			},
// 		}

// 		update := bson.M{"enabled": false}
// 		_, err := b.mongoDalbulkService.UpdateOne(ctx, filter, update)
// 		if err != nil {
// 			if err == mongo.ErrNoDocuments {
// 				continue
// 			}
// 			b.logger.Errorf("failed to enable parent %s: %v", parentKey, err)
// 			return nil, fmt.Errorf("FAILED_TO_UPDATE_PARENT")
// 		}
// 	}

// 	return action, nil
// }
