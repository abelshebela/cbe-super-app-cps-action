package mappers

import (
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToNotificationDomain(doc *model.NotificationDocument) *entities.Notification {
	return &entities.Notification{
		ID:                doc.ID.Hex(),
		NotificationType:  doc.NotificationType,
		NotificationBody:  doc.NotificationBody,
		IsPublic:          doc.IsPublic,
		For:               entities.NotificationFor(doc.For),
		CreatedBy:         doc.CreatedBy,
		NotificationParts: doc.NotificationParts,
		Seen:              doc.Seen,
		Title:             doc.Title,
		Status:            entities.NotificationStatus(strings.ToUpper(doc.Status)),
		Enabled:           doc.Enabled,
		IsDeleted:         doc.IsDeleted,
		CreatedAt:         doc.CreatedAt,
		LastModified:      doc.LastModified,
	}
}

func ToNotificationDocument(domain *entities.Notification) (*model.NotificationDocument, error) {
	var ID bson.ObjectID

	if domain.ID != "" {
		id, err := bson.ObjectIDFromHex(domain.ID)
		if err != nil {
			return nil, err
		}

		ID = id
	} else {
		ID = bson.NewObjectID()
	}

	return &model.NotificationDocument{
		ID:                ID,
		NotificationType:  domain.NotificationType,
		NotificationBody:  domain.NotificationBody,
		IsPublic:          domain.IsPublic,
		For:               string(domain.For),
		CreatedBy:         domain.CreatedBy,
		NotificationParts: domain.NotificationParts,
		Seen:              domain.Seen,
		Title:             domain.Title,
		Status:            string(domain.Status),
		Enabled:           domain.Enabled,
		IsDeleted:         domain.IsDeleted,
		CreatedAt:         domain.CreatedAt,
		LastModified:      domain.LastModified,
	}, nil
}
