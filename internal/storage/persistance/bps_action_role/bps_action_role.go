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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BPSActionRoleRepository struct {
	client        *mongo.Client
	mongoDal      dal.MongoDal[imodel.BPSActionRole, imodel.BPSActionRole]
	approverDal   dal.MongoDal[imodel.BPSActionApproveIndex, imodel.BPSActionApproveIndex]
	actionListDal dal.MongoDal[imodel.BPSActionList, imodel.BPSActionList]
	logger        utils.Logger
	collection    *mongo.Collection
}

func NewBPSActionRoleRepository(client *mongo.Client, cfg *config.VaultConfig, database string, collection []string, logger utils.Logger) storage.BPSActionRoleRepository {
	return &BPSActionRoleRepository{
		client:        client,
		mongoDal:      dal.NewMongoDal[imodel.BPSActionRole, imodel.BPSActionRole](client, cfg, database, collection[0]),
		actionListDal: dal.NewMongoDal[imodel.BPSActionList, imodel.BPSActionList](client, cfg, database, collection[1]),
		approverDal:   dal.NewMongoDal[imodel.BPSActionApproveIndex, imodel.BPSActionApproveIndex](client, cfg, database, collection[2]),
		logger:        logger,
		collection:    client.Database(database).Collection(collection[0]),
	}
}

