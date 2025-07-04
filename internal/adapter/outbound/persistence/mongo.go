package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	userPort "cbe-super-app-member-users/internal/domain/users"
	accountPort "cbe-super-app-member-users/internal/port/outbound/account"

	"cbe-super-app-member-users/pkgs/constants"
	"cbe-super-app-member-users/pkgs/entities/enums"
	"cbe-super-app-member-users/pkgs/entities/type_definition"

	entities "cbe-super-app-member-users/internal/adapter/outbound/model"

	dal "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	localModel "cbe-super-app-member-users/pkgs/entities"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrNotFound = errors.New("not found")

type MongoRepository struct {
	userDal          dal.MongoDal[entities.User, entities.User]
	linkedAccountDal dal.MongoDal[entities.LinkedAccount, entities.LinkedAccount]
	otpDal           dal.MongoDal[entities.OTP, entities.OTP]
	hqDal            dal.MongoDal[entities.HQ, entities.HQ]
	client           *mongo.Client
	dbName           string
}

func NewMongoRepository(client *mongo.Client, dbName string) *MongoRepository {
	return &MongoRepository{
		userDal:          dal.NewMongoDal[entities.User, entities.User](client, dbName, "user"),
		linkedAccountDal: dal.NewMongoDal[entities.LinkedAccount, entities.LinkedAccount](client, dbName, "linked_accounts"),
		otpDal:           dal.NewMongoDal[entities.OTP, entities.OTP](client, dbName, "otp"),
		hqDal:            dal.NewMongoDal[entities.HQ, entities.HQ](client, dbName, "hq"),
		client:           client,
		dbName:           dbName,
	}
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*userPort.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("INVALID_USER_ID: %w", err)
	}

	userFilter := bson.M{
		"_id":        objectID,
		"is_deleted": bson.M{"$ne": true},
	}
	userEntity, err := r.userDal.FindOne(ctx, userFilter, nil)
	if err != nil {

		return nil, fmt.Errorf("AUTH_USER_NOT_FOUND: %w", err)
	}

	return &userPort.User{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		FullName:  userEntity.FullName,
		IsDeleted: userEntity.IsDeleted,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: userEntity.Device.DeviceUUID,
			AppVersion: userEntity.Device.AppVersion,
		},
		LoginPIN: type_definition.LoginPIN{
			PIN:              userEntity.LoginPIN.PIN,
			PINHistory:       userEntity.LoginPIN.PINHistory,
			LastPINCreatedAt: userEntity.LoginPIN.LastPINCreatedAt,
		},
	}, nil
}

func (r *MongoRepository) GetOneHQ(ctx context.Context, req map[string]interface{}) (*userPort.HQ, error) {
	hqFilter := bson.M{"enabled": true}
	// sort := bson.M{"$orderBy": bson.M{"last_modified": -1}}
	hqEntity, err := r.hqDal.FindOne(ctx, hqFilter, nil)
	if err != nil {
		return nil, err
	}

	return &userPort.HQ{
		// ID:                   hqEntity.ID,
		UniqueID:             hqEntity.UniqueID,
		Name:                 hqEntity.Name,
		Address:              hqEntity.Address,
		PhoneNumber:          hqEntity.PhoneNumber,
		Email:                hqEntity.Email,
		LinkedAccounts:       nil,
		LatestiOSVersion:     hqEntity.LatestiOSVersion,
		LatestAndroidVersion: hqEntity.LatestAndroidVersion,
		ArchiveExpiry:        hqEntity.ArchiveExpiry,
		BlockTime:            hqEntity.BlockTime,
		BlockTimeStatus:      hqEntity.BlockTimeStatus,
		ArchiveTime:          hqEntity.ArchiveTime,
		ArchiveTimeStatus:    hqEntity.ArchiveTimeStatus,
		Enabled:              hqEntity.Enabled,
		IsDeleted:            hqEntity.IsDeleted,
		CreatedAt:            hqEntity.CreatedAt,
		LastModified:         hqEntity.LastModified,
	}, nil
}

