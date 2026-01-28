package account_block

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AccountBlockStorage struct {
	accountBlock  dal.MongoDal[model.AccountBlock, model.AccountBlock]
	client        *mongo.Client
	dbName        string
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewAccountBlockRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.AccountBlockRepository {
	return &AccountBlockStorage{
		accountBlock:  dal.NewMongoDal[model.AccountBlock, model.AccountBlock](client, cfg, dbName, collection),
		client:        client,
		dbName:        dbName,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (a *AccountBlockStorage) GetBranchByIds(ctx context.Context, id string) (*model.AccountBlock, error) {
	a.logger.Infof("[GetBranchByIds] fetching branch by id: %s", id)

	collection := a.client.Database(a.dbName).Collection("account_block")

	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": obj}
	filter["type"] = "B"

	branch, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, 1, a.logger)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[GetBranchByIds] branch not found")
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		a.logger.Errorf("[GetBranchByIds] failed to fetch branch: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(branch) == 0 {
		a.logger.Errorf("[GetBranchByIds] branch not found")
		return nil, errors.New(localization.ErrorBranchNotFound.Code)
	}

	a.logger.Infof("[GetBranchByIds] branch retrieved successfully")
	return branch[0], nil
}

func (a *AccountBlockStorage) CreateBranch(ctx context.Context, branch *model.AccountBlock) error {
	a.logger.Infof("[CreateBranch] creating branch")
	if branch.ID.IsZero() {
		branch.ID = bson.NewObjectID()
	}

	now := time.Now()
	branch.CreatedAt = now
	branch.UpdatedAt = now

	newBranch, err := a.accountBlock.InsertOne(ctx, *branch)
	if err != nil {
		a.logger.Errorf("[CreateBranch] failed to create branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[CreateBranch] branch created successfully")

	a.kafkaProducer.PublishMessage(ctx, newBranch, string(constants.ClientOrchestrationAccountBlockTopic), string(constants.ClientOrchestrationAccountBlockTopic), "create new branch")
	return nil
}

func (a *AccountBlockStorage) DeleteBranch(ctx context.Context, id string) error {
	a.logger.Infof("[DeleteBranch] deleting branch for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[DeleteBranch] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.accountBlock.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("[DeleteBranch] failed to delete branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[DeleteBranch] branch deleted successfully")
	return nil
}

func (a *AccountBlockStorage) FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"name", "code", "address", "is_enabled", "is_deleted", "city_id", "district_id", "region_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": bson.M{"$regex": "^" + filterParam.Search + "$", "$options": "i"}},
			{"code": searchRegex},
			// {"address": searchRegex},
			{"city_id": searchRegex},
			{"district_id": searchRegex},
			{"region_id": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["type"] = "B"
	filter["is_deleted"] = false

	ApplyIDFilter(filter, filterParam.Filters, "city_id")
	ApplyIDFilter(filter, filterParam.Filters, "district_id")
	ApplyIDFilter(filter, filterParam.Filters, "region_id")

	collection := a.client.Database(a.dbName).Collection("account_block")
	results, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, skip, limit, a.logger)
	if err != nil {
		a.logger.Errorf("[FindAllBranchesWithPagination] failed to find branches: %v", err)
		return nil, err
	}

	total, err := a.accountBlock.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("[FindAllBranchesWithPagination] failed to count branches: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	a.logger.Infof("[FindAllBranchesWithPagination] retrieved %d branches", len(results))

	return &types.PaginatedResponse[[]*model.AccountBlock]{
		Data: results,
		Meta: meta,
	}, nil
}

func (a *AccountBlockStorage) EnableOrDisableBranches(ctx context.Context, ids []string, reason string, enabled bool) error {
	var objIDs []bson.ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		objIDs = append(objIDs, objID)
	}

	collection := a.client.Database(a.dbName).Collection("account_block")

	filter := bson.M{"_id": bson.M{"$in": objIDs}, "type": "B"}
	update := bson.M{"$set": bson.M{"is_enabled": enabled, "updated_at": time.Now()}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	descendantsFilter := bson.M{"region_id": bson.M{"$in": objIDs}}
	_, err = collection.UpdateMany(ctx, descendantsFilter, update)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	var updatedResponse []*model.AccountBlock
	for _, id := range objIDs {
		updatedResponse = append(updatedResponse, &model.AccountBlock{
			ID:        id,
			Type:      "B",
			IsEnabled: enabled,
		})
	}

	a.kafkaProducer.PublishMessage(
		ctx,
		updatedResponse,
		string(constants.ClientOrchestrationServicesTopic),
		"account-block-updated",
		"account block enable status updated",
	)

	return nil
}

func (a *AccountBlockStorage) GetBranchesByIds(ctx context.Context, ids []string) ([]*model.AccountBlock, error) {
	var objIDs []bson.ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		objIDs = append(objIDs, objID)
	}

	collection := a.client.Database(a.dbName).Collection("account_block")
	filter := bson.M{"_id": bson.M{"$in": objIDs}, "type": "B"}

	block, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, int64(len(ids)), a.logger)
	if err != nil {
		return nil, err
	}
	if len(block) == 0 {
		return nil, errors.New(localization.ErrorBranchNotFound.Code)
	}
	return block, nil
}

func (a *AccountBlockStorage) CreateRegion(ctx context.Context, region *model.AccountBlock) error {
	if region.ID.IsZero() {
		region.ID = bson.NewObjectID()
	}
	now := time.Now()
	region.CreatedAt = now
	region.UpdatedAt = now

	newRegion, err := a.accountBlock.InsertOne(ctx, *region)
	if err != nil {
		a.logger.Errorf("Error creating region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.kafkaProducer.PublishMessage(ctx, newRegion, string(constants.ClientOrchestrationAccountBlockTopic), string(constants.ClientOrchestrationAccountBlockTopic), "create new region")

	return nil
}

func (a *AccountBlockStorage) DeleteRegion(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	filter["type"] = "R"

	err = a.accountBlock.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	allowedKeys := []string{"name", "code", "address", "is_enabled", "is_deleted"}

	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
			{"address": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["type"] = "R"

	collection := a.client.Database(a.dbName).Collection("account_block")
	results, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, skip, limit, a.logger)
	if err != nil {
		a.logger.Errorf("Error finding cities with pagination: %v", err)
		return nil, err
	}
	total, err := a.accountBlock.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.AccountBlock]{
		Data: results,
		Meta: meta,
	}, nil
}

func (a *AccountBlockStorage) EnableOrDisableRegions(ctx context.Context, ids []string, reason string, enabled bool) error {
	var objIDs []bson.ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		objIDs = append(objIDs, objID)
	}

	collection := a.client.Database(a.dbName).Collection("account_block")

	filter := bson.M{"_id": bson.M{"$in": objIDs}, "type": "R"}
	update := bson.M{"$set": bson.M{"is_enabled": enabled, "updated_at": time.Now()}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	descendantsFilter := bson.M{"region_id": bson.M{"$in": objIDs}}
	_, err = collection.UpdateMany(ctx, descendantsFilter, update)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	var updatedResponse []*model.AccountBlock
	for _, id := range objIDs {
		updatedResponse = append(updatedResponse, &model.AccountBlock{
			ID:        id,
			Type:      "R",
			IsEnabled: enabled,
		})
	}

	a.kafkaProducer.PublishMessage(
		ctx,
		updatedResponse,
		string(constants.ClientOrchestrationServicesTopic),
		"account-block-updated",
		"account block enable status updated",
	)

	return nil
}

func (a *AccountBlockStorage) GetRegionsByIds(ctx context.Context, ids []string) ([]*model.AccountBlock, error) {
	var objIDs []bson.ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		objIDs = append(objIDs, objID)
	}

	collection := a.client.Database(a.dbName).Collection("account_block")
	filter := bson.M{"_id": bson.M{"$in": objIDs}, "type": "R"}

	block, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, int64(len(ids)), a.logger)
	if err != nil {
		return nil, err
	}
	if len(block) == 0 {
		return nil, errors.New(localization.ErrorRegionNotFound.Code)
	}
	return block, nil

}

