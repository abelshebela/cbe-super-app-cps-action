package cps_action

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	cps_action_core "cbe-super-app-cps-action/internal/storage/persistance/cps_action/core"
	"context"
	"errors"
	"fmt"

	local_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSActionStorage struct {
	dal        dal.MongoDal[model.CPSAction, model.CPSAction]
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
}

func NewCPSActionRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.CPSActionRepository {
	return &CPSActionStorage{
		dal:        dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collection),
		client:     client,
		logger:     logger,
		collection: client.Database(dbName).Collection(collection),
	}
}

// Ensure CPSActionRepository implements the storage.CPSActionRepository interface

func (r *CPSActionStorage) Save(ctx context.Context, cpsAction *model.CPSAction) error {
	r.logger.Infof("Attempting to save CPSAction: %+v", cpsAction)

	cps, err := r.dal.InsertOne(ctx, *cpsAction)
	if err != nil {
		r.logger.Errorf("Failed to save CPSAction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	fmt.Printf("Cps file of id is:%s", cps.ID)
	r.logger.Infof("Successfully saved CPSAction with ActionCode: %s", cpsAction.ActionCode)
	return nil
}

func (s *CPSActionStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (types.PaginatedResponse[[]model.CPSAction], error) {
	s.logger.Infof("Finding all CPSActions with pagination. Department: %s, Filter: %+v", department, filterParam)
	filter := bson.M{
		"is_deleted": false,
		"department": department,
	}
	searchKeys := bson.M{}

	allowedKeys := []string{"action_status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["field1"] = searchRegex // choose your searchable field(s)
		searchKeys["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"checker_name": searchRegex},
			{"checker_phone_number": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
			{"action_status": searchRegex},
		}
		s.logger.Infof("Search applied with regex: %v", searchRegex)
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	s.logger.Debugf("Mongo filter: %+v, skip: %d, limit: %d", filter, skip, limit)
	Filter := dal.FilterOp{
		Filter: filter,
		Limit: limit,
		Projection: Projection,
	}

	data, err := s.dal.FindAllWithCursorBasedPagination(ctx, Filter)
	if err != nil {
		s.logger.Errorf("Error fetching paginated CPSActions: %v", err)
		return types.PaginatedResponse[[]model.CPSAction]{}, errors.New(localization.ErrorUnexpectedError.Message)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("Error counting total CPSActions: %v", err)
		return types.PaginatedResponse[[]model.CPSAction]{}, errors.New(localization.ErrorUnexpectedError.Message)
	}

	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	s.logger.Infof("Successfully fetched paginated CPSActions. Total: %d", total)
return types.PaginatedResponse[[]model.CPSAction]{
    Data: data,
    Meta: meta,
}, nil

}

func (r *CPSActionStorage) FindOne(ctx context.Context, filter bson.M) (*model.CPSAction, error) {
	r.logger.Infof("Finding one CPSAction with filter: %+v", filter)
	filterMap := filter

	data, err := r.dal.FindOne(ctx, filterMap, Projection)
	if err != nil {
		r.logger.Errorf("Error finding CPSAction: %v", err)
		code, _ := local_utils.HandleMongoError(err)
		return nil, errors.New(code)
	}
	r.logger.Infof("Successfully found CPSAction: %+v", data)
	return data, nil
}

func (r *CPSActionStorage) Update(ctx context.Context, actionCode string, update model.CPSAction) (*model.CPSAction, error) {
	r.logger.Infof("Updating CPSAction with ActionCode: %s, Update: %+v", actionCode, update)
	filterMap := BuildCPSActionFilter(update)
	updateMap := BuildCPSActionUpdateMap(update)

	data, err := r.dal.UpdateOne(ctx, filterMap, updateMap)
	if err != nil {

		r.logger.Errorf("Error updating CPSAction: %v", err)
		code, _ := local_utils.HandleMongoError(err)
		return nil, errors.New(code)
	}
	r.logger.Infof("Successfully updated CPSAction with ActionCode: %s", actionCode)
	return &data, nil
}

func (r *CPSActionStorage) UpdateCustome(ctx context.Context, filter, update bson.M) error {
	r.logger.Infof("Updating CPSAction with ActionCode: %s, Update: %+v", filter, update)

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("Error updating CPSAction: %v", err)
		code, _ := local_utils.HandleMongoError(err)
		return errors.New(code)
	}
	r.logger.Infof("Successfully updated CPSAction with ")
	return nil
}

func (r *CPSActionStorage) Delete(ctx context.Context, id string) error {
	r.logger.Infof("Deleting CPSAction with ID: %s", id)
	idObj, ok := local_utils.StringToObjectID(id)
	if !ok {
		r.logger.Errorf("Invalid ObjectID for deletion: %s", id)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	filterMap := BuildCPSActionFilter(model.CPSAction{ID: idObj})

	err := r.dal.DeleteOne(ctx, filterMap)
	if err != nil {
		r.logger.Errorf("Error deleting CPSAction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	r.logger.Infof("Successfully deleted CPSAction with ID: %s", id)
	return nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	r.logger.Infof("Finding all CPSActions with pagination. Department: %s, Filter: %+v", department, filterParam)
	// 1. Base filter (only active records)
	baseFilter := bson.M{
		"is_deleted": false,
		"department": department,
	}
	searchKeys := bson.M{}

	allowedKeys := []string{"action_status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["field1"] = searchRegex
		searchKeys["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"checker_name": searchRegex},
			{"checker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"action_type": searchRegex},

		}
		r.logger.Infof("Search applied with regex: %v", searchRegex)
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	filter := dynamicFilter

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value:bson.D{{Key: "created_at", Value: -1}} }},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		cps_action_core.SanitizePipeline(exclude),
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}
	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("Error counting total CPSActions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	r.logger.Infof("Successfully fetched paginated CPSActions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) SanitizedFindOne(ctx context.Context, filter bson.M) (*model.CPSAction, error) {
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		cps_action_core.SanitizePipeline(exclude),
		{{Key: "$limit", Value: 1}},
		{{Key: "$project", Value: Projection}},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return nil, nil
	}

	var result model.CPSAction
	if err := cur.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
