package users

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	local_utils "cbe-super-app-member-users/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UserService struct {
	// repository userRepoPort.UserRepositoryPort
	repository UserRepository
	logger     utils.Logger
	minIO      config.MinioClientInterface
	cfg        *config.VaultConfig
}

func NewUserService(repository UserRepository, logger utils.Logger, minIO config.MinioClientInterface, cfg *config.VaultConfig) *UserService {
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
			UserID:         user.ID,
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
		UserID:         user.ID,
		FullName:       user.FullName,
		LinkedAccounts: linkedAccountsResponse,
	}

	s.logger.Infof("Successfully fetched active linked accounts for user ID %s", id)
	return response, nil
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

	encryptedOTPCode, _, err := local_utils.LocalEncryptPassword(otp, "otp", "", "", s.cfg)
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

	decryptedOTP, err := local_utils.LocalDecryptPassword(record.OTP, s.cfg)
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
	if err := s.repository.DeleteOtp(ctx, record.ID); err != nil {
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
	if currentUser.Device != nil && currentUser.Device.DeviceUUID == DeviceID {

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
s.logger.Infof(userData.LoginPIN.PIN ,ChangePinRequest.OldPin)
	
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
	newLoginPIN := LoginPIN{
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
