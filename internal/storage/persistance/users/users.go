package users

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepository struct {
	userDal    dal.MongoDal[model.User, model.User]
	client     *mongo.Client
	logger     utils.Logger
	dbName     string
	collection string
}

func NewUserRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.UserRepository {
	return &userRepository{
		userDal:    dal.NewMongoDal[model.User, model.User](client, dbName, collection),
		logger:     logger,
		client:     client,
		dbName:     dbName,
		collection: collection,
	}
}

func (r *userRepository) Save(ctx context.Context, user *model.User) error {

	if _, err := r.userDal.InsertOne(ctx, *user); err != nil {
		r.logger.Errorf("failed to insert user")
		return errors.New(localization.ErrorInternalServerError.Code)
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
			return nil, errors.New(localization.ErrorUserNotFound.Code)
		}
		r.logger.Errorf("unexpected error during FindById")
		return nil, errors.New(localization.ErrorInternalServerError.Code)
	}

	r.logger.Infof("user found by id")
	return user, nil
}

func (r *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.User, error) {

	projection := UserProjection()
	filter := UserPhoneFilterAttachment(phoneNumber)

	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("failed to find user by phone number", err)

			return nil, errors.New(localization.ErrorUserNotFound.Code)
		}
		r.logger.Errorf("failed to find user by phone number")
		return nil, errors.New(localization.ErrorInternalServerError.Code)
	}

	r.logger.Infof("user found by phone number")
	return user, nil
}

func (r *userRepository) FindByDeviceUUID(ctx context.Context, deviceUUID string) (*model.User, error) {

	projection := UserProjection()

	filter := UserDeviceUUIDAttachment(deviceUUID)

	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Warnf("no user found for the provided deviceUUID")
			return nil, errors.New(localization.ErrorUserNotFound.Code)
		}
		r.logger.Errorf("unexpected error during FindByDeviceUUID: %v", err)
		return nil, errors.New(localization.ErrorInternalServerError.Code)
	}

	r.logger.Infof("user found by deviceUUID")
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, id string, update *model.User) error {

	filter, _ := UserIdFilterAttachMent(id)
	req := UserBuilder(*update)
	_, err := r.userDal.UpdateOne(ctx, filter, req)
	if err != nil {
		r.logger.Errorf("failed to update user")
		return errors.New(localization.ErrorInternalServerError.Code)
	}

	r.logger.Infof("user updated successfully")
	return nil
}

func (r *userRepository) FindByUserCode(ctx context.Context, userCode string) (*model.User, error) {
	filter := bson.M{
		"user_code": userCode,
	}
	projection := UserProjection()
	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorUserNotFound.Code)
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByCustomerNumber(ctx context.Context, customerNumber string) (*model.User, error) {
	filter := bson.M{
		"customer_number": customerNumber,
	}
	projection := UserProjection()
	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorUserNotFound.Code)
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByAccount(ctx context.Context, accNumber string) (*model.User, error) {
	linkedAccountCollection := r.client.Database(r.dbName).Collection("linked_accounts")
	linkedAccountFilter := bson.M{"account_number": accNumber}
	var linkedAccount struct {
		CustomerNumber string `bson:"customer_number"`
	}
	err := linkedAccountCollection.FindOne(ctx, linkedAccountFilter).Decode(&linkedAccount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Warnf("no linked account found for the provided account number")
			return nil, errors.New(localization.ErrorUserNotFound.Code)
		}
		r.logger.Errorf("unexpected error during GetUserByAccount (linked account): %v", err)
		return nil, errors.New(localization.ErrorInternalServerError.Code)
	}

	userFilter := bson.M{
		"customer_number": linkedAccount.CustomerNumber,
	}
	projection := UserProjection()
	user, err := r.userDal.FindOne(ctx, userFilter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Warnf("no user found for the provided customer number")
			return nil, errors.New(localization.ErrorUserNotFound.Code)
		}
		r.logger.Errorf("unexpected error during GetUserByAccount (user): %v", err)
		return nil, errors.New(localization.ErrorInternalServerError.Code)
	}

	r.logger.Infof("user found by account number via linked account")
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	db := r.client.Database(r.dbName)
	collection := db.Collection(r.collection)
	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Errorf("failed to hard delete document from %s: %v", r.collection, err)
		return errors.New(localization.ErrorInternalServerError.Code)
	}

	r.logger.Infof("Successfully hard deleted document from %s with id: %v", r.collection, id)
	return nil

}
