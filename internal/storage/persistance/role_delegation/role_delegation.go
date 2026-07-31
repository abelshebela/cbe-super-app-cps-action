package role_delegation_repo

import (
	cpsuser "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type roleDelegationRepository struct {
	repo              dal.MongoDal[imodel.RoleDelegation, imodel.RoleDelegation]
	cpsRepo           dal.MongoDal[imodel.CPSUser, imodel.CPSUser]
	bpsRepo           dal.MongoDal[imodel.BPSUser, imodel.BPSUser]
	logger            utils.Logger
	collection        *mongo.Collection
	cpsUserCollection *mongo.Collection
	bpsUserCollection *mongo.Collection
}

// Create implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) CreateWithExistingUser(ctx context.Context, role *imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegationRepository][CreateWithExistingUser] creating role delegation for user=%s job_title=%s rest of data: %+v", role.DelegatedUserID, role.DelegatedUserJobTitle, role)
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Create] failed to start session: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (any, error) {
		role.ID = bson.NewObjectID()
		role.Enable = true
		result, err := r.collection.InsertOne(sc, role)
		if err != nil {
			log.Errorf("[RoleDelegationRepository][Create] failed to create role delegation: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		insertedID, ok := result.InsertedID.(bson.ObjectID)
		if !ok {
			log.Errorf("[RoleDelegationRepository][Create] inserted id has unexpected type: %T", result.InsertedID)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		updatePipeline := mongo.Pipeline{
			bson.D{{
				Key: "$set",
				Value: bson.M{
					"is_delegation_active": true,
					"delegation_id":        insertedID,
					"enabled":              true,
				},
			}},
		}

		if role.DelegatedUserUserType == "CPS" {
			_, err = r.cpsUserCollection.UpdateOne(sc, bson.M{"username": role.DelegatedUserID, "is_deleted": false}, updatePipeline)
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to update cps user delegate: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		} else {
			_, err = r.bpsUserCollection.UpdateOne(sc, bson.M{"username": role.DelegatedUserID, "is_deleted": false}, updatePipeline)
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to update cps user delegate: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		}

		if role.RevokeExistingDelegation {
			_, err = r.collection.UpdateMany(sc, bson.M{"delegated_user_id": role.DelegatedUserID, "_id": bson.M{"$ne": insertedID}}, bson.M{"$set": bson.M{"enable": false}})
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to disable existing delegations: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		}

		return nil, nil
	})
	if err != nil {
		r.logger.Errorf("[RoleDelegationRepository][Create] transaction failed: %v", err)
		return localization.ErrorUnexpectedError
	}

	log.Infof("[RoleDelegationRepository][Create] role delegation created successfully")
	return nil
}

// Create implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) CreateWithNewUser(ctx context.Context, role *imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegationRepository][CreateNewUser] creating role delegation for user=%s job_title=%s restOfData:%+v", role.DelegatedUserID, role.DelegatedUserJobTitle, role)
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Create] failed to start session: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (any, error) {
		role.ID = bson.NewObjectID()
		role.Enable = true
		// if role.DelegatedUserUserType == "CPS" {
		// 	role.DelegatedUserUserCode = local_util.GenerateCPSUserCode()
		// } else {
		// 	role.DelegatedUserUserCode = local_util.GenerateBPSUserCode()
		// }
		result, err := r.collection.InsertOne(sc, role)
		if err != nil {
			log.Errorf("[RoleDelegationRepository][Create] failed to create role delegation: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		insertedID, ok := result.InsertedID.(bson.ObjectID)
		if !ok {
			log.Errorf("[RoleDelegationRepository][Create] inserted id has unexpected type: %T", result.InsertedID)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		now := time.Now()
		if role.DelegatedUserUserType == "CPS" {
			assignedDepartment, err := bson.ObjectIDFromHex(role.DelegatedUserDepartmentOrBranch)
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] invalid department id: %v", err)
				return nil, localization.ErrorUnexpectedError
			}
			_, err = r.cpsUserCollection.InsertOne(sc, model.CPSUser{
				UserCode:           role.DelegatedUserUserCode,
				FullName:           role.DelegatedUserFullName,
				Role:               role.DelegatedUserExistingRole,
				PhoneNumber:        role.DelegatedUserPhoneNumber,
				Email:              role.DelegatedUserEmail,
				UserName:           role.DelegatedUserID,
				JobTitle:           role.DelegatedUserJobTitle,
				Department:         assignedDepartment,
				Enabled:            true,
				DateJoined:         &now,
				CreatedAt:          now,
				IsDelegationActive: true,
				DelegationID:       insertedID,
				CreatedBy:          "CPS Portal",
			})
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to update cps user delegate: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		} else {
			_, err = r.bpsUserCollection.InsertOne(sc, model.BPSUser{
				UserCode:         role.DelegatedUserUserCode,
				FullName:         role.DelegatedUserFullName,
				Role:             role.DelegatedUserExistingRole,
				PhoneNumber:      role.DelegatedUserPhoneNumber,
				Email:            role.DelegatedUserEmail,
				UserName:         role.DelegatedUserID,
				JobTitle:         role.DelegatedUserJobTitle,
				BranchCode:       []string{role.DelegatedUserDepartmentOrBranch},
				Enabled:          true,
				CreatedAt:        now,
				IsFirstTimeLogin: true,
				ChangedBy:        "CPS Portal",
			})
			if err != nil {
				log.Errorf("[RoleDelegationRepository][Create] failed to update bps user delegate: %v", err)
				return nil, local_util.HandleDBError(err)
			}
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	log.Infof("[RoleDelegationRepository][Create] role delegation created successfully")
	return nil
}

// EnableOrDisable implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegationRepository][EnableOrDisable] updating role delegation id=%s enabled=%v", id, enable)
	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enable": enable, "updated_at": time.Now()}

	_, err = r.repo.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][EnableOrDisable] failed to update role delegation: %v", err)
		return local_util.HandleDBError(err)
	}

	log.Infof("[RoleDelegationRepository][EnableOrDisable] role delegation updated successfully")
	return nil
}

