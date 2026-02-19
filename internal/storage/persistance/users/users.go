package users

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepository struct {
	userDal       dal.MongoDal[member.User, member.User]
	client        *mongo.Client
	logger        utils.Logger
	dbName        string
	collection    string
	kafkaProducer kafka.ClientOrchestrationProducer
}

func NewUserRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, clientOrchestrationProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.UserRepository {
	return &userRepository{
		userDal:       dal.NewMongoDal[member.User, member.User](client, cfg, dbName, collection),
		logger:        logger,
		client:        client,
		dbName:        dbName,
		collection:    collection,
		kafkaProducer: clientOrchestrationProducer,
	}
}

func (r *userRepository) Save(ctx context.Context, user *member.User) error {

	if _, err := r.userDal.InsertOne(ctx, *user); err != nil {
		r.logger.Errorf("failed to insert user ")
		return errors.New(localization.ErrorInternalServerError.Code)
	}
	r.logger.Infof("user saved successfully")
	return nil
}

func (r *userRepository) FindById(ctx context.Context, id string) (*member.User, error) {
	// projection := UserProjection()
	filter, err := UserIdFilterAttachMent(id)
	if err != nil {
		r.logger.Errorf("invalid user id for FindById")
		return nil, err
	}

	user, err := r.userDal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	r.logger.Infof("user found by id")
	return user, nil
}

func (r *userRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*member.User, error) {

	projection := UserProjection()
	filter := UserPhoneFilterAttachment(phoneNumber)

	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	r.logger.Infof("user found by phone number")
	return user, nil
}

func (r *userRepository) FindByDeviceUUID(ctx context.Context, deviceUUID string) (*member.User, error) {

	projection := UserProjection()

	filter := UserDeviceUUIDAttachment(deviceUUID)

	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	r.logger.Infof("user found by deviceUUID")
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, id string, update *member.User) error {

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

func (r *userRepository) FindByUserCode(ctx context.Context, userCode string) (*member.User, error) {
	filter := bson.M{
		"user_code": userCode,
	}
	// projection := UserProjection()
	user, err := r.userDal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return user, nil
}

func (r *userRepository) FindByCustomerNumber(ctx context.Context, customerNumber string) (*member.User, error) {
	filter := bson.M{
		"customer_number": customerNumber,
	}
	projection := UserProjection()
	user, err := r.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return user, nil
}

func (r *userRepository) GetUserByAccount(ctx context.Context, accNumber string) (*member.User, error) {
	linkedAccountCollection := r.client.Database(r.dbName).Collection("linked_accounts")
	linkedAccountFilter := bson.M{"account_number": accNumber}
	var linkedAccount struct {
		CustomerNumber string `bson:"customer_number"`
	}
	err := linkedAccountCollection.FindOne(ctx, linkedAccountFilter).Decode(&linkedAccount)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	userFilter := bson.M{
		"customer_number": linkedAccount.CustomerNumber,
	}
	projection := UserProjection()
	user, err := r.userDal.FindOne(ctx, userFilter, projection)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	r.logger.Infof("user found by account number via linked account")
	return user, nil
}

func (r *userRepository) DeleteHard(ctx context.Context, id string) error {
	// objId, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	return localization.ErrorUnexpectedError
	// }
	filter := bson.M{
		"_id": id,
	}
	return r.userDal.DeleteOneH(ctx, filter)
}
func (r *userRepository) Delete(ctx context.Context, id string) error {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return localization.ErrorUnexpectedError
	}
	db := r.client.Database(r.dbName)
	collection := db.Collection(r.collection)
	filter := bson.M{"_id": objId}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Errorf("failed to hard delete document from %s: %v", r.collection, err)
		return errors.New(localization.ErrorInternalServerError.Code)
	}

	r.logger.Infof("Successfully hard deleted document from %s with id: %v", r.collection, id)
	return nil

}
