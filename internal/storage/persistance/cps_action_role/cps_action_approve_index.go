package action_role_repo

import (
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSActionApproveIndexRepository struct {
	client     *mongo.Client
	logger     utils.Logger
	collection *mongo.Collection
}

func NewCPSActionApproveIndexRepository(client *mongo.Client, database string, collection string, logger utils.Logger) storage.CPSActionApproveIndexRepository {
	return &CPSActionApproveIndexRepository{
		client:     client,
		logger:     logger,
		collection: client.Database(database).Collection(collection),
	}
}

func (r *CPSActionApproveIndexRepository) SaveIndices(ctx context.Context, indices []model.CPSActionApproveIndex) error {
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
		return err
	}
	r.logger.Infof("SaveIndices: Successfully saved %d indices", len(indices))
	return nil
}

func (r *CPSActionApproveIndexRepository) SyncIndices(ctx context.Context, oldActionName string, newIndices []model.CPSActionApproveIndex) error {
	r.logger.Infof("SyncIndices: Syncing %d indices for oldActionName: %s", len(newIndices), oldActionName)
	cursor, err := r.collection.Find(ctx, bson.M{"action_name": oldActionName})
	if err != nil {
		return err
	}
	var oldIndices []model.CPSActionApproveIndex
	if err := cursor.All(ctx, &oldIndices); err != nil {
		return err
	}

	oldMap := make(map[string]model.CPSActionApproveIndex)
	for _, idx := range oldIndices {
		oldMap[idx.RoleId] = idx
	}
	var writes []mongo.WriteModel

	for _, newIdx := range newIndices {
		if oldIdx, exists := oldMap[newIdx.RoleId]; exists {
			// Update existing index
			update := bson.M{
				"action_name":   newIdx.ActionName,
				"maker_index":   newIdx.MakerIndex,
				"checker_index": newIdx.CheckerIndex,
				"auditor_index": newIdx.AuditorIndex,
				"updated_at":    time.Now(),
			}
			r.logger.Infof("SyncIndices: Updating RoleID %s with %v", newIdx.RoleId, update)
			writes = append(writes, mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": oldIdx.ID}).
				SetUpdate(bson.M{"$set": update}))
			delete(oldMap, newIdx.RoleId)
		} else {
			// Insert new index
			r.logger.Infof("SyncIndices: Inserting new index for RoleID %s", newIdx.RoleId)
			writes = append(writes, mongo.NewInsertOneModel().SetDocument(newIdx))
		}
	}

	for _, oldIdx := range oldMap {
		// Delete removed index
		r.logger.Infof("SyncIndices: Deleting index for RoleID %s", oldIdx.RoleId)
		writes = append(writes, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": oldIdx.ID}))
	}

	if len(writes) > 0 {
		_, err := r.collection.BulkWrite(ctx, writes)
		return err
	}
	return nil
}
func (r *CPSActionApproveIndexRepository) InsertMany(
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
					"role_id":     id,
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
						"role_id":     id,
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
					"role_id":     id,
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
		return err
	}

	r.logger.Infof("InsertMany: Successfully processed indices for RoleCode: %s", roleCode)
	return nil
}

// func (r *CPSActionApproveIndexRepository) InsertMany(ctx context.Context, makerIndex []bson.ObjectID, checkerIndex [][]bson.ObjectID, auditorIndex []bson.ObjectID, roleCode string) error {
// 	r.logger.Infof("InsertMany: Inserting indices for RoleCode: %s", roleCode)

// 	var models []mongo.WriteModel
// 	now := time.Now()

// 	// 1. Insert makers
// 	for _, id := range makerIndex {
// 		models = append(models, mongo.NewInsertOneModel().SetDocument(bson.M{
// 			"role_id":       id,
// 			"action_name":   roleCode,
// 			"maker_index":   1, // marker for maker
// 			"checker_index": nil,
// 			"auditor_index": nil,
// 			"created_at":    now,
// 			"updated_at":    now,
// 		}))
// 	}