func (r *MongoRepository) GetOneUser(ctx context.Context, req map[string]interface{}) (*localModel.User, error) {
	// fmt.Println("hello there ")
	filter := bson.M{}

	filter["is_deleted"] = false
	for key, value := range req {
		if key == "id" {
			// oid, err := bson.ObjectIDFromHex(value.(string))
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid id type: expected string")
			}
			oid, err := bson.ObjectIDFromHex(strValue)

			if err != nil {
				return nil, ErrNotFound
			}
			filter["_id"] = oid
		} else {
			filter[key] = value
		}
	}

	user, err := r.userDal.FindOne(ctx, filter, nil)

	if err != nil {
		return nil, err
	}

	return &localModel.User{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		IsDeleted: user.IsDeleted,
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

func (r *MongoRepository) FindOTP(ctx context.Context, userID, otpFor string) (*userPort.OTPRecord, error) {
	filter := bson.M{
		"user_code": userID,
		"otp_for":   otpFor,
	}
	fmt.Println("filter", filter)
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

	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) CreateLinkedAccount(ctx context.Context, acc *accountPort.LinkedAccount) error {
	oid, err := bson.ObjectIDFromHex(acc.UserID)
	if err != nil {
		return ErrNotFound
	}

	linkedAccount := entities.LinkedAccount{
		UserID:            oid,
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

	userFilter := bson.M{
		"_id":        oid,
		"is_deleted": bson.M{"$ne": true},
	}

	updateProfileURL := bson.M{
		"avatar": imageURL,
	}

	_, err = r.userDal.UpdateOne(ctx, userFilter, updateProfileURL)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) DeleteOtp(ctx context.Context, userCode, otpCode, otpFor string) error {

	filter := bson.M{"otp_code": otpCode, "otp_for": otpFor, "user_code": userCode}
	fmt.Println("filter", filter)
	err := r.otpDal.DeleteOne(ctx, filter)

	if err != nil {
		fmt.Printf("Error deleting OTP records: %v\n", err)
		return fmt.Errorf("failed to delete OTP records: %v", err)
	}
	return nil
}

func (r *MongoRepository) UnlinkDevice(ctx context.Context, userID string, deviceID string) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrNotFound
	}

	filter := bson.M{
		"_id": oid,
	}

	update := bson.M{
		"$unset": bson.M{"device": ""},
		"$set": bson.M{
			"device_status":    "unlinked",
			"last_modified_at": time.Now().UTC(),
		},
	}

	_, err = r.userDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to unlink device: %w", err)
	}

	return nil
}

func (r *MongoRepository) ChangePin(ctx context.Context, userID string, loginPIN type_definition.LoginPIN) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrNotFound
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"login_pin":        loginPIN,
		"is_verified":      true,
		"last_modified_at": time.Now().UTC(),
	}

	_, err = r.userDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to change pin: %w", err)
	}

	return nil
}

// OTP CRUD
func (r *MongoRepository) CreateOtp(ctx context.Context, otp *userPort.OTPRecord) error {
	otpEntity := entities.OTP{
		UserCode:  otp.UserCode,
		OTPFor:    entities.OTPFor(otp.OTPFor), // or make this a parameter if needed
		OTPCode:   otp.OTP,
		Email:     otp.Email,
		ExpiresAt: otp.ExpiresAt,
		CreatedAt: otp.CreatedAt,
		Status:    "",
		IsDeleted: false,
	}
	_, err := r.otpDal.InsertOne(ctx, otpEntity)
	return err
}

func (r *MongoRepository) GetOtpByID(ctx context.Context, id string) (*userPort.OTPRecord, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}
	filter := bson.M{"_id": oid, "is_deleted": bson.M{"$ne": true}}
	otpEntity, err := r.otpDal.FindOne(ctx, filter, nil)
	if err != nil {
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

func (r *MongoRepository) UpdateOtp(ctx context.Context, otp *userPort.OTPRecord) error {
	oid, err := bson.ObjectIDFromHex(otp.ID)
	if err != nil {
		return ErrNotFound
	}
	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"otp_code":   otp.OTP,
			"expires_at": otp.ExpiresAt,
		},
	}

	_, err = r.otpDal.UpdateOne(ctx, filter, update)
	return err
}

