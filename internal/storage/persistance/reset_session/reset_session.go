package reset_session

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	local_util "github.com/CBE-Super-App/cbe-super-app-member-auth/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ResetSessionRepository struct {
	resetSessionDal dal.MongoDal[model.PinResetSession, model.PinResetSession]
	logger          utils.Logger
}

func NewResetSessionRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ResetSessionRepository {
	return &ResetSessionRepository{
		resetSessionDal: dal.NewMongoDal[model.PinResetSession, model.PinResetSession](client, dbName, collection),
		logger:          logger,
	}
}

func (r *ResetSessionRepository) Save(ctx context.Context, session *model.PinResetSession) error {
	if session == nil {
		r.logger.Errorf("Save reset session failed: session is nil")
		return errors.ErrUnexpected
	}

	_, err := r.resetSessionDal.InsertOne(ctx, *session)
	if err != nil {
		r.logger.Errorf("Unexpected error while saving reset session. error=%v, session=%+v", err, session)
		return errors.ErrUnexpected
	}
	r.logger.Infof("Reset session saved successfully. session=%+v", session)
	return nil
}

func (r *ResetSessionRepository) FindById(ctx context.Context, id string) (*model.PinResetSession, error) {
	if id == "" {
		r.logger.Errorf("FindById reset session failed: id is empty")
		return nil, errors.ErrIdEmpty
	}

	projection := ResetSessionProjection()
	filter := ResetSessionIdFilterAttachment(id)

	session, err := r.resetSessionDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Warnf("No reset session found for the provided id")
			return nil, errors.ErrNoResetSession
		}
		r.logger.Errorf("Unexpected error while finding reset session by id. error=%v, id=%s", err, id)
		return nil, errors.ErrUnexpected
	}
	r.logger.Infof("Reset session found by id successfully. id=%s", id)
	return session, nil
}

func (r *ResetSessionRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.PinResetSession, error) {
	if phoneNumber == "" {
		r.logger.Errorf("FindByPhoneNumber reset session failed: phoneNumber is empty")
		return nil, errors.ErrPhoneNumberCanNotBeEmpty
	}

	projection := ResetSessionProjection()
	filter := ResetSessionPhoneFilterAttachment(phoneNumber)

	session, err := r.resetSessionDal.FindOne(ctx, filter, projection)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			r.logger.Warnf("No reset session found for the provided phone number")
			return nil, errors.ErrNoResetSession
		}
		r.logger.Errorf("Unexpected error while finding reset session by phone number. error=%v, phoneNumber=%s", err, phoneNumber)
		return nil, errors.ErrUnexpected
	}
	r.logger.Infof("Reset session found by phone number successfully. phoneNumber=%s", phoneNumber)
	return session, nil
}

func (r *ResetSessionRepository) FindByDeviceUUID(ctx context.Context, deviceUUID string) (*model.PinResetSession, error) {
	if deviceUUID == "" {
		r.logger.Errorf("FindByDeviceUUID reset session failed: deviceUUID is empty")
		return nil, errors.ErrDeviceUUIDCanNotBeNull
	}

	projection := ResetSessionProjection()
	filter := bson.M{}
	ResetSessionDeviceUUIDFilterAttachment(deviceUUID, filter)

	session, err := r.resetSessionDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == errors.ErrNoMongoDocument {
			r.logger.Warnf("No reset session found for the provided deviceUUID")
			return nil, errors.ErrUnexpected
		}
		r.logger.Errorf("Unexpected error while finding reset session by deviceUUID. error=%v, deviceUUID=%s", err, deviceUUID)
		return nil, errors.ErrUnexpected
	}
	r.logger.Infof("Reset session found by deviceUUID successfully. deviceUUID=%s", deviceUUID)
	return session, nil
}

func (r *ResetSessionRepository) Update(ctx context.Context, id string, update *model.PinResetSession) error {
	if update == nil {
		r.logger.Errorf("Update reset session failed: update is nil")
		return errors.ErrUnexpected
	}

	filter := ResetSessionIdFilterAttachment(id)

	updateDoc := bson.M{
		"user_id":           update.UserID,
		"phone_number":      update.PhoneNumber,
		"device_uuid":       update.DeviceUUID,
		"expires_at":        update.ExpiresAt,
		"created_at":        update.CreatedAt,
		"attempts":          update.Attempts,
		"max_attempts":      update.MaxAttempts,
		"otp":               update.OTP,
		"otp_for":           update.OTPFor,
		"verified_at":       update.VerifiedAt,
		"completed_at":      update.CompletedAt,
		"enabled":           update.Enabled,
		"access_restricted": update.AccessRestricted,
		"status":            update.Status,
		"restrictions":      update.Restrictions,
	}

	_, err := r.resetSessionDal.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		r.logger.Errorf("Unexpected error while updating reset session. error=%v, id=%s", err, id)
		return errors.ErrUnexpected
	}
	r.logger.Infof("Reset session updated successfully. id=%s", id)
	return nil
}

func (r *ResetSessionRepository) Delete(ctx context.Context, id string) error {

	filter, err := local_util.FilterIdFor(id)

	if err != nil {
		r.logger.Errorf("Delete OTP failed: error creating filter. error=%v, id=%s", err, id)
		return err
	}
	if err := r.resetSessionDal.DeleteOne(ctx, filter); err != nil {
		r.logger.Errorf("Delete OTP failed: error deleting OTP. error=%v, id=%s", err, id)
		return err
	}
	r.logger.Infof("OTP deleted successfully. id=%s", id)
	return nil
}
