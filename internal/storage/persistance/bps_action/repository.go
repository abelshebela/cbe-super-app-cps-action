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

	bps_action "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	// bps_action "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
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
	"maker_reason":          1,
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
	// Additional fields needed for auditor view
	"user_information":     1,
	"business":             1,
	"checkers_needed":      1,
	"checkers_approved":    1,
	"checker_id":           1,
	"maker_user":           1,
	"checker_name":         1,
	"checker_phone_number": 1,
	"action_reason":        1,
	"maker_mid":            1,
	"checker_mid":          1,
	"auditor_mid":          1,
	"checker_name_list":    1,
	"auditor_name_list":    1,
	"auditors":             1,
	"checker_time":         1,
	"auditor_time":         1,
	"value":                1,
	"home_branch":          1,
	"account_branch_code":  1,
	"district_code":        1,
	"branch_code":          1,
	"linked_district_code": 1,
	"account_number":       1,
	"account_holder_name":  1,
	"service_name":         1,
	"is_auditor_approved":  1,
}

type bpsActionRepository struct {
	client     *mongo.Client
	actionDal  dal.MongoDal[types.BPSActionDocument, types.BPSActionDocument]
	collection *mongo.Collection
	logger     utils.Logger
}

func NewBPSActionRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger, cfg *config.VaultConfig) storage.BPSActionRepository {
	return &bpsActionRepository{
		client:     client,
		actionDal:  dal.NewMongoDal[types.BPSActionDocument, types.BPSActionDocument](client, cfg, dbName, collection),
		collection: client.Database(dbName).Collection(collection),
		logger:     logger,
	}
}

// GetBPSActionByUserID implements [storage.BPSActionRepository].
func (b *bpsActionRepository) GetBPSActionByUserID(ctx context.Context, userID string, filterParam types.Filter) (types.PaginatedResponse[[]bps_action.BPSAction], error) {
	b.logger.Infof("[BPSAction][GetBPSActionByUserID] fetching BPS actions for user ID: %s", userID)
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, nil)
	filter["user_information.user_code"] = userID

	b.logger.Infof("[BPSAction][GetBPSActionByUserID] constructed filter: %v", filter)
	dbActions, err := b.actionDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		b.logger.Errorf("[BPSAction][GetBPSActionByUserID] failed to fetch BPS actions: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Map BPSActionDocument to BPSAction entities
	var actions []bps_action.BPSAction
	for _, dbAction := range dbActions {
		action := mapBpsActionToEntity(dbAction)
		actions = append(actions, action)
	}

	filter = bson.M{"user_information.user_code": userID}
	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[BPSAction][GetBPSActionByUserID] failed to fetch total count: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return types.PaginatedResponse[[]bps_action.BPSAction]{
		Data: actions,
		Meta: meta,
	}, nil
}