// FindForExport implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) FindForExport(ctx context.Context, startDate, endDate time.Time) ([]imodel.RoleDelegation, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	filter := bson.M{
		"start_at": bson.M{"$gte": startDate},
		"end_at":   bson.M{"$lte": endDate},
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"new_department_obj_id": bson.M{
				"$convert": bson.M{
					"input":   "$new_department_or_branch",
					"to":      "objectId",
					"onError": nil,
					"onNull":  nil,
				},
			},
			"delegated_department_obj_id": bson.M{
				"$convert": bson.M{
					"input":   "$delegated_user_department_or_branch",
					"to":      "objectId",
					"onError": nil,
					"onNull":  nil,
				},
			},
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "department",
			"localField":   "new_department_obj_id",
			"foreignField": "_id",
			"as":           "new_department_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "department",
			"localField":   "delegated_department_obj_id",
			"foreignField": "_id",
			"as":           "delegated_department_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "job_roles",
			"let":  bson.M{"role_value": "$new_role_id"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$ne": bson.A{"$is_deleted", true}},
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{"$code", "$$role_value"}},
							bson.M{"$eq": bson.A{"$role", "$$role_value"}},
							bson.M{"$eq": bson.A{"$name", "$$role_value"}},
							bson.M{"$eq": bson.A{"$job_title", "$$role_value"}},
						}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
			},
			"as": "new_role_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "job_roles",
			"let":  bson.M{"role_value": "$delegated_user_existing_role"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$ne": bson.A{"$is_deleted", true}},
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{"$code", "$$role_value"}},
							bson.M{"$eq": bson.A{"$role", "$$role_value"}},
							bson.M{"$eq": bson.A{"$name", "$$role_value"}},
							bson.M{"$eq": bson.A{"$job_title", "$$role_value"}},
						}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
			},
			"as": "delegated_role_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "job_roles",
			"let":  bson.M{"role_value": "$delegator_user_role"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$ne": bson.A{"$is_deleted", true}},
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{"$code", "$$role_value"}},
							bson.M{"$eq": bson.A{"$role", "$$role_value"}},
							bson.M{"$eq": bson.A{"$name", "$$role_value"}},
							bson.M{"$eq": bson.A{"$job_title", "$$role_value"}},
						}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
			},
			"as": "delegator_role_info",
		}}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"new_department_or_branch_name": bson.M{
				"$cond": bson.A{
					bson.M{"$eq": bson.A{"$delegation_type", "BPS"}},
					"$new_department_or_branch",
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$new_department_info.department", 0}}, ""}},
				},
			},
			"delegated_user_department_or_branch_name": bson.M{
				"$cond": bson.A{
					bson.M{"$eq": bson.A{"$delegated_user_user_type", "BPS"}},
					"$delegated_user_department_or_branch",
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$delegated_department_info.department", 0}}, ""}},
				},
			},
			"new_role_id_name": bson.M{"$ifNull": bson.A{
				bson.M{"$arrayElemAt": bson.A{"$new_role_info.name", 0}},
				bson.M{"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$new_role_info.job_title", 0}},
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$new_role_info.role_name", 0}}, ""}},
				}},
			}},
			"delegated_user_existing_role_name": bson.M{"$ifNull": bson.A{
				bson.M{"$arrayElemAt": bson.A{"$delegated_role_info.name", 0}},
				bson.M{"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$delegated_role_info.job_title", 0}},
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$delegated_role_info.role_name", 0}}, ""}},
				}},
			}},
			"delegator_user_role_name": bson.M{"$ifNull": bson.A{
				bson.M{"$arrayElemAt": bson.A{"$delegator_role_info.name", 0}},
				bson.M{"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$delegator_role_info.job_title", 0}},
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$delegator_role_info.role_name", 0}}, ""}},
				}},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"new_department_obj_id":       0,
			"delegated_department_obj_id": 0,
			"new_department_info":         0,
			"delegated_department_info":   0,
			"new_role_info":               0,
			"delegated_role_info":         0,
			"delegator_role_info":         0,
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindForExport] failed to fetch role delegations: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var data []imodel.RoleDelegation
	if err = cursor.All(ctx, &data); err != nil {
		log.Errorf("[RoleDelegationRepository][FindForExport] failed to decode role delegations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return data, nil
}

// Find implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) Find(ctx context.Context, filter bson.M) (*imodel.RoleDelegation, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	result, err := r.repo.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Find] failed to find role delegation: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

// FindAll implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) FindAll(ctx context.Context) (*[]imodel.RoleDelegation, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	data, err := r.repo.FindAllWithCursorBasedPagination(ctx, dal.FilterOp{})
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindAll] failed to fetch role delegations: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return &data, nil
}

