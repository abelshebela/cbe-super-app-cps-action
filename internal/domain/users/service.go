package users

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"cbe-super-app-member-users/internal/application/dto"
	"cbe-super-app-member-users/pkgs/entities"
	"cbe-super-app-member-users/pkgs/entities/enums"
	"cbe-super-app-member-users/pkgs/entities/type_definition"
	"cbe-super-app-member-users/pkgs/utils"

	"github.com/google/uuid"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService struct {
	// repository userRepoPort.UserRepositoryPort
	repository UserRepository
	logger     shared_utils.Logger
	minIO      config.MinioClientInterface
	cfg        *config.VaultConfig
}

func NewUserService(repository UserRepository, logger shared_utils.Logger, minIO config.MinioClientInterface, cfg *config.VaultConfig) *UserService {
	return &UserService{
		repository: repository,
		logger:     logger,
		minIO:      minIO,
		cfg:        cfg,
	}
}

func (s *UserService) UpdateProfilePicture(ctx context.Context, id string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			s.logger.Errorf("User with ID %s not found", id)
			return "", fmt.Errorf("NOT_FOUND")
		}
		s.logger.Errorf("Failed to find user with ID %s: %v", id, err)
		return "", fmt.Errorf("UPLOAD_FAILED")
	}

	// Create a temporary file to store the uploaded content
	tempFile, err := os.CreateTemp("", "profile-*.tmp")
	if err != nil {
		s.logger.Errorf("Failed to create temporary file: %v", err)
		return "", fmt.Errorf("UPLOAD_FAILED")
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the uploaded file content to the temporary file
	if _, err := io.Copy(tempFile, file); err != nil {
		s.logger.Errorf("Failed to copy file content to temporary file: %v", err)
		return "", fmt.Errorf("UPLOAD_FAILED")
	}

	objectName := fmt.Sprintf("profile-pictures/%s/%s", id, filepath.Base(fileHeader.Filename))

	resp, err := s.minIO.SaveObject(ctx, config.SaveObjectBody{
		BucketName:  "user-profile-pictures",
		ObjectName:  objectName,
		File:        tempFile.Name(),
		ContentType: "jpeg",
	})
	if err != nil {
		s.logger.Errorf("Failed to save object to minio: %v", err)
		return "", fmt.Errorf("UPLOAD_FAILED")
	}

	return resp.Key, nil
}

