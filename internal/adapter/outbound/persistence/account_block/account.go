package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/account_block"
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
func (o *outboundAccountBlockStore) FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {

	if region == "" || district == "" {
		return nil, errors.New("region and district are required")
	}

	filter := bson.M{"branch_region": region, "district_name": district}
	data, err := o.MongoDalBranch.FindAll(ctx, filter, nil)
	if err != nil {
		fmt.Printf("failed to fetch branches: %v\n", err)
		return nil, fmt.Errorf("failed to fetch branches: %w", err)
	}

	var result []action.Branch
	for _, b := range data {
		if b == nil {
			fmt.Println("Skipping nil branch document")
			continue
		}
		ab := action.Branch{
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
		}
		result = append(result, ab)
	}

	return result, nil
}
func (o *outboundAccountBlockStore) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error) {
	department, _ := ctx.Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		return "", errors.New("department is required in context")
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
		return "", fmt.Errorf("A pending or approved disable action already exists for this branch")
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
		RequestAction:    string(model.RequestDisableMultiBranches),
		PreviosAction:    prevAction,
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
		return errors.New("actionID is required")
	}

	filter := bson.M{"action_code": actionID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found")
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
		return errors.New("A pending disable action already exists for this branch")
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
			return errors.New("branch code is required in action data")
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
func (o *outboundAccountBlockStore) FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	if region == "" {
		return nil, errors.New("region is required")
	}

	filter := bson.M{"branch_region": region}
	if district != "" {
		filter["district_name"] = district
	}

	data, err := o.MongoDalBranch.FindAll(ctx, filter, nil)
	if err != nil {
		fmt.Printf("failed to fetch branches: %v\n", err)
		return nil, fmt.Errorf("failed to fetch branches: %w", err)
	}

	var result []action.Branch
	for _, b := range data {
		if b == nil {
			fmt.Println("Skipping nil branch document")
			continue
		}
		ab := action.Branch{
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
		}
		result = append(result, ab)
	}

	return result, nil
}
func (o *outboundAccountBlockStore) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error) {
	department, _ := ctx.Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		return "", errors.New("department is required in context")
	}

	var branchCodes []string
	for _, branch := range branches {
		branchCodes = append(branchCodes, branch.BranchCode)
	}

	filter := bson.M{
		"unique_id":      bson.M{"$in": branchCodes},
		"department":     department,
		"action_type":    "DELETE",
		"request_action": "DISABLE_MULTI_BRANCHES",
		"action_status":  bson.M{"$in": []string{"PENDING", "APPROVED"}},
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", fmt.Errorf("A pending or approved disable action already exists for one or more branches")
	}

	currAction, _ := json.Marshal(branches)
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
		PreviosAction:    prevAction,
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
func (o *outboundAccountBlockStore) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if strings.TrimSpace(actionID) == "" {
		return errors.New("actionID is required")
	}

	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, bson.M{"action_code": actionID}, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found")
	}

	uniqueIDs := strings.Split(actionDoc.UniqueId, ",")
	statusCheck := "APPROVED"
	if !approve {
		statusCheck = "REJECTED"
	}
	dupFilter := bson.M{
		"unique_id":      bson.M{"$in": uniqueIDs},
		"action_status":  statusCheck,
		"action_type":    "DELETE",
		"request_action": "DISABLE_MULTI_BRANCHES",
	}
	alreadyProcessed, _ := o.MongoDalCPSAction.FindOne(ctx, dupFilter, bson.M{})
	if alreadyProcessed != nil {
		return errors.New("This bulk action has already been processed for one or more branches")
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

	_, err = o.MongoDalCPSAction.UpdateOne(ctx, bson.M{"action_code": actionID}, update)
	if err != nil {
		return err
	}

	if approve {
		var branches []action.Branch
		switch v := actionDoc.CurrentAction.(type) {
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
				return errors.New("branch code is required in action data")
			}
			branchUpdate := bson.M{
				"enabled":    false,
				"updated_at": time.Now(),
			}
			_, err := o.MongoDalBranch.UpdateOne(ctx, bson.M{"branch_code": branch.BranchCode}, branchUpdate)
			if err != nil {
				return fmt.Errorf("failed to update branch %s: %w", branch.BranchCode, err)
			}
		}
	}

	return nil
}
func (o *outboundAccountBlockStore) GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error) {
	filter := bson.M{"branch_code": branchCode}
	b, err := o.MongoDalBranch.FindOne(ctx, filter, nil)
	if err != nil || b == nil {
		return action.Branch{}, err
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
		return "", errors.New("department is required in context")
	}

	filter := bson.M{
		"unique_id":      regionCode,
		"department":     department,
		"action_type":    "DELETE",
		"request_action": "DISABLE_REGION",
		"action_status":  bson.M{"$in": []string{"PENDING", "APPROVED"}},
	}
	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return "", fmt.Errorf("A pending or approved disable action already exists for this region")
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
		PreviosAction:    prevAction,
		CurrentAction:    currAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", fmt.Errorf("database error on InsertOne CPSAction: %w", err)
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
func (o *outboundAccountBlockStore) ApproveRegionBlock(
	ctx context.Context,
	actionID string,
	approve bool,
	reason *string,
	checker action.User,
) error {
	if strings.TrimSpace(actionID) == "" {
		return errors.New("actionID is required")
	}

	filter := bson.M{"action_code": actionID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		filter2 := bson.M{"action_code": actionID, "request_action": "REQUEST_DISABLE_MULTI_BRANCHES"}
		actionDoc, err = o.MongoDalCPSAction.FindOne(ctx, filter2, nil)
		if err != nil || actionDoc == nil {
			return errors.New("action not found after update")
		}
	}

	uniqueID := actionDoc.UniqueId
	statusCheck := "APPROVED"
	if !approve {
		statusCheck = "REJECTED"
	}
	dupFilter := bson.M{
		"unique_id":      uniqueID,
		"action_status":  statusCheck,
		"action_type":    "DELETE",
		"request_action": "DISABLE_REGION",
	}
	alreadyProcessed, _ := o.MongoDalCPSAction.FindOne(ctx, dupFilter, bson.M{})
	if alreadyProcessed != nil {
		return errors.New("This region action has already been processed")
	}

	update := bson.M{
		"last_modified_at":     time.Now(),
		"checker_id":           checker.UserID,
		"checker_name":         checker.FullName,
		"checker_phone_number": checker.PhoneNumber,
		"checker_action_time":  time.Now(),
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

	if !approve {
		return nil
	}

	var regionCode string
	switch v := actionDoc.CurrentAction.(type) {
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
		return errors.New("region_code missing in action")
	}

	region, err := o.GetRegionByCode(ctx, regionCode)
	if err != nil {
		return fmt.Errorf("region not found: %w", err)
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
		return fmt.Errorf("failed to fetch branches in region: %w", err)
	}
	for _, branch := range branches {
		if branch == nil {
			continue
		}
		_, err := o.MongoDalBranch.UpdateOne(ctx, bson.M{"_id": branch.ID}, branchUpdate)
		if err != nil {
			return fmt.Errorf("failed to disable branch %s: %w", branch.BranchCode, err)
		}
	}

	regionUpdate := bson.M{
		"enabled":    false,
		"updated_at": time.Now(),
	}
	_, err = o.MongoDalRegion.UpdateOne(ctx, bson.M{"region_code": regionCode}, regionUpdate)
	if err != nil {
		return fmt.Errorf("failed to disable region %s: %w", regionCode, err)
	}

	return nil
}
func (o *outboundAccountBlockStore) GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error) {
	if regionCode == "" {
		return action.Region{}, fmt.Errorf("regionCode is required")
	}
	filter := bson.M{"region_code": regionCode}
	regionDoc, err := o.MongoDalRegion.FindOne(ctx, filter, nil)
	if err != nil || regionDoc == nil {
		return action.Region{}, fmt.Errorf("failed to get region by code: %w", err)
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
		PreviosAction:    prevAction,
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
func (o *outboundAccountBlockStore) ApproveBlockDistrict(
	ctx context.Context,
	actionID string,
	approve bool,
	reason *string,
	checker action.User,
) error {
	if strings.TrimSpace(actionID) == "" {
		return errors.New("action_id is required")
	}

	filter := bson.M{"action_code": actionID, "request_action": "BLOCK_DISTRICT"}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found after update")
	}

	uniqueID := actionDoc.UniqueId
	statusCheck := "APPROVED"
	if !approve {
		statusCheck = "REJECTED"
	}
	dupFilter := bson.M{
		"unique_id":      uniqueID,
		"action_status":  statusCheck,
		"action_type":    "DELETE",
		"request_action": "BLOCK_DISTRICT",
	}
	alreadyProcessed, _ := o.MongoDalCPSAction.FindOne(ctx, dupFilter, bson.M{})
	if alreadyProcessed != nil {
		return errors.New(" This district action has already been processed")
	}

	update := bson.M{
		"last_modified_at":     time.Now(),
		"checker_id":           checker.UserID,
		"checker_name":         checker.FullName,
		"checker_phone_number": checker.PhoneNumber,
		"checker_action_time":  time.Now(),
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
		return fmt.Errorf("failed to update CPSAction: %w", err)
	}

	if !approve {
		return nil
	}

	var districtCode string
	switch v := actionDoc.CurrentAction.(type) {
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
		return errors.New("district_code missing in action")
	}
	districtCode = strings.TrimSpace(districtCode)
	if districtCode == "" {
		return errors.New("district_code missing in action")
	}

	districtUpdate := bson.M{
		"enabled":    false,
		"updated_at": time.Now(),
	}
	_, err = o.MongoDalDistrict.UpdateOne(ctx, bson.M{"district_code": districtCode}, districtUpdate)
	if err != nil {
		return fmt.Errorf("failed to disable district %s: %w", districtCode, err)
	}

	return nil
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
		PreviosAction:    prevAction,
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
func (o *outboundAccountBlockStore) ApproveBlockCity(
	ctx context.Context,
	actionID string,
	approve bool,
	reason *string,
	checker action.User,
) error {
	if strings.TrimSpace(actionID) == "" {
		return errors.New("action_id is required")
	}

	filter := bson.M{"action_code": actionID, "request_action": "BLOCK_CITY"}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found after update")
	}

	uniqueID := actionDoc.UniqueId
	statusCheck := "APPROVED"
	if !approve {
		statusCheck = "REJECTED"
	}
	dupFilter := bson.M{
		"unique_id":      uniqueID,
		"action_status":  statusCheck,
		"action_type":    "DELETE",
		"request_action": "BLOCK_CITY",
	}
	alreadyProcessed, _ := o.MongoDalCPSAction.FindOne(ctx, dupFilter, bson.M{})
	if alreadyProcessed != nil {
		return errors.New("This city action has already been processed")
	}

	update := bson.M{
		"last_modified_at":     time.Now(),
		"checker_id":           checker.UserID,
		"checker_name":         checker.FullName,
		"checker_phone_number": checker.PhoneNumber,
		"checker_action_time":  time.Now(),
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
		return fmt.Errorf("failed to update CPSAction: %w", err)
	}

	if !approve {
		return nil
	}

	var cityCode string
	switch v := actionDoc.CurrentAction.(type) {
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
		return errors.New("city_code missing in action")
	}
	cityCode = strings.TrimSpace(cityCode)
	if cityCode == "" {
		return errors.New("city_code missing in action")
	}

	cityUpdate := bson.M{
		"enabled":    false,
		"updated_at": time.Now(),
	}
	_, err = o.MongoDalCity.UpdateOne(ctx, bson.M{"city_code": cityCode}, cityUpdate)
	if err != nil {
		return fmt.Errorf("failed to disable city %s: %w", cityCode, err)
	}

	return nil
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
        PreviosAction:    prevAction,
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
func (o *outboundAccountBlockStore) ApproveBlockUser(
    ctx context.Context,
    actionID string,
    approve bool,
    reason *string,
    checker action.User,
) error {
    if strings.TrimSpace(actionID) == "" {
        return errors.New("action_id is required")
    }

    filter := bson.M{"action_code": actionID, "request_action": "BLOCK_USER"}
    actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
    if err != nil || actionDoc == nil {
        return errors.New("action not found after update")
    }

    uniqueID := actionDoc.UniqueId
    statusCheck := "APPROVED"
    if !approve {
        statusCheck = "REJECTED"
    }
    dupFilter := bson.M{
        "unique_id":      uniqueID,
        "action_status":  statusCheck,
        "action_type":    "DELETE",
        "request_action": "BLOCK_USER",
    }
    alreadyProcessed, _ := o.MongoDalCPSAction.FindOne(ctx, dupFilter, bson.M{})
    if alreadyProcessed != nil {
        return errors.New("This user block action has already been processed")
    }

    update := bson.M{
        "last_modified_at":     time.Now(),
        "checker_id":           checker.UserID,
        "checker_name":         checker.FullName,
        "checker_phone_number": checker.PhoneNumber,
        "checker_action_time":  time.Now(),
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
        return fmt.Errorf("failed to update CPSAction: %w", err)
    }

    if !approve {
        return nil
    }

    var userCode string
    switch v := actionDoc.CurrentAction.(type) {
    case string:
        userCode = v
    case []byte:
        _ = json.Unmarshal(v, &userCode)
    case bson.Binary:
        _ = json.Unmarshal(v.Data, &userCode)
    }
    userCode = strings.TrimSpace(userCode)
    if userCode == "" {
        return errors.New("user_code missing in action")
    }

    userUpdate := bson.M{
        "is_account_blocked": true,
        "last_modified_at":   time.Now(),
    }
    _, err = o.MongoDalUser.UpdateOne(ctx, bson.M{"user_code": userCode}, userUpdate)
    if err != nil {
        return fmt.Errorf("failed to block user %s: %w", userCode, err)
    }

    return nil
}