// FindAllWithPagination implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"delegated_user_job_title": searchRegex},
			{"delegated_user_id": searchRegex},
			{"delegated_user_email": searchRegex},
			{"delegated_user_phone_number": searchRegex},
			{"delegated_user_full_name": searchRegex},
		}
	}

	allowedKeys := []string{"enable", "delegated_user_email", "delegated_user_phone_number", "delegated_user_job_title", "delegated_user_id", "delegated_user_full_name", "start_at", "end_at"}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	if val, ok := filterParam.Filters["portal"]; ok {
		if val == "BPS" || val == "CPS" {
			filter["delegation_type"] = val
		}
	}

	if val, ok := filterParam.Filters["status"]; ok {
		if val == "Pending" {
			filter["enable"] = true
			filter["start_at"] = bson.M{"$gt": time.Now()}
		}
		if val == "Expired" {
			filter["enable"] = true
			filter["end_at"] = bson.M{"$lte": time.Now()}
		}
		if val == "Revoked" {
			filter["enable"] = false
		}
		if val == "Active" {
			filter["enable"] = true
			filter["start_at"] = bson.M{"$lte": time.Now()}
			filter["end_at"] = bson.M{"$gt": time.Now()}
		}
		if val == "Suspended" {
			filter["enable"] = true
			filter["start_at"] = bson.M{"$gt": time.Now()}
		}
	}

	// start_date/end_date are validated in handler; apply them as additional bounds.
	rangeAndClauses := bson.A{}
	if startDateRaw, ok := filterParam.Filters["start_date"]; ok {
		if startDate, ok := startDateRaw.(time.Time); ok {
			rangeAndClauses = append(rangeAndClauses, bson.M{"start_at": bson.M{"$gte": startDate}})
		}
	}

	if endDateRaw, ok := filterParam.Filters["end_date"]; ok {
		if endDate, ok := endDateRaw.(time.Time); ok {
			rangeAndClauses = append(rangeAndClauses, bson.M{"end_at": bson.M{"$lte": endDate}})
		}
	}

	if len(rangeAndClauses) > 0 {
		if existingAnd, exists := filter["$and"]; exists {
			switch v := existingAnd.(type) {
			case bson.A:
				filter["$and"] = append(v, rangeAndClauses...)
			case []bson.M:
				andClauses := make(bson.A, 0, len(v)+len(rangeAndClauses))
				for _, clause := range v {
					andClauses = append(andClauses, clause)
				}
				andClauses = append(andClauses, rangeAndClauses...)
				filter["$and"] = andClauses
			default:
				filter["$and"] = rangeAndClauses
			}
		} else {
			filter["$and"] = rangeAndClauses
		}
	}

	total, err := r.repo.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to count role delegations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if total <= 0 || skip >= total {
		return &types.PaginatedResponse[[]imodel.RoleDelegation]{}, nil
	}

	pipeline := mongo.Pipeline{
		// 1. Initial Filtering, Sorting, and Pagination
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},

		// 2. CONDITIONAL CPS LOOKUP (Only runs if type is CPS)
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": r.cpsUserCollection.Name(),
			"let": bson.M{
				"user_type": "$delegated_user_user_type",
				"user_code": "$delegated_user_user_code",
			},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$$user_type", "CPS"}}, // IF condition
						bson.M{"$eq": bson.A{"$user_code", "$$user_code"}},
						bson.M{"$ne": bson.A{"$is_deleted", true}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{
					"user_code":                           1,
					"full_name":                           1,
					"username":                            1,
					"email":                               1,
					"phone_number":                        1,
					"job_title":                           1,
					"role":                                1,
					"delegated_user_department_or_branch": bson.M{"$toString": "$department"},
				}}},
			},
			"as": "cps_source",
		}}},

		// 3. CONDITIONAL BPS LOOKUP (Only runs if type is BPS)
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": r.bpsUserCollection.Name(),
			"let": bson.M{
				"user_type": "$delegated_user_user_type",
				"user_code": "$delegated_user_user_code",
			},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$$user_type", "BPS"}}, // IF condition
						bson.M{"$eq": bson.A{"$user_code", "$$user_code"}},
						bson.M{"$ne": bson.A{"$is_deleted", true}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{
					"user_code":    1,
					"full_name":    1,
					"username":     1,
					"email":        1,
					"phone_number": 1,
					"job_title":    1,
					"role":         1,
					"delegated_user_department_or_branch": bson.M{"$ifNull": bson.A{
						"$branch_name",
						bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$branch_code", 0}}, ""}},
					}},
				}}},
			},
			"as": "bps_source",
		}}},

		// 4. Merge the source arrays into a single unified object variable
		bson.D{{Key: "$addFields", Value: bson.M{
			"delegated_user_source": bson.M{
				"$cond": bson.A{
					bson.M{"$eq": bson.A{"$delegated_user_user_type", "CPS"}},
					bson.M{"$arrayElemAt": bson.A{"$cps_source", 0}},
					bson.M{"$cond": bson.A{
						bson.M{"$eq": bson.A{"$delegated_user_user_type", "BPS"}},
						bson.M{"$arrayElemAt": bson.A{"$bps_source", 0}},
						nil,
					}},
				},
			},
		}}},

		// Clean up temporary lookups right away
		bson.D{{Key: "$project", Value: bson.M{"cps_source": 0, "bps_source": 0}}},

		// 5. Map fields exactly like you did before (using the unified delegated_user_source)
		bson.D{{Key: "$addFields", Value: bson.M{
			"delegated_user_user_code":            bson.M{"$ifNull": bson.A{"$delegated_user_source.user_code", "$delegated_user_user_code"}},
			"delegated_user_full_name":            bson.M{"$ifNull": bson.A{"$delegated_user_source.full_name", "$delegated_user_full_name"}},
			"delegated_user_id":                   bson.M{"$ifNull": bson.A{"$delegated_user_source.username", "$delegated_user_id"}},
			"delegated_user_email":                bson.M{"$ifNull": bson.A{"$delegated_user_source.email", "$delegated_user_email"}},
			"delegated_user_phone_number":         bson.M{"$ifNull": bson.A{"$delegated_user_source.phone_number", "$delegated_user_phone_number"}},
			"delegated_user_job_title":            bson.M{"$ifNull": bson.A{"$delegated_user_source.job_title", "$delegated_user_job_title"}},
			"delegated_user_existing_role":        bson.M{"$ifNull": bson.A{"$delegated_user_source.role", "$delegated_user_existing_role"}},
			"delegated_user_department_or_branch": bson.M{"$ifNull": bson.A{"$delegated_user_source.delegated_user_department_or_branch", "$delegated_user_department_or_branch"}},
		}}},
		bson.D{{Key: "$project", Value: bson.M{"delegated_user_source": 0}}},

		// 6. Rest of your pipeline (Conversions, Department lookups, Job role lookups, etc.)
		bson.D{{Key: "$addFields", Value: bson.M{
			"new_department_obj_id":       bson.M{"$convert": bson.M{"input": "$new_department_or_branch", "to": "objectId", "onError": nil, "onNull": nil}},
			"delegated_department_obj_id": bson.M{"$convert": bson.M{"input": "$delegated_user_department_or_branch", "to": "objectId", "onError": nil, "onNull": nil}},
		}}},
		// ... Keep the exact same code you had from your $lookup department stages down to the end

		// pipeline := mongo.Pipeline{
		// 	bson.D{{Key: "$match", Value: filter}},
		// 	bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		// 	bson.D{{Key: "$skip", Value: skip}},
		// 	bson.D{{Key: "$limit", Value: limit}},
		// 	bson.D{{Key: "$facet", Value: bson.M{
		// 		"cps_data": mongo.Pipeline{
		// 			bson.D{{Key: "$match", Value: bson.M{"delegated_user_user_type": "CPS"}}},
		// 			bson.D{{Key: "$lookup", Value: bson.M{
		// 				"from": r.cpsUserCollection.Name(),
		// 				"let":  bson.M{"delegated_user_user_code": "$delegated_user_user_code"},
		// 				"pipeline": mongo.Pipeline{
		// 					bson.D{{Key: "$match", Value: bson.M{
		// 						"$expr": bson.M{"$and": bson.A{
		// 							bson.M{"$eq": bson.A{"$user_code", "$$delegated_user_user_code"}},
		// 							bson.M{"$ne": bson.A{"$is_deleted", true}},
		// 						}},
		// 					}}},
		// 					bson.D{{Key: "$limit", Value: 1}},
		// 					bson.D{{Key: "$project", Value: bson.M{
		// 						"user_code":                           1,
		// 						"full_name":                           1,
		// 						"username":                            1,
		// 						"email":                               1,
		// 						"phone_number":                        1,
		// 						"job_title":                           1,
		// 						"role":                                1,
		// 						"delegated_user_department_or_branch": bson.M{"$toString": "$department"},
		// 					}}},
		// 				},
		// 				"as": "delegated_user_source",
		// 			}}},
		// 			bson.D{{Key: "$addFields", Value: bson.M{
		// 				"delegated_user_source":               bson.M{"$arrayElemAt": bson.A{"$delegated_user_source", 0}},
		// 				"delegated_user_user_code":            bson.M{"$ifNull": bson.A{"$delegated_user_source.user_code", "$delegated_user_user_code"}},
		// 				"delegated_user_full_name":            bson.M{"$ifNull": bson.A{"$delegated_user_source.full_name", "$delegated_user_full_name"}},
		// 				"delegated_user_id":                   bson.M{"$ifNull": bson.A{"$delegated_user_source.username", "$delegated_user_id"}},
		// 				"delegated_user_email":                bson.M{"$ifNull": bson.A{"$delegated_user_source.email", "$delegated_user_email"}},
		// 				"delegated_user_phone_number":         bson.M{"$ifNull": bson.A{"$delegated_user_source.phone_number", "$delegated_user_phone_number"}},
		// 				"delegated_user_job_title":            bson.M{"$ifNull": bson.A{"$delegated_user_source.job_title", "$delegated_user_job_title"}},
		// 				"delegated_user_existing_role":        bson.M{"$ifNull": bson.A{"$delegated_user_source.role", "$delegated_user_existing_role"}},
		// 				"delegated_user_department_or_branch": bson.M{"$ifNull": bson.A{"$delegated_user_source.delegated_user_department_or_branch", "$delegated_user_department_or_branch"}},
		// 			}}},
		// 			bson.D{{Key: "$project", Value: bson.M{"delegated_user_source": 0}}},
		// 		},
		// 		"bps_data": mongo.Pipeline{
		// 			bson.D{{Key: "$match", Value: bson.M{"delegated_user_user_type": "BPS"}}},
		// 			bson.D{{Key: "$lookup", Value: bson.M{
		// 				"from": r.bpsUserCollection.Name(),
		// 				"let":  bson.M{"delegated_user_user_code": "$delegated_user_user_code"},
		// 				"pipeline": mongo.Pipeline{
		// 					bson.D{{Key: "$match", Value: bson.M{
		// 						"$expr": bson.M{"$and": bson.A{
		// 							bson.M{"$eq": bson.A{"$user_code", "$$delegated_user_user_code"}},
		// 							bson.M{"$ne": bson.A{"$is_deleted", true}},
		// 						}},
		// 					}}},
		// 					bson.D{{Key: "$limit", Value: 1}},
		// 					bson.D{{Key: "$project", Value: bson.M{
		// 						"user_code":    1,
		// 						"full_name":    1,
		// 						"username":     1,
		// 						"email":        1,
		// 						"phone_number": 1,
		// 						"job_title":    1,
		// 						"role":         1,
		// 						"delegated_user_department_or_branch": bson.M{"$ifNull": bson.A{
		// 							"$branch_name",
		// 							bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$branch_code", 0}}, ""}},
		// 						}},
		// 					}}},
		// 				},
		// 				"as": "delegated_user_source",
		// 			}}},
		// 			bson.D{{Key: "$addFields", Value: bson.M{
		// 				"delegated_user_source":               bson.M{"$arrayElemAt": bson.A{"$delegated_user_source", 0}},
		// 				"delegated_user_user_code":            bson.M{"$ifNull": bson.A{"$delegated_user_source.user_code", "$delegated_user_user_code"}},
		// 				"delegated_user_full_name":            bson.M{"$ifNull": bson.A{"$delegated_user_source.full_name", "$delegated_user_full_name"}},
		// 				"delegated_user_id":                   bson.M{"$ifNull": bson.A{"$delegated_user_source.username", "$delegated_user_id"}},
		// 				"delegated_user_email":                bson.M{"$ifNull": bson.A{"$delegated_user_source.email", "$delegated_user_email"}},
		// 				"delegated_user_phone_number":         bson.M{"$ifNull": bson.A{"$delegated_user_source.phone_number", "$delegated_user_phone_number"}},
		// 				"delegated_user_job_title":            bson.M{"$ifNull": bson.A{"$delegated_user_source.job_title", "$delegated_user_job_title"}},
		// 				"delegated_user_existing_role":        bson.M{"$ifNull": bson.A{"$delegated_user_source.role", "$delegated_user_existing_role"}},
		// 				"delegated_user_department_or_branch": bson.M{"$ifNull": bson.A{"$delegated_user_source.delegated_user_department_or_branch", "$delegated_user_department_or_branch"}},
		// 			}}},
		// 			bson.D{{Key: "$project", Value: bson.M{"delegated_user_source": 0}}},
		// 		},
		// 		"other_data": mongo.Pipeline{
		// 			bson.D{{Key: "$match", Value: bson.M{"delegated_user_user_type": bson.M{"$nin": bson.A{"CPS", "BPS"}}}}},
		// 		},
		// 	}}},
		// 	bson.D{{Key: "$project", Value: bson.M{
		// 		"merged_docs": bson.M{"$concatArrays": bson.A{"$cps_data", "$bps_data", "$other_data"}},
		// 	}}},
		// 	bson.D{{Key: "$unwind", Value: "$merged_docs"}},
		// 	bson.D{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$merged_docs"}}},
		// 	bson.D{{Key: "$addFields", Value: bson.M{
		// 		"new_department_obj_id": bson.M{
		// 			"$convert": bson.M{
		// 				"input":   "$new_department_or_branch",
		// 				"to":      "objectId",
		// 				"onError": nil,
		// 				"onNull":  nil,
		// 			},
		// 		},
		// 		"delegated_department_obj_id": bson.M{
		// 			"$convert": bson.M{
		// 				"input":   "$delegated_user_department_or_branch",
		// 				"to":      "objectId",
		// 				"onError": nil,
		// 				"onNull":  nil,
		// 			},
		// 		},
		// 	}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "department",
			"localField":   "new_department_obj_id",
			"foreignField": "_id",
			"as":           "new_department_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "department",
			"localField":   "delegated_department_obj_id",
			"foreignField": "_id",
			"as":           "delegated_department_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "job_roles",
			"let":  bson.M{"role_value": "$new_role_id"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$ne": bson.A{"$is_deleted", true}},
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{"$code", "$$role_value"}},
							bson.M{"$eq": bson.A{"$role", "$$role_value"}},
							bson.M{"$eq": bson.A{"$name", "$$role_value"}},
							bson.M{"$eq": bson.A{"$job_title", "$$role_value"}},
						}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
			},
			"as": "new_role_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "job_roles",
			"let":  bson.M{"role_value": "$delegated_user_existing_role"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$ne": bson.A{"$is_deleted", true}},
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{"$code", "$$role_value"}},
							bson.M{"$eq": bson.A{"$role", "$$role_value"}},
							bson.M{"$eq": bson.A{"$name", "$$role_value"}},
							bson.M{"$eq": bson.A{"$job_title", "$$role_value"}},
						}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
			},
			"as": "delegated_role_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "job_roles",
			"let":  bson.M{"role_value": "$delegator_user_role"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$ne": bson.A{"$is_deleted", true}},
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{"$code", "$$role_value"}},
							bson.M{"$eq": bson.A{"$role", "$$role_value"}},
							bson.M{"$eq": bson.A{"$name", "$$role_value"}},
							bson.M{"$eq": bson.A{"$job_title", "$$role_value"}},
						}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
			},
			"as": "delegator_role_info",
		}}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"new_department_or_branch_name": bson.M{
				"$cond": bson.A{
					bson.M{"$eq": bson.A{"$delegation_type", "BPS"}},
					"$new_department_or_branch",
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$new_department_info.department", 0}}, ""}},
				},
			},
			"delegated_user_department_or_branch_name": bson.M{
				"$cond": bson.A{
					bson.M{"$eq": bson.A{"$delegated_user_user_type", "BPS"}},
					"$delegated_user_department_or_branch",
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$delegated_department_info.department", 0}}, ""}},
				},
			},
			"new_role_id_name": bson.M{"$ifNull": bson.A{
				bson.M{"$arrayElemAt": bson.A{"$new_role_info.name", 0}},
				bson.M{"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$new_role_info.job_title", 0}},
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$new_role_info.role_name", 0}}, ""}},
				}},
			}},
			"delegated_user_existing_role_name": bson.M{"$ifNull": bson.A{
				bson.M{"$arrayElemAt": bson.A{"$delegated_role_info.name", 0}},
				bson.M{"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$delegated_role_info.job_title", 0}},
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$delegated_role_info.role_name", 0}}, ""}},
				}},
			}},
			"delegator_user_role_name": bson.M{"$ifNull": bson.A{
				bson.M{"$arrayElemAt": bson.A{"$delegator_role_info.name", 0}},
				bson.M{"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$delegator_role_info.job_title", 0}},
					bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$delegator_role_info.role_name", 0}}, ""}},
				}},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"new_department_obj_id":       0,
			"delegated_department_obj_id": 0,
			"new_department_info":         0,
			"delegated_department_info":   0,
			"new_role_info":               0,
			"delegated_role_info":         0,
			"delegator_role_info":         0,
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to fetch role delegations: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var data []imodel.RoleDelegation
	if err = cursor.All(ctx, &data); err != nil {
		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to decode role delegations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]imodel.RoleDelegation]{
		Data: data,
		Meta: meta,
	}, nil
}

