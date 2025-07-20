package faydaaccount

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FaydaAccountRepo struct {
	client      *mongo.Client
	logger      utils.Logger
	cpsDal      dal.MongoDal[entity.CPSAction, entity.CPSAction]
	customerDal dal.MongoDal[member.User, member.User]
}

var _ outbound.FaydaAccountRepository = (*FaydaAccountRepo)(nil)

func InitFaydaAccountPersistence(client *mongo.Client, database string,
	cpsCollection []string, logger utils.Logger) *FaydaAccountRepo {
	cpsDal := dal.NewMongoDal[entity.CPSAction, entity.CPSAction](client, database, cpsCollection[0])
	customerDal := dal.NewMongoDal[member.User, member.User](client, database, cpsCollection[1])

	return &FaydaAccountRepo{
		client:      client,
		logger:      logger,
		cpsDal:      cpsDal,
		customerDal: customerDal,
	}
}

func (f *FaydaAccountRepo) InitiateDisableFaydaAccount(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {

	filter := bson.M{
		"maker_phone_number": req.MakerPhoneNumber,
		"action_status":      entity.ActionPending,
		"department":         req.Department,
	}

	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	fmt.Println("check filter--------", filter)
	faydaAccount, err := f.cpsDal.FindOne(ctx, filter, projection)

	if err != nil && err != mongo.ErrNoDocuments {
		f.logger.Errorf("failed to get fayda account", err)

		return nil, fmt.Errorf("internal server error")
	}

	if faydaAccount != nil {
		f.logger.Infof("pending cps action present", req.MakerName, req.MakerID, req.Department)
		return nil, fmt.Errorf("pending cps action present")
	}

	actionData := make(map[string]interface{})
	byte, err := json.Marshal(req.CurrentAction)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(byte, &actionData); err != nil {
		return nil, err
	}
	customerFilter := bson.M{
		"phone_number": actionData["phone_number"],
	}

	customerProjection := bson.M{
		"is_account_blocked": 1,
		"user_code":          1,
		"full_name":          1,
		"phone_number":       1,
	}

	customer, err := f.customerDal.FindOne(ctx, customerFilter, customerProjection)
	if err != nil {
		f.logger.Errorf("failed to get customer account", err)
		err = fmt.Errorf("failed to get customer account %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	if customer.IsAccountBlocked {
		f.logger.Errorf("the user already disabled")
		return nil, fmt.Errorf("ACCOUNT_ALREADY_DISABLED")
	}

	req.ActionStatus = entity.ActionPending
	req.RequestAction = entity.RequestDisableFaydaAccount
	req.ActionType = entity.ActionCreate
	req.PreviousAction = map[string]any{
		"user_code":          customer.UserCode,
		"full_name":          customer.FullName,
		"phone_number":       customer.PhoneNumber,
		"is_account_blocked": customer.IsAccountBlocked,
	}
	req.CurrentAction = map[string]any{
		"is_account_blocked": true,
	}

	req.MakerActionTime = time.Now()
	cpsAction, err := f.cpsDal.InsertOne(ctx, req)
	if err != nil {
		f.logger.Errorf("failed to create cps action", err)

		return nil, fmt.Errorf("FAILED_TO_INSERT_CPS_ACTION")
	}

	return &cpsAction, nil
}

func (f *FaydaAccountRepo) AuthorizeFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {

	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": entity.ActionPending,
	}

	update := bson.M{
		"checker_name":         req.CheckerName,
		"checker_id":           req.CheckerID,
		"checker_phone_number": req.CheckerPhoneNumber,
		"action_status":        entity.ActionApproved,
		"checker_action_time":  time.Now(),
	}

	cpsAction, err := f.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		f.logger.Errorf("failed to update cps action", err)
		return nil, fmt.Errorf("FAILED_TO_FIND_CPS_ACTION")
	}

	actionData := make(map[string]interface{})
	byte, err := json.Marshal(cpsAction.PreviousAction)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byte, &actionData); err != nil {
		return nil, err
	}

	customerFilter := bson.M{
		"phone_number": actionData["phone_number"],
		"kyc.level":    1,
	}
	customerUpdate := bson.M{
		"is_account_blocked": true,
	}

	_, err = f.customerDal.UpdateOne(ctx, customerFilter, customerUpdate)
	if err != nil {
		f.logger.Errorf("failed to update  action", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_USER_DATA")
	}

	return &cpsAction, nil
}

func (f *FaydaAccountRepo) RejectFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error) {
	actionData := make(map[string]interface{})
	byte, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byte, &actionData); err != nil {
		return nil, err
	}
	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": entity.ActionPending,
	}

	update := bson.M{
		"checker_name":           req.CheckerName,
		"checker_id":             req.CheckerID,
		"checker_phone_number":   req.CheckerPhoneNumber,
		"action_status":          entity.ActionRejected,
		"rejected_action_reason": req.RejectionReason,
		"checker_action_time":    time.Now(),
	}

	cpsAction, err := f.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		f.logger.Errorf("failed to update fayda customer status", err)

		return nil, fmt.Errorf("FAILED_TO_UPDATE_USER_DATA")
	}
	return &cpsAction, nil
}

func (f *FaydaAccountRepo) GetAllFaydaAccounts(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	filter := bson.M{
		"is_deleted": false,
	}
	filter["kyc.level"] = int32(1) // Only KYC level 1 users (Fayda accounts)

	projection := bson.M{}

	// Apply search if provided
	if filterParams.Search != "" {
		filter["$or"] = []bson.M{
			{"full_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"phone_number": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"user_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"fayda.id_number": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	// Apply filters if provided
	if filterParams.Filters != "" {
		filter["account_status"] = filterParams.Filters
	}

	// Calculate pagination
	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	// Get total count
	totalDocs, err := f.customerDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated data
	users, err := f.customerDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &common_util.PaginatedResponse[[]*member.User]{
				Data: []*member.User{},
				Meta: common_util.BuildPaginationMeta(0, filterParams.Page, filterParams.PerPage),
			}, nil
		}
		return nil, fmt.Errorf("failed to get fayda accounts: %w", err)
	}

	// Build pagination metadata using utility function
	meta := common_util.BuildPaginationMeta(totalDocs, filterParams.Page, filterParams.PerPage)

	return &common_util.PaginatedResponse[[]*member.User]{
		Data: users,
		Meta: meta,
	}, nil
}
