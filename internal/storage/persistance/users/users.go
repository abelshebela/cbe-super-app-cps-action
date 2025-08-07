package users

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

type userRepository struct {
	userDal dal.MongoDal[model.User, model.User]
	logger  utils.Logger
}

func NewUserRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.UserRepository {
	return &userRepository{
		userDal: dal.NewMongoDal[model.User, model.User](client, dbName, collection),
		logger:  logger,
	}
}

func (r *userRepository) Save(ctx context.Context, user *model.User) error {
	if user == nil {
		r.logger.Errorf("attempted to save nil user")
		return errors.ErrTryToSaveEmptyUser
	}

	if _, err := r.userDal.InsertOne(ctx, *user); err != nil {
		r.logger.Errorf("failed to insert user")
		return errors.ErrUnexpected
	}
	r.logger.Infof("user saved successfully")
	return nil
}

func (r *userRepository) FindById(ctx context.Context, id string) (*model.User, error) {
	projection := UserProjection()
	filter, err := UserIdFilterAttachMent(id)
	if err != nil {
		r.logger.Errorf("invalid user id for FindById")
		return nil, err
	}

	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Warnf("no user found for the provided id")
			return nil, errors.ErrUserNotFound
		}
		r.logger.Errorf("unexpected error during FindById")
		return nil, errors.ErrUnexpected
	}

	r.logger.Infof("user found by id")
	return user, nil
}

func (r *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.User, error) {
	if phoneNumber == "" {
		r.logger.Errorf("phone number is empty in FindByPhoneNumber")
		return nil, errors.ErrPhoneNumberCanNotBeEmpty
	}

	projection := UserProjection()
	filter := UserPhoneFilterAttachment(phoneNumber)

	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("failed to find user by phone number", err)

			return nil, errors.ErrUserNotFound
		}
		r.logger.Errorf("failed to find user by phone number")
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	r.logger.Infof("user found by phone number")
	return user, nil
}

func (r *userRepository) FindByDeviceUUID(ctx context.Context, deviceUUID string) (*model.User, error) {
	if deviceUUID == "" {
		r.logger.Errorf("deviceUUID is empty in FindByDeviceUUID")
		return nil, errors.ErrDeviceUUIDCanNotBeNull
	}

	projection := UserProjection()

	filter := UserDeviceUUIDAttachment(deviceUUID)

	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Warnf("no user found for the provided deviceUUID")
			return nil, errors.ErrUserNotFound
		}
		r.logger.Errorf("unexpected error during FindByDeviceUUID: %v", err)
		return nil, errors.ErrUnexpected
	}

	r.logger.Infof("user found by deviceUUID")
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, id string, update *model.User) error {
	if update == nil {
		r.logger.Errorf("attempted to update with nil user")
		return errors.ErrTryToSaveEmptyUser
	}
	filter, _ := UserIdFilterAttachMent(id)
	req := UserBuilder(*update)

	_, err := r.userDal.UpdateOne(ctx, filter, req)
	if err != nil {
		r.logger.Errorf("failed to update user")
		return err
	}

	r.logger.Infof("user updated successfully")
	return nil
}

func (r *userRepository) UpdateLoginAttemp(ctx context.Context, id string, update bson.M) error {
	if update == nil {

		return errors.ErrEmptyEmptyData
	}

	filter, _ := UserIdFilterAttachMent(id)
	_, err := r.userDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