// // FindAllWithPagination implements [storage.RoleDelegationRepository].
// func (r *roleDelegationRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
// 	log := local_util.LoggerFromCtx(ctx, r.logger)

// 	searchKeys := bson.M{}
// 	if filterParam.Search != "" {
// 		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
// 		searchKeys["$or"] = []bson.M{
// 			{"delegated_user_job_title": searchRegex},
// 			{"delegated_user_id": searchRegex},
// 			{"delegated_user_email": searchRegex},
// 			{"delegated_user_phone_number": searchRegex},
// 			{"delegated_user_full_name": searchRegex},
// 		}
// 	}

// 	allowedKeys := []string{"enable", "delegated_user_email", "delegated_user_phone_number", "delegated_user_job_title", "delegated_user_id", "delegated_user_full_name", "start_at", "end_at"}
// 	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

// 	if val, ok := filterParam.Filters["portal"]; ok {
// 		if val == "BPS" || val == "CPS" {
// 			filter["delegation_type"] = val
// 		}
// 	}

// 	if val, ok := filterParam.Filters["status"]; ok {
// 		if val == "Expired" {
// 			filter["enable"] = true
// 			filter["end_at"] = bson.M{"$lte": time.Now()}
// 		}
// 		if val == "Revoked" {
// 			filter["enable"] = false
// 		}
// 		if val == "Active" {
// 			filter["enable"] = true
// 			filter["start_at"] = bson.M{"$lte": time.Now()}
// 			filter["end_at"] = bson.M{"$gt": time.Now()}
// 		}
// 		if val == "Suspended" {
// 			filter["enable"] = true
// 			filter["start_at"] = bson.M{"$gt": time.Now()}
// 			filter["end_at"] = bson.M{"$gt": time.Now()}
// 		}
// 	}

