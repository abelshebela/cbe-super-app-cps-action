package device_history

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func DeviceLinkHistoryProjection() bson.M {
	return bson.M{
		"_id":               1,
		"user_id":           1,
		"device_type":       1,
		"device_uuid":       1,
		"device_name":       1,
		"device_os_version": 1,
		"linked_at":         1,
		"unlinked_at":       1,
	}
}

func DeviceLinkHistoryIdFilterAttachment(id string, filter bson.M) error {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrUnexpected
	}
	filter["_id"] = objId
	return nil
}

func DeviceLinkHistoryUserFilterAttachment(userID string, filter bson.M) error {
	objId, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return errors.ErrUnexpected
	}
	filter["user_id"] = objId
	return nil
}

func DeviceLinkHistoryDeviceUUIDFilterAttachment(deviceUUID string, filter bson.M) {
	filter["device_uuid"] = deviceUUID
}
