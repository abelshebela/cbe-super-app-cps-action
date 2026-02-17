package action_role_repo

import (
	"cbe-super-app-cps-action/internal/constants/model"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

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
		collection: client.Database(database).Collection("cps_action_approver_index"),
	}
}

func (r *CPSActionApproveIndexRepository) PopulateUserApproverAllocations(ctx context.Context, role_id string) ([]string, []string, []string, []string, error) {
	r.logger.Infof("PopulateUserApproverAllocations: Populating approver allocations for RoleID: %s", role_id)
	var makerAllocations []string
	var checkerAllocations []string
	var auditorAllocations []string
	var portalCard []string

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"role_id": role_id,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "cps_action_roles",
			"localField":   "action_name",
			"foreignField": "action_name",
			"as":           "action_role_info",
		}}},
		{{Key: "$match", Value: bson.M{
			"action_role_info": bson.M{
				"$elemMatch": bson.M{"enabled": true},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("PopulateUserApproverAllocations: Aggregate failed: %v", err)
		return nil, nil, nil, nil, err
	}
	defer cursor.Close(ctx)
	var results []imodel.CPSActionApproveIndex
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("PopulateUserApproverAllocations: Cursor.All failed: %v", err)
		return nil, nil, nil, nil, err
	}

	for _, v := range results {
		if v.MakerIndex != nil {
			makerAllocations = append(makerAllocations, v.ActionName)
		}
		if v.CheckerIndex != nil {
			checkerAllocations = append(checkerAllocations, v.ActionName)
		}
		if v.AuditorIndex != nil {
			auditorAllocations = append(auditorAllocations, v.ActionName)
		}
		portalCard = append(portalCard, v.PortalCardName)
	}
	r.logger.Infof("PopulateUserApproverAllocations: Found %d maker, %d checker, %d auditor allocations for RoleID: %s", len(makerAllocations), len(checkerAllocations), len(auditorAllocations), role_id)
	return makerAllocations, checkerAllocations, auditorAllocations, portalCard, nil
}
func (r *CPSActionApproveIndexRepository) SaveIndices(ctx context.Context, indices []imodel.CPSActionApproveIndex) error {
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

func (r *CPSActionApproveIndexRepository) SyncIndices(ctx context.Context, oldActionName, portalCard string, newIndices []imodel.CPSActionApproveIndex, isVersionChanged bool) error {
	r.logger.Infof("SyncIndices: Syncing %d indices for oldActionName: %s", len(newIndices), oldActionName)

	// if !isVersionChanged && len(newIndices) > 0 {
	// 	if _, err := r.collection.DeleteMany(ctx, bson.M{"action_name": oldActionName, "version": newIndices[0].Version, "portal_card_name": portalCard}); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	if _, err := r.collection.DeleteMany(ctx, bson.M{"action_name": oldActionName, "portal_card_name": portalCard}); err != nil {
	// 		return err
	// 	}
	// }

	if _, err := r.collection.DeleteMany(ctx, bson.M{"action_name": oldActionName, "portal_card_name": portalCard}); err != nil {
		return err
	}
	if err := r.SaveIndices(ctx, newIndices); err != nil {
		return err
	}

	return nil
}

func (r *CPSActionApproveIndexRepository) FindMakerAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.CPSActionApproveIndex, error) {
	r.logger.Infof("FindMakerAllocationsByRoleID: Finding maker allocations for RoleID: %s", roleID.Hex())

	cursor, err := r.collection.Find(ctx, bson.M{
		"role_id":     roleID.Hex(),
		"maker_index": bson.M{"$ne": nil},
	})
	if err != nil {
		r.logger.Errorf("FindMakerAllocationsByRoleID: Find failed: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []imodel.CPSActionApproveIndex
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("FindMakerAllocationsByRoleID: Cursor.All failed: %v", err)
		return nil, err
	}

	r.logger.Infof("FindMakerAllocationsByRoleID: Found %d maker allocations for RoleID: %s", len(results), roleID.Hex())
	return results, nil
}

func (r *CPSActionApproveIndexRepository) FindCheckerAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.CPSActionApproveIndex, error) {
	r.logger.Infof("FindCheckerAllocationsByRoleID: Finding checker allocations for RoleID: %s", roleID.Hex())

	cursor, err := r.collection.Find(ctx, bson.M{
		"role_id":       roleID.Hex(),
		"checker_index": bson.M{"$ne": nil},
	})
	if err != nil {
		r.logger.Errorf("FindCheckerAllocationsByRoleID: Find failed: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []imodel.CPSActionApproveIndex
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("FindCheckerAllocationsByRoleID: Cursor.All failed: %v", err)
		return nil, err
	}

	r.logger.Infof("FindCheckerAllocationsByRoleID: Found %d checker allocations for RoleID: %s", len(results), roleID.Hex())
	return results, nil
}

func (r *CPSActionApproveIndexRepository) FindAuditorAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.CPSActionApproveIndex, error) {
	r.logger.Infof("FindAuditorAllocationsByRoleID: Finding auditor allocations for RoleID: %s", roleID.Hex())

	cursor, err := r.collection.Find(ctx, bson.M{
		"role_id":       roleID.Hex(),
		"auditor_index": bson.M{"$ne": nil},
	})
	if err != nil {
		r.logger.Errorf("FindAuditorAllocationsByRoleID: Find failed: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []imodel.CPSActionApproveIndex
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("FindAuditorAllocationsByRoleID: Cursor.All failed: %v", err)
		return nil, err
	}

	r.logger.Infof("FindAuditorAllocationsByRoleID: Found %d auditor allocations for RoleID: %s", len(results), roleID.Hex())
	return results, nil
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

func (r *CPSActionApproveIndexRepository) InsertAll(ctx context.Context, new model.CPSActionRole) error {
	r.logger.Infof("InsertAll: Inserting indices for ActionName: %s", new.ActionName)
	var indices []imodel.CPSActionApproveIndex
	now := time.Now()

	// Makers: 1-based index
	for i, makerId := range new.AssignedMakersRoles {
		idx := int64(i + 1)
		indices = append(indices, imodel.CPSActionApproveIndex{
			ID:           bson.NewObjectID(),
			RoleId:       makerId,
			ActionName:   new.ActionName,
			MakerIndex:   &idx,
			CheckerIndex: nil,
			AuditorIndex: nil,
			UpdatedAt:    now,
			CreatedAt:    now,
		})
	}

	// Checkers: 2D slice, float index (e.g., 1.1, 2.1, ...)
	for i, checkerGroup := range new.AssignedCheckerRoles {
		for j, checkerId := range checkerGroup {
			idx := float64(i+1) + float64(j+1)*0.1
			indices = append(indices, imodel.CPSActionApproveIndex{
				ID:           bson.NewObjectID(),
				RoleId:       checkerId,
				ActionName:   new.ActionName,
				MakerIndex:   nil,
				CheckerIndex: &idx,
				AuditorIndex: nil,
				UpdatedAt:    now,
				CreatedAt:    now,
			})
		}
	}

	for i, auditorGroup := range new.AssignedAuditorRoles {
		for j, auditorId := range auditorGroup {
			idx := float64(i+1) + float64(j+1)*0.1
			indices = append(indices, imodel.CPSActionApproveIndex{
				ID:           bson.NewObjectID(),
				RoleId:       auditorId,
				ActionName:   new.ActionName,
				MakerIndex:   nil,
				CheckerIndex: nil,
				AuditorIndex: &idx,
				UpdatedAt:    now,
				CreatedAt:    now,
			})
		}
	}

	if len(indices) == 0 {
		r.logger.Infof("InsertAll: No indices to insert for ActionName: %s", new.ActionName)
		return nil
	}

	docs := make([]interface{}, len(indices))
	for i, v := range indices {
		docs[i] = v
	}
	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		r.logger.Errorf("InsertAll: InsertMany failed: %v", err)
		return err
	}
	r.logger.Infof("InsertAll: Successfully inserted %d indices for ActionName: %s", len(indices), new.ActionName)
	return nil
}

func (r *CPSActionApproveIndexRepository) DeleteAll(ctx context.Context, prev imodel.CPSActionRoleResposne) error {
	r.logger.Infof("DeleteAll: Deleting indices for ActionName: %s", prev.ActionName)
	var models []mongo.WriteModel

	// Makers
	if prev.AssignedMakersRoles != nil && len(prev.AssignedMakersRoles) > 0 {
		var makerIds []bson.ObjectID
		for _, v := range prev.AssignedMakersRoles {
			if v == nil {
				continue
			}
			switch val := v.(type) {
			case string:
				if objId, err := bson.ObjectIDFromHex(val); err == nil {
					makerIds = append(makerIds, objId)
				}
			case bson.ObjectID:
				makerIds = append(makerIds, val)
			case map[string]interface{}:
				if idVal, ok := val["_id"]; ok {
					switch id := idVal.(type) {
					case bson.ObjectID:
						makerIds = append(makerIds, id)
					case map[string]interface{}:
						if oid, ok := id["$oid"].(string); ok {
							if objId, err := bson.ObjectIDFromHex(oid); err == nil {
								makerIds = append(makerIds, objId)
							}
						}
					}
				}
			}
		}
		if len(makerIds) > 0 {
			models = append(models, mongo.NewDeleteManyModel().SetFilter(bson.M{
				"_id":         bson.M{"$in": makerIds},
				"action_name": prev.ActionName,
			}))
		}
	}

	// Checkers
	if prev.AssignedCheckersRoles != nil && len(prev.AssignedCheckersRoles) > 0 {
		for _, group := range prev.AssignedCheckersRoles {
			var checkerIds []bson.ObjectID
			for _, v := range group {
				if v == nil {
					continue
				}
				switch val := v.(type) {
				case string:
					if objId, err := bson.ObjectIDFromHex(val); err == nil {
						checkerIds = append(checkerIds, objId)
					}
				case bson.ObjectID:
					checkerIds = append(checkerIds, val)
				case map[string]interface{}:
					if idVal, ok := val["_id"]; ok {
						switch id := idVal.(type) {
						case bson.ObjectID:
							checkerIds = append(checkerIds, id)
						case map[string]interface{}:
							if oid, ok := id["$oid"].(string); ok {
								if objId, err := bson.ObjectIDFromHex(oid); err == nil {
									checkerIds = append(checkerIds, objId)
								}
							}
						}
					}
				}
			}
			if len(checkerIds) > 0 {
				models = append(models, mongo.NewDeleteManyModel().SetFilter(bson.M{
					"_id":         bson.M{"$in": checkerIds},
					"action_name": prev.ActionName,
				}))
			}
		}
	}

	// Auditors
	if prev.AssignedAuditorRoles != nil && len(prev.AssignedAuditorRoles) > 0 {
		var auditorIds []bson.ObjectID
		for _, v := range prev.AssignedAuditorRoles {
			if v == nil {
				continue
			}
			switch val := v.(type) {
			case string:
				if objId, err := bson.ObjectIDFromHex(val); err == nil {
					auditorIds = append(auditorIds, objId)
				}
			case bson.ObjectID:
				auditorIds = append(auditorIds, val)
			case map[string]interface{}:
				if idVal, ok := val["_id"]; ok {
					switch id := idVal.(type) {
					case bson.ObjectID:
						auditorIds = append(auditorIds, id)
					case map[string]interface{}:
						if oid, ok := id["$oid"].(string); ok {
							if objId, err := bson.ObjectIDFromHex(oid); err == nil {
								auditorIds = append(auditorIds, objId)
							}
						}
					}
				}
			}
		}
		if len(auditorIds) > 0 {
			models = append(models, mongo.NewDeleteManyModel().SetFilter(bson.M{
				"_id":         bson.M{"$in": auditorIds},
				"action_name": prev.ActionName,
			}))
		}
	}

	if len(models) == 0 {
		r.logger.Infof("DeleteAll: No delete models generated for ActionName: %s", prev.ActionName)
		return nil
	}

	_, err := r.collection.BulkWrite(ctx, models)
	if err != nil {
		r.logger.Errorf("DeleteAll: BulkWrite failed: %v", err)
		return err
	}
	r.logger.Infof("DeleteAll: Successfully deleted indices for ActionName: %s", prev.ActionName)
	return nil
}