// 	total, err := r.repo.TotalCount(ctx, filter)
// 	if err != nil {
// 		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to count role delegations: %v", err)
// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
// 	}

// 	if total <= 0 || skip >= total {
// 		return &types.PaginatedResponse[[]imodel.RoleDelegation]{}, nil
// 	}

// 	data, err := r.repo.FindAllWithPaginationD(ctx, dal.FilterParam{
// 		Filter:     filter,
// 		Projection: bson.M{},
// 		Sort:       bson.D{bson.E{Key: "created_at", Value: -1}},
// 		Skip:       skip,
// 		Limit:      limit,
// 	})
// 	if err != nil {
// 		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to fetch role delegations: %v", err)
// 		return nil, local_util.HandleDBError(err)
// 	}

// 	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
// 	return &types.PaginatedResponse[[]imodel.RoleDelegation]{
// 		Data: data,
// 		Meta: meta,
// 	}, nil
// }

func (r *roleDelegationRepository) FindByUsername(ctx context.Context, id string, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	allowedKeys := []string{}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowedKeys)
	filter["delegated_user_id"] = id

	total, err := r.repo.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindAllWithPagination] failed to count role delegations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if total <= 0 || skip >= total {
		return &types.PaginatedResponse[[]imodel.RoleDelegation]{}, nil
	}

	data, err := r.repo.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindAllWithPagination] failed to fetch role delegations: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &types.PaginatedResponse[[]imodel.RoleDelegation]{}, localization.ErrRoleDelegationForUserNotFound
		}
		return nil, localization.ErrorUnexpectedError
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]imodel.RoleDelegation]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *roleDelegationRepository) FindByUserCode(ctx context.Context, usercode string) (*cpsuser.CpsUserPopulatedResponse, error) {
	r.logger.Infof("[RoleDelegationRepository][FindByUsername] fetching active delegation for user code: %s", usercode)
	log := local_util.LoggerFromCtx(ctx, r.logger)
	currentTime := time.Now()

	data, err := r.repo.FindOne(ctx, bson.M{
		"delegated_user_user_code": usercode,
		"start_at":                 bson.M{"$lte": currentTime},
		"end_at":                   bson.M{"$gt": currentTime},
		"enable":                   true,
	}, bson.M{})
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByUsername] failed to fetch role delegations: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, localization.ErrorUserNotFound
		}
		return nil, localization.ErrorUnexpectedError
	}

	return &cpsuser.CpsUserPopulatedResponse{
		ID:       data.ID,
		UserCode: data.DelegatedUserUserCode,
		FullName: data.DelegatedUserFullName,
		Role: cpsuser.RoleResponse{
			Code: data.NewRoleID,
		},
		JobTitle: data.DelegatedUserJobTitle,
		Department: &cpsuser.DepartmentResponse{
			Name: data.NewDepartmentOrBranch,
		},
		PhoneNumber:        data.DelegatedUserPhoneNumber,
		Email:              data.DelegatedUserEmail,
		UserName:           data.DelegatedUserID,
		Realm:              data.DelegatedUserUserType,
		IsDelegationActive: true,
		DelegatedRole:      data.NewRoleID,
	}, nil
}

