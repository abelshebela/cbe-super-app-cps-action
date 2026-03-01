package cps_action

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	cps_action_core "cbe-super-app-cps-action/internal/storage/persistance/cps_action/core"
	"context"
	"errors"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	local_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CPSActionStorage struct {
	dal        dal.MongoDal[model.CPSAction, model.CPSAction]
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
}

func NewCPSActionRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.CPSActionRepository {
	return &CPSActionStorage{
		dal:        dal.NewMongoDal[model.CPSAction, model.CPSAction](client, cfg, dbName, collection),
		client:     client,
		logger:     logger,
		collection: client.Database(dbName).Collection(collection),
	}
}

// Ensure CPSActionRepository implements the storage.CPSActionRepository interface

func (r *CPSActionStorage) Save(ctx context.Context, cpsAction *model.CPSAction) error {
	r.logger.Infof("[CPSAction][Save] saving CPS action")

	cps, err := r.dal.InsertOne(ctx, *cpsAction)
	if err != nil {
		r.logger.Errorf("[CPSAction][Save] failed to save CPS action: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	r.logger.Infof("[CPSAction][Save] CPS action saved successfully with id: %s", cps.ID.Hex())
	return nil
}

func (s *CPSActionStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (types.PaginatedResponse[[]model.CPSAction], error) {
	s.logger.Infof("[CPSAction][FindAllWithPagination] fetching CPS actions with pagination for department: %s", department)
	filter := bson.M{
		"is_deleted": false,
		"department": department,
	}
	searchKeys := bson.M{}

	allowedKeys := []string{"unique_id", "action_status", "action_code", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id", "created_at"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"unique_id": searchRegex},
			{"current_action": searchRegex},
			{"previous_action": searchRegex},
			{"maker_phone_number": searchRegex},
			{"checker_name": searchRegex},
			{"checker_phone_number": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
			{"action_status": searchRegex},
			{"action_code": searchRegex},
		}
	}

	filter, _, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	Filter := dal.FilterOp{
		Filter:     filter,
		Limit:      limit,
		Projection: Projection,
	}

	data, err := s.dal.FindAllWithCursorBasedPagination(ctx, Filter)
	if err != nil {
		s.logger.Errorf("[CPSAction][FindAllWithPagination] failed to fetch CPS actions: %v", err)
		return types.PaginatedResponse[[]model.CPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[CPSAction][FindAllWithPagination] failed to count CPS actions: %v", err)
		return types.PaginatedResponse[[]model.CPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	s.logger.Infof("[CPSAction][FindAllWithPagination] retrieved %d CPS actions", len(data))
	return types.PaginatedResponse[[]model.CPSAction]{
		Data: data,
		Meta: meta,
	}, nil

}

func (r *CPSActionStorage) FindOne(ctx context.Context, filter bson.M) (*model.CPSAction, error) {
	r.logger.Infof("[CPSAction][FindOne] fetching CPS action")
	filterMap := filter

	data, err := r.dal.FindOne(ctx, filterMap, Projection)
	if err != nil {
		r.logger.Errorf("[CPSAction][FindOne] failed to find CPS action: %v", err)
		return nil, local_utils.HandleDBError(err)
	}
	r.logger.Infof("[CPSAction][FindOne] CPS action retrieved successfully")
	return data, nil
}

func (r *CPSActionStorage) Update(ctx context.Context, actionCode string, update model.CPSAction) (*model.CPSAction, error) {
	r.logger.Infof("[CPSAction][Update] updating CPS action for action code: %s", actionCode)
	filterMap := BuildCPSActionFilter(update)
	updateMap := BuildCPSActionUpdateMap(update)
	// filterMap := bson.M{}
	data, err := r.dal.UpdateOne(ctx, filterMap, updateMap)
	if err != nil {
		r.logger.Errorf("[CPSAction][Update] failed to update CPS action: %v", err)
		return nil, local_utils.HandleDBError(err)
	}
	r.logger.Infof("[CPSAction][Update] CPS action updated successfully")
	return &data, nil
}

func (r *CPSActionStorage) UpdateByActionCode(ctx context.Context, actionCode string, update model.CPSAction) (*model.CPSAction, error) {
	r.logger.Infof("[CPSAction][UpdateByActionCode] updating CPS action for action code: %s", actionCode)
	updateMap := BuildCPSActionUpdateMap(update)
	filterMap := bson.M{"action_code": actionCode}
	data, err := r.dal.UpdateOne(ctx, filterMap, updateMap)
	if err != nil {
		r.logger.Errorf("[CPSAction][UpdateByActionCode] failed to update CPS action: %v", err)
		return nil, local_utils.HandleDBError(err)
	}
	r.logger.Infof("[CPSAction][UpdateByActionCode] CPS action updated successfully")
	return &data, nil
}

func (r *CPSActionStorage) UpdateCustome(ctx context.Context, filter, update bson.M) error {
	r.logger.Infof("[CPSAction][UpdateCustome] updating CPS action")

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[CPSAction][UpdateCustome] failed to update CPS action: %v", err)
		return local_utils.HandleDBError(err)
	}
	r.logger.Infof("[CPSAction][UpdateCustome] CPS action updated successfully")
	return nil
}

func (r *CPSActionStorage) Delete(ctx context.Context, id string) error {
	r.logger.Infof("[CPSAction][Delete] deleting CPS action with ID: %s", id)
	idObj, ok := local_utils.StringToObjectID(id)
	if !ok {
		r.logger.Errorf("[CPSAction][Delete] invalid object id: %s", id)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filterMap := BuildCPSActionFilter(model.CPSAction{ID: idObj})

	err := r.dal.DeleteOne(ctx, filterMap)
	if err != nil {
		r.logger.Errorf("[CPSAction][Delete] failed to delete CPS action: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	r.logger.Infof("[CPSAction][Delete] CPS action deleted successfully with ID: %s", id)
	return nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	r.logger.Infof("[CPSAction][SanitizedFindAllWithPagination] finding all CPS actions. Department: %s, Filter: %+v", department, filterParam)
	// 1. Base filter (only active records)
	baseFilter := bson.M{
		"is_deleted": false,
		// "department": department,
	}
	searchKeys := bson.M{}

	allowedKeys := []string{"action_status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	delete(dynamicFilter, "created_at")
	filter := dynamicFilter

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		cps_action_core.SanitizePipeline(exclude),
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPagination] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPagination] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPagination] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	r.logger.Infof("[CPSAction][SanitizedFindAllWithPagination] successfully fetched paginated CPS actions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPaginationForApprover(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	r.logger.Infof("[CPSAction][SanitizedFindAllWithPaginationForApprover] finding all CPS actions. RAList: %s, Filter: %+v", RAList, filterParam)
	// 1. Base filter (only active records)
	baseFilter := bson.M{
		"is_deleted": false,
	}
	searchKeys := bson.M{}

	allowedKeys := []string{"action_status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	delete(dynamicFilter, "created_at")
	filter := dynamicFilter
	filter["request_action"] = bson.M{"$in": RAList}

	userFilter := bson.M{
		"$or": []bson.M{
			{"maker_id": userID},
			{"checker_users.checker_id": userID},
		},
	}

	var finalMatch bson.M
	if filterParam.Filters["action_status"] == "PENDING" {
		finalMatch = filter
	} else {
		finalMatch = bson.M{
			"$and": []bson.M{
				filter,
				userFilter,
			},
		}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: finalMatch}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		cps_action_core.SanitizePipeline(exclude),
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationForApprover] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationForApprover] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := r.dal.TotalCount(ctx, finalMatch)
	if err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationForApprover] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	r.logger.Infof("[CPSAction][SanitizedFindAllWithPaginationForApprover] successfully fetched paginated CPS actions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPaginationForAuditor(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	r.logger.Infof("[CPSAction][SanitizedFindAllWithPaginationForAuditor] finding all CPS actions. RAList: %s, Filter: %+v", RAList, filterParam)
	inboxOr := interface{}(nil)
	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["__inbox_or"]; ok {
			inboxOr = v
			delete(filterParam.Filters, "__inbox_or")
		}
	}
	// 1. Base filter (only active records)
	baseFilter := bson.M{
		"is_deleted": false,
	}
	searchKeys := bson.M{}

	allowedKeys := []string{"action_code", "action_status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id", "auditor_status", "auditor_users.auditor_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"action_code": searchRegex},
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"auditor_status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
		}

	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	delete(dynamicFilter, "created_at")
	filter := dynamicFilter

	if inboxOr == nil {
		if RAList != nil {
			RAList = local_utils.RemoveDuplicates(RAList)
		}
		if RAList != nil && len(RAList) > 0 {
			filter["request_action"] = bson.M{"$in": RAList}
		}
	} else {
		var ors []bson.M
		switch t := inboxOr.(type) {
		case []bson.M:
			ors = t
		case []interface{}:
			ors = make([]bson.M, 0, len(t))
			for _, it := range t {
				if m, ok := it.(bson.M); ok {
					ors = append(ors, m)
				}
			}
		}

		and := []bson.M{}
		for k, v := range baseFilter {
			if k == "$or" {
				continue
			}
			and = append(and, bson.M{k: v})
		}
		if baseFilter["$or"] != nil {
			and = append(and, bson.M{"$or": baseFilter["$or"]})
		}
		if len(filter) > 0 {
			and = append(and, filter)
		}
		if len(ors) > 0 {
			and = append(and, bson.M{"$or": ors})
		}
		filter = bson.M{"$and": and}
	}

	if filter["action_status"] == "" || filter["action_status"] == constants.Pending {
		filter["action_status"] = bson.M{"$in": []string{string(constants.Approved), string(constants.Rejected)}}
	}
	// filter["action_status"] = bson.M{"$in": []string{string(constants.Approved), string(constants.Rejected)}}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		cps_action_core.SanitizePipeline(exclude),
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationForAuditor] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationForAuditor] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationForAuditor] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	r.logger.Infof("[CPSAction][SanitizedFindAllWithPaginationForAuditor] successfully fetched paginated CPS actions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPaginationCPSActions(ctx context.Context, userID, role string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	r.logger.Infof("[CPSAction][SanitizedFindAllWithPaginationCPSActions] finding all CPS actions for User: %s, RAList: %s, Filter: %+v", userID, RAList, filterParam)

	if strings.TrimSpace(userID) == "" {
		meta := local_utils.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage)
		return &types.PaginatedResponse[[]*model.CPSAction]{
			Data: []*model.CPSAction{},
			Meta: meta,
		}, nil
	}

	baseFilter := bson.M{
		"is_deleted": false,
	}

	searchKeys := bson.M{}

	// Exclude maker_id and checker_id from allowedKeys so request cannot override userFilter
	allowedKeys := []string{"action_status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "checker_name", "checker_phone_number", "auditor_status"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"auditor_status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	delete(dynamicFilter, "created_at")
	filter := dynamicFilter

	if role != "maker" {
		if RAList != nil {
			RAList = local_utils.RemoveDuplicates(RAList)
		} else {
			RAList = []string{}
		}
		filter["request_action"] = bson.M{"$in": RAList}
	}

	var userFilter bson.M
	if role == "maker" {
		userFilter = bson.M{"maker_id": userID}
	}
	if role == "checker" {
		userFilter = bson.M{
			"$or": []bson.M{
				{"maker_id": userID},
				{"checker_users.checker_id": userID},
			},
		}
	}
	// if role == "auditor" {
	// 	userFilter = bson.M{"auditor_users.auditor_id": userID}
	// }

	var finalMatch bson.M
	if role == "checker" && filterParam.Filters["action_status"] == "PENDING" || role == "auditor" {
		finalMatch = filter
	} else {
		finalMatch = bson.M{
			"$and": []bson.M{
				filter,
				userFilter,
			},
		}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: finalMatch}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		cps_action_core.SanitizePipeline(exclude),
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationCPSActions] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationCPSActions] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := r.dal.TotalCount(ctx, finalMatch)
	if err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindAllWithPaginationCPSActions] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	r.logger.Infof("[CPSAction][SanitizedFindAllWithPaginationCPSActions] successfully fetched paginated CPS actions. Total: %d", total)
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
		r.logger.Errorf("[CPSAction][SanitizedFindOne] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}
	defer func() {
		_ = cur.Close(ctx)
	}()

	// Handle empty cursor
	if cur == nil || !cur.Next(ctx) {
		r.logger.Warnf("[CPSAction][SanitizedFindOne] no CPS action found for filter")
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	var result model.CPSAction
	if err := cur.Decode(&result); err != nil {
		r.logger.Errorf("[CPSAction][SanitizedFindOne] failed to decode document: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &result, nil
}

func (r *CPSActionStorage) GetCountByDepartment(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error) {
	pipeline := mongo.Pipeline{
		// match stage
		{{Key: "$match", Value: bson.M{
			"is_deleted": false,
			"department": department,
		}}},
		// group stage
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "pending", Value: bson.D{
				{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{bson.D{{Key: "$eq", Value: bson.A{"$action_status", "PENDING"}}}, 1, 0}}}},
			}},
			{Key: "approved", Value: bson.D{
				{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{bson.D{{Key: "$eq", Value: bson.A{"$action_status", "APPROVED"}}}, 1, 0}}}},
			}},
			{Key: "rejected", Value: bson.D{
				{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{bson.D{{Key: "$eq", Value: bson.A{"$action_status", "REJECTED"}}}, 1, 0}}}},
			}},
			{Key: "canceled", Value: bson.D{
				{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{bson.D{{Key: "$eq", Value: bson.A{"$action_status", "CANCELED"}}}, 1, 0}}}},
			}},
		}}},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[CPSAction][GetCountByDepartment] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() {
		_ = cur.Close(ctx)
	}()

	var result actionDto.CPSActionCountResponse
	if cur.Next(ctx) {
		if err := cur.Decode(&result); err != nil {
			r.logger.Errorf("[CPSAction][GetCountByDepartment] failed to decode document: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
	} else {
		// If no documents match, return zero counts instead of error
		return &actionDto.CPSActionCountResponse{
			Pending:  0,
			Approved: 0,
			Rejected: 0,
			Canceled: 0,
		}, nil
	}

	return &result, nil
}

// memory intensice hold all in memory at one good for small  amount of data  memory user O(n)
// func (r *CPSActionStorage) FindByDateRange(ctx context.Context,start_date, end_date time.Time) ([]*model.CPSAction, error) {
//     filter := bson.M{
//         "created_at": bson.M{
//             "$gte": start_date,
//             "$lte": end_date,
//         },
//     }

//     // opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})

//     // cursor, err := r.collection.Find(ctx, filter, opts)
// 	actions,err := r.dal.FindAll(ctx,filter,nil)
//     if err != nil {
//         return nil, err
//     }
//     // defer cursor.Close(ctx)

//     // var actions []*model.CPSAction
//     // if err := cursor.All(ctx, &actions); err != nil {
//     //     return nil, err
//     // }

//     return actions, nil

// }

// stream process constant memrory useage O(1)
func (r *CPSActionStorage) StreamByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	handler func(*model.CPSAction) error,
) error {

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetBatchSize(1000) // very important

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var action model.CPSAction
		if err := cursor.Decode(&action); err != nil {
			return err
		}

		if err := handler(&action); err != nil {
			return err
		}
	}

	return cursor.Err()
}
