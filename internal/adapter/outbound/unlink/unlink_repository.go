package unlink

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/unlink"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UnlinkRepo struct {
	client     *mongo.Client
	user       dal.MongoDal[member.User, member.User]
	actionRepo dal.MongoDal[entities.CPSAction, entities.CPSAction]
	otpRepo    dal.MongoDal[member.OTP, member.OTP]
	logger     utils.Logger
}

var _ outbound.UnlinkRepository = (*UnlinkRepo)(nil)

func NewUnlinkInfrastructure(client *mongo.Client, dbName string, collectionName []string, logger utils.Logger) *UnlinkRepo {
	user := dal.NewMongoDal[member.User, member.User](client, dbName, collectionName[0])
	otpRepo := dal.NewMongoDal[member.OTP, member.OTP](client, dbName, collectionName[1])
	cpsRepo := dal.NewMongoDal[entities.CPSAction, entities.CPSAction](client, dbName, collectionName[2])

	return &UnlinkRepo{
		user:       user,
		actionRepo: cpsRepo,
		otpRepo:    otpRepo,
		logger:     logger,
	}
}

func (u *UnlinkRepo) UnlinkDevice(ctx context.Context, userCode string, cpsAction entities.CPSAction) (string, error) {
	filter := bson.M{"user_code": userCode}

	user, err := u.user.FindOne(ctx, filter, bson.M{})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			u.logger.Warnf("[UnlinkDevice] user not found for userCode=%s", userCode)
			return "", fmt.Errorf(error_codes.AuthUserNotFound)
		}

		u.logger.Warnf("[UnlinkDevice] query failed for userCode=%s: %v", userCode, err)
		return "", fmt.Errorf(error_codes.UnhandledServerError)
	}

	if user.IsAccountBlocked {
		u.logger.Warnf("[UnlinkDevice] user %s is blocked", userCode)
		return "", fmt.Errorf(error_codes.AccountBlocked)
	}

	pendingFilter := bson.M{
		"action_code":    cpsAction.ActionCode,
		"action_status":  "PENDING",
		"request_action": string(entities.UnlinkDevice),
	}

	existingPendingAction, err := u.actionRepo.FindOne(ctx, pendingFilter, bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		u.logger.Errorf("[UnlinkDevice] failed to check existing pending actions: %v", err)
		return "", fmt.Errorf(error_codes.PendingActionCheckFailed)
	}

	if existingPendingAction != nil {
		if existingPendingAction.Department == cpsAction.Department {
			u.logger.Warnf("[UnlinkDevice] pending request already exists in same department: %s", userCode)
			return "", fmt.Errorf(error_codes.PendingRequestExists)
		}

		updateFilter := bson.M{
			"action_code":   cpsAction.ActionCode,
			"action_status": "PENDING",
		}
		update := bson.M{
			"action_status": "REJECTED",
		}
		_, err := u.actionRepo.UpdateOne(ctx, updateFilter, update)
		if err != nil {
			u.logger.Errorf("[UnlinkDevice] failed to reject old pending action for user %s: %v", userCode, err)
			return "", fmt.Errorf(error_codes.PendingActionRejectionFailed)
		}
	}

	cpsAction.CreatedAt = time.Now()
	cpsAction.LastModifiedAt = time.Now()

	cpsAction.CurrentAction = map[string]interface{}{
		"action_type": "UNLINK_DEVICE",
		"user_code":   userCode,
		"user_id":     user.ID,
		"full_name":   user.FullName,
	}

	cpsData, err := u.actionRepo.InsertOne(ctx, cpsAction)
	if err != nil {
		u.logger.Errorf("[UnlinkDevice] failed to insert unlink action for userCode=%s: %v", userCode, err)
		return "", fmt.Errorf(error_codes.UnlinkActionFailed)
	}

	return cpsData.ActionCode, nil
}

func (u *UnlinkRepo) Authorize(ctx context.Context, cpsAction *cps_entities.CPSAction) (*cps_entities.CPSAction, error) {
	curAction, err := u.ExtractUnlinkDeviceAction(cpsAction)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"user_code": curAction.UserCode}
	user, err := u.user.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			u.logger.Warnf("[Authorize] user not found for userCode=%s", cpsAction.ActionCode)
			return nil, fmt.Errorf(error_codes.AuthUserNotFound)
		}
		u.logger.Warnf("[Authorize] failed to find user: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	updateUser := bson.M{
		"device.device_uuid":    "",
		"device_status":         "UNLINKED",
		"login_pin.pin":         "",
		"bps_reject_status":     "AUTHORIZED",
		"login_pin.pin_history": user.LoginPIN.PIN,
	}

	if _, err = u.user.UpdateOne(ctx, filter, updateUser); err != nil {
		u.logger.Errorf("[Authorize] failed to update user info: %v", err)
		return nil, fmt.Errorf(error_codes.ActionApprovalFailed)
	}

	u.logger.Infof("[Authorize] unlink approved for user: %s", cpsAction.ActionCode)
	return cpsAction, nil
}

func (u *UnlinkRepo) ExtractUnlinkDeviceAction(cpsAction *cps_entities.CPSAction) (entities.UnlinkDeviceAction, error) {
	var action entities.UnlinkDeviceAction

	raw, err := bson.Marshal(cpsAction.CurrentAction)
	if err != nil {
		return action, fmt.Errorf("marshal current action", err.Error(), error_codes.InvalidActionData)
	}

	if err := bson.Unmarshal(raw, &action); err != nil {
		return action, fmt.Errorf("unmarshal current action", err.Error(), error_codes.InvalidActionData)
	}

	return action, nil
}