func (s *UserService) FetchLinkedAccounts(ctx context.Context, id string) (*LinkedAccountResponse, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			s.logger.Warnf("User with ID %s not found", id)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		s.logger.Errorf("Failed to fetch user with ID %s: %v", id, err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	if user.IsDeleted {
		s.logger.Infof("User with ID %s is deleted", id)
		return &LinkedAccountResponse{
			UserID:         user.ID.Hex(),
			FullName:       user.FullName,
			LinkedAccounts: []LinkedAccountDetail{},
		}, nil
	}

	linkedAccounts, err := s.repository.FindActiveLinkedAccounts(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch linked accounts for user ID %s: %v", id, err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}
	var linkedAccountsResponse []LinkedAccountDetail
	for _, account := range linkedAccounts {
		linkedAccountsResponse = append(linkedAccountsResponse, LinkedAccountDetail{
			AccountNumber:     account.AccountNumber,
			AccountBranchCode: account.AccountBranchCode,
			LinkedBranch:      account.LinkedBranch,
			IsAccountActive:   account.IsAccountActive,
			LinkedStatus:      account.LinkedStatus,
			CurrencyCode:      account.CurrencyCode,
		})
	}

	response := &LinkedAccountResponse{
		UserID:         user.ID.Hex(),
		FullName:       user.FullName,
		LinkedAccounts: linkedAccountsResponse,
	}

	s.logger.Infof("Successfully fetched active linked accounts for user ID %s", id)
	return response, nil
}

func (s *UserService) GetOneHQ(ctx context.Context, req map[string]interface{}) (*HQ, error) {
	hq, err := s.repository.GetOneHQ(ctx, req)
	if err != nil {
		if err == ErrNotFound {
			s.logger.Warnf("HQ data not found", req)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		s.logger.Errorf("Failed to fetch HQ data: %v", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	s.logger.Infof("Successfully fetched HQ data", req)
	return hq, nil
}

func (s *UserService) GetOneUser(ctx context.Context, req map[string]interface{}) (*entities.User, error) {
	user, err := s.repository.GetOneUser(ctx, req)
	if err != nil {

		if err == ErrNotFound {
			s.logger.Warnf("User data not found", req)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		s.logger.Errorf("Failed to fetch User data: %v", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	s.logger.Infof("Successfully fetched User data", req)
	return user, nil
}

func (s *UserService) GenerateEmailOTP(ctx context.Context, req OTPRequest) (string, error) {
	currentUser, err := s.repository.FindByID(ctx, req.UserID)
	if err != nil {
		s.logger.Errorf("Failed to find user by ID %s: %v", req.UserID, err)
		return "UNHANDLED_SERVER_ERROR", fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	if currentUser != nil && currentUser.Email != "" {
		s.logger.Errorf("User %s already has an email address", req.UserID)
		return "USER_ALREADY_HAS_EMAIL", fmt.Errorf("USER_ALREADY_HAS_EMAIL")
	}

	existingUser, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil && err.Error() != "not found" {
		s.logger.Errorf("Failed to find user by email %s: %v", req.Email, err)
		return "UNHANDLED_SERVER_ERROR", fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	if existingUser != nil {
		s.logger.Errorf("Email %s is already in use", req.Email)
		return "EMAIL_IN_USE", fmt.Errorf("EMAIL_IN_USE")
	}

	record, err := s.repository.FindOTP(ctx, req.UserID, req.Email)
	if record != nil && !record.ExpiresAt.IsZero() && time.Now().Before(record.ExpiresAt) {
		s.logger.Errorf("wait until the previous otp expired", req.UserID, err)
		return "WAIT_FOR_PREVIOUS_OTP_EXPIRATION", fmt.Errorf("WAIT_FOR_PREVIOUS_OTP_EXPIRATION")
	}

	otp := utils.OTPGenerator(6)

	encryptedOTPCode, _, err := utils.LocalEncryptPassword(otp, "otp", "", "", s.cfg)
	if err != nil {
		s.logger.Errorf("OTP encryption failed: %v", err)
		return "UNHANDLED_SERVER_ERROR", fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	otp_record := &OTPRecord{
		UserID:    req.UserID,
		Email:     req.Email,
		OTP:       encryptedOTPCode,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := s.repository.StoreOTP(ctx, otp_record); err != nil {
		s.logger.Errorf("OTP storage failed: %v", err)
		return "UNHANDLED_SERVER_ERROR", fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	return otp, nil
}

func (s *UserService) VerifyEmailOTP(ctx context.Context, verification OTPVerification) error {
	record, err := s.repository.FindOTP(ctx, verification.UserID, verification.Email)
	s.logger.Infof("record: %+v, err: %v", record, err)
	if record == nil {
		s.logger.Errorf("OTP record not found for user %s, email %s", verification.UserID, verification.Email)
		return fmt.Errorf("INVALID_OTP")
	}

	if err != nil {
		s.logger.Errorf("OTP lookup failed for user %s: %v", verification.UserID, err)
		return fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	decryptedOTP, err := utils.LocalDecryptPassword(record.OTP, s.cfg)
	if err != nil {
		s.logger.Errorf("OTP decryption failed for user %s: %v", verification.UserID, err)
		return fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	if decryptedOTP != verification.OTP {
		s.logger.Warnf("Invalid OTP provided for user %s (expected: %s, got: %s)",
			verification.UserID, decryptedOTP, verification.OTP)
		return fmt.Errorf("INVALID_OTP")
	}

	if time.Now().After(record.ExpiresAt) {
		s.logger.Warnf("Expired OTP for user %s (expired at: %v)", verification.UserID, record.ExpiresAt)
		return fmt.Errorf("EXPIRED_OTP")
	}

	if err := s.repository.UpdateUserEmail(ctx, verification.UserID, verification.Email); err != nil {
		s.logger.Errorf("Email update failed for user %s: %v", verification.UserID, err)
		return fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	// TODO Delete the OTPb record --->Hard delete
	if err := s.repository.DeleteOtp(ctx, record.OTP, record.UserCode, record.OTPFor); err != nil {
		s.logger.Errorf("Failed to delete OTP record for user %s: %v", verification.UserID, err)
		return fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	s.logger.Infof("Email OTP verified successfully for user %s", verification.UserID)
	return nil
}

func (s *UserService) UnlinkDevice(ctx context.Context, UserID string, DeviceID string) error {
	currentUser, err := s.repository.FindByID(ctx, UserID)

	if err != nil {
		s.logger.Errorf("Failed to find user by ID %s: %v", UserID, err)
		return fmt.Errorf("NOT_FOUND")
	}
	// s.logger.Infof("mmmm",currentUser.Device)
	if currentUser.Device.DeviceUUID == DeviceID {

		err = s.repository.UnlinkDevice(ctx, UserID, DeviceID)
		if err != nil {
			s.logger.Errorf("Failed to unlink device %s for user %s: %v", DeviceID, UserID, err)
			return fmt.Errorf("COULD_NOT_UNLINK_DEVICE")

		}
		s.logger.Infof("Successfully unlinked device %s for user %s", DeviceID, UserID)
		return nil
	} else {
		s.logger.Errorf("User %s has no linked devices", UserID)
		return fmt.Errorf("NO_LINKED_DEVICES")
	}

}

func (s *UserService) FindByID(ctx context.Context, id string) (*User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			s.logger.Warnf("User with ID %s not found", id)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		s.logger.Errorf("Failed to find user with ID %s: %v", id, err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}
	return user, nil
}
func (a *UserService) ValidatePin(ctx context.Context, pin string) (bool, error) {

	if pin == "" {
		return false, nil
	}
	if len(pin) != 6 {
		return false, fmt.Errorf("PIN_LIMIT")
	}

	for _, c := range pin {
		if c < '0' || c > '9' {
			return false, fmt.Errorf("PIN_OLY_DIG")
		}
	}

	maxRedundant := 1
	count := 1
	for i := 1; i < len(pin); i++ {
		if pin[i] == pin[i-1] {
			count++
			if count > maxRedundant {
				maxRedundant = count
			}
		} else {
			count = 1
		}
	}
	if maxRedundant > 4 {
		return false, fmt.Errorf("PIN_REDANDANT")
	}

	for i := 0; i <= len(pin)-4; i++ {
		asc, desc := true, true
		for j := 1; j < 4; j++ {
			if pin[i+j]-pin[i+j-1] != 1 {
				asc = false
			}
			if pin[i+j-1]-pin[i+j] != 1 {
				desc = false
			}
		}
		if asc || desc {
			return false, fmt.Errorf("PIN_SEQ")
		}
	}

	return true, nil
}
func (s *UserService) ChangePin(ctx context.Context, ChangePinRequest ChangePinRequest) error {
	userData, err := s.repository.FindByID(ctx, ChangePinRequest.UserID)
	if err != nil {
		return fmt.Errorf("NOT_FOUND")
	}
	s.logger.Infof(userData.LoginPIN.PIN, ChangePinRequest.OldPin)

	if userData.LoginPIN.PIN != ChangePinRequest.OldPin {
		return fmt.Errorf("OLD_PIN_MISMATCH")
	}

	isValid, err := s.ValidatePin(ctx, ChangePinRequest.NewPin)
	if err != nil {
		return err
	}
	if !isValid {
		return err
	}

	if ChangePinRequest.OldPin == ChangePinRequest.NewPin {
		return fmt.Errorf("SAME_PIN")
	}

	for _, historyPin := range userData.LoginPIN.PINHistory {
		if historyPin == ChangePinRequest.NewPin {
			return fmt.Errorf("PIN_IN_HISTORY")
		}
	}

	var newHistory [4]string
	copy(newHistory[1:], userData.LoginPIN.PINHistory[:3])
	newHistory[0] = userData.LoginPIN.PIN
	newLoginPIN := type_definition.LoginPIN{
		PIN:              ChangePinRequest.NewPin,
		PINHistory:       newHistory,
		LastPINCreatedAt: time.Now(),
	}

	err = s.repository.ChangePin(ctx, ChangePinRequest.UserID, newLoginPIN)
	if err != nil {
		return fmt.Errorf("ERROR_CHANGING_PIN")
	}
	return nil
}

func (s *UserService) CreateOtp(ctx context.Context, otp OTPRecord) error {
	return s.repository.CreateOtp(ctx, &otp)
}

func (s *UserService) VerifyOtp(ctx context.Context, userID, otp string, deviceUUID *string, userRealm, otpFor string) (*dto.VerifyOtpResponse, error) {
	// Validate input parameters
	if err := s.validateVerifyOtpInputs(userID, otp, userRealm, otpFor); err != nil {
		return nil, err
	}

	// Clear sensitive data from memory after function execution
	defer func() {
		otp = ""
	}()

	// Encrypt OTP for comparison
	encryptedOtp, _, err := utils.LocalEncryptPassword(otp, "otp", "", "", s.cfg)
	if err != nil {
		s.logger.Errorf("Failed to encrypt OTP: %v", err)
		return nil, fmt.Errorf("OTP_PROCESSING_FAILED")
	}

	// Verify OTP
	err = s.verifyOtpInternal(ctx, userID, encryptedOtp, deviceUUID, userRealm, otpFor)
	if err != nil {
		s.logger.Errorf("OTP verification failed: %v", err)
		return nil, err
	}

	// Get user details for token generation
	var phone, fullName string
	if user, err := s.FindByID(ctx, userID); err == nil && user != nil {
		phone = user.PhoneNumber
		fullName = user.FullName
	} else if deviceUUID != nil {
		phone = ""
		fullName = ""
	}

	// Generate verify OTP token using TempTokenMaker
	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Errorf("Invalid user ID format: %v", err)
		return nil, fmt.Errorf("INVALID_USER_ID")
	}

	userEntity := &entities.User{
		ID:          objectID,
		FullName:    fullName,
		PhoneNumber: phone,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: *deviceUUID,
		},
	}

	permissions := []string{"verify_otp", "set_pin"}
	additional := map[string]interface{}{
		"otp_for":    otpFor,
		"token_type": "verify_otp",
	}

	token, err := utils.TempTokenMaker(userEntity, permissions, otpFor, additional, "verify_otp", s.cfg)
	if err != nil {
		s.logger.Errorf("Failed to generate verify OTP token: %v", err)
		return nil, fmt.Errorf("TOKEN_GENERATION_FAILED")
	}

	// Determine next step based on OTP purpose
	nextStep := "set_pin"
	if otpFor == "registration" {
		nextStep = "complete_registration"
	}

	response := &dto.VerifyOtpResponse{
		UserID:      userID,
		PhoneNumber: phone,
		OTPVerified: true,
		Token:       token,
		TokenType:   "verify_otp",
		TokenExpiry: time.Now().Add(10 * time.Minute),
		NextStep:    nextStep,
	}

	s.logger.Infof("OTP verified successfully for user %s, purpose: %s", userID, otpFor)
	return response, nil
}

func (s *UserService) verifyOtpInternal(ctx context.Context, userID, otp string, deviceUUID *string, userRealm, otpFor string) error {
	// For registration flow, we can use the FindPendingRegistration method
	if otpFor == "registration" && deviceUUID != nil {
		registration, err := s.repository.FindPendingRegistration(ctx, "", *deviceUUID)
		if err != nil {
			s.logger.Errorf("Failed to find registration record: %v", err)
			return fmt.Errorf("OTP_NOT_FOUND")
		}

		// Check if OTP is expired
		if time.Now().After(registration.ExpiresAt) {
			s.logger.Warnf("OTP expired for registration")
			return fmt.Errorf("EXPIRED_OTP")
		}

		// Compare OTP
		if registration.OTP != otp {
			s.logger.Warnf("Invalid OTP for registration")
			return fmt.Errorf("INVALID_OTP")
		}

		return nil
	}

	// For other flows, we'll need to implement proper OTP verification
	// For now, return an error indicating this needs to be implemented
	return fmt.Errorf("OTP_VERIFICATION_NOT_IMPLEMENTED")
}

func (s *UserService) SetPin(ctx context.Context, userID, newPin, otp string, deviceUUID *string, userRealm, otpFor string) (*dto.SetPinResponse, error) {
	// Validate input parameters
	if err := s.validateSetPinInputs(userID, newPin, otp, userRealm, otpFor); err != nil {
		return nil, err
	}

	// Clear sensitive data from memory after function execution
	defer func() {
		newPin = ""
		otp = ""
	}()

	// Encrypt OTP for comparison
	encryptedOtp, _, err := utils.LocalEncryptPassword(otp, "otp", "", "", s.cfg)
	if err != nil {
		s.logger.Errorf("Failed to encrypt OTP: %v", err)
		return nil, fmt.Errorf("OTP_PROCESSING_FAILED")
	}

	// Verify OTP first
	err = s.verifyOtpInternal(ctx, userID, encryptedOtp, deviceUUID, userRealm, otpFor)
	if err != nil {
		s.logger.Errorf("OTP verification failed for PIN setting: %v", err)
		return nil, err
	}

	// Find user
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to find user %s: %v", userID, err)
		return nil, fmt.Errorf("USER_NOT_FOUND")
	}

	// Check PIN history to prevent reuse
	if err := s.checkPinHistory(user, newPin); err != nil {
		s.logger.Warnf("PIN history check failed for user %s: %v", userID, err)
		return nil, err
	}

	// Update user's PIN
	if err := s.updateUserPin(ctx, userID, newPin, user); err != nil {
		s.logger.Errorf("Failed to update user PIN: %v", err)
		return nil, fmt.Errorf("PIN_UPDATE_FAILED")
	}

	// Generate permanent token using TokenMaker
	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Errorf("Invalid user ID format: %v", err)
		return nil, fmt.Errorf("INVALID_USER_ID")
	}

	userEntity := entities.User{
		ID:          objectID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: *deviceUUID,
		},
	}

	permissions := []string{"access", "transfer", "balance_check"}
	token, err := utils.TokenMaker(userEntity, permissions, "permanent")
	if err != nil {
		s.logger.Errorf("Failed to generate permanent token: %v", err)
		return nil, fmt.Errorf("TOKEN_GENERATION_FAILED")
	}

	response := &dto.SetPinResponse{
		UserID:      userID,
		PinSet:      true,
		Token:       token,
		TokenType:   "permanent",
		TokenExpiry: time.Now().Add(24 * time.Hour),
		NextStep:    "login",
	}

	s.logger.Infof("PIN set successfully for user %s", userID)
	return response, nil
}

// validateSetPinInputs validates the input parameters for setting PIN
func (s *UserService) validateSetPinInputs(userID, newPin, otp, userRealm, otpFor string) error {
	if userID == "" {
		return fmt.Errorf("INVALID_INPUT: user ID is required")
	}
	if newPin == "" {
		return fmt.Errorf("INVALID_INPUT: new PIN is required")
	}
	if otp == "" {
		return fmt.Errorf("INVALID_INPUT: OTP is required")
	}
	if userRealm == "" {
		return fmt.Errorf("INVALID_INPUT: user realm is required")
	}
	if otpFor == "" {
		return fmt.Errorf("INVALID_INPUT: OTP purpose is required")
	}
	return nil
}

// checkPinHistory checks if the new PIN is not in the user's PIN history
func (s *UserService) checkPinHistory(user *User, newPin string) error {
	for _, historyPin := range user.LoginPIN.PINHistory {
		if historyPin == newPin {
			s.logger.Warnf("New PIN found in history for user %s", user.ID.Hex())
			return fmt.Errorf("PIN_IN_HISTORY")
		}
	}
	return nil
}

// updateUserPin updates the user's PIN and PIN history
func (s *UserService) updateUserPin(ctx context.Context, userID, newPin string, user *User) error {
	// Update PIN history: move current PIN to history and add new PIN
	var newHistory [4]string
	copy(newHistory[1:], user.LoginPIN.PINHistory[:3])
	newHistory[0] = user.LoginPIN.PIN

	newLoginPIN := type_definition.LoginPIN{
		PIN:              newPin,
		PINHistory:       newHistory,
		LastPINCreatedAt: time.Now(),
	}

	if err := s.repository.ChangePin(ctx, userID, newLoginPIN); err != nil {
		s.logger.Errorf("Failed to update PIN for user %s: %v", userID, err)
		return fmt.Errorf("ERROR_SETTING_PIN")
	}

	return nil
}

func (s *UserService) Register(ctx context.Context, phone, deviceUUID, platform string) (*dto.RegisterResponse, error) {
	// Validate input parameters
	if err := s.validateRegistrationInputs(phone, deviceUUID, platform); err != nil {
		return nil, err
	}

	// Format phone number
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" {
		return nil, ErrInvalidPhoneNumber
	}

	// Check if user already exists
	existingUser, err := s.FindUserByPhone(ctx, formattedPhone)
	if err == nil && existingUser != nil {
		s.logger.Warnf("User with phone %s already exists", formattedPhone)
		return nil, ErrPhoneAlreadyExists
	}

	// Check if device is already registered
	existingDeviceUser, err := s.FindUserByDevice(ctx, deviceUUID)
	if err == nil && existingDeviceUser != nil {
		s.logger.Warnf("Device %s is already registered to user %s", deviceUUID, existingDeviceUser.ID.Hex())
		return nil, ErrDeviceAlreadyRegistered
	}

	// Check for pending registration
	pendingRegistration, err := s.FindPendingRegistration(ctx, formattedPhone, deviceUUID)
	if err == nil && pendingRegistration != nil {
		// Check if existing registration is still valid
		if time.Now().Before(pendingRegistration.ExpiresAt) && pendingRegistration.Status == "incomplete" {
			s.logger.Warnf("Registration already in progress for phone %s", formattedPhone)
			return nil, ErrRegistrationInProgress
		}
		// Delete expired registration
		if err := s.DeletePendingRegistration(ctx, pendingRegistration.ID); err != nil {
			s.logger.Errorf("Failed to delete expired registration: %v", err)
		}
	}

	// Generate OTP
	otpCode := utils.GenerateRandom(6)
	expirationTime := 10 * time.Minute
	wait := int(expirationTime.Minutes())

	// Encrypt OTP
	encOtpCode, _, err := utils.LocalEncryptPassword(otpCode, "otp", "", "", s.cfg)
	if err != nil {
		s.logger.Errorf("Failed to encrypt registration OTP: %v", err)
		return nil, ErrRegistrationFailed
	}

	// Create registration record
	registrationID := uuid.New().String()
	registration := RegistrationRecord{
		ID:          registrationID,
		PhoneNumber: formattedPhone,
		DeviceUUID:  deviceUUID,
		Platform:    platform,
		OTP:         encOtpCode,
		OTPFor:      "registration",
		Status:      "incomplete",
		ExpiresAt:   time.Now().Add(expirationTime),
		CreatedAt:   time.Now(),
		Attempts:    0,
		MaxAttempts: 3,
	}

	// Store registration record
	if err := s.CreatePendingRegistration(ctx, registration); err != nil {
		s.logger.Errorf("Failed to create registration record: %v", err)
		return nil, ErrRegistrationFailed
	}

	// Send OTP via SMS (async)
	go func() {
		message := fmt.Sprintf("Your CBE Super App registration OTP is: %s. Valid for %d minutes.", otpCode, wait)
		if err := utils.AxiosSendSms(ctx, formattedPhone, message); err != nil {
			s.logger.Errorf("Failed to send registration SMS: %v", err)
		}
	}()

	// Generate permanent token for registration using TokenMaker
	userEntity := entities.User{
		ID:          bson.NewObjectID(),
		FullName:    "",
		PhoneNumber: formattedPhone,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: deviceUUID,
		},
	}

	permissions := []string{"registration"}
	token, err := utils.TokenMaker(userEntity, permissions, "registration")
	if err != nil {
		s.logger.Errorf("Failed to generate registration token: %v", err)
		// Don't fail the registration if token generation fails
	}

	response := &dto.RegisterResponse{
		RegistrationID:     registrationID,
		PhoneNumber:        formattedPhone,
		DeviceUUID:         deviceUUID,
		Platform:           platform,
		OTPSent:            true,
		OTPExpiryMinutes:   wait,
		Token:              token,
		TokenType:          "permanent",
		TokenExpiry:        time.Now().Add(24 * time.Hour),
		NextStep:           "verify_otp",
		RegistrationStatus: "incomplete",
	}

	s.logger.Infof("Registration initiated successfully for phone %s, registration ID: %s", formattedPhone, registrationID)
	return response, nil
}

// validateRegistrationInputs validates the registration input parameters
func (s *UserService) validateRegistrationInputs(phone, deviceUUID, platform string) error {
	if phone == "" {
		return fmt.Errorf("INVALID_INPUT: phone number is required")
	}
	if deviceUUID == "" {
		return fmt.Errorf("INVALID_INPUT: device UUID is required")
	}
	if platform == "" {
		return fmt.Errorf("INVALID_INPUT: platform is required")
	}

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return fmt.Errorf("INVALID_PHONE_NUMBER")
	}

	// Validate platform
	validPlatforms := map[string]bool{"android": true, "ios": true, "web": true}
	if !validPlatforms[platform] {
		return fmt.Errorf("INVALID_PLATFORM")
	}

	// Validate device UUID format (basic validation)
	if len(deviceUUID) < 10 {
		return fmt.Errorf("INVALID_DEVICE_UUID")
	}

	return nil
}

// FindUserByPhone finds a user by phone number
func (s *UserService) FindUserByPhone(ctx context.Context, phone string) (*User, error) {
	return s.repository.FindUserByPhone(ctx, phone)
}

// FindUserByDevice finds a user by device UUID
func (s *UserService) FindUserByDevice(ctx context.Context, deviceUUID string) (*User, error) {
	return s.repository.FindUserByDevice(ctx, deviceUUID)
}

// FindPendingRegistration finds a pending registration by phone and device
func (s *UserService) FindPendingRegistration(ctx context.Context, phone, deviceUUID string) (*RegistrationRecord, error) {
	return s.repository.FindPendingRegistration(ctx, phone, deviceUUID)
}

// CreatePendingRegistration creates a new pending registration
func (s *UserService) CreatePendingRegistration(ctx context.Context, registration RegistrationRecord) error {
	return s.repository.CreatePendingRegistration(ctx, &registration)
}

// DeletePendingRegistration deletes a pending registration
func (s *UserService) DeletePendingRegistration(ctx context.Context, registrationID string) error {
	return s.repository.DeletePendingRegistration(ctx, registrationID)
}

// UpdatePendingRegistration updates a pending registration
func (s *UserService) UpdatePendingRegistration(ctx context.Context, registration RegistrationRecord) error {
	return s.repository.UpdatePendingRegistration(ctx, &registration)
}

func (s *UserService) Login(ctx context.Context, phone, deviceUUID, pin string) (*dto.LoginResponse, error) {
	// Validate input parameters
	if err := s.validateLoginInputs(phone, deviceUUID, pin); err != nil {
		return nil, err
	}

	// Format phone number
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" {
		return nil, ErrInvalidPhoneNumber
	}

	// Find user by phone number
	user, err := s.FindUserByPhone(ctx, formattedPhone)
	if err != nil {
		s.logger.Errorf("Failed to find user with phone %s: %v", formattedPhone, err)
		return nil, ErrUserNotFound
	}

	// Check if user is deleted
	if user.IsDeleted {
		s.logger.Warnf("Login attempt for deleted user %s", user.ID.Hex())
		return nil, ErrAccountDeleted
	}

	// Check if account is blocked
	if user.IsAccountBlocked {
		s.logger.Warnf("Login attempt for blocked user %s", user.ID.Hex())
		return nil, ErrAccountBlocked
	}

	// Check if device is linked to this user
	if user.Device.DeviceUUID != deviceUUID {
		s.logger.Warnf("Device mismatch for login. User %s, Expected: %s, Got: %s", user.ID.Hex(), user.Device.DeviceUUID, deviceUUID)
		return nil, ErrDeviceNotLinked
	}

	// Check if PIN is set
	if user.LoginPIN.PIN == "" {
		s.logger.Warnf("Login attempt for user %s without PIN", user.ID.Hex())
		return nil, ErrPinNotSet
	}

	// Check rate limiting
	if user.LoginAttemptCount >= 5 {
		// Check if enough time has passed since last attempt
		if time.Since(user.LastLoginAttempt) < 15*time.Minute {
			s.logger.Warnf("Too many login attempts for user %s", user.ID.Hex())
			return nil, ErrTooManyLoginAttempts
		}
		// Reset attempt count if enough time has passed
		if err := s.resetLoginAttempts(ctx, user.ID.Hex()); err != nil {
			s.logger.Errorf("Failed to reset login attempts: %v", err)
		}
	}

	// Validate PIN
	if user.LoginPIN.PIN != pin {
		// Increment login attempts
		if err := s.incrementLoginAttempts(ctx, user.ID.Hex()); err != nil {
			s.logger.Errorf("Failed to increment login attempts: %v", err)
		}
		s.logger.Warnf("Invalid PIN for user %s", user.ID.Hex())
		return nil, ErrInvalidPin
	}

	// Reset login attempts on successful login
	if err := s.resetLoginAttempts(ctx, user.ID.Hex()); err != nil {
		s.logger.Errorf("Failed to reset login attempts: %v", err)
	}

	// Update last login time
	if err := s.updateLastLogin(ctx, user.ID.Hex()); err != nil {
		s.logger.Errorf("Failed to update last login: %v", err)
	}

	// Generate permanent token using TokenMaker
	objectID, err := bson.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		s.logger.Errorf("Invalid user ID format: %v", err)
		return nil, fmt.Errorf("INVALID_USER_ID")
	}

	userEntity := entities.User{
		ID:          objectID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: deviceUUID,
		},
	}

	permissions := []string{"access", "transfer", "balance_check"}
	token, err := utils.TokenMaker(userEntity, permissions, "login")
	if err != nil {
		s.logger.Errorf("Failed to generate permanent token: %v", err)
		return nil, ErrTokenGenerationFailed
	}

	// Calculate session expiry
	sessionExpires := time.Now().Add(24 * time.Hour)

	// Prepare login response
	loginResponse := &dto.LoginResponse{
		Token:          token,
		UserID:         user.ID.Hex(),
		UserCode:       user.UserCode,
		FullName:       user.FullName,
		PhoneNumber:    user.PhoneNumber,
		KYCLevel:       user.KYC.KYCLevel,
		IsVerified:     user.IsVerified,
		DeviceUUID:     deviceUUID,
		LoginTime:      time.Now(),
		SessionExpires: sessionExpires,
		LastLogin:      user.LastLogin,
		LoginAttempts:  int(user.LoginAttemptCount),
	}

	s.logger.Infof("Login successful for user %s", user.ID.Hex())
	return loginResponse, nil
}

// validateLoginInputs validates the login input parameters
func (s *UserService) validateLoginInputs(phone, deviceUUID, pin string) error {
	if phone == "" {
		return fmt.Errorf("INVALID_INPUT: phone number is required")
	}
	if deviceUUID == "" {
		return fmt.Errorf("INVALID_INPUT: device UUID is required")
	}
	if pin == "" {
		return fmt.Errorf("INVALID_INPUT: PIN is required")
	}

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return ErrInvalidPhoneNumber
	}

	// Validate PIN format (6 digits)
	if len(pin) != 6 {
		return ErrInvalidPin
	}

	// Validate PIN contains only digits
	for _, char := range pin {
		if char < '0' || char > '9' {
			return ErrInvalidPin
		}
	}

	// Validate device UUID format (basic validation)
	if len(deviceUUID) < 10 {
		return ErrInvalidDeviceUUID
	}

	return nil
}

// generateAccessToken generates an access token for the user
func (s *UserService) generateAccessToken(ctx context.Context, user *User) (string, error) {
	// Create user entity for token generation
	userEntity := entities.User{
		ID:          user.ID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: user.Device.DeviceUUID,
		},
	}

	permissions := []string{"access"}
	token, err := utils.TokenMaker(userEntity, permissions, "login")
	if err != nil {
		return "", err
	}

	return token, nil
}

// incrementLoginAttempts increments the login attempt count for a user
func (s *UserService) incrementLoginAttempts(ctx context.Context, userID string) error {
	return s.repository.IncrementLoginAttempts(ctx, userID)
}

// resetLoginAttempts resets the login attempt count for a user
func (s *UserService) resetLoginAttempts(ctx context.Context, userID string) error {
	return s.repository.ResetLoginAttempts(ctx, userID)
}

// updateLastLogin updates the last login time for a user
func (s *UserService) updateLastLogin(ctx context.Context, userID string) error {
	return s.repository.UpdateLastLogin(ctx, userID)
}

func (s *UserService) ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) (*dto.ForgetPinSendOtpResponse, error) {
	// Validate input parameters
	if err := s.validateForgetPinInputs(phone, deviceUUID); err != nil {
		return nil, err
	}

	// Format phone number
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" {
		return nil, ErrInvalidPhoneNumber
	}

	// Find user by phone number
	user, err := s.FindUserByPhone(ctx, formattedPhone)
	if err != nil {
		s.logger.Errorf("Failed to find user with phone %s: %v", formattedPhone, err)
		return nil, ErrPinResetUserNotFound
	}

	// Check if user is deleted
	if user.IsDeleted {
		s.logger.Warnf("PIN reset attempt for deleted user %s", user.ID.Hex())
		return nil, ErrAccountDeleted
	}

	// Check if account is blocked
	if user.IsAccountBlocked {
		s.logger.Warnf("PIN reset attempt for blocked user %s", user.ID.Hex())
		return nil, ErrAccountBlocked
	}

	// Check if device is linked to this user
	if user.Device.DeviceUUID != deviceUUID {
		s.logger.Warnf("Device mismatch for PIN reset. User %s, Expected: %s, Got: %s", user.ID.Hex(), user.Device.DeviceUUID, deviceUUID)
		return nil, ErrPinResetDeviceMismatch
	}

	// Check if PIN reset is already in progress
	existingSession, err := s.repository.FindPinResetSessionByPhone(ctx, formattedPhone, deviceUUID)
	if err == nil && existingSession != nil {
		// Check if existing session is still valid (not expired)
		if time.Now().Before(existingSession.ExpiresAt) && existingSession.Status == "pending" {
			s.logger.Warnf("PIN reset already in progress for user %s", user.ID.Hex())
			return nil, ErrPinResetAlreadyInProgress
		}
		// Delete expired session
		if err := s.repository.DeletePinResetSession(ctx, existingSession.ID); err != nil {
			s.logger.Errorf("Failed to delete expired PIN reset session: %v", err)
		}
	}

	// Generate OTP
	otpCode := utils.GenerateRandom(6)
	expirationTime := 10 * time.Minute // 10 minutes expiry
	wait := int(expirationTime.Minutes())

	// Encrypt OTP
	encOtpCode, _, err := utils.LocalEncryptPassword(otpCode, "otp", "", "", s.cfg)
	if err != nil {
		s.logger.Errorf("Failed to encrypt PIN reset OTP: %v", err)
		return nil, ErrPinResetFailed
	}

	// Create PIN reset session
	sessionID := uuid.New().String()
	pinResetSession := &PinResetSession{
		ID:               sessionID,
		UserID:           user.ID.Hex(),
		PhoneNumber:      formattedPhone,
		DeviceUUID:       deviceUUID,
		OTP:              encOtpCode,
		OTPFor:           "pin_reset",
		Status:           "pending",
		ExpiresAt:        time.Now().Add(expirationTime),
		CreatedAt:        time.Now(),
		Attempts:         0,
		MaxAttempts:      3,
		AccessRestricted: true,
		Restrictions:     []string{"transfer", "balance_check"},
	}

	// Store PIN reset session
	if err := s.repository.CreatePinResetSession(ctx, pinResetSession); err != nil {
		s.logger.Errorf("Failed to create PIN reset session: %v", err)
		return nil, ErrPinResetFailed
	}

	// Send OTP via SMS (async)
	go func() {
		message := fmt.Sprintf("Your CBE Super App PIN reset OTP is: %s. Valid for %d minutes. Access will be restricted after reset.", otpCode, wait)
		if err := utils.AxiosSendSms(ctx, formattedPhone, message); err != nil {
			s.logger.Errorf("Failed to send PIN reset SMS: %v", err)
		}
	}()

	// Generate forget pin token using TempTokenMaker
	objectID, err := bson.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		s.logger.Errorf("Invalid user ID format: %v", err)
		return nil, fmt.Errorf("INVALID_USER_ID")
	}

	userEntity := &entities.User{
		ID:          objectID,
		FullName:    user.FullName,
		PhoneNumber: formattedPhone,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: deviceUUID,
		},
	}

	permissions := []string{"forget_pin", "verify_otp"}
	additional := map[string]interface{}{
		"reset_session_id": sessionID,
		"token_type":       "forget_pin",
	}

	token, err := utils.TempTokenMaker(userEntity, permissions, "pin_reset", additional, "forget_pin", s.cfg)
	if err != nil {
		s.logger.Errorf("Failed to generate forget pin token: %v", err)
		return nil, ErrPinResetFailed
	}

	// Prepare response
	response := &dto.ForgetPinSendOtpResponse{
		PhoneNumber:      formattedPhone,
		DeviceUUID:       deviceUUID,
		OTPSent:          true,
		OTPExpiryMinutes: wait,
		ResetSessionID:   sessionID,
		Token:            token,
		TokenType:        "forget_pin",
		TokenExpiry:      time.Now().Add(10 * time.Minute),
		NextStep:         "verify_otp",
	}

	s.logger.Infof("PIN reset OTP sent successfully for user %s, session: %s", user.ID.Hex(), sessionID)
	return response, nil
}

