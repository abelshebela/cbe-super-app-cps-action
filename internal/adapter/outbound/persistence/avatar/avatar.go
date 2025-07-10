package avatar

import (
	"context"
	"errors"
	"fmt"

	// "net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/avatar"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AvatarPersistence struct {
	cpsActionDal dal.MongoDal[model.CPSAction, model.CPSAction]
	avatarDal    dal.MongoDal[dto.Avatar, dto.Avatar]
	logger       utils.Logger
}

func InitAvatarPersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) avatar.AvatarOutbound {
	cpsActionDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collections[0])
	avatarDal := dal.NewMongoDal[dto.Avatar, dto.Avatar](client, dbName, collections[1])
	return &AvatarPersistence{
		cpsActionDal: cpsActionDal,
		avatarDal:    avatarDal,
		logger:       logger,
	}
}

func (a *AvatarPersistence) CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
		"request_action":          cpsReq.RequestAction,
	}

	projection := bson.M{
		"action_code": 1,
	}

	exists, err := a.cpsActionDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get avatar", err)
		return fmt.Errorf("FAILED_TO_GET_AVATER")
	} else if exists != nil {

		a.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return fmt.Errorf("PENDING_CPS_ACTION_PRESENT")
	}

	return nil
}

func (a *AvatarPersistence) CreateAvatar(ctx context.Context, cpsActionReq model.CreateCPSAction) (model.CPSAction, error) {
	cpsAction, err := a.cpsActionDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsActionReq.MakerUser.UserCode,
		MakerName:        cpsActionReq.MakerUser.FullName,
		MakerPhoneNumber: cpsActionReq.MakerUser.PhoneNumber,
		Department:       cpsActionReq.Department,
		ActionStatus:     string(model.ActionPending),
		RequestAction:    string(model.RequestCreateAvatar),
		ActionType:       string(model.ActionCreate),
		CurrentAction:    cpsActionReq.ActionData,
		MakerActionTime:  time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)

		return model.CPSAction{}, fmt.Errorf("FAILED_TO_CRATE_CPS_ACTION")
	}

	return cpsAction, nil
}

func (a *AvatarPersistence) DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CPSAction, error) {
	filter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	projection := bson.M{
		"label":      1,
		"avatar_url": 1,
		"is_deleted": 1,
	}

	avatar, err := a.avatarDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found %v", err)

			return model.CPSAction{}, fmt.Errorf("FAILED_TO_GET_AVATER")
		}
		a.logger.Errorf("failed to get avatar", err)

		return model.CPSAction{}, fmt.Errorf("NO_AVATER_DATA_FOUND")
	}

	cpsAction, err := a.cpsActionDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsActionReq.MakerUser.UserCode,
		MakerName:        cpsActionReq.MakerUser.FullName,
		MakerPhoneNumber: cpsActionReq.MakerUser.PhoneNumber,
		Department:       cpsActionReq.Department,
		ActionStatus:     string(model.ActionPending),
		RequestAction:    string(model.RequestDeleteAvatar),
		ActionType:       string(model.ActionDelete),
		CurrentAction:    cpsActionReq.ActionData,
		PreviosAction: map[string]any{
			"label":      avatar.Label,
			"avatar":     avatar.Avatar,
			"is_deleted": avatar.IsDeleted,
		},
		MakerActionTime: time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		return model.CPSAction{}, fmt.Errorf("FAILED_TO_CRATE_CPS_ACTION")
	}

	return cpsAction, nil
}

func (a *AvatarPersistence) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CPSAction, error) {
	var avatar dto.Avatar
	var err error

	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": model.ActionPending,
	}

	update := bson.M{

		"checker_id":           req.CheckerUser.UserCode,
		"checker_phone_number": req.CheckerUser.PhoneNumber,
		"checker_name":         req.CheckerUser.FullName,
		"action_status":        model.ActionApproved,
		"checker_action_time":  time.Now(),
	}

	cpsAction, err := a.cpsActionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("failed to update cps action", err)
		return model.CPSAction{}, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	var actionData dto.Avatar
	data, err := bson.Marshal(cpsAction.CurrentAction)
	if err != nil {
		a.logger.Errorf("failed to marshal bson: %v", err)

		return model.CPSAction{}, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	if err := bson.Unmarshal(data, &actionData); err != nil {
		a.logger.Errorf("failed to unmarshal into avatar: %v", err)

		return model.CPSAction{}, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	if cpsAction.ActionType == string(model.ActionCreate) {
		req := dto.Avatar{
			ID:        bson.NewObjectID().Hex(),
			Avatar:    actionData.Avatar,
			Label:     actionData.Label,
			Enable:    true,
			CreatedAt: time.Now(),
		}

		avatar, err = a.avatarDal.InsertOne(ctx, req)
		if err != nil {
			a.logger.Errorf("failed to create avatar", err)

			return model.CPSAction{}, fmt.Errorf("FAIL_TO_CREATE_AVATAR")
		}

		cpsAction.CurrentAction = avatar

		return cpsAction, nil

	}

	if cpsAction.ActionType == string(model.ActionUpdate) {
		filter := bson.M{
			"id":         actionData.ID,
			"is_deleted": false,
		}
		update := bson.M{}

		if actionData.Label != "" {
			update["label"] = actionData.Label
		}

		if actionData.Avatar != "" {
			update["avatar"] = actionData.Avatar
		}

		if cpsAction.RequestAction == string(model.RequestEnableAvatar) {
			update["enable"] = true
		}

		if cpsAction.RequestAction == string(model.RequestDisableAvatar) {
			update["enable"] = false
		}

		update["last_modified_at"] = time.Now()

		avatar, err = a.avatarDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("failed to update avatar", err)

			return model.CPSAction{}, fmt.Errorf("FAILED_TO_UPDATE_ADVERT")
		}

		cpsAction.CurrentAction = avatar

		return cpsAction, nil
	}

	if cpsAction.ActionType == string(model.ActionDelete) {
		filter := bson.M{
			"id":         actionData.ID,
			"is_deleted": false,
		}

		update := bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}

		avatar, err = a.avatarDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("failed to update avatar", err)

			return model.CPSAction{}, fmt.Errorf("FAILED_TO_UPDATE_ADVERT")
		}
		cpsAction.CurrentAction = avatar

		return cpsAction, nil
	}

	return cpsAction, nil
}

