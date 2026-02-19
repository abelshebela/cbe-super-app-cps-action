package action_role_repo

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BPSActionApproveIndexRepository struct {
	client     *mongo.Client
	logger     utils.Logger
	collection *mongo.Collection
}

func NewBPSActionApproveIndexRepository(client *mongo.Client, database string, collection string, logger utils.Logger) storage.BPSActionApproveIndexRepository {
	return &BPSActionApproveIndexRepository{
		client:     client,
		logger:     logger,
		collection: client.Database(database).Collection(collection),
	}
}

func (r *BPSActionApproveIndexRepository) PopulateUserApproverAllocations(ctx context.Context, role_id string) ([]string, []string, []string, error) {
	r.logger.Infof("PopulateUserApproverAllocations: Populating approver allocations for RoleID: %s", role_id)
	makerSet := map[string]struct{}{}
	checkerSet := map[string]struct{}{}
	auditorSet := map[string]struct{}{}

	cursor, err := r.collection.Find(ctx, bson.M{
		"role_id": role_id,
	})
	if err != nil {
		r.logger.Errorf("PopulateUserApproverAllocations: Find failed: %v", err)
		return nil, nil, nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)
	var results []imodel.BPSActionApproveIndex
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("PopulateUserApproverAllocations: Cursor.All failed: %v", err)
		return nil, nil, nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	for _, v := range results {
		actionName := v.ActionName
		if v.MakerIndex != nil {
			makerSet[actionName] = struct{}{}
		}
		if v.CheckerIndex != nil {
			checkerSet[actionName] = struct{}{}
		}
		if v.AuditorIndex != nil {
			auditorSet[actionName] = struct{}{}
		}
	}

	makerAllocations := make([]string, 0, len(makerSet))
	for k := range makerSet {
		makerAllocations = append(makerAllocations, k)
	}
	checkerAllocations := make([]string, 0, len(checkerSet))
	for k := range checkerSet {
		checkerAllocations = append(checkerAllocations, k)
	}
	auditorAllocations := make([]string, 0, len(auditorSet))
	for k := range auditorSet {
		auditorAllocations = append(auditorAllocations, k)
	}

	r.logger.Infof("PopulateUserApproverAllocations: Found %d maker, %d checker, %d auditor allocations for RoleID: %s", len(makerAllocations), len(checkerAllocations), len(auditorAllocations), role_id)
	return makerAllocations, checkerAllocations, auditorAllocations, nil
}
func (r *BPSActionApproveIndexRepository) SaveIndices(ctx context.Context, indices []imodel.BPSActionApproveIndex) error {
	r.logger.Infof("SaveIndices: Saving %d indices to DB: %s, Collection: %s", len(indices), r.collection.Database().Name(), r.collection.Name())
	if len(indices) == 0 {
		return nil
	}
	docs := make([]interface{}, len(indices))
	for i, v := range indices {
		docs[i] = v
	}
	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		r.logger.Errorf("SaveIndices: InsertMany failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	r.logger.Infof("SaveIndices: Successfully saved %d indices", len(indices))
	return nil
}

func (r *BPSActionApproveIndexRepository) SyncIndices(ctx context.Context, oldActionName string, newIndices []imodel.BPSActionApproveIndex) error {
	r.logger.Infof("SyncIndices: Syncing %d indices for oldActionName: %s", len(newIndices), oldActionName)

	if _, err := r.collection.DeleteMany(ctx, bson.M{"action_name": oldActionName}); err != nil {
		r.logger.Errorf("SyncIndices: DeleteMany failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if err := r.SaveIndices(ctx, newIndices); err != nil {
		return err
	}

	return nil
}

func (r *BPSActionApproveIndexRepository) FindMakerAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.BPSActionApproveIndex, error) {
	r.logger.Infof("FindMakerAllocationsByRoleID: Finding maker allocations for RoleID: %s", roleID.Hex())

	cursor, err := r.collection.Find(ctx, bson.M{
		"role_id":     roleID.Hex(),
		"maker_index": bson.M{"$ne": nil},
	})
	if err != nil {
		r.logger.Errorf("FindMakerAllocationsByRoleID: Find failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var results []imodel.BPSActionApproveIndex
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("FindMakerAllocationsByRoleID: Cursor.All failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	r.logger.Infof("FindMakerAllocationsByRoleID: Found %d maker allocations for RoleID: %s", len(results), roleID.Hex())
	return results, nil
}

func (r *BPSActionApproveIndexRepository) FindCheckerAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.BPSActionApproveIndex, error) {
	r.logger.Infof("FindCheckerAllocationsByRoleID: Finding checker allocations for RoleID: %s", roleID.Hex())

	cursor, err := r.collection.Find(ctx, bson.M{
		"role_id":       roleID.Hex(),
		"checker_index": bson.M{"$ne": nil},
	})
	if err != nil {
		r.logger.Errorf("FindCheckerAllocationsByRoleID: Find failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var results []imodel.BPSActionApproveIndex
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("FindCheckerAllocationsByRoleID: Cursor.All failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	r.logger.Infof("FindCheckerAllocationsByRoleID: Found %d checker allocations for RoleID: %s", len(results), roleID.Hex())
	return results, nil
}

func (r *BPSActionApproveIndexRepository) FindAuditorAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.BPSActionApproveIndex, error) {
	r.logger.Infof("FindAuditorAllocationsByRoleID: Finding auditor allocations for RoleID: %s", roleID.Hex())

	cursor, err := r.collection.Find(ctx, bson.M{
		"role_id":       roleID.Hex(),
		"auditor_index": bson.M{"$ne": nil},
	})
	if err != nil {
		r.logger.Errorf("FindAuditorAllocationsByRoleID: Find failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var results []imodel.BPSActionApproveIndex
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("FindAuditorAllocationsByRoleID: Cursor.All failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	r.logger.Infof("FindAuditorAllocationsByRoleID: Found %d auditor allocations for RoleID: %s", len(results), roleID.Hex())
	return results, nil
}
func (r *BPSActionApproveIndexRepository) InsertMany(
	ctx context.Context,
	makerIndex []bson.ObjectID,
	checkerIndex [][]bson.ObjectID,
	auditorIndex []bson.ObjectID,
	roleCode string,
) error {

	r.logger.Infof("InsertMany: Insert or update indices for RoleCode: %s", roleCode)

	var models []mongo.WriteModel
	now := time.Now()

	// Makers
	for _, id := range makerIndex {
		models = append(models,
			mongo.NewUpdateOneModel().
				SetFilter(bson.M{
					"role_id":     id.Hex(),
					"action_name": roleCode,
				}).
				SetUpdate(bson.M{
					"$set": bson.M{
						"maker_index": 1,
						"updated_at":  now,
					},
					"$setOnInsert": bson.M{
						"created_at": now,
					},
				}).
				SetUpsert(true),
		)
	}

	// Checkers
	for i, ids := range checkerIndex {
		for _, id := range ids {
			models = append(models,
				mongo.NewUpdateOneModel().
					SetFilter(bson.M{
						"role_id":     id.Hex(),
						"action_name": roleCode,
					}).
					SetUpdate(bson.M{
						"$set": bson.M{
							"checker_index": i + 1,
							"updated_at":    now,
						},
						"$setOnInsert": bson.M{
							"created_at": now,
						},
					}).
					SetUpsert(true),
			)
		}
	}

	// Auditors
	for _, id := range auditorIndex {
		models = append(models,
			mongo.NewUpdateOneModel().
				SetFilter(bson.M{
					"role_id":     id.Hex(),
					"action_name": roleCode,
				}).
				SetUpdate(bson.M{
					"$set": bson.M{
						"auditor_index": 1,
						"updated_at":    now,
					},
					"$setOnInsert": bson.M{
						"created_at": now,
					},
				}).
				SetUpsert(true),
		)
	}

	if len(models) == 0 {
		return nil
	}

	_, err := r.collection.BulkWrite(ctx, models)
	if err != nil {
		r.logger.Errorf("InsertMany: BulkWrite failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	r.logger.Infof("InsertMany: Successfully processed indices for RoleCode: %s", roleCode)
	return nil
}

func (r *BPSActionApproveIndexRepository) DeleteMany(
	ctx context.Context,
	makerIndex []bson.ObjectID,
	checkerIndex [][]bson.ObjectID,
	auditorIndex []bson.ObjectID,
	roleCode string,
) error {

	r.logger.Infof("DeleteMany: Deleting indices for RoleCode: %s", roleCode)

	var models []mongo.WriteModel

	// 1. Delete makers (only if NOT checker or auditor)
	if len(makerIndex) > 0 {
		makerIDs := make([]string, 0, len(makerIndex))
		for _, id := range makerIndex {
			makerIDs = append(makerIDs, id.Hex())
		}
		models = append(models,
			mongo.NewDeleteManyModel().SetFilter(bson.M{
				"role_id":     bson.M{"$in": makerIDs},
				"action_name": roleCode,

				// must be pure maker
				"checker_index": bson.M{"$eq": nil},
				"auditor_index": bson.M{"$eq": nil},
			}),
		)
	}

	// 2. Delete checkers (only if NOT maker or auditor)
	for i, ids := range checkerIndex {
		if len(ids) == 0 {
			continue
		}
		checkerIDs := make([]string, 0, len(ids))
		for _, id := range ids {
			checkerIDs = append(checkerIDs, id.Hex())
		}

		models = append(models,
			mongo.NewDeleteManyModel().SetFilter(bson.M{
				"role_id":       bson.M{"$in": checkerIDs},
				"action_name":   roleCode,
				"checker_index": i + 1, // 1-based index

				// must be pure checker
				"maker_index":   bson.M{"$eq": nil},
				"auditor_index": bson.M{"$eq": nil},
			}),
		)
	}

	// 3. Delete auditors (only if NOT maker or checker)
	if len(auditorIndex) > 0 {
		auditorIDs := make([]string, 0, len(auditorIndex))
		for _, id := range auditorIndex {
			auditorIDs = append(auditorIDs, id.Hex())
		}
		models = append(models,
			mongo.NewDeleteManyModel().SetFilter(bson.M{
				"role_id":     bson.M{"$in": auditorIDs},
				"action_name": roleCode,

				// must be pure auditor
				"maker_index":   bson.M{"$eq": nil},
				"checker_index": bson.M{"$eq": nil},
			}),
		)
	}

	if len(models) == 0 {
		r.logger.Infof("DeleteMany: No delete models generated for RoleCode: %s", roleCode)
		return nil
	}

	_, err := r.collection.BulkWrite(ctx, models)
	if err != nil {
		r.logger.Errorf("DeleteMany: BulkWrite failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	r.logger.Infof("DeleteMany: Successfully deleted indices for RoleCode: %s", roleCode)
	return nil
}

func (r *BPSActionApproveIndexRepository) InsertAll(ctx context.Context, new imodel.BPSActionApproveIndex) error {
	if new.ID.IsZero() {
		new.ID = bson.NewObjectID()
	}
	if new.CreatedAt.IsZero() {
		now := time.Now()
		new.CreatedAt = now
		new.UpdatedAt = now
	}
	_, err := r.collection.InsertOne(ctx, new)
	if err != nil {
		r.logger.Errorf("InsertAll: InsertOne failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *BPSActionApproveIndexRepository) DeleteAll(ctx context.Context, prev imodel.BPSActionApproveIndex) error {
	filter := bson.M{}
	if !prev.ID.IsZero() {
		filter["_id"] = prev.ID
	} else {
		filter["role_id"] = prev.RoleId
		filter["action_name"] = prev.ActionName
	}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.Errorf("DeleteAll: DeleteMany failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
