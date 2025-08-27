package account_block

import (
	"cbe-super-app-cps-action/internal/constants"
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
	branchDal   dal.MongoDal[model.Branch, model.Branch]
	actionDal   dal.MongoDal[model.CPSAction, model.CPSAction]
	cityDal     dal.MongoDal[model.City, model.City]
	regionDal   dal.MongoDal[model.Region, model.Region]
	districtDal dal.MongoDal[model.District, model.District]
	client      *mongo.Client
	dbName      string
	logger      utils.Logger
}

func NewAccountBlockRepository(client *mongo.Client, dbName string, branchCollection string, actionCollection string, logger utils.Logger) storage.AccountBlockRepository {
	return &AccountBlockStorage{
		branchDal:   dal.NewMongoDal[model.Branch, model.Branch](client, dbName, branchCollection),
		actionDal:   dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, actionCollection),
		cityDal:     dal.NewMongoDal[model.City, model.City](client, dbName, "cities"),
		regionDal:   dal.NewMongoDal[model.Region, model.Region](client, dbName, "regions"),
		districtDal: dal.NewMongoDal[model.District, model.District](client, dbName, "districts"),
		client:      client,
		dbName:      dbName,
		logger:      logger,
	}
}

func (a *AccountBlockStorage) GetBranchByCode(ctx context.Context, branchCode string) (*model.Branch, error) {
	filter := bson.M{"branch_code": branchCode, "is_deleted": false}

	result, err := a.branchDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result, nil
}

func (a *AccountBlockStorage) GetAllBranches(ctx context.Context, region, district string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Branch], error) {
	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage
	projection := bson.M{}

	filter := bson.M{
		"branch_region": bson.M{"$regex": region, "$options": "i"},
		"district_name": bson.M{"$regex": district, "$options": "i"},
	}

	branches, err := a.branchDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		a.logger.Errorf("failed to fetch branches: %v", err)
		return nil, errors.New("FAILED_TO_FETCH_BRANCHES")
	}

	if len(branches) == 0 {
		a.logger.Warnf("No branches found for region: %s, district: %s", region, district)
		return nil, errors.New("BRANCH_NOT_FOUND")
	}

	total, err := a.branchDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := local_util.BuildPaginationMeta(total, filterParams.Page, limit)

	return &types.PaginatedResponse[[]*model.Branch]{
		Data: branches,
		Meta: meta,
	}, nil
}

func (a *AccountBlockStorage) CreateBranch(ctx context.Context, branch *model.Branch) error {
	if branch.ID.IsZero() {
		branch.ID = bson.NewObjectID()
	}
	branch.CreatedAt = time.Now()
	branch.UpdatedAt = time.Now()

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
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := BranchMapper(*branch)

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
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.branchDal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting branch: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableBranch(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable, "updated_at": time.Now()}}

	_, err = a.branchDal.UpdateOne(ctx, filter, update)
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
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("Error finding branch by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return branch, nil
}

