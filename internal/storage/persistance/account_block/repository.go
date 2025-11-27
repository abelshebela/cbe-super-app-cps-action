package account_block

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AccountBlockStorage struct {
	accountBlock dal.MongoDal[model.AccountBlock, model.AccountBlock]
	client       *mongo.Client
	dbName       string
	logger       utils.Logger
}

func NewAccountBlockRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.AccountBlockRepository {
	return &AccountBlockStorage{
		accountBlock: dal.NewMongoDal[model.AccountBlock, model.AccountBlock](client, dbName, collection),
		client:       client,
		dbName:       dbName,
		logger:       logger,
	}
}

// Standard CRUD operations for Branch

func (a *AccountBlockStorage) GetBranchByCode(ctx context.Context, branchCode string) (*model.AccountBlock, error) {
	collection := a.client.Database(a.dbName).Collection("account_block")
	// result, err := FindAccountBlockByCodeWithParentPopulated(ctx, collection, branchCode, "B", a.logger)
	filter := bson.M{"code": branchCode}
	filter["type"] = "B"
	result, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, 1, a.logger)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result[0], nil
}

func (a *AccountBlockStorage) CreateBranch(ctx context.Context, branch *model.AccountBlock) error {
	if branch.ID.IsZero() {
		branch.ID = bson.NewObjectID()
	}

	now := time.Now()
	branch.CreatedAt = now
	branch.UpdatedAt = now

	_, err := a.accountBlock.InsertOne(ctx, *branch)
	if err != nil {
		a.logger.Errorf("Error creating branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) UpdateBranch(ctx context.Context, id string, branch *model.AccountBlock) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := AccountBlockMapperForUpdate(*branch)

	_, err = a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) DeleteBranch(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.accountBlock.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableBranch(ctx context.Context, code string, reason string, enabled bool) error {
	filter := bson.M{"code": code}
	update := bson.M{"is_enabled": enabled, "updated_at": time.Now()}

	_, err := a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling/disabling branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindBranchByID(ctx context.Context, id string) (*model.AccountBlock, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	branch, err := a.accountBlock.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		a.logger.Errorf("Error finding branch by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return branch, nil
}

func (a *AccountBlockStorage) FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"name", "code", "address", "is_enabled", "is_deleted", "city_id", "district_id", "region_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
			{"address": searchRegex},
			{"city_id": searchRegex},
			{"district_id": searchRegex},
			{"region_id": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["type"] = "B"
	filter["is_deleted"] = false

	collection := a.client.Database(a.dbName).Collection("account_block")
	results, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, skip, limit, a.logger)
	if err != nil {
		a.logger.Errorf("Error finding cities with pagination: %v", err)
		return nil, err
	}

	total, err := a.accountBlock.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.AccountBlock]{
		Data: results,
		Meta: meta,
	}, nil
}

func (a *AccountBlockStorage) GetRegionByCode(ctx context.Context, regionCode string) (*model.AccountBlock, error) {
	collection := a.client.Database(a.dbName).Collection("account_block")
	filter := bson.M{"code": regionCode}
	filter["type"] = "R"
	result, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, 1, a.logger)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorRegionNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result[0], nil
}

func (a *AccountBlockStorage) CreateRegion(ctx context.Context, region *model.AccountBlock) error {
	if region.ID.IsZero() {
		region.ID = bson.NewObjectID()
	}
	now := time.Now()
	region.CreatedAt = now
	region.UpdatedAt = now

	_, err := a.accountBlock.InsertOne(ctx, *region)
	if err != nil {
		a.logger.Errorf("Error creating region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) UpdateRegion(ctx context.Context, id string, region *model.AccountBlock) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := AccountBlockMapperForUpdate(*region)

	_, err = a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) DeleteRegion(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.accountBlock.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableRegion(ctx context.Context, code string, reason string, enabled bool) error {
	filter := bson.M{"code": code}
	update := bson.M{"is_enabled": enabled, "updated_at": time.Now()}

	_, err := a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling/disabling region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindRegionByID(ctx context.Context, id string) (*model.AccountBlock, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	region, err := a.accountBlock.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorRegionNotFound.Code)
		}
		a.logger.Errorf("Error finding region by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return region, nil
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

	// Build filter + pagination
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["type"] = "R"

	// Fetch paginated data
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

	// Return paginated response (always return results, even if Search is empty)
	return &types.PaginatedResponse[[]*model.AccountBlock]{
		Data: results,
		Meta: meta,
	}, nil
}

// Standard CRUD operations for District

func (a *AccountBlockStorage) GetDistrictByCode(ctx context.Context, districtCode string) (*model.AccountBlock, error) {
	collection := a.client.Database(a.dbName).Collection("account_block")
	filter := bson.M{"code": districtCode}
	filter["type"] = "D"
	result, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, 1, a.logger)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorDistrictNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result[0], nil
}

func (a *AccountBlockStorage) CreateDistrict(ctx context.Context, district *model.AccountBlock) error {
	if district.ID.IsZero() {
		district.ID = bson.NewObjectID()
	}
	now := time.Now()
	district.CreatedAt = now
	district.UpdatedAt = now

	_, err := a.accountBlock.InsertOne(ctx, *district)
	if err != nil {
		a.logger.Errorf("Error creating district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) UpdateDistrict(ctx context.Context, id string, district *model.AccountBlock) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := AccountBlockMapperForUpdate(*district)

	_, err = a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

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

func (a *AccountBlockStorage) EnableOrDisableDistrict(ctx context.Context, code string, reason string, enabled bool) error {
	filter := bson.M{"code": code}
	update := bson.M{"is_enabled": enabled, "updated_at": time.Now()}

	_, err := a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling/disabling district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindDistrictByID(ctx context.Context, id string) (*model.AccountBlock, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	district, err := a.accountBlock.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorDistrictNotFound.Code)
		}
		a.logger.Errorf("Error finding district by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return district, nil
}

func (a *AccountBlockStorage) FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
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

	// Build filter + pagination using FilterBuilder
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["type"] = "D"

	// Fetch paginated results
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

	// Return proper paginated response
	return &types.PaginatedResponse[[]*model.AccountBlock]{
		Data: results,
		Meta: meta,
	}, nil
}

// Standard CRUD operations for City

func (a *AccountBlockStorage) GetCityByCode(ctx context.Context, cityCode string) (*model.AccountBlock, error) {
	collection := a.client.Database(a.dbName).Collection("account_block")
	filter := bson.M{"code": cityCode}
	filter["type"] = "C"
	result, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, 1, a.logger)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorCityNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result[0], nil
}

func (a *AccountBlockStorage) CreateCity(ctx context.Context, city *model.AccountBlock) error {
	if city.ID.IsZero() {
		city.ID = bson.NewObjectID()
	}

	now := time.Now()
	city.CreatedAt = now
	city.UpdatedAt = now

	_, err := a.accountBlock.InsertOne(ctx, *city)
	if err != nil {
		a.logger.Errorf("Error creating city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) UpdateCity(ctx context.Context, id string, city *model.AccountBlock) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := AccountBlockMapperForUpdate(*city)

	_, err = a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

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

func (a *AccountBlockStorage) EnableOrDisableCity(ctx context.Context, code string, reason string, enabled bool) error {
	filter := bson.M{"code": code}
	update := bson.M{"is_enabled": enabled, "updated_at": time.Now()}

	_, err := a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling/disabling city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindCityByID(ctx context.Context, id string) (*model.AccountBlock, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	city, err := a.accountBlock.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorCityNotFound.Code)
		}
		a.logger.Errorf("Error finding city by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return city, nil
}

func (a *AccountBlockStorage) FindAllCitiesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"name", "code", "address", "is_enabled", "is_deleted"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
			{"address": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["type"] = "C"

	collection := a.client.Database(a.dbName).Collection("account_block")
	results, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, skip, limit, a.logger)
	if err != nil {
		a.logger.Errorf("Error finding cities with pagination: %v", err)
		return nil, err
	}

	// Get total count
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
