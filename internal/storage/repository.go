package storage

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"go.mongodb.org/mongo-driver/bson"
)

type OTPRepository interface {
	Save(ctx context.Context, otp *dto.OTP) error
	Find(ctx context.Context, filte bson.M) (*dto.OTP, error)
	Delete(ctx context.Context, id string) error
}

type UserRepository interface {
	Save(ctx context.Context, user *dto.User) error
	FindById(ctx context.Context, id string) (*dto.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*dto.User, error)
	FindByDeviceUUID(ctx context.Context, deviceUUID string) (*dto.User, error)
	Update(ctx context.Context, id string, update *dto.User) error
}

type DeviceLinkHistoryRepository interface {
	Save(ctx context.Context, deviceLinkHistory *dto.DeviceLinkHistroy) error
	Find(ctx context.Context) ([]dto.DeviceLinkHistroy, error)
	FindById(ctx context.Context, filter bson.M) (*dto.DeviceLinkHistroy, error)
	Update(ctx context.Context, update *dto.DeviceLinkHistroy) error
	Delete(ctx context.Context, id string) error
}
