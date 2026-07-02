package bps_user

import (
	// "cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	// "time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/local_model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	bpsUserDto "cbe-super-app-cps-action/internal/constants/dto/bps_user"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// bpsUserActiveFilter matches non–soft-deleted users. Use $ne: true so documents
// without is_deleted (legacy) or with false are included; only true is excluded.
func bpsUserActiveFilter() bson.M {
	return bson.M{"is_deleted": bson.M{"$ne": true}}
}

type BPSUserStorage struct {
	dal        dal.MongoDal[bps_model.BPSUser, bps_model.BPSUser]
	client     *mongo.Client
	logger     utils.Logger
	collection *mongo.Collection
}

func NewBPSUserRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.BPSUserRepository {
	return &BPSUserStorage{
		dal:        dal.NewMongoDal[bps_model.BPSUser, bps_model.BPSUser](client, cfg, dbName, collection),
		client:     client,
		logger:     logger,
		collection: client.Database(dbName).Collection(collection),
	}
}

func (b *BPSUserStorage) FindForExport(ctx context.Context, startDate, endDate time.Time, userName string) ([]imodel.ExportBPSUser, error) {
	filter := bpsUserActiveFilter()
	filter["created_at"] = bson.M{"$gte": startDate, "$lte": endDate}
	if userName != "" {
		filter["username"] = bson.M{"$regex": "^" + regexp.QuoteMeta(userName) + "$", "$options": "i"}
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "roles",
			"let":  bson.M{"job_title_value": "$job_title"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$or": bson.A{
						bson.M{"$eq": bson.A{"$job_title", "$$job_title_value"}},
						bson.M{"$eq": bson.A{"$name", "$$job_title_value"}},
						bson.M{"$eq": bson.A{"$role", "$$job_title_value"}},
						bson.M{"$eq": bson.A{"$code", "$$job_title_value"}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"role": 1, "name": 1, "job_title": 1, "code": 1}}},
			},
			"as": "role_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "role_delegations",
			"let":  bson.M{"uname": "$username"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$eq": []interface{}{bson.M{"$toLower": "$delegated_user_id"}, bson.M{"$toLower": "$$uname"}}},
						bson.M{"$eq": []interface{}{"$enable", true}},
						bson.M{"$lte": []interface{}{"$start_at", "$$NOW"}},
						bson.M{"$gt": []interface{}{"$end_at", "$$NOW"}},
					}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"end_at": 1}}},
			},
			"as": "delegation_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "cps_actions",
			"let":  bson.M{"user_id": "$user_code"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$in": bson.A{"$request_action", bson.A{
							string(constants.RequestCreateBPSUser),
							string(constants.RequestBpsUserUpdate),
							string(constants.RequestBpsUserDelete),
							string(constants.RequestEnableBPSUser),
							string(constants.RequestDisableBPSUser),
						}}},
						bson.M{"$eq": bson.A{"$previous_action.user_code", "$$user_id"}},
					}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"request_action": 1, "created_at": 1}}},
			},
			"as": "last_action_info",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "cps_actions",
			"let":  bson.M{"user_id": "$user_code"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$request_action", string(constants.RequestCreateBPSUser)}},
						bson.M{"$eq": bson.A{"$action_status", "APPROVED"}},
						bson.M{"$eq": bson.A{"$current_action.user_code", "$$user_id"}},
					}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"maker_id": 1, "checker_users": 1}}},
			},
			"as": "create_action_info",
		}}},

		bson.D{{Key: "$unwind", Value: bson.M{"path": "$role_info", "preserveNullAndEmptyArrays": true}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$delegation_info", "preserveNullAndEmptyArrays": true}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$last_action_info", "preserveNullAndEmptyArrays": true}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$create_action_info", "preserveNullAndEmptyArrays": true}}},

		bson.D{{Key: "$project", Value: bson.M{
			"_id":          0,
			"first_name":   "$full_name",
			"phone_number": 1,
			"email":        1,
			"branch_name":  bson.M{"$ifNull": bson.A{"$branch_name", bson.M{"$arrayElemAt": bson.A{"$branch_code", 0}}}},
			"branch_code": bson.M{
				"$ifNull": bson.A{
					bson.M{
						"$arrayElemAt": bson.A{"$branch_code", 0},
					},
					"",
				},
			},
			"job_title": 1,
			"role": bson.M{"$ifNull": bson.A{
				"$role_info.role",
				bson.M{"$ifNull": bson.A{
					"$role_info.name",
					bson.M{"$ifNull": bson.A{
						"$role_info.code",
						"$role",
					}},
				}},
			}},
			"username":                   1,
			"created_at":                 1,
			"expiry_date_for_delegation": "$delegation_info.end_at",
			"last_modification_action":   "$last_action_info.request_action",
			"last_modified":              "$last_action_info.created_at",
			"created_by": bson.M{"$cond": bson.M{
				"if":   bson.M{"$ne": bson.A{"$create_action_info.maker_id", ""}},
				"then": "$create_action_info.maker_id",
				"else": "SSO",
			}},
			"approved_by": bson.M{"$ifNull": bson.A{
				bson.M{"$arrayElemAt": bson.A{"$create_action_info.checker_users.checker_id", 0}},
				"",
			}},
			"enabled":    1,
			"last_login": 1,
		}}},
	}

	cursor, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var users []imodel.ExportBPSUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return users, nil
}