// Login-related methods
func (r *MongoRepository) FindUserByPhone(ctx context.Context, phone string) (*userPort.User, error) {
	filter := bson.M{
		"phone_number": phone,
		"is_deleted":   false,
	}

	userEntity, err := r.userDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return r.mapUserEntityToDomain(userEntity), nil
}

func (r *MongoRepository) FindUserByPhoneForLogin(ctx context.Context, phone string, pin string) (*userPort.User, error) {
	filter := bson.M{
		"phone_number":  phone,
		"login_pin.pin": pin,
		"is_deleted":    false,
	}

	userEntity, err := r.userDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return r.mapUserEntityToDomain(userEntity), nil
}

func (r *MongoRepository) FindUserByDevice(ctx context.Context, deviceUUID string) (*userPort.User, error) {
	filter := bson.M{
		"device.device_uuid": deviceUUID,
		"is_deleted":         false,
	}

	fmt.Println("filter", filter)
	userEntity, err := r.userDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return r.mapUserEntityToDomain(userEntity), nil
}

func (r *MongoRepository) IncrementLoginAttempts(ctx context.Context, userID string) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrNotFound
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$inc": bson.M{
			"login_attempt_count": 1,
		},
		"$set": bson.M{
			"last_login_attempt": time.Now(),
		},
	}

	_, err = r.userDal.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) ResetLoginAttempts(ctx context.Context, userID string) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrNotFound
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"login_attempt_count": 0,
			"last_login_attempt":  time.Time{},
		},
	}

	_, err = r.userDal.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrNotFound
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"last_login":          time.Now(),
			"login_attempt_count": 0,
			"last_login_attempt":  time.Time{},
		},
	}

	_, err = r.userDal.UpdateOne(ctx, filter, update)
	return err
}

// Helper method to map user entity to domain model
func (r *MongoRepository) mapUserEntityToDomain(userEntity *entities.User) *userPort.User {
	return &userPort.User{
		ID:                userEntity.ID,
		UserCode:          userEntity.UserCode,
		FullName:          userEntity.FullName,
		PhoneNumber:       userEntity.PhoneNumber,
		Email:             userEntity.Email,
		IsDeleted:         userEntity.IsDeleted,
		IsAccountBlocked:  userEntity.IsAccountBlocked,
		IsVerified:        userEntity.IsVerified,
		LoginAttemptCount: userEntity.LoginAttemptCount,
		LastLoginAttempt:  userEntity.LastLoginAttempt,
		LastLogin:         userEntity.LastLogin,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: userEntity.Device.DeviceUUID,
			AppVersion: userEntity.Device.AppVersion,
		},
		KYC: struct {
			KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
			KYCStatus            enums.KYCStatus     `json:"kyc_status" bson:"kyc_status"`
			KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
			KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
			KYCActivityBy        map[string]struct{} `json:"kyc_activity_by" bson:"kyc_activity_by"`
			KYCLevel             uint8               `json:"level" bson:"level"`
		}{
			KYCLevel: userEntity.KYC.KYCLevel,
		},
		LoginPIN: type_definition.LoginPIN{
			PIN:              userEntity.LoginPIN.PIN,
			PINHistory:       userEntity.LoginPIN.PINHistory,
			LastPINCreatedAt: userEntity.LoginPIN.LastPINCreatedAt,
		},
		CreatedAt:      userEntity.CreatedAt,
		LastModifiedAt: userEntity.LastModifiedAt,
	}
}

