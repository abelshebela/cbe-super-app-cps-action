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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	notification_dto "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/dto"
	shared_producer "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/producer"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationStorage struct {
	dal                 dal.MongoDal[model.Notification, model.Notification]
	client              *mongo.Client
	kafkaProducer       kafka.ClientOrchestrationProducer
	sharedKafkaProducer *shared_producer.NotificationProducer
	logger              utils.Logger
}

func NewNotificationRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, sharedProducer *shared_producer.NotificationProducer, logger utils.Logger) storage.NotificationRepository {
	return &NotificationStorage{
		dal:                 dal.NewMongoDal[model.Notification, model.Notification](client, cfg, dbName, collection),
		client:              client,
		kafkaProducer:       kafkaProducer,
		sharedKafkaProducer: sharedProducer,
		logger:              logger,
	}
}

func (n *NotificationStorage) Create(ctx context.Context, notification *model.Notification) error {
	// notification.IsFromCPS = true
	newNotification, err := n.dal.InsertOne(ctx, *notification)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	inAppMessage := notification_dto.InAppKafkaMessage{
		Type:    "in_app",
		Title:   "notification created",
		Message: "new notification created",
		Data:    NotificationMapper(newNotification),
	}

	err = n.sharedKafkaProducer.PublishInAppMessage(ctx, inAppMessage)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][Create] failed to send in app notification %v", err)
	}
	n.logger.Infof("[NotificationStorage][Create] notification sent successfully")

	// n.kafkaProducer.PublishMessage(ctx, newNotification, "new_notification_created", "inapp-notifications", "new notification created")
	return nil
}

func (n *NotificationStorage) Update(ctx context.Context, id string, notification *model.Notification) error {
	n.logger.Infof("[NotificationStorage][Update] updating notification for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := NotificationMapper(*notification)

	updatedNotification, err := n.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return local_util.HandleDBError(err)
	}

	inAppMessage := notification_dto.InAppKafkaMessage{
		Type:    "in_app",
		Title:   "notification updated",
		Message: "notification has been updated",
		Data:    NotificationMapper(updatedNotification),
	}

	err = n.sharedKafkaProducer.PublishInAppMessage(ctx, inAppMessage)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][Update] failed to send notification %v", err)
	}
	n.logger.Infof("[NotificationStorage][Update] in-app notification sent successfully")

	// n.kafkaProducer.PublishMessage(ctx, updatedNotification, string(constants.ClientOrchestrationNotificationTopic), string(constants.ClientOrchestrationNotificationTopic), "notification updated")

	n.logger.Infof("[NotificationStorage][Update] notification updated successfully")
	return nil
}

func (n *NotificationStorage) Delete(ctx context.Context, id string) error {
	n.logger.Infof("[NotificationStorage][Delete] deleting notification for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = n.dal.DeleteOne(ctx, filter)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][Delete] failed to delete notification: %v", err)
		return local_util.HandleDBError(err)
	}
	n.logger.Infof("[NotificationStorage][Delete] notification deleted successfully")
	return nil
}

func (n *NotificationStorage) FindByID(ctx context.Context, id string) (*model.Notification, error) {
	n.logger.Infof("[NotificationStorage][FindByID] fetching notification by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := n.dal.FindOne(ctx, filter, nil)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][FindByID] failed to find notification: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	n.logger.Infof("[NotificationStorage][FindByID] notification retrieved successfully")
	return result, nil
}

func (n *NotificationStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Notification], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search", "is_public", "notification_type", "notification_code", "for", "seen", "enabled", "title"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"title": searchRegex},
			{"for": searchRegex},
			{"enabled": searchRegex},
			{"is_public": searchRegex},
			{"notification_code": searchRegex},
			{"notification_body": searchRegex},
			{"notification_type": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	data, err := n.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][FindAllWithPagination] failed to fetch notifications: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := n.dal.TotalCount(ctx, filter)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][FindAllWithPagination] failed to count notifications: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	n.logger.Infof("[NotificationStorage][FindAllWithPagination] retrieved %d notifications", len(data))

	return &types.PaginatedResponse[[]model.Notification]{
		Data: data,
		Meta: meta,
	}, nil
}

func (n *NotificationStorage) NotificationExists(ctx context.Context, notificationType string, forValue constants.NotificationFor, id *string) (bool, error) {
	n.logger.Infof("[NotificationStorage][NotificationExists] checking notification existence")
	if notificationType == "" || forValue == "" {
		n.logger.Errorf("[NotificationStorage][NotificationExists] invalid input parameters")
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
			n.logger.Errorf("[NotificationStorage][NotificationExists] invalid object id: %v", err)
			return false, errors.New(localization.ErrorInvalidID.Code)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	count, err := n.dal.TotalCount(ctx, filter)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][NotificationExists] failed to check notification existence: %v", err)
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}
	n.logger.Infof("[NotificationStorage][NotificationExists] notification existence check completed, count: %d", count)
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
		return nil, local_util.HandleDBError(err)
	}

	status := "disabled"
	if enable {
		status = "enabled"
	}

	inAppMessage := notification_dto.InAppKafkaMessage{
		Type:    "in_app",
		Title:   fmt.Sprintf("notification %s", status),
		Message: "notification has been " + status,
		Data:    NotificationMapper(updatedNotification),
	}

	err = n.sharedKafkaProducer.PublishInAppMessage(ctx, inAppMessage)
	if err != nil {
		n.logger.Errorf("[NotificationStorage][EnableDisableNotification] failed to %s notification %v", status, err)
	}
	n.logger.Infof("[NotificationStorage][EnableDisableNotification] notification sent successfully")

	// n.kafkaProducer.PublishMessage(ctx, updatedNotification, string(constants.ClientOrchestrationNotificationTopic), string(constants.ClientOrchestrationNotificationTopic), "notification enabled/disabled")

	return &updatedNotification, nil
}
