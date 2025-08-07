package device_history

import (
	"context"

	"cbe-super-app-member-auth/internal/constants/errors"
	"cbe-super-app-member-auth/internal/constants/model"
	"cbe-super-app-member-auth/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DeviceLinkHistoryRepository struct {
	deviceHistoryDal dal.MongoDal[model.DeviceLinkHistroy, model.DeviceLinkHistroy]
	logger           utils.Logger
}

func NewDeviceLinkHistoryRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.DeviceLinkHistoryRepository {
	return &DeviceLinkHistoryRepository{
		deviceHistoryDal: dal.NewMongoDal[model.DeviceLinkHistroy, model.DeviceLinkHistroy](client, dbName, collection),
		logger:           logger,
	}
}

func (d *DeviceLinkHistoryRepository) Save(ctx context.Context, deviceLinkHistory *model.DeviceLinkHistroy) error {
	if deviceLinkHistory == nil {
		d.logger.Errorf("Save device link history failed: deviceLinkHistory is nil")
		return errors.ErrUnexpected
	}

	_, err := d.deviceHistoryDal.InsertOne(ctx, *deviceLinkHistory)
	if err != nil {
		d.logger.Errorf("Unexpected error while saving device link history. error=%v, deviceLinkHistory=%+v", err, deviceLinkHistory)
		return errors.ErrUnexpected
	}
	d.logger.Infof("Device link history saved successfully. deviceLinkHistory=%+v", deviceLinkHistory)
	return nil
}

func (d *DeviceLinkHistoryRepository) Update(ctx context.Context, update *model.DeviceLinkHistroy) error {
	if update == nil {
		d.logger.Errorf("Update device link history failed: update is nil")
		return errors.ErrUnexpected
	}

	filter := bson.M{"_id": update.ID}
	updateDoc := bson.M{
		"$set": bson.M{
			"user_id":           update.UserID,
			"device_type":       update.DeviceType,
			"device_uuid":       update.DeviceUUID,
			"device_name":       update.DeviceName,
			"device_os_version": update.DeviceOSVersion,
			"linked_at":         update.LinkedAt,
			"unlinked_at":       update.UnLinkedAt,
		},
	}

	_, err := d.deviceHistoryDal.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		d.logger.Errorf("Unexpected error while updating device link history. error=%v, update=%+v", err, update)
		return errors.ErrUnexpected
	}
	d.logger.Infof("Device link history updated successfully. update=%+v", update)
	return nil
}