func (a *AccountBlockStorage) CreateDistrict(ctx context.Context, district *model.AccountBlock) error {
	if district.ID.IsZero() {
		district.ID = bson.NewObjectID()
	}
	now := time.Now()
	district.CreatedAt = now
	district.UpdatedAt = now

	newDistrict, err := a.accountBlock.InsertOne(ctx, *district)
	if err != nil {
		a.logger.Errorf("Error creating district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.kafkaProducer.PublishMessage(ctx, newDistrict, string(constants.ClientOrchestrationAccountBlockTopic), string(constants.ClientOrchestrationAccountBlockTopic), "create new district")

	return nil
}

func (a *AccountBlockStorage) DeleteDistrict(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.accountBlock.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	allowedKeys := []string{"region_id", "name", "code", "address", "is_enabled", "is_deleted"}

	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
			{"address": searchRegex},
			{"region_id": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["type"] = "D"

	if regionId, ok := filterParam.Filters["region_id"].(string); ok && regionId != "" {
		objID, err := bson.ObjectIDFromHex(regionId)
		if err == nil {
			filter["region_id"] = objID
		}
	}

	collection := a.client.Database(a.dbName).Collection("account_block")
	results, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, skip, limit, a.logger)
	if err != nil {
		a.logger.Errorf("Error finding cities with pagination: %v", err)
		return nil, err
	}

	total, err := a.accountBlock.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.AccountBlock]{
		Data: results,
		Meta: meta,
	}, nil
}

func (a *AccountBlockStorage) EnableOrDisableDistricts(ctx context.Context, ids []string, reason string, enabled bool) error {
	var objIDs []bson.ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		objIDs = append(objIDs, objID)
	}

	collection := a.client.Database(a.dbName).Collection("account_block")

	filter := bson.M{"_id": bson.M{"$in": objIDs}, "type": "D"}
	update := bson.M{"$set": bson.M{"is_enabled": enabled, "updated_at": time.Now()}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	descendantsFilter := bson.M{"district_id": bson.M{"$in": objIDs}}
	_, err = collection.UpdateMany(ctx, descendantsFilter, update)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	var updatedResponse []*model.AccountBlock
	for _, id := range objIDs {
		updatedResponse = append(updatedResponse, &model.AccountBlock{
			ID:        id,
			Type:      "D",
			IsEnabled: enabled,
		})
	}

	a.kafkaProducer.PublishMessage(
		ctx,
		updatedResponse,
		string(constants.ClientOrchestrationServicesTopic),
		"account-block-updated",
		"account block enable status updated",
	)

	return nil
}

func (a *AccountBlockStorage) GetDistrictsByIds(ctx context.Context, ids []string) ([]*model.AccountBlock, error) {
	var objIDs []bson.ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		objIDs = append(objIDs, objID)
	}

	collection := a.client.Database(a.dbName).Collection("account_block")
	filter := bson.M{"_id": bson.M{"$in": objIDs}, "type": "D"}

	block, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, int64(len(ids)), a.logger)
	if err != nil {
		return nil, err
	}
	if len(block) == 0 {
		return nil, errors.New(localization.ErrorDistrictNotFound.Code)
	}
	return block, nil
}