// FindByID implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) FindByID(ctx context.Context, id string) (*imodel.RoleDelegation, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	result, err := r.repo.FindOne(ctx, bson.M{"_id": objID}, nil)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][FindByID] failed to find role delegation: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	if result.DelegatedUserUserType == "CPS" {
		cpsUser, err := r.cpsRepo.FindOne(ctx, bson.M{"user_code": result.DelegatedUserUserCode}, bson.M{})
		if err != nil {
			log.Errorf("[RoleDelegationRepository][FindByID] failed to fetch cps user user_code %s err: %v", result.DelegatedUserUserCode, err)
			return nil, localization.ErrorUnexpectedError
		}
		result.DelegatedUserFullName = cpsUser.FullName
		result.DelegatedUserDepartmentOrBranch = cpsUser.Department.Hex()
		result.DelegatedUserEmail = cpsUser.Email
		result.DelegatedUserExistingRole = cpsUser.Role
		result.DelegatedUserID = cpsUser.UserName
		result.DelegatedUserJobTitle = cpsUser.JobTitle
		result.DelegatedUserPhoneNumber = cpsUser.PhoneNumber
	} else {
		bpsUser, err := r.bpsRepo.FindOne(ctx, bson.M{"user_code": result.DelegatedUserUserCode}, bson.M{})
		if err != nil {
			log.Errorf("[RoleDelegationRepository][FindByID] failed to fetch bps user user_code %s err: %v", result.DelegatedUserUserCode, err)
			return nil, localization.ErrorUnexpectedError
		}
		result.DelegatedUserFullName = bpsUser.FullName
		if len(bpsUser.BranchCode) > 0 {
			result.DelegatedUserDepartmentOrBranch = bpsUser.BranchCode[0]
		}
		result.DelegatedUserDepartmentOrBranchName = bpsUser.BranchName
		result.DelegatedUserEmail = bpsUser.Email
		result.DelegatedUserExistingRole = bpsUser.Role
		result.DelegatedUserExistingRoleName = bpsUser.RoleName
		result.DelegatedUserID = bpsUser.UserName
		result.DelegatedUserJobTitle = bpsUser.JobTitle
		result.DelegatedUserPhoneNumber = bpsUser.PhoneNumber
	}
	return result, nil
}

