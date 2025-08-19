package adapter

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	ctx_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type outboundAccountBlockStore struct {
	MongoDalBranch    *infra_mongo.MongoDal[model.Branch, model.Branch]
	MongoDalRegion    *infra_mongo.MongoDal[model.Region, model.Region]
	MongoDalDistrict  *infra_mongo.MongoDal[model.District, model.District]
	MongoDalCity      *infra_mongo.MongoDal[model.City, model.City]
	MongoDalCPSAction *infra_mongo.MongoDal[model.CPSAction, model.CPSAction]
	MongoDalUser      *infra_mongo.MongoDal[member.User, member.User]
	Logger            utils.Logger
}

func NewOutboundAccountBlockStore(
	client *mongo.Client,
	dbName string,
	branchCollection, regionCollection, cpsActionCollection, districtCollection, userCollection, cityCollection string,
	logger utils.Logger,
) account_block.AccountBlockOutboundPort {
	return &outboundAccountBlockStore{
		MongoDalBranch:    infra_mongo.NewMongoDal[model.Branch, model.Branch](client, dbName, branchCollection),
		MongoDalRegion:    infra_mongo.NewMongoDal[model.Region, model.Region](client, dbName, regionCollection),
		MongoDalDistrict:  infra_mongo.NewMongoDal[model.District, model.District](client, dbName, districtCollection),
		MongoDalCity:      infra_mongo.NewMongoDal[model.City, model.City](client, dbName, cityCollection),
		MongoDalCPSAction: infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, cpsActionCollection),
		MongoDalUser:      infra_mongo.NewMongoDal[member.User, member.User](client, dbName, userCollection),
		Logger:            logger,
	}
}

func (o *outboundAccountBlockStore) GetBranch(ctx context.Context, branchCode string, filterParams *constant.Filter) (*model.Branch, error) {
	filter := bson.M{"branch_code": branchCode}

	branch, err := o.MongoDalBranch.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("BRANCH_NOT_FOUND")
		}
		o.Logger.Errorf("Unexpected database error occured while fetching branch by branch_code: %v", branchCode)
		return nil, fmt.Errorf("")
	}

	return &model.Branch{
		ID:            branch.ID,
		BranchCode:    branch.BranchCode,
		BranchName:    branch.BranchName,
		BranchAddress: branch.BranchAddress,
		DistrictCode:  branch.DistrictCode,
		DistrictName:  branch.DistrictName,
		BranchRegion:  branch.BranchRegion,
		RecordStat:    branch.RecordStat,
		CreatedAt:     branch.CreatedAt,
		UpdatedAt:     branch.UpdatedAt,
		Enabled:       branch.Enabled,
	}, nil

}

func (o *outboundAccountBlockStore) GetAllBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage
	projection := bson.M{}

	filter := bson.M{
		"branch_region": bson.M{"$regex": region, "$options": "i"},
		"district_name": bson.M{"$regex": district, "$options": "i"},
	}

	branches, err := o.MongoDalBranch.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		o.Logger.Errorf("failed to fetch branches: %v", err)
		return nil, common.DefineError.Branch["FAILED_TO_FETCH_BRANCHES"]
	}

	if len(branches) == 0 {
		o.Logger.Warnf("No branches found for region: %s, district: %s", region, district)
		return nil, common.DefineError.Branch["BRANCH_NOT_FOUND"]
	}

	total, err := o.MongoDalBranch.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := constant_utils.BuildPaginationMeta(total, filterParams.Page, limit)

	return &constant_utils.PaginatedResponse[[]*model.Branch]{
		Data: branches,
		Meta: meta,
	}, nil
}

func (o *outboundAccountBlockStore) GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error) {
	if strings.TrimSpace(regionCode) == "" {
		return action.Region{}, common.DefineError.Branch["REGION_CODE_REQUIRED"]
	}
	filter := bson.M{"region_code": regionCode}
	regionDoc, err := o.MongoDalRegion.FindOne(ctx, filter, nil)
	if err != nil || regionDoc == nil {
		return action.Region{}, common.DefineError.Branch["REGION_NOT_FOUND"]
	}
	return action.Region{
		ID:            regionDoc.ID.Hex(),
		RegionCode:    regionDoc.RegionCode,
		RegionName:    regionDoc.RegionName,
		RegionAddress: regionDoc.RegionAddress,
		CreatedAt:     regionDoc.CreatedAt,
		UpdatedAt:     regionDoc.UpdatedAt,
		Enabled:       regionDoc.Enabled,
	}, nil
}

