package action_role_repo

import (
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BPSActionRoleRepository struct {
	client     *mongo.Client
	mongoDal   dal.MongoDal[model.ActionRole, model.ActionRole]
	logger     utils.Logger
	collection *mongo.Collection
}

func NewBPSActionRoleRepository(client *mongo.Client, database, collection string, logger utils.Logger) storage.BPSActionRoleRepository {
	return &BPSActionRoleRepository{
		client:     client,
		mongoDal:   dal.NewMongoDal[model.ActionRole, model.ActionRole](client, database, collection),
		logger:     logger,
		collection: client.Database(database).Collection(collection),
	}
}

func (r *BPSActionRoleRepository) Create(ctx context.Context, actionRole *model.ActionRole) error {
	r.logger.Infof("[Create] creating BPS action role")
	_, err := r.mongoDal.InsertOne(ctx, *actionRole)
	if err != nil {
		r.logger.Errorf("[Create] failed to create BPS action role: %v", err)
		return err
	}
	r.logger.Infof("[Create] BPS action role created successfully")
	return nil
}

func (r *BPSActionRoleRepository) UpdateByActionCode(ctx context.Context, actionCode string, actionRole *model.ActionRole) error {
	r.logger.Infof("[UpdateByActionCode] updating BPS action role for action code: %s", actionCode)
	update := bson.M{
		"action_name":             actionRole.ActionName,
		"assigned_makers_roles":   actionRole.AssignedMakersRoles,
		"assigned_checkers_roles": actionRole.AssignedCheckerRoles,
		"assigned_auditors_roles": actionRole.AssignedAuditorRoles,
		"enabled":                 actionRole.Enabled,
		"is_maker_only":           actionRole.IsMakerOnly,
		"approver_count":          actionRole.ApproverCount,
		"updated_at":              actionRole.UpdatedAt,
	}
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, update)
	if err != nil {
		r.logger.Errorf("[UpdateByActionCode] failed to update BPS action role: %v", err)
		return err
	}
	r.logger.Infof("[UpdateByActionCode] BPS action role updated successfully")
	return nil
}

func (r *BPSActionRoleRepository) EnableOrDisableByActionCode(ctx context.Context, actionCode string, enable bool) error {
	r.logger.Infof("[EnableOrDisableByActionCode] processing BPS action role enable/disable for action code: %s, enabled: %v", actionCode, enable)
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, bson.M{"enabled": enable})
	if err != nil {
		r.logger.Errorf("[EnableOrDisableByActionCode] failed to enable/disable BPS action role: %v", err)
		return err
	}
	r.logger.Infof("[EnableOrDisableByActionCode] BPS action role enable/disable completed successfully")
	return nil
}

