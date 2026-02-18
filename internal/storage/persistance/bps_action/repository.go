package bps_action

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	bps_action_core "cbe-super-app-cps-action/internal/storage/persistance/cps_action/core"
	"context"
	"errors"
	"strings"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	bpsActionDto "cbe-super-app-cps-action/internal/constants/dto/bps_action"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	bps_action "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Projection = bson.M{
	"action_code":           1,
	"action_name":           1,
	"maker_id":              1,
	"maker_name":            1,
	"maker_phone_number":    1,
	"checker_users":         1,
	"checker_count":         1,
	"current_checker_index": 1,
	"action_description":    1,
	"action_type":           1,
	"status":                1,
	"auditor_status":        1,
	"auditor_users":         1,
	"auditor_count":         1,
	"current_auditor_index": 1,
	"request_action":        1,
	"action_created_at":     1,
	"maker_action_time":     1,
	"checker_action_time":   1,
	"previous_action":       1,
	"current_action":        1,
	"department":            1,
	"rejection_reason":      1,
	"unique_id":             1,
	"created_at":            1,
	"last_modified_at":      1,
	"reversed_by_role_id":   1,
	"reversed_by_id":        1,
	"reversed_by_name":      1,
	"reversed_at":           1,
}

type bpsActionRepository struct {
	client     *mongo.Client
	actionDal  dal.MongoDal[bps_action.BPSAction, bps_action.BPSAction]
	collection *mongo.Collection
	logger     utils.Logger
}

func NewBPSActionRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger, cfg *config.VaultConfig) storage.BPSActionRepository {
	return &bpsActionRepository{
		client:     client,
		actionDal:  dal.NewMongoDal[bps_action.BPSAction, bps_action.BPSAction](client, cfg, dbName, collection),
		collection: client.Database(dbName).Collection(collection),
		logger:     logger,
	}
}