func (a *AccountBlockStorage) FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Branch], error) {
	searchKeys := bson.M{}
	filter := bson.M{"is_deleted": false}
	allowedKeys := []string{"branch_address", "district_name", "branch_region", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["category_name"] = searchRegex
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

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

// Standard CRUD operations for City

func (a *AccountBlockStorage) GetCityByCode(ctx context.Context, cityCode string) (*model.City, error) {

	filter := bson.M{"city_code": cityCode}
	cityDoc, err := a.cityDal.FindOne(ctx, filter, nil)

	if err != nil || cityDoc == nil {
		return &model.City{}, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return &model.City{
		ID:           cityDoc.ID,
		City:         cityDoc.City,
		CityCode:     cityDoc.CityCode,
		CityName:     cityDoc.CityName,
		DistrictID:   cityDoc.DistrictID,
		DistrictName: cityDoc.DistrictName,
		RegionID:     cityDoc.RegionID,
		RegionName:   cityDoc.RegionName,
		CreatedAt:    cityDoc.CreatedAt,
		UpdatedAt:    cityDoc.UpdatedAt,
		Enabled:      cityDoc.Enabled,
	}, nil
}

func (a *AccountBlockStorage) CreateCity(ctx context.Context, city *model.City) error {
	if city.ID.IsZero() {
		city.ID = bson.NewObjectID()
	}
	city.CreatedAt = time.Now()
	city.UpdatedAt = time.Now()

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
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := CityMapper(*city)

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
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.cityDal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting city: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableCity(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}

	_, err = a.cityDal.UpdateOne(ctx, filter, update)
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
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("Error finding city by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return city, nil
}

func (a *AccountBlockStorage) FindAllCitiesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.City], error) {
	searchKeys := bson.M{}
	filter := bson.M{"is_deleted": false}
	allowedKeys := []string{"city_address", "city_name", "city_region", "region_name", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"city_name": searchRegex},
			{"city_code": searchRegex},
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

	if filterParam.Search == "" {
		return &types.PaginatedResponse[[]*model.City]{
			Data: []*model.City{},
		}, nil
	}

	meta := local_util.BuildPaginationMeta(totalCount, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.City]{
		Data: results,
		Meta: meta,
	}, nil
}

// Standard CRUD operations for Region

func (a *AccountBlockStorage) GetRegionByCode(ctx context.Context, regionCode string) (*model.Region, error) {
	filter := bson.M{"region_code": regionCode}
	regionDoc, err := a.regionDal.FindOne(ctx, filter, nil)
	if err != nil || regionDoc == nil {
		return &model.Region{}, errors.New("REGION_NOT_FOUND")
	}
	return &model.Region{
		ID:            regionDoc.ID,
		RegionCode:    regionDoc.RegionCode,
		RegionName:    regionDoc.RegionName,
		RegionAddress: regionDoc.RegionAddress,
		CreatedAt:     regionDoc.CreatedAt,
		UpdatedAt:     regionDoc.UpdatedAt,
		Enabled:       regionDoc.Enabled,
	}, nil
}

func (a *AccountBlockStorage) CreateRegion(ctx context.Context, region *model.Region) error {
	if region.ID.IsZero() {
		region.ID = bson.NewObjectID()
	}
	region.CreatedAt = time.Now()
	region.UpdatedAt = time.Now()

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
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"region_code":    region.RegionCode,
		"region_name":    region.RegionName,
		"region_address": region.RegionAddress,
		"enabled":        region.Enabled,
		"updated_at":     time.Now(),
	}

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
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.regionDal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting region: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableRegion(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable, "updated_at": time.Now()}}

	_, err = a.regionDal.UpdateOne(ctx, filter, update)
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
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("Error finding region by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return region, nil
}

func (a *AccountBlockStorage) FindAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Region], error) {
	filter := bson.M{}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"region_name": searchRegex},
			{"region_code": searchRegex},
		}
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	// Get total count
	totalCount, err := a.regionDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error counting regions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	results, err := a.regionDal.FindAllWithPagination(ctx, filter, nil, skip, limit)
	if err != nil {
		a.logger.Errorf("Error finding regions with pagination: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// If filterParam.Search is empty, return empty data and meta using local_utils
	if filterParam.Search == "" {
		return &types.PaginatedResponse[[]*model.Region]{
			Data: []*model.Region{},
			Meta: types.PaginationMeta{TotalDocs: totalCount},
		}, nil
	}

	return &types.PaginatedResponse[[]*model.Region]{
		Data: results,
		Meta: types.PaginationMeta{TotalDocs: totalCount},
	}, nil
}

// Standard CRUD operations for District

func (a *AccountBlockStorage) GetDistrictByCode(ctx context.Context, districtCode string) (*model.District, error) {
	filter := bson.M{"district_code": districtCode}
	districtDoc, err := a.districtDal.FindOne(ctx, filter, nil)
	if err != nil || districtDoc == nil {
		return &model.District{}, errors.New("FAILED_TO_GET_DISTRICT")
	}
	return &model.District{
		ID:              districtDoc.ID,
		DistrictCode:    districtDoc.DistrictCode,
		DistrictName:    districtDoc.DistrictName,
		DistrictAddress: districtDoc.DistrictAddress,
		RegionID:        districtDoc.RegionID,
		RegionName:      districtDoc.RegionName,
		CreatedAt:       districtDoc.CreatedAt,
		UpdatedAt:       districtDoc.UpdatedAt,
		Enabled:         districtDoc.Enabled,
	}, nil
}

func (a *AccountBlockStorage) CreateDistrict(ctx context.Context, district *model.District) error {
	if district.ID.IsZero() {
		district.ID = bson.NewObjectID()
	}
	district.CreatedAt = time.Now()
	district.UpdatedAt = time.Now()

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
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"district_code":    district.DistrictCode,
		"district_name":    district.DistrictName,
		"district_address": district.DistrictAddress,
		"region_id":        district.RegionID,
		"region_name":      district.RegionName,
		"enabled":          district.Enabled,
		"updated_at":       time.Now(),
	}

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
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	err = a.districtDal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error deleting district: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (a *AccountBlockStorage) EnableOrDisableDistrict(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable, "updated_at": time.Now()}}

	_, err = a.districtDal.UpdateOne(ctx, filter, update)
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
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("Error finding district by ID: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return district, nil
}

