package avatar

import (
	"context"
	"errors"
	"fmt"

	// "net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/avatar"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AvatarPersistence struct {
	cpsActionDal dal.MongoDal[model.CPSAction, model.CPSAction]
	avatarDal    dal.MongoDal[AvatarDocument, AvatarDocument]
	logger       utils.Logger
}

func InitAvatarPersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) avatar.AvatarOutbound {
	cpsActionDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collections[0])
	avatarDal := dal.NewMongoDal[AvatarDocument, AvatarDocument](client, dbName, collections[1])
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

func (a *AvatarPersistence) CreateAvatar(ctx context.Context, cpsActionReq model.CreateCPSAction) (*dto.CPSAction, error) {
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
		LastModifiedAt:   time.Now(),
		CreatedAt:        time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)

		return nil, fmt.Errorf("FAILED_TO_CRATE_CPS_ACTION")
	}

	return ToCPSAction(&cpsAction), nil
}

func (a *AvatarPersistence) DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (*dto.CPSAction, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("invalid object id: %v", err)
		return nil, fmt.Errorf("INVALID_OBJECT_ID")
	}
	filter := bson.M{
		"_id":        objectID,
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
			return nil, fmt.Errorf("NOT_FOUND")
		}
		a.logger.Errorf("failed to get avatar", err)
		return nil, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
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
		PreviousAction: map[string]any{
			"label":      avatar.Label,
			"avatar":     avatar.Avatar,
			"is_deleted": avatar.IsDeleted,
		},
		MakerActionTime: time.Now(),
		LastModifiedAt:  time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf("FAILED_TO_CRATE_CPS_ACTION")
	}

	return ToCPSAction(&cpsAction), nil
}

func (a *AvatarPersistence) Authorize(ctx context.Context, cpsAction *entities.CPSAction) (bool, error) {

	var err error

	var actionData dto.Avatar
	data, err := bson.Marshal(cpsAction.CurrentAction)
	if err != nil {
		a.logger.Errorf("failed to marshal bson: %v", err)

		return false, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	if err := bson.Unmarshal(data, &actionData); err != nil {
		a.logger.Errorf("failed to unmarshal into avatar: %v", err)

		return false, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	if string(cpsAction.ActionType) == string(model.ActionCreate) {
		req, err := ToAvatarDocument(dto.Avatar{
			ID:             bson.NewObjectID().Hex(),
			Avatar:         actionData.Avatar,
			Label:          actionData.Label,
			IsDeleted:      false,
			Enable:         true,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		})

		if err != nil {
			a.logger.Errorf("failed to convert avatar to document", err)
			return false, fmt.Errorf("FAILED_TO_CONVERT_AVATAR_TO_DOCUMENT")
		}

		avatarDoc, err := a.avatarDal.InsertOne(ctx, *req)
		if err != nil {
			a.logger.Errorf("failed to create avatar", err)

			return false, fmt.Errorf("FAIL_TO_CREATE_AVATAR")
		}

		cpsAction.CurrentAction = avatarDoc.toModel()

		return true, nil

	}

	if string(cpsAction.ActionType) == string(model.ActionUpdate) {
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

		if string(cpsAction.RequestAction) == string(model.RequestEnableAvatar) {
			update["enable"] = true
		}

		if string(cpsAction.RequestAction) == string(model.RequestDisableAvatar) {
			update["enable"] = false
		}

		update["last_modified_at"] = time.Now()

		avatarDoc, err := a.avatarDal.UpdateOne(ctx, filter, update)
		if err != nil {
			a.logger.Errorf("failed to update avatar", err)
			if errors.Is(err, mongo.ErrNoDocuments) {
				a.logger.Errorf("avatar not found", err)
				return false, fmt.Errorf("NOT_FOUND")
			}

			return false, fmt.Errorf("FAILED_TO_UPDATE_AVATAR")
		}

		cpsAction.CurrentAction = avatarDoc.toModel()

		return true, nil
	}

	if string(cpsAction.ActionType) == string(model.ActionDelete) {
		filter := bson.M{
			"id":         actionData.ID,
			"is_deleted": false,
		}

		update := bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}

		avatarDoc, err := a.avatarDal.UpdateOne(ctx, filter, update)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				a.logger.Errorf("avatar not found", err)
				return false, fmt.Errorf("NOT_FOUND")
			}
			a.logger.Errorf("failed to update avatar", err)

			return false, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
		}
		cpsAction.CurrentAction = avatarDoc

		return true, nil
	}

	return true, nil
}

