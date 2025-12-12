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

type CPSActionRoleRepository struct {
	client     *mongo.Client
	mongoDal   dal.MongoDal[model.CPSActionRole, model.CPSActionRole]
	logger     utils.Logger
	collection *mongo.Collection
}

func NewCPSActionRoleRepository(client *mongo.Client, database, collection string, logger utils.Logger) storage.CPSActionRoleRepository {
	return &CPSActionRoleRepository{
		client:     client,
		mongoDal:   dal.NewMongoDal[model.CPSActionRole, model.CPSActionRole](client, database, collection),
		logger:     logger,
		collection: client.Database(database).Collection(collection),
	}
}

func (r *CPSActionRoleRepository) Create(ctx context.Context, actionRole *model.CPSActionRole) error {
	_, err := r.mongoDal.InsertOne(ctx, *actionRole)
	return err
}

func (r *CPSActionRoleRepository) UpdateByActionCode(ctx context.Context, actionCode string, actionRole *model.CPSActionRole) error {
	update := bson.M{
		"action_name":             actionRole.ActionName,
		"assigned_makers_roles":   actionRole.AssignedMakersRoles,
		"assigned_checkers_roles": actionRole.AssignedCheckersRoles,
		"enabled":                 actionRole.Enabled,
		"updated_at":              actionRole.UpdatedAt,
	}
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, update)
	return err
}

func (r *CPSActionRoleRepository) EnableOrDisableByActionCode(ctx context.Context, actionCode string, enable bool) error {
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, bson.M{"enabled": enable})
	return err
}

func (r *CPSActionRoleRepository) FindByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error) {
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
    {Key: "_id", Value: 1},
    {Key: "action_code", Value: 1},
    {Key: "action_name", Value: 1},
    {Key: "enabled", Value: 1},
    {Key: "updated_at", Value: 1},
    {Key: "created_at", Value: 1},
    {Key: "assigned_makers_roles", Value: 1},
    {Key: "assigned_checkers_roles", Value: 1},
    {Key: "assigned_auditor_roles", Value: 1}, // FIXED spelling
}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		
		return nil, err
	}
	var results []*actionrole_dto.GetActionRoleByActionCodeRes
	if err := cursor.All(ctx, &results); err != nil {
		fmt.Println("/////// pipline error ",err)
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New(localization.ErrorBpsActionRoleNotFound.Code)
	}
	return results[0], nil
}
func (r *CPSActionRoleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.CPSActionRole], error) {
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
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.CPSActionRole]{
		Data: data,
		Meta: meta,
	}, nil
}
func (r *CPSActionRoleRepository) FindByActionName(ctx context.Context, actionName string) (*model.CPSActionRole, error) {
	return r.mongoDal.FindOne(ctx, bson.M{"action_name": actionName}, bson.M{})
}
func (r *CPSActionRoleRepository) FindByActionCodeOne(ctx context.Context, actionCode string) (*model.CPSActionRole, error) {
	return r.mongoDal.FindOne(ctx, bson.M{"action_code": actionCode}, bson.M{})
}