// ResetPin resets PIN using forget pin token and generates permanent token
func (s *UserService) ResetPin(ctx context.Context, resetSessionID, phone, deviceUUID, otp, newPin string) (*dto.ResetPinResponse, error) {
	// Validate input parameters
	if err := s.validateResetPinInputs(resetSessionID, phone, deviceUUID, otp, newPin); err != nil {
		return nil, err
	}

	// Format phone number
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" {
		return nil, ErrInvalidPhoneNumber
	}

	// Find PIN reset session
	session, err := s.repository.FindPinResetSession(ctx, resetSessionID)
	if err != nil {
		s.logger.Errorf("Failed to find PIN reset session %s: %v", resetSessionID, err)
		return nil, ErrPinResetSessionNotFound
	}

	// Verify session is still valid
	if time.Now().After(session.ExpiresAt) {
		s.logger.Warnf("PIN reset session %s has expired", resetSessionID)
		return nil, ErrPinResetSessionExpired
	}

	if session.Status != "pending" {
		s.logger.Warnf("PIN reset session %s is not in pending status: %s", resetSessionID, session.Status)
		return nil, ErrPinResetSessionInvalid
	}

	// Verify phone and device match
	if session.PhoneNumber != formattedPhone || session.DeviceUUID != deviceUUID {
		s.logger.Warnf("Phone or device mismatch for PIN reset session %s", resetSessionID)
		return nil, ErrPinResetSessionInvalid
	}

	// Check attempt limits
	if session.Attempts >= session.MaxAttempts {
		s.logger.Warnf("Too many attempts for PIN reset session %s", resetSessionID)
		return nil, ErrPinResetTooManyAttempts
	}

	// Verify OTP
	decryptedOTP, err := utils.LocalDecryptPassword(session.OTP, s.cfg)
	if err != nil {
		s.logger.Errorf("Failed to decrypt PIN reset OTP: %v", err)
		return nil, ErrPinResetFailed
	}

	if decryptedOTP != otp {
		// Increment attempts
		if err := s.repository.IncrementPinResetAttempts(ctx, resetSessionID); err != nil {
			s.logger.Errorf("Failed to increment PIN reset attempts: %v", err)
		}
		s.logger.Warnf("Invalid OTP for PIN reset session %s", resetSessionID)
		return nil, ErrPinResetOTPInvalid
	}

	// Find user
	user, err := s.FindUserByPhone(ctx, formattedPhone)
	if err != nil {
		s.logger.Errorf("Failed to find user for PIN reset: %v", err)
		return nil, ErrPinResetUserNotFound
	}

	// Check PIN history to prevent reuse
	if err := s.checkPinHistory(user, newPin); err != nil {
		s.logger.Warnf("PIN history check failed for user %s: %v", user.ID.Hex(), err)
		return nil, err
	}

	// Update user's PIN
	user.LoginPIN.PIN = newPin
	user.LoginPIN.LastPINCreatedAt = time.Now()

	// Update PIN history (shift and add new PIN)
	for i := len(user.LoginPIN.PINHistory) - 1; i > 0; i-- {
		user.LoginPIN.PINHistory[i] = user.LoginPIN.PINHistory[i-1]
	}
	user.LoginPIN.PINHistory[0] = newPin

	// Update user in database
	if err := s.repository.ChangePin(ctx, user.ID.Hex(), user.LoginPIN); err != nil {
		s.logger.Errorf("Failed to update user PIN: %v", err)
		return nil, ErrPinResetFailed
	}

	// Update session status to completed
	session.Status = "completed"
	session.VerifiedAt = time.Now()
	session.CompletedAt = time.Now()

	if err := s.repository.UpdatePinResetSession(ctx, session); err != nil {
		s.logger.Errorf("Failed to update PIN reset session: %v", err)
	}

	// Generate permanent token using TokenMaker
	objectID, err := bson.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		s.logger.Errorf("Invalid user ID format: %v", err)
		return nil, fmt.Errorf("INVALID_USER_ID")
	}

	userEntity := entities.User{
		ID:          objectID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: deviceUUID,
		},
	}

	permissions := []string{"access", "transfer", "balance_check"}
	token, err := utils.TokenMaker(userEntity, permissions, "permanent")
	if err != nil {
		s.logger.Errorf("Failed to generate permanent token: %v", err)
		return nil, fmt.Errorf("TOKEN_GENERATION_FAILED")
	}

	// Prepare response
	response := &dto.ResetPinResponse{
		UserID:           user.ID.Hex(),
		UserCode:         user.UserCode,
		FullName:         user.FullName,
		PhoneNumber:      user.PhoneNumber,
		PinReset:         true,
		ResetTime:        time.Now(),
		AccessRestricted: session.AccessRestricted,
		Restrictions:     session.Restrictions,
		Token:            token,
		TokenType:        "permanent",
		TokenExpiry:      time.Now().Add(24 * time.Hour),
		NextStep:         "login_with_new_pin",
	}

	s.logger.Infof("PIN reset completed successfully for user %s, session: %s", user.ID.Hex(), resetSessionID)
	return response, nil
}

