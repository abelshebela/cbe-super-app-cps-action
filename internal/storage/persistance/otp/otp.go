package otp

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OTPRepository struct {
	otpDal dal.MongoDal[dto.OTP, dto.OTP]
	logger utils.Logger
}

func NewOtpRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.OTPRepository {
	return &OTPRepository{
		otpDal: dal.MongoDal[dto.OTP, dto.OTP],
		logger: utils.Logger,
	}
}

func (o *OTPRepository) Find(ctx context.Context, filter bson.M) (*dto.OTP, error) {
	if filter != nil {
		return nil, errors.ErrEmptyFilterParam
	}
	projection := OtpProjection()
	otp, err := o.otpDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == errors.ErrOTPNotFound {
			return nil, errors.ErrOTPNotFound
		}
		return nil, errors.ErrUnexpected
	}
	return otp, nil
}

func (o *OTPRepository) Save(ctx context.Context, otp *dto.OTP) error {
	if otp == nil {
		return errors.ErrInvalidOTP
	}

	_, err := o.otpDal.InsertOne(ctx, *otp)
	if err != nil {
		return errors.ErrUnexpected
	}
	return nil
}
func (o *OTPRepository) Delete(ctx context.Context, id string) error {

	if id == "" {
		return errors.ErrIdEmpty
	}
	filter, err := FilterIdForOtp(id)
	if err != nil {
		return err
	}
	if err := o.otpDal.DeleteOne(ctx, filter); err != nil {
		return err
	}
	return nil
}
