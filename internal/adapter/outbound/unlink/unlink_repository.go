package unlink

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/unlink"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

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

func (u *UnlinkRepo) UnlinkDevice(userCode string, cpsAction entities.CPSAction) error {
	u.logger.Infof("[UnlinkDevice] initiated for user: %s by %s", userCode, cpsAction.MakerID)
	ctx := context.Background()
	filter := bson.M{"user_code": userCode}
	user, err := u.user.FindOne(ctx, filter, bson.M{})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			u.logger.Warnf("[UnlinkDevice] user not found for userCode=%s", userCode)
			return fmt.Errorf(error_codes.AuthUserNotFound)
		}

		u.logger.Warnf("[UnlinkDevice] query failed for userCode=%s: %v", userCode, err)
		return fmt.Errorf(error_codes.UnhandledServerError)
	}

	pendingFilter := bson.M{
		"value":         userCode,
		"action_status": "PENDING",
	}

	existingPendingAction, err := u.actionRepo.FindOne(ctx, pendingFilter, bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		u.logger.Errorf("[UnlinkDevice] failed to check existing pending actions: %v", err)
		return fmt.Errorf(error_codes.PendingActionCheckFailed)
	}

	if existingPendingAction != nil {
		if existingPendingAction.Department == cpsAction.Department {
			u.logger.Warnf("[UnlinkDevice] pending request already exists in same department: %s", userCode)
			return fmt.Errorf(error_codes.PendingRequestExists)
		}

		updateFilter := bson.M{
			"action_code":   userCode,
			"action_status": "PENDING",
		}
		update := bson.M{
			"action_status": "REJECTED",
		}
		_, err := u.actionRepo.UpdateOne(ctx, updateFilter, update)
		if err != nil {
			u.logger.Errorf("[UnlinkDevice] failed to reject old pending action for user %s: %v", userCode, err)
			return fmt.Errorf(error_codes.PendingActionRejectionFailed)
		}
	}

	cpsAction.ActionCode = userCode
	cpsAction.CreatedAt = time.Now()
	cpsAction.LastModifiedAt = time.Now()

	cpsAction.CurrentAction = map[string]interface{}{
		"action_type": "UNLINK_DEVICE",
		"user_code":   userCode,
		"user_id":     user.ID,
		"full_name":   user.FullName,
	}

	if _, err := u.actionRepo.InsertOne(ctx, cpsAction); err != nil {
		u.logger.Errorf("[UnlinkDevice] failed to insert unlink action for userCode=%s: %v", userCode, err)
		return fmt.Errorf(error_codes.UnlinkActionFailed)
	}

	return nil
}

func (u *UnlinkRepo) ApproveOrDecline(userCode, decision, reason string, cpsAction entities.CPSAction) error {
	u.logger.Infof("[ApproveOrDecline] decision: %s for user: %s by %s", decision, userCode, cpsAction.CheckerName)
	ctx := context.Background()
	filter := bson.M{"user_code": userCode}
	user, err := u.user.FindOne(ctx, filter, bson.M{})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			u.logger.Warnf("[ApproveOrDecline] user not found for userCode=%s", userCode)
			return fmt.Errorf(error_codes.AuthUserNotFound)
		}
		u.logger.Warnf("[ApproveOrDecline] failed to find user: %v", err)
		return fmt.Errorf(error_codes.UnhandledServerError)
	}

	pendingFilter := bson.M{
		"action_code":    userCode,
		"action_status":  "PENDING",
		"request_action": "UNLINK_DEVICE",
	}
	action, err := u.actionRepo.FindOne(ctx, pendingFilter, bson.M{})
	if err != nil || action == nil {
		u.logger.Warnf("[ApproveOrDecline] unlink action not found or processed for userCode=%s", userCode)
		return fmt.Errorf(error_codes.ActionNotFound)
	}

	switch decision {
	case "AUTHORIZED":
		updateAction := bson.M{
			"action_status":        "APPROVED",
			"checker_name":         cpsAction.CheckerName,
			"checker_id":           cpsAction.CheckerID,
			"checker_phone_number": cpsAction.CheckerPhoneNumber,
			"checker_action_time":  time.Now(),
		}
		if _, err := u.actionRepo.UpdateOne(ctx, pendingFilter, updateAction); err != nil {
			u.logger.Errorf("[ApproveOrDecline] failed to approve action: %v", err)
			return fmt.Errorf(error_codes.ActionApprovalFailed)
		}

		updateUser := bson.M{
			"device.device_uuid":    "",
			"device_status":         "UNLINKED",
			"login_pin.pin":         "reset_pin_requested",
			"bps_reject_status":     "AUTHORIZED",
			"login_pin.pin_history": user.LoginPIN.PIN,
		}

		if _, err = u.user.UpdateOne(ctx, filter, updateUser); err != nil {
			u.logger.Errorf("[ApproveOrDecline] failed to update user info: %v", err)
			return fmt.Errorf(error_codes.ActionApprovalFailed)
		}

		u.logger.Infof("[ApproveOrDecline] unlink approved for user: %s", userCode)
		return nil

	case "DENIED":
		updateAction := bson.M{
			"action_status":        "REJECTED",
			"checker_name":         cpsAction.CheckerName,
			"checker_id":           cpsAction.CheckerID,
			"checker_phone_number": cpsAction.CheckerPhoneNumber,
			"checker_action_time":  time.Now(),
			"rejection_reason":     reason,
		}
		if _, err := u.actionRepo.UpdateOne(ctx, pendingFilter, updateAction); err != nil {
			u.logger.Errorf("[ApproveOrDecline] failed to reject action: %v", err)
			return fmt.Errorf(error_codes.ActionRejectionFailed)
		}

		if _, err = u.user.UpdateOne(ctx, filter, bson.M{"bps_reject_status": "DENIED"}); err != nil {
			u.logger.Errorf("[ApproveOrDecline] failed to update user status: %v", err)
			return fmt.Errorf(error_codes.UserStatusUpdateFailed)
		}

		u.logger.Infof("[ApproveOrDecline] unlink denied for user: %s", userCode)
		return nil

	default:
		u.logger.Warnf("[ApproveOrDecline] invalid decision: %s", decision)
		return fmt.Errorf(error_codes.InvalidDecison)
	}
}