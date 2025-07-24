package avatar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	// "net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/avatar"
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

func (a *AvatarPersistence) GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*dto.Avatar], error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	// Apply enable filter if provided
	if filterParams.Filters != "" {
		filter["enable"] = filterParams.Filters == "true"
	}

	// Apply search filter on avatar and label fields
	if filterParams.Search != "" {
		filter["$or"] = []bson.M{
			{"avatar": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"label": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	avatarDocs, err := a.avatarDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		a.logger.Errorf("failed to get avatars: %v", err)
		return nil, fmt.Errorf("FAILED_TO_GET_AVATARS")
	}

	total, err := a.avatarDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("failed to get avatar count: %v", err)
		return nil, fmt.Errorf("FAILED_TO_GET_COUNT")
	}

	avatars := make([]*dto.Avatar, 0, len(avatarDocs))
	for _, doc := range avatarDocs {
		converted := doc.toModel()
		avatars = append(avatars, &converted)
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

func (a *AvatarPersistence) AuthorizeCreateAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	actionData, err := getAvatar(cpsAction)
	if err != nil {
		return nil, err
	}

	req, err := ToAvatarDocument(dto.Avatar{
		Avatar:         actionData.Avatar,
		Label:          actionData.Label,
		IsDeleted:      false,
		Enable:         true,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	})
	if err != nil {
		a.logger.Errorf("failed to convert avatar to document: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CONVERT_AVATAR_TO_DOCUMENT")
	}

	avatarDoc, err := a.avatarDal.InsertOne(ctx, *req)
	if err != nil {
		a.logger.Errorf("failed to create avatar: %v", err)
		return nil, fmt.Errorf("FAIL_TO_CREATE_AVATAR")
	}

	cpsAction.CurrentAction = avatarDoc.toModel()
	return cpsAction, nil
}

func (a *AvatarPersistence) AuthorizeUpdateAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	actionData, err := getAvatar(cpsAction)
	if err != nil {
		return nil, err
	}

	objId, err := bson.ObjectIDFromHex(actionData.ID)
	if err != nil {
		return nil, fmt.Errorf("INVALID_ID")
	}
	filter := bson.M{
		"_id":         objId,
		"is_deleted": false,
	}
	update := bson.M{}

	if actionData.Label != "" {
		update["label"] = actionData.Label
	}
	if actionData.Avatar != "" {
		update["avatar"] = actionData.Avatar
	}

	switch string(cpsAction.RequestAction) {
	case string(model.RequestEnableAvatar):
		update["enable"] = true
	case string(model.RequestDisableAvatar):
		update["enable"] = false
	}

	update["last_modified_at"] = time.Now()

	avatarDoc, err := a.avatarDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found: %v", err)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		a.logger.Errorf("failed to update avatar: %v", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_AVATAR")
	}

	cpsAction.CurrentAction = avatarDoc.toModel()
	return cpsAction, nil
}

func (a *AvatarPersistence) AuthorizeDeleteAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	actionData, err := getAvatar(cpsAction)
	if err != nil {
		return nil, err
	}

	objId, err := bson.ObjectIDFromHex(actionData.ID)
	if err != nil {
		return nil, fmt.Errorf("INVALID_ID")
	}

	filter := bson.M{
		"_id":         objId,
		"is_deleted": false,
	}
	update := bson.M{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}

	avatarDoc, err := a.avatarDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("avatar not found: %v", err)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		a.logger.Errorf("failed to delete avatar: %v", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	cpsAction.CurrentAction = avatarDoc
	return cpsAction, nil
}

func getAvatar(cpsAction *entities.CPSAction) (*dto.Avatar, error) {
	var actionData dto.Avatar

	bytes, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %w", err)
	}

	if err := json.Unmarshal(bytes, &actionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CurrentAction into Advert: %w", err)
	}

	return &actionData, nil
}