func (r *BPSActionRoleRepository) FindByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error) {
	const rolesCollection = "roles"

	pipeline := mongo.Pipeline{
		// Stage 1: Match the document by actionCode
		{{Key: "$match", Value: bson.D{{Key: "action_code", Value: actionCode}}}},

		// --- Stages for AssignedMakers (Simple Lookup) ---
		// Stage 2: Lookup the Roles for 'assigned_makers' (1:N relationship)
		{{Key: "$lookup", Value: bson.M{
			"from":         rolesCollection,
			"localField":   "assigned_makers_roles",
			"foreignField": "_id",
			"as":           "assigned_makers_roles", // Overwrites the IDs with the full Role objects
		}}},

		// --- Stages for AssignedCheckersRoles ([[]ObjectID] -> [[]Role]) ---

		// Stage 3: Unwind the OUTER array of assigned_checkers, preserving the index
		{{Key: "$unwind", Value: bson.M{
			"path":              "$assigned_checkers_roles",
			"includeArrayIndex": "outer_index",
		}}},

		// Stage 4: Unwind the INNER array of ObjectID, preserving the index
		{{Key: "$unwind", Value: bson.M{
			"path":              "$assigned_checkers_roles",
			"includeArrayIndex": "inner_index",
		}}},

		// Stage 5: Lookup the Role document for the now single ObjectID
		{{Key: "$lookup", Value: bson.M{
			"from":         rolesCollection,
			"localField":   "assigned_checkers_roles", // This is now a single ObjectID
			"foreignField": "_id",
			"as":           "checker_role_doc", // Temporary field for the fetched role (still an array [Role])
		}}},

		// Stage 6: Convert the temporary single-element array to a single object
		// This prepares the role for grouping.
		{{Key: "$unwind", Value: "$checker_role_doc"}},

		// Stage 7: Group back the inner array ([Role])
		// Group by the original document ID AND the outer_index to rebuild the inner array.
		// IMPORTANT: Use $first to carry forward ALL original fields.
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "outer_index", Value: "$outer_index"},
			}},
			// --- Carry forward all original fields ---
			{Key: "action_code", Value: bson.D{{Key: "$first", Value: "$action_code"}}},
			{Key: "action_name", Value: bson.D{{Key: "$first", Value: "$action_name"}}},
			{Key: "enabled", Value: bson.D{{Key: "$first", Value: "$enabled"}}},
			{Key: "updated_at", Value: bson.D{{Key: "$first", Value: "$updated_at"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "assigned_makers_roles", Value: bson.D{{Key: "$first", Value: "$assigned_makers_roles"}}},
			// ----------------------------------------

			// Reconstruct the inner array of roles
			{Key: "inner_roles_array", Value: bson.D{{Key: "$push", Value: "$checker_role_doc"}}},
		}}},

		// Stage 8: Group back the outer array ([[]Role])
		// Group by the original document ID to rebuild the outer array and final document.
		// IMPORTANT: Use $first to carry forward ALL original fields.
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id._id"}, // The original document ID

			// --- Carry forward all original fields (using the field names from the previous stage) ---
			{Key: "action_code", Value: bson.D{{Key: "$first", Value: "$action_code"}}},
			{Key: "action_name", Value: bson.D{{Key: "$first", Value: "$action_name"}}},
			{Key: "enabled", Value: bson.D{{Key: "$first", Value: "$enabled"}}},
			{Key: "updated_at", Value: bson.D{{Key: "$first", Value: "$updated_at"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "assigned_makers_roles", Value: bson.D{{Key: "$first", Value: "$assigned_makers_roles"}}},
			// ----------------------------------------

			// Reconstruct the outer array of role arrays
			{Key: "assigned_checkers_roles", Value: bson.D{{Key: "$push", Value: "$inner_roles_array"}}},
		}}},

		// Stage 9: Final Projection (Optional but recommended for clean output)
		// This renames the fields to match the DTO struct's JSON/BSON tags if necessary
		// and ensures the output document is clean.
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1}, // keep _id
			{Key: "action_code", Value: 1},
			{Key: "action_name", Value: 1},
			{Key: "enabled", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "assigned_makers_roles", Value: 1},
			{Key: "assigned_checkers_roles", Value: 1},
			{Key: "assigned_auditors_roles", Value: 1},
		}}},
	}

	r.logger.Infof("[FindByActionCode] fetching BPS action role by action code: %s", actionCode)
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[FindByActionCode] failed to aggregate BPS action role: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []*actionrole_dto.GetActionRoleByActionCodeRes
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("[FindByActionCode] failed to decode aggregation results: %v", err)
		return nil, err
	}
	if len(results) == 0 {
		r.logger.Errorf("[FindByActionCode] BPS action role not found")
		return nil, fmt.Errorf(localization.ErrorBpsActionRoleNotFound.Code)
	}
	r.logger.Infof("[FindByActionCode] BPS action role retrieved successfully")
	return results[0], nil
}

func (r *BPSActionRoleRepository) FindByActionName(ctx context.Context, actionName string) (*model.ActionRole, error) {
	return r.mongoDal.FindOne(ctx, bson.M{"action_name": actionName}, bson.M{})
}

func (r *BPSActionRoleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ActionRole], error) {
	r.logger.Infof("[FindAllWithPagination] fetching BPS action roles with pagination")
	searchKeys := bson.M{}
	allowedKeys := []string{"action_code", "action_name", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"action_code": searchRegex},
			{"action_name": searchRegex},
		}
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := r.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		r.logger.Errorf("[FindAllWithPagination] failed to fetch BPS action roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[FindAllWithPagination] failed to count BPS action roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	r.logger.Infof("[FindAllWithPagination] retrieved %d BPS action roles", len(data))

	return &types.PaginatedResponse[[]*model.ActionRole]{
		Data: data,
		Meta: meta,
	}, nil
}
