package otp

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OTPRepository struct {
	otpDal dal.MongoDal[model.OTP, model.OTP]
	logger utils.Logger
}

func NewOtpRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.OTPRepository {
	return &OTPRepository{
		otpDal: dal.NewMongoDal[model.OTP, model.OTP](client, dbName, collection),
		logger: logger,
	}
}

func (o *OTPRepository) Find(ctx context.Context, filter bson.M) (*model.OTP, error) {
	if filter != nil {
		o.logger.Errorf("Find OTP failed: filter is not nil. filter=%+v", filter)
		return nil, errors.ErrEmptyFilterParam
	}
	projection := OtpProjection()
	otp, err := o.otpDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == errors.ErrOTPNotFound {
			o.logger.Errorf("OTP not found. filter=%+v", filter)
			return nil, errors.ErrOTPNotFound
		}
		o.logger.Errorf("Unexpected error while finding OTP. error=%v, filter=%+v", err, filter)
		return nil, errors.ErrUnexpected
	}
	o.logger.Infof("OTP found successfully. otp=%+v", otp)
	return otp, nil
}

func (o *OTPRepository) Save(ctx context.Context, otp *model.OTP) error {
	if otp == nil {
		o.logger.Errorf("Save OTP failed: otp is nil")
		return errors.ErrInvalidOTP
	}

	_, err := o.otpDal.InsertOne(ctx, *otp)
	if err != nil {
		o.logger.Errorf("Unexpected error while saving OTP. error=%v, otp=%+v", err, otp)
		return errors.ErrUnexpected
	}
	o.logger.Infof("OTP saved successfully. otp=%+v", otp)
	return nil
}

func (o *OTPRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		o.logger.Errorf("Delete OTP failed: id is empty")
		return errors.ErrIdEmpty
	}
	filter, err := FilterIdForOtp(id)
	if err != nil {
		o.logger.Errorf("Delete OTP failed: error creating filter. error=%v, id=%s", err, id)
		return err
	}
	if err := o.otpDal.DeleteOne(ctx, filter); err != nil {
		o.logger.Errorf("Delete OTP failed: error deleting OTP. error=%v, id=%s", err, id)
		return err
	}
	o.logger.Infof("OTP deleted successfully. id=%s", id)
	return nil
}