// Registration-related methods
func (r *MongoRepository) FindPendingRegistration(ctx context.Context, userID, deviceUUID string) (*userPort.RegistrationRecord, error) {
	filter := bson.M{
		"user_code":   userID,
		"device_uuid": deviceUUID,
		"status":      string(entities.Pending),
		"is_deleted":  false,
		"expires_at":  bson.M{"$gt": time.Now()},
	}

	fmt.Println("filter", filter)
	otpEntity, err := r.otpDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	deviceUUIDStr := getDeviceUUIDString(otpEntity.DeviceUUID)

	return &userPort.RegistrationRecord{
		ID:          otpEntity.ID.Hex(),
		PhoneNumber: otpEntity.PhoneNumber,
		DeviceUUID:  deviceUUIDStr,
		Platform:    constants.DefaultPlatform, // Default platform since it's not in OTP model
		OTP:         otpEntity.OTPCode,
		OTPFor:      string(otpEntity.OTPFor),
		Status:      string(otpEntity.Status),
		ExpiresAt:   otpEntity.ExpiresAt,
		CreatedAt:   otpEntity.CreatedAt,
		Attempts:    constants.DefaultAttempts,
		MaxAttempts: constants.DefaultMaxAttempts,
	}, nil
}

func getDeviceUUIDString(deviceUUID *string) string {
	if deviceUUID != nil {
		return *deviceUUID
	}
	return ""
}
func (r *MongoRepository) FindPendingRegistrationByID(ctx context.Context, registrationID string) (*userPort.RegistrationRecord, error) {
	oid, err := bson.ObjectIDFromHex(registrationID)
	if err != nil {
		return nil, ErrNotFound
	}

	filter := bson.M{"_id": oid}
	otpEntity, err := r.otpDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Handle nullable DeviceUUID
	deviceUUIDStr := ""
	if otpEntity.DeviceUUID != nil {
		deviceUUIDStr = *otpEntity.DeviceUUID
	}

	return &userPort.RegistrationRecord{
		ID:          otpEntity.ID.Hex(),
		PhoneNumber: otpEntity.PhoneNumber,
		DeviceUUID:  deviceUUIDStr,
		Platform:    "android", // Default platform since it's not in OTP model
		OTP:         otpEntity.OTPCode,
		OTPFor:      string(otpEntity.OTPFor),
		Status:      string(otpEntity.Status),
		ExpiresAt:   otpEntity.ExpiresAt,
		CreatedAt:   otpEntity.CreatedAt,
		Attempts:    0, // Default since it's not in OTP model
		MaxAttempts: 3, // Default since it's not in OTP model
	}, nil
}

func (r *MongoRepository) CreatePendingRegistration(ctx context.Context, registration *userPort.RegistrationRecord) error {
	otpEntity := entities.OTP{
		PhoneNumber: registration.PhoneNumber,
		DeviceUUID:  &registration.DeviceUUID,
		OTPCode:     registration.OTP,
		OTPFor:      entities.OTPFor(registration.OTPFor),
		Status:      entities.OTPStatus(registration.Status),
		ExpiresAt:   registration.ExpiresAt,
		CreatedAt:   registration.CreatedAt,
	}

	_, err := r.otpDal.InsertOne(ctx, otpEntity)
	return err
}

func (r *MongoRepository) DeletePendingRegistration(ctx context.Context, registrationID string) error {
	oid, err := bson.ObjectIDFromHex(registrationID)
	if err != nil {
		return ErrNotFound
	}

	filter := bson.M{"_id": oid}
	err = r.otpDal.DeleteOne(ctx, filter)
	return err
}

