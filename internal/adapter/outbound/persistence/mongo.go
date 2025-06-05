package persistence
import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities"
	"cbe-super-app-member-users/internal/domain/users"
	"cbe-super-app-member-users/internal/port/outbound"
)

type MongoRepository struct {
	userDal          dal.MongoDal[entities.User, entities.User]
	linkedAccountDal dal.MongoDal[entities.LinkedAccount, entities.LinkedAccount]
	otpDal           dal.MongoDal[entities.OTP, entities.OTP]
}

func NewMongoRepository(client *mongo.Client, dbName string) outbound.UserRepository {
	return &MongoRepository{
		userDal:          dal.NewMongoDal[entities.User, entities.User](client, dbName, "users"),
		linkedAccountDal: dal.NewMongoDal[entities.LinkedAccount, entities.LinkedAccount](client, dbName, "linked_accounts"),
		otpDal: dal.NewMongoDal[entities.OTP, entities.OTP](client, dbName, "otps"),
		
	}
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*users.User, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	userFilter := bson.M{
		"_id":        oid,
		"is_deleted": bson.M{"$ne": true},
	}
	userEntity, err := r.userDal.FindOne(ctx, userFilter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	if userEntity == nil {
		return nil, mongo.ErrNoDocuments
	}

	return &users.User{
    ID:        userEntity.ID.Hex(),
    FullName:  userEntity.FullName, 
    IsDeleted: userEntity.IsDeleted,
}, nil
}

func (r *MongoRepository) FindActiveLinkedAccounts(ctx context.Context, userID string) ([]users.LinkedAccountDetail, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	accountFilter := bson.M{
		"user_id":          oid,
		"is_account_active": true,
		"linked_status":    true,
		"is_deleted":       bson.M{"$ne": true},
	}
	accountEntities, err := r.linkedAccountDal.FindAll(ctx, accountFilter, nil)
	if err != nil {
		return nil, err
	}

	linkedAccounts := make([]users.LinkedAccountDetail, 0, len(accountEntities))
	for _, acc := range accountEntities {
		linkedAccounts = append(linkedAccounts, users.LinkedAccountDetail{
			AccountNumber:     acc.AccountNumber,
			AccountBranchCode: acc.AccountBranchCode,
			LinkedBranch:      acc.LinkedBranch,
			IsAccountActive:   acc.IsAccountActive,
			LinkedStatus:      acc.LinkedStatus,
			CurrencyCode:      acc.CurrencyCode,
		})
	}

	return linkedAccounts, nil
}



func (r *MongoRepository) StoreOTP(ctx context.Context, otp *users.OTPRecord) error {
	// userID, err := bson.ObjectIDFromHex(otp.UserID)
	// if err != nil {
	// 	return err
	// }

	otpEntity := entities.OTP{
		UserRealm:   entities.MemberRealm,
		UserCode:    otp.UserID,
		OTPFor:      entities.OTPForChangeEmail,
		OTPCode:     otp.OTP,
		Email:       otp.Email,
		ExpiresAt:   otp.ExpiresAt,
		CreatedAt:   time.Now(),
		Status:      "",
		IsDeleted:   false,
	}

	_, err := r.otpDal.InsertOne(ctx, otpEntity)
	return err
}

func (r *MongoRepository) FindOTP(ctx context.Context, userID, email string) (*users.OTPRecord, error) {
	filter := bson.M{
		"user_code": userID,
		"email":     email,
		"otp_for":   entities.OTPForChangeEmail,
		"is_deleted": bson.M{"$ne": true},
		"status":    bson.M{"$ne": entities.Verified},
	}

	otpEntity, err := r.otpDal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	return &users.OTPRecord{
		UserID:    otpEntity.UserCode,
		Email:     otpEntity.Email,
		OTP:       otpEntity.OTPCode,
		CreatedAt: otpEntity.CreatedAt,
		ExpiresAt: otpEntity.ExpiresAt,
	}, nil
}

func (r *MongoRepository) UpdateUserEmail(ctx context.Context, userID, email string) error {
	oid, err :=bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"email":            email,
			"last_modified_at": time.Now(),
		},
	}

	_, err = r.userDal.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (*users.UserEmail, error) {
	filter := bson.M{
		"email":      email,
		"is_deleted": bson.M{"$ne": true},
	}

	userEntity, err := r.userDal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	return &users.UserEmail{
		ID: userEntity.ID.Hex(),
		Email: userEntity.Email,
	}, nil
}