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

	filter := bson.M{"branchRegion": region, "districtName": district}
	data, err := o.MongoDalBranch.FindAll(ctx, filter, nil)
	if err != nil {
		fmt.Printf("failed to fetch branches: %v\n", err)
		return nil, fmt.Errorf("failed to fetch branches: %w", err)
	}

	var branches []model.Branch
	for _, b := range data {
		if b == nil {
			fmt.Println("Skipping nil branch document")
			continue
		}
		branch := model.Branch{
			ID:            b.ID,
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
		if branch.ID.Hex() == "" {
			fmt.Printf("Warning: Branch has empty ID or unmapped fields: %+v\n", b)
		}
		branches = append(branches, branch)
	}

	if len(branches) == 0 {
	} else {
	}

	var result []action.Branch
	for _, b := range branches {
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
func (o *outboundAccountBlockStore) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error {
	department, _ := ctx.Value("department").(string)
	if strings.TrimSpace(department) == "" {
		return errors.New("department is required in context")
	}

	filter := bson.M{
		"department":                 department,
		"maker_user.user_code":       maker.UserID,
		"action_status":              model.ActionPending,
		"action_type":                model.ActionDelete,
		"request_action":             model.RequestDisableSingleBranch,
		"current_action.branch_code": branch.BranchCode,
	}

	existing, err := o.MongoDalCPSAction.FindOne(ctx, filter, bson.M{})
	if err == nil && existing != nil {
		return fmt.Errorf("pending disable action already exists for this branch and maker")
	}

	prevBranchPtr, err := o.MongoDalBranch.FindOne(ctx, bson.M{"branch_code": branch.BranchCode}, bson.M{})
	var prevAction json.RawMessage
	if err == nil && prevBranchPtr != nil {
		prevAction, _ = json.Marshal(prevBranchPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(branch)

	cpsAction := model.CPSAction{
		ActionCode: utils.RandomGenerator(24),
		MakerUser: model.User{
			UserCode:    maker.UserID,
			FullName:    maker.FullName,
			PhoneNumber: maker.PhoneNumber},
		Department:     department,
		ActionStatus:   model.ActionPending,
		ActionType:     model.ActionDelete,
		RequestAction:  model.RequestDisableSingleBranch,
		PreviosAction:  prevAction,
		CurrentAction:  currAction,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	return err
}
func (o *outboundAccountBlockStore) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	filter := bson.M{"action_code": actionID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found")
	}
	if !approve {
		update := bson.M{"$set": bson.M{"action_status": "REJECTED", "last_modified_at": time.Now()}}
		_, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
		return err
	}
	var branch action.Branch
	_ = json.Unmarshal(actionDoc.CurrentAction, &branch)
	update := bson.M{"$set": bson.M{"enabled": false, "updated_at": time.Now()}}
	_, err = o.MongoDalBranch.UpdateOne(ctx, bson.M{"branch_code": branch.BranchCode}, update)
	if err != nil {
		return err
	}
	updateAction := bson.M{"$set": bson.M{"action_status": "APPROVED", "last_modified_at": time.Now()}}
	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, updateAction)
	return err
}

func (o *outboundAccountBlockStore) FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	filter := bson.M{"branch_region": region}
	if district != "" {
		filter["district_name"] = district
	}
	data, err := o.MongoDalBranch.FindAll(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	var result []action.Branch
	for _, b := range data {
		if b == nil {
			continue
		}
		result = append(result, action.Branch{
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
		})
	}
	return result, nil
}
func (o *outboundAccountBlockStore) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error {
	department, _ := ctx.Value("department").(string)
	if strings.TrimSpace(department) == "" {
		return errors.New("department is required in context")
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
	cpsAction := model.CPSAction{
		ActionCode: utils.RandomGenerator(24),
		MakerUser: model.User{
			UserCode:    maker.UserID,
			FullName:    maker.FullName,
			PhoneNumber: maker.PhoneNumber},
		Department:     department,
		ActionStatus:   model.ActionPending,
		ActionType:     model.ActionDelete,
		RequestAction:  model.RequestDisableMultiBranches,
		PreviosAction:  prevAction,
		CurrentAction:  currAction,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
	_, err := o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	return err
}

func (o *outboundAccountBlockStore) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	filter := bson.M{"action_code": actionID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found")
	}
	if !approve {
		update := bson.M{"$set": bson.M{"action_status": "REJECTED", "last_modified_at": time.Now()}}
		_, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
		return err
	}
	var branches []action.Branch
	_ = json.Unmarshal(actionDoc.CurrentAction, &branches)
	for _, branch := range branches {
		update := bson.M{"$set": bson.M{"enabled": false, "updated_at": time.Now()}}
		_, err := o.MongoDalBranch.UpdateOne(ctx, bson.M{"branch_code": branch.BranchCode}, update)
		if err != nil {
			return err
		}
	}
	updateAction := bson.M{"$set": bson.M{"action_status": "APPROVED", "last_modified_at": time.Now()}}
	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, updateAction)
	return err
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
func (o *outboundAccountBlockStore) BlockRegion(ctx context.Context, regionID string, maker action.CPSAction) error {
	department, _ := ctx.Value("department").(string)
	if strings.TrimSpace(department) == "" {
		return errors.New("department is required in context")
	}

	prevRegionPtr, err := o.MongoDalRegion.FindOne(ctx, bson.M{"id": regionID}, bson.M{})
	var prevAction json.RawMessage
	if err == nil && prevRegionPtr != nil {
		prevAction, _ = json.Marshal(prevRegionPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(regionID)

	cpsAction := model.CPSAction{
		ActionCode: utils.RandomGenerator(24),
		MakerUser: model.User{
			UserCode:    maker.Maker.UserID,
			FullName:    maker.Maker.FullName,
			PhoneNumber: maker.Maker.PhoneNumber,
		},
		CheckerUser: model.User{
			UserCode:    maker.Checker.UserID,
			FullName:    maker.Checker.FullName,
			PhoneNumber: maker.Checker.PhoneNumber,
		},
		UniqueId:        regionID,
		Department:      department,
		RejectionReason: maker.RejectionReason,
		PreviosAction:   prevAction,
		CurrentAction:   currAction,
		ActionStatus:    model.ActionPending,
		ActionType:      model.ActionDelete,
		RequestAction:   model.RequestBlockRegion,
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	return err
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
func (o *outboundAccountBlockStore) ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	filter := bson.M{"action_code": actionID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found")
	}

	checkerInfo := bson.M{
		"checker_code": checker.UserCode,
		"checker_name": checker.FullName,
		"checked_at":   time.Now(),
	}

	if !approve {
		update := bson.M{
			"$set": bson.M{
				"action_status":    "REJECTED",
				"last_modified_at": time.Now(),
				"reason":           reason,
			},
			"$push": bson.M{
				"checker_history": checkerInfo,
			},
		}
		_, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
		return err
	}

	var regionID string
	if err := json.Unmarshal(actionDoc.CurrentAction, &regionID); err != nil {
		return errors.New("failed to parse regionID from action")
	}

	updateRegion := bson.M{
		"$set": bson.M{
			"enabled":    false,
			"updated_at": time.Now(),
		},
	}
	_, err = o.MongoDalRegion.UpdateOne(ctx, bson.M{"id": regionID}, updateRegion)
	if err != nil {
		return err
	}

	updateAction := bson.M{
		"$set": bson.M{
			"action_status":    "APPROVED",
			"last_modified_at": time.Now(),
			"reason":           reason,
		},
		"$push": bson.M{
			"checker_history": checkerInfo,
		},
	}
	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, updateAction)
	return err
}
func (o *outboundAccountBlockStore) GetRegionByID(ctx context.Context, regionID string) (action.Region, error) {
	if regionID == "" {
		return action.Region{}, fmt.Errorf("regionID is required")
	}
	filter := bson.M{"id": regionID}
	regionDoc, err := o.MongoDalRegion.FindOne(ctx, filter, nil)
	if err != nil || regionDoc == nil {
		return action.Region{}, fmt.Errorf("failed to get region by ID: %w", err)
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

func (o *outboundAccountBlockStore) BlockDistrict(ctx context.Context, districtID string, maker action.CPSAction) error {
	department, _ := ctx.Value("department").(string)
	if strings.TrimSpace(department) == "" {
		return errors.New("department is required in context")
	}

	prevDistrictPtr, err := o.MongoDalDistrict.FindOne(ctx, bson.M{"id": districtID}, bson.M{})
	var prevAction json.RawMessage
	if err == nil && prevDistrictPtr != nil {
		prevAction, _ = json.Marshal(prevDistrictPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(districtID)

	cpsAction := model.CPSAction{
		ActionCode: utils.RandomGenerator(24),
		MakerUser: model.User{
			UserCode:    maker.Maker.UserID,
			FullName:    maker.Maker.FullName,
			PhoneNumber: maker.Maker.PhoneNumber},
		CheckerUser: model.User{
			UserCode:    maker.Checker.UserID,
			FullName:    maker.Checker.FullName,
			PhoneNumber: maker.Checker.PhoneNumber},
		UniqueId:        districtID,
		Department:      department,
		RejectionReason: maker.RejectionReason,
		PreviosAction:   prevAction,
		CurrentAction:   currAction,
		ActionStatus:    model.ActionPending,
		ActionType:      model.ActionDelete,
		RequestAction:   model.RequestAction(model.RequestBlockDistrict),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	return err
}
func (o *outboundAccountBlockStore) GetDistrictByID(ctx context.Context, districtID string) (action.District, error) {
	if districtID == "" {
		return action.District{}, fmt.Errorf("districtID is required")
	}
	filter := bson.M{"id": districtID}
	districtDoc, err := o.MongoDalDistrict.FindOne(ctx, filter, nil)
	if err != nil || districtDoc == nil {
		return action.District{}, fmt.Errorf("failed to get district by ID: %w", err)
	}
	return action.District{
		ID:           districtDoc.ID.Hex(),
		DistrictCode: districtDoc.DistrictCode,
		DistrictName: districtDoc.DistrictName,
		CreatedAt:    districtDoc.CreatedAt,
		UpdatedAt:    districtDoc.UpdatedAt,
		Enabled:      districtDoc.Enabled,
	}, nil
}
func (o *outboundAccountBlockStore) ApproveBlockDistrict(ctx context.Context, districtID string, checker action.CPSAction) error {
	filter := bson.M{"action_code": districtID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found")
	}

	checkerInfo := bson.M{
		"checker_code": checker.Checker.UserID,
		"checker_name": checker.Checker.FullName,
		"checked_at":   time.Now(),
	}

	update := bson.M{
		"$set": bson.M{
			"action_status":    "APPROVED",
			"last_modified_at": time.Now(),
		},
		"$push": bson.M{
			"checker_history": checkerInfo,
		},
	}
	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	return err
}

func (o *outboundAccountBlockStore) BlockCity(ctx context.Context, cityID string, maker action.CPSAction) error {
	department, _ := ctx.Value("department").(string)
	if strings.TrimSpace(department) == "" {
		return errors.New("department is required in context")
	}

	prevCityPtr, err := o.MongoDalCity.FindOne(ctx, bson.M{"id": cityID}, bson.M{})
	var prevAction json.RawMessage
	if err == nil && prevCityPtr != nil {
		prevAction, _ = json.Marshal(prevCityPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(cityID)

	cpsAction := model.CPSAction{
		ActionCode: utils.RandomGenerator(24),
		MakerUser: model.User{
			UserCode: maker.Maker.UserID, FullName: maker.Maker.FullName,
			PhoneNumber: maker.Maker.PhoneNumber},
		CheckerUser: model.User{
			UserCode:    maker.Checker.UserID,
			FullName:    maker.Checker.FullName,
			PhoneNumber: maker.Checker.PhoneNumber},
		UniqueId:        cityID,
		Department:      department,
		RejectionReason: maker.RejectionReason,
		PreviosAction:   prevAction,
		CurrentAction:   currAction,
		ActionStatus:    model.ActionPending,
		ActionType:      model.ActionDelete,
		RequestAction:   model.RequestAction(model.RequestBlockCity),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	return err
}
func (o *outboundAccountBlockStore) GetCityByID(ctx context.Context, cityID string) (action.City, error) {
	if cityID == "" {
		return action.City{}, fmt.Errorf("cityID is required")
	}
	filter := bson.M{"id": cityID}
	cityDoc, err := o.MongoDalCity.FindOne(ctx, filter, nil)
	if err != nil || cityDoc == nil {
		return action.City{}, fmt.Errorf("failed to get city by ID: %w", err)
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
func (o *outboundAccountBlockStore) ApproveBlockCity(ctx context.Context, cityID string, checker action.CPSAction) error {
	filter := bson.M{"action_code": cityID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found")
	}

	checkerInfo := bson.M{
		"checker_code": checker.Checker.UserID,
		"checker_name": checker.Checker.FullName,
		"checked_at":   time.Now(),
	}

	update := bson.M{
		"$set": bson.M{
			"action_status":    "APPROVED",
			"last_modified_at": time.Now(),
		},
		"$push": bson.M{
			"checker_history": checkerInfo,
		},
	}
	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	return err
}

func (o *outboundAccountBlockStore) BlockUser(ctx context.Context, userID string, maker action.CPSAction) error {
	department, _ := ctx.Value("department").(string)
	if strings.TrimSpace(department) == "" {
		return errors.New("department is required in context")
	}

	prevUserPtr, err := o.MongoDalUser.FindOne(ctx, bson.M{"user_id": userID}, bson.M{})
	var prevAction json.RawMessage
	if err == nil && prevUserPtr != nil {
		prevAction, _ = json.Marshal(prevUserPtr)
	} else {
		prevAction = json.RawMessage("null")
	}

	currAction, _ := json.Marshal(userID)

	cpsAction := model.CPSAction{
		ActionCode: utils.RandomGenerator(24),
		MakerUser: model.User{
			UserCode:    maker.Maker.UserID,
			FullName:    maker.Maker.FullName,
			PhoneNumber: maker.Maker.PhoneNumber},
		CheckerUser: model.User{
			UserCode:    maker.Checker.UserID,
			FullName:    maker.Checker.FullName,
			PhoneNumber: maker.Checker.PhoneNumber},
		UniqueId:        userID,
		Department:      department,
		RejectionReason: maker.RejectionReason,
		PreviosAction:   prevAction,
		CurrentAction:   currAction,
		ActionStatus:    model.ActionPending,
		ActionType:      model.ActionDelete,
		RequestAction:   model.RequestBlockUser,
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	_, err = o.MongoDalCPSAction.InsertOne(ctx, cpsAction)
	return err
}

func (o *outboundAccountBlockStore) GetUserByID(ctx context.Context, userID string, maker action.CPSAction) (member.User, error) {
	if userID == "" {
		return member.User{}, fmt.Errorf("userID is required")
	}
	filter := bson.M{"user_id": userID}
	userDoc, err := o.MongoDalUser.FindOne(ctx, filter, nil)
	if err != nil || userDoc == nil {
		return member.User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return *userDoc, nil
}
func (o *outboundAccountBlockStore) ApproveBlockUser(ctx context.Context, userID string, checker action.CPSAction) error {
	filter := bson.M{"action_code": userID}
	actionDoc, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil || actionDoc == nil {
		return errors.New("action not found")
	}

	checkerInfo := bson.M{
		"checker_code": checker.Checker.UserID,
		"checker_name": checker.Checker.FullName,
		"checked_at":   time.Now(),
	}

	update := bson.M{
		"$set": bson.M{
			"action_status":    "APPROVED",
			"last_modified_at": time.Now(),
		},
		"$push": bson.M{
			"checker_history": checkerInfo,
		},
	}
	_, err = o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	return err
}