func (r *MongoRepository) UpdatePendingRegistration(ctx context.Context, registration *userPort.RegistrationRecord) error {
	oid, err := bson.ObjectIDFromHex(registration.ID)
	if err != nil {
		return ErrNotFound
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"status":     entities.OTPStatus(registration.Status),
			"expires_at": registration.ExpiresAt,
		},
	}

	_, err = r.otpDal.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) CreateUser(ctx context.Context, user *userPort.User) error {
	userEntity := entities.User{
		ID:                user.ID,
		UserCode:          user.UserCode,
		FullName:          user.FullName,
		PhoneNumber:       user.PhoneNumber,
		Email:             user.Email,
		IsDeleted:         user.IsDeleted,
		IsAccountBlocked:  user.IsAccountBlocked,
		IsVerified:        user.IsVerified,
		LoginAttemptCount: user.LoginAttemptCount,
		LastLoginAttempt:  user.LastLoginAttempt,
		LastLogin:         user.LastLogin,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: user.Device.DeviceUUID,
			AppVersion: user.Device.AppVersion,
		},
		KYC: struct {
			KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
			KYCStatus            entities.KYCStatus  `json:"kyc_status" bson:"kyc_status"`
			KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
			KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
			KYCActivityBy        map[string]struct{} `json:"kyc_activity_by" bson:"kyc_activity_by"`
			KYCLevel             uint8               `json:"level" bson:"level"`
		}{
			KYCRejectReasonField: user.KYC.KYCRejectReasonField,
			KYCStatus:            entities.KYCStatus(user.KYC.KYCStatus),
			KYCRejectReason:      user.KYC.KYCRejectReason,
			KYCApproved:          user.KYC.KYCApproved,
			KYCActivityBy:        user.KYC.KYCActivityBy,
			KYCLevel:             user.KYC.KYCLevel,
		},
		LoginPIN: entities.LoginPIN{
			PIN:              user.LoginPIN.PIN,
			PINHistory:       user.LoginPIN.PINHistory,
			LastPINCreatedAt: user.LoginPIN.LastPINCreatedAt,
		},
		CreatedAt:      user.CreatedAt,
		LastModifiedAt: user.LastModifiedAt,
	}

	_, err := r.userDal.InsertOne(ctx, userEntity)
	return err
}

// PIN Reset methods
func (r *MongoRepository) CreatePinResetSession(ctx context.Context, session *userPort.PinResetSession) error {
	collection := r.client.Database(r.dbName).Collection("pin_reset_sessions")

	// Convert to BSON document
	doc := bson.M{
		"_id":               session.ID,
		"user_id":           session.UserID,
		"phone_number":      session.PhoneNumber,
		"device_uuid":       session.DeviceUUID,
		"otp":               session.OTP,
		"otp_for":           session.OTPFor,
		"status":            session.Status,
		"expires_at":        session.ExpiresAt,
		"created_at":        session.CreatedAt,
		"attempts":          session.Attempts,
		"max_attempts":      session.MaxAttempts,
		"verified_at":       session.VerifiedAt,
		"completed_at":      session.CompletedAt,
		"access_restricted": session.AccessRestricted,
		"restrictions":      session.Restrictions,
	}

	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to create PIN reset session: %w", err)
	}

	return nil
}

func (r *MongoRepository) FindPinResetSession(ctx context.Context, sessionID string) (*userPort.PinResetSession, error) {
	collection := r.client.Database(r.dbName).Collection("pin_reset_sessions")

	var doc bson.M
	err := collection.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("PIN_RESET_SESSION_NOT_FOUND")
		}
		return nil, fmt.Errorf("failed to find PIN reset session: %w", err)
	}

	session := &userPort.PinResetSession{}

	if v, ok := doc["_id"].(string); ok {
		session.ID = v
	}
	if v, ok := doc["user_id"].(string); ok {
		session.UserID = v
	}
	if v, ok := doc["phone_number"].(string); ok {
		session.PhoneNumber = v
	}
	if v, ok := doc["device_uuid"].(string); ok {
		session.DeviceUUID = v
	}
	if v, ok := doc["otp"].(string); ok {
		session.OTP = v
	}
	if v, ok := doc["otp_for"].(string); ok {
		session.OTPFor = v
	}
	if v, ok := doc["status"].(string); ok {
		session.Status = v
	}
	if v, ok := doc["expires_at"].(bson.DateTime); ok {
		session.ExpiresAt = v.Time()
	}
	if v, ok := doc["created_at"].(bson.DateTime); ok {
		session.CreatedAt = v.Time()
	}
	if v, ok := doc["attempts"].(int32); ok {
		session.Attempts = int(v)
	}
	if v, ok := doc["max_attempts"].(int32); ok {
		session.MaxAttempts = int(v)
	}
	if v, ok := doc["access_restricted"].(bool); ok {
		session.AccessRestricted = v
	}
	if v, ok := doc["restrictions"].(bson.A); ok {
		session.Restrictions = convertToStringSlice(v)
	}

	// Handle optional time fields
	if verifiedAt, ok := doc["verified_at"]; ok && verifiedAt != nil {
		session.VerifiedAt = verifiedAt.(bson.DateTime).Time()
	}
	if completedAt, ok := doc["completed_at"]; ok && completedAt != nil {
		session.CompletedAt = completedAt.(bson.DateTime).Time()
	}

	return session, nil
}

