package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/entity"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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
func (o *outboundAccountBlockStore) FilterSingleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	if region == "" || district == "" {
		return nil, common.DefineError.Branch["BRANCH_REGION_AND_DISTRICT_REQUIRED"]
	}

	filter := bson.M{"branch_region": region, "district_name": district}
	projection := bson.M{}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

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
	meta := constant_utils.BuildPaginationMeta(total, limit, filterParams.Page)

	return &constant_utils.PaginatedResponse[[]*model.Branch]{
		Data: branches,
		Meta: meta,
	}, nil
}
func (o *outboundAccountBlockStore) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error) {
	department, _ := ctx.Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		return "", common.DefineError.General["INCOMPLETE_USER_INFO"]
	}

	filter := bson.M{
		"unique_id":      branch.BranchCode,
		"department":     department,
		"action_type":    "DELETE",
		"request_action": "DISABLE_SINGLE_BRANCH",
		"action_status":  bson.M{"$in": []string{"PENDING", "APPROVED"}},
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", common.DefineError.Branch["BRANCH_DISABLE_ACTION_ALREADY_EXISTS"]
	}

	prevBranchPtr, err := o.MongoDalBranch.FindOne(ctx, bson.M{"branch_code": branch.BranchCode}, bson.M{})
	var prevAction json.RawMessage
	if err == nil && prevBranchPtr != nil {
		prevAction, _ = json.Marshal(prevBranchPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(branch)
	actionCode := utils.RandomGenerator(24)
	cpsAction := model.CPSAction{
		ActionCode:       actionCode,
		MakerID:          maker.UserID,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       department,
		UniqueId:         branch.BranchCode,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionDelete),
		RequestAction:    string(model.RequestDisableSingleBranch),
		PreviousAction:   prevAction,
		CurrentAction:    currAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return actionCode, nil
}

func (o *outboundAccountBlockStore) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if strings.TrimSpace(actionID) == "" {
		return common.DefineError.General["ACTION_ID_IS_REQUIRED"]
	}

	filter := bson.M{"action_code": actionID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return common.DefineError.General["ACTION_NOT_FOUND"]
	}

	uniqueID := actionDoc.UniqueId
	dupCheck := bson.M{
		"unique_id":      uniqueID,
		"action_status":  "APPROVED",
		"action_type":    "DELETE",
		"request_action": "DISABLE_SINGLE_BRANCH",
	}
	alreadyApproved, err := o.MongoDalCPSAction.FindOne(ctx, dupCheck, bson.M{})
	if approve && err == nil && alreadyApproved != nil {
		return common.DefineError.Branch["BRANCH_DISABLE_ACTION_ALREADY_EXISTS"]
	}

	update := bson.M{
		"last_modified_at": time.Now(),
	}
	if approve {
		update["action_status"] = "APPROVED"
	} else {
		update["action_status"] = "REJECTED"
		if reason != nil {
			update["rejection_reason"] = *reason
		}
	}

	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if approve {
		var branch action.Branch
		switch v := actionDoc.CurrentAction.(type) {
		case []byte:
			if err := json.Unmarshal(v, &branch); err != nil {
				return fmt.Errorf("failed to unmarshal branch ([]byte): %w", err)
			}
		case string:
			if err := json.Unmarshal([]byte(v), &branch); err != nil {
				return fmt.Errorf("failed to unmarshal branch (string): %w", err)
			}
		case bson.Binary:
			if err := json.Unmarshal(v.Data, &branch); err != nil {
				return fmt.Errorf("failed to unmarshal branch (bson.Binary): %w", err)
			}
		case map[string]interface{}:
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to marshal branch (map): %w", err)
			}
			if err := json.Unmarshal(b, &branch); err != nil {
				return fmt.Errorf("failed to unmarshal branch (map): %w", err)
			}
		default:
			return fmt.Errorf("CurrentAction is not a supported type, got %T", v)
		}

		if branch.BranchCode == "" {
			return common.DefineError.Branch["BRANCH_ID_REQUIRED"]
		}
		branchUpdate := bson.M{
			"enabled":    false,
			"updated_at": time.Now(),
		}
		_, err = o.MongoDalBranch.UpdateOne(ctx, bson.M{"branch_code": branch.BranchCode}, branchUpdate)
		if err != nil {
			return fmt.Errorf("failed to update branch: %w", err)
		}
	}

	return nil
}
func (o *outboundAccountBlockStore) FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	if region == "" {
		return nil, common.DefineError.Branch["BRANCH_REGION_AND_DISTRICT_REQUIRED"]
	}

	filter := bson.M{"branch_region": region}
	if district != "" {
		filter["district_name"] = district
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage
	projection := bson.M{}

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
	meta := constant_utils.BuildPaginationMeta(total, limit, filterParams.Page)

	return &constant_utils.PaginatedResponse[[]*model.Branch]{
		Data: branches,
		Meta: meta,
	}, nil
}

