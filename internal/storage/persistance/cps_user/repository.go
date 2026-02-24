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

	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSUserStorage struct {
	dal               dal.MongoDal[imodel.CPSUser, imodel.CPSUser]
	cpsAction         dal.MongoDal[imodel.CPSAction, imodel.CPSAction]
	client            *mongo.Client
	collection        *mongo.Collection
	redisRepository   storage.RedisRepository
	relatedCollection []string
	logger            utils.Logger
}

func NewCPSUserRepository(client *mongo.Client, redisRepository storage.RedisRepository, cfg *config.VaultConfig, dbName string, collection string, relatedCollection []string, logger utils.Logger) storage.CpsUserRepository {
	return &CPSUserStorage{
		dal:               dal.NewMongoDal[imodel.CPSUser, imodel.CPSUser](client, cfg, dbName, collection),
		cpsAction:         dal.NewMongoDal[imodel.CPSAction, imodel.CPSAction](client, cfg, dbName, "cps_actions"),
		client:            client,
		redisRepository:   redisRepository,
		collection:        client.Database(dbName).Collection(collection),
		relatedCollection: relatedCollection,
		logger:            logger,
	}
}

// Implement actual repository methods for CPS action authorization
func (r *CPSUserStorage) Create(ctx context.Context, cpsUser *imodel.CPSUser) error {
	// cpsUserMap := CPSUserMapper(*cpsUser)
	_, err := r.dal.InsertOne(ctx, *cpsUser)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][Create] failed to create CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *CPSUserStorage) Update(ctx context.Context, userCode string, cpsUser *imodel.CPSUser) error {
	r.logger.Infof("[CPSUserStorage][Update] updating CPS user")
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := CPSUserUpdateMapper(cpsUser)

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][Update] failed to update CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	r.logger.Infof("[CPSUserStorage][Update] CPS user updated successfully")
	return nil
}

func (r *CPSUserStorage) Delete(ctx context.Context, userCode string) error {
	r.logger.Infof("[CPSUserStorage][Delete] deleting CPS user")
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"is_deleted": true}

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][Delete] failed to delete CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	r.logger.Infof("[CPSUserStorage][Delete] CPS user deleted successfully")
	return nil
}

func (r *CPSUserStorage) EnableOrDisable(ctx context.Context, userCode string, enable bool) error {
	r.logger.Infof("[CPSUserStorage][EnableOrDisable] processing CPS user enable/disable, enabled: %v", enable)
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified": time.Now()}
	if enable {
		update["login_attempt_count"] = 0
		update["is_first_time_login"] = true
	}
	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][EnableOrDisable] failed to enable/disable CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	// remove user device id from redis
	user, err := r.FindByID(ctx, userCode)
	if err == nil {
		if !enable {
			objID := user.ID.Hex()
			if err := r.redisRepository.Delete(ctx, fmt.Sprintf("%s:%s", constants.RedisCPSUserDeviceIDPrefix, objID)); err != nil {
				r.logger.Errorf("[CPSUserStorage][EnableOrDisable] failed to delete user device id from redis: %v", err)
			}
		}
	}
	r.logger.Infof("[CPSUserStorage][EnableOrDisable] CPS user enable/disable completed successfully")
	return nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByUsername(ctx context.Context, username string) (*imodel.CPSUser, error) {
	r.logger.Infof("[CPSUserStorage][FindByUsername] searching for CPS user by username")
	filter := bson.M{"username": username}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][FindByUsername] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	r.logger.Infof("[CPSUserStorage][FindByUsername] CPS user retrieved successfully")
	return result, nil
}
func (r *CPSUserStorage) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*imodel.CPSUser, error) {
	r.logger.Infof("[CPSUserStorage][FindByPhoneNumber] searching for CPS user by phone number")
	filter := bson.M{"phone_number": phoneNumber}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Infof("[CPSUserStorage][FindByPhoneNumber] CPS user not found")
			return nil, nil
		}
		r.logger.Errorf("[CPSUserStorage][FindByPhoneNumber] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	r.logger.Infof("[CPSUserStorage][FindByPhoneNumber] CPS user retrieved successfully")
	return result, nil
}
func (r *CPSUserStorage) FindByEmail(ctx context.Context, email string) (*imodel.CPSUser, error) {
	r.logger.Infof("[CPSUserStorage][FindByEmail] searching for CPS user by email")
	filter := bson.M{"email": email}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Infof("[CPSUserStorage][FindByEmail] CPS user not found")
			return nil, nil
		}
		r.logger.Errorf("[CPSUserStorage][FindByEmail] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	r.logger.Infof("[CPSUserStorage][FindByEmail] CPS user retrieved successfully")
	return result, nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByID(ctx context.Context, id string) (*imodel.CPSUser, error) {
	r.logger.Infof("[CPSUserStorage][FindByID] fetching CPS user by id")
	filter := bson.M{"user_code": id, "is_deleted": false}

	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][FindByID] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	r.logger.Infof("[CPSUserStorage][FindByID] CPS user retrieved successfully")
	return result, nil
}