// validateForgetPinInputs validates the forget PIN input parameters
func (s *UserService) validateForgetPinInputs(phone, deviceUUID string) error {
	if phone == "" {
		return fmt.Errorf("INVALID_INPUT: phone number is required")
	}
	if deviceUUID == "" {
		return fmt.Errorf("INVALID_INPUT: device UUID is required")
	}

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return ErrInvalidPhoneNumber
	}

	// Validate device UUID format (basic validation)
	if len(deviceUUID) < 10 {
		return ErrInvalidDeviceUUID
	}

	return nil
}

// validateResetPinInputs validates the reset PIN input parameters
func (s *UserService) validateResetPinInputs(resetSessionID, phone, deviceUUID, otp, newPin string) error {
	if resetSessionID == "" {
		return fmt.Errorf("INVALID_INPUT: reset session ID is required")
	}
	if phone == "" {
		return fmt.Errorf("INVALID_INPUT: phone number is required")
	}
	if deviceUUID == "" {
		return fmt.Errorf("INVALID_INPUT: device UUID is required")
	}
	if otp == "" {
		return fmt.Errorf("INVALID_INPUT: OTP is required")
	}
	if newPin == "" {
		return fmt.Errorf("INVALID_INPUT: new PIN is required")
	}

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return ErrInvalidPhoneNumber
	}

	// Validate OTP format (6 digits)
	if len(otp) != 6 {
		return ErrPinResetOTPInvalid
	}

	// Validate OTP contains only digits
	for _, char := range otp {
		if char < '0' || char > '9' {
			return ErrPinResetOTPInvalid
		}
	}

	// Validate new PIN format (6 digits)
	if len(newPin) != 6 {
		return ErrInvalidPin
	}

	// Validate new PIN contains only digits
	for _, char := range newPin {
		if char < '0' || char > '9' {
			return ErrInvalidPin
		}
	}

	// Validate device UUID format (basic validation)
	if len(deviceUUID) < 10 {
		return ErrInvalidDeviceUUID
	}

	return nil
}

