package faydaaccount

import (
	"context"
	// "encoding/json"
	"errors"

	// "errors"
	"fmt"
	// "strings"
	"time"

	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/fayda_account"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type FaydaAccountRepo struct {
	client            *mongo.Client
	logger            utils.Logger
	mongoDalCpsAction dal.MongoDal[model.CPSAction, model.CPSAction]
	mongoDalCustomer  dal.MongoDal[member.User, member.User]
}

func InitFaydaAccountPersistence(client *mongo.Client, database string, cpsCollection []string, logger utils.Logger) outbound.FaydaRepository {

	mongoDalCpsAction := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, cpsCollection[0])
	mongoDalCustomer := dal.NewMongoDal[member.User, member.User](client, database, cpsCollection[1])

	return &FaydaAccountRepo{
		client:            client,
		logger:            logger,
		mongoDalCpsAction: mongoDalCpsAction,
		mongoDalCustomer:  mongoDalCustomer,
	}
}

func (f *FaydaAccountRepo) EnableFaydaAccount(ctx context.Context, req *entities.CPSAction, userID string) (string, error) {
	objID, err := bson.ObjectIDFromHex(userID)

	if err != nil {
		f.logger.Errorf("[persistance.EnableFaydaAccount]can't convert string ID to objectID", err)
		return "", err
	}
	filter := bson.M{
		"_id": objID,
	}

	userData, err := f.mongoDalCustomer.FindOne(ctx, filter, bson.M{})
	if err != nil {
		f.logger.Errorf("[persistance.enableFaydaAccount]failed to find user")
		return "", fmt.Errorf("FAILED_TO_FIND_USER")
	}
	if userData.KYCLevel != 1 {
		f.logger.Errorf("user not a fayda user")
		return "", fmt.Errorf("USER_IS_NOT_FAYDA_USER")
	}
	if !userData.IsBlocked {
		f.logger.Infof("user already enabled")
		return "", fmt.Errorf("FAYDA_USER_ALREADY_ENABLED")
	}

	filterAction := bson.M{
		"department":     req.Department,
		"action_type":    "UPDATE",
		"action_status":  "PENDING",
		"request_action": constant.RequestEnableFaydaAccount,
	}

	existing, err := f.mongoDalCpsAction.FindOne(ctx, filterAction, bson.M{})
	if err != nil && !(errors.Is(err, mongo.ErrNoDocuments)) {
		f.logger.Errorf("[persistance.enableFaydaAccount]unable to feach cps action ")
		return "", fmt.Errorf("FAILED_TO_GET_FAYDA_ACCOUNT")
	}

	if existing != nil {
		f.logger.Errorf("[persistance.enableFaydaAccount]pending action already exist")
		return "", fmt.Errorf("PENDING_CPS_ACTION_PRESENT")
	}

	previous := *userData
	userData.IsBlocked = false

	CPSAction := model.CPSAction{
		ID:                 bson.NewObjectID(),
		ActionCode:         req.ActionCode,
		MakerID:            req.MakerID,
		MakerName:          req.MakerName,
		MakerPhoneNumber:   req.MakerPhoneNumber,
		CheckerID:          req.CheckerID,
		CheckerName:        req.CheckerName,
		CheckerPhoneNumber: req.CheckerPhoneNumber,
		Department:         req.Department,
		RejectionReason:    req.RejectionReason,
		PreviousAction:     previous,
		CurrentAction:      userData,
		ActionStatus:       string(req.ActionStatus),
		ActionType:         string(req.ActionType),
		RequestAction:      string(req.RequestAction),
		CreatedAt:          req.CreatedAt,
		LastModifiedAt:     req.LastModifiedAt,
		MakerActionTime:    req.MakerActionTime,
	}

	createdAction, err := f.mongoDalCpsAction.InsertOne(ctx, CPSAction)

	if err != nil {
		return "", err
	}

	return createdAction.ActionCode, nil
}

func (f *FaydaAccountRepo) DisableFaydaAccount(ctx context.Context, req *entities.CPSAction, userID string) (string, error) {
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		f.logger.Errorf("[persistance.EnableFaydaAccount]can't convert string ID to objectID", err)
		return "", err
	}
	filter := bson.M{
		"_id": objID,
	}

	userData, err := f.mongoDalCustomer.FindOne(ctx, filter, bson.M{})
	if err != nil {
		f.logger.Errorf("[persistance.enableFaydaAccount]fail to find user")
		return "", fmt.Errorf("FAILED_TO_FIND_USER")
	}
	if userData.KYCLevel != 1 {
		f.logger.Errorf("user not a fayda user")
		return "", fmt.Errorf("USER_IS_NOT_FAYDA_USER")
	}
	if userData.IsBlocked {
		f.logger.Infof("user already disabled")
		return "", fmt.Errorf("FAYDA_USER_ALREADY_DISABLED")
	}

	filterAction := bson.M{
		"department":     req.Department,
		"action_type":    "UPDATE",
		"action_status":  "PENDING",
		"request_action": constant.RequestDisableFaydaAccount,
	}

	existing, err := f.mongoDalCpsAction.FindOne(ctx, filterAction, bson.M{})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		f.logger.Errorf("[persistance.enableFaydaAccount]unable to feach cps action ")
		return "", fmt.Errorf("FAILED_TO_GET_FAYDA_ACCOUNT")
	}

	if existing != nil {
		f.logger.Errorf("[persistance.enableFaydaAccount]pending action already exist")
		return "", fmt.Errorf("PENDING_CPS_ACTION_PRESENT")
	}

	previous := *userData
	userData.IsBlocked = false

	CPSAction := model.CPSAction{
		ID:                 bson.NewObjectID(),
		ActionCode:         req.ActionCode,
		MakerID:            req.MakerID,
		MakerName:          req.MakerName,
		MakerPhoneNumber:   req.MakerPhoneNumber,
		CheckerID:          req.CheckerID,
		CheckerName:        req.CheckerName,
		CheckerPhoneNumber: req.CheckerPhoneNumber,
		Department:         req.Department,
		RejectionReason:    req.RejectionReason,
		PreviousAction:     previous,
		CurrentAction:      userData,
		ActionStatus:       string(req.ActionStatus),
		ActionType:         string(req.ActionType),
		RequestAction:      string(req.RequestAction),
		CreatedAt:          req.CreatedAt,
		LastModifiedAt:     req.LastModifiedAt,
		MakerActionTime:    req.MakerActionTime,
	}

	createdAction, err := f.mongoDalCpsAction.InsertOne(ctx, CPSAction)

	if err != nil {
		return "", err
	}

	return createdAction.ActionCode, nil

}

