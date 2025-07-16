package faydaaccount

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"

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
	req.PreviosAction = map[string]any{
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
	byte, err := json.Marshal(cpsAction.PreviosAction)

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
