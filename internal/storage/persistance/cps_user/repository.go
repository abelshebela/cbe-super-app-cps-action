package cps_user

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSUserStorage struct {
	dal       dal.MongoDal[model.CPSUser, model.CPSUser]
	cpsAction dal.MongoDal[model.CPSAction, model.CPSAction]
	client    *mongo.Client
	logger    utils.Logger
}

func NewCPSUserRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.CpsUserRepository {
	return &CPSUserStorage{
		dal:       dal.NewMongoDal[model.CPSUser, model.CPSUser](client, dbName, collection),
		cpsAction: dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, "cps_actions"),
		client:    client,
		logger:    logger,
	}
}

// Implement actual repository methods for CPS action authorization
func (r *CPSUserStorage) Create(ctx context.Context, cpsUser *model.CPSUser) error {
	_, err := r.dal.InsertOne(ctx, *cpsUser)
	if err != nil {
		r.logger.Errorf("failed to create CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

func (r *CPSUserStorage) Update(ctx context.Context, userCode string, cpsUser *model.CPSUser) error {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := CPSUserUpdateMapper(cpsUser)

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to update CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

func (r *CPSUserStorage) Delete(ctx context.Context, userCode string) error {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"is_deleted": true}

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to delete CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

func (r *CPSUserStorage) EnableOrDisable(ctx context.Context, userCode string, enable bool) error {
	filter := bson.M{"user_code": userCode, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified": time.Now()}

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to enable/disable CPS user: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}

// FindByID supports both ObjectID and user_code lookups
func (r *CPSUserStorage) FindByID(ctx context.Context, id string) (*model.CPSUser, error) {
	var filter bson.M
	if objID, ok := local_util.StringToObjectID(id); ok {
		filter = bson.M{"_id": objID, "is_deleted": false}
	} else {
		filter = bson.M{"user_code": id, "is_deleted": false}
	}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	return result, nil
}

func (r *CPSUserStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.CPSUser], error) {
	filter := bson.M{}
	searchKeys := bson.M{}

	allowedKeys := []string{"enabled", "department", "role"}
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
	filter["is_deleted"] = false
	fmt.Println(filter)
	fmt.Println("999999999999999999999999999999999999-----------------")

	data, err := r.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*model.CPSUser]{
		Data: data,
		Meta: meta,
	}, nil
}