func (a *AccountBlockStorage) CreateCity(ctx context.Context, city *model.AccountBlock) error {
	if city.ID.IsZero() {
		city.ID = bson.NewObjectID()
	}

	now := time.Now()
	city.CreatedAt = now
	city.UpdatedAt = now

	newCity, err := a.accountBlock.InsertOne(ctx, *city)
	if err != nil {
		a.logger.Errorf("Error creating city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	a.kafkaProducer.PublishMessage(ctx, newCity, string(constants.ClientOrchestrationAccountBlockTopic), string(constants.ClientOrchestrationAccountBlockTopic), "create new city")

	return nil
}

func (a *AccountBlockStorage) DeleteCity(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.accountBlock.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindAllCitiesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"name", "code", "district_id", "address", "is_enabled", "is_deleted"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
			{"district_id": searchRegex},
			{"address": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["type"] = "C"

	if districtID, ok := filterParam.Filters["district_id"].(string); ok && districtID != "" {
		objID, err := bson.ObjectIDFromHex(districtID)
		if err == nil {
			filter["district_id"] = objID
		}
	}

	collection := a.client.Database(a.dbName).Collection("account_block")
	results, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, skip, limit, a.logger)
	if err != nil {
		a.logger.Errorf("Error finding cities with pagination: %v", err)
		return nil, err
	}

	totalCount, err := a.accountBlock.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error counting cities: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(totalCount, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.AccountBlock]{
		Data: results,
		Meta: meta,
	}, nil
}

func (a *AccountBlockStorage) EnableOrDisableCities(ctx context.Context, ids []string, reason string, enabled bool) error {
	var objIDs []bson.ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		objIDs = append(objIDs, objID)
	}

	collection := a.client.Database(a.dbName).Collection("account_block")

	filter := bson.M{"_id": bson.M{"$in": objIDs}, "type": "C"}
	update := bson.M{"$set": bson.M{"is_enabled": enabled, "updated_at": time.Now()}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	descendantsFilter := bson.M{"city_id": bson.M{"$in": objIDs}}
	_, err = collection.UpdateMany(ctx, descendantsFilter, update)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	var updatedResponse []*model.AccountBlock
	for _, id := range objIDs {
		updatedResponse = append(updatedResponse, &model.AccountBlock{
			ID:        id,
			Type:      "C",
			IsEnabled: enabled,
		})
	}

	a.kafkaProducer.PublishMessage(
		ctx,
		updatedResponse,
		string(constants.ClientOrchestrationServicesTopic),
		"account-block-updated",
		"account block enable status updated",
	)

	return nil
}

func (a *AccountBlockStorage) GetCitiesByIds(ctx context.Context, ids []string) ([]*model.AccountBlock, error) {
	var objIDs []bson.ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		objIDs = append(objIDs, objID)
	}

	collection := a.client.Database(a.dbName).Collection("account_block")
	filter := bson.M{"_id": bson.M{"$in": objIDs}, "type": "C"}

	block, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, int64(len(ids)), a.logger)
	if err != nil {
		return nil, err
	}
	if len(block) == 0 {
		return nil, errors.New(localization.ErrorCityNotFound.Code)
	}
	return block, nil

}