func (b *bpsActionRepository) Save(ctx context.Context, action *bps_action.BPSAction) error {
	b.logger.Infof("[BPSAction][Save] saving BPS action")

	// Convert BPSAction to BPSActionDocument for storage
	dbAction := types.BPSActionDocument{
		ID:                bson.NewObjectID(),
		ActionCode:        action.ActionCode,
		IsAuditorApproved: action.IsAuditorApproved,
		UserInformation: struct {
			UserID         interface{} `json:"user_id" bson:"user_id"`
			UserCode       string      `json:"user_code" bson:"user_code"`
			FullName       string      `json:"full_name" bson:"full_name"`
			AccountNumbers []string    `json:"account_numbers" bson:"account_numbers"`
			PhoneNumbers   string      `json:"phone_numbers" bson:"phone_numbers"`
			BranchCode     string      `json:"branch_code" bson:"branch_code"`
		}{
			UserID:         action.UserInformation.UserID,
			UserCode:       action.UserInformation.UserCode,
			FullName:       action.UserInformation.FullName,
			AccountNumbers: action.UserInformation.AccountNumbers,
			PhoneNumbers:   action.UserInformation.PhoneNumbers,
			BranchCode:     action.UserInformation.BranchCode,
		},
		BusinessInformation: struct {
			BusinessID   interface{} `json:"business_id" bson:"business_id"`
			TILLNumber   string      `json:"till_number" bson:"till_number"`
			BusinessName string      `json:"business_name" bson:"business_name"`
		}{
			BusinessID:   action.BusinessInformation.BusinessID,
			TILLNumber:   action.BusinessInformation.TILLNumber,
			BusinessName: action.BusinessInformation.BusinessName,
		},
		CheckersNeeded:     action.CheckersNeeded,
		CheckersApproved:   action.CheckersApproved,
		CheckerID:          action.CheckerID,
		MakerID:            action.MakerID,
		MakerName:          action.MakerName,
		MakerReason:        action.MakerReason,
		MakerPhoneNumber:   action.MakerPhoneNumber,
		CheckerName:        action.CheckerName,
		CheckerPhoneNumber: action.CheckerPhoneNumber,
		ActionReason: struct {
			ActionType string `json:"action_type" bson:"action_type"`
			ActionNote string `json:"action_note" bson:"action_note"`
			Identifier string `json:"identifier" bson:"identifier"`
		}{
			ActionType: string(action.ActionReason.ActionType),
			ActionNote: action.ActionReason.ActionNote,
			Identifier: action.ActionReason.Identifier,
		},
		MakerMID:        action.MakerMID,
		CheckerMID:      action.CheckerMID,
		AuditorMID:      action.AuditorMID,
		CheckerNameList: action.CheckerNameList,
		AuditorNameList: action.AuditorNameList,
		Auditors: struct {
			AuditorName        string   `json:"auditor_name" bson:"auditor_name"`
			AuditorPhoneNumber string   `json:"auditor_phone_number" bson:"auditor_phone_number"`
			AuditorsRequired   int      `json:"auditors_required" bson:"auditors_required"`
			AuditorID          []string `json:"auditor_id" bson:"auditor_id"`
			Audited            bool     `json:"audited" bson:"audited"`
			AuditorApproval    bool     `json:"auditor_approval" bson:"auditor_approval"`
			Reason             string   `json:"reason" bson:"reason"`
		}{
			AuditorName:        action.Auditors.AuditorName,
			AuditorPhoneNumber: action.Auditors.AuditorPhoneNumber,
			AuditorsRequired:   action.Auditors.AuditorsRequired,
			AuditorID:          action.Auditors.AuditorID,
			Audited:            action.Auditors.Audited,
			AuditorApproval:    action.Auditors.AuditorApproval,
			Reason:             action.Auditors.Reason,
		},
		CheckerTime:        action.CheckerTime,
		AuditorTime:        action.AuditorTime,
		RequestAction:      string(action.RequestAction),
		EntityIdentifyer:   action.EntityIdentifyer,
		HomeBranch:         action.HomeBranch,
		AccountBranchCode:  action.AccountBranchCode,
		DistrictCode:       action.DistrictCode,
		BranchCode:         action.BranchCode,
		LinkedDistrictCode: action.LinkedDistrictCode,
		AccountNumber:      action.AccountNumber,
		AccountHolderName:  action.AccountHolderName,
		ServiceName:        action.ServiceName,
		CurrentAction:      action.CurrentAction,
		PreviousAction:     action.PreviousAction,
		VerifiedAt:         action.VerifiedAt,
		Status:             action.Status,
		CreatedAt:          action.CreatedAt,
		LastModifiedAt:     action.LastModifiedAt,
	}

	_, err := b.actionDal.InsertOne(ctx, dbAction)
	if err != nil {
		b.logger.Errorf("[BPSAction][Save] failed to save BPS action: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (b *bpsActionRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (types.PaginatedResponse[[]bps_action.BPSAction], error) {
	b.logger.Infof("[BPSAction][FindAllWithPagination] fetching BPS actions with pagination for department: %s", department)

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

	dbActions, err := b.actionDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		b.logger.Errorf("[BPSAction][FindAllWithPagination] failed to fetch BPS actions: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Map BPSActionDocument to BPSAction entities
	var actions []bps_action.BPSAction
	for _, dbAction := range dbActions {
		action := mapBpsActionToEntity(dbAction)
		actions = append(actions, action)
	}

	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[BPSAction][FindAllWithPagination] failed to count BPS actions: %v", err)
		return types.PaginatedResponse[[]bps_action.BPSAction]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return types.PaginatedResponse[[]bps_action.BPSAction]{
		Data: actions,
		Meta: meta,
	}, nil
}

func (b *bpsActionRepository) FindOne(ctx context.Context, filter bson.M) (*bps_action.BPSAction, error) {
	b.logger.Infof("[BPSAction][FindOne] fetching BPS action")

	dbAction, err := b.actionDal.FindOne(ctx, filter, Projection)
	if err != nil {
		b.logger.Errorf("[BPSAction][FindOne] failed to find BPS action: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// Map BPSActionDocument to BPSAction entity
	action := mapBpsActionToEntity(*dbAction)
	return &action, nil
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
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPagination] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var dbResults []types.BPSActionDocument
	if err := cur.All(ctx, &dbResults); err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPagination] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Map BPSActionDocument to BPSAction entities
	var results []*bps_action.BPSAction
	for _, dbAction := range dbResults {
		action := mapBpsActionToEntity(dbAction)
		results = append(results, &action)
	}

	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPagination] failed to count BPS actions: %v", err)
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

	baseFilter := bson.M{}

	if strings.TrimSpace(filterParam.Search) != "" {
		regex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		baseFilter["$or"] = []bson.M{
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

	for k, v := range baseFilter {
		dynamicFilter[k] = v
	}
	delete(dynamicFilter, "created_at")
	filter := dynamicFilter
	if len(RAList) > 0 {
		filter["request_action"] = bson.M{"$in": RAList}
	} else {
		return &types.PaginatedResponse[[]*bps_action.BPSAction]{
			Data: []*bps_action.BPSAction{},
			Meta: local_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage),
		}, nil
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

	b.logger.Infof("[BPSAction][SanitizedFindAllWithPaginationForApprover] userID: %s, filter: %v", userID, filter)
	cur, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationForApprover] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cur.Close(ctx)

	var dbResults []types.BPSActionDocument
	if err := cur.All(ctx, &dbResults); err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationForApprover] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Map BPSActionDocument to BPSAction entities
	var results []*bps_action.BPSAction
	for _, dbAction := range dbResults {
		action := mapBpsActionToEntity(dbAction)
		results = append(results, &action)
	}

	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationForApprover] failed to count BPS actions: %v", err)
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
	if len(RAList) > 0 {
		filter["request_action"] = bson.M{"$in": RAList}
	} else {
		return &types.PaginatedResponse[[]*bps_action.BPSAction]{}, nil
	}

	switch filter["auditor_status"] {
	case "NOTCHECKED":
		filter["auditors.audited"] = false
		filter["status"] = "APPROVED" // Only show APPROVED actions for auditors
	case "CHECKED":
		filter["auditors.audited"] = true
	}
	delete(filter, "auditor_status")

	exclude := []string{"password", "first_password_set", "login_attempt_count", "is_deleted", "otp_verfy_count", "otp_last_tried_at", "otp_last_verified_at", "permission_group", "permissions", "last_login_attempt", "next_login_attempt", "is_first_time_login", "last_login"}

	b.logger.Infof("[BPSAction][SanitizedFindAllWithPaginationForAuditor] filter*************: %v", filter)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
		// Add type conversion stage to handle byte array to string conversion
		{{Key: "$set", Value: bson.M{
			"action_reason.action_note": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$isArray": "$action_reason.action_note"},
					"then": bson.M{"$toString": "$action_reason.action_note"},
					"else": "$action_reason.action_note",
				},
			},
		}}},
		{{Key: "$project", Value: Projection}},
		bps_action_core.SanitizePipeline(exclude),
	}

	b.logger.Infof("[BPSAction][SanitizedFindAllWithPaginationForAuditor] filter*************: %v", filter)

	cur, err := b.collection.Aggregate(ctx, pipeline)
	if err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationForAuditor] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var dbResults []types.BPSActionDocument
	if err := cur.All(ctx, &dbResults); err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationForAuditor] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Map BPSActionDocument to BPSAction entities
	var results []*bps_action.BPSAction
	for _, dbAction := range dbResults {
		action := mapBpsActionToEntity(dbAction)
		results = append(results, &action)
	}

	total, err := b.actionDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationForAuditor] failed to count BPS actions: %v", err)
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
	if role != "maker" {
		filter["request_action"] = bson.M{"$in": RAList}
	}

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
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationBPSActions] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var dbResults []types.BPSActionDocument
	if err := cur.All(ctx, &dbResults); err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationBPSActions] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Map BPSActionDocument to BPSAction entities
	var results []*bps_action.BPSAction
	for _, dbAction := range dbResults {
		action := mapBpsActionToEntity(dbAction)
		results = append(results, &action)
	}

	total, err := b.actionDal.TotalCount(ctx, finalMatch)
	if err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindAllWithPaginationBPSActions] failed to count BPS actions: %v", err)
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
		b.logger.Errorf("[BPSAction][SanitizedFindOne] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	if cur == nil || !cur.Next(ctx) {
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	var dbResult types.BPSActionDocument
	if err := cur.Decode(&dbResult); err != nil {
		b.logger.Errorf("[BPSAction][SanitizedFindOne] failed to decode document: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Map BPSActionDocument to BPSAction entity
	result := mapBpsActionToEntity(dbResult)
	return &result, nil
}

func (b *bpsActionRepository) Update(ctx context.Context, actionCode string, update bps_action.BPSAction) (*bps_action.BPSAction, error) {
	filter := bson.M{"action_code": actionCode}

	// Convert BPSAction to BPSActionDocument for update
	dbUpdate := types.BPSActionDocument{
		ActionCode:        update.ActionCode,
		IsAuditorApproved: update.IsAuditorApproved,
		UserInformation: struct {
			UserID         interface{} `json:"user_id" bson:"user_id"`
			UserCode       string      `json:"user_code" bson:"user_code"`
			FullName       string      `json:"full_name" bson:"full_name"`
			AccountNumbers []string    `json:"account_numbers" bson:"account_numbers"`
			PhoneNumbers   string      `json:"phone_numbers" bson:"phone_numbers"`
			BranchCode     string      `json:"branch_code" bson:"branch_code"`
		}{
			UserID:         update.UserInformation.UserID,
			UserCode:       update.UserInformation.UserCode,
			FullName:       update.UserInformation.FullName,
			AccountNumbers: update.UserInformation.AccountNumbers,
			PhoneNumbers:   update.UserInformation.PhoneNumbers,
			BranchCode:     update.UserInformation.BranchCode,
		},
		BusinessInformation: struct {
			BusinessID   interface{} `json:"business_id" bson:"business_id"`
			TILLNumber   string      `json:"till_number" bson:"till_number"`
			BusinessName string      `json:"business_name" bson:"business_name"`
		}{
			BusinessID:   update.BusinessInformation.BusinessID,
			TILLNumber:   update.BusinessInformation.TILLNumber,
			BusinessName: update.BusinessInformation.BusinessName,
		},
		CheckersNeeded:     update.CheckersNeeded,
		CheckersApproved:   update.CheckersApproved,
		CheckerID:          update.CheckerID,
		MakerID:            update.MakerID,
		MakerName:          update.MakerName,
		MakerReason:        update.MakerReason,
		MakerPhoneNumber:   update.MakerPhoneNumber,
		CheckerName:        update.CheckerName,
		CheckerPhoneNumber: update.CheckerPhoneNumber,
		ActionReason: struct {
			ActionType string `json:"action_type" bson:"action_type"`
			ActionNote string `json:"action_note" bson:"action_note"`
			Identifier string `json:"identifier" bson:"identifier"`
		}{
			ActionType: string(update.ActionReason.ActionType),
			ActionNote: update.ActionReason.ActionNote,
			Identifier: update.ActionReason.Identifier,
		},
		MakerMID:        update.MakerMID,
		CheckerMID:      update.CheckerMID,
		AuditorMID:      update.AuditorMID,
		CheckerNameList: update.CheckerNameList,
		AuditorNameList: update.AuditorNameList,
		Auditors: struct {
			AuditorName        string   `json:"auditor_name" bson:"auditor_name"`
			AuditorPhoneNumber string   `json:"auditor_phone_number" bson:"auditor_phone_number"`
			AuditorsRequired   int      `json:"auditors_required" bson:"auditors_required"`
			AuditorID          []string `json:"auditor_id" bson:"auditor_id"`
			Audited            bool     `json:"audited" bson:"audited"`
			AuditorApproval    bool     `json:"auditor_approval" bson:"auditor_approval"`
			Reason             string   `json:"reason" bson:"reason"`
		}{
			AuditorName:        update.Auditors.AuditorName,
			AuditorPhoneNumber: update.Auditors.AuditorPhoneNumber,
			AuditorsRequired:   update.Auditors.AuditorsRequired,
			AuditorID:          update.Auditors.AuditorID,
			Audited:            update.Auditors.Audited,
			AuditorApproval:    update.Auditors.AuditorApproval,
			Reason:             update.Auditors.Reason,
		},
		CheckerTime:        update.CheckerTime,
		AuditorTime:        update.AuditorTime,
		RequestAction:      string(update.RequestAction),
		EntityIdentifyer:   update.EntityIdentifyer,
		HomeBranch:         update.HomeBranch,
		AccountBranchCode:  update.AccountBranchCode,
		DistrictCode:       update.DistrictCode,
		BranchCode:         update.BranchCode,
		LinkedDistrictCode: update.LinkedDistrictCode,
		AccountNumber:      update.AccountNumber,
		AccountHolderName:  update.AccountHolderName,
		ServiceName:        update.ServiceName,
		CurrentAction:      update.CurrentAction,
		PreviousAction:     update.PreviousAction,
		VerifiedAt:         update.VerifiedAt,
		Status:             update.Status,
		CreatedAt:          update.CreatedAt,
		LastModifiedAt:     update.LastModifiedAt,
	}

	updateMap := bson.M{"$set": dbUpdate}
	data, err := b.actionDal.UpdateOne(ctx, filter, updateMap)
	if err != nil {
		b.logger.Errorf("[BPSAction][Update] failed to update BPS action: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// Map result back to BPSAction
	result := mapBpsActionToEntity(data)
	return &result, nil
}

func (b *bpsActionRepository) UpdateByActionCode(ctx context.Context, actionCode string, update bps_action.BPSAction) (*bps_action.BPSAction, error) {
	filter := bson.M{"action_code": actionCode}

	// Convert BPSAction to BPSActionDocument for update
	dbUpdate := types.BPSActionDocument{
		ActionCode:        update.ActionCode,
		IsAuditorApproved: update.IsAuditorApproved,
		UserInformation: struct {
			UserID         interface{} `json:"user_id" bson:"user_id"`
			UserCode       string      `json:"user_code" bson:"user_code"`
			FullName       string      `json:"full_name" bson:"full_name"`
			AccountNumbers []string    `json:"account_numbers" bson:"account_numbers"`
			PhoneNumbers   string      `json:"phone_numbers" bson:"phone_numbers"`
			BranchCode     string      `json:"branch_code" bson:"branch_code"`
		}{
			UserID:         update.UserInformation.UserID,
			UserCode:       update.UserInformation.UserCode,
			FullName:       update.UserInformation.FullName,
			AccountNumbers: update.UserInformation.AccountNumbers,
			PhoneNumbers:   update.UserInformation.PhoneNumbers,
			BranchCode:     update.UserInformation.BranchCode,
		},
		BusinessInformation: struct {
			BusinessID   interface{} `json:"business_id" bson:"business_id"`
			TILLNumber   string      `json:"till_number" bson:"till_number"`
			BusinessName string      `json:"business_name" bson:"business_name"`
		}{
			BusinessID:   update.BusinessInformation.BusinessID,
			TILLNumber:   update.BusinessInformation.TILLNumber,
			BusinessName: update.BusinessInformation.BusinessName,
		},
		CheckersNeeded:     update.CheckersNeeded,
		CheckersApproved:   update.CheckersApproved,
		CheckerID:          update.CheckerID,
		MakerID:            update.MakerID,
		MakerName:          update.MakerName,
		MakerReason:        update.MakerReason,
		MakerPhoneNumber:   update.MakerPhoneNumber,
		CheckerName:        update.CheckerName,
		CheckerPhoneNumber: update.CheckerPhoneNumber,
		ActionReason: struct {
			ActionType string `json:"action_type" bson:"action_type"`
			ActionNote string `json:"action_note" bson:"action_note"`
			Identifier string `json:"identifier" bson:"identifier"`
		}{
			ActionType: string(update.ActionReason.ActionType),
			ActionNote: update.ActionReason.ActionNote,
			Identifier: update.ActionReason.Identifier,
		},
		MakerMID:        update.MakerMID,
		CheckerMID:      update.CheckerMID,
		AuditorMID:      update.AuditorMID,
		CheckerNameList: update.CheckerNameList,
		AuditorNameList: update.AuditorNameList,
		Auditors: struct {
			AuditorName        string   `json:"auditor_name" bson:"auditor_name"`
			AuditorPhoneNumber string   `json:"auditor_phone_number" bson:"auditor_phone_number"`
			AuditorsRequired   int      `json:"auditors_required" bson:"auditors_required"`
			AuditorID          []string `json:"auditor_id" bson:"auditor_id"`
			Audited            bool     `json:"audited" bson:"audited"`
			AuditorApproval    bool     `json:"auditor_approval" bson:"auditor_approval"`
			Reason             string   `json:"reason" bson:"reason"`
		}{
			AuditorName:        update.Auditors.AuditorName,
			AuditorPhoneNumber: update.Auditors.AuditorPhoneNumber,
			AuditorsRequired:   update.Auditors.AuditorsRequired,
			AuditorID:          update.Auditors.AuditorID,
			Audited:            update.Auditors.Audited,
			AuditorApproval:    update.Auditors.AuditorApproval,
			Reason:             update.Auditors.Reason,
		},
		CheckerTime:        update.CheckerTime,
		AuditorTime:        update.AuditorTime,
		RequestAction:      string(update.RequestAction),
		EntityIdentifyer:   update.EntityIdentifyer,
		HomeBranch:         update.HomeBranch,
		AccountBranchCode:  update.AccountBranchCode,
		DistrictCode:       update.DistrictCode,
		BranchCode:         update.BranchCode,
		LinkedDistrictCode: update.LinkedDistrictCode,
		AccountNumber:      update.AccountNumber,
		AccountHolderName:  update.AccountHolderName,
		ServiceName:        update.ServiceName,
		CurrentAction:      update.CurrentAction,
		PreviousAction:     update.PreviousAction,
		VerifiedAt:         update.VerifiedAt,
		Status:             update.Status,
		CreatedAt:          update.CreatedAt,
		LastModifiedAt:     update.LastModifiedAt,
	}

	updateMap := bson.M{"$set": dbUpdate}
	data, err := b.actionDal.UpdateOne(ctx, filter, updateMap)
	if err != nil {
		b.logger.Errorf("[BPSAction][UpdateByActionCode] failed to update BPS action: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// Map result back to BPSAction
	result := mapBpsActionToEntity(data)
	return &result, nil
}

func (b *bpsActionRepository) UpdateCustome(ctx context.Context, filter, update bson.M) error {
	_, err := b.actionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("[BPSAction][UpdateCustome] failed to update BPS action: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (b *bpsActionRepository) Delete(ctx context.Context, id string) error {
	objID, ok := local_util.StringToObjectID(id)
	if !ok {
		b.logger.Errorf("[BPSAction][Delete] invalid object id: %s", id)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if err := b.actionDal.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		b.logger.Errorf("[BPSAction][Delete] failed to delete BPS action: %v", err)
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