func (r *MongoRepository) FindPinResetSessionByPhone(ctx context.Context, phone, deviceUUID string) (*userPort.PinResetSession, error) {
	collection := r.client.Database(r.dbName).Collection("pin_reset_sessions")

	var doc bson.M
	err := collection.FindOne(ctx, bson.M{
		"phone_number": phone,
		"device_uuid":  deviceUUID,
		"status":       "pending",
	}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No active session found
		}
		return nil, fmt.Errorf("failed to find PIN reset session by phone: %w", err)
	}

	session := &userPort.PinResetSession{
		ID:               doc["_id"].(string),
		UserID:           doc["user_id"].(string),
		PhoneNumber:      doc["phone_number"].(string),
		DeviceUUID:       doc["device_uuid"].(string),
		OTP:              doc["otp"].(string),
		OTPFor:           doc["otp_for"].(string),
		Status:           doc["status"].(string),
		ExpiresAt:        doc["expires_at"].(bson.DateTime).Time(),
		CreatedAt:        doc["created_at"].(bson.DateTime).Time(),
		Attempts:         int(doc["attempts"].(int32)),
		MaxAttempts:      int(doc["max_attempts"].(int32)),
		AccessRestricted: doc["access_restricted"].(bool),
		Restrictions:     convertToStringSlice(doc["restrictions"].(bson.A)),
	}

	// Handle optional time fields
	if verifiedAt, ok := doc["verified_at"]; ok && verifiedAt != nil {
		session.VerifiedAt = verifiedAt.(bson.DateTime).Time()
	}
	if completedAt, ok := doc["completed_at"]; ok && completedAt != nil {
		session.CompletedAt = completedAt.(bson.DateTime).Time()
	}

	return session, nil
}

func (r *MongoRepository) UpdatePinResetSession(ctx context.Context, session *userPort.PinResetSession) error {
	collection := r.client.Database(r.dbName).Collection("pin_reset_sessions")

	update := bson.M{
		"$set": bson.M{
			"status":            session.Status,
			"attempts":          session.Attempts,
			"verified_at":       session.VerifiedAt,
			"completed_at":      session.CompletedAt,
			"access_restricted": session.AccessRestricted,
			"restrictions":      session.Restrictions,
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": session.ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update PIN reset session: %w", err)
	}

	return nil
}

func (r *MongoRepository) DeletePinResetSession(ctx context.Context, sessionID string) error {
	collection := r.client.Database(r.dbName).Collection("pin_reset_sessions")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": sessionID})
	if err != nil {
		return fmt.Errorf("failed to delete PIN reset session: %w", err)
	}

	return nil
}

func (r *MongoRepository) IncrementPinResetAttempts(ctx context.Context, sessionID string) error {
	collection := r.client.Database(r.dbName).Collection("pin_reset_sessions")

	update := bson.M{
		"$inc": bson.M{"attempts": 1},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": sessionID}, update)
	if err != nil {
		return fmt.Errorf("failed to increment PIN reset attempts: %w", err)
	}

	return nil
}

func (r *MongoRepository) ResetPinResetAttempts(ctx context.Context, sessionID string) error {
	collection := r.client.Database(r.dbName).Collection("pin_reset_sessions")

	update := bson.M{
		"$set": bson.M{"attempts": 0},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": sessionID}, update)
	if err != nil {
		return fmt.Errorf("failed to reset PIN reset attempts: %w", err)
	}

	return nil
}

// Helper function to convert bson.A to []string
func convertToStringSlice(a bson.A) []string {
	result := make([]string, 0, len(a))
	for _, v := range a {
		if str, ok := v.(string); ok {
			result = append(result, str)
		}
	}

	return result

}