func (a *AvatarPersistence) Reject(ctx context.Context, req model.RejectCPSAction) (model.CPSAction, error) {
	filter := bson.M{
		"action_code": req.ActionCode,
		"department":  req.Department,
		"status":      model.ActionPending,
	}

	update := bson.M{
		"checker_user": bson.M{
			"full_name":    req.CheckerUser.FullName,
			"phone_number": req.CheckerUser.PhoneNumber,
			"user_code":    req.CheckerUser.UserCode,
		},
		"status":              model.ActionRejected,
		"rejected_reason":     req.RejectedReason,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := a.cpsActionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("failed to update avatar status", err)
		return model.CPSAction{}, fmt.Errorf("FAILED_TO_UPDATE_ADVERT")
	}
	return cpsAction, nil
}

func (a *AvatarPersistence) EnableOrDisableAvatar(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	avatarFilter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	avatarProjection := bson.M{
		"avatar": 1,
		"label":  1,
		"enable": 1,
	}

	avatar, err := a.avatarDal.FindOne(ctx, avatarFilter, avatarProjection)
	if err != nil {
		a.logger.Errorf("failed to get avatar", err)

		return nil, fmt.Errorf("FAILED_TO_UPDATE_AVATAR")
	}

	cps, err := a.cpsActionDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsReq.MakerUser.UserCode,
		MakerName:        cpsReq.MakerUser.UserCode,
		MakerPhoneNumber: cpsReq.MakerUser.PhoneNumber,
		Department:       cpsReq.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionUpdate),
		RequestAction:    string(requestAction),
		PreviosAction: map[string]any{
			"avatar": avatar.Avatar,
			"label":  avatar.Label,
			"enable": avatar.Enable,
		},
		CurrentAction:   cpsReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)

		return nil, fmt.Errorf("FAILED_TO_CRATE_CPS_ACTION")
	}

	return &cps, nil
}

func (a *AvatarPersistence) GetAllAvatar(ctx context.Context, filterParams constant.Filter) (dto.AvatarResponse, error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Filters != "" {
		filter["enable"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	avatar, err := a.avatarDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found")

			return dto.AvatarResponse{}, fmt.Errorf("NO_AVATER_DATA_FOUND")
		}
		a.logger.Errorf("failed to get avatar", err)

		return dto.AvatarResponse{}, fmt.Errorf("FAILED_TO_GET_AVATER")
	}

	total, err := a.avatarDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("failed to get avatar total counts", err)
		return dto.AvatarResponse{}, fmt.Errorf("FAILED_TO_GET_COUNT")
	}

	return dto.AvatarResponse{
		Page:    filterParams.Page,
		Avatars: avatar,
		Limit:   constant.DefaultPerPage,
		Total:   total,
	}, nil
}

func (a *AvatarPersistence) GetAvatar(ctx context.Context, id string) (*dto.Avatar, error) {
	filter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	projection := bson.M{}

	avatar, err := a.avatarDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found")

			return nil, fmt.Errorf("NO_AVATER_DATA_FOUND")
		}
		a.logger.Errorf("failed to get avatar", err)

		return nil, fmt.Errorf("FAILED_TO_GET_AVATER")
	}
	return avatar, nil
}

func (a *AvatarPersistence) UpdateAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CPSAction, error) {
	filter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	projection := bson.M{
		"avatar": 1,
	}

	avatar, err := a.avatarDal.FindOne(ctx, filter, projection)
	if err != nil {
		a.logger.Errorf("failed to get avatar", err)

		return model.CPSAction{}, fmt.Errorf("FAILED_TO_GET_AVATER")
	}

	cps, err := a.cpsActionDal.InsertOne(ctx, model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsActionReq.MakerUser.UserCode,
		MakerName:        cpsActionReq.MakerUser.FullName,
		MakerPhoneNumber: cpsActionReq.MakerUser.PhoneNumber,
		Department:       cpsActionReq.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionUpdate),
		RequestAction:    string(model.RequestUpdateAvatar),
		PreviosAction: map[string]any{
			"avatar": avatar.Avatar,
		},
		CurrentAction:   cpsActionReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)

		return model.CPSAction{}, fmt.Errorf("FAILED_TO_CRATE_CPS_ACTION")
	}

	return cps, nil
}
