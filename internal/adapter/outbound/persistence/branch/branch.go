package persistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/repository"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BranchPersistence struct {
	branchDal dal.MongoDal[entities.Branch, entities.Branch]
	cpsDal    dal.MongoDal[entities.CPSAction, entities.CPSAction]
	logger    utils.Logger
}

var _ repository.BulkCustomerRepo = (*BranchPersistence)(nil)

func NewBranchPersistence(client *mongo.Client, dbName string, 
	branchCollection,cpsCollection string, logger utils.Logger) *BranchPersistence {
	branchDal := dal.NewMongoDal[entities.Branch, entities.Branch](client, dbName, branchCollection)
	cpsDal := dal.NewMongoDal[entities.CPSAction, entities.CPSAction](client, dbName, cpsCollection)
	return &BranchPersistence{
		branchDal: branchDal,
		cpsDal:    cpsDal,
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
func (b *BranchPersistence) DisableSingleBranch(ctx context.Context, branchCode string, cpsData entities.CPSAction) error {
	if branchCode == "" {
		return errors.New("branchCode is required")
	}

	filter := bson.M{
		"action_code":    branchCode,
		"action_status":  entities.ActionPending,
		"action_type":    entities.ActionDisable,
		"request_action": entities.RequestDisableSingleBranch,
	}
	projection := bson.M{"_id": 1}
	existing, err := b.cpsDal.FindOne(ctx, filter, projection)
	if err == nil && existing != nil {
		return errors.New(common.DefineError.General["CONFLICT_KEY"].Message)
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	cpsData.ActionType = entities.ActionDisable
	cpsData.RequestAction = entities.RequestDisableSingleBranch
	cpsData.ActionStatus = entities.ActionPending
	cpsData.ActionCode = branchCode
	cpsData.CreatedAt = time.Now()
	cpsData.LastModifiedAt = time.Now()

	_, err = b.cpsDal.InsertOne(ctx, cpsData)
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
		"action_status":    entities.ActionRejected,
		"last_modified_at": time.Now(),
	}
	if approve {
		update["action_status"] = entities.ActionApproved
	}
	if reason != nil {
		update["rejection_reason"] = *reason
	}
	_, err := b.cpsDal.UpdateOne(ctx, filter, update)
	return err
}
func (b *BranchPersistence) FilterMultipleBranches(ctx context.Context, region, district string) ([]string, error) {
	filter := bson.M{
		"branchRegion": strings.TrimSpace(region),
		"districtName": strings.TrimSpace(district),
	}
	branchPtrs, err := b.branchDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	branches := make([]string, 0, len(branchPtrs))
	for _, ptr := range branchPtrs {
		if ptr != nil {
			branches = append(branches, ptr.BranchCode)
		}
	}
	return branches, nil
}

func (b *BranchPersistence) DisableMultipleBranches(ctx context.Context, branchCodes []string) (*entities.CPSAction, error) {
	if len(branchCodes) == 0 {
		return nil, errors.New("branchCodes are required")
	}
	cpsAction := entities.CPSAction{
		ActionCode:     strings.Join(branchCodes, ","),
		RequestAction:  entities.RequestDisableMultiUsers,
		ActionStatus:   entities.ActionPending,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
	_, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		return nil, err
	}
	return &cpsAction, nil
}

func (b *BranchPersistence) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return errors.New("actionID is required")
	}
	filter := bson.M{"action_code": actionID}
	update := bson.M{
		"action_status":    entities.ActionRejected,
		"last_modified_at": time.Now(),
	}
	if approve {
		update["action_status"] = entities.ActionApproved
	}
	if reason != nil {
		update["rejection_reason"] = *reason
	}
	_, err := b.cpsDal.UpdateOne(ctx, filter, update)
	return err
}
