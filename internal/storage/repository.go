package storage

import (
	"context"

	"cbe-super-app-member-auth/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OTPRepository interface {
	Save(ctx context.Context, otp *model.OTP) error
	Find(ctx context.Context, filte bson.M) (*model.OTP, error)
	Delete(ctx context.Context, id string) error
}

type UserRepository interface {
	Save(ctx context.Context, user *model.User) (*model.User, error)
	FindById(ctx context.Context, id string) (*model.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.User, error)
	FindByDeviceUUID(ctx context.Context, deviceUUID string) (*model.User, error)
	Update(ctx context.Context, id string, update *model.User) error
}

type HQRepository interface {
	FindOne(ctx context.Context, filter bson.M) (*model.HQ, error)
}

type DeviceLinkHistoryRepository interface {
	Save(ctx context.Context, deviceLinkHistory *model.DeviceLinkHistroy) error
	Update(ctx context.Context, update *model.DeviceLinkHistroy) error
}

type ResetSessionRepository interface {
	Save(ctx context.Context, session *model.PinResetSession) error
	FindById(ctx context.Context, id string) (*model.PinResetSession, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.PinResetSession, error)
	FindByDeviceUUID(ctx context.Context, deviceUUID string) (*model.PinResetSession, error)
	Update(ctx context.Context, id string, update *model.PinResetSession) error
	Delete(ctx context.Context, id string) error
}
