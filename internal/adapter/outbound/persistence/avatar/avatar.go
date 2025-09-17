package avatar

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	// "net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/avatar"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AvatarPersistence struct {
	avatarDal dal.MongoDal[model.AvatarDocument, model.AvatarDocument]
	logger    utils.Logger
}

func InitAvatarPersistence(client *mongo.Client, dbName string, collection string, logger utils.Logger) outbound.AvatarRepository {
	return &AvatarPersistence{
		avatarDal: dal.NewMongoDal[model.AvatarDocument, model.AvatarDocument](client, dbName, collection),
		logger:    logger,
	}
}

func (a *AvatarPersistence) CreateAvatar(ctx context.Context, avatar avatar.Avatar) (*avatar.Avatar, error) {
	avatarDoc, err := mappers.ToAvatarDocument(avatar)
	if err != nil {
		a.logger.Errorf("Failed to convert avatar to document: %v", err)
		return nil, err
	}

	res, err := a.avatarDal.InsertOne(ctx, *avatarDoc)
	if err != nil {
		a.logger.Errorf("Failed to insert avatar: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBInsertFailed)
	}

	result := mappers.ToAvatarModel(&res)
	return &result, nil
}

func (a *AvatarPersistence) UpdateAvatar(ctx context.Context, avatar avatar.Avatar) (*avatar.Avatar, error) {
	objID, err := common_util.ParsePrimitiveObjectID(avatar.ID)
	if err != nil {
		a.logger.Errorf("Invalid avatar ID: %v", err)
		return nil, fmt.Errorf(common_util.InvalidID)
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{}
	if avatar.Label != "" {
		update["label"] = avatar.Label
	}
	if avatar.Avatar != "" {
		update["avatar"] = avatar.Avatar
	}
	update["last_modified_at"] = time.Now()

	if len(update) == 1 {
		a.logger.Warnf("No data provided for update on avatar ID %s", avatar.ID)
		return nil, fmt.Errorf(common_util.NoDataProvidedForUpdate)
	}

	avatarDoc, err := a.avatarDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("Avatar with ID %s not found", avatar.ID)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		a.logger.Errorf("Failed to update avatar ID %s: %v", avatar.ID, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToAvatarModel(&avatarDoc)

	return &result, nil
}

func (a *AvatarPersistence) DeleteAvatar(ctx context.Context, id string) (*avatar.Avatar, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		a.logger.Errorf("Invalid avatar ID: %v", err)
		return nil, fmt.Errorf(common_util.InvalidID)
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}
	update := bson.M{
		"is_deleted":       true,
		"deleted_at":       time.Now(),
		"last_modified_at": time.Now(),
	}

	avatarDoc, err := a.avatarDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("Avatar with ID %s not found", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		a.logger.Errorf("Failed to delete avatar ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToAvatarModel(&avatarDoc)

	return &result, nil
}

func (a *AvatarPersistence) EnableDisableAvatar(ctx context.Context, id string, enable bool) (*avatar.Avatar, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		a.logger.Errorf("Invalid avatar ID: %v", err)
		return nil, fmt.Errorf(common_util.InvalidID)
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{
		"enable":           enable,
		"last_modified_at": time.Now(),
	}

	avatarDoc, err := a.avatarDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("Avatar with ID %s not found for enable/disable", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		a.logger.Errorf("Failed to update enabled state for avatar ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToAvatarModel(&avatarDoc)

	return &result, nil
}

func (a *AvatarPersistence) GetAvatar(ctx context.Context, id string) (*avatar.Avatar, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		a.logger.Errorf("Invalid avatar ID: %v", err)
		return nil, fmt.Errorf(common_util.InvalidID)
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	avatarDoc, err := a.avatarDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("Avatar with ID %s not found", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		a.logger.Errorf("Failed to get avatar ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	result := mappers.ToAvatarModel(avatarDoc)

	return &result, nil
}

func (a *AvatarPersistence) GetAvatarByLabel(ctx context.Context, label string) (*avatar.Avatar, error) {
	filter := bson.M{
		"label":      bson.M{"$regex": "^" + regexp.QuoteMeta(label) + "$", "$options": "i"},
		"is_deleted": false,
	}

	avatarDoc, err := a.avatarDal.FindOne(ctx, filter, nil)
	if err != nil {
		a.logger.Errorf("failed to get avatar by label %s: %v", label, err)
		return nil, err
	}
	result := mappers.ToAvatarModel(avatarDoc)

	return &result, nil
}

func (a *AvatarPersistence) GetAllAvatar(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*avatar.Avatar], error) {
	filter := bson.M{"is_deleted": false}

	if filterParams != nil {
		if filterParams.Filters != "" {
			filter["enable"] = filterParams.Filters == "true"
		}
		if filterParams.Search != "" {
			searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
			filter["$or"] = []bson.M{
				{"avatar": searchRegex},
				{"label": searchRegex},
			}
		}
	}

	page := int64(1)
	limit := int64(10)
	if filterParams != nil {
		page = int64(filterParams.Page)
		limit = int64(filterParams.PerPage)
	}
	skip := (page - 1) * limit

	avatarDocs, err := a.avatarDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		a.logger.Errorf("Failed to get avatars: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	total, err := a.avatarDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("Failed to get avatar count: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	avatars := make([]*avatar.Avatar, 0, len(avatarDocs))
	for _, doc := range avatarDocs {
		converted := mappers.ToAvatarModel(doc)
		avatars = append(avatars, &converted)
	}

	meta := common_util.BuildPaginationMeta(total, filterParams.Page, int(limit))
	return &common_util.PaginatedResponse[[]*avatar.Avatar]{
		Data: avatars,
		Meta: meta,
	}, nil
}
