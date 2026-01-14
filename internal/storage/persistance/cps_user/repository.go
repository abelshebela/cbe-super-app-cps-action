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
		r.logger.Errorf("failed to create CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

func (r *CPSUserStorage) Update(ctx context.Context, userCode string, cpsUser *imodel.CPSUser) error {
	r.logger.Infof("[Update] updating CPS user")
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := CPSUserUpdateMapper(cpsUser)

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[Update] failed to update CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	r.logger.Infof("[Update] CPS user updated successfully")
	return nil
}

func (r *CPSUserStorage) Delete(ctx context.Context, userCode string) error {
	r.logger.Infof("[Delete] deleting CPS user")
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"is_deleted": true}

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[Delete] failed to delete CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	r.logger.Infof("[Delete] CPS user deleted successfully")
	return nil
}

func (r *CPSUserStorage) EnableOrDisable(ctx context.Context, userCode string, enable bool) error {
	r.logger.Infof("[EnableOrDisable] processing CPS user enable/disable, enabled: %v", enable)
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified": time.Now()}
	if enable {
		update["login_attempt_count"] = 0
		update["is_first_time_login"] = true
	}
	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[EnableOrDisable] failed to enable/disable CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	// remove user device id from redis
	user, err := r.FindByID(ctx, userCode)
	if err == nil {
		if !enable {
			objID := user.ID.Hex()
			if err := r.redisRepository.Delete(ctx, fmt.Sprintf("%s:%s", constants.RedisCPSUserDeviceIDPrefix, objID)); err != nil {
				r.logger.Errorf("[EnableOrDisable] failed to delete user device id from redis: %v", err)
			}
		}
	}
	r.logger.Infof("[EnableOrDisable] CPS user enable/disable completed successfully")
	return nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByUsername(ctx context.Context, username string) (*imodel.CPSUser, error) {
	r.logger.Infof("[FindByUsername] searching for CPS user by username")
	filter := bson.M{"username": username}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Infof("[FindByUsername] CPS user not found")
			return nil, nil // Return nil, nil when no document found (not an error)
		}
		r.logger.Errorf("[FindByUsername] failed to find CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	r.logger.Infof("[FindByUsername] CPS user retrieved successfully")
	return result, nil
}
func (r *CPSUserStorage) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*imodel.CPSUser, error) {
	r.logger.Infof("[FindByPhoneNumber] searching for CPS user by phone number")
	filter := bson.M{"phone_number": phoneNumber}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Infof("[FindByPhoneNumber] CPS user not found")
			return nil, nil
		}
		r.logger.Errorf("[FindByPhoneNumber] failed to find CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	r.logger.Infof("[FindByPhoneNumber] CPS user retrieved successfully")
	return result, nil
}
func (r *CPSUserStorage) FindByEmail(ctx context.Context, email string) (*imodel.CPSUser, error) {
	r.logger.Infof("[FindByEmail] searching for CPS user by email")
	filter := bson.M{"email": email}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Infof("[FindByEmail] CPS user not found")
			return nil, nil
		}
		r.logger.Errorf("[FindByEmail] failed to find CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	r.logger.Infof("[FindByEmail] CPS user retrieved successfully")
	return result, nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByID(ctx context.Context, id string) (*imodel.CPSUser, error) {
	r.logger.Infof("[FindByID] fetching CPS user by id")
	filter := bson.M{"user_code": id, "is_deleted": false}

	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("[FindByID] CPS user not found")
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		r.logger.Errorf("[FindByID] failed to find CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	r.logger.Infof("[FindByID] CPS user retrieved successfully")
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
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
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
			"from":         "department",
			"localField":   "department",
			"foreignField": "_id",
			"as":           "department_info",
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$project", Value: bson.M{
					"_id":        1,
					"department": 1,
				}}},
			},
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$department_info",
			"preserveNullAndEmptyArrays": true,
		}}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"department": bson.M{
				"id":   "$department_info._id",
				"name": "$department_info.department",
			},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"department_info": 0,
		}}},

		// Use $facet for concurrent data fetching and counting
		bson.D{{Key: "$facet", Value: bson.M{
			"data": []bson.D{
				{{Key: "$skip", Value: skip}},
				{{Key: "$limit", Value: limit}},
			},
			"total": []bson.D{
				{{Key: "$count", Value: "count"}},
			},
		}}},
	}

	r.logger.Infof("[FindAllWithPagination] fetching CPS users with pagination")
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[FindAllWithPagination] failed to execute aggregation pipeline: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data  []*cpsuser.CPSUserWithDepartment `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("[FindAllWithPagination] failed to decode aggregation results: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
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
	r.logger.Infof("[FindAllWithPagination] retrieved %d CPS users", len(results[0].Data))
	return &types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment]{
		Data: results[0].Data,
		Meta: meta,
	}, nil
}

func (r *CPSUserStorage) GetPopulatedByID(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error) {
	r.logger.Infof("[GetPopulatedByID] fetching populated CPS user")
	// pipeline := PipelineBuilder(userCode, r.relatedCollection[0], r.relatedCollection[3], r.relatedCollection[1], r.relatedCollection[2])
	pipeline := PipelineBuilder(userCode)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[GetPopulatedByID] failed to aggregate CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		r.logger.Errorf("[GetPopulatedByID] CPS user not found")
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}

	var resp cpsuser.CpsUserResponse
	if err := cursor.Decode(&resp); err != nil {
		r.logger.Errorf("[GetPopulatedByID] failed to decode CPS user response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	fmt.Println("Decoded CPS User Response:", resp)
	r.logger.Infof("[GetPopulatedByID] populated CPS user retrieved successfully")
	return &resp, nil
}

func (r *CPSUserStorage) GetPopulatedWithRole(ctx context.Context, userCode string) (*cpsuser.CpsUserPopulatedResponse, error) {
	r.logger.Infof("[GetPopulatedByID] fetching populated CPS user")
	// relatedCollection: [DepartmentsCollection, PermissionCollection, PermissionCategoryCollection, PermissionGroupsCollection, RolesCollection, JobRolesCollection]
	pipeline := PipelineBuilderWithRole(userCode, r.relatedCollection[0], r.relatedCollection[4], r.relatedCollection[5])

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[GetPopulatedByID] failed to aggregate CPS user: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		r.logger.Errorf("[GetPopulatedByID] CPS user not found")
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}

	var resp cpsuser.CpsUserPopulatedResponse
	if err := cursor.Decode(&resp); err != nil {
		r.logger.Errorf("[GetPopulatedByID] failed to decode CPS user response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	fmt.Println("Decoded CPS User Response:", resp)
	r.logger.Infof("[GetPopulatedByID] populated CPS user retrieved successfully")
	return &resp, nil
}