func (a *AccountBlockStorage) FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.District], error) {
	filter := bson.M{}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"district_name": searchRegex},
			{"district_code": searchRegex},
		}
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	// Get total count
	totalCount, err := a.districtDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("Error counting districts: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	results, err := a.districtDal.FindAllWithPagination(ctx, filter, nil, skip, limit)
	if err != nil {
		a.logger.Errorf("Error finding districts with pagination: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// If filterParam.Search is empty, return empty data and meta using local_utils
	if filterParam.Search == "" {
		return &types.PaginatedResponse[[]*model.District]{
			Data: []*model.District{},
			Meta: types.PaginationMeta{TotalDocs: totalCount},
		}, nil
	}

	return &types.PaginatedResponse[[]*model.District]{
		Data: results,
		Meta: types.PaginationMeta{TotalDocs: totalCount},
	}, nil
}


func (a *AccountBlockStorage) EnableOrDisable(ctx context.Context, blockType string, codes []string, cpsAction model.CPSAction, requestType constants.RequestAction) error {
	var filter bson.M

	switch blockType {
	case "BRANCH":
		filter = bson.M{"branch_code": bson.M{"$in": codes}}

		// Check branches exists
		_, err := a.branchDal.FindOne(ctx, filter, bson.M{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return errors.New("BRANCH_NOT_FOUND")
			}
			return errors.New("GENERAL_DB_QUERY_FAILED")
		}

		// Check if duplicate action is requested
		var boolStatus bool
		switch requestType {
		case constants.RequestEnableBranches:
			boolStatus = true
		case constants.RequestDisableBranches:
			boolStatus = false
		}

		for _, code := range codes {
			dupFilter := bson.M{"branch_code": code, "enabled": boolStatus}
			branch, err := a.branchDal.FindOne(ctx, dupFilter, nil)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					continue
				}
				return errors.New("database error while finding user")
			}

			if branch.Enabled == boolStatus {
				return errors.New("DUPLICATE_ACTION")
			}
		}
	case "REGION":
		filter = bson.M{"region_code": bson.M{"$in": codes}}

		_, err := a.regionDal.FindOne(ctx, filter, bson.M{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return errors.New("REGION_NOT_FOUND")
			}
			return errors.New("GENERAL_DB_QUERY_FAILED")
		}

		var boolStatus bool
		switch requestType {
		case constants.RequestEnableRegion:
			boolStatus = true
		case constants.RequestDisableRegion:
			boolStatus = false
		}

		for _, code := range codes {
			dupFilter := bson.M{"region_code": code, "enabled": boolStatus}
			branch, err := a.regionDal.FindOne(ctx, dupFilter, nil)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					continue
				}
				return errors.New("database error while finding user")
			}

			if branch.Enabled == boolStatus {
				return errors.New("DUPLICATE_ACTION")
			}
		}
	case "DISTRICT":
		filter = bson.M{"district_code": bson.M{"$in": codes}}

		_, err := a.districtDal.FindOne(ctx, filter, bson.M{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return errors.New("DISTRICT_NOT_FOUND")
			}
			return errors.New("GENERAL_DB_QUERY_FAILED")
		}

		var boolStatus bool
		switch requestType {
		case constants.RequestEnableDistrict:
			boolStatus = true
		case constants.RequestDisableDistrict:
			boolStatus = false
		}

		for _, code := range codes {
			dupFilter := bson.M{"district_code": code, "enabled": boolStatus}
			branch, err := a.districtDal.FindOne(ctx, dupFilter, nil)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					continue
				}
				return errors.New("database error while finding user")
			}

			if branch.Enabled == boolStatus {
				return errors.New("DUPLICATE_ACTION")
			}
		}
	case "CITY":
		filter = bson.M{"city_code": bson.M{"$in": codes}}

		_, err := a.cityDal.FindOne(ctx, filter, bson.M{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return errors.New("CITY_NOT_FOUND")
			}
			return errors.New("GENERAL_DB_QUERY_FAILED")
		}
		var boolStatus bool
		switch requestType {
		case constants.RequestEnableCity:
			boolStatus = true
		case constants.RequestDisableCity:
			boolStatus = false
		}

		for _, code := range codes {
			dupFilter := bson.M{"city_code": code, "enabled": boolStatus}
			branch, err := a.cityDal.FindOne(ctx, dupFilter, nil)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					continue
				}
				return errors.New("database error while finding user")
			}

			if branch.Enabled == boolStatus {
				return errors.New("DUPLICATE_ACTION")
			}
		}
	}

	// Check for any pending action
	maker := local_util.ExtractUserFromContext(ctx)
	pendingFilter := bson.M{
		"maker_id":       maker.UserID,
		"department":     maker.Department,
		"action_status":  "PENDING",
		"request_action": requestType,
	}

	pendingAction, err := a.actionDal.FindOne(ctx, pendingFilter, bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.New("DATABASE_ERROR_CHECKING_PENDING_ACTION")
	}
	if pendingAction != nil {
		return errors.New("PENDING_ACTION_EXISTS")
	}

	_, err = a.actionDal.InsertOne(ctx, cpsAction)
	if err != nil {
		return errors.New("database errorwhile creating CPS action")
	}

	return nil
}