func (o *outboundAccountBlockStore) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error) {
	department, _ := ctx.Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		return "", common.DefineError.General["INCOMPLETE_USER_INFO"]
	}

	branchCodeSet := make(map[string]struct{})
	for _, branch := range branches {
		code := strings.TrimSpace(branch.BranchCode)
		if code != "" {
			branchCodeSet[code] = struct{}{}
		}
	}
	if len(branchCodeSet) == 0 {
		return "", common.DefineError.Branch["BRANCH_ID_REQUIRED"]
	}
	var branchCodes []string
	for code := range branchCodeSet {
		branchCodes = append(branchCodes, code)
	}
	sort.Strings(branchCodes)

	currActionBytes, _ := json.Marshal(branches)

	filter := bson.M{
		"department":     department,
		"action_type":    "DELETE",
		"action_status":  bson.M{"$in": []string{"PENDING", "APPROVED"}},
		"request_action": "REQUEST_DISABLE_MULTI_BRANCHES",
		"current_action": currActionBytes,
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", common.DefineError.Branch["BRANCH_DISABLE_MULTI_ACTION_ALREADY_EXISTS"]
	}

	var prevBranches []model.Branch
	for _, branch := range branches {
		prevBranchPtr, err := o.MongoDalBranch.FindOne(ctx, bson.M{"branch_code": branch.BranchCode}, bson.M{})
		if err == nil && prevBranchPtr != nil {
			prevBranches = append(prevBranches, *prevBranchPtr)
		}
	}
	prevAction, _ := json.Marshal(prevBranches)

	actionCode := utils.RandomGenerator(24)
	cpsAction := model.CPSAction{
		ActionCode:       actionCode,
		MakerID:          maker.UserID,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       department,
		UniqueId:         strings.Join(branchCodes, ","),
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionDelete),
		RequestAction:    string(model.RequestDisableMultiBranches),
		PreviousAction:   prevAction,
		CurrentAction:    currActionBytes,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}
	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return actionCode, nil
}