func (r *CPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment], error) {
	searchKeys := bson.M{}

	allowedKeys := []string{"enabled", "department", "role"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"email": searchRegex},
			{"username": searchRegex},
			{"user_code": searchRegex},
			{"phone_number": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "date_joined", Value: -1}}}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":           1,
			"user_code":     1,
			"full_name":     1,
			"job_title":     1,
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
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$role_info",
			"preserveNullAndEmptyArrays": true,
		}}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"role": "$role_info.role",
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"role_info": 0,
		}}},

		// Add a $group stage to ensure unique results
		bson.D{{Key: "$group", Value: bson.M{
			"_id": "$user_code",
			"doc": bson.M{"$first": "$$ROOT"},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.M{
			"newRoot": "$doc",
		}}},

		// Use $facet for concurrent data fetching and counting
		bson.D{{Key: "$facet", Value: bson.M{
			"data": []bson.D{
				{{Key: "$sort", Value: bson.D{{Key: "date_joined", Value: -1}}}}, // Ensure sorting within the facet
				{{Key: "$skip", Value: skip}},
				{{Key: "$limit", Value: limit}},
			},
			"total": []bson.D{
				{{Key: "$count", Value: "count"}},
			},
		}}},
	}

	r.logger.Infof("[CPSUserStorage][FindAllWithPagination] fetching CPS users with pagination")
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][FindAllWithPagination] failed to execute aggregation pipeline: %v", err)
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
		r.logger.Errorf("[CPSUserStorage][FindAllWithPagination] failed to decode aggregation results: %v", err)
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
	r.logger.Infof("[CPSUserStorage][FindAllWithPagination] retrieved %d CPS users", len(results[0].Data))
	return &types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment]{
		Data: results[0].Data,
		Meta: meta,
	}, nil
}

func (r *CPSUserStorage) GetPopulatedByID(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error) {
	r.logger.Infof("[CPSUserStorage][GetPopulatedByID] fetching populated CPS user")
	// pipeline := PipelineBuilder(userCode, r.relatedCollection[0], r.relatedCollection[3], r.relatedCollection[1], r.relatedCollection[2])
	pipeline := PipelineBuilder(userCode)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][GetPopulatedByID] failed to aggregate CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		r.logger.Errorf("[CPSUserStorage][GetPopulatedByID] CPS user not found")
		return nil, errors.New(localization.ErrorUserNotFound.Code)
	}

	var resp cpsuser.CpsUserResponse
	if err := cursor.Decode(&resp); err != nil {
		r.logger.Errorf("[CPSUserStorage][GetPopulatedByID] failed to decode CPS user response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	fmt.Println("Decoded CPS User Response:", resp)
	r.logger.Infof("[CPSUserStorage][GetPopulatedByID] populated CPS user retrieved successfully")
	return &resp, nil
}

func (r *CPSUserStorage) GetPopulatedWithRole(ctx context.Context, userCode string) (*cpsuser.CpsUserPopulatedResponse, error) {
	r.logger.Infof("[CPSUserStorage][GetPopulatedWithRole] fetching populated CPS user")
	// relatedCollection: [DepartmentsCollection, PermissionCollection, PermissionCategoryCollection, PermissionGroupsCollection, RolesCollection, JobRolesCollection]
	pipeline := PipelineBuilderWithRole(userCode, r.relatedCollection[0], r.relatedCollection[4], r.relatedCollection[5])

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][GetPopulatedWithRole] failed to aggregate CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		r.logger.Errorf("[CPSUserStorage][GetPopulatedWithRole] CPS user not found")
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}

	var resp cpsuser.CpsUserPopulatedResponse
	if err := cursor.Decode(&resp); err != nil {
		r.logger.Errorf("[CPSUserStorage][GetPopulatedWithRole] failed to decode CPS user response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	r.logger.Infof("[CPSUserStorage][GetPopulatedWithRole] populated CPS user retrieved successfully")
	return &resp, nil
}

// FindByEmailOrPhoneNumberOrUserName implements [storage.CpsUserRepository].
func (r *CPSUserStorage) FindByEmailOrPhoneNumberOrUserName(ctx context.Context, email string, phoneNumber string, username string) (*imodel.CPSUser, error) {
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
		return nil, errors.New("at least one of email, phoneNumber, or username must be provided")
	}
	filter := bson.M{"$or": orFilters}

	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[CPSUserStorage][FindByEmailOrPhoneNumberOrUserName] failed to find CPS user: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}