// GetBPSActionByUserID implements [storage.BPSActionRepository].
func (b *bpsActionRepository) GetBPSActionByUserID(ctx context.Context, userID string, filterParam types.Filter) (types.PaginatedResponse[[]bps_action.BPSAction], error) {
	b.logger.Infof("[GetBPSActionByUserID] fetching BPS actions for user ID: %s", userID)
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		b.logger.Errorf("[GetBPSActionByUserID] invalid user ID: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorInvalidID.Code)
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, nil)
	filter["user_information.user_id"] = objID

	cus, err := b.actionDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		b.logger.Errorf("[GetBPSActionByUserID] failed to fetch BPS actions: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter = bson.M{"user_information.user_id": objID}
	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[GetBPSActionByUserID] failed to fetch total count: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return types.PaginatedResponse[[]bps_action.BPSAction]{
		Data: cus,
		Meta: meta,
	}, nil
}

func (b *bpsActionRepository) Save(ctx context.Context, action *bps_action.BPSAction) error {
	b.logger.Infof("[Save] saving BPS action")

	_, err := b.actionDal.InsertOne(ctx, *action)
	if err != nil {
		b.logger.Errorf("[Save] failed to save BPS action: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (b *bpsActionRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (types.PaginatedResponse[[]bps_action.BPSAction], error) {
	b.logger.Infof("[FindAllWithPagination] fetching BPS actions with pagination for department: %s", department)

	baseFilter := bson.M{
		"is_deleted": false,
		"department": department,
	}
	searchKeys := bson.M{}
	allowedKeys := []string{"unique_id", "status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"unique_id": searchRegex},
			{"maker_phone_number": searchRegex},
			{"checker_name": searchRegex},
			{"checker_phone_number": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
			{"status": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	for k, v := range baseFilter {
		filter[k] = v
	}

	data, err := b.actionDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		b.logger.Errorf("[FindAllWithPagination] failed to fetch BPS actions: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[FindAllWithPagination] failed to count BPS actions: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return types.PaginatedResponse[[]bps_action.BPSAction]{
		Data: data,
		Meta: meta,
	}, nil
}

func (b *bpsActionRepository) FindOne(ctx context.Context, filter bson.M) (*bps_action.BPSAction, error) {
	b.logger.Infof("[FindOne] fetching BPS action")

	data, err := b.actionDal.FindOne(ctx, filter, Projection)
	if err != nil {
		b.logger.Errorf("[FindOne] failed to find BPS action: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return data, nil
}

func (b *bpsActionRepository) SanitizedFindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*bps_action.BPSAction], error) {
	baseFilter := bson.M{
		"is_deleted": false,
	}
	if strings.TrimSpace(department) != "" {
		baseFilter["department"] = department
	}

	allowedKeys := []string{"status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id", "auditor_status"}
	searchKeys := bson.M{}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"status": searchRegex},
			{"auditor_status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	delete(dynamicFilter, "created_at")
	filter := dynamicFilter

	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		bps_action_core.SanitizePipeline(exclude),
	}

	cur, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPagination] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []*bps_action.BPSAction
	if err := cur.All(ctx, &results); err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPagination] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPagination] failed to count BPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*bps_action.BPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (b *bpsActionRepository) SanitizedFindAllWithPaginationForApprover(
	ctx context.Context,
	userID string,
	filterParam types.Filter,
	RAList []string,
) (*types.PaginatedResponse[[]*bps_action.BPSAction], error) {

	if strings.TrimSpace(userID) == "" {
		meta := local_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage)
		return &types.PaginatedResponse[[]*bps_action.BPSAction]{
			Data: []*bps_action.BPSAction{},
			Meta: meta,
		}, nil
	}

	filter := bson.M{}

	if len(RAList) > 0 {
		filter["request_action"] = bson.M{"$in": RAList}
	}

	if strings.TrimSpace(filterParam.Search) != "" {
		regex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"maker_name": regex},
			{"maker_phone_number": regex},
			{"status": regex},
			{"auditor_status": regex},
			{"action_type": regex},
			{"request_action": regex},
		}
	}

	allowedKeys := []string{
		"status",
		"action_type",
		"request_action",
		"maker_phone_number",
		"checker_phone_number",
		"maker_name",
		"checker_name",
		"auditor_status",
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, nil, allowedKeys)

	for k, v := range dynamicFilter {
		if k != "created_at" { // avoid broken range leftovers
			filter[k] = v
		}
	}

	exclude := []string{
		"password",
		"first_password_set",
		"login_attempt_count",
		"is_deleted",
		"otp_verfy_count",
		"otp_last_tried_at",
		"otp_last_verified_at",
		"permission_group",
		"permissions",
		"last_login_attempt",
		"next_login_attempt",
		"is_first_time_login",
		"last_login",
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		bps_action_core.SanitizePipeline(exclude),
	}

	cur, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationForApprover] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*bps_action.BPSAction
	if err := cur.All(ctx, &results); err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationForApprover] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationForApprover] failed to count BPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*bps_action.BPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (b *bpsActionRepository) SanitizedFindAllWithPaginationForAuditor(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*bps_action.BPSAction], error) {
	if strings.TrimSpace(userID) == "" {
		meta := local_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage)
		return &types.PaginatedResponse[[]*bps_action.BPSAction]{Data: []*bps_action.BPSAction{}, Meta: meta}, nil
	}

	baseFilter := bson.M{}
	allowedKeys := []string{"status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "checker_name", "checker_phone_number", "checker_id", "auditor_status", "home_branch", "branch_code"}
	searchKeys := bson.M{}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
			{"home_branch": searchRegex},
			{"branch_code": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	delete(dynamicFilter, "created_at")
	filter := dynamicFilter
	filter["request_action"] = bson.M{"$in": RAList}

	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		bps_action_core.SanitizePipeline(exclude),
	}

	cur, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationForAuditor] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []*bps_action.BPSAction
	if err := cur.All(ctx, &results); err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationForAuditor] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationForAuditor] failed to count BPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*bps_action.BPSAction]{Data: results, Meta: meta}, nil
}

