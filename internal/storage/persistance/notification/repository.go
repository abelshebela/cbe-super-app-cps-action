package notification

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationStorage struct {
	dal           dal.MongoDal[model.Notification, model.Notification]
	client        *mongo.Client
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewNotificationRepository(client *mongo.Client, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.NotificationRepository {
	return &NotificationStorage{
		dal:           dal.NewMongoDal[model.Notification, model.Notification](client, dbName, collection),
		client:        client,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (n *NotificationStorage) Create(ctx context.Context, notification *model.Notification) error {
	// notification.IsFromCPS = true
	newNotification, err := n.dal.InsertOne(ctx, *notification)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	n.kafkaProducer.PublishMessage(ctx, newNotification, string(constants.ClientOrchestrationNotificationTopic), string(constants.ClientOrchestrationNotificationTopic), "new notification created")
	return nil
}

func (n *NotificationStorage) Update(ctx context.Context, id string, notification *model.Notification) error {
	n.logger.Infof("[Update] updating notification for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		n.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := NotificationMapper(*notification)

	updatedNotification, err := n.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			n.logger.Errorf("[Update] notification not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		n.logger.Errorf("[Update] failed to update notification: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	n.kafkaProducer.PublishMessage(ctx, updatedNotification, string(constants.ClientOrchestrationNotificationTopic), string(constants.ClientOrchestrationNotificationTopic), "notification updated")

	n.logger.Infof("[Update] notification updated successfully")
	return nil
}

func (n *NotificationStorage) Delete(ctx context.Context, id string) error {
	n.logger.Infof("[Delete] deleting notification for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		n.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = n.dal.DeleteOne(ctx, filter)
	if err != nil {
		n.logger.Errorf("[Delete] failed to delete notification: %v", err)
		return err
	}
	n.logger.Infof("[Delete] notification deleted successfully")
	return nil
}

func (n *NotificationStorage) FindByID(ctx context.Context, id string) (*model.Notification, error) {
	n.logger.Infof("[FindByID] fetching notification by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		n.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := n.dal.FindOne(ctx, filter, nil)
	if err != nil {
		n.logger.Errorf("[FindByID] failed to find notification: %v", err)
		return nil, err
	}
	n.logger.Infof("[FindByID] notification retrieved successfully")
	return result, nil
}

func (n *NotificationStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Notification], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"is_public", "notification_type", "notification_code", "for", "seen", "enabled", "title"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"title": searchRegex},
			{"notification_code": searchRegex},
			{"notification_body": searchRegex},
			{"notification_type": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	data, err := n.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		n.logger.Errorf("[FindAllWithPagination] failed to fetch notifications: %v", err)
		return nil, err
	}

	total, err := n.dal.TotalCount(ctx, filter)
	if err != nil {
		n.logger.Errorf("[FindAllWithPagination] failed to count notifications: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	n.logger.Infof("[FindAllWithPagination] retrieved %d notifications", len(data))

	return &types.PaginatedResponse[[]*model.Notification]{
		Data: data,
		Meta: meta,
	}, nil
}

func (n *NotificationStorage) NotificationExists(ctx context.Context, notificationType string, forValue constants.NotificationFor, id *string) (bool, error) {
	n.logger.Infof("[NotificationExists] checking notification existence")
	if notificationType == "" || forValue == "" {
		n.logger.Errorf("[NotificationExists] invalid input parameters")
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
			n.logger.Errorf("[NotificationExists] invalid object id: %v", err)
			return false, errors.New(localization.ErrorInvalidID.Code)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	count, err := n.dal.TotalCount(ctx, filter)
	if err != nil {
		n.logger.Errorf("[NotificationExists] failed to check notification existence: %v", err)
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}
	n.logger.Infof("[NotificationExists] notification existence check completed, count: %d", count)
	return count > 0, nil
}

func (n *NotificationStorage) EnableDisableNotification(ctx context.Context, id string, enable bool) (*model.Notification, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"enabled":          enable,
		"last_modified_at": time.Now(),
	}
	updatedNotification, err := n.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		n.logger.Errorf("EnableOrDisable Event failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	n.kafkaProducer.PublishMessage(ctx, updatedNotification, string(constants.ClientOrchestrationNotificationTopic), string(constants.ClientOrchestrationNotificationTopic), "notification enabled/disabled")

	return nil, nil
}
