package action_role_repo

import (
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSActionRoleRepository struct {
	client        *mongo.Client
	mongoDal      dal.MongoDal[model.CPSActionRole, model.CPSActionRole]
	actionListDal dal.MongoDal[imodel.CPSActionList, imodel.CPSActionList]
	logger        utils.Logger
	collection    *mongo.Collection
}

func NewCPSActionRoleRepository(client *mongo.Client, database string, collection []string, logger utils.Logger) storage.CPSActionRoleRepository {
	return &CPSActionRoleRepository{
		client:        client,
		mongoDal:      dal.NewMongoDal[model.CPSActionRole, model.CPSActionRole](client, database, collection[0]),
		actionListDal: dal.NewMongoDal[imodel.CPSActionList, imodel.CPSActionList](client, database, collection[1]),
		logger:        logger,
		collection:    client.Database(database).Collection(collection[0]),
	}
}

func (a *CPSActionRoleRepository) UpdateActionList(ctx context.Context, actionCode string, status bool) error {
	_, err := a.actionListDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, bson.M{"is_configured": status})
	return err
}
func (a *CPSActionRoleRepository) FindAllAccessListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.CPSActionList], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"action_name", "action_code"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["action_name"] = searchRegex
		searchKeys["action_code"] = searchRegex
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := a.actionListDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		a.logger.Errorf("[FindAllWithPagination] failed to fetch access lists: %v", err)
		return nil, err
	}

	total, err := a.actionListDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("[FindAllWithPagination] failed to count access lists: %v", err)
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	a.logger.Infof("[FindAllWithPagination] retrieved %d access lists", len(data))

	return &types.PaginatedResponse[[]*imodel.CPSActionList]{
		Data: data,
		Meta: meta,
	}, nil
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
		"assigned_auditor_roles":  actionRole.AssignedAuditorRoles,
		"approver_count":          actionRole.ApproverCount,
		"is_maker_only":           actionRole.IsMakerOnly,
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
	const rolesCollection = "job_roles"
	pipeline := mongo.Pipeline{
		// Stage 1: Match by actionCode
		{{Key: "$match", Value: bson.D{{Key: "action_code", Value: actionCode}}}},

		// Stage 2: Lookup assigned_makers_roles (simple array)
		{{Key: "$lookup", Value: bson.M{
			"from":         rolesCollection,
			"localField":   "assigned_makers_roles",
			"foreignField": "code",
			"as":           "assigned_makers_roles",
		}}},

		// Stage 3: Lookup assigned_checkers_roles (nested array [[]ObjectID])
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$assigned_checkers_roles",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "checker_outer_index",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$assigned_checkers_roles",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "checker_inner_index",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         rolesCollection,
			"localField":   "assigned_checkers_roles",
			"foreignField": "code",
			"as":           "checker_role_doc",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$checker_role_doc",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "checker_outer_index", Value: "$checker_outer_index"},
			}},
			{Key: "action_code", Value: bson.D{{Key: "$first", Value: "$action_code"}}},
			{Key: "action_name", Value: bson.D{{Key: "$first", Value: "$action_name"}}},
			{Key: "is_maker_only", Value: bson.D{{Key: "$first", Value: "$is_maker_only"}}},
			{Key: "approver_count", Value: bson.D{{Key: "$first", Value: "$approver_count"}}},
			{Key: "enabled", Value: bson.D{{Key: "$first", Value: "$enabled"}}},
			{Key: "updated_at", Value: bson.D{{Key: "$first", Value: "$updated_at"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "assigned_makers_roles", Value: bson.D{{Key: "$first", Value: "$assigned_makers_roles"}}},
			// FIX: Pass the auditor IDs forward
			{Key: "assigned_auditor_roles", Value: bson.D{{Key: "$first", Value: "$assigned_auditor_roles"}}},
			{Key: "inner_checkers", Value: bson.D{{Key: "$push", Value: "$checker_role_doc"}}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id._id"},
			{Key: "action_code", Value: bson.D{{Key: "$first", Value: "$action_code"}}},
			{Key: "action_name", Value: bson.D{{Key: "$first", Value: "$action_name"}}},
			{Key: "approver_count", Value: bson.D{{Key: "$first", Value: "$approver_count"}}},
			{Key: "is_maker_only", Value: bson.D{{Key: "$first", Value: "$is_maker_only"}}},
			{Key: "enabled", Value: bson.D{{Key: "$first", Value: "$enabled"}}},
			{Key: "updated_at", Value: bson.D{{Key: "$first", Value: "$updated_at"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "assigned_makers_roles", Value: bson.D{{Key: "$first", Value: "$assigned_makers_roles"}}},
			// FIX: Pass the auditor IDs forward again
			{Key: "assigned_auditor_roles", Value: bson.D{{Key: "$first", Value: "$assigned_auditor_roles"}}},
			{Key: "assigned_checkers_roles", Value: bson.D{{Key: "$push", Value: "$inner_checkers"}}},
		}}},

		// Stage 4: Lookup assigned_auditor_roles (similar to makers)
		{{Key: "$lookup", Value: bson.M{
			"from":         rolesCollection,
			"localField":   "assigned_auditor_roles",
			"foreignField": "code",
			"as":           "assigned_auditor_roles",
		}}},

		// Stage 5: Project final fields
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "action_code", Value: 1},
			{Key: "action_name", Value: 1},
			{Key: "is_maker_only", Value: 1},
			{Key: "approver_count", Value: 1},
			{Key: "enabled", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "assigned_makers_roles", Value: 1},
			{Key: "assigned_checkers_roles", Value: 1},
			{Key: "assigned_auditor_roles", Value: 1},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {

		return nil, err
	}
	var results []*actionrole_dto.GetActionRoleByActionCodeRes
	if err := cursor.All(ctx, &results); err != nil {
		fmt.Println("/////// pipline error ", err)
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New(localization.ErrorBpsActionRoleNotFound.Code)
	}
	return results[0], nil
}

