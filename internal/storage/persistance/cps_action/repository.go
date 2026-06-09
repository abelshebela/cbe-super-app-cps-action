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

	local_util "cbe-super-app-cps-action/pkgs/utils"

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

func (r *CPSActionStorage) Save(ctx context.Context, cpsAction *model.CPSAction) (model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][Save] saving CPS action")

	cps, err := r.dal.InsertOne(ctx, *cpsAction)
	if err != nil {
		log.Errorf("[CPSAction][Save] failed to save CPS action: %v", err)
		return model.CPSAction{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsAction.ActionCode = cps.ActionCode
	if md := types.GetMetadata(ctx); md != nil {
		md.CPSActionCode = cps.ActionCode
	}
	log.Infof("[CPSAction][Save] CPS action saved successfully with code: %s", cps.ActionCode)
	return cps, nil
}

func (s *CPSActionStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (types.PaginatedResponse[[]model.CPSAction], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	log.Infof("[CPSAction][FindAllWithPagination] fetching CPS actions with pagination for department: %s", department)

	searchKeys := bson.M{}

	allowedKeys := []string{"unique_id", "action_status", "action_code", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_id", "created_at"}

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

	// Always enforce base constraints — FilterBuilder returns a fresh bson.M
	// and does not carry over is_deleted or department.
	filter["is_deleted"] = false
	filter["department"] = department

	applyActionStatusFilter(filterParam.Filters, filter)
	applyActionCodeFilter(filterParam.Filters, filter)
	Filter := dal.FilterOp{
		Filter:     filter,
		Limit:      limit,
		Projection: Projection,
	}

	log.Infof("Update cps action filter: %v", filter)
	data, err := s.dal.FindAllWithCursorBasedPagination(ctx, Filter)
	if err != nil {
		log.Errorf("[CPSAction][FindAllWithPagination] failed to fetch CPS actions: %v", err)
		return types.PaginatedResponse[[]model.CPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[CPSAction][FindAllWithPagination] failed to count CPS actions: %v", err)
		return types.PaginatedResponse[[]model.CPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	log.Infof("[CPSAction][FindAllWithPagination] retrieved %d CPS actions", len(data))
	return types.PaginatedResponse[[]model.CPSAction]{
		Data: data,
		Meta: meta,
	}, nil

}

func (r *CPSActionStorage) FindOne(ctx context.Context, filter bson.M) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][FindOne] fetching CPS action")
	filterMap := filter

	data, err := r.dal.FindOne(ctx, filterMap, Projection)
	if err != nil {
		log.Errorf("[CPSAction][FindOne] failed to find CPS action: %v", err)
		return nil, local_utils.HandleDBError(err)
	}
	log.Infof("[CPSAction][FindOne] CPS action retrieved successfully")
	return data, nil
}

func (r *CPSActionStorage) Update(ctx context.Context, actionCode string, update model.CPSAction, Group string, RequestActionGroups map[string][]constants.RequestAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][Update] updating CPS action for action code: %s", actionCode)
	filterMap := BuildCPSActionFilter(update)
	updateMap := BuildCPSActionUpdateMap(update)
	modelData := []mongo.WriteModel{}

	if strings.Contains(update.RequestAction, constants.DELETE) || strings.Contains(update.RequestAction, constants.DISABLE) {
		RAUpdateList := local_utils.GetRAListForUpdateAction(constants.RequestAction(update.RequestAction), Group, RequestActionGroups)

		log.Infof("[CPSActionRepo][Update] list of request action get RAUpdateList: %v", RAUpdateList)
		modelData = append(modelData, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"unique_id": update.UniqueId, "is_deleted": false, "request_action": bson.M{"$in": RAUpdateList}}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"action_status":   constants.Canceled,
					"canceled_reason": constants.CanceledBySystemDueToLinkedRequestAction,
				},
			}),
		)

		modelData = append(modelData, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"action_code": actionCode}).
			SetUpdate(bson.M{
				"$set": update, // assuming update is a struct or bson.M
			}),
		)

		_, err := r.collection.BulkWrite(ctx, modelData, options.BulkWrite().SetOrdered(false))
		if err != nil {
			log.Errorf("[CpsAction][Update] failed to make bulk update")
			return nil, local_utils.HandleDBError(err)
		}

		return &update, nil
	}

	// filterMap := bson.M{}
	data, err := r.dal.UpdateOne(ctx, filterMap, updateMap)
	if err != nil {
		log.Errorf("[CPSAction][Update] failed to update CPS action: %v", err)
		return nil, local_utils.HandleDBError(err)
	}

	log.Infof("[CPSAction][Update] CPS action updated successfully")
	return &data, nil
}

