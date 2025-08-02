package storage

import (
	"cbe-super-app-budget/internal/constants/dto"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type OTPRepository interface {
	Save(ctx context.Context, otp *dto.OTP) error
	Find(ctx context.Context) ([]dto.OTP, error)
	FindById(ctx context.Context, filter bson.M) (*dto.OTP, error)
	Update(ctx context.Context, update *dto.OTP) error
	Delete(ctx context.Context, id string) error
}

type UserRepository interface {
	Save(ctx context.Context, user *dto.User) error
	Find(ctx context.Context) ([]dto.User, error)
	FindById(ctx context.Context, filter bson.M) (*dto.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*dto.User, error)
	FindByUserCode(ctx context.Context, userCode string) (*dto.User, error)
	FindByDeviceUUID(ctx context.Context, deviceUUID string) (*dto.User, error)
	Update(ctx context.Context, update *dto.User) error
}

type DeviceLinkHistoryRepository interface {
	Save(ctx context.Context, deviceLinkHistory *dto.DeviceLinkHistory) error
	Find(ctx context.Context) ([]dto.DeviceLinkHistory, error)
	FindById(ctx context.Context, filter bson.M) (*dto.DeviceLinkHistory, error)
	Update(ctx context.Context, update *dto.DeviceLinkHistory) error
	Delete(ctx context.Context, id string) error
}
