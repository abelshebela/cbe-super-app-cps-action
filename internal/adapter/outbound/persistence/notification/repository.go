package notification

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	notification_outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/notification"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// NotificationPersistence implements the NotificationRepository interface
type NotificationPersistence struct {
	notificationDal dal.MongoDal[model.NotificationDocument, model.NotificationDocument]
	logger          utils.Logger
	client          *mongo.Client
}

// InitNotificationPersistence initializes the notification persistence layer
func InitNotificationPersistence(client *mongo.Client, dbName string, collection string, logger utils.Logger) notification_outbound.NotificationRepository {
	return &NotificationPersistence{
		notificationDal: dal.NewMongoDal[model.NotificationDocument, model.NotificationDocument](client, dbName, collection),
		logger:          logger,
		client:          client,
	}
}

// CreateNotification creates a new notification in the database
func (n *NotificationPersistence) CreateNotification(ctx context.Context, notification *entities.Notification) (*entities.Notification, error) {
	notificationDoc, err := mappers.ToNotificationDocument(notification)
	if err != nil {
		n.logger.Errorf("Failed to map notification to document: %v", err)
		return nil, err
	}

	res, err := n.notificationDal.InsertOne(ctx, *notificationDoc)
	if err != nil {
		n.logger.Errorf("Failed to insert notification: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBInsertFailed)
	}

	result := mappers.ToNotificationDomain(&res)
	return result, nil
}

// FetchNotificationByID fetches a notification by its ID
func (n *NotificationPersistence) FetchNotificationByID(ctx context.Context, id string) (*entities.Notification, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		n.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	notificationDoc, err := n.notificationDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			n.logger.Warnf("Notification with ID %s not found", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		n.logger.Errorf("Failed to fetch notification by ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	result := mappers.ToNotificationDomain(notificationDoc)
	return result, nil
}

// FetchNotifications fetches notifications with pagination and filtering
func (n *NotificationPersistence) FetchNotifications(ctx context.Context, filterParam *constant.MongoFilter) (*common_util.PaginatedResponse[[]*entities.Notification], error) {
	filter := bson.M{"is_deleted": false}

	fmt.Println("I get called")

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"notification_type": searchRegex},
			{"notification_body": searchRegex},
			{"for": searchRegex},
			{"title": searchRegex},
		}
	}

	if filterParam.Filters != nil {
		allowedKeys := []string{"enabled", "notification_type"}
		handlers := map[string]func(interface{}) interface{}{
			"is_deleted": func(value interface{}) interface{} {
				if str, ok := value.(string); ok {
					if parsed, err := strconv.ParseBool(str); err == nil {
						return parsed
					}
				}
				return value
			},
		}

		enhancedFilter := common_util.BuildMongoFilterWithHandlers(filterParam.Filters, allowedKeys, handlers)
		for key, value := range enhancedFilter {
			filter[key] = value
		}
	}

	skip := (filterParam.Page - 1) * filterParam.PerPage
	limit := filterParam.PerPage

	notificationDocs, err := n.notificationDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		n.logger.Errorf("Failed to fetch notifications: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	var notifications []*entities.Notification
	for _, doc := range notificationDocs {
		converted := mappers.ToNotificationDomain(doc)
		notifications = append(notifications, converted)
	}

	total, err := n.notificationDal.TotalCount(ctx, filter)
	if err != nil {
		n.logger.Errorf("Failed to count notifications: %v", err)
		return nil, err
	}

	meta := common_util.BuildPaginationMeta(total, filterParam.Page, limit)
	return &common_util.PaginatedResponse[[]*entities.Notification]{
		Data: notifications,
		Meta: meta,
	}, nil
}

// UpdateNotification updates an existing notification
func (n *NotificationPersistence) UpdateNotification(ctx context.Context, notification *entities.Notification) (*entities.Notification, error) {
	objID, err := common_util.ParsePrimitiveObjectID(notification.ID)
	if err != nil {
		n.logger.Errorf("Invalid ID format: %s, error: %v", notification.ID, err)
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{}
	if notification.NotificationType != "" {
		update["notification_type"] = notification.NotificationType
	}
	if notification.NotificationBody != "" {
		update["notification_body"] = notification.NotificationBody
	}
	if notification.IsPublic {
		update["is_public"] = notification.IsPublic
	}
	if notification.For != "" {
		update["for"] = string(notification.For)
	}

	if notification.Title != "" {
		update["title"] = notification.Title
	}
	if notification.CreatedBy != "" {
		update["created_by"] = notification.CreatedBy
	}
	if notification.NotificationParts != nil {
		update["notification_parts"] = notification.NotificationParts
	}
	update["last_modified"] = time.Now()

	if len(update) == 1 {
		n.logger.Warnf("No data provided for update for notification ID %s", notification.ID)
		return nil, fmt.Errorf(common_util.NoDataProvidedForUpdate)
	}

	notificationDoc, err := n.notificationDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			n.logger.Warnf("Notification with ID %s not found for update", notification.ID)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		n.logger.Errorf("Failed to update notification ID %s: %v", notification.ID, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToNotificationDomain(&notificationDoc)
	return result, nil
}

// DeleteNotification soft-deletes a notification
func (n *NotificationPersistence) DeleteNotification(ctx context.Context, id string) (*entities.Notification, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		n.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return nil, fmt.Errorf(common_util.InvalidID)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	notificationDoc, err := n.notificationDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			n.logger.Warnf("Notification with ID %s not found for deletion", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		n.logger.Errorf("Failed to delete notification ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToNotificationDomain(&notificationDoc)
	return result, nil
}

// EnableDisableNotification enables or disables a notification
func (n *NotificationPersistence) EnableDisableNotification(ctx context.Context, id string, enable bool) (*entities.Notification, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		n.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{
		"enabled":      enable,
		"lastModified": time.Now(),
	}

	notificationDoc, err := n.notificationDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			n.logger.Errorf("Notification with ID %s not found for enable/disable", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		n.logger.Errorf("Failed to update enabled state for notification ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToNotificationDomain(&notificationDoc)
	return result, nil
}

// NotificationExists checks if a notification with the given type and target exists
func (n *NotificationPersistence) NotificationExists(ctx context.Context, notificationType string, forValue entities.NotificationFor, id *string) (bool, error) {
	if notificationType == "" || forValue == "" {
		n.logger.Warnf("Invalid input for notification existence check: type=%s, for=%s", notificationType, forValue)
		return false, fmt.Errorf(common_util.InvalidInput)
	}

	filter := bson.M{
		"notificationType": bson.M{"$regex": fmt.Sprintf("^%s$", notificationType), "$options": "i"},
		"for":              bson.M{"$regex": fmt.Sprintf("^%s$", forValue), "$options": "i"},
		"is_deleted":       false,
	}

	if id != nil {
		objID, err := common_util.ParsePrimitiveObjectID(*id)
		if err != nil {
			n.logger.Errorf("Invalid ID format: %s, error: %v", *id, err)
			return false, fmt.Errorf(common_util.InvalidID)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	count, err := n.notificationDal.TotalCount(ctx, filter)
	if err != nil {
		n.logger.Errorf("Failed to check notification existence for type %s, for %s: %v", notificationType, forValue, err)
		return false, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	return count > 0, nil
}

func (p *NotificationPersistence) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	p.logger.Debugf("Starting MongoDB session for transaction")

	session, err := p.client.StartSession()
	if err != nil {
		p.logger.Errorf("failed to start MongoDB session: %v", err)
		return fmt.Errorf(common_util.UnhandledServerError)
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(txCtx context.Context) error {
		p.logger.Debugf("Starting MongoDB transaction")

		if err := session.StartTransaction(); err != nil {
			p.logger.Errorf("failed to start transaction: %v", err)
			return fmt.Errorf(common_util.UnhandledServerError)
		}

		err := fn(txCtx)
		if err != nil {
			p.logger.Errorf("transaction logic failed: %v", err)
			if abortErr := session.AbortTransaction(txCtx); abortErr != nil {
				p.logger.Errorf("failed to abort transaction: %v", abortErr)
			} else {
				p.logger.Debugf("Transaction aborted successfully")
			}
			return err
		}

		if err := session.CommitTransaction(txCtx); err != nil {
			p.logger.Errorf("failed to commit transaction: %v", err)
			return fmt.Errorf(common_util.UnhandledServerError)
		}

		p.logger.Debugf("Transaction committed successfully")
		return nil
	})
}
