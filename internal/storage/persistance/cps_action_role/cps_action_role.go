package action_role_repo

import (
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/cps_action_role"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSActionRoleRepository struct {
	client        *mongo.Client
	mongoDal      dal.MongoDal[imodel.CPSActionRole, imodel.CPSActionRole]
	approverDal   dal.MongoDal[imodel.CPSActionApproveIndex, imodel.CPSActionApproveIndex]
	actionListDal dal.MongoDal[imodel.CPSActionList, imodel.CPSActionList]
	logger        utils.Logger
	collection    *mongo.Collection
}

func NewCPSActionRoleRepository(client *mongo.Client, cfg *config.VaultConfig, database string, collection []string, logger utils.Logger) storage.CPSActionRoleRepository {
	return &CPSActionRoleRepository{
		client:        client,
		mongoDal:      dal.NewMongoDal[imodel.CPSActionRole, imodel.CPSActionRole](client, cfg, database, "cps_action_roles"),
		actionListDal: dal.NewMongoDal[imodel.CPSActionList, imodel.CPSActionList](client, cfg, database, collection[1]),
		approverDal:   dal.NewMongoDal[imodel.CPSActionApproveIndex, imodel.CPSActionApproveIndex](client, cfg, database, "cps_action_approver_index"),
		logger:        logger,
		collection:    client.Database(database).Collection(collection[0]),
	}
}

func (a *CPSActionRoleRepository) UpdateActionList(ctx context.Context, actionCode, portalCard string, status bool) error {
	_, err := a.actionListDal.UpdateOne(ctx, bson.M{"action_code": actionCode, "portal_card_name": portalCard}, bson.M{"is_configured": status})
	return err
}
func (a *CPSActionRoleRepository) FindAllAccessListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CPSActionList], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"action_name", "action_code"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}

		searchKeys["$or"] = []bson.M{
			{"action_name": searchRegex},
			{"action_code": searchRegex},
			{"portal_card_name": searchRegex},
		}
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

	return &types.PaginatedResponse[[]imodel.CPSActionList]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *CPSActionRoleRepository) Create(ctx context.Context, actionRole *imodel.CPSActionRole) error {
	_, err := r.mongoDal.InsertOne(ctx, *actionRole)
	return err
}

func (r *CPSActionRoleRepository) UpdateByActionCode(ctx context.Context, actionCode string, actionRole *imodel.CPSActionRole) error {
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
	return err
}

func (r *CPSActionRoleRepository) EnableOrDisableByActionCode(ctx context.Context, actionCode string, enable bool) error {
	_, err := r.mongoDal.UpdateOne(ctx, bson.M{"action_code": actionCode}, bson.M{"enabled": enable})
	return err
}
func (r *CPSActionRoleRepository) FindByActionCode(
	ctx context.Context,
	actionCode string,
) (*actionrole_dto.GetActionRoleByActionCodeRes, error) {

	const rolesCollection = "job_roles"

	pipeline := mongo.Pipeline{

		// 1️⃣ Match action
		{{
			Key: "$match",
			Value: bson.M{
				"action_code": actionCode,
			},
		}},

		// 2️⃣ Normalize arrays (SAFETY)
		{{
			Key: "$addFields",
			Value: bson.M{
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
			},
		}},

		// 3️⃣ VIEWERS (simple []string)
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         rolesCollection,
				"localField":   "assigned_viewers_roles",
				"foreignField": "code",
				"as":           "assigned_viewers_roles",
			},
		}},
		{{
			Key: "$addFields",
			Value: bson.M{
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
			},
		}},

		// 4️⃣ MAKERS (simple []string)
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         rolesCollection,
				"localField":   "assigned_makers_roles",
				"foreignField": "code",
				"as":           "assigned_makers_roles",
			},
		}},
		{{
			Key: "$addFields",
			Value: bson.M{
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
			},
		}},

		// 5️⃣ CHECKERS — flatten [][]string
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

		// 6️⃣ REBUILD CHECKERS [][]JobRole
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

		// 7️⃣ AUDITORS — flatten [][]string
		{{
			Key: "$lookup",
			Value: bson.M{
				"from": rolesCollection,
				"let": bson.M{
					"allAuditorCodes": bson.M{
						"$reduce": bson.M{
							"input":        "$assigned_auditor_roles",
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
							"$in": []interface{}{"$code", "$$allAuditorCodes"},
						},
					}}},
				},
				"as": "auditor_roles_all",
			},
		}},

		// 8️⃣ REBUILD AUDITORS [][]JobRole
		{{
			Key: "$addFields",
			Value: bson.M{
				"assigned_auditor_roles": bson.M{
					"$map": bson.M{
						"input": "$assigned_auditor_roles",
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
																"input": "$auditor_roles_all",
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

		// 9️⃣ CLEANUP
		{{
			Key: "$project",
			Value: bson.M{
				"checker_roles_all": 0,
				"auditor_roles_all": 0,
			},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []*actionrole_dto.GetActionRoleByActionCodeRes
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("action role not found")
	}

	return results[0], nil
}

func (r *CPSActionRoleRepository) FindAllWithPagination(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]*model.CPSActionRoleResposne], error) {

	// ----------------------------------
	// Build search filter
	// ----------------------------------
	searchKeys := bson.M{}
	allowedKeys := []string{"action_code", "action_name", "portal_card_name", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"action_code": searchRegex, "$options": "i"},
			{"action_name": searchRegex, "$options": "i"},
			{"portal_card_name": searchRegex, "$options": "i"},
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
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var data []*model.CPSActionRoleResposne
	if err := cursor.All(ctx, &data); err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// ----------------------------------
	// Total count
	// ----------------------------------
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

func (r *CPSActionRoleRepository) FindByActionName(ctx context.Context, actionName string) (*imodel.CPSActionRole, error) {

	roleData, err := r.mongoDal.FindOne(ctx, bson.M{"action_name": actionName}, bson.M{})
	if err != nil {
		r.logger.Errorf("error finding role by action name: %v error: %v", actionName, err)
		return nil, err
	}

	return roleData, nil
}

func (r *CPSActionRoleRepository) FindByActionNameAndPortalCard(ctx context.Context, actionName, portalCard string) (*imodel.CPSActionRole, error) {

	roleData, err := r.mongoDal.FindOne(ctx, bson.M{"action_name": actionName, "portal_card_name": portalCard}, bson.M{})
	if err != nil {
		r.logger.Errorf("error finding role by action name: %v error: %v", actionName, err)
		return nil, err
	}

	return roleData, nil
}

func (r *CPSActionRoleRepository) FindApproverByActionName(ctx context.Context, actionName, role_code string) (imodel.CPSActionApproveIndex, error) {
	approverModal, err := r.approverDal.FindOne(ctx, bson.M{"action_name": actionName, "role_id": role_code}, bson.M{})
	if err != nil {
		r.logger.Errorf("error finding approver index: %v", err)
		return imodel.CPSActionApproveIndex{}, err
	}

	return *approverModal, nil
}
func (r *CPSActionRoleRepository) FindByActionCodeOne(ctx context.Context, actionCode string) (*imodel.CPSActionRole, error) {
	return r.mongoDal.FindOne(ctx, bson.M{"action_code": actionCode}, bson.M{})
}