func (s *UserService) DeviceLookup(ctx context.Context, deviceUUID, platform, appVersion string) (*dto.DeviceLookupResponse, error) {
	// Validate input parameters
	if deviceUUID == "" {
		return nil, fmt.Errorf("INVALID_INPUT: device UUID is required")
	}
	if platform == "" {
		return nil, fmt.Errorf("INVALID_INPUT: platform is required")
	}

	// Check if device is linked to any user
	user, err := s.FindUserByDevice(ctx, deviceUUID)
	userFound := err == nil && user != nil

	// Generate device lookup token using TempTokenMaker
	userEntity := &entities.User{
		ID:          bson.NewObjectID(),
		FullName:    "",
		PhoneNumber: "",
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: deviceUUID,
			AppVersion: appVersion,
		},
	}

	permissions := []string{"device_lookup"}
	additional := map[string]interface{}{
		"platform":   platform,
		"token_type": "device_lookup",
	}

	token, err := utils.TempTokenMaker(userEntity, permissions, "device_lookup", additional, "device_lookup", s.cfg)
	if err != nil {
		s.logger.Errorf("Failed to generate device lookup token: %v", err)
		return nil, fmt.Errorf("TOKEN_GENERATION_FAILED")
	}

	// Get HQ data for app version check
	hqData, err := s.repository.GetOneHQ(ctx, map[string]interface{}{})
	isLatest := true
	if err == nil && hqData != nil {
		isLatest = s.isAppVersionLatest(platform, appVersion, hqData)
	}

	// Determine next step based on user found status
	nextStep := "register"
	if userFound {
		nextStep = "login"
	}

	response := &dto.DeviceLookupResponse{
		DeviceUUID:  deviceUUID,
		Platform:    platform,
		AppVersion:  appVersion,
		IsLatest:    isLatest,
		UserFound:   userFound,
		Token:       token,
		TokenType:   "device_lookup",
		TokenExpiry: time.Now().Add(5 * time.Minute),
		NextStep:    nextStep,
	}

	// Generate OTP and include in response for dev and uat environments when user is found
	if userFound && (s.cfg.GoEnv == "dev" || s.cfg.GoEnv == "uat") {
		otpCode := utils.OTPGenerator(6)
		response.OTPCode = otpCode

		// Encrypt OTP for storage
		encOtpCode, _, err := utils.LocalEncryptPassword(otpCode, "otp", "", "", s.cfg)
		if err != nil {
			s.logger.Errorf("Failed to encrypt OTP: %v", err)
			return nil, fmt.Errorf("OTP_ENCRYPTION_FAILED")
		}

		// Calculate OTP expiration time
		wait, err := strconv.Atoi(s.cfg.OtpWaitingTime)
		if err != nil {
			wait = 10 // Default to 10 minutes
		}
		expirationTime := time.Duration(wait) * time.Minute

		// Create OTP record
		otpRecord := OTPRecord{
			UserID:     user.ID.Hex(),
			OTP:        encOtpCode,
			OTPFor:     "device_lookup",
			UserRealm:  "member",
			ExpiresAt:  time.Now().Add(expirationTime),
			CreatedAt:  time.Now(),
			DeviceUUID: &deviceUUID,
		}

		// Store OTP in database
		if err := s.CreateOtp(ctx, otpRecord); err != nil {
			s.logger.Errorf("Failed to create OTP record: %v", err)
			return nil, fmt.Errorf("OTP_CREATION_FAILED")
		}

		// Send OTP via SMS (async)
		go func() {
			message := fmt.Sprintf("Your device lookup OTP is: %s. Valid for %d minutes.", otpCode, wait)
			if err := utils.AxiosSendSms(ctx, user.PhoneNumber, message); err != nil {
				s.logger.Errorf("Failed to send SMS: %v", err)
			}
		}()

		s.logger.Infof("OTP generated and sent for device lookup - Device: %s, User: %s, OTP: %s", deviceUUID, user.PhoneNumber, otpCode)
	}

	s.logger.Infof("Device lookup completed for device %s, user found: %v", deviceUUID, userFound)
	return response, nil
}

