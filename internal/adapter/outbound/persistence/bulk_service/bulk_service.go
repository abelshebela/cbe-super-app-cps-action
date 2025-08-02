package bulkservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulk_service"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bulk_service"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BulkServicePersistence struct {
	mongoDalCpsAction   dal.MongoDal[model.CPSAction, model.CPSAction]
	mongoDalbulkService dal.MongoDal[domain.APPAccessList, domain.APPAccessList]
	logger              utils.Logger
}

func InitBulkServicePersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) outbound.BulkServiceRepository {
	mongoDalCpsAction := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collections[0])
	mongoDalbulkService := dal.NewMongoDal[domain.APPAccessList, domain.APPAccessList](client, dbName, collections[1])
	return &BulkServicePersistence{
		mongoDalCpsAction:   mongoDalCpsAction,
		mongoDalbulkService: mongoDalbulkService,
		logger:              logger,
	}
}

func (b BulkServicePersistence) GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*domain.APPAccessList], error) {
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
		return nil, err
	}

	totalDocs, err := b.mongoDalbulkService.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := common_util.BuildPaginationMeta(totalDocs, page, limit)

	return &common_util.PaginatedResponse[[]*domain.APPAccessList]{
		Data: bulkServices,
		Meta: meta,
	}, nil
}

func (b BulkServicePersistence) EnableOrDisableBulkService(ctx context.Context, keys []string, cpsAction model.CPSAction, requestActionType model.RequestAction) error {
	// Check if there is a pending action
	makerData := ctx_util.ExtractContext(ctx)
	pendingFilter := bson.M{
		"maker_id":       makerData.UserID,
		"department":     makerData.Department,
		"action_status":  "PENDING",
		"request_action": requestActionType,
	}
	projection := bson.M{}
	pendingAction, err := b.mongoDalCpsAction.FindOne(ctx, pendingFilter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		b.logger.Errorf("database error occured while checking a pending action: %v\n", err)
		return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}
	if pendingAction != nil {
		return fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	// Check if the keys do exist
	allAccessLists, err := b.mongoDalbulkService.FindAll(ctx, bson.M{}, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			b.logger.Errorf("no resource found: %v\n", err)
			return fmt.Errorf("NO_RESOURCE_FOUND")
		}
		b.logger.Errorf("database error occured: %v\n", err)
		return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}

	validKeys := make(map[string]bool)
	for _, access := range allAccessLists {
		validKeys[access.Key] = true
		for _, sub := range access.SubAccessList {
			validKeys[sub.Key] = true
		}
	}

	// Now validate the incoming keys
	var notFoundKeys []string
	for _, key := range keys {
		if !validKeys[key] {
			notFoundKeys = append(notFoundKeys, key)
		}
	}

	if len(notFoundKeys) > 0 {
		return fmt.Errorf("SURVICE_NOT_FOUND")
	}

	// Check if one or more key is already enable or disabled.
	dupAction := []string{}
	for _, k := range keys {
		var boolStatus bool
		switch strings.ToUpper(cpsAction.ActionType) {
		case "ENABLE":
			boolStatus = true
		case "DISABLE":
			boolStatus = false
		}

		// Find the key in the parent
		pFilter := bson.M{"key": k, "enabled": boolStatus}
		_, err := b.mongoDalbulkService.FindOne(ctx, pFilter, bson.M{})
		if err == nil {
			dupAction = append(dupAction, k)
		} else if err != mongo.ErrNoDocuments {
			b.logger.Errorf("database error occured: %v\n", err)
			return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
		}

		cFilter := bson.M{
			"subAccessList": bson.M{
				"$elemMatch": bson.M{
					"key":     k,
					"enabled": boolStatus,
				},
			},
		}
		_, subErr := b.mongoDalbulkService.FindOne(ctx, cFilter, bson.M{})
		if subErr == nil {
			dupAction = append(dupAction, k)
			continue
		} else if subErr != mongo.ErrNoDocuments {
			b.logger.Errorf("database error occured: %v\n", err)
			return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
		}
	}

	if len(dupAction) > 0 {
		b.logger.Errorf("Duplicate action sent for keys: %v\n", dupAction)
		return fmt.Errorf("DUPLICATE_ACTION")
	}

	// Create the CPS action
	_, err = b.mongoDalCpsAction.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("database error occured: %v\n", err)
		return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}

	return nil
}

