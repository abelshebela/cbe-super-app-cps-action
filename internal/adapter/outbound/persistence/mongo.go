package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	userPort "cbe-super-app-member-users/internal/domain/users"
	accountPort "cbe-super-app-member-users/internal/port/outbound/account"

	entities "cbe-super-app-member-users/internal/adapter/outbound/model"

	dal "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrNotFound = errors.New("not found")

type MongoRepository struct {
	userDal          dal.MongoDal[entities.User, entities.User]
	linkedAccountDal dal.MongoDal[entities.LinkedAccount, entities.LinkedAccount]
	otpDal           dal.MongoDal[entities.OTP, entities.OTP]
}

func NewMongoRepository(client *mongo.Client, dbName string) *MongoRepository {
	return &MongoRepository{
		userDal:          dal.NewMongoDal[entities.User, entities.User](client, dbName, "user"),
		linkedAccountDal: dal.NewMongoDal[entities.LinkedAccount, entities.LinkedAccount](client, dbName, "linked_accounts"),
		otpDal:           dal.NewMongoDal[entities.OTP, entities.OTP](client, dbName, "otp"),
	}
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*userPort.User, error) {
	// fmt.Println("hello there ")
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}

	userFilter := bson.M{
		"_id":        oid,
		"is_deleted": bson.M{"$ne": true},
	}
	userEntity, err := r.userDal.FindOne(ctx, userFilter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &userPort.User{
		ID:        userEntity.ID.Hex(),
		Email:     userEntity.Email,
		FullName:  userEntity.FullName,
		IsDeleted: userEntity.IsDeleted,
	}, nil
}