func (a *BPSActionRoleRepository) UpdateActionList(ctx context.Context, actionCode string, status bool) error {
	_, err := a.actionListDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, bson.M{"is_configured": status})
	if err != nil {
		a.logger.Errorf("[BPSActionRoleRepository][UpdateActionList] failed to update action list: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}
func (a *BPSActionRoleRepository) FindAllAccessListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.BPSActionList], error) {
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
		a.logger.Errorf("[BPSActionRoleRepository][FindAllAccessListWithPagination] failed to fetch access lists: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := a.actionListDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("[BPSActionRoleRepository][FindAllAccessListWithPagination] failed to count access lists: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	a.logger.Infof("[BPSActionRoleRepository][FindAllAccessListWithPagination] retrieved %d access lists", len(data))

	return &types.PaginatedResponse[[]imodel.BPSActionList]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *BPSActionRoleRepository) Create(ctx context.Context, actionRole *imodel.BPSActionRole) error {
	_, err := r.mongoDal.InsertOne(ctx, *actionRole)
	if err != nil {
		if mongo.IsTimeout(err) {
			r.logger.Errorf("[BPSActionRoleRepository][Create] timeout creating action role: %v", err)
			return errors.New(localization.ErrorInternalServerTimeout.Code)
		}
		r.logger.Errorf("[BPSActionRoleRepository][Create] failed to create action role: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *BPSActionRoleRepository) UpdateByActionCode(ctx context.Context, actionCode string, actionRole *imodel.BPSActionRole) error {
	update := bson.M{
		"action_name":             actionRole.ActionName,
		"assigned_makers_roles":   actionRole.AssignedMakersRoles,
		"assigned_checkers_roles": actionRole.AssignedCheckerRoles,
		"assigned_auditor_roles":  actionRole.AssignedAuditorRoles,
		"approver_count":          actionRole.ApproverCount,
		"is_maker_only":           actionRole.IsMakerOnly,
		"enabled":                 actionRole.Enabled,
		"updated_at":              actionRole.UpdatedAt,
	}
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, update)
	if err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][UpdateByActionCode] failed to update action role: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *BPSActionRoleRepository) EnableOrDisableByActionCode(ctx context.Context, actionCode string, enable bool) error {
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, bson.M{"enabled": enable})
	if err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][EnableOrDisableByActionCode] failed to enable/disable action code %s: %v", actionCode, err)
		return local_util.HandleDBError(err)
	}
	return nil
}
func (r *BPSActionRoleRepository) FindByActionCode(
	ctx context.Context,
	actionCode string,
) (*actionrole_dto.GetActionRoleByActionCodeRes, error) {

	const rolesCollection = "job_roles"

	pipeline := mongo.Pipeline{

		{{Key: "$match", Value: bson.M{
			"action_code": actionCode,
		}}},

		{{Key: "$addFields", Value: bson.M{
			"assigned_viewers_roles": bson.M{
				"$cond": bson.A{
					bson.M{"$isArray": "$assigned_viewers_roles"},
					"$assigned_viewers_roles",
					bson.A{},
				},
			},
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

		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         rolesCollection,
				"localField":   "assigned_viewers_roles",
				"foreignField": "code",
				"as":           "assigned_viewers_roles",
			},
		}},
		{{Key: "$addFields", Value: bson.M{
			"assigned_viewers_roles": bson.M{
				"$map": bson.M{
					"input": "$assigned_viewers_roles",
					"as":    "r",
					"in": bson.M{
						"_id":          "$$r._id",
						"code":         "$$r.code",
						"name":         "$$r.name",
						"portal_cards": "$$r.portal_cards",
						"created_at":   "$$r.created_at",
						"updated_at":   "$$r.updated_at",
					},
				},
			},
		}}},
		// 3️⃣ Lookup MAKERS (by code) + FORCE projection

		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         rolesCollection,
				"localField":   "assigned_makers_roles",
				"foreignField": "code",
				"as":           "assigned_makers_roles",
			},
		}},
		{{Key: "$addFields", Value: bson.M{
			"assigned_makers_roles": bson.M{
				"$map": bson.M{
					"input": "$assigned_makers_roles",
					"as":    "r",
					"in": bson.M{
						"_id":          "$$r._id",
						"code":         "$$r.code",
						"name":         "$$r.name",
						"portal_cards": "$$r.portal_cards",
						"created_at":   "$$r.created_at",
						"updated_at":   "$$r.updated_at",
					},
				},
			},
		}}},

		// 4️⃣ Lookup ALL CHECKER ROLES (flatten by code)
		{{
			Key: "$lookup",
			Value: bson.M{
				"from": rolesCollection,
				"let": bson.M{
					"allCheckerCodes": bson.M{
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
						"$expr": bson.M{
							"$in": []interface{}{"$code", "$$allCheckerCodes"},
						},
					}}},
				},
				"as": "checker_roles_all",
			},
		}},

		// 5️⃣ Rebuild CHECKERS [][]JobRole (FORCED projection)
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
										"as":    "code",
										"in": bson.M{
											"$let": bson.M{
												"vars": bson.M{
													"role": bson.M{
														"$first": bson.M{
															"$filter": bson.M{
																"input": "$checker_roles_all",
																"as":    "r",
																"cond": bson.M{
																	"$eq": []interface{}{"$$r.code", "$$code"},
																},
															},
														},
													},
												},
												"in": bson.M{
													"_id":          "$$role._id",
													"code":         "$$role.code",
													"name":         "$$role.name",
													"portal_cards": "$$role.portal_cards",
													"created_at":   "$$role.created_at",
													"updated_at":   "$$role.updated_at",
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

		// 6️⃣ Lookup AUDITORS (by code) + FORCE projection
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         rolesCollection,
				"localField":   "assigned_auditor_roles",
				"foreignField": "code",
				"as":           "assigned_auditor_roles",
			},
		}},
		{{Key: "$addFields", Value: bson.M{
			"assigned_auditor_roles": bson.M{
				"$map": bson.M{
					"input": "$assigned_auditor_roles",
					"as":    "r",
					"in": bson.M{
						"_id":          "$$r._id",
						"code":         "$$r.code",
						"name":         "$$r.name",
						"portal_cards": "$$r.portal_cards",
						"created_at":   "$$r.created_at",
						"updated_at":   "$$r.updated_at",
					},
				},
			},
		}}},

		// 7️⃣ Cleanup helper field
		{{Key: "$project", Value: bson.M{
			"checker_roles_all": 0,
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][FindByActionCode] failed to aggregate: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var results []*actionrole_dto.GetActionRoleByActionCodeRes
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][FindByActionCode] failed to decode results: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(results) == 0 {
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	return results[0], nil
}

func (r *BPSActionRoleRepository) FindAllWithPagination(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]*imodel.BPSActionRoleResposne], error) {

	// ----------------------------------
	// Build search filter
	// ----------------------------------
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

	const rolesCollection = "job_roles"

	// ----------------------------------
	// Aggregation pipeline
	// ----------------------------------
	pipeline := mongo.Pipeline{

		// 1️⃣ Match base filter
		{{Key: "$match", Value: filter}},

		// 2️⃣ Normalize arrays (schema safety)
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

		// 3️⃣ Lookup makers by CODE
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         rolesCollection,
				"localField":   "assigned_makers_roles",
				"foreignField": "code",
				"as":           "assigned_makers_roles",
			},
		}},

		// 4️⃣ Lookup ALL checker roles (flatten → code-based)
		{{
			Key: "$lookup",
			Value: bson.M{
				"from": rolesCollection,
				"let": bson.M{
					"allCheckerCodes": bson.M{
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
						"$expr": bson.M{
							"$in": []interface{}{"$code", "$$allCheckerCodes"},
						},
					}}},
				},
				"as": "checker_roles_all",
			},
		}},

		// 5️⃣ Rebuild nested [][]JobRole (checkers)
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
										"as":    "code",
										"in": bson.M{
											"$first": bson.M{
												"$filter": bson.M{
													"input": "$checker_roles_all",
													"as":    "role",
													"cond": bson.M{
														"$eq": []interface{}{"$$role.code", "$$code"},
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

		// 6️⃣ Lookup auditors by CODE
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         rolesCollection,
				"localField":   "assigned_auditor_roles",
				"foreignField": "code",
				"as":           "assigned_auditor_roles",
			},
		}},

		// 7️⃣ Cleanup helper field
		{{Key: "$project", Value: bson.M{
			"checker_roles_all": 0,
		}}},

		// 8️⃣ Pagination
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	// ----------------------------------
	// Execute aggregation
	// ----------------------------------
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][FindAllWithPagination] failed to aggregate action roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var data []*imodel.BPSActionRoleResposne
	if err := cursor.All(ctx, &data); err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][FindAllWithPagination] failed to decode action roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// ----------------------------------
	// Total count
	// ----------------------------------
	total, err := r.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][FindAllWithPagination] failed to count action roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(
		total,
		filterParam.Page,
		filterParam.PerPage,
	)

	return &types.PaginatedResponse[[]*imodel.BPSActionRoleResposne]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *BPSActionRoleRepository) FindByActionName(ctx context.Context, actionName string) (*imodel.BPSActionRole, error) {

	roleData, err := r.mongoDal.FindOne(ctx, bson.M{"action_name": actionName}, bson.M{})
	if err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][FindByActionName] error finding role by action name: %v error: %v", actionName, err)
		return nil, local_util.HandleDBError(err)
	}

	return roleData, nil
}

func (r *BPSActionRoleRepository) FindApproverByActionName(ctx context.Context, actionName, role_code string) (imodel.BPSActionApproveIndex, error) {
	approverModal, err := r.approverDal.FindOne(ctx, bson.M{"action_name": actionName, "role_id": role_code}, bson.M{})
	if err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][FindApproverByActionName] error finding approver index: %v", err)
		return imodel.BPSActionApproveIndex{}, local_util.HandleDBError(err)
	}

	return *approverModal, nil
}
func (r *BPSActionRoleRepository) FindByActionCodeOne(ctx context.Context, actionCode string) (*imodel.BPSActionRole, error) {
	result, err := r.mongoDal.FindOne(ctx, bson.M{"action_code": actionCode}, bson.M{})
	if err != nil {
		r.logger.Errorf("[BPSActionRoleRepository][FindByActionCodeOne] failed to find action role by code %s: %v", actionCode, err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}