func (b *bpsActionRepository) SanitizedFindAllWithPaginationBPSActions(ctx context.Context, userID, role string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*bps_action.BPSAction], error) {
	if strings.TrimSpace(userID) == "" {
		meta := local_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage)
		return &types.PaginatedResponse[[]*bps_action.BPSAction]{Data: []*bps_action.BPSAction{}, Meta: meta}, nil
	}

	baseFilter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}
	allowedKeys := []string{"status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "checker_name", "checker_phone_number", "auditor_status"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"status": searchRegex},
			{"auditor_status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	delete(dynamicFilter, "created_at")
	filter := dynamicFilter
	filter["request_action"] = bson.M{"$in": RAList}

	var userFilter bson.M
	if role == "maker" {
		userFilter = bson.M{"maker_id": userID}
	}
	if role == "auditor" {
		userFilter = bson.M{"auditor_users.auditor_id": userID}
	}

	var finalMatch bson.M
	if role == "checker" && filterParam.Filters != nil && filterParam.Filters["status"] == "PENDING" {
		finalMatch = filter
	} else if userFilter != nil {
		finalMatch = bson.M{"$and": []bson.M{filter, userFilter}}
	} else {
		finalMatch = filter
	}

	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: finalMatch}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		bps_action_core.SanitizePipeline(exclude),
	}

	cur, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationBPSActions] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []*bps_action.BPSAction
	if err := cur.All(ctx, &results); err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationBPSActions] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.actionDal.TotalCount(ctx, finalMatch)
	if err != nil {
		b.logger.Errorf("[SanitizedFindAllWithPaginationBPSActions] failed to count BPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*bps_action.BPSAction]{Data: results, Meta: meta}, nil
}

func (b *bpsActionRepository) SanitizedFindOne(ctx context.Context, filter bson.M) (*bps_action.BPSAction, error) {
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		bps_action_core.SanitizePipeline(exclude),
		{{Key: "$limit", Value: 1}},
		{{Key: "$project", Value: Projection}},
	}

	cur, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		b.logger.Errorf("[SanitizedFindOne] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	if cur == nil || !cur.Next(ctx) {
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	var result bps_action.BPSAction
	if err := cur.Decode(&result); err != nil {
		b.logger.Errorf("[SanitizedFindOne] failed to decode document: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &result, nil
}

func (b *bpsActionRepository) Update(ctx context.Context, actionCode string, update bps_action.BPSAction) (*bps_action.BPSAction, error) {
	filter := bson.M{"action_code": actionCode}
	updateMap := bson.M{"$set": update}
	data, err := b.actionDal.UpdateOne(ctx, filter, updateMap)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return &data, nil
}

func (b *bpsActionRepository) UpdateByActionCode(ctx context.Context, actionCode string, update bps_action.BPSAction) (*bps_action.BPSAction, error) {
	filter := bson.M{"action_code": actionCode}
	updateMap := bson.M{"$set": update}
	data, err := b.actionDal.UpdateOne(ctx, filter, updateMap)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return &data, nil
}

func (b *bpsActionRepository) UpdateCustome(ctx context.Context, filter, update bson.M) error {
	_, err := b.actionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	return nil
}

func (b *bpsActionRepository) Delete(ctx context.Context, id string) error {
	objID, ok := local_util.StringToObjectID(id)
	if !ok {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if err := b.actionDal.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *bpsActionRepository) GetCountByDepartment(ctx context.Context, department string) (*bpsActionDto.BPSActionCountResponse, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"is_deleted": false, "department": department}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "approved", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{bson.D{{Key: "$eq", Value: bson.A{"$status", "APPROVED"}}}, 1, 0}}}}}},
			{Key: "rejected", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{bson.D{{Key: "$eq", Value: bson.A{"$status", "REJECTED"}}}, 1, 0}}}}}},
			{Key: "inprogress", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{bson.D{{Key: "$eq", Value: bson.A{"$auditor_status", "INPROGRESS"}}}, 1, 0}}}}}},
			{Key: "completed", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{bson.D{{Key: "$eq", Value: bson.A{"$auditor_status", "AUDITORNOTCHECKED"}}}, 1, 0}}}}}},
		}}},
	}

	cur, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		b.logger.Errorf("[GetCountByDepartment] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var result bpsActionDto.BPSActionCountResponse
	if cur.Next(ctx) {
		if err := cur.Decode(&result); err != nil {
			b.logger.Errorf("[GetCountByDepartment] failed to decode document: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return &result, nil
	}

	return &bpsActionDto.BPSActionCountResponse{Approved: 0, Rejected: 0, Inprogress: 0, Completed: 0}, nil
}