func (b *BPSUserStorage) FindByOr(ctx context.Context, phone, email, username string) (*bps_model.BPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][FindByOr] searching for BPS user by phone, email, or username")
	// Build conditions dynamically, only for non-empty parameters
	conditions := []bson.M{}

	if phone != "" {
		conditions = append(conditions, bson.M{
			"phone_number": bson.M{"$regex": "^" + regexp.QuoteMeta(phone) + "$", "$options": "i"},
		})
	}

	if username != "" {
		conditions = append(conditions, bson.M{
			"username": bson.M{"$regex": "^" + regexp.QuoteMeta(username) + "$", "$options": "i"},
		})
	}

	if email != "" {
		conditions = append(conditions, bson.M{
			"email": bson.M{"$regex": "^" + regexp.QuoteMeta(email) + "$", "$options": "i"},
		})
	}
	if len(conditions) == 0 {
		log.Infof("[BPSUserStorage][FindByOr] no search parameters provided, returning nil")
		return nil, nil
	}
	filter := bson.M{"$or": conditions}
	data, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[BPSUserStorage][FindByOr] failed to find bps user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return data, nil
}

func (b *BPSUserStorage) GetByUserCode(ctx context.Context, userCode string) (*bpsUserDto.BPSUserResposenDTO, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][GetByUserCode] fetching BPS user by user code")

	pipeline := mongo.Pipeline{
		// Match by user_code and not deleted
		bson.D{{Key: "$match", Value: bson.M{
			"user_code":  userCode,
			"is_deleted": bson.M{"$ne": true},
		}}},

		// Lookup role from roles collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "roles"},
			{Key: "localField", Value: "job_title"},
			{Key: "foreignField", Value: "job_title"},
			{Key: "as", Value: "roles"},
		}}},

		// Extract role code
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "role_code", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$roles.role", 0}},
			}},
		}}},

		// Lookup job role details
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "job_roles"},
			{Key: "localField", Value: "role_code"},
			{Key: "foreignField", Value: "code"},
			{Key: "as", Value: "job_roles"},
		}}},

		// Add user_role field
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "role", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$job_roles.name", 0}},
			}},
		}}},

		// Clean up temporary fields
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "roles", Value: 0},
			{Key: "job_roles", Value: 0},
			{Key: "role_code", Value: 0},
		}}},
	}

	cursor, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[BPSUserStorage][GetByUserCode] failed aggregation: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var results []bpsUserDto.BPSUserResposenDTO
	if err := cursor.All(ctx, &results); err != nil {
		log.Errorf("[BPSUserStorage][GetByUserCode] failed to decode user: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	if len(results) == 0 {
		log.Infof("[BPSUserStorage][GetByUserCode] no BPS user found with user_code: %s", userCode)
		return nil, errors.New(localization.ErrorUserNotFound.Code)
	}

	log.Infof("[BPSUserStorage][GetByUserCode] BPS user retrieved successfully with role")
	return &results[0], nil
}
func (b *BPSUserStorage) GetByUserID(ctx context.Context, userID string) (*bps_model.BPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][GetByUserID] fetching BPS user by user ID")
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Errorf("[BPSUserStorage][GetByUserID] invalid user ID format: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": bson.M{"$ne": true}}
	result, err := b.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		log.Errorf("[BPSUserStorage][GetByUserID] failed to find BPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[BPSUserStorage][GetByUserID] BPS user retrieved successfully")
	return result, nil
}

