package persistence

import (
    "context"
    "errors"
    "strings"
    "time"

    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/entities"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/repository"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

type BranchPersistence struct {
    branchDal dal.MongoDal[entities.Branch, entities.Branch]
    cpsDal    dal.MongoDal[entities.CPSAction, entities.CPSAction]
    timeout   time.Duration
}

var _ repository.BulkCustomerRepo = (*BranchPersistence)(nil)

func NewBranchPersistence(client *mongo.Client, dbName string, timeout time.Duration) *BranchPersistence {
    branchDal := dal.NewMongoDal[entities.Branch, entities.Branch](client, dbName, "branches")
    cpsDal := dal.NewMongoDal[entities.CPSAction, entities.CPSAction](client, dbName, "cps_actions")
    return &BranchPersistence{
        branchDal: branchDal,
        cpsDal:    cpsDal,
        timeout:   timeout,
    }
}

func (b *BranchPersistence) FilterSingleBranches(ctx context.Context, region, district string) ([]string, error) {
    filter := bson.M{
        "branchRegion":  strings.TrimSpace(region),
        "districtName":  strings.TrimSpace(district),
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
func (b *BranchPersistence) DisableSingleBranch(ctx context.Context, branchCode, cpsData string) (*entities.CPSAction, error) {
    if branchCode == "" || cpsData == "" {
        return nil, errors.New("branchCode and cpsData are required")
    }
    cpsAction := entities.CPSAction{
        ActionCode:    branchCode,
        RequestAction: entities.RequestDisableSingleBranch,
        ActionStatus:  entities.ActionPending,
        CreatedAt:     time.Now(),
        LastModifiedAt: time.Now(),
    }
    _, err := b.cpsDal.InsertOne(ctx, cpsAction)
    if err != nil {
        return nil, err
    }
    return &cpsAction, nil
}

func (b *BranchPersistence) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
    if actionID == "" {
        return errors.New("actionID is required")
    }
    filter := bson.M{"action_code": actionID}
    update := bson.M{
        "action_status": entities.ActionRejected,
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
        "branchRegion":  strings.TrimSpace(region),
        "districtName":  strings.TrimSpace(district),
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

func (b *BranchPersistence) DisableMultipleBranches(ctx context.Context, branchCodes []string, cpsData string) (*entities.CPSAction, error) {
    if len(branchCodes) == 0 || cpsData == "" {
        return nil, errors.New("branchCodes and cpsData are required")
    }
    cpsAction := entities.CPSAction{
        ActionCode:    strings.Join(branchCodes, ","),
        RequestAction: entities.RequestDisableMultiUsers,
        ActionStatus:  entities.ActionPending,
        CreatedAt:     time.Now(),
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
        "action_status": entities.ActionRejected,
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