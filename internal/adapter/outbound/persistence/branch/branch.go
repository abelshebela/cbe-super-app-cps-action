package persistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bulkcustomer/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bulkcustomer/repository"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BranchPersistence struct {
	client    *mongo.Client
	branchDal dal.MongoDal[entities.Branch, entities.Branch]
	cpsDal    dal.MongoDal[entities.CPSAction, entities.CPSAction]
	timeout   time.Duration
	logger    utils.Logger
}

var _ repository.BulkCustomerRepo = (*BranchPersistence)(nil)

func NewBranchPersistence(client *mongo.Client, dbName string, collectionNames []string, timeout time.Duration, logger utils.Logger) *BranchPersistence {
	branchDal := dal.NewMongoDal[entities.Branch, entities.Branch](client, dbName, collectionNames[0])
	cpsDal := dal.NewMongoDal[entities.CPSAction, entities.CPSAction](client, dbName, collectionNames[1])
	return &BranchPersistence{
		client:    client,
		branchDal: branchDal,
		cpsDal:    cpsDal,
		timeout:   timeout,
		logger:    logger,
	}
}

func (b *BranchPersistence) FilterSingleBranches(ctx context.Context, region, district string) ([]entities.Branch, error) {
	filter := bson.M{
		"branchRegion": strings.TrimSpace(region),
		"districtName": strings.TrimSpace(district),
	}
	branchPtrs, err := b.branchDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	branches := make([]entities.Branch, 0, len(branchPtrs))
	for _, ptr := range branchPtrs {
		if ptr != nil {
			branches = append(branches, *ptr)
		}
	}
	return branches, nil
}

func (b *BranchPersistence) DisableSingleBranch(ctx context.Context, branch entities.Branch, maker entities.User) error {
	actionCode := utils.RandomGenerator(24)
	actionData := entities.ActionData{
		BranchCode:    branch.BranchCode,
		BranchName:    branch.BranchName,
		BranchAddress: branch.BranchAddress,
		DistrictCode:  branch.DistrictCode,
		DistrictName:  branch.DistrictName,
		BranchRegion:  branch.BranchRegion,
	}

	cpsAction := entities.CPSAction{
		ID:              bson.NewObjectID().Hex(),
		ActionCode:      actionCode,
		MakerUser:       maker,
		CheckerUser:     entities.User{},
		RejectedReason:  "",
		Department:      "",
		Status:          entities.ActionPending,
		PreviousAction:  nil,
		RequestAction:   entities.RequestDisableSingleBranch,
		ActionType:      entities.ActionUpdate,
		ActionData:      actionData,
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	filter := bson.M{
		"action_data.branch_code": branch.BranchCode,
		"status":                  entities.ActionPending,
		"action_type":             entities.ActionUpdate,
		"request_action":          entities.RequestDisableSingleBranch,
	}
	projection := bson.M{"_id": 1}
	existing, err := b.cpsDal.FindOne(ctx, filter, projection)
	if err == nil && existing != nil {
		return errors.New(common.DefineError.General["CONFLICT_KEY"].Message)
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	_, err = b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) || (err != nil && strings.Contains(err.Error(), "E11000")) {
			return errors.New(common.DefineError.General["CONFLICT_KEY"].Message)
		}
		return err
	}
	return nil
}

func (b *BranchPersistence) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return errors.New("actionID is required")
	}
	filter := bson.M{"action_code": actionID}
	update := bson.M{
		"last_modified_at": time.Now(),
	}
	if approve {
		update["status"] = entities.ActionApproved
	} else {
		update["status"] = entities.ActionRejected
	}
	if reason != nil {
		update["rejection_reason"] = *reason
	}
	_, err := b.cpsDal.UpdateOne(ctx, filter, update)
	return err
}
func (b *BranchPersistence) FilterMultipleBranches(ctx context.Context, region, district string) ([]entities.Branch, error) {
	filter := bson.M{
		"branchRegion": strings.TrimSpace(region),
		"districtName": strings.TrimSpace(district),
	}
	branchPtrs, err := b.branchDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	branches := make([]entities.Branch, 0, len(branchPtrs))
	for _, ptr := range branchPtrs {
		if ptr != nil {
			branches = append(branches, *ptr)
		}
	}
	return branches, nil
}
func (b *BranchPersistence) DisableMultipleBranches(ctx context.Context, branches []entities.Branch, maker entities.User) error {
	if len(branches) == 0 {
		return errors.New("branches are required")
	}

	firstBranch := branches[0]
	actionData := entities.ActionData{
		BranchCode:    firstBranch.BranchCode,
		BranchName:    firstBranch.BranchName,
		BranchAddress: firstBranch.BranchAddress,
		DistrictCode:  firstBranch.DistrictCode,
		DistrictName:  firstBranch.DistrictName,
		BranchRegion:  firstBranch.BranchRegion,
	}

	var branchCodes []string
	for _, branch := range branches {
		branchCodes = append(branchCodes, branch.BranchCode)
	}

	actionCode := utils.RandomGenerator(24)

	cpsAction := entities.CPSAction{
		ID:              bson.NewObjectID().Hex(),
		ActionCode:      actionCode,
		MakerUser:       maker,
		CheckerUser:     entities.User{},
		RejectedReason:  "",
		Department:      "",
		Status:          entities.ActionPending,
		PreviousAction:  nil,
		RequestAction:   entities.RequestDisableMultiUsers,
		ActionType:      entities.ActionUpdate,
		ActionData:      actionData,
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	filter := bson.M{
		"action_data.branch_code": bson.M{"$in": branchCodes},
		"status":                  entities.ActionPending,
		"action_type":             entities.ActionUpdate,
		"request_action":          entities.RequestDisableMultiUsers,
	}
	projection := bson.M{"_id": 1}
	existing, err := b.cpsDal.FindOne(ctx, filter, projection)
	if err == nil && existing != nil {
		return errors.New(common.DefineError.General["CONFLICT_KEY"].Message)
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	_, err = b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) || (err != nil && strings.Contains(err.Error(), "E11000")) {
			return errors.New(common.DefineError.General["CONFLICT_KEY"].Message)
		}
		return err
	}
	return nil
}
func (b *BranchPersistence) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return errors.New("actionID is required")
	}
	filter := bson.M{"action_code": actionID}
	update := bson.M{
		"last_modified_at": time.Now(),
	}
	if approve {
		update["status"] = entities.ActionApproved
	} else {
		update["status"] = entities.ActionRejected
	}
	if reason != nil {
		update["rejection_reason"] = *reason
	}
	_, err := b.cpsDal.UpdateOne(ctx, filter, update)
	return err
}
func (b *BranchPersistence) GetBranchByCode(ctx context.Context, branchCode string) (entities.Branch, error) {
	filter := bson.M{"branchCode": branchCode}
	result, err := b.branchDal.FindOne(ctx, filter, nil)
	if err != nil {
		return entities.Branch{}, err
	}
	if result == nil {
		return entities.Branch{}, mongo.ErrNoDocuments
	}
	return *result, nil
}
