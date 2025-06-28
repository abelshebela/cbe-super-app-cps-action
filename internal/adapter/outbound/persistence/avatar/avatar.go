package avatar

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/avatar"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/avatar"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AvatarPersistence struct {
	cpsActionDal dal.MongoDal[model.CpsAction, model.CpsAction]
	avatarDal    dal.MongoDal[dto.Avatar, dto.Avatar]
	logger       utils.Logger
}

func InitAvatarPersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) avatar.AvatarOutbound {
	cpsActionDal := dal.NewMongoDal[model.CpsAction, model.CpsAction](client, dbName, collections[0])
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
		err = fmt.Errorf("failed to get avatar %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return err
	} else if exists != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		a.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return err
	}

	return nil
}

func (a *AvatarPersistence) CreateAvatar(ctx context.Context, cpsActionReq model.CreateCPSAction) (model.CpsAction, error) {
	cpsAction, err := a.cpsActionDal.InsertOne(ctx, model.CpsAction{
		ID:              bson.NewObjectID().Hex(),
		ActionCode:      utils.RandomGenerator(20),
		MakerUser:       cpsActionReq.MakerUser,
		Department:      cpsActionReq.Department,
		Status:          model.ActionPending,
		RequestAction:   model.RequestCreateAvatar,
		ActionType:      model.ActionCreate,
		ActionData:      cpsActionReq.ActionData,
		MakerActionTime: time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarPersistence) DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CpsAction, error) {
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
			err = fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "avatar not found",
			})
			return model.CpsAction{}, err
		}
		a.logger.Errorf("failed to get avatar", err)
		err = fmt.Errorf("failed to get avatar %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return model.CpsAction{}, err
	}

	cpsAction, err := a.cpsActionDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsActionReq.MakerUser,
		Department:    cpsActionReq.Department,
		Status:        model.ActionPending,
		RequestAction: model.RequestDeleteAvatar,
		ActionType:    model.ActionDelete,
		ActionData:    cpsActionReq.ActionData,
		PreviousData: map[string]any{
			"label":      avatar.Label,
			"avatar": avatar.Avatar,
			"is_deleted": avatar.IsDeleted,
		},
		CurrentData: map[string]any{
			"is_deleted": true,
		},
		MakerActionTime: time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarPersistence) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CpsAction, error) {
	var avatar dto.Avatar
	var err error

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
		"status":              model.ActionApproved,
		"checker_action_time": time.Now(),
	}

	cpsAction, err := a.cpsActionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("failed to update cps action", err)
		err = fmt.Errorf("failed to update cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return model.CpsAction{}, err
	}

	var actionData dto.Avatar
	data, err := bson.Marshal(cpsAction.ActionData)
	if err != nil {
		a.logger.Errorf("failed to marshal bson: %v", err)
		err = fmt.Errorf("failed to update cps action %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
		return model.CpsAction{}, err
	}

	if err := bson.Unmarshal(data, &actionData); err != nil {
		a.logger.Errorf("failed to unmarshal into avatar: %v", err)
		err = fmt.Errorf("failed to update cps action %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
		return model.CpsAction{}, err
	}

	if cpsAction.ActionType == model.ActionCreate {
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
			err = fmt.Errorf(" %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return model.CpsAction{}, err
		}

		cpsAction.ActionData = avatar

		return cpsAction, nil

	}

	if cpsAction.ActionType == model.ActionUpdate {
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

		if cpsAction.RequestAction == model.RequestEnableAvatar {
			update["enable"] = true
		}

		if cpsAction.RequestAction == model.RequestDisableAvatar {
			update["enable"] = false
		}

		update["last_modified_at"] = time.Now()

		avatar, err = a.avatarDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("failed to update avatar", err)
			err = fmt.Errorf("failed to update avatar %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return model.CpsAction{}, err
		}

		cpsAction.ActionData = avatar

		return cpsAction, nil
	}

	if cpsAction.ActionType == model.ActionDelete {
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
			err = fmt.Errorf("failed to update avatar %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
			return model.CpsAction{}, err
		}
		cpsAction.ActionData = avatar

		return cpsAction, nil
	}

	return cpsAction, nil
}

func (a *AvatarPersistence) Reject(ctx context.Context, req model.RejectCPSAction) (model.CpsAction, error) {
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
		err = fmt.Errorf("failed to update avatar status %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return model.CpsAction{}, err
	}
	return cpsAction, nil
}

func (a *AvatarPersistence) EnableOrDisableAvatar(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
	avatarFilter := bson.M{
		"id":         id,
		"is_deleted": false,
	}

	avatarProjection := bson.M{
		"avatar": 1,
		"label":      1,
		"enable":     1,
	}

	avatar, err := a.avatarDal.FindOne(ctx, avatarFilter, avatarProjection)
	if err != nil {
		a.logger.Errorf("failed to get avatar", err)
		err = fmt.Errorf("failed to get avatar %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	cps, err := a.cpsActionDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsReq.MakerUser,
		Department:    cpsReq.Department,
		Status:        model.ActionPending,
		ActionType:    model.ActionUpdate,
		ActionData:    cpsReq.ActionData,
		RequestAction: requestAction,
		PreviousData: map[string]any{
			"avatar": avatar.Avatar,
			"label":  avatar.Label,
			"enable": avatar.Enable,
		},
		CurrentData:     cpsReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cps, nil
}

func (a *AvatarPersistence) GetAllAvatar(ctx context.Context, filterParams constant.Filter) (dto.AvatarResponse, error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Filters != "" {
		filter["status"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	avatar, err := a.avatarDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found")
			err = fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "avatar not found",
			})
			return dto.AvatarResponse{}, err
		}
		a.logger.Errorf("failed to get avatar", err)
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return dto.AvatarResponse{}, err
	}

	total, err := a.avatarDal.TotalCount(ctx, bson.M{})
	if err != nil {
		a.logger.Errorf("failed to get avatar total counts", err)
		err := fmt.Errorf("failed to get avatar total counts %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return dto.AvatarResponse{}, err
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
			err = fmt.Errorf("%w", constant.ErrorDefinition{
				Code:    http.StatusNotFound,
				Message: "avatar not found",
			})
			return nil, err
		}
		a.logger.Errorf("failed to get avatar", err)
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}
	return avatar, nil
}

func (a *AvatarPersistence) UpdateAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CpsAction, error) {
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
		err = fmt.Errorf("failed to get avatar %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return model.CpsAction{}, err
	}

	cps, err := a.cpsActionDal.InsertOne(ctx, model.CpsAction{
		ID:            bson.NewObjectID().Hex(),
		ActionCode:    utils.RandomGenerator(20),
		MakerUser:     cpsActionReq.MakerUser,
		Department:    cpsActionReq.Department,
		Status:        model.ActionPending,
		ActionType:    model.ActionUpdate,
		ActionData:    cpsActionReq.ActionData,
		RequestAction: model.RequestUpdateAvatar,
		PreviousData: map[string]any{
			"avatar": avatar.Avatar,
		},
		CurrentData:     cpsActionReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return model.CpsAction{}, err
	}

	return cps, nil
}
