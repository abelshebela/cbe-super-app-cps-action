package mappers

import (
	"errors"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/in-app-notification"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToInAppNotificationDocument(notification *entities.InAppNotification) (*model.InAppNotificationDocument, error) {
	objID := bson.NewObjectID()
	creatorObj, err := bson.ObjectIDFromHex(notification.CreatedBy)
	if err != nil {
		return nil, errors.New(error_codes.InvalidID)
	}
	return &model.InAppNotificationDocument{
		ID:                objID,
		NotificationType:  notification.NotificationType,
		NotificationBody:  notification.NotificationBody,
		IsPublic:          notification.IsPublic,
		For:               notification.For,
		CreatedBy:         creatorObj,
		NotificationParts: notification.NotificationParts,
		Seen:              false,
		Enabled:           true,
		IsDeleted:         false,
		CreatedAt:         time.Now(),
		LastModified:      time.Now(),
	}, nil
}
