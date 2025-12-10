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
	a.logger.Infof("[GetBranchByCode] fetching branch by code: %s", branchCode)
	collection := a.client.Database(a.dbName).Collection("account_block")
	// result, err := FindAccountBlockByCodeWithParentPopulated(ctx, collection, branchCode, "B", a.logger)
	filter := bson.M{"code": branchCode}
	filter["type"] = "B"
	result, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, 1, a.logger)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[GetBranchByCode] branch not found")
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		a.logger.Errorf("[GetBranchByCode] failed to fetch branch: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[GetBranchByCode] branch retrieved successfully")
	return result[0], nil
}

func (a *AccountBlockStorage) GetBranchByIds(ctx context.Context, id string) (*model.AccountBlock, error) {
	a.logger.Infof("[GetBranchByIds] fetching branch by id: %s", id)
	collection := a.client.Database(a.dbName).Collection("account_block")
	filter := bson.M{"_id": id}
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

	_, err := a.accountBlock.InsertOne(ctx, *branch)
	if err != nil {
		a.logger.Errorf("[CreateBranch] failed to create branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[CreateBranch] branch created successfully")
	return nil
}

func (a *AccountBlockStorage) UpdateBranch(ctx context.Context, id string, branch *model.AccountBlock) error {
	a.logger.Infof("[UpdateBranch] updating branch for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[UpdateBranch] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := AccountBlockMapperForUpdate(*branch)

	_, err = a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("[UpdateBranch] failed to update branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[UpdateBranch] branch updated successfully")
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

func (a *AccountBlockStorage) EnableOrDisableBranch(ctx context.Context, code string, reason string, enabled bool) error {
	a.logger.Infof("[EnableOrDisableBranch] processing branch enable/disable for code: %s, enabled: %v", code, enabled)
	filter := bson.M{"code": code}
	update := bson.M{"is_enabled": enabled, "updated_at": time.Now()}

	_, err := a.accountBlock.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("[EnableOrDisableBranch] failed to enable/disable branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[EnableOrDisableBranch] branch enable/disable completed successfully")
	return nil
}

func (a *AccountBlockStorage) FindBranchByID(ctx context.Context, id string) (*model.AccountBlock, error) {
	a.logger.Infof("[FindBranchByID] fetching branch by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[FindBranchByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	branch, err := a.accountBlock.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[FindBranchByID] branch not found")
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		a.logger.Errorf("[FindBranchByID] failed to find branch: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[FindBranchByID] branch retrieved successfully")
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

	// Override ID fields to exact match (not regex) when they're in Filters
	// This is needed because BuildMongoFilterWithKeys converts strings to regex
	if filterParam.Filters != nil {
		if districtID, ok := filterParam.Filters["district_id"].(string); ok && districtID != "" {
			filter["district_id"] = districtID // Exact match, not regex
		}
		if cityID, ok := filterParam.Filters["city_id"].(string); ok && cityID != "" {
			filter["city_id"] = cityID // Exact match, not regex
		}
		if regionID, ok := filterParam.Filters["region_id"].(string); ok && regionID != "" {
			filter["region_id"] = regionID // Exact match, not regex
		}
	}

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

func (a *AccountBlockStorage) GetRegionByIds(ctx context.Context, id string) (*model.AccountBlock, error) {
	// Implement this later when mongodb id is object_id
	// objID, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }

	filter := bson.M{"_id": id}

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

func (a *AccountBlockStorage) EnableOrDisableRegion(ctx context.Context, id string, reason string, enabled bool) error {
	// Get the branch itsef enabled/disable
	_, err := a.accountBlock.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"is_enabled": enabled, "updated_at": time.Now()})
	if err != nil {
		return err
	}

	// Get the districts within that region id
	districtFilter := types.Filter{
		Filters: map[string]interface{}{"region_id": id},
		Page:    1,
		PerPage: 1000,
	}

	districts, err := a.FindAllDistrictsWithPagination(ctx, districtFilter)
	if err != nil {
		return err
	}
	// Iterate over those districts and get all districts enabled/disabled
	for _, district := range districts.Data {
		filter := bson.M{"_id": district.ID.Hex()}
		update := bson.M{"is_enabled": enabled, "updated_at": time.Now()}

		_, err := a.accountBlock.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("Error enabling/disabling district: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		// That one district has many branches, enable/disable them all too
		branchFilter := types.Filter{
			Filters: map[string]interface{}{"district_id": district.ID.Hex()},
			Page:    1,
			PerPage: 1000,
		}
		branches, err := a.FindAllBranchesWithPagination(ctx, branchFilter)
		if err != nil {
			return err
		}

		for _, branch := range branches.Data {
			filter := bson.M{"_id": branch.ID.Hex()}
			branchUpdate := bson.M{"is_enabled": enabled, "updated_at": time.Now()}

			_, err := a.accountBlock.UpdateOne(ctx, filter, branchUpdate)
			if err != nil {
				a.logger.Errorf("Error enabling/disabling branches: %v", err)
				return errors.New(localization.ErrorUnexpectedError.Code)
			}
		}
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

func (a *AccountBlockStorage) GetDistrictById(ctx context.Context, id string) (*model.AccountBlock, error) {
	// collection := a.client.Database(a.dbName).Collection("account_block")

	// obj, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get object id: %v", err)
	// }
	filter := bson.M{"_id": id}
	filter["type"] = "D"
	// districts, err := FindAccountBlocksWithParentPopulatedRecursive(ctx, collection, filter, 0, 1, a.logger)
	district, err := a.accountBlock.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorDistrictNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return district, nil
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

func (a *AccountBlockStorage) EnableOrDisableDistrict(ctx context.Context, id string, reason string, enabled bool) error {
	filterParam := types.Filter{
		Filters: map[string]interface{}{"district_id": id},
		Page:    1,
		PerPage: 1000,
	}
	branches, err := a.FindAllBranchesWithPagination(ctx, filterParam)
	if err != nil {
		return err
	}

	for _, branch := range branches.Data {
		filter := bson.M{"_id": branch.ID.Hex()}
		update := bson.M{"is_enabled": enabled, "updated_at": time.Now()}

		_, err := a.accountBlock.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("Error enabling/disabling branches: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	// Needs to be changed to object id when real mongo data is used
	districtUpdate := bson.M{"is_enabled": enabled, "updated_at": time.Now()}
	_, err = a.accountBlock.UpdateOne(ctx, bson.M{"_id": id}, districtUpdate)
	if err != nil {
		return err
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
