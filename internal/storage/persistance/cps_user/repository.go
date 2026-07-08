package cps_user

import (
	"cbe-super-app-cps-action/internal/constants"
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"regexp"

	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSUserStorage struct {
	dal                  dal.MongoDal[imodel.CPSUser, imodel.CPSUser]
	role_delegations_dal dal.MongoDal[imodel.RoleDelegation, imodel.RoleDelegation]
	cpsAction            dal.MongoDal[imodel.CPSAction, imodel.CPSAction]
	client               *mongo.Client
	collection           *mongo.Collection
	redisRepository      storage.RedisRepository
	relatedCollection    []string
	logger               utils.Logger
}

func NewCPSUserRepository(client *mongo.Client, redisRepository storage.RedisRepository, cfg *config.VaultConfig, dbName string, collection string, relatedCollection []string, logger utils.Logger) storage.CpsUserRepository {
	return &CPSUserStorage{
		dal:                  dal.NewMongoDal[imodel.CPSUser, imodel.CPSUser](client, cfg, dbName, collection),
		role_delegations_dal: dal.NewMongoDal[imodel.RoleDelegation, imodel.RoleDelegation](client, cfg, dbName, relatedCollection[6]),
		cpsAction:            dal.NewMongoDal[imodel.CPSAction, imodel.CPSAction](client, cfg, dbName, "cps_actions"),
		client:               client,
		redisRepository:      redisRepository,
		collection:           client.Database(dbName).Collection(collection),
		relatedCollection:    relatedCollection,
		logger:               logger,
	}
}

func (r *CPSUserStorage) FindForExport(ctx context.Context, startDate, endDate time.Time, userName string) ([]imodel.ExportCPSUser, error) {
	matchFilter := bson.M{
		"created_at": bson.M{"$gte": startDate, "$lte": endDate},
		"is_deleted": bson.M{"$ne": true},
	}
	if userName != "" {
		matchFilter["username"] = bson.M{"$regex": "^" + regexp.QuoteMeta(userName) + "$", "$options": "i"}
	}

	pipeline := mongo.Pipeline{
		// Stage 1: Match users by date range
		bson.D{{Key: "$match", Value: matchFilter}},

		// Stage 2: Lookup department name
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "department",
			"let":  bson.M{"dept_id": "$department"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$eq": []interface{}{"$_id", "$$dept_id"}},
				}}},
				bson.D{{Key: "$project", Value: bson.M{"department": 1}}},
			},
			"as": "department_info",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$department_info",
			"preserveNullAndEmptyArrays": true,
		}}},

		// Stage 3: Lookup active delegation (enable==true and not yet expired)
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "role_delegations",
			"let":  bson.M{"uid": "$username"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{
						"$and": bson.A{
							bson.M{"$eq": []interface{}{bson.M{"$toLower": "$delegated_user_id"}, bson.M{"$toLower": "$$uid"}}},
							bson.M{"$eq": []interface{}{"$enable", true}},
							bson.M{"$lte": []interface{}{"$start_at", "$$NOW"}},
							bson.M{"$gt": []interface{}{"$end_at", "$$NOW"}},
						},
					},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"end_at": 1}}},
			},
			"as": "delegation_info",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$delegation_info",
			"preserveNullAndEmptyArrays": true,
		}}},

		// Stage 4: Lookup last modification action (UPDATE/DELETE/ENABLE/DISABLE after CREATE)
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "cps_actions",
			"let":  bson.M{"ucode": "$user_code"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{
						"$and": bson.A{
							bson.M{"$in": []interface{}{"$request_action", bson.A{
								"CREATE_CPS_USER", "UPDATE_CPS_USER", "DELETE_CPS_USER",
								"ENABLE_CPS_USER", "DISABLE_CPS_USER",
							}}},
							bson.M{"$eq": []interface{}{"$previous_action.user_code", "$$ucode"}},
						},
					},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"request_action": 1, "created_at": 1}}},
			},
			"as": "last_modification_info",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$last_modification_info",
			"preserveNullAndEmptyArrays": true,
		}}},

		// Stage 5: Lookup the original CREATE action for created_by and approved_by
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "cps_actions",
			"let":  bson.M{"ucode": "$user_code"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{
						"$and": bson.A{
							bson.M{"$eq": []interface{}{"$request_action", "CREATE_CPS_USER"}},
							bson.M{"$eq": []interface{}{"$action_status", "APPROVED"}},
							bson.M{"$eq": []interface{}{"$current_action.user_code", "$$ucode"}},
						},
					},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"maker_id": 1, "checker_users": 1}}},
			},
			"as": "create_action_info",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$create_action_info",
			"preserveNullAndEmptyArrays": true,
		}}},

		// Stage 6: Project into ExportCPSUser shape
		bson.D{{Key: "$project", Value: bson.M{
			"_id":                        0,
			"first_name":                 "$full_name",
			"phone_number":               1,
			"email":                      1,
			"department":                 "$department_info.department",
			"job_title":                  1,
			"role":                       1,
			"username":                   1,
			"created_at":                 1,
			"enabled":                    1,
			"last_login":                 1,
			"expiry_date_for_delegation": "$delegation_info.end_at",
			"last_modification_action":   "$last_modification_info.request_action",
			"last_modified":              "$last_modification_info.created_at",
			"created_by": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$ne": bson.A{"$create_action_info.maker_id", ""}},
					"then": "$create_action_info.maker_id",
					"else": "SSO",
				},
			},
			"approved_by": bson.M{
				"$ifNull": []interface{}{
					bson.M{"$arrayElemAt": []interface{}{"$create_action_info.checker_users.checker_id", 0}},
					"",
				},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var users []imodel.ExportCPSUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return users, nil
}

// Implement actual repository methods for CPS action authorization
func (r *CPSUserStorage) Create(ctx context.Context, cpsUser *imodel.CPSUser) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	// cpsUserMap := CPSUserMapper(*cpsUser)
	current_time := time.Now()
	cpsUser.DateJoined = &current_time
	_, err := r.dal.InsertOne(ctx, *cpsUser)
	if err != nil {
		log.Errorf("[CPSUserStorage][Create] failed to create CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *CPSUserStorage) Update(ctx context.Context, userCode string, cpsUser *imodel.CPSUser) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][Update] updating CPS user")
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := CPSUserUpdateMapper(cpsUser)

	session, err := r.client.StartSession()
	if err != nil {
		log.Errorf("[CPSUserStorage][Update] failed to start session: %v", err)
		return localization.ErrorUnexpectedError
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (any, error) {
		if _, err := r.dal.UpdateOne(sc, filter, update); err != nil {
			log.Errorf("[CPSUserStorage][Update] failed to update CPS user: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		if _, err := r.role_delegations_dal.UpdateOne(sc, bson.M{"delegated_user_user_code": cpsUser.UserCode}, bson.M{
			"delegated_user_department_or_branch": cpsUser.DelegationID,
			"delegated_user_email":                cpsUser.Email,
			"delegated_user_existing_role":        cpsUser.Role,
			"delegated_user_full_name":            cpsUser.FullName,
			"delegated_user_id":                   cpsUser.UserName,
			"delegated_user_job_title":            cpsUser.JobTitle,
			"delegated_user_phone_number":         cpsUser.PhoneNumber,
		}); err != nil {
			log.Errorf("[CPSUserStorage][Update] failed to update role delegation: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		return nil, nil
	})
	if err != nil {
		log.Errorf("[CPSUserStorage][Update] transaction failed: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[CPSUserStorage][Update] CPS user updated successfully")
	return nil
}

func (r *CPSUserStorage) Delete(ctx context.Context, userCode string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][Delete] deleting CPS user")
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"is_deleted": true}

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[CPSUserStorage][Delete] failed to delete CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Infof("[CPSUserStorage][Delete] CPS user deleted successfully")
	return nil
}

func (r *CPSUserStorage) EnableOrDisable(ctx context.Context, userCode string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][EnableOrDisable] processing CPS user enable/disable, enabled: %v", enable)
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified": time.Now()}
	if enable {
		update["login_attempt_count"] = 0
		update["is_first_time_login"] = true
	}
	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[CPSUserStorage][EnableOrDisable] failed to enable/disable CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	// remove user device id from redis
	user, err := r.FindByID(ctx, userCode)
	if err == nil {
		if !enable {
			objID := user.ID.Hex()
			id := local_util.FirstHex24(objID)
			if err := r.redisRepository.Delete(ctx, fmt.Sprintf("%s:%s", constants.RedisCPSUserDeviceIDPrefix, id)); err != nil {
				log.Errorf("[CPSUserStorage][EnableOrDisable] failed to delete user device id from redis: %v", err)
			}
		}
	}
	log.Infof("[CPSUserStorage][EnableOrDisable] CPS user enable/disable completed successfully")
	return nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByUsername(ctx context.Context, username string) (*imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][FindByUsername] searching for CPS user by username")

	filter := bson.M{
		"username": bson.M{
			"$regex":   "^" + regexp.QuoteMeta(username) + "$",
			"$options": "i",
		},
	}

	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[CPSUserStorage][FindByUsername] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[CPSUserStorage][FindByUsername] CPS user retrieved successfully")
	return result, nil
}
func (r *CPSUserStorage) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][FindByPhoneNumber] searching for CPS user by phone number")
	filter := bson.M{"phone_number": phoneNumber}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Infof("[CPSUserStorage][FindByPhoneNumber] CPS user not found")
			return nil, nil
		}
		log.Errorf("[CPSUserStorage][FindByPhoneNumber] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[CPSUserStorage][FindByPhoneNumber] CPS user retrieved successfully")
	return result, nil
}
func (r *CPSUserStorage) FindByEmail(ctx context.Context, email string) (*imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][FindByEmail] searching for CPS user by email")
	filter := bson.M{"email": email}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Infof("[CPSUserStorage][FindByEmail] CPS user not found")
			return nil, nil
		}
		log.Errorf("[CPSUserStorage][FindByEmail] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[CPSUserStorage][FindByEmail] CPS user retrieved successfully")
	return result, nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByID(ctx context.Context, id string) (*imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][FindByID] fetching CPS user by id")
	filter := bson.M{"user_code": id, "is_deleted": false}

	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[CPSUserStorage][FindByID] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[CPSUserStorage][FindByID] CPS user retrieved successfully")
	return result, nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByUserID(ctx context.Context, id string) (*imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[CPSUserStorage][FindByID] fetching CPS user by id")
	log.Infof("[CPSUserStorage][FindByID] fetching CPS user by id")
	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[CPSUserStorage][FindByID] failed to convert user id to object id")
		return nil, localization.ErrorUnexpectedError
	}
	filter := bson.M{"_id": obj, "is_deleted": false}

	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[CPSUserStorage][FindByID] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[CPSUserStorage][FindByID] CPS user retrieved successfully")
	return result, nil
}