func (b BulkServicePersistence) AuthorizeBulkServiceEnable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	// Extract keys from CurrentAction
	m, ok := action.CurrentAction.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("INVALID_CURRENT_ACTION")
	}

	rawKeys, ok := m["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("INVALID_KEY_FORMAT")
	}

	var keys []string
	for _, v := range rawKeys {
		if str, ok := v.(string); ok {
			keys = append(keys, str)
		}
	}

	// Process each key
	for _, key := range keys {
		// Try updating parent
		parentFilter := bson.M{"key": key}
		parentUpdate := bson.M{"enabled": true}

		result, err := b.mongoDalbulkService.UpdateOne(ctx, parentFilter, parentUpdate)
		if err != nil {
			// If parent doesn't exist, try updating child
			if err == mongo.ErrNoDocuments {
				childFilter := bson.M{"subAccessList.key": key}
				childUpdate := bson.M{"subAccessList.$.enabled": true}

				_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
				if err != nil {
					b.logger.Errorf("failed to update child: %v\n", err)
					return nil, fmt.Errorf("FAILED_TO_UPDATE_CHILD")
				}
				continue
			}

			// Other errors
			b.logger.Errorf("failed to update parent: %v\n", err)
			return nil, fmt.Errorf("FAILED_TO_UPDATE_PARENT")
		}

		// Parent exists — now disable all sub-keys
		for _, sub := range result.SubAccessList {
			childFilter := bson.M{"subAccessList.key": sub.Key}
			childUpdate := bson.M{"subAccessList.$.enabled": true}

			_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
			if err != nil {
				b.logger.Errorf("failed to update child: %v\n", err)
				return nil, fmt.Errorf("FAILED_TO_UPDATE_CHILD")
			}
		}
	}

	return action, nil
}

func (b BulkServicePersistence) AuthorizeBulkServiceDisable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	// Extract keys from CurrentAction
	m, ok := action.CurrentAction.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("INVALID_CURRENT_ACTION")
	}

	rawKeys, ok := m["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("INVALID_KEY_FORMAT")
	}

	var keys []string
	for _, v := range rawKeys {
		if str, ok := v.(string); ok {
			keys = append(keys, str)
		}
	}

	// Process each key
	for _, key := range keys {
		// Try updating parent
		parentFilter := bson.M{"key": key}
		parentUpdate := bson.M{"enabled": false}

		result, err := b.mongoDalbulkService.UpdateOne(ctx, parentFilter, parentUpdate)
		if err != nil {
			// If parent doesn't exist, try updating child
			if err == mongo.ErrNoDocuments {
				childFilter := bson.M{"subAccessList.key": key}
				childUpdate := bson.M{"subAccessList.$.enabled": false}

				_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
				if err != nil {
					b.logger.Errorf("failed to update child: %v\n", err)
					return nil, fmt.Errorf("FAILED_TO_UPDATE_CHILD")
				}
				continue
			}

			// Other errors
			b.logger.Errorf("failed to update parent: %v\n", err)
			return nil, fmt.Errorf("FAILED_TO_UPDATE_PARENT")
		}

		// Parent exists — now disable all sub-keys
		for _, sub := range result.SubAccessList {
			childFilter := bson.M{"subAccessList.key": sub.Key}
			childUpdate := bson.M{"subAccessList.$.enabled": false}

			_, err := b.mongoDalbulkService.UpdateOne(ctx, childFilter, childUpdate)
			if err != nil {
				b.logger.Errorf("failed to update child: %v\n", err)
				return nil, fmt.Errorf("FAILED_TO_UPDATE_CHILD")
			}
		}
	}

	return action, nil
}