func (o *outboundAccountBlockStore) GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error) {
	if strings.TrimSpace(districtCode) == "" {
		return action.District{}, fmt.Errorf("district_code is required")
	}
	filter := bson.M{"district_code": districtCode}
	districtDoc, err := o.MongoDalDistrict.FindOne(ctx, filter, nil)
	if err != nil || districtDoc == nil {
		return action.District{}, fmt.Errorf("FAILED_TO_GET_DISTRICT")
	}
	return action.District{
		ID:              districtDoc.ID.Hex(),
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

func (o *outboundAccountBlockStore) GetCityByCode(ctx context.Context, cityCode string) (action.City, error) {
	if strings.TrimSpace(cityCode) == "" {
		return action.City{}, fmt.Errorf("city_code is required")
	}
	filter := bson.M{"city_code": cityCode}
	cityDoc, err := o.MongoDalCity.FindOne(ctx, filter, nil)

	if err != nil || cityDoc == nil {
		return action.City{}, fmt.Errorf("failed to get city by code")
	}
	return action.City{
		ID:           cityDoc.ID.Hex(),
		CityAddress:  cityDoc.City,
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

func (o *outboundAccountBlockStore) GetAllCities(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, common.DefineError.General["INVALID_PAGINATION_PARAMS"]
	}
	filter := bson.M{}
	projection := bson.M{}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage
	// Support search and filter for GetAllCities
	if filterParams.Search != "" {
		// Assuming search on city_name or city_code
		filter["$or"] = []bson.M{
			{"city_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"city_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}
	// Additional filters (if any) from filterParams.Filters map
	if filterParams.Filters != "" {
		filter[filterParams.Filters] = filterParams.Filters

		// for k, v := range filterParams.Filters {
		// 	ks, ok1 := k.(string)
		// 	vs, ok2 := v.(string)
		// 	if ok1 && ok2 && ks != "" && vs != "" {
		// 	}
		// }
	}

	cities, err := o.MongoDalCity.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		o.Logger.Errorf("failed to fetch cities: %v", err)
		return nil, common.DefineError.General["FAILED_TO_FETCH"]
	}

	total, err := o.MongoDalCity.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := constant_utils.BuildPaginationMeta(total, filterParams.Page, limit)

	return &constant_utils.PaginatedResponse[[]*model.City]{
		Data: cities,
		Meta: meta,
	}, nil

}

func (o *outboundAccountBlockStore) GetAllDistricts(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, common.DefineError.General["INVALID_PAGINATION_PARAMS"]
	}
	filter := bson.M{}
	projection := bson.M{}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	if filterParams.Search != "" {
		// Assuming search on city_name or city_code
		filter["$or"] = []bson.M{
			{"district_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"district_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	districts, err := o.MongoDalDistrict.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		o.Logger.Errorf("failed to fetch districts: %v", err)
		return nil, common.DefineError.General["FAILED_TO_FETCH"]
	}

	total, err := o.MongoDalDistrict.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := constant_utils.BuildPaginationMeta(total, filterParams.Page, limit)

	return &constant_utils.PaginatedResponse[[]*model.District]{
		Data: districts,
		Meta: meta,
	}, nil
}

func (o *outboundAccountBlockStore) GetAllRegions(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, common.DefineError.General["INVALID_PAGINATION_PARAMS"]
	}
	filter := bson.M{}
	projection := bson.M{}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	if filterParams.Search != "" {
		// Assuming search on city_name or city_code
		filter["$or"] = []bson.M{
			{"region_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"region_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	regions, err := o.MongoDalRegion.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		o.Logger.Errorf("failed to fetch regions: %v", err)
		return nil, common.DefineError.General["FAILED_TO_FETCH"]
	}

	total, err := o.MongoDalRegion.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := constant_utils.BuildPaginationMeta(total, filterParams.Page, limit)

	return &constant_utils.PaginatedResponse[[]*model.Region]{
		Data: regions,
		Meta: meta,
	}, nil
}

func (o *outboundAccountBlockStore) EnableOrDisable(ctx context.Context, blockType string, codes []string, cpsAction model.CPSAction, requestType model.RequestAction) error {
	var filter bson.M

	switch blockType {
	case "BRANCH":
		filter = bson.M{"branch_code": bson.M{"$in": codes}}

		// Check branches exists
		_, err := o.MongoDalBranch.FindOne(ctx, filter, bson.M{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fmt.Errorf("BRANCH_NOT_FOUND")
			}
			return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
		}

		// Check if duplicate action is requested
		var boolStatus bool
		switch requestType {
		case model.RequestEnableBranches:
			boolStatus = true
		case model.RequestDisableBranches:
			boolStatus = false
		}

		for _, code := range codes {
			dupFilter := bson.M{"branch_code": code, "enabled": boolStatus}
			branch, err := o.MongoDalBranch.FindOne(ctx, dupFilter, nil)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					continue
				}
				return fmt.Errorf("database error while finding user")
			}

			if branch.Enabled == boolStatus {
				return fmt.Errorf("DUPLICATE_ACTION")
			}
		}
	case "REGION":
		filter = bson.M{"region_code": bson.M{"$in": codes}}

		_, err := o.MongoDalRegion.FindOne(ctx, filter, bson.M{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fmt.Errorf("REGION_NOT_FOUND")
			}
			return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
		}

		var boolStatus bool
		switch requestType {
		case model.RequestEnableRegion:
			boolStatus = true
		case model.RequestDisableRegion:
			boolStatus = false
		}

		for _, code := range codes {
			dupFilter := bson.M{"region_code": code, "enabled": boolStatus}
			branch, err := o.MongoDalRegion.FindOne(ctx, dupFilter, nil)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					continue
				}
				return fmt.Errorf("database error while finding user")
			}

			if branch.Enabled == boolStatus {
				return fmt.Errorf("DUPLICATE_ACTION")
			}
		}
	case "DISTRICT":
		filter = bson.M{"district_code": bson.M{"$in": codes}}

		_, err := o.MongoDalDistrict.FindOne(ctx, filter, bson.M{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fmt.Errorf("DISTRICT_NOT_FOUND")
			}
			return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
		}

		var boolStatus bool
		switch requestType {
		case model.RequestEnableDistrict:
			boolStatus = true
		case model.RequestDisableDistrict:
			boolStatus = false
		}

		for _, code := range codes {
			dupFilter := bson.M{"district_code": code, "enabled": boolStatus}
			branch, err := o.MongoDalDistrict.FindOne(ctx, dupFilter, nil)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					continue
				}
				return fmt.Errorf("database error while finding user")
			}

			if branch.Enabled == boolStatus {
				return fmt.Errorf("DUPLICATE_ACTION")
			}
		}
	case "CITY":
		filter = bson.M{"city_code": bson.M{"$in": codes}}

		_, err := o.MongoDalCity.FindOne(ctx, filter, bson.M{})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fmt.Errorf("CITY_NOT_FOUND")
			}
			return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
		}
		var boolStatus bool
		switch requestType {
		case model.RequestEnableCity:
			boolStatus = true
		case model.RequestDisableCity:
			boolStatus = false
		}

		for _, code := range codes {
			dupFilter := bson.M{"city_code": code, "enabled": boolStatus}
			branch, err := o.MongoDalCity.FindOne(ctx, dupFilter, nil)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					continue
				}
				return fmt.Errorf("database error while finding user")
			}

			if branch.Enabled == boolStatus {
				return fmt.Errorf("DUPLICATE_ACTION")
			}
		}
	}

	// Check for any pending action
	maker := ctx_utils.ExtractContext(ctx)
	pendingFilter := bson.M{
		"maker_id":       maker.UserID,
		"department":     maker.Department,
		"action_status":  "PENDING",
		"request_action": requestType,
	}

	pendingAction, err := o.MongoDalCPSAction.FindOne(ctx, pendingFilter, bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("DATABASE_ERROR_CHECKING_PENDING_ACTION")
	}
	if pendingAction != nil {
		return fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return fmt.Errorf("database errorwhile creating CPS action")
	}

	return nil
}

func (o *outboundAccountBlockStore) AuthorizeEnableBranches(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	data, err := constant_utils.JsonUnmarshal[[]model.Branch](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"branch_code": service.BranchCode}
		update := bson.M{"enabled": true}
		_, err = o.MongoDalBranch.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (o *outboundAccountBlockStore) AuthorizeDisableBranches(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	data, err := constant_utils.JsonUnmarshal[[]model.Branch](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"branch_code": service.BranchCode}
		update := bson.M{"enabled": false}
		_, err = o.MongoDalBranch.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (o *outboundAccountBlockStore) AuthorizeEnableRegions(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	data, err := constant_utils.JsonUnmarshal[[]model.Region](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"region_code": service.RegionCode}
		update := bson.M{"enabled": true}
		_, err = o.MongoDalRegion.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (o *outboundAccountBlockStore) AuthorizeDisableRegions(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	data, err := constant_utils.JsonUnmarshal[[]model.Region](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"region_code": service.RegionCode}
		update := bson.M{"enabled": false}
		_, err = o.MongoDalRegion.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (o *outboundAccountBlockStore) AuthorizeEnableDistricts(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	data, err := constant_utils.JsonUnmarshal[[]model.District](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"district_code": service.DistrictCode}
		update := bson.M{"enabled": true}
		_, err = o.MongoDalDistrict.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (o *outboundAccountBlockStore) AuthorizeDisableDistrict(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	data, err := constant_utils.JsonUnmarshal[[]model.District](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"district_code": service.DistrictCode}
		update := bson.M{"enabled": false}
		_, err = o.MongoDalDistrict.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (o *outboundAccountBlockStore) AuthorizeEnableCities(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	data, err := constant_utils.JsonUnmarshal[[]model.City](action.CurrentAction)
	fmt.Printf("Data: %v\n", data)
	if err != nil {
		return nil, err
	}
	for _, code := range *data {
		filter := bson.M{"city_code": code.CityCode}
		update := bson.M{"enabled": true}
		fmt.Println("Filter:", filter)
		fmt.Println("Update:", update)
		_, err = o.MongoDalCity.UpdateOne(ctx, filter, update)
		if err != nil {
			fmt.Println("Error", err)
			return nil, err
		}
	}

	return action, nil
}

func (o *outboundAccountBlockStore) AuthorizeDisableCities(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	data, err := constant_utils.JsonUnmarshal[[]model.City](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"city_code": service.CityCode}
		update := bson.M{"enabled": false}
		_, err = o.MongoDalCity.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}