func (r *CPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	searchKeys := bson.M{}

	allowedKeys := []string{"enabled", "department", "role", "created_at"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"email": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
			{"job_title": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},

		bson.D{{Key: "$project", Value: bson.M{
			"_id":           1,
			"user_code":     1,
			"full_name":     1,
			"job_title":     1,
			"department":    1,
			"role":          1,
			"gender":        1,
			"phone_number":  1,
			"email":         1,
			"username":      1,
			"enabled":       1,
			"date_joined":   1,
			"last_modified": 1,
			"country":       1,
			"region":        1,
		}}},

		// 🔥 FIX: convert department string → ObjectId
		bson.D{{Key: "$addFields", Value: bson.M{
			"department_obj_id": bson.M{
				"$convert": bson.M{
					"input":   "$department",
					"to":      "objectId",
					"onError": nil,
					"onNull":  nil,
				},
			},
		}}},

		// 🔹 Role lookup
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "roles",
			"localField":   "job_title",
			"foreignField": "job_title",
			"as":           "role_info",
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$project", Value: bson.M{
					"_id":  1,
					"role": 1,
				}}},
			},
		}}},

		// 🔹 Department lookup (FIXED)
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "department",
			"localField":   "department_obj_id", // ✅ FIXED
			"foreignField": "_id",
			"as":           "department_info",
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$project", Value: bson.M{
					"_id":        1,
					"department": 1,
				}}},
			},
		}}},

		// Active delegation lookup by username
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "role_delegations",
			"let":  bson.M{"uid": "$username"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{
						"$and": bson.A{
							bson.M{"$eq": []interface{}{bson.M{"$toLower": "$delegated_user_id"}, bson.M{"$toLower": "$$uid"}}},
							bson.M{"$lte": []interface{}{"$start_at", "$$NOW"}},
							bson.M{"$gt": []interface{}{"$end_at", "$$NOW"}},
						},
					},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"new_role_id": 1}}},
			},
			"as": "delegation_info",
		}}},

		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$role_info",
			"preserveNullAndEmptyArrays": true,
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$department_info",
			"preserveNullAndEmptyArrays": true,
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$delegation_info",
			"preserveNullAndEmptyArrays": true,
		}}},

		bson.D{{Key: "$addFields", Value: bson.M{
			"role":           "$role_info.role",
			"delegated_role": bson.M{"$ifNull": []interface{}{"$delegation_info.new_role_id", ""}},
			"department": bson.M{
				"id":   "$department_info._id",
				"name": "$department_info.department",
			},
		}}},

		// 🔹 Clean temp + lookup fields
		bson.D{{Key: "$project", Value: bson.M{
			"role_info":         0,
			"department_info":   0,
			"delegation_info":   0,
			"department_obj_id": 0, // ✅ remove temp field
		}}},

		// 🔹 Ensure unique users
		bson.D{{Key: "$group", Value: bson.M{
			"_id": "$user_code",
			"doc": bson.M{"$first": "$$ROOT"},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.M{
			"newRoot": "$doc",
		}}},

		// 🔹 Pagination
		bson.D{{Key: "$facet", Value: bson.M{
			"data": []bson.D{
				{{Key: "$sort", Value: bson.D{{Key: "date_joined", Value: -1}}}},
				{{Key: "$skip", Value: skip}},
				{{Key: "$limit", Value: limit}},
			},
			"total": []bson.D{
				{{Key: "$count", Value: "count"}},
			},
		}}},
	}

	log.Infof("[CPSUserStorage][FindAllWithPagination] fetching CPS users with pagination")
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[CPSUserStorage][FindAllWithPagination] failed to execute aggregation pipeline: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data  []*cpsuser.CPSUserWithDepartment `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		log.Errorf("[CPSUserStorage][FindAllWithPagination] failed to decode aggregation results: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(results) == 0 {
		return &types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment]{
			Data: []*cpsuser.CPSUserWithDepartment{},
			Meta: local_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage),
		}, nil
	}

	var total int64
	if len(results[0].Total) > 0 {
		total = results[0].Total[0].Count
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	log.Infof("[CPSUserStorage][FindAllWithPagination] retrieved %d CPS users", len(results[0].Data))
	return &types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment]{
		Data: results[0].Data,
		Meta: meta,
	}, nil
}

func (r *CPSUserStorage) GetPopulatedByID(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][GetPopulatedByID] fetching populated CPS user")
	// pipeline := PipelineBuilder(userCode, r.relatedCollection[0], r.relatedCollection[3], r.relatedCollection[1], r.relatedCollection[2])
	pipeline := PipelineBuilder(userCode)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[CPSUserStorage][GetPopulatedByID] failed to aggregate CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		log.Errorf("[CPSUserStorage][GetPopulatedByID] CPS user not found")
		return nil, errors.New(localization.ErrorUserNotFound.Code)
	}

	var resp cpsuser.CpsUserResponse
	if err := cursor.Decode(&resp); err != nil {
		log.Errorf("[CPSUserStorage][GetPopulatedByID] failed to decode CPS user response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	fmt.Println("Decoded CPS User Response:", resp)
	log.Infof("[CPSUserStorage][GetPopulatedByID] populated CPS user retrieved successfully")
	return &resp, nil
}

func (r *CPSUserStorage) GetPopulatedWithRole(ctx context.Context, userCode string) (*cpsuser.CpsUserPopulatedResponse, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][GetPopulatedWithRole] fetching populated CPS user")
	// relatedCollection: [DepartmentsCollection, PermissionCollection, PermissionCategoryCollection, PermissionGroupsCollection, RolesCollection, JobRolesCollection]
	pipeline := PipelineBuilderWithRole(userCode, r.relatedCollection[0], r.relatedCollection[4], r.relatedCollection[5])

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[CPSUserStorage][GetPopulatedWithRole] failed to aggregate CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		log.Errorf("[CPSUserStorage][GetPopulatedWithRole] CPS user not found")
		return nil, localization.ErrorUserNotFound
	}

	var resp cpsuser.CpsUserPopulatedResponse
	if err := cursor.Decode(&resp); err != nil {
		log.Errorf("[CPSUserStorage][GetPopulatedWithRole] failed to decode CPS user response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	log.Infof("[CPSUserStorage][GetPopulatedWithRole] populated CPS user retrieved successfully")
	return &resp, nil
}

func (r *CPSUserStorage) GetByUserCode(ctx context.Context, userCode string) (imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[CPSUserStorage][GetPopulatedByUserName] fetching populated CPS user")

	data, err := r.dal.FindOne(ctx, bson.M{"user_code": userCode, "is_deleted": false}, nil)
	if err != nil {
		log.Errorf("[CPSUserStorage][GetPopulatedByUserName] failed to find CPS user: %v", err)
		return imodel.CPSUser{}, local_util.HandleDBError(err)
	}

	return *data, nil
}

func (r *CPSUserStorage) GetPopulatedWithRoleByUserName(ctx context.Context, userName string) (*cpsuser.CpsUserPopulatedResponse, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][GetPopulatedWithRoleByUserName] fetching populated CPS user")
	// relatedCollection: [DepartmentsCollection, PermissionCollection, PermissionCategoryCollection, PermissionGroupsCollection, RolesCollection, JobRolesCollection]
	pipeline := PipelineBuilderWithRoleByUserName(userName, r.relatedCollection[0], r.relatedCollection[4], r.relatedCollection[5])

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[CPSUserStorage][GetPopulatedWithRoleByUserName] failed to aggregate CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		log.Errorf("[CPSUserStorage][GetPopulatedWithRoleByUserName] CPS user not found")
		return nil, errors.New(localization.ErrorUserNotFound.Code)
	}

	var resp cpsuser.CpsUserPopulatedResponse
	if err := cursor.Decode(&resp); err != nil {
		log.Errorf("[CPSUserStorage][GetPopulatedWithRoleByUserName] failed to decode CPS user response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	log.Infof("[CPSUserStorage][GetPopulatedWithRoleByUserName] populated CPS user retrieved successfully")
	return &resp, nil
}

// FindByEmailOrPhoneNumberOrUserName implements [storage.CpsUserRepository].
func (r *CPSUserStorage) FindByEmailOrPhoneNumberOrUserName(ctx context.Context, email string, phoneNumber string, username string) (*imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][FindByEmailOrPhoneNumberOrUserName] searching for CPS user by email, phone number, or username")
	var orFilters []bson.M
	if email != "" {
		orFilters = append(orFilters, bson.M{"email": email})
	}
	if phoneNumber != "" {
		orFilters = append(orFilters, bson.M{"phone_number": phoneNumber})
	}
	if username != "" {
		orFilters = append(orFilters, bson.M{"username": username})
	}
	if len(orFilters) == 0 {
		return nil, nil
	}

	if len(orFilters) == 0 {
		log.Infof("[CPSUserStorage][FindByEmailOrPhoneNumberOrUserName] no search parameters provided, returning nil")
		return nil, nil
	}
	filter := bson.M{"$or": orFilters}

	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[CPSUserStorage][FindByEmailOrPhoneNumberOrUserName] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (r *CPSUserStorage) UpdateCpsUsersJobTitle(ctx context.Context, oldJobTitle, newJobTitle string) error {
	// Define the filter to match documents with the old job title
	filter := bson.M{"job_title": oldJobTitle}

	// Define the update operation to set the new job title
	update := bson.M{"$set": bson.M{"job_title": newJobTitle}}

	// Perform the update operation
	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	// Log the number of documents updated (optional)
	fmt.Printf("Updated %d documents in collection %s\n", result.ModifiedCount, r.collection.Name())

	return nil
}

func (r *CPSUserStorage) GetUserByDepartment(ctx context.Context, department string) (*imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][GetUserByDepartment] fetching CPS user by department: %s", department)
	objID, err := bson.ObjectIDFromHex(department)
	if err != nil {
		log.Errorf("[CPSUserStorage][GetUserByDepartment] invalid department id: %v", err)
		return nil, localization.ErrorUnexpectedError
	}

	filter := bson.M{"department": objID, "is_deleted": false}

	cpsUser, err := r.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Infof("[CPSUserStorage][GetUserByDepartment] no CPS users found for department: %s", department)
			return nil, nil
		}
		return nil, err
	}
	return cpsUser, nil

}

func (r *CPSUserStorage) GetUserByJobTitle(ctx context.Context, jobTitle string) (*imodel.CPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSUserStorage][GetUserByJobTitle] fetching CPS user by job title: %s", jobTitle)

	filter := bson.M{"job_title": jobTitle, "is_deleted": false}

	cpsUser, err := r.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Infof("[CPSUserStorage][GetUserByJobTitle] no CPS users found for job title: %s", jobTitle)
			return nil, nil
		}
		return nil, err
	}
	return cpsUser, nil

}