func (r *CPSActionStorage) UpdateByActionCode(ctx context.Context, actionCode string, update model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][UpdateByActionCode] updating CPS action for action code: %s", actionCode)
	updateMap := BuildCPSActionUpdateMap(update)
	filterMap := bson.M{"action_code": actionCode}

	data, err := r.dal.UpdateOne(ctx, filterMap, updateMap)
	if err != nil {
		log.Errorf("[CPSAction][UpdateByActionCode] failed to update CPS action: %v", err)
		return nil, local_utils.HandleDBError(err)
	}
	log.Infof("[CPSAction][UpdateByActionCode] CPS action updated successfully")
	return &data, nil
}

func (r *CPSActionStorage) UpdateCustome(ctx context.Context, filter, update bson.M) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][UpdateCustome] updating CPS action")

	_, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[CPSAction][UpdateCustome] failed to update CPS action: %v", err)
		return local_utils.HandleDBError(err)
	}
	log.Infof("[CPSAction][UpdateCustome] CPS action updated successfully")
	return nil
}

func (r *CPSActionStorage) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][Delete] deleting CPS action with ID: %s", id)
	idObj, ok := local_utils.StringToObjectID(id)
	if !ok {
		log.Errorf("[CPSAction][Delete] invalid object id: %s", id)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filterMap := BuildCPSActionFilter(model.CPSAction{ID: idObj})

	err := r.dal.DeleteOne(ctx, filterMap)
	if err != nil {
		log.Errorf("[CPSAction][Delete] failed to delete CPS action: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Infof("[CPSAction][Delete] CPS action deleted successfully with ID: %s", id)
	return nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][SanitizedFindAllWithPagination] finding all CPS actions. Department: %s, Filter: %+v", department, filterParam)
	// 1. Base filter (only active records)
	baseFilter := bson.M{
		"is_deleted": false,
	}

	from, okFrom := filterParam.Filters["created_at_from"].(string)
	to, okTo := filterParam.Filters["created_at_to"].(string)

	parseFlexible := func(s string) (time.Time, bool) {
		if s == "" {
			return time.Time{}, false
		}
		// try RFC3339 first
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t, true
		}
		// try date-only YYYY-MM-DD
		if t, err := time.Parse("2006-01-02", s); err == nil {
			return t, true
		}
		// try without timezone
		if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
			return t, true
		}
		return time.Time{}, false
	}

	if okFrom && okTo && from != "" && to != "" {
		fromTime, ok1 := parseFlexible(from)
		toTime, ok2 := parseFlexible(to)

		if ok1 && ok2 {
			// If inputs were date-only (length 10), expand to full-day bounds.
			if len(from) == 10 {
				fromTime = time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, time.UTC)
			}
			if len(to) == 10 {
				toTime = time.Date(toTime.Year(), toTime.Month(), toTime.Day(), 23, 59, 59, int(time.Millisecond*999), time.UTC)
			}
			baseFilter["created_at"] = bson.M{
				"$gte": fromTime,
				"$lte": toTime,
			}
		} else {
			delete(filterParam.Filters, "created_at_from")
			delete(filterParam.Filters, "created_at_to")
		}
	}

	searchKeys := bson.M{}

	allowedKeys := []string{"action_status", "action_code", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id", "unique_id", "created_at"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"unique_id": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"action_type": searchRegex},
			{"action_code": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	// keep created_at from baseFilter so date range filtering is applied
	filter := dynamicFilter

	// Filter by action codes resolved from user_action_log (set by service).
	// Empty slice means no matching log entries — use never-match to return 0 results.
	if codes, ok := filterParam.Filters["action_code_in"].([]string); ok {
		if len(codes) > 0 {
			filter["action_code"] = bson.M{"$in": codes}
		} else {
			filter["_id"] = bson.M{"$exists": false}
		}
	}
	delete(filter, "action_code_in")

	applyActionStatusFilter(filterParam.Filters, filter)
	applyActionCodeFilter(filterParam.Filters, filter)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		cps_action_core.SanitizePipeline(exclude),
	}

	action, ok := filterParam.Filters["action"]
	if ok && action == "export" {
		pipeline = mongo.Pipeline{
			{{Key: "$match", Value: filter}},
			{{Key: "$project", Value: Projection}},
			cps_action_core.SanitizePipeline(exclude),
		}
	}
	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPagination] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPagination] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPagination] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	log.Infof("[CPSAction][SanitizedFindAllWithPagination] successfully fetched paginated CPS actions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPaginationForApprover(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][SanitizedFindAllWithPaginationForApprover] finding all CPS actions. RAList: %s, Filter: %+v", RAList, filterParam)
	// 1. Base filter (only active records)
	baseFilter := bson.M{
		"is_deleted": false,
	}
	from, okFrom := filterParam.Filters["created_at_from"].(string)
	to, okTo := filterParam.Filters["created_at_to"].(string)

	// Flexible parsing: accept RFC3339, date-only (YYYY-MM-DD), or datetime without TZ.
	parseFlexible := func(s string) (time.Time, bool) {
		if s == "" {
			return time.Time{}, false
		}
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t, true
		}
		if t, err := time.Parse("2006-01-02", s); err == nil {
			return t, true
		}
		if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
			return t, true
		}
		return time.Time{}, false
	}

	if okFrom && okTo && from != "" && to != "" {
		fromTime, ok1 := parseFlexible(from)
		toTime, ok2 := parseFlexible(to)
		if ok1 && ok2 {
			if len(from) == 10 {
				fromTime = time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, time.UTC)
			}
			if len(to) == 10 {
				toTime = time.Date(toTime.Year(), toTime.Month(), toTime.Day(), 23, 59, 59, int(time.Millisecond*999), time.UTC)
			}
			baseFilter["created_at"] = bson.M{
				"$gte": fromTime,
				"$lte": toTime,
			}
		} else {
			delete(filterParam.Filters, "created_at_from")
			delete(filterParam.Filters, "created_at_to")
		}
	}
	searchKeys := bson.M{}

	allowedKeys := []string{"action_status", "action_code", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id", "unique_id", "created_at"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"unique_id": searchRegex},
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"action_type": searchRegex},
			{"action_code": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	// keep created_at from baseFilter so date range filtering is applied
	filter := dynamicFilter
	filter["request_action"] = bson.M{"$in": RAList}

	userFilter := bson.M{
		"$or": []bson.M{
			{"maker_id": userID},
			{"checker_users.checker_id": userID},
		},
	}

	isPendingOnly := func() bool {
		switch v := filterParam.Filters["action_status"].(type) {
		case string:
			return v == "PENDING"
		case []string:
			return len(v) == 1 && v[0] == "PENDING"
		}
		return false
	}()

	var finalMatch bson.M
	if isPendingOnly {
		finalMatch = filter
	} else {
		finalMatch = bson.M{
			"$and": []bson.M{
				filter,
				userFilter,
			},
		}
	}

	applyActionStatusFilter(filterParam.Filters, finalMatch)
	applyActionCodeFilter(filterParam.Filters, finalMatch)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: finalMatch}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		cps_action_core.SanitizePipeline(exclude),
	}

	action, ok := filterParam.Filters["action"]
	if ok && action == "export" {
		pipeline = mongo.Pipeline{
			{{Key: "$match", Value: finalMatch}},
			{{Key: "$project", Value: Projection}},
			cps_action_core.SanitizePipeline(exclude),
		}
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationForApprover] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationForApprover] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := r.dal.TotalCount(ctx, finalMatch)
	if err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationForApprover] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	log.Infof("[CPSAction][SanitizedFindAllWithPaginationForApprover] successfully fetched paginated CPS actions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPaginationForAuditor(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][SanitizedFindAllWithPaginationForAuditor] finding all CPS actions. RAList: %s, Filter: %+v", RAList, filterParam)
	inboxOr := interface{}(nil)
	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["__inbox_or"]; ok {
			inboxOr = v
			delete(filterParam.Filters, "__inbox_or")
		}
	}
	baseFilter := bson.M{
		"is_deleted": false,
	}

	from, okFrom := filterParam.Filters["created_at_from"].(string)
	to, okTo := filterParam.Filters["created_at_to"].(string)

	parseFlexible := func(s string) (time.Time, bool) {
		if s == "" {
			return time.Time{}, false
		}
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t, true
		}
		if t, err := time.Parse("2006-01-02", s); err == nil {
			return t, true
		}
		if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
			return t, true
		}
		return time.Time{}, false
	}

	if okFrom && okTo && from != "" && to != "" {
		fromTime, ok1 := parseFlexible(from)
		toTime, ok2 := parseFlexible(to)

		if ok1 && ok2 {
			if len(from) == 10 {
				fromTime = time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, time.UTC)
			}
			if len(to) == 10 {
				toTime = time.Date(toTime.Year(), toTime.Month(), toTime.Day(), 23, 59, 59, int(time.Millisecond*999), time.UTC)
			}
			baseFilter["created_at"] = bson.M{
				"$gte": fromTime,
				"$lte": toTime,
			}
		} else {
			delete(filterParam.Filters, "created_at_from")
			delete(filterParam.Filters, "created_at_to")
		}
	}

	searchKeys := bson.M{}

	allowedKeys := []string{"action_code", "action_status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id", "auditor_status", "auditor_users.auditor_id", "unique_id", "created_at"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"action_code": searchRegex},
			{"unique_id": searchRegex},
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

	// delete(dynamicFilter, "created_at")
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

	as := filter["action_status"]
	if as == nil || as == "" || as == constants.Pending {
		filter["action_status"] = bson.M{"$in": []string{string(constants.Approved), string(constants.Rejected)}}
	}

	applyActionStatusFilter(filterParam.Filters, filter)
	applyActionCodeFilter(filterParam.Filters, filter)
	// filter["action_status"] = bson.M{"$in": []string{string(constants.Approved), string(constants.Rejected)}}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$project", Value: Projection}},
		cps_action_core.SanitizePipeline(exclude),
	}

	action, ok := filterParam.Filters["action"]
	if ok && action == "export" {
		pipeline = mongo.Pipeline{
			{{Key: "$match", Value: filter}},
			{{Key: "$project", Value: Projection}},
			cps_action_core.SanitizePipeline(exclude),
		}
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationForAuditor] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationForAuditor] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationForAuditor] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	log.Infof("[CPSAction][SanitizedFindAllWithPaginationForAuditor] successfully fetched paginated CPS actions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) SanitizedFindAllWithPaginationCPSActions(ctx context.Context, userID, role string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][SanitizedFindAllWithPaginationCPSActions] finding all CPS actions for User: %s, RAList: %s, Filter: %+v", userID, RAList, filterParam)

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
	allowedKeys := []string{"action_status", "action_code", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "checker_name", "checker_phone_number", "auditor_status", "unique_id", "created_at"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"unique_id": searchRegex},
			{"action_code": searchRegex},
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
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationCPSActions] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationCPSActions] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := r.dal.TotalCount(ctx, finalMatch)
	if err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationCPSActions] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	log.Infof("[CPSAction][SanitizedFindAllWithPaginationCPSActions] successfully fetched paginated CPS actions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) SanitizedFindOne(ctx context.Context, filter bson.M) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		cps_action_core.SanitizePipeline(exclude),
		{{Key: "$limit", Value: 1}},
		{{Key: "$project", Value: Projection}},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[CPSAction][SanitizedFindOne] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}
	defer func() {
		_ = cur.Close(ctx)
	}()

	// Handle empty cursor
	if cur == nil || !cur.Next(ctx) {
		log.Warnf("[CPSAction][SanitizedFindOne] no CPS action found for filter")
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	var result model.CPSAction
	if err := cur.Decode(&result); err != nil {
		log.Errorf("[CPSAction][SanitizedFindOne] failed to decode document: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &result, nil
}

func (r *CPSActionStorage) GetCountByDepartment(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

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
		log.Errorf("[CPSAction][GetCountByDepartment] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() {
		_ = cur.Close(ctx)
	}()

	var result actionDto.CPSActionCountResponse
	if cur.Next(ctx) {
		if err := cur.Decode(&result); err != nil {
			log.Errorf("[CPSAction][GetCountByDepartment] failed to decode document: %v", err)
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
func (r *CPSActionStorage) FindByDateRange(ctx context.Context, filterParam *types.Filter) ([]*model.CPSAction, error) {
	// filter := buildCPSActionDateRangeFilter(start_date, end_date)

	baseFilter := bson.M{
		"is_deleted": false,
	}
	searchKeys := bson.M{}
	//---------------------------------------

	allowedKeys := []string{"action_status", "action_code", "action_type", "request_action", "maker_phone_number", "maker_name", "checker_name", "checker_phone_number", "auditor_status", "unique_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"unique_id": searchRegex},
			{"action_code": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"auditor_status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(*filterParam, searchKeys, allowedKeys)

	//---------------------------------------
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetSkip(skip).
		SetLimit(limit).
		SetProjection(Projection)

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}

	cursor, err := r.collection.Find(ctx, dynamicFilter, opts)
	// actions,err := r.dal.FindAll(ctx,filter,nil)
	if err != nil {
		return nil, err
	}
	// defer cursor.Close(ctx)

	var actions []*model.CPSAction
	if err := cursor.All(ctx, &actions); err != nil {
		return nil, err
	}

	return actions, nil

}

// stream process constant memrory useage O(1)
func (r *CPSActionStorage) StreamByDateRange(
	ctx context.Context,
	filterParam *types.Filter,
	handler func(*model.CPSAction) error,
) error {
	// filter := buildCPSActionDa/teRangeFilter(startDate, endDate)

	// opts := options.Find().
	// 	SetSort(bson.D{{Key: "created_at", Value: 1}}).
	// 	SetBatchSize(1000) // very important

	baseFilter := bson.M{
		"is_deleted": false,
	}
	searchKeys := bson.M{}
	//---------------------------------------

	allowedKeys := []string{"action_status", "action_code", "action_type", "request_action", "maker_phone_number", "maker_name", "checker_name", "checker_phone_number", "auditor_status", "unique_id"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"unique_id": searchRegex},
			{"action_code": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"auditor_status": searchRegex},
			{"action_type": searchRegex},
			{"request_action": searchRegex},
		}
	}

	fieldProjection := filterParam.Filters["fields"]
	if fieldProjection == nil {
		fieldProjection = Projection
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(*filterParam, searchKeys, allowedKeys)

	//---------------------------------------
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetSkip(skip).
		SetLimit(limit).
		SetProjection(fieldProjection)

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}

	cursor, err := r.collection.Find(ctx, dynamicFilter, opts)
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

func (r *CPSActionStorage) ActionByDateRange(ctx context.Context, filterParam types.Filter, RAList []string) ([]*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CPSAction][SanitizedFindAllWithPaginationForApprover] finding all CPS actions. RAList: %s, Filter: %+v", RAList, filterParam)
	userFilter := bson.M{}
	userData := local_utils.ExtractUserFromContext(ctx)
	userID := userData.UserID
	// 1. Base filter (only active records)
	baseFilter := bson.M{
		"is_deleted": false,
	}
	searchKeys := bson.M{}

	allowedKeys := []string{"action_status", "action_code", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id", "unique_id", "created_at"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
			{"unique_id": searchRegex},
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"action_status": searchRegex},
			{"action_type": searchRegex},
			{"action_code": searchRegex},
			{"request_action": searchRegex},
		}
	}

	dynamicFilter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	// delete(dynamicFilter, "created_at")
	filter := dynamicFilter
	if RAList != nil && len(RAList) > 0 {
		filter["request_action"] = bson.M{"$in": RAList}
	}

	if filterParam.Filters["actor"] == constants.Maker {
		userFilter = bson.M{"maker_id": userID}
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
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationForApprover] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var results []*model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[CPSAction][SanitizedFindAllWithPaginationForApprover] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 8. Return standard paginated response
	return results, nil
}

// applyActionStatusFilter converts action_status to a $in query when the caller
// supplies a slice, leaving single-string values handled by FilterBuilder as-is.
func applyActionStatusFilter(filters map[string]interface{}, filter bson.M) {
	raw, ok := filters["action_status"]
	if !ok {
		return
	}
	switch v := raw.(type) {
	case string:
		if v != "" {
			filter["action_status"] = bson.M{"$in": []string{v}}
		}
	case []string:
		if len(v) > 0 {
			filter["action_status"] = bson.M{"$in": v}
		}
	case []interface{}:
		statuses := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				statuses = append(statuses, s)
			}
		}
		if len(statuses) > 0 {
			filter["action_status"] = bson.M{"$in": statuses}
		}
	}
}

func applyActionCodeFilter(filters map[string]interface{}, filter bson.M) {
	raw, ok := filters["action_code"]
	if !ok {
		return
	}
	switch v := raw.(type) {
	case string:
		if v != "" {
			filter["action_code"] = bson.M{"$in": []string{v}}
		}
	case []string:
		if len(v) > 0 {
			filter["action_code"] = bson.M{"$in": v}
		}
	case []interface{}:
		statuses := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				statuses = append(statuses, s)
			}
		}
		if len(statuses) > 0 {
			filter["action_code"] = bson.M{"$in": statuses}
		}
	}
}

func buildCPSActionDateRangeFilter(startDate, endDate time.Time) bson.M {
	rangeFilter := bson.M{
		"$gte": startDate,
		"$lte": endDate,
	}

	// Some records use different timestamp fields; include all known variants.
	return bson.M{
		"$or": []bson.M{
			{"created_at": rangeFilter},
			{"action_created_at": rangeFilter},
			{"maker_action_time": rangeFilter},
			{"last_modified_at": rangeFilter},
		},
	}
}
