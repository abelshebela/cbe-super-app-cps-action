package archived_user

import (
	unlink_dto "cbe-super-app-cps-action/internal/constants/dto/unlink"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// 	"cbe-super-app-cps-action/internal/constants/localization"
// 	local_util "cbe-super-app-cps-action/pkgs/utils"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// )

type archivedUserStorage struct {
	dal    dal.MongoDal[model.ArchivedUser, model.ArchivedUser]
	client *mongo.Client
	// used for aggregation pipelines
	collection *mongo.Collection
	logger     utils.Logger
}

func NewArchivedUserRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.ArchivedUserRepository {
	return &archivedUserStorage{
		dal:        dal.NewMongoDal[model.ArchivedUser, model.ArchivedUser](client, cfg, dbName, collection),
		client:     client,
		collection: client.Database(dbName).Collection(collection),
		logger:     logger,
	}
}

func (a *archivedUserStorage) Create(ctx context.Context, user *member.User) error {

	archivedUser := UserToArchivedUser(user)
	_, err := a.dal.InsertOne(ctx, *archivedUser)
	if err != nil {
		a.logger.Errorf("[ArchivedUserStorage][Create] failed to create archived user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *archivedUserStorage) FindByID(ctx context.Context, id string) (*model.ArchivedUser, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[ArchivedUserStorage][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	result, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("[ArchivedUserStorage][FindByID] failed to find archived user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (s *archivedUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ArchivedUser], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}
	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"enabled", "kyc_level", "is_blocked", "is_verified"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
		}
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("[ArchivedUserStorage][FindAllWithPagination] failed to fetch archived users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[ArchivedUserStorage][FindAllWithPagination] failed to count archived users: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.ArchivedUser]{
		Data: data,
		Meta: meta,
	}, nil
}

func (s *archivedUserStorage) FindAllArchievedUsersWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*unlink_dto.ArchivedUserResponse], error) {
	s.logger.Infof("[ArchivedUserStorage][FindAllArchievedUsersWithPagination] fetching archived users with pipeline")

	// base and dynamic filters
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}
	allowedKeys := []string{"enabled", "kyc_level", "is_blocked", "is_verified"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	// ensure soft-delete filter is always applied
	filter["is_deleted"] = false

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         "archived_linked_account",
				"localField":   "_id",
				"foreignField": "user_id",
				"as":           "linked_account",
			},
		}},
		{{Key: "$unwind", Value: bson.M{"path": "$linked_account", "preserveNullAndEmptyArrays": true}}},
		{{
			Key: "$lookup",
			Value: bson.M{
				"from": "account_block",
				"let":  bson.M{"branch_code": "$linked_account.branch_code"},
				"pipeline": mongo.Pipeline{
					{{
						Key: "$match",
						Value: bson.D{{
							Key: "$expr",
							Value: bson.D{{
								Key: "$and",
								Value: bson.A{
									bson.D{{Key: "$eq", Value: bson.A{"$code", "$$branch_code"}}},
									bson.D{{Key: "$eq", Value: bson.A{"$type", "B"}}},
									bson.D{{Key: "$eq", Value: bson.A{"$is_deleted", false}}},
								},
							}},
						}},
					}},
					{{
						Key: "$lookup",
						Value: bson.M{
							"from":         "account_block",
							"localField":   "parent_id",
							"foreignField": "_id",
							"as":           "parent_city",
						},
					}},
					{{Key: "$unwind", Value: bson.M{"path": "$parent_city", "preserveNullAndEmptyArrays": true}}},
					{{
						Key: "$lookup",
						Value: bson.M{
							"from":         "account_block",
							"localField":   "parent_city.parent_id",
							"foreignField": "_id",
							"as":           "district",
						},
					}},
					{{Key: "$unwind", Value: bson.M{"path": "$district", "preserveNullAndEmptyArrays": true}}},
					{{
						Key: "$project",
						Value: bson.M{
							"_id":      1,
							"name":     1,
							"district": "$district",
						},
					}},
				},
				"as": "branch",
			},
		}},
		{{Key: "$unwind", Value: bson.M{"path": "$branch", "preserveNullAndEmptyArrays": true}}},
		{{
			Key: "$project",
			Value: bson.M{
				"_id":                           1,
				"customer_code":                 "$user_code",
				"customer_name":                 "$full_name",
				"branch_code":                   bson.M{"$ifNull": bson.A{"$linked_account.branch_code", "$branch_code"}},
				"branch":                        bson.M{"id": bson.M{"$ifNull": bson.A{bson.M{"$toString": "$branch._id"}, ""}}, "name": bson.M{"$ifNull": bson.A{"$branch.name", ""}}},
				"district":                      bson.M{"id": bson.M{"$ifNull": bson.A{bson.M{"$toString": "$branch.district._id"}, ""}}, "name": bson.M{"$ifNull": bson.A{"$branch.district.name", ""}}},
				"phone_number":                  "$phone_number",
				"account_number":                "$linked_account.account_number",
				"language":                      "$language",
				"avatar":                        "$avatar",
				"email":                         "$email",
				"push_token":                    "$push_token",
				"customer_number":               "$customer_number",
				"user_category":                 "$user_category",
				"industry":                      "$industry",
				"sector":                        "$sector",
				"ownership":                     "$ownership",
				"customer_segment":              "$customer_segment",
				"blocked_reason":                "$blocked_reason",
				"device_uuid":                   "$device_uuid",
				"app_version":                   "$app_version",
				"gender":                        "$gender",
				"account_branch_type":           "$account_branch_type",
				"platform":                      "$platform",
				"device_status":                 "$device_status",
				"onboarding_method":             "$onboarding_method",
				"blocked_on":                    "$blocked_on",
				"enabled_channels":              "$enabled_channels",
				"login_attempt_count":           "$login_attempt_count",
				"last_login_attempt":            "$last_login_attempt",
				"last_login":                    "$last_login",
				"application_installation_date": "$application_installation_date",
				"created_at":                    "$created_at",
				"expiry_at":                     "$expiry_at",
				"last_modified_at":              "$last_modified_at",
				"is_blocked":                    "$is_blocked",
				"enabled":                       "$enabled",
				"first_pin_set":                 "$first_pin_set",
				"is_activated":                  "$is_activated",
			},
		}},
		{{
			Key: "$facet",
			Value: bson.M{
				"data": mongo.Pipeline{
					{{Key: "$skip", Value: skip}},
					{{Key: "$limit", Value: limit}},
				},
				"total": mongo.Pipeline{
					{{Key: "$count", Value: "count"}},
				},
			},
		}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		s.logger.Errorf("[ArchivedUserStorage][FindAllArchievedUsersWithPagination] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data  []*unlink_dto.ArchivedUserResponse `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		s.logger.Errorf("[ArchivedUserStorage][FindAllArchievedUsersWithPagination] decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Handle empty aggregation result
	if len(results) == 0 {
		meta := local_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage)
		return &types.PaginatedResponse[[]*unlink_dto.ArchivedUserResponse]{Data: []*unlink_dto.ArchivedUserResponse{}, Meta: meta}, nil
	}

	total := int64(0)
	if len(results[0].Total) > 0 {
		total = results[0].Total[0].Count
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*unlink_dto.ArchivedUserResponse]{
		Data: results[0].Data,
		Meta: meta,
	}, nil
}