// Update implements [storage.RoleDelegationRepository].
func (r *roleDelegationRepository) Update(ctx context.Context, id string, role *imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegationRepository][Update] updating role delegation id=%s", id)
	objID, err := local_util.ParseObjectID(id)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	update := bson.M{
		"start_at":   role.StartAt,
		"end_at":     role.EndAt,
		"updated_at": time.Now(),
	}

	_, err = r.repo.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Update] failed to update role delegation: %v", err)
		return local_util.HandleDBError(err)
	}

	return nil
}

func (r *roleDelegationRepository) CheckIfDelegationAlreadyExists(ctx context.Context, userID string, start, end time.Time) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	filter := bson.M{
		"delegated_user_id": userID,
		"enabled":           true,
		"$and": []bson.M{
			{
				"start_at": bson.M{"$lte": start},
				"end_at":   bson.M{"$gte": start},
			},
		},
	}
	count, err := r.repo.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][CheckIfDelegationAlreadyExists] failed to check existing delegation: %v", err)
		return false, local_util.HandleDBError(err)
	}
	return count > 0, nil
}

func (r *roleDelegationRepository) Delete(ctx context.Context, userID string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{
		"_id": objID,
	}

	err = r.repo.DeleteOneH(ctx, filter)
	if err != nil {
		log.Errorf("[RoleDelegationRepository][Delete] failed to delete delegations: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func NewRoleDelegationRepository(client *mongo.Client, cfg *config.VaultConfig, database string, collection []string, logger utils.Logger) storage.RoleDelegationRepository {
	return &roleDelegationRepository{
		repo:              dal.NewMongoDal[imodel.RoleDelegation, imodel.RoleDelegation](client, cfg, database, collection[0]),
		cpsRepo:           dal.NewMongoDal[imodel.CPSUser, imodel.CPSUser](client, cfg, database, collection[1]),
		bpsRepo:           dal.NewMongoDal[imodel.BPSUser, imodel.BPSUser](client, cfg, database, collection[2]),
		logger:            logger,
		collection:        client.Database(database).Collection(collection[0]),
		cpsUserCollection: client.Database(database).Collection(collection[1]),
		bpsUserCollection: client.Database(database).Collection(collection[2]),
	}
}