// 	// 2. Insert checkers (by index)
// 	for i, ids := range checkerIndex {
// 		for _, id := range ids {
// 			models = append(models, mongo.NewInsertOneModel().SetDocument(bson.M{
// 				"role_id":       id,
// 				"action_name":   roleCode,
// 				"maker_index":   nil,
// 				"checker_index": i + 1, // 1-based
// 				"auditor_index": nil,
// 				"created_at":    now,
// 				"updated_at":    now,
// 			}))
// 		}
// 	}

// 	// 3. Insert auditors
// 	for _, id := range auditorIndex {
// 		models = append(models, mongo.NewInsertOneModel().SetDocument(bson.M{
// 			"role_id":       id,
// 			"action_name":   roleCode,
// 			"maker_index":   nil,
// 			"checker_index": nil,
// 			"auditor_index": 1, // marker for auditor
// 			"created_at":    now,
// 			"updated_at":    now,
// 		}))
// 	}

// 	if len(models) == 0 {
// 		return nil
// 	}

// 	_, err := r.collection.BulkWrite(ctx, models)
// 	if err != nil {
// 		r.logger.Errorf("InsertMany: BulkWrite failed: %v", err)
// 		return err
// 	}
// 	r.logger.Infof("InsertMany: Successfully inserted indices for RoleCode: %s", roleCode)
// 	return nil
// }

// func (r *CPSActionApproveIndexRepository) DeleteMany(ctx context.Context, makerIndex []bson.ObjectID, checkerIndex [][]bson.ObjectID, auditorIndex []bson.ObjectID, roleCode string) error {
// 	r.logger.Infof("DeleteMany: Deleting indices for RoleCode: %s", roleCode)

// 	var models []mongo.WriteModel

// 	// 1. Delete makers
// 	if len(makerIndex) > 0 {
// 		models = append(models, mongo.NewDeleteManyModel().SetFilter(bson.M{
// 			"role_id":     bson.M{"$in": makerIndex},
// 			"action_name": roleCode,
// 		}))
// 	}

// 	// 2. Delete checkers (by index)
// 	for i, ids := range checkerIndex {
// 		if len(ids) == 0 {
// 			continue
// 		}
// 		models = append(models, mongo.NewDeleteManyModel().SetFilter(bson.M{
// 			"role_id":       bson.M{"$in": ids},
// 			"action_name":   roleCode,
// 			"checker_index": i + 1, // 1-based
// 		}))
// 	}

// 	// 3. Delete auditors
// 	if len(auditorIndex) > 0 {
// 		models = append(models, mongo.NewDeleteManyModel().SetFilter(bson.M{
// 			"role_id":     bson.M{"$in": auditorIndex},
// 			"action_name": roleCode,
// 		}))
// 	}

// 	if len(models) == 0 {
// 		return nil
// 	}

//		_, err := r.collection.BulkWrite(ctx, models)
//		if err != nil {
//			r.logger.Errorf("DeleteMany: BulkWrite failed: %v", err)
//			return err
//		}
//		r.logger.Infof("DeleteMany: Successfully deleted indices for RoleCode: %s", roleCode)
//		return nil
//	}
func (r *CPSActionApproveIndexRepository) DeleteMany(
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
		models = append(models,
			mongo.NewDeleteManyModel().SetFilter(bson.M{
				"role_id":     bson.M{"$in": makerIndex},
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

		models = append(models,
			mongo.NewDeleteManyModel().SetFilter(bson.M{
				"role_id":       bson.M{"$in": ids},
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
		models = append(models,
			mongo.NewDeleteManyModel().SetFilter(bson.M{
				"role_id":     bson.M{"$in": auditorIndex},
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
		return err
	}

	r.logger.Infof("DeleteMany: Successfully deleted indices for RoleCode: %s", roleCode)
	return nil
}