func (o *outboundAccountBlockStore) GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error) {
	filter := bson.M{"branch_code": branchCode}
	b, err := o.MongoDalBranch.FindOne(ctx, filter, nil)
	if err != nil || b == nil {
		return action.Branch{}, common.DefineError.Branch["BRANCH_NOT_FOUND"]
	}
	return action.Branch{
		ID:            b.ID.Hex(),
		BranchCode:    b.BranchCode,
		BranchName:    b.BranchName,
		BranchAddress: b.BranchAddress,
		DistrictCode:  b.DistrictCode,
		DistrictName:  b.DistrictName,
		BranchRegion:  b.BranchRegion,
		RecordStat:    b.RecordStat,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		Version:       b.Version,
		Enabled:       b.Enabled,
	}, nil
}
func (o *outboundAccountBlockStore) BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error) {
	department, _ := ctx.Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		return "", common.DefineError.General["INCOMPLETE_USER_INFO"]
	}

	filter := bson.M{
		"unique_id":      regionCode,
		"department":     department,
		"action_type":    "DELETE",
		"request_action": "BLOCK_REGION",
		"action_status":  bson.M{"$in": []string{"PENDING", "APPROVED"}},
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", common.DefineError.Branch["BLOCK_REGION_ALREADY_EXISTS"]
	}

	prevRegionPtr, err := o.MongoDalRegion.FindOne(ctx, bson.M{"region_code": regionCode}, bson.M{})
	var prevAction json.RawMessage
	if err == nil && prevRegionPtr != nil {
		prevAction, _ = json.Marshal(prevRegionPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(map[string]string{"region_code": regionCode})

	cpsAction := model.CPSAction{
		ActionCode:       utils.RandomGenerator(24),
		MakerID:          maker.MakerID,
		MakerName:        maker.MakerName,
		MakerPhoneNumber: maker.MakerPhoneNumber,
		Department:       department,
		UniqueId:         regionCode,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionDelete),
		RequestAction:    string(model.RequestBlockRegion),
		PreviousAction:   prevAction,
		CurrentAction:    currAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", common.DefineError.General["DATABASE_ERROR"]
	}
	return cpsAction.ActionCode, nil
}
func (o *outboundAccountBlockStore) UpdateRegion(ctx context.Context, region action.Region) error {
	filter := bson.M{"id": region.ID}
	update := bson.M{
		"$set": bson.M{
			"region_code":    region.RegionCode,
			"region_name":    region.RegionName,
			"region_address": region.RegionAddress,
			"created_at":     region.CreatedAt,
			"updated_at":     region.UpdatedAt,
			"enabled":        region.Enabled,
		},
	}
	_, err := o.MongoDalRegion.UpdateOne(ctx, filter, update)
	return err
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
func (o *outboundAccountBlockStore) BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error) {
	department := maker.Department
	if strings.TrimSpace(department) == "" {
		return "", errors.New("department is required in context")
	}

	filter := bson.M{
		"unique_id":      districtCode,
		"department":     department,
		"action_type":    "DELETE",
		"request_action": "BLOCK_DISTRICT",
		"action_status":  bson.M{"$in": []string{"PENDING", "APPROVED"}},
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", fmt.Errorf("A pending or approved block action already exists for this district")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return "", fmt.Errorf("database error on FindOne CPSAction: %w", err)
	}

	prevDistrictPtr, err := o.MongoDalDistrict.FindOne(ctx, bson.M{"district_code": districtCode}, bson.M{})
	var prevAction json.RawMessage
	if prevDistrictPtr != nil {
		prevAction, _ = json.Marshal(prevDistrictPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(map[string]string{"district_code": districtCode})

	cpsAction := model.CPSAction{
		ActionCode:       utils.RandomGenerator(24),
		MakerID:          maker.MakerID,
		MakerName:        maker.MakerName,
		MakerPhoneNumber: maker.MakerPhoneNumber,
		Department:       department,
		UniqueId:         districtCode,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionDelete),
		RequestAction:    string(model.RequestBlockDistrict),
		PreviousAction:   prevAction,
		CurrentAction:    currAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		RejectionReason: func() string {
			if maker.RejectionReason != nil {
				return *maker.RejectionReason
			}
			return ""
		}(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", fmt.Errorf("database error on InsertOne CPSAction: %w", err)
	}
	return cpsAction.ActionCode, nil
}
func (o *outboundAccountBlockStore) GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error) {
	if strings.TrimSpace(districtCode) == "" {
		return action.District{}, fmt.Errorf("district_code is required")
	}
	filter := bson.M{"district_code": districtCode}
	districtDoc, err := o.MongoDalDistrict.FindOne(ctx, filter, nil)
	if err != nil || districtDoc == nil {
		return action.District{}, fmt.Errorf("failed to get district by code: %w", err)
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

func (o *outboundAccountBlockStore) BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error) {
	department := maker.Department
	if strings.TrimSpace(department) == "" {
		return "", errors.New("department is required in context")
	}

	filter := bson.M{
		"unique_id":      cityCode,
		"department":     department,
		"action_type":    "DELETE",
		"request_action": "BLOCK_CITY",
		"action_status":  bson.M{"$in": []string{"PENDING", "APPROVED"}},
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", fmt.Errorf("A pending or approved block action already exists for this city")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return "", fmt.Errorf("database error on FindOne CPSAction: %w", err)
	}

	prevCityPtr, err := o.MongoDalCity.FindOne(ctx, bson.M{"city_code": cityCode}, bson.M{})
	var prevAction json.RawMessage
	if prevCityPtr != nil {
		prevAction, _ = json.Marshal(prevCityPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(map[string]string{"city_code": cityCode})

	cpsAction := model.CPSAction{
		ActionCode:       utils.RandomGenerator(24),
		MakerID:          maker.MakerID,
		MakerName:        maker.MakerName,
		MakerPhoneNumber: maker.MakerPhoneNumber,
		Department:       department,
		UniqueId:         cityCode,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionDelete),
		RequestAction:    string(model.RequestBlockCity),
		PreviousAction:   prevAction,
		CurrentAction:    currAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		RejectionReason: func() string {
			if maker.RejectionReason != nil {
				return *maker.RejectionReason
			}
			return ""
		}(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", fmt.Errorf("database error on InsertOne CPSAction: %w", err)
	}
	return cpsAction.ActionCode, nil
}
func (o *outboundAccountBlockStore) GetCityByCode(ctx context.Context, cityCode string) (action.City, error) {
	if strings.TrimSpace(cityCode) == "" {
		return action.City{}, fmt.Errorf("city_code is required")
	}
	filter := bson.M{"city_code": cityCode}
	cityDoc, err := o.MongoDalCity.FindOne(ctx, filter, nil)
	if err != nil || cityDoc == nil {
		return action.City{}, fmt.Errorf("failed to get city by code: %w", err)
	}
	return action.City{
		ID:           cityDoc.ID.Hex(),
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
func (o *outboundAccountBlockStore) BlockUser(ctx context.Context, phoneNumber string, maker action.CPSAction) (string, error) {
	department, _ := ctx.Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		return "", errors.New("department is required in context")
	}

	user, err := o.GetUserByPhone(ctx, phoneNumber, maker)
	if err != nil || strings.TrimSpace(user.UserCode) == "" {
		o.Logger.Errorf("BlockUser: user not found for phone_number='%s', err=%v", phoneNumber, err)
		return "", fmt.Errorf("user not found")
	}
	userCode := user.UserCode

	filter := bson.M{
		"unique_id":      userCode,
		"department":     department,
		"action_type":    "DELETE",
		"request_action": "BLOCK_USER",
		"action_status":  bson.M{"$in": []string{"PENDING", "APPROVED"}},
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", fmt.Errorf("BLOCKED_ACTION_USER_ALREADY_EXIST")
	}

	prevAction, _ := json.Marshal(user)
	currAction, _ := json.Marshal(userCode)
	actionCode := utils.RandomGenerator(24)
	cpsAction := model.CPSAction{
		ActionCode:       actionCode,
		MakerID:          maker.MakerID,
		MakerName:        maker.MakerName,
		MakerPhoneNumber: maker.MakerPhoneNumber,
		Department:       department,
		UniqueId:         userCode,
		PreviousAction:   prevAction,
		CurrentAction:    currAction,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionDelete),
		RequestAction:    string(model.RequestBlockUser),
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return actionCode, nil
}
func (o *outboundAccountBlockStore) GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error) {
	if strings.TrimSpace(phoneNumber) == "" {
		return member.User{}, fmt.Errorf("phoneNumber is required")
	}
	filter := bson.M{"phone_number": phoneNumber}

	userDoc, err := o.MongoDalUser.FindOne(ctx, filter, nil)
	if err == nil && userDoc != nil {
		return *userDoc, nil
	}

	if strings.HasPrefix(phoneNumber, "+") {
		filter = bson.M{"phone_number": phoneNumber[1:]}
	} else {
		filter = bson.M{"phone_number": "+" + phoneNumber}
	}
	userDoc, err = o.MongoDalUser.FindOne(ctx, filter, nil)
	if err == nil && userDoc != nil {
		return *userDoc, nil
	}

	return member.User{}, fmt.Errorf("failed to get user by phone number: %w", err)
}

// ********************************************************************
func (o *outboundAccountBlockStore) AuthorizeBlockUser(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	var actionData entity.BankDocument
	raw, _ := bson.Marshal(cpsAction.CurrentAction)
	if err := bson.Unmarshal(raw, &actionData); err != nil {
		o.Logger.Errorf("failed to unmarshal action data for authorization, action_code: %s", cpsAction.ActionCode)
		return nil, fmt.Errorf(error_codes.InvalidActionData)
	}

	var userCode string
	switch v := cpsAction.CurrentAction.(type) {
	case string:
		userCode = v
	case []byte:
		_ = json.Unmarshal(v, &userCode)
	case bson.Binary:
		_ = json.Unmarshal(v.Data, &userCode)
	}
	userCode = strings.TrimSpace(userCode)
	if userCode == "" {
		return nil, common.DefineError.General["USER_CODE_IS_REQUIRED"]
	}

	userUpdate := bson.M{
		"is_account_blocked": true,
		"last_modified_at":   time.Now(),
	}
	_, err := o.MongoDalUser.UpdateOne(ctx, bson.M{"user_code": userCode}, userUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to block user %s: %w", userCode, err)
	}

	return cpsAction, nil

}
func (o *outboundAccountBlockStore) AuthorizeBlockCity(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	var cityCode string
	switch v := cpsAction.CurrentAction.(type) {
	case string:
		var tmp struct {
			CityCode string `json:"city_code"`
		}
		_ = json.Unmarshal([]byte(v), &tmp)
		cityCode = tmp.CityCode
	case []byte:
		var tmp struct {
			CityCode string `json:"city_code"`
		}
		_ = json.Unmarshal(v, &tmp)
		cityCode = tmp.CityCode
	case bson.Binary:
		var tmp struct {
			CityCode string `json:"city_code"`
		}
		_ = json.Unmarshal(v.Data, &tmp)
		cityCode = tmp.CityCode
	case map[string]interface{}:
		if cc, ok := v["city_code"].(string); ok {
			cityCode = cc
		}
	default:
		return nil, errors.New("city_code missing in action")
	}
	cityCode = strings.TrimSpace(cityCode)
	if cityCode == "" {
		return nil, errors.New("city_code missing in action")
	}

	cityUpdate := bson.M{
		"enabled":    false,
		"updated_at": time.Now(),
	}
	_, err := o.MongoDalCity.UpdateOne(ctx, bson.M{"city_code": cityCode}, cityUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to disable city %s: %w", cityCode, err)
	}

	return cpsAction, nil

}
func (o *outboundAccountBlockStore) AuthorizeBlockDistrict(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	var districtCode string
	switch v := cpsAction.CurrentAction.(type) {
	case string:
		var tmp struct {
			DistrictCode string `json:"district_code"`
		}
		_ = json.Unmarshal([]byte(v), &tmp)
		districtCode = tmp.DistrictCode
	case []byte:
		var tmp struct {
			DistrictCode string `json:"district_code"`
		}
		_ = json.Unmarshal(v, &tmp)
		districtCode = tmp.DistrictCode
	case bson.Binary:
		var tmp struct {
			DistrictCode string `json:"district_code"`
		}
		_ = json.Unmarshal(v.Data, &tmp)
		districtCode = tmp.DistrictCode
	case map[string]interface{}:
		if dc, ok := v["district_code"].(string); ok {
			districtCode = dc
		}
	default:
		return nil, errors.New("district_code missing in action")
	}
	districtCode = strings.TrimSpace(districtCode)
	if districtCode == "" {
		return nil, errors.New("district_code missing in action")
	}

	districtUpdate := bson.M{
		"enabled":    false,
		"updated_at": time.Now(),
	}
	_, err := o.MongoDalDistrict.UpdateOne(ctx, bson.M{"district_code": districtCode}, districtUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to disable district %s: %w", districtCode, err)
	}

	return cpsAction, nil

}
func (o *outboundAccountBlockStore) AuthorizeRegionBlock(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	var regionCode string
	switch v := cpsAction.CurrentAction.(type) {
	case string:
		var tmp struct {
			RegionCode string `json:"region_code"`
		}
		_ = json.Unmarshal([]byte(v), &tmp)
		regionCode = tmp.RegionCode
	case []byte:
		var tmp struct {
			RegionCode string `json:"region_code"`
		}
		_ = json.Unmarshal(v, &tmp)
		regionCode = tmp.RegionCode
	case bson.Binary:
		var tmp struct {
			RegionCode string `json:"region_code"`
		}
		_ = json.Unmarshal(v.Data, &tmp)
		regionCode = tmp.RegionCode
	case map[string]interface{}:
		if rc, ok := v["region_code"].(string); ok {
			regionCode = rc
		}
	}
	regionCode = strings.TrimSpace(regionCode)
	if regionCode == "" {
		return nil, common.DefineError.Branch["REGION_CODE_REQUIRED"]
	}

	region, err := o.GetRegionByCode(ctx, regionCode)
	if err != nil {
		return nil, common.DefineError.Branch["REGION_NOT_FOUND"]
	}

	branchUpdate := bson.M{
		"$set": bson.M{
			"enabled":    false,
			"updated_at": time.Now(),
		},
	}
	branchFilter := bson.M{
		"$or": []bson.M{
			{"branch_region_code": region.RegionCode},
			{"branch_region": region.RegionName},
		},
	}
	branches, err := o.MongoDalBranch.FindAll(ctx, branchFilter, nil)
	if err != nil {
		return nil, common.DefineError.Branch["FAILED_TO_FETCH_BRANCHES"]
	}
	for _, branch := range branches {
		if branch == nil {
			continue
		}
		_, err := o.MongoDalBranch.UpdateOne(ctx, bson.M{"_id": branch.ID}, branchUpdate)
		if err != nil {
			return nil, common.DefineError.Branch["FAILED_TO_DISABLE_BRANCH"]
		}
	}

	regionUpdate := bson.M{
		"enabled":    false,
		"updated_at": time.Now(),
	}
	_, err = o.MongoDalRegion.UpdateOne(ctx, bson.M{"region_code": regionCode}, regionUpdate)
	if err != nil {
		return nil, common.DefineError.Branch["FAILED_TO_DISABLE_REGION"]
	}

	return cpsAction, nil

}
func (o *outboundAccountBlockStore) AuthorizeBulkBranchesDisable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	var branches []action.Branch
	switch v := cpsAction.CurrentAction.(type) {
	case []byte:
		_ = json.Unmarshal(v, &branches)
	case string:
		_ = json.Unmarshal([]byte(v), &branches)
	case bson.Binary:
		_ = json.Unmarshal(v.Data, &branches)
	case map[string]interface{}:
		b, _ := json.Marshal(v)
		_ = json.Unmarshal(b, &branches)
	}
	for _, branch := range branches {
		if branch.BranchCode == "" {
			return nil, errors.New("branch code is required in action data")
		}
		branchUpdate := bson.M{
			"enabled":    false,
			"updated_at": time.Now(),
		}
		_, err := o.MongoDalBranch.UpdateOne(ctx, bson.M{"branch_code": branch.BranchCode}, branchUpdate)
		if err != nil {
			return nil, fmt.Errorf("failed to update branch %s: %w", branch.BranchCode, err)
		}
	}

	return cpsAction, nil

}
func (o *outboundAccountBlockStore) AuthorizeSingleBranchDisable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	var branch action.Branch
	switch v := cpsAction.CurrentAction.(type) {
	case []byte:
		if err := json.Unmarshal(v, &branch); err != nil {
			return nil, fmt.Errorf("failed to unmarshal branch ([]byte): %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(v), &branch); err != nil {
			return nil, fmt.Errorf("failed to unmarshal branch (string): %w", err)
		}
	case bson.Binary:
		if err := json.Unmarshal(v.Data, &branch); err != nil {
			return nil, fmt.Errorf("failed to unmarshal branch (bson.Binary): %w", err)
		}
	case map[string]interface{}:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal branch (map): %w", err)
		}
		if err := json.Unmarshal(b, &branch); err != nil {
			return nil, fmt.Errorf("failed to unmarshal branch (map): %w", err)
		}
	default:
		return nil, fmt.Errorf("CurrentAction is not a supported type, got %T", v)
	}

	if branch.BranchCode == "" {
		return nil, common.DefineError.Branch["BRANCH_ID_REQUIRED"]
	}
	branchUpdate := bson.M{
		"enabled":    false,
		"updated_at": time.Now(),
	}
	_, err := o.MongoDalBranch.UpdateOne(ctx, bson.M{"branch_code": branch.BranchCode}, branchUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to update branch: %w", err)
	}
	return cpsAction, nil

}
