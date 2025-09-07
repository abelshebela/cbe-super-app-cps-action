package notification

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationStorage struct {
	dal    dal.MongoDal[model.Notification, model.Notification]
	client *mongo.Client
	logger utils.Logger
}

func NewNotificationRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.NotificationRepository {
	return &NotificationStorage{
		dal:    dal.NewMongoDal[model.Notification, model.Notification](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (n *NotificationStorage) Create(ctx context.Context, notification *model.Notification) error {
	_, err := n.dal.InsertOne(ctx, *notification)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (n *NotificationStorage) Update(ctx context.Context, id string, notification *model.Notification) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := NotificationMapper(*notification)

	_, err = n.dal.UpdateOne(ctx, filter, bson.M{"$set": updateData})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (n *NotificationStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return n.dal.DeleteOne(ctx, filter)
}

func (n *NotificationStorage) FindByID(ctx context.Context, id string) (*model.Notification, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := n.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (n *NotificationStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Notification], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"title": searchRegex},
			{"notification_body": searchRegex},
		}
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := n.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := n.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.Notification]{
		Data: data,
		Meta: meta,
	}, nil
}

func (n *NotificationStorage) NotificationExists(ctx context.Context, notificationType string, forValue constants.NotificationFor, id *string) (bool, error) {
	if notificationType == "" || forValue == "" {
		n.logger.Warnf("Invalid input for notification existence check: type=%s, for=%s", notificationType, forValue)
		return false, errors.New(localization.ErrorInvalidRequest.Code)
	}

	filter := bson.M{
		"notification_type": bson.M{"$regex": fmt.Sprintf("^%s$", notificationType), "$options": "i"},
		"for":               bson.M{"$regex": fmt.Sprintf("^%s$", forValue), "$options": "i"},
		"is_deleted":        false,
	}

	if id != nil {
		objID, err := bson.ObjectIDFromHex(*id)
		if err != nil {
			n.logger.Errorf("Invalid ID format: %s, error: %v", *id, err)
			return false, errors.New(localization.ErrorInvalidID.Code)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	count, err := n.dal.TotalCount(ctx, filter)
	if err != nil {
		n.logger.Errorf("Failed to check notification existence for type %s, for %s: %v", notificationType, forValue, err)
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return count > 0, nil
}

func (n *NotificationStorage) EnableDisableNotification(ctx context.Context, id string, enable bool) (*model.Notification, error) {
	objID, err := local_util.ParsePrimitiveObjectID(id)
	if err != nil {
		n.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorInvalidIDFormat.Code)
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{
		"enabled":       enable,
		"last_modified": time.Now(),
	}

	_, err = n.dal.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			n.logger.Errorf("Notification with ID %s not found for enable/disable", id)
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		n.logger.Errorf("Failed to update enabled state for notification ID %s: %v", id, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Fetch and return the updated document
	updated, err := n.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
