package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// OTPSessionRepository implements the OTPSessionRepository interface
type OTPSessionRepository struct {
	redisRepo storage.RedisRepository
	logger    utils.Logger
}

// NewOTPSessionRepository creates a new OTP session repository instance
func NewOTPSessionRepository(redisRepo storage.RedisRepository, logger utils.Logger) storage.OTPSessionRepository {
	return &OTPSessionRepository{
		redisRepo: redisRepo,
		logger:    logger,
	}
}

// SaveOTP saves an OTP for a phone number
func (o *OTPSessionRepository) SaveOTP(ctx context.Context, phoneNumber string, otp string, expiration time.Duration) error {
	key := fmt.Sprintf("otp:%s", phoneNumber)

	err := o.redisRepo.Set(ctx, key, otp, expiration)
	if err != nil {
		o.logger.Errorf("Failed to save OTP for phone number %s: %v", phoneNumber, err)
		return err
	}

	o.logger.Infof("Successfully saved OTP for phone number: %s", phoneNumber)
	return nil
}

// GetOTP retrieves an OTP for a phone number
func (o *OTPSessionRepository) GetOTP(ctx context.Context, phoneNumber string) (string, error) {
	key := fmt.Sprintf("otp:%s", phoneNumber)

	otp, err := o.redisRepo.Get(ctx, key)
	if err != nil {
		o.logger.Errorf("Failed to get OTP for phone number %s: %v", phoneNumber, err)
		return "", err
	}

	o.logger.Infof("Successfully retrieved OTP for phone number: %s", phoneNumber)
	return otp, nil
}

// DeleteOTP deletes an OTP for a phone number
func (o *OTPSessionRepository) DeleteOTP(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("otp:%s", phoneNumber)

	err := o.redisRepo.Delete(ctx, key)
	if err != nil {
		o.logger.Errorf("Failed to delete OTP for phone number %s: %v", phoneNumber, err)
		return err
	}

	o.logger.Infof("Successfully deleted OTP for phone number: %s", phoneNumber)
	return nil
}

// IncrementOTPAttempts increments the OTP attempt counter for a phone number
func (o *OTPSessionRepository) IncrementOTPAttempts(ctx context.Context, phoneNumber string) (int, error) {
	key := fmt.Sprintf("otp_attempts:%s", phoneNumber)

	// Get current attempts
	currentAttemptsStr, err := o.redisRepo.Get(ctx, key)
	if err != nil {
		// If key doesn't exist, start with 1
		err = o.redisRepo.Set(ctx, key, "1", 24*time.Hour) // Reset after 24 hours
		if err != nil {
			o.logger.Errorf("Failed to set initial OTP attempts for phone number %s: %v", phoneNumber, err)
			return 0, err
		}
		o.logger.Infof("Set initial OTP attempts for phone number: %s", phoneNumber)
		return 1, nil
	}

	currentAttempts, err := strconv.Atoi(currentAttemptsStr)
	if err != nil {
		o.logger.Errorf("Failed to parse OTP attempts for phone number %s: %v", phoneNumber, err)
		return 0, err
	}

	newAttempts := currentAttempts + 1
	err = o.redisRepo.Set(ctx, key, strconv.Itoa(newAttempts), 24*time.Hour)
	if err != nil {
		o.logger.Errorf("Failed to increment OTP attempts for phone number %s: %v", phoneNumber, err)
		return 0, err
	}

	o.logger.Infof("Incremented OTP attempts for phone number %s: %d", phoneNumber, newAttempts)
	return newAttempts, nil
}

// GetOTPAttempts gets the current OTP attempt count for a phone number
func (o *OTPSessionRepository) GetOTPAttempts(ctx context.Context, phoneNumber string) (int, error) {
	key := fmt.Sprintf("otp_attempts:%s", phoneNumber)

	attemptsStr, err := o.redisRepo.Get(ctx, key)
	if err != nil {
		o.logger.Warnf("No OTP attempts found for phone number %s", phoneNumber)
		return 0, nil
	}

	attempts, err := strconv.Atoi(attemptsStr)
	if err != nil {
		o.logger.Errorf("Failed to parse OTP attempts for phone number %s: %v", phoneNumber, err)
		return 0, err
	}

	o.logger.Infof("Retrieved OTP attempts for phone number %s: %d", phoneNumber, attempts)
	return attempts, nil
}

// ResetOTPAttempts resets the OTP attempt counter for a phone number
func (o *OTPSessionRepository) ResetOTPAttempts(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("otp_attempts:%s", phoneNumber)

	err := o.redisRepo.Delete(ctx, key)
	if err != nil {
		o.logger.Errorf("Failed to reset OTP attempts for phone number %s: %v", phoneNumber, err)
		return err
	}

	o.logger.Infof("Successfully reset OTP attempts for phone number: %s", phoneNumber)
	return nil
}

// SaveOTPVerification saves OTP verification status
func (o *OTPSessionRepository) SaveOTPVerification(ctx context.Context, phoneNumber string, verified bool, expiration time.Duration) error {
	key := fmt.Sprintf("otp_verified:%s", phoneNumber)

	err := o.redisRepo.Set(ctx, key, strconv.FormatBool(verified), expiration)
	if err != nil {
		o.logger.Errorf("Failed to save OTP verification for phone number %s: %v", phoneNumber, err)
		return err
	}

	o.logger.Infof("Successfully saved OTP verification for phone number: %s", phoneNumber)
	return nil
}

// GetOTPVerification gets the OTP verification status
func (o *OTPSessionRepository) GetOTPVerification(ctx context.Context, phoneNumber string) (bool, error) {
	key := fmt.Sprintf("otp_verified:%s", phoneNumber)

	verifiedStr, err := o.redisRepo.Get(ctx, key)
	if err != nil {
		o.logger.Warnf("No OTP verification found for phone number %s", phoneNumber)
		return false, nil
	}

	verified, err := strconv.ParseBool(verifiedStr)
	if err != nil {
		o.logger.Errorf("Failed to parse OTP verification for phone number %s: %v", phoneNumber, err)
		return false, err
	}

	o.logger.Infof("Retrieved OTP verification for phone number %s: %t", phoneNumber, verified)
	return verified, nil
}

// DeleteOTPVerification deletes the OTP verification status
func (o *OTPSessionRepository) DeleteOTPVerification(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("otp_verified:%s", phoneNumber)

	err := o.redisRepo.Delete(ctx, key)
	if err != nil {
		o.logger.Errorf("Failed to delete OTP verification for phone number %s: %v", phoneNumber, err)
		return err
	}

	o.logger.Infof("Successfully deleted OTP verification for phone number: %s", phoneNumber)
	return nil
}
