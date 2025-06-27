package unlink

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/unlink"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UnlinkRepo struct {
	client     *mongo.Client
	user       dal.MongoDal[member.User, member.User]
	actionRepo dal.MongoDal[bps.BPSAction, bps.BPSAction]
	otpRepo    dal.MongoDal[member.OTP, member.OTP]
	logger     utils.Logger
}

var _ outbound.UnlinkRepository = (*UnlinkRepo)(nil)

func NewUnlinkInfrastructure(client *mongo.Client, dbName string, collectionName []string, logger utils.Logger) *UnlinkRepo {
	user := dal.NewMongoDal[member.User, member.User](client, dbName, collectionName[0])
	otpRepo := dal.NewMongoDal[member.OTP, member.OTP](client, dbName, collectionName[1])
	actionRepo := dal.NewMongoDal[bps.BPSAction, bps.BPSAction](client, dbName, collectionName[2])

	return &UnlinkRepo{
		user:       user,
		actionRepo: actionRepo,
		otpRepo:    otpRepo,
		logger:     logger,
	}
}

func hasCommonElement(a, b []string) bool {
	for _, v := range a {
		for _, w := range b {
			if v == w {
				return true
			}
		}
	}
	return false
}

func (u *UnlinkRepo) UnlinkDevice(userCode, makerUser string, branchCode []string, homeBranch string) error {
	u.logger.Infof("[UnlinkDevice] initiated for user: %s by %s", userCode, makerUser)
	ctx := context.Background()
	filter := bson.M{"user_code": userCode}
	user, err := u.user.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			u.logger.Warnf("[UnlinkDevice] user not found for userCode=%s", userCode)
			return fmt.Errorf("user not found")
		}
		u.logger.Warnf("[UnlinkDevice] query failed for userCode=%s: %v", userCode, err)
		return fmt.Errorf("Internal server error")
	}

	pendingFilter := bson.M{
		"value":  userCode,
		"status": "PENDING",
	}

	existingPendingAction, err := u.actionRepo.FindOne(ctx, pendingFilter, bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		u.logger.Errorf("[UnlinkDevice] failed to check existing pending actions: %v", err)
		return fmt.Errorf("failed to check pending actions")
	}

	if existingPendingAction != nil && len(existingPendingAction.BranchCode) > 0 {
		if hasCommonElement([]string{existingPendingAction.BranchCode}, branchCode) {
			u.logger.Warnf("[UnlinkDevice] pending request already exists in same branch: %s", userCode)
			return fmt.Errorf("pending request already exists in same branch")
		}

		updateFilter := bson.M{
			"value":  userCode,
			"status": "PENDING",
		}
		update := bson.M{
			"status": "REJECTED",
		}
		_, err := u.actionRepo.UpdateOne(ctx, updateFilter, update)
		if err != nil {
			u.logger.Errorf("[UnlinkDevice] failed to reject old pending action for user %s: %v", userCode, err)
			return fmt.Errorf("failed to reject pending action")
		}
	}

	now := time.Now()
	action := bps.BPSAction{
		ActionCode:       userCode,
		RequestAction:    "UNLINK_DEVICE",
		Status:           "PENDING",
		HomeBranch:       homeBranch,
		BranchCode:       strings.Join(branchCode, ","),
		CreatedAt:        now,
		LastModifiedAt:   now,
		EntityIdentifyer: userCode,
		MakerID:          makerUser,
		ActionReason:     bps.BPSAction{}.ActionReason,
	}
	action.ActionReason.ActionType = "UNLINK_DEVICE"
	action.ActionReason.ActionNote = "Requested unlink action"
	action.UserInformation.UserCode = userCode
	action.UserInformation.UserID = user.ID
	action.UserInformation.FullName = user.FullName

	if _, err := u.actionRepo.InsertOne(ctx, action); err != nil {
		u.logger.Errorf("[UnlinkDevice] failed to insert unlink action for userCode=%s: %v", userCode, err)
		return fmt.Errorf("failed to request unlink action")
	}

	if _, err := u.user.UpdateOne(ctx, filter, bson.M{"bps_reject_status": "PENDING"}); err != nil {
		u.logger.Errorf("[UnlinkDevice] failed to update user bps_reject_status for userCode=%s: %v", userCode, err)
		return fmt.Errorf("failed to request action")
	}

	return nil
}

func (u *UnlinkRepo) ApproveOrDecline(userCode, decision, reason, checkerUser string) error {
	u.logger.Infof("[ApproveOrDecline] decision: %s for user: %s by %s", decision, userCode, checkerUser)
	ctx := context.Background()
	filter := bson.M{"user_code": userCode}
	user, err := u.user.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			u.logger.Warnf("[ApproveOrDecline] user not found for userCode=%s", userCode)
			return fmt.Errorf("user not found")
		}
		u.logger.Warnf("[ApproveOrDecline] failed to find user: %v", err)
		return fmt.Errorf("Internal server error")
	}

	pendingFilter := bson.M{
		"value":          userCode,
		"status":         "PENDING",
		"request_action": "UNLINK_DEVICE",
	}
	action, err := u.actionRepo.FindOne(ctx, pendingFilter, bson.M{})
	if err != nil || action == nil {
		u.logger.Warnf("[ApproveOrDecline] unlink action not found or processed for userCode=%s", userCode)
		return fmt.Errorf("unlink action not found")
	}

	switch decision {
	case "AUTHORIZED":
		updateAction := bson.M{
			"status":       "APPROVED",
			"checker_user": checkerUser,
			"time":         time.Now(),
		}
		if _, err := u.actionRepo.UpdateOne(ctx, pendingFilter, updateAction); err != nil {
			u.logger.Errorf("[ApproveOrDecline] failed to approve action: %v", err)
			return fmt.Errorf("failed to approve action")
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
			return fmt.Errorf("failed to approve action")
		}

		u.logger.Infof("[ApproveOrDecline] unlink approved for user: %s", userCode)
		return nil

	case "DENIED":
		updateAction := bson.M{
			"status":       "REJECTED",
			"checker_user": checkerUser,
			"reason":       reason,
		}
		if _, err := u.actionRepo.UpdateOne(ctx, pendingFilter, updateAction); err != nil {
			u.logger.Errorf("[ApproveOrDecline] failed to reject action: %v", err)
			return fmt.Errorf("failed to reject action")
		}

		if _, err = u.user.UpdateOne(ctx, filter, bson.M{"bps_reject_status": "DENIED"}); err != nil {
			u.logger.Errorf("[ApproveOrDecline] failed to update user status: %v", err)
			return fmt.Errorf("failed to update user status")
		}

		u.logger.Infof("[ApproveOrDecline] unlink denied for user: %s", userCode)
		return nil

	default:
		u.logger.Warnf("[ApproveOrDecline] invalid decision: %s", decision)
		return fmt.Errorf("invalid decision")
	}
}
