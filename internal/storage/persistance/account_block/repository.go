package account_block

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AccountBlockStorage struct {
	branchDal   dal.MongoDal[model.Branch, model.Branch]
	regionDal   dal.MongoDal[model.Region, model.Region]
	districtDal dal.MongoDal[model.District, model.District]
	cityDal     dal.MongoDal[model.City, model.City]
	client      *mongo.Client
	dbName      string
	logger      utils.Logger
}

func NewAccountBlockRepository(client *mongo.Client, dbName string, logger utils.Logger) storage.AccountBlockRepository {
	return &AccountBlockStorage{
		branchDal:   dal.NewMongoDal[model.Branch, model.Branch](client, dbName, "branches"),
		regionDal:   dal.NewMongoDal[model.Region, model.Region](client, dbName, "regions"),
		districtDal: dal.NewMongoDal[model.District, model.District](client, dbName, "districts"),
		cityDal:     dal.NewMongoDal[model.City, model.City](client, dbName, "cities"),
		client:      client,
		dbName:      dbName,
		logger:      logger,
	}
}

// Standard CRUD operations for Branch

func (a *AccountBlockStorage) GetBranchByCode(ctx context.Context, branchCode string) (*model.Branch, error) {
	filter := bson.M{"branch_code": branchCode}

	result, err := a.branchDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result, nil
}

func (a *AccountBlockStorage) CreateBranch(ctx context.Context, branch *model.Branch) error {
	if branch.ID.IsZero() {
		branch.ID = bson.NewObjectID()
	}

	now := time.Now()
	branch.CreatedAt = now
	branch.UpdatedAt = now

	_, err := a.branchDal.InsertOne(ctx, *branch)
	if err != nil {
		a.logger.Errorf("Error creating branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) UpdateBranch(ctx context.Context, id string, branch *model.Branch) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := BranchMapperForUpdate(*branch)

	_, err = a.branchDal.UpdateOne(ctx, filter, update)
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

	err = a.branchDal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableBranch(ctx context.Context, code string, reason string, enabled bool) error {
	filter := bson.M{"branch_code": code}
	update := bson.M{"enabled": enabled, "enable_or_disable_reason": reason, "updated_at": time.Now()}

	_, err := a.branchDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling/disabling branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindBranchByID(ctx context.Context, id string) (*model.Branch, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	branch, err := a.branchDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		a.logger.Errorf("Error finding branch by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return branch, nil
}

func (a *AccountBlockStorage) FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Branch], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"branch_code", "branch_name", "branch_address", "region_name", "district_name", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"branch_code": searchRegex},
			{"branch_name": searchRegex},
			{"branch_address": searchRegex},
			{"region_name": searchRegex},
			{"district_name": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	// Always exclude deleted branches
	filter["is_deleted"] = false

	data, err := a.branchDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	total, err := a.branchDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.Branch]{
		Data: data,
		Meta: meta,
	}, nil
}

// Standard CRUD operations for Region

func (a *AccountBlockStorage) GetRegionByCode(ctx context.Context, regionCode string) (*model.Region, error) {
	filter := bson.M{"region_code": regionCode}
	result, err := a.regionDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorRegionNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result, nil
}

func (a *AccountBlockStorage) CreateRegion(ctx context.Context, region *model.Region) error {
	if region.ID.IsZero() {
		region.ID = bson.NewObjectID()
	}
	now := time.Now()
	region.CreatedAt = now
	region.UpdatedAt = now

	_, err := a.regionDal.InsertOne(ctx, *region)
	if err != nil {
		a.logger.Errorf("Error creating region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) UpdateRegion(ctx context.Context, id string, region *model.Region) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := RegionMapperForUpdate(*region)

	_, err = a.regionDal.UpdateOne(ctx, filter, update)
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

	err = a.regionDal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableRegion(ctx context.Context, code string, reason string, enabled bool) error {
	filter := bson.M{"region_code": code}
	update := bson.M{"enabled": enabled, "enable_or_disable_reason": reason, "updated_at": time.Now()}

	_, err := a.regionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling/disabling region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindRegionByID(ctx context.Context, id string) (*model.Region, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	region, err := a.regionDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorRegionNotFound.Code)
		}
		a.logger.Errorf("Error finding region by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return region, nil
}

func (a *AccountBlockStorage) FindAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Region], error) {
	allowedKeys := []string{"region_name", "region_code", "enabled"}

	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"region_name": searchRegex},
			{"region_code": searchRegex},
		}
	}

	// Build filter + pagination
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	fmt.Println("Filter constructed:", filter)
	fmt.Println("Skip:", skip, "Limit:", limit)

	// Fetch paginated data
	results, err := a.regionDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		a.logger.Errorf("Error finding regions with pagination: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := a.regionDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// Return paginated response (always return results, even if Search is empty)
	return &types.PaginatedResponse[[]*model.Region]{
		Data: results,
		Meta: meta,
	}, nil
}

// Standard CRUD operations for District

func (a *AccountBlockStorage) GetDistrictByCode(ctx context.Context, districtCode string) (*model.District, error) {
	filter := bson.M{"district_code": districtCode}
	result, err := a.districtDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorDistrictNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result, nil
}

func (a *AccountBlockStorage) CreateDistrict(ctx context.Context, district *model.District) error {
	if district.ID.IsZero() {
		district.ID = bson.NewObjectID()
	}
	now := time.Now()
	district.CreatedAt = now
	district.UpdatedAt = now

	_, err := a.districtDal.InsertOne(ctx, *district)
	if err != nil {
		a.logger.Errorf("Error creating district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) UpdateDistrict(ctx context.Context, id string, district *model.District) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := DistrictMapperForUpdate(*district)

	_, err = a.districtDal.UpdateOne(ctx, filter, update)
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

	err = a.districtDal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableDistrict(ctx context.Context, code string, reason string, enabled bool) error {
	filter := bson.M{"district_code": code}
	update := bson.M{"enabled": enabled, "enable_or_disable_reason": reason, "updated_at": time.Now()}

	_, err := a.districtDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling/disabling district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindDistrictByID(ctx context.Context, id string) (*model.District, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	district, err := a.districtDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorDistrictNotFound.Code)
		}
		a.logger.Errorf("Error finding district by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return district, nil
}

func (a *AccountBlockStorage) FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.District], error) {
	allowedKeys := []string{"district_name", "district_code", "enabled", "region_name"}

	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"district_name": searchRegex},
			{"district_code": searchRegex},
		}
	}

	// Build filter + pagination using FilterBuilder
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// Fetch paginated results
	results, err := a.districtDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		a.logger.Errorf("Error finding districts with pagination: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := a.regionDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// Return proper paginated response
	return &types.PaginatedResponse[[]*model.District]{
		Data: results,
		Meta: meta,
	}, nil
}

// Standard CRUD operations for City

func (a *AccountBlockStorage) GetCityByCode(ctx context.Context, cityCode string) (*model.City, error) {

	filter := bson.M{"city_code": cityCode}
	result, err := a.cityDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorCityNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result, nil
}

func (a *AccountBlockStorage) CreateCity(ctx context.Context, city *model.City) error {
	if city.ID.IsZero() {
		city.ID = bson.NewObjectID()
	}

	now := time.Now()
	city.CreatedAt = now
	city.UpdatedAt = now

	_, err := a.cityDal.InsertOne(ctx, *city)
	if err != nil {
		a.logger.Errorf("Error creating city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) UpdateCity(ctx context.Context, id string, city *model.City) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := CityMapperForUpdate(*city)

	_, err = a.cityDal.UpdateOne(ctx, filter, update)
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

	err = a.cityDal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableCity(ctx context.Context, code string, reason string, enabled bool) error {
	filter := bson.M{"city_code": code}
	update := bson.M{"enabled": enabled, "enable_or_disable_reason": reason, "updated_at": time.Now()}

	_, err := a.cityDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling/disabling city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) FindCityByID(ctx context.Context, id string) (*model.City, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	city, err := a.cityDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorCityNotFound.Code)
		}
		a.logger.Errorf("Error finding city by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return city, nil
}

func (a *AccountBlockStorage) FindAllCitiesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.City], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"city_address", "city_name", "region_name", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"city_name": searchRegex},
			{"city_code": searchRegex},
			{"city_address": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// Get total count
	totalCount, err := a.cityDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error counting cities: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	results, err := a.cityDal.FindAllWithPagination(ctx, filter, nil, skip, limit)
	if err != nil {
		a.logger.Errorf("Error finding cities with pagination: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(totalCount, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.City]{
		Data: results,
		Meta: meta,
	}, nil
}
