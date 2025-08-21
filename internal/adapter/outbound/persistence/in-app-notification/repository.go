package inappnotification

import (
	"context"
	"errors"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/in-app-notification"
	inappnotification_outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/in-app-notification"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// NotificationPersistence implements the NotificationRepository interface
type InAppNotificationPersistence struct {
	inAppNotificationDal dal.MongoDal[model.InAppNotificationDocument, model.InAppNotificationDocument]
	logger               utils.Logger
}

func InitInAppNotificationPersistence(client *mongo.Client, dbName string, collection string, logger utils.Logger) inappnotification_outbound.InAppNotificationRepository {
	return &InAppNotificationPersistence{
		inAppNotificationDal: dal.NewMongoDal[model.InAppNotificationDocument, model.InAppNotificationDocument](client, dbName, collection),
		logger:               logger,
	}
}

func (n *InAppNotificationPersistence) CreateInAppNotification(ctx context.Context, notification *entities.InAppNotification) error {
	inAppNotificationDoc, err := mappers.ToInAppNotificationDocument(notification)
	if err != nil {
		n.logger.Errorf("Failed to map notification to document: %v", err)
		return err
	}

	_, err = n.inAppNotificationDal.InsertOne(ctx, *inAppNotificationDoc)
	if err != nil {
		n.logger.Errorf("Failed to insert notification: %v", err)
		return errors.New(error_codes.GeneralDBInsertFailed)
	}

	return nil
}

func (n *InAppNotificationPersistence) DeleteInAppNotification(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		n.logger.Errorf("Failed to convert id to object id: %v", err)
		return errors.New(error_codes.InvalidID)
	}

	err = n.inAppNotificationDal.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			n.logger.Errorf("Notification not found: %v", err)
			return errors.New(error_codes.NotFound)
		}
		n.logger.Errorf("Failed to delete notification: %v", err)
		return errors.New(error_codes.UnhandledServerError)
	}
	return nil
}

func (n *InAppNotificationPersistence) EnableDisableInAppNotification(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		n.logger.Errorf("Failed to convert id to object id: %v", err)
		return errors.New(error_codes.InvalidID)
	}

	_, err = n.inAppNotificationDal.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"enabled": enable})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			n.logger.Errorf("Notification not found: %v", err)
			return errors.New(error_codes.NotFound)
		}
		n.logger.Errorf("Failed to enable/disable notification: %v", err)
		return errors.New(error_codes.UnhandledServerError)
	}
	return nil
}