func (a *AvatarPersistence) Reject(ctx context.Context, req model.RejectCPSAction) (*dto.CPSAction, error) {
	checkerUser := contexts.ExtractContext(ctx)
	filter := bson.M{
		"action_code":   req.ActionCode,
		"department":    req.Department,
		"action_status": model.ActionPending,
	}

	update := bson.M{
		"checker_id":           checkerUser.UserCode,
		"checker_name":         checkerUser.FullName,
		"checker_phone_number": checkerUser.PhoneNumber,
		"action_status":        model.ActionRejected,
		"rejection_reason":     req.RejectedReason,
		"checker_action_time":  time.Now(),
	}

	cpsAction, err := a.cpsActionDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("cps action not found", err)
			return nil, fmt.Errorf("ACTION_NOT_FOUND")
		}
		a.logger.Errorf("failed to update avatar status", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}
	return ToCPSAction(&cpsAction), nil
}

func (a *AvatarPersistence) EnableOrDisableAvatar(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*dto.CPSAction, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("invalid object id: %v", err)
		return nil, fmt.Errorf("INVALID_OBJECT_ID")
	}
	avatarFilter := bson.M{
		"_id":        objectID,
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found", err)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		return nil, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
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
		PreviousAction: map[string]any{
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

	return ToCPSAction(&cps), nil
}

func (a *AvatarPersistence) GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*dto.Avatar], error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Filters != "" {
		filter["enable"] = filterParams.Filters
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (filterParams.Page - 1) * filterParams.PerPage

	avatarDocs, err := a.avatarDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		a.logger.Errorf("failed to get avatar", err)
		return nil, fmt.Errorf("FAILED_TO_GET_AVATER")
	}

	total, err := a.avatarDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("failed to get avatar total counts", err)
		return nil, fmt.Errorf("FAILED_TO_GET_COUNT")
	}

	var avatars []*dto.Avatar
	for _, doc := range avatarDocs {
		n := doc.toModel()
		avatars = append(avatars, &n)
	}

	meta := common_util.BuildPaginationMeta(total, page, limit)
	return &common_util.PaginatedResponse[[]*dto.Avatar]{
		Data: avatars,
		Meta: meta,
	}, nil
}

func (a *AvatarPersistence) GetAvatar(ctx context.Context, id string) (*dto.Avatar, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("invalid object id: %v", err)
		return nil, fmt.Errorf("INVALID_ID")
	}
	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	projection := bson.M{}

	avatarDoc, err := a.avatarDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found")
			return nil, fmt.Errorf("NOT_FOUND")
		}
		a.logger.Errorf("failed to get avatar", err)

		return nil, fmt.Errorf("FAILED_TO_GET_AVATER")
	}
	avatar := avatarDoc.toModel()
	return &avatar, nil
}

func (a *AvatarPersistence) UpdateAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (*dto.CPSAction, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("invalid object id: %v", err)
		return nil, fmt.Errorf("INVALID_ID")
	}
	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	projection := bson.M{
		"avatar": 1,
	}

	avatar, err := a.avatarDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found", err)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		a.logger.Errorf("failed to get avatar", err)

		return nil, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
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
		PreviousAction: map[string]any{
			"avatar": avatar.Avatar,
		},
		CurrentAction:   cpsActionReq.ActionData,
		MakerActionTime: time.Now(),
	})
	if err != nil {
		a.logger.Errorf("failed to create cps action", err)

		return nil, fmt.Errorf("FAILED_TO_CRATE_CPS_ACTION")
	}

	return ToCPSAction(&cps), nil
}