// Helper method to check app version
func (s *UserService) isAppVersionLatest(platform, appVersion string, hqData *HQ) bool {
	if hqData == nil {
		return true
	}

	switch platform {
	case "android":
		return appVersion >= hqData.LatestAndroidVersion
	case "ios":
		return appVersion >= hqData.LatestiOSVersion
	default:
		return true
	}
}

// validateVerifyOtpInputs validates the verify OTP input parameters
func (s *UserService) validateVerifyOtpInputs(userID, otp, userRealm, otpFor string) error {
	if otp == "" {
		return fmt.Errorf("INVALID_INPUT: OTP is required")
	}
	if userRealm == "" {
		return fmt.Errorf("INVALID_INPUT: user realm is required")
	}
	if otpFor == "" {
		return fmt.Errorf("INVALID_INPUT: OTP purpose is required")
	}

	// Validate OTP format (6 digits)
	if len(otp) != 6 {
		return fmt.Errorf("INVALID_OTP: OTP must be exactly 6 digits")
	}

	// Validate OTP contains only digits
	for _, char := range otp {
		if char < '0' || char > '9' {
			return fmt.Errorf("INVALID_OTP: OTP must contain only digits")
		}
	}

	return nil
}

// CompleteRegistration completes registration and generates permanent token
func (s *UserService) CompleteRegistration(ctx context.Context, registrationID, phone, deviceUUID, platform, fullName string) (*dto.CompleteRegistrationResponse, error) {
	// Find the pending registration
	registration, err := s.repository.FindPendingRegistrationByID(ctx, registrationID)
	if err != nil {
		s.logger.Errorf("Failed to find registration %s: %v", registrationID, err)
		return nil, fmt.Errorf("REGISTRATION_NOT_FOUND")
	}

	// Verify the registration is still valid
	if time.Now().After(registration.ExpiresAt) {
		s.logger.Warnf("Registration %s has expired", registrationID)
		return nil, fmt.Errorf("REGISTRATION_EXPIRED")
	}

	if registration.Status != "incomplete" {
		s.logger.Warnf("Registration %s is not in incomplete status: %s", registrationID, registration.Status)
		return nil, fmt.Errorf("INVALID_REGISTRATION_STATUS")
	}

	// Update registration status to complete in otp_collection
	registration.Status = "complete"

	if err := s.repository.UpdatePendingRegistration(ctx, registration); err != nil {
		s.logger.Errorf("Failed to update registration status: %v", err)
		return nil, fmt.Errorf("REGISTRATION_UPDATE_FAILED")
	}

	// Create user in user_collection with kyc_level: 0
	userCode := utils.GenerateRandom(8) // Generate 8-digit user code
	user := &User{
		ID:          bson.NewObjectID(),
		UserCode:    userCode,
		FullName:    fullName,
		PhoneNumber: phone,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: deviceUUID,
			AppVersion: "",
		},
		KYC: struct {
			KYCRejectReasonField map[string]struct{} `json:"kyc_reject_reason_failed" bson:"kyc_reject_reason_failed"`
			KYCStatus            enums.KYCStatus     `json:"kyc_status" bson:"kyc_status"`
			KYCRejectReason      string              `json:"kyc_reject_reason" bson:"kyc_reject_reason"`
			KYCApproved          bool                `json:"kyc_approved" bson:"kyc_approved"`
			KYCActivityBy        map[string]struct{} `json:"kyc_activity_by" bson:"kyc_activity_by"`
			KYCLevel             uint8               `json:"level" bson:"level"`
		}{
			KYCLevel: 0, // Set KYC level to 0 as specified
		},
		IsVerified:     false,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
		LoginPIN: type_definition.LoginPIN{
			PIN:              "",
			PINHistory:       [4]string{},
			LastPINCreatedAt: time.Time{},
		},
	}

	// Save user to user_collection
	if err := s.repository.CreateUser(ctx, user); err != nil {
		s.logger.Errorf("Failed to create user: %v", err)
		return nil, fmt.Errorf("USER_CREATION_FAILED")
	}

	// Generate permanent token using TokenMaker
	objectID, err := bson.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		s.logger.Errorf("Invalid user ID format: %v", err)
		return nil, fmt.Errorf("INVALID_USER_ID")
	}

	userEntity := entities.User{
		ID:          objectID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Device: struct {
			DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
			AppVersion string `json:"app_version" bson:"app_version"`
		}{
			DeviceUUID: deviceUUID,
		},
	}

	permissions := []string{"access", "transfer", "balance_check"}
	token, err := utils.TokenMaker(userEntity, permissions, "permanent")
	if err != nil {
		s.logger.Errorf("Failed to generate permanent token: %v", err)
		return nil, fmt.Errorf("TOKEN_GENERATION_FAILED")
	}

	response := &dto.CompleteRegistrationResponse{
		UserID:      user.ID.Hex(),
		UserCode:    user.UserCode,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		KYCLevel:    user.KYC.KYCLevel,
		Status:      "complete",
		Token:       token,
		TokenType:   "permanent",
		TokenExpiry: time.Now().Add(24 * time.Hour),
		NextStep:    "set_pin",
	}

	s.logger.Infof("Registration completed successfully for phone %s, user ID: %s, KYC Level: 0", phone, user.ID.Hex())
	return response, nil
}

// ValidateRegistrationInputs validates the registration input parameters (public for testing)
func (s *UserService) ValidateRegistrationInputs(phone, deviceUUID, platform string) error {
	return s.validateRegistrationInputs(phone, deviceUUID, platform)
}