func (a *AccountBlockStorage) AuthorizeEnableBranches(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	data, err := local_util.JsonUnmarshal[[]model.Branch](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"branch_code": service.BranchCode}
		update := bson.M{"enabled": true}
		_, err = a.branchDal.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (a *AccountBlockStorage) AuthorizeDisableBranches(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	data, err := local_util.JsonUnmarshal[[]model.Branch](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"branch_code": service.BranchCode}
		update := bson.M{"enabled": false}
		_, err = a.branchDal.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (a *AccountBlockStorage) AuthorizeEnableRegions(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	data, err := local_util.JsonUnmarshal[[]model.Region](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"region_code": service.RegionCode}
		update := bson.M{"enabled": true}
		_, err = a.regionDal.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (a *AccountBlockStorage) AuthorizeDisableRegions(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	data, err := local_util.JsonUnmarshal[[]model.Region](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"region_code": service.RegionCode}
		update := bson.M{"enabled": false}
		_, err = a.regionDal.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (a *AccountBlockStorage) AuthorizeEnableDistricts(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	data, err := local_util.JsonUnmarshal[[]model.District](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"district_code": service.DistrictCode}
		update := bson.M{"enabled": true}
		_, err = a.districtDal.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (a *AccountBlockStorage) AuthorizeDisableDistrict(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	data, err := local_util.JsonUnmarshal[[]model.District](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"district_code": service.DistrictCode}
		update := bson.M{"enabled": false}
		_, err = a.districtDal.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (a *AccountBlockStorage) AuthorizeEnableCities(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	data, err := local_util.JsonUnmarshal[[]model.City](action.CurrentAction)
	a.logger.Infof("Data: %v\n", data)
	if err != nil {
		return nil, err
	}
	for _, code := range *data {
		filter := bson.M{"city_code": code.CityCode}
		update := bson.M{"enabled": true}
		a.logger.Infof("Filter:", filter)
		a.logger.Infof("Update:", update)
		_, err = a.cityDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("Error", err)
			return nil, err
		}
	}

	return action, nil
}

func (a *AccountBlockStorage) AuthorizeDisableCities(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	data, err := local_util.JsonUnmarshal[[]model.City](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"city_code": service.CityCode}
		update := bson.M{"enabled": false}
		_, err = a.cityDal.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

// The rest of the code remains unchanged (bulk enable/disable/approve methods)