func (r *CPSActionRoleRepository) FindAllWithPagination(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]*model.CPSActionRoleResposne], error) {

	// -----------------------------
	// Build search filter
	// -----------------------------
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

	// -----------------------------
	// Aggregation pipeline
	// -----------------------------
	pipeline := mongo.Pipeline{

		// 1️⃣ Match base filter
		{{Key: "$match", Value: filter}},

		// 2️⃣ Normalize arrays to avoid schema-drift crashes
		{{Key: "$addFields", Value: bson.M{
			"assigned_makers_roles": bson.M{
				"$cond": bson.A{
					bson.M{"$isArray": "$assigned_makers_roles"},
					"$assigned_makers_roles",
					bson.A{},
				},
			},
			"assigned_checkers_roles": bson.M{
				"$cond": bson.A{
					bson.M{"$isArray": "$assigned_checkers_roles"},
					"$assigned_checkers_roles",
					bson.A{},
				},
			},
			"assigned_auditor_roles": bson.M{
				"$cond": bson.A{
					bson.M{"$isArray": "$assigned_auditor_roles"},
					"$assigned_auditor_roles",
					bson.A{},
				},
			},
		}}},

		// 3️⃣ Lookup maker roles
		{{
			Key: "$lookup",
			Value: bson.M{
				"from": "roles",
				"let":  bson.M{"roleIds": "$assigned_makers_roles"},
				"pipeline": mongo.Pipeline{
					{{Key: "$match", Value: bson.M{
						"$expr": bson.M{"$in": []interface{}{"$_id", "$$roleIds"}},
					}}},
				},
				"as": "assigned_makers_roles",
			},
		}},

		// 4️⃣ Lookup ALL checker roles (flatten safely)
		{{
			Key: "$lookup",
			Value: bson.M{
				"from": "roles",
				"let": bson.M{
					"allCheckerIds": bson.M{
						"$reduce": bson.M{
							"input":        "$assigned_checkers_roles",
							"initialValue": bson.A{},
							"in": bson.M{
								"$concatArrays": bson.A{
									"$$value",
									bson.M{
										"$cond": bson.A{
											bson.M{"$isArray": "$$this"},
											"$$this",
											bson.A{},
										},
									},
								},
							},
						},
					},
				},
				"pipeline": mongo.Pipeline{
					{{Key: "$match", Value: bson.M{
						"$expr": bson.M{"$in": []interface{}{"$_id", "$$allCheckerIds"}},
					}}},
				},
				"as": "checker_roles_all",
			},
		}},

		// 5️⃣ Rebuild nested checker arrays → full role docs
		{{
			Key: "$addFields",
			Value: bson.M{
				"assigned_checkers_roles": bson.M{
					"$map": bson.M{
						"input": "$assigned_checkers_roles",
						"as":    "level",
						"in": bson.M{
							"$cond": bson.A{
								bson.M{"$isArray": "$$level"},
								bson.M{
									"$map": bson.M{
										"input": "$$level",
										"as":    "rid",
										"in": bson.M{
											"$first": bson.M{
												"$filter": bson.M{
													"input": "$checker_roles_all",
													"as":    "role",
													"cond": bson.M{
														"$eq": []interface{}{"$$role._id", "$$rid"},
													},
												},
											},
										},
									},
								},
								bson.A{},
							},
						},
					},
				},
			},
		}},

		// 6️⃣ Lookup auditor roles
		{{
			Key: "$lookup",
			Value: bson.M{
				"from": "roles",
				"let":  bson.M{"roleIds": "$assigned_auditor_roles"},
				"pipeline": mongo.Pipeline{
					{{Key: "$match", Value: bson.M{
						"$expr": bson.M{"$in": []interface{}{"$_id", "$$roleIds"}},
					}}},
				},
				"as": "assigned_auditor_roles",
			},
		}},

		// 7️⃣ Cleanup helper field
		{{Key: "$project", Value: bson.M{"checker_roles_all": 0}}},

		// 8️⃣ Pagination
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	// -----------------------------
	// Execute aggregation
	// -----------------------------
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var data []*model.CPSActionRoleResposne
	if err := cursor.All(ctx, &data); err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// -----------------------------
	// Total count
	// -----------------------------
	total, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(
		total,
		filterParam.Page,
		filterParam.PerPage,
	)

	return &types.PaginatedResponse[[]*model.CPSActionRoleResposne]{
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