func (b *BPSUserStorage) GetByUsername(ctx context.Context, userName string) (*bps_model.BPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][GetByUsername] fetching BPS user by username: %s", userName)
	filter := bson.M{"username": userName, "is_deleted": bson.M{"$ne": true}}
	result, err := b.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		log.Errorf("[BPSUserStorage][GetByUsername] failed to find BPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[BPSUserStorage][GetByUsername] BPS user retrieved successfully")
	return result, nil
}

func (s *BPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	// Base match: only active users (include legacy docs without is_deleted)
	match := bpsUserActiveFilter()
	// Search filters
	if enabledVal, ok := filterParam.Filters["enabled"]; ok {
		match["enabled"] = enabledVal
	}
	filterParam.Search = strings.TrimSpace(filterParam.Search)
	if strings.HasPrefix(filterParam.Search, "09") || strings.HasPrefix(filterParam.Search, "07") {
		filterParam.Search = filterParam.Search[1:]
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		match["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
			{"email": searchRegex},
			{"branch_name": searchRegex},
		}
	}

	// Count total documents (for pagination metadata)
	total, err := s.collection.CountDocuments(ctx, match)
	if err != nil {
		log.Errorf("[BPSUserStorage][FindAllWithPagination] failed to count users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	offset := filterParam.PerPage * (filterParam.Page - 1)
	if offset >= int(total) {
		// No data for this page, return empty result
		meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
		log.Infof("[BPSUserStorage][FindAllWithPagination] offset >= total, returning empty result")
		return &types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO]{
			Data: []bpsUserDto.BPSUserResposenDTO{},
			Meta: meta,
		}, nil
	}
	if offset+filterParam.PerPage > int(total) {
		filterParam.PerPage = int(total) - offset
	}

	pipeline := mongo.Pipeline{

		bson.D{{Key: "$match", Value: match}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "roles"},
			{Key: "localField", Value: "job_title"},
			{Key: "foreignField", Value: "job_title"},
			{Key: "as", Value: "roles"},
		}}},

		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "role_code", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$roles.role", 0}}, // Make sure field is "role" not "code"
			}},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "job_roles"},
			{Key: "localField", Value: "role_code"}, // MUST use the extracted role_code
			{Key: "foreignField", Value: "code"},    // matches job_roles.code
			{Key: "as", Value: "job_roles"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "role_delegations",
			"let":  bson.M{"uname": "$username"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{"$and": bson.A{
						bson.M{"$eq": []interface{}{bson.M{"$toLower": "$delegated_user_id"}, bson.M{"$toLower": "$$uname"}}},
						bson.M{"$lte": []interface{}{"$start_at", "$$NOW"}},
						bson.M{"$gt": []interface{}{"$end_at", "$$NOW"}},
					}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.M{"new_role_id": 1}}},
			},
			"as": "delegation_info",
		}}},

		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$delegation_info",
			"preserveNullAndEmptyArrays": true,
		}}},

		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "user_role", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$job_roles.name", 0}},
			}},
			{Key: "delegated_role", Value: bson.D{
				{Key: "$ifNull", Value: bson.A{"$delegation_info.new_role_id", ""}},
			}},
		}}},

		bson.D{{Key: "$project", Value: bson.D{
			{Key: "roles", Value: 0},
			{Key: "job_roles", Value: 0},
			{Key: "role_code", Value: 0},
			{Key: "delegation_info", Value: 0},
		}}},

		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "created_at", Value: -1},
		}}},
		bson.D{{Key: "$skip", Value: filterParam.PerPage * (filterParam.Page - 1)}},
		bson.D{{Key: "$limit", Value: filterParam.PerPage}},
	}

	//  Execute aggregation
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[BPSUserStorage][FindAllWithPagination] failed aggregation: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var users []bpsUserDto.BPSUserResposenDTO
	if err := cursor.All(ctx, &users); err != nil {
		log.Errorf("[BPSUserStorage][FindAllWithPagination] failed to decode users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	//  Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	log.Infof("[BPSUserStorage][FindAllWithPagination] retrieved %d BPS users", len(users))

	//  Return paginated response
	return &types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO]{
		Data: users,
		Meta: meta,
	}, nil
}