func (r *MongoRepository) FindActiveLinkedAccounts(ctx context.Context, userID string) ([]userPort.LinkedAccountDetail, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrNotFound
	}

	accountFilter := bson.M{
		"user_id":           oid,
		"is_account_active": true,
		"linked_status":     true,
	}
	accountEntities, err := r.linkedAccountDal.FindAll(ctx, accountFilter, nil)
	if err != nil {
		return nil, err
	}

	linkedAccounts := make([]userPort.LinkedAccountDetail, 0, len(accountEntities))
	for _, acc := range accountEntities {
		linkedAccounts = append(linkedAccounts, userPort.LinkedAccountDetail{
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

func (r *MongoRepository) StoreOTP(ctx context.Context, otp *userPort.OTPRecord) error {
	otpEntity := entities.OTP{
		UserCode:  otp.UserID,
		OTPFor:    entities.OTPForChangeEmail,
		OTPCode:   otp.OTP,
		Email:     otp.Email,
		ExpiresAt: otp.ExpiresAt,
		CreatedAt: time.Now(),
		Status:    "",
		IsDeleted: false,
	}

	_, err := r.otpDal.InsertOne(ctx, otpEntity)
	return err
}

func (r *MongoRepository) FindOTP(ctx context.Context, userID, email string) (*userPort.OTPRecord, error) {
	filter := bson.M{
		"user_code":  userID,
		"email":      email,
		"otp_for":    entities.OTPForChangeEmail,
		"is_deleted": bson.M{"$ne": true},
		"status":     bson.M{"$ne": entities.Verified},
		"expires_at": bson.M{"$gt": time.Now()},
	}

	otpEntity, err := r.otpDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &userPort.OTPRecord{
		ID:        otpEntity.ID.Hex(),
		UserCode:  otpEntity.UserCode,
		UserID:    otpEntity.UserCode,
		Email:     otpEntity.Email,
		OTP:       otpEntity.OTPCode,
		CreatedAt: otpEntity.CreatedAt,
		ExpiresAt: otpEntity.ExpiresAt,
	}, nil
}

func (r *MongoRepository) UpdateUserEmail(ctx context.Context, userID, email string) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrNotFound
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"email":            email,
		"last_modified_at": time.Now(),
	}

	_, err = r.userDal.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (*userPort.UserEmail, error) {
	filter := bson.M{
		"email": email,
		// "is_deleted": bson.M{"$ne": true},
	}

	userEntity, err := r.userDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &userPort.UserEmail{
		ID:    userEntity.ID.Hex(),
		Email: userEntity.Email,
	}, nil
}

func (r *MongoRepository) FindAccountUserByID(ctx context.Context, id string) (*accountPort.AccountUser, error) {
	// fmt.Printf("lalalal", id)
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}

	userFilter := bson.M{
		"_id":        oid,
		"is_deleted": bson.M{"$ne": true},
	}
	userEntity, err := r.userDal.FindOne(ctx, userFilter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &accountPort.AccountUser{
		ID:               userEntity.ID.Hex(),
		PhoneNumber:      userEntity.PhoneNumber,
		KYCLevel:         userEntity.KYC.KYCLevel,
		BranchCode:       userEntity.BranchCode,
		FullName:         userEntity.FullName,
		RegistrationType: "",
		AndOrStatus:      true,
	}, nil
}

func (r *MongoRepository) UpdateUserCustomerNumber(ctx context.Context, userId string, customerNumber string) error {
	oid, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Printf("Error converting userId to ObjectID: %v\n", err)
		return ErrNotFound
	}

	filter := bson.M{
		"_id": oid,
	}
	update := bson.M{
		"customer_number":  customerNumber,
		"last_modified_at": time.Now().UTC(),
	}
	_, err = r.userDal.UpdateOne(ctx, filter, update)
	fmt.Printf("Error converting userId to ObjectID: %v\n", err)
	return err
}

func (r *MongoRepository) CreateLinkedAccount(ctx context.Context, acc *accountPort.LinkedAccount) error {
	oid, err := bson.ObjectIDFromHex(acc.UserID)
	if err != nil {
		return ErrNotFound
	}

	linkedAccount := entities.LinkedAccount{
		UserID:            bson.ObjectID(oid),
		CustomerNumber:    acc.CustomerNumber,
		AccountNumber:     acc.AccountNumber,
		AccountHolderName: acc.AccountHolderName,
		AccountType:       acc.AccountType,
		BranchCode:        acc.BranchCode,
		LinkedStatus:      true,
		LastLinkedStatus:  false,
		LinkedAt:          time.Now().UTC(),
		IsAccountActive:   true,
		AndOrStatus:       acc.AndOrStatus,
		AccountBranchCode: acc.BranchCode,
		CurrencyCode:      acc.CurrencyCode,
		IsMain:            acc.IsMain,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	_, err = r.linkedAccountDal.InsertOne(ctx, linkedAccount)
	return err
}

func (r *MongoRepository) UpdateProfileImageURL(ctx context.Context, id string, imageURL string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}
	fmt.Printf("UpdateProfileImageURL called with id: %s (type: %T)\n", oid, oid)
	fmt.Printf("oid type: %T, oid value: %v\n", oid, oid)
	userFilter := bson.M{
		"_id": oid,
		// "is_deleted": bson.M{"$ne": true},
	}

	updateProfileURL := bson.M{
		"avatar": imageURL,
	}
	fmt.Printf("Updating avatar with URL: %s\n", imageURL)

	_, err = r.userDal.UpdateOne(ctx, userFilter, updateProfileURL)
	if err != nil {
		return err
	}
	user, _ := r.userDal.FindOne(ctx, bson.M{"_id": oid}, bson.M{})
	fmt.Println("Stored avatar URL:", user.Avatar)

	return nil
}

func (r *MongoRepository) DeleteOtp(ctx context.Context, ID string) error {
	oid, err := bson.ObjectIDFromHex(ID)
	if err != nil {
		return ErrNotFound
	}
	filter := bson.M{"_id": oid}

	err = r.otpDal.DeleteOne(
		ctx,
		filter,
	)
	if err != nil {
		fmt.Printf("Error deleting OTP records: %v\n", err)
		return fmt.Errorf("failed to delete OTP records: %v", err)
	}
	return nil
}
