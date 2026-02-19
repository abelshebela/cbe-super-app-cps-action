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

	account_block_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AccountBlockStorage struct {
	cfg           *config.VaultConfig
	accountBlock  dal.MongoDal[model.AccountBlock, model.AccountBlock]
	cpsActionRepo dal.MongoDal[model.CPSAction, model.CPSAction]
	client        *mongo.Client
	dbName        string
	abCollection  string
	cpsCollection string
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewAccountBlockRepository(client *mongo.Client, cfg *config.VaultConfig, dbName, abCollection, cpsCollection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.AccountBlockRepository {
	return &AccountBlockStorage{
		cfg:           cfg,
		accountBlock:  dal.NewMongoDal[model.AccountBlock, model.AccountBlock](client, cfg, dbName, abCollection),
		cpsActionRepo: dal.NewMongoDal[model.CPSAction, model.CPSAction](client, cfg, dbName, cpsCollection),
		client:        client,
		dbName:        dbName,
		abCollection:  abCollection,
		cpsCollection: cpsCollection,
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
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := a.accountBlock.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("[FindAllBranchesWithPagination] failed to count branches: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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
	update := bson.M{"$set": bson.M{"is_enabled": enabled, "disabled_reason": reason, "updated_at": time.Now()}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("[EnableOrDisableBranches] failed to update branches: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	descendantsFilter := bson.M{"region_id": bson.M{"$in": objIDs}}
	_, err = collection.UpdateMany(ctx, descendantsFilter, update)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	updatedBranches, err := a.GetBranchesByIds(ctx, ids)

	a.kafkaProducer.PublishMessage(
		ctx,
		updatedBranches,
		string(constants.ClientOrchestrationServicesTopic),
		a.cfg.AccountBlockUpdate,
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
		a.logger.Errorf("[GetBranchesByIds] failed to fetch branches: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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
		a.logger.Errorf("[FindAllRegionsWithPagination] failed to find regions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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
	update := bson.M{"$set": bson.M{"is_enabled": enabled, "disabled_reason": reason, "updated_at": time.Now()}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("[EnableOrDisableRegions] failed to update regions: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// descendantsFilter := bson.M{"region_id": bson.M{"$in": objIDs}}
	// _, err = collection.UpdateMany(ctx, descendantsFilter, update)
	// if err != nil {
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// regions, err := a.GetRegionsByIds(ctx, ids)
	// if err != nil {
	// 	a.logger.Errorf("Error fetching updated regions: %v", err)
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// a.kafkaProducer.PublishMessage(
	// 	ctx,
	// 	regions,
	// 	string(constants.ClientOrchestrationServicesTopic),
	// 	a.cfg.AccountBlockUpdate,
	// 	"account block enable status updated",
	// )

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
		a.logger.Errorf("[GetRegionsByIds] failed to fetch regions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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

	ApplyIDFilter(filter, filterParam.Filters, "region_id")

	collection := a.client.Database(a.dbName).Collection("account_block")
	results, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, skip, limit, a.logger)
	if err != nil {
		a.logger.Errorf("[FindAllDistrictsWithPagination] failed to find districts: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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
	update := bson.M{"$set": bson.M{"is_enabled": enabled, "disabled_reason": reason, "updated_at": time.Now()}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("[EnableOrDisableDistricts] failed to update districts: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// descendantsFilter := bson.M{"district_id": bson.M{"$in": objIDs}}
	// _, err = collection.UpdateMany(ctx, descendantsFilter, update)
	// if err != nil {
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// districts, err := a.GetDistrictsByIds(ctx, ids)
	// if err != nil {
	// 	a.logger.Errorf("Error fetching updated regions: %v", err)
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// a.kafkaProducer.PublishMessage(
	// 	ctx,
	// 	districts,
	// 	string(constants.ClientOrchestrationServicesTopic),
	// 	a.cfg.AccountBlockUpdate,
	// 	"account block enable status updated",
	// )

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
		a.logger.Errorf("[GetDistrictsByIds] failed to fetch districts: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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
		a.logger.Errorf("[FindAllCitiesWithPagination] failed to find cities: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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
	update := bson.M{"$set": bson.M{"is_enabled": enabled, "disabled_reason": reason, "updated_at": time.Now()}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("[EnableOrDisableCities] failed to update cities: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// descendantsFilter := bson.M{"city_id": bson.M{"$in": objIDs}}
	// _, err = collection.UpdateMany(ctx, descendantsFilter, update)
	// if err != nil {
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// cities, err := a.GetCitiesByIds(ctx, ids)
	// if err != nil {
	// 	a.logger.Errorf("Error fetching updated regions: %v", err)
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// a.kafkaProducer.PublishMessage(
	// 	ctx,
	// 	cities,
	// 	string(constants.ClientOrchestrationServicesTopic),
	// 	a.cfg.AccountBlockUpdate,
	// 	"account block enable status updated",
	// )

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
		a.logger.Errorf("[GetCitiesByIds] failed to fetch cities: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if len(block) == 0 {
		return nil, errors.New(localization.ErrorCityNotFound.Code)
	}
	return block, nil

}

func (a *AccountBlockStorage) GetAccountBlockDetails(ctx context.Context, id string) ([]account_block_dto.AccountBlockActionResponse, error) {
	a.logger.Infof("[GetAccountBlockDetails] fetching CPS actions for account block id: %s", id)

	if _, err := bson.ObjectIDFromHex(id); err != nil {
		a.logger.Errorf("[GetAccountBlockDetails] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	cpsCollection := a.client.Database(a.dbName).Collection(a.cpsCollection)

	accountBlockRequestActions := []string{
		string(constants.RequestEnableBranches),
		string(constants.RequestDisableBranches),
		string(constants.RequestEnableCities),
		string(constants.RequestDisableCities),
		string(constants.RequestEnableDistricts),
		string(constants.RequestDisableDistricts),
		string(constants.RequestEnableRegions),
		string(constants.RequestDisableRegions),
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"is_deleted":         false,
			"previous_action.id": id,
			"request_action":     bson.M{"$in": accountBlockRequestActions},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
	}

	cur, err := cpsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		a.logger.Errorf("[GetAccountBlockDetails] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		a.logger.Errorf("[GetAccountBlockDetails] failed to decode CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(results) == 0 {
		a.logger.Errorf("[GetAccountBlockDetails] no CPS actions found for id: %s", id)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	// Map CPSAction to AccountBlockActionResponse
	response := make([]account_block_dto.AccountBlockActionResponse, 0, len(results))
	for _, cpsAction := range results {
		actionResponse := account_block_dto.AccountBlockActionResponse{
			ID:                  cpsAction.ID,
			ActionCode:          cpsAction.ActionCode,
			UniqueId:            cpsAction.UniqueId,
			MakerID:             cpsAction.MakerID,
			MakerName:           cpsAction.MakerName,
			MakerPhoneNumber:    cpsAction.MakerPhoneNumber,
			CheckerUsers:        convertCheckers(cpsAction.CheckerUsers),
			AuditorUsers:        convertAuditors(cpsAction.AuditorUsers),
			AuditorCount:        cpsAction.AuditorCount,
			AuditorStatus:       account_block_dto.AuditorStatus(cpsAction.AuditorStatus),
			CurrentAuditorIndex: cpsAction.CurrentAuditorIndex,
			CheckerCount:        cpsAction.CheckerCount,
			CurrentCheckerIndex: cpsAction.CurrentCheckerIndex,
			RoleCode:            cpsAction.RoleCode,
			RejectionReason:     cpsAction.RejectionReason,
			CanceledReason:      cpsAction.CanceledReason,
			ActionStatus:        cpsAction.ActionStatus,
			ActionType:          cpsAction.ActionType,
			IsDeleted:           cpsAction.IsDeleted,
			RequestAction:       cpsAction.RequestAction,
			Version:             cpsAction.Version,
			ReversedByRoleID:    cpsAction.ReversedByRoleID,
			ReversedByID:        cpsAction.ReversedByID,
			ReversedByName:      cpsAction.ReversedByName,
			ReversedAt:          cpsAction.ReversedAt,
			CreatedAt:           cpsAction.CreatedAt,
			LastModifiedAt:      cpsAction.LastModifiedAt,
			MakerActionTime:     cpsAction.MakerActionTime,
		}

		// Map previous_action based on ActionStatus
		var previousAction interface{}
		if cpsAction.ActionStatus == string(constants.ActionApproved) {
			// If APPROVED, use current_action filtered by matching id
			previousAction = getMatchingAction(cpsAction.CurrentAction, id, a.logger)
		} else if cpsAction.ActionStatus == string(constants.ActionPending) || cpsAction.ActionStatus == string(constants.ActionRejected) {
			// If PENDING or REJECTED, use previous_action filtered by matching id
			previousAction = getMatchingAction(cpsAction.PreviousAction, id, a.logger)
		} else {
			// For other statuses, use previous_action as fallback
			previousAction = getMatchingAction(cpsAction.PreviousAction, id, a.logger)
		}

		actionResponse.PreviousAction = previousAction
		response = append(response, actionResponse)
	}

	a.logger.Infof("[GetAccountBlockDetails] successfully mapped %d CPS actions", len(response))
	return response, nil
}

// getMatchingAction extracts the action that matches the given account block id
// from either previous_action or current_action (which can be arrays or single objects)
func getMatchingAction(actionData interface{}, accountBlockID string, logger utils.Logger) interface{} {
	if actionData == nil {
		return nil
	}

	// Try to unmarshal as array of EnableDisableAction
	actions, err := local_util.JsonUnmarshal[[]types.EnableDisableAction](actionData)
	if err == nil && actions != nil {
		// Find the action matching the account block id
		for _, action := range *actions {
			if action.ID == accountBlockID {
				return action
			}
		}
		// If no match found, return nil
		return nil
	}

	// Try to unmarshal as single EnableDisableAction
	singleAction, err := local_util.JsonUnmarshal[types.EnableDisableAction](actionData)
	if err == nil && singleAction != nil {
		if singleAction.ID == accountBlockID {
			return *singleAction
		}
		return nil
	}

	// If unmarshaling fails, check if it's already a map/object with id field
	if actionMap, ok := actionData.(map[string]interface{}); ok {
		if id, exists := actionMap["id"]; exists {
			if idStr, ok := id.(string); ok && idStr == accountBlockID {
				return actionMap
			}
		}
	}

	logger.Warnf("[getMatchingAction] failed to extract matching action for id: %s", accountBlockID)
	return nil
}

// convertCheckers converts model.Checker to account_block_dto.Checker
func convertCheckers(checkers []model.Checker) []account_block_dto.Checker {
	result := make([]account_block_dto.Checker, 0, len(checkers))
	for _, c := range checkers {
		result = append(result, account_block_dto.Checker{
			CheckerID:          c.CheckerID,
			RoleID:             c.RoleID,
			CheckerIndex:       c.CheckerIndex,
			CheckerName:        c.CheckerName,
			CheckerPhoneNumber: c.CheckerPhoneNumber,
			ApprovedAt:         c.ApprovedAt,
		})
	}
	return result
}

// convertAuditors converts model.Auditor to account_block_dto.Auditor
func convertAuditors(auditors []model.Auditor) []account_block_dto.Auditor {
	result := make([]account_block_dto.Auditor, 0, len(auditors))
	for _, a := range auditors {
		result = append(result, account_block_dto.Auditor{
			AuditorID:          a.AuditorID,
			RoleID:             a.RoleID,
			AuditorIndex:       a.AuditorIndex,
			AuditorName:        a.AuditorName,
			AuditorPhoneNumber: a.AuditorPhoneNumber,
			AuditorReason:      a.AuditorReason,
			AuditorMark:        account_block_dto.AuditorMark(a.AuditorMark),
			ApprovedAt:         a.ApprovedAt,
		})
	}
	return result
}
