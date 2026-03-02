package otp

import (
	"context"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OTPRepository struct {
	otpDal dal.MongoDal[model.OTP, model.OTP]
	logger utils.Logger
}

func NewOtpRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.OTPRepository {
	return &OTPRepository{
		otpDal: dal.NewMongoDal[model.OTP, model.OTP](client, cfg, dbName, collection),
		logger: logger,
	}
}

func (o *OTPRepository) Find(ctx context.Context, filter bson.M) (*model.OTP, error) {

	projection := OtpProjection()
	otp, err := o.otpDal.FindOne(ctx, filter, projection)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	o.logger.Infof("OTP found successfully. otp=%+v", otp)
	return otp, nil
}

func (o *OTPRepository) Save(ctx context.Context, otp *model.OTP) error {

	_, err := o.otpDal.InsertOne(ctx, *otp)
	if err != nil {
		o.logger.Errorf("Unexpected error while saving OTP. error=%v, otp=%+v", err, otp)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	o.logger.Infof("OTP saved successfully. otp=%+v", otp)
	return nil
}

func (o *OTPRepository) Delete(ctx context.Context, id string) error {

	filter, err := local_util.FilterIdFor(id)
	if err != nil {
		o.logger.Errorf("Delete OTP failed: error creating filter. error=%v, id=%s", err, id)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	if err := o.otpDal.DeleteOne(ctx, filter); err != nil {
		o.logger.Errorf("Delete OTP failed: error deleting OTP. error=%v, id=%s", err, id)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	o.logger.Infof("OTP deleted successfully. id=%s", id)
	return nil
}
