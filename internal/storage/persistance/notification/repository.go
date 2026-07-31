package notification

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"
	"fmt"
	"time"

	local_model "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	notification_dto "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/dto"
	shared_notification "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/dto"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationStorage struct {
	dal             dal.MongoDal[local_model.NotificationDocument, local_model.NotificationDocument]
	notificationDal dal.MongoDal[local_model.NotificationDocument, local_model.NotificationDocument]
	client          *mongo.Client
	kafkaProducer   kafka.NotificationProducer
	cfg             *config.VaultConfig
	logger          utils.Logger
}

func NewNotificationRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.NotificationProducer, logger utils.Logger) storage.NotificationRepository {
	return &NotificationStorage{
		dal:             dal.NewMongoDal[local_model.NotificationDocument, local_model.NotificationDocument](client, cfg, dbName, collection),
		notificationDal: dal.NewMongoDal[local_model.NotificationDocument, local_model.NotificationDocument](client, cfg, dbName, collection),
		client:          client,
		kafkaProducer:   kafkaProducer,
		cfg:             cfg,
		logger:          logger,
	}
}

func toSharedNotificationMessage(notification *local_model.NotificationDocument) shared_notification.BroadcastInAppNotificationMessage {
	if notification == nil {
		return shared_notification.BroadcastInAppNotificationMessage{}
	}

	return shared_notification.BroadcastInAppNotificationMessage{
		Title:         notification.Title,
		Message:       notification.NotificationBody,
		Category:      notification.Category,
		BroadcastType: notification.For,
	}
}

func (n *NotificationStorage) Create(ctx context.Context, notification *local_model.NotificationDocument) error {
	log := local_util.LoggerFromCtx(ctx, n.logger)
	if notification == nil {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}

	data := *notification
	if data.NotificationCode == "" {
		data.NotificationCode = time.Now().Format("20060102150405")
	}
	if data.CreatedAt.IsZero() {
		data.CreatedAt = time.Now()
	}

	newNotification, err := n.dal.InsertOne(ctx, data)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	publishData := toSharedNotificationMessage(&newNotification)

	err = n.kafkaProducer.PublishMessage(ctx, publishData, newNotification.For, n.cfg.KafkaInAppBordcastTopic, "inapp-notifications")
	if err != nil {
		log.Errorf("[NotificationStorage][Create] failed to send in app notification %v", err)
	}
	log.Infof("[NotificationStorage][Create] notification sent successfully")

	return nil
}

func (n *NotificationStorage) Update(ctx context.Context, id string, notification *local_model.NotificationDocument) error {
	log := local_util.LoggerFromCtx(ctx, n.logger)
	if notification == nil {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}

	log.Infof("[NotificationStorage][Update] updating notification for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[NotificationStorage][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	data := *notification
	updateData := NotificationMapper(data)

	updatedNotification, err := n.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	publishData := toSharedNotificationMessage(&updatedNotification)

	err = n.kafkaProducer.PublishMessage(ctx, publishData, publishData.BroadcastType, n.cfg.KafkaInAppTopic, "notification updated")
	if err != nil {
		log.Errorf("[NotificationStorage][Update] failed to update notification %v", err)
	}
	log.Infof("[NotificationStorage][Update] notification updated successfully")
	return nil
}

func (n *NotificationStorage) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	log.Infof("[NotificationStorage][Delete] deleting notification for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[NotificationStorage][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	err = n.dal.DeleteOneH(ctx, filter)
	if err != nil {
		log.Errorf("[NotificationStorage][Delete] failed to delete notification: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[NotificationStorage][Delete] notification deleted successfully")
	return nil
}

func (n *NotificationStorage) FindByID(ctx context.Context, id string) (*local_model.NotificationDocument, error) {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	log.Infof("[NotificationStorage][FindByID] fetching notification by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[NotificationStorage][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := n.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[NotificationStorage][FindByID] failed to find notification: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	log.Infof("[NotificationStorage][FindByID] notification retrieved successfully")
	return result, nil
}

func (n *NotificationStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.NotificationDocument], error) {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	searchKeys := bson.M{}
	allowedKeys := []string{"search", "title", "message", "for", "type"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"title": searchRegex},
			{"message": searchRegex},
			{"type": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := n.notificationDal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		log.Errorf("[NotificationStorage][FindAllWithPagination] failed to fetch notifications: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := n.notificationDal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[NotificationStorage][FindAllWithPagination] failed to count notifications: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	log.Infof("[NotificationStorage][FindAllWithPagination] retrieved %d notifications", len(data))

	return &types.PaginatedResponse[[]local_model.NotificationDocument]{
		Data: data,
		Meta: meta,
	}, nil
}

func (n *NotificationStorage) NotificationExists(ctx context.Context, notificationType string, forValue constants.NotificationFor, id *string) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	log.Infof("[NotificationStorage][NotificationExists] checking notification existence")
	if notificationType == "" || forValue == "" {
		log.Errorf("[NotificationStorage][NotificationExists] invalid input parameters")
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
			log.Errorf("[NotificationStorage][NotificationExists] invalid object id: %v", err)
			return false, errors.New(localization.ErrorInvalidID.Code)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	count, err := n.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[NotificationStorage][NotificationExists] failed to check notification existence: %v", err)
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}
	log.Infof("[NotificationStorage][NotificationExists] notification existence check completed, count: %d", count)
	return count > 0, nil
}

func (n *NotificationStorage) EnableDisableNotification(ctx context.Context, id string, enable bool) (*local_model.NotificationDocument, error) {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
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
		Type:    "OTHER",
		Title:   fmt.Sprintf("notification %s", status),
		Message: "notification has been " + status,
		Data:    NotificationMapper(updatedNotification),
	}

	err = n.kafkaProducer.PublishMessage(ctx, inAppMessage, "in_app_broadcast", n.cfg.KafkaInAppTopic, "notification enabled/disabled")
	if err != nil {
		log.Errorf("[NotificationStorage][EnableDisableNotification] failed to %s notification %v", status, err)
	}
	log.Infof("[NotificationStorage][EnableDisableNotification] notification sent successfully")

	// n.kafkaProducer.PublishMessage(ctx, updatedNotification, string(constants.ClientOrchestrationNotificationTopic), string(constants.ClientOrchestrationNotificationTopic), "notification enabled/disabled")

	return &updatedNotification, nil
}