func (f *FaydaAccountRepo) AuthorizeBlockFaydaUser(ctx context.Context, cpsAction *entities.CPSAction) (*model.CPSAction, error) {

	var actionData member.User
	raw, _ := bson.Marshal(cpsAction.CurrentAction)
	if err := bson.Unmarshal(raw, &actionData); err != nil {
		f.logger.Errorf("failed to unmarshal action data for authorization, action_code: %s", cpsAction.ActionCode)
		return nil, fmt.Errorf(error_codes.InvalidActionData)
	}

	currActionData := cpsAction.CurrentAction

	dataMap, ok := currActionData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("CurrentAction is not a map[string]interface{}")
	}

	idVal, ok := dataMap["_id"]
	if !ok {
		return nil, fmt.Errorf("_id not found in CurrentAction")
	}

	idStr, ok := idVal.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("_id is not a string, got: %T", idVal)
	}

	userUpdate := bson.M{
		"is_blocked":       true,
		"last_modified_at": time.Now(),
	}

	_, err := f.mongoDalCustomer.UpdateOne(ctx, bson.M{"_id": idStr}, userUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to block user ")
	}

	CPSAction := model.CPSAction{
		ID:                 idStr,
		ActionCode:         cpsAction.ActionCode,
		MakerID:            cpsAction.MakerID,
		MakerName:          cpsAction.MakerName,
		MakerPhoneNumber:   cpsAction.MakerPhoneNumber,
		CheckerID:          cpsAction.CheckerID,
		CheckerName:        cpsAction.CheckerName,
		CheckerPhoneNumber: cpsAction.CheckerPhoneNumber,
		Department:         cpsAction.Department,
		RejectionReason:    cpsAction.RejectionReason,
		PreviousAction:     cpsAction.PreviousAction,
		CurrentAction:      cpsAction.CurrentAction,
		ActionStatus:       string(cpsAction.ActionStatus),
		ActionType:         string(cpsAction.ActionType),
		RequestAction:      string(cpsAction.RequestAction),
		CreatedAt:          cpsAction.CreatedAt,
		LastModifiedAt:     cpsAction.LastModifiedAt,
		MakerActionTime:    cpsAction.MakerActionTime,
	}
	return &CPSAction, nil

}
func (f *FaydaAccountRepo) AuthorizeEnableFaydaUser(ctx context.Context, cpsAction *entities.CPSAction) (*model.CPSAction, error) {

	var actionData member.User
	raw, _ := bson.Marshal(cpsAction.CurrentAction)
	if err := bson.Unmarshal(raw, &actionData); err != nil {
		f.logger.Errorf("failed to unmarshal action data for authorization, action_code: %s", cpsAction.ActionCode)
		return nil, fmt.Errorf(error_codes.InvalidActionData)
	}

	currActionData := cpsAction.CurrentAction

	dataMap, ok := currActionData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("CurrentAction is not a map[string]interface{}")
	}

	idVal, ok := dataMap["_id"]
	if !ok {
		return nil, fmt.Errorf("_id not found in CurrentAction")
	}

	idStr, ok := idVal.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("_id is not a string, got: %T", idVal)
	}

	userUpdate := bson.M{
		"is_blocked":       false,
		"last_modified_at": time.Now(),
	}

	_, err := f.mongoDalCustomer.UpdateOne(ctx, bson.M{"_id": idStr}, userUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to block user ")
	}

	CPSAction := model.CPSAction{
		ID:                 idStr,
		ActionCode:         cpsAction.ActionCode,
		MakerID:            cpsAction.MakerID,
		MakerName:          cpsAction.MakerName,
		MakerPhoneNumber:   cpsAction.MakerPhoneNumber,
		CheckerID:          cpsAction.CheckerID,
		CheckerName:        cpsAction.CheckerName,
		CheckerPhoneNumber: cpsAction.CheckerPhoneNumber,
		Department:         cpsAction.Department,
		RejectionReason:    cpsAction.RejectionReason,
		PreviousAction:     cpsAction.PreviousAction,
		CurrentAction:      cpsAction.CurrentAction,
		ActionStatus:       string(cpsAction.ActionStatus),
		ActionType:         string(cpsAction.ActionType),
		RequestAction:      string(cpsAction.RequestAction),
		CreatedAt:          cpsAction.CreatedAt,
		LastModifiedAt:     cpsAction.LastModifiedAt,
		MakerActionTime:    cpsAction.MakerActionTime,
	}
	return &CPSAction, nil
}