func (b *BPSUserStorage) Update(ctx context.Context, BpsUser *bps_model.BPSUser) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][Update] updating BPS user")
	filter := bson.M{"user_code": BpsUser.UserCode, "is_deleted": bson.M{"$ne": true}}
	_, err := b.dal.UpdateOne(ctx, filter, BPSUserMapper(*BpsUser))
	if err != nil {
		log.Errorf("[BPSUserStorage][Update] failed to update BPS user: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[BPSUserStorage][Update] BPS user updated successfully")
	return nil

}

func (b *BPSUserStorage) EnableDisableBPSUser(ctx context.Context, userCode string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)
	log.Infof("[BPSUserStorage][Update] updating BPS user")
	// objID, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	log.Errorf("[BPSUserStorage][Update] invalid ObjectID: %v", err)
	// 	return errors.New(localization.ErrorInvalidID.Code)
	// }

	filter := bson.M{"user_code": userCode, "is_deleted": bson.M{"$ne": true}}
	var update bson.M
	if enable {
		update = bson.M{
			"enabled":             true,
			"last_modified_at":    time.Now(),
			"login_attempt_count": 0,
			"otp_verify_count":    0,
			"max_pass_attempt":    false,
		}
	} else {
		update = bson.M{"enabled": false, "last_modified_at": time.Now()}
	}
	_, err := b.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[BPSUserStorage][Update] failed to update BPS user: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[BPSUserStorage][Update] BPS user updated successfully")
	return nil

}

func (b *BPSUserStorage) Create(ctx context.Context, req bps_model.BPSUser) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateBPSUser", "BPS User", "CreateBPSUser")
	defer span.End()

	log.Infof("[BPSUserStorage][Create] creating BPS user with user_code: %s", req.UserCode)
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()
	req.IsFirstTimeLogin = true
	req.IsDeleted = false

	raw, err := bson.Marshal(req)
	if err != nil {
		log.Errorf("[BPSUserStorage][Create] failed to marshal BPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	var doc bson.M
	if err := bson.Unmarshal(raw, &doc); err != nil {
		log.Errorf("[BPSUserStorage][Create] failed to unmarshal BPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	doc["created_by"] = "CPS Portal"

	if _, err := b.collection.InsertOne(ctx, doc); err != nil {
		log.Errorf("[BPSUserStorage][Create] failed to create BPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	log.Infof("[BPSUserStorage][Create] BPS user created successfully for user_code: %s", req.UserCode)
	return nil
}

func (b *BPSUserStorage) FindByFilterKey(ctx context.Context, field, value string) (*bps_model.BPSUser, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	log.Infof("[BPSUserStorage][FindByFilterKey] searching BPS user by %s: %s", field, value)
	var filter bson.M

	if field == "id" {
		field = "_id"
		objID, err := bson.ObjectIDFromHex(value)
		if err != nil {
			log.Errorf("[BPSUserStorage][FindByFilterKey] invalid ObjectID: %v", err)
			return nil, errors.New(localization.ErrorInvalidID.Code)
		}
		filter = bson.M{field: objID, "is_deleted": bson.M{"$ne": true}}
	} else {
		filter = bson.M{field: value, "is_deleted": bson.M{"$ne": true}}
	}

	result, err := b.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		log.Errorf("[BPSUserStorage][FindByFilterKey] failed to find BPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

// func BPSUserMapper(user bps_model.BPSUser) bson.M {
