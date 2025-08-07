package redis

import (
	"context"
	"fmt"
	"time"

	"cbe-super-app-member-auth/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// UserSessionRepository implements the UserSessionRepository interface
type UserSessionRepository struct {
	redisRepo storage.RedisRepository
	logger    utils.Logger
}

// NewUserSessionRepository creates a new user session repository instance
func NewUserSessionRepository(redisRepo storage.RedisRepository, logger utils.Logger) storage.UserSessionRepository {
	return &UserSessionRepository{
		redisRepo: redisRepo,
		logger:    logger,
	}
}

// SaveUserSession saves user session data
func (u *UserSessionRepository) SaveUserSession(ctx context.Context, sessionID string, userData map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("user_session:%s", sessionID)

	for field, value := range userData {
		err := u.redisRepo.HSet(ctx, key, field, value)
		if err != nil {
			u.logger.Errorf("Failed to save user session field %s for session %s: %v", field, sessionID, err)
			return err
		}
	}

	// Set expiration for the entire hash
	err := u.redisRepo.Expire(ctx, key, expiration)
	if err != nil {
		u.logger.Errorf("Failed to set expiration for user session %s: %v", sessionID, err)
		return err
	}

	u.logger.Infof("Successfully saved user session: %s", sessionID)
	return nil
}

// GetUserSession retrieves user session data
func (u *UserSessionRepository) GetUserSession(ctx context.Context, sessionID string) (map[string]string, error) {
	key := fmt.Sprintf("user_session:%s", sessionID)

	sessionData, err := u.redisRepo.HGetAll(ctx, key)
	if err != nil {
		u.logger.Errorf("Failed to get user session %s: %v", sessionID, err)
		return nil, err
	}

	u.logger.Infof("Successfully retrieved user session: %s", sessionID)
	return sessionData, nil
}

// DeleteUserSession deletes user session data
func (u *UserSessionRepository) DeleteUserSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("user_session:%s", sessionID)

	err := u.redisRepo.Delete(ctx, key)
	if err != nil {
		u.logger.Errorf("Failed to delete user session %s: %v", sessionID, err)
		return err
	}

	u.logger.Infof("Successfully deleted user session: %s", sessionID)
	return nil
}

// UpdateUserSession updates specific fields in user session data
func (u *UserSessionRepository) UpdateUserSession(ctx context.Context, sessionID string, updates map[string]interface{}) error {
	key := fmt.Sprintf("user_session:%s", sessionID)

	for field, value := range updates {
		err := u.redisRepo.HSet(ctx, key, field, value)
		if err != nil {
			u.logger.Errorf("Failed to update user session field %s for session %s: %v", field, sessionID, err)
			return err
		}
	}

	u.logger.Infof("Successfully updated user session: %s", sessionID)
	return nil
}

// SaveUserDeviceMapping saves user device mapping data
func (u *UserSessionRepository) SaveUserDeviceMapping(ctx context.Context, userID string, deviceUUID string, sessionData map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("user_device:%s", userID)

	// Save device UUID as a field
	err := u.redisRepo.HSet(ctx, key, "device_uuid", deviceUUID)
	if err != nil {
		u.logger.Errorf("Failed to save device UUID for user %s: %v", userID, err)
		return err
	}

	// Save additional session data
	for field, value := range sessionData {
		err := u.redisRepo.HSet(ctx, key, field, value)
		if err != nil {
			u.logger.Errorf("Failed to save session data field %s for user %s: %v", field, userID, err)
			return err
		}
	}

	// Set expiration
	err = u.redisRepo.Expire(ctx, key, expiration)
	if err != nil {
		u.logger.Errorf("Failed to set expiration for user device mapping %s: %v", userID, err)
		return err
	}

	u.logger.Infof("Successfully saved user device mapping: %s", userID)
	return nil
}

// GetUserDeviceMapping retrieves user device mapping data
func (u *UserSessionRepository) GetUserDeviceMapping(ctx context.Context, userID string) (map[string]string, error) {
	key := fmt.Sprintf("user_device:%s", userID)

	deviceData, err := u.redisRepo.HGetAll(ctx, key)
	if err != nil {
		u.logger.Errorf("Failed to get user device mapping %s: %v", userID, err)
		return nil, err
	}

	u.logger.Infof("Successfully retrieved user device mapping: %s", userID)
	return deviceData, nil
}

// DeleteUserDeviceMapping deletes user device mapping data
func (u *UserSessionRepository) DeleteUserDeviceMapping(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user_device:%s", userID)

	err := u.redisRepo.Delete(ctx, key)
	if err != nil {
		u.logger.Errorf("Failed to delete user device mapping %s: %v", userID, err)
		return err
	}

	u.logger.Infof("Successfully deleted user device mapping: %s", userID)
	return nil
}

// SetUserOnline sets user online status
func (u *UserSessionRepository) SetUserOnline(ctx context.Context, userID string, deviceUUID string, expiration time.Duration) error {
	// Add user to online users set
	err := u.redisRepo.SAdd(ctx, "online_users", userID)
	if err != nil {
		u.logger.Errorf("Failed to add user %s to online users: %v", userID, err)
		return err
	}

	// Set user online status with device info
	key := fmt.Sprintf("user_online:%s", userID)
	onlineData := map[string]interface{}{
		"device_uuid": deviceUUID,
		"online_at":   time.Now().Unix(),
		"status":      "online",
	}

	for field, value := range onlineData {
		err := u.redisRepo.HSet(ctx, key, field, value)
		if err != nil {
			u.logger.Errorf("Failed to set online status field %s for user %s: %v", field, userID, err)
			return err
		}
	}

	// Set expiration
	err = u.redisRepo.Expire(ctx, key, expiration)
	if err != nil {
		u.logger.Errorf("Failed to set expiration for user online status %s: %v", userID, err)
		return err
	}

	u.logger.Infof("Successfully set user online: %s", userID)
	return nil
}

// SetUserOffline sets user offline status
func (u *UserSessionRepository) SetUserOffline(ctx context.Context, userID string, deviceUUID string) error {
	// Remove user from online users set
	err := u.redisRepo.SRem(ctx, "online_users", userID)
	if err != nil {
		u.logger.Errorf("Failed to remove user %s from online users: %v", userID, err)
		return err
	}

	// Update user online status
	key := fmt.Sprintf("user_online:%s", userID)
	offlineData := map[string]interface{}{
		"device_uuid": deviceUUID,
		"offline_at":  time.Now().Unix(),
		"status":      "offline",
	}

	for field, value := range offlineData {
		err := u.redisRepo.HSet(ctx, key, field, value)
		if err != nil {
			u.logger.Errorf("Failed to set offline status field %s for user %s: %v", field, userID, err)
			return err
		}
	}

	u.logger.Infof("Successfully set user offline: %s", userID)
	return nil
}

// IsUserOnline checks if a user is online
func (u *UserSessionRepository) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	isMember, err := u.redisRepo.SIsMember(ctx, "online_users", userID)
	if err != nil {
		u.logger.Errorf("Failed to check if user %s is online: %v", userID, err)
		return false, err
	}

	u.logger.Infof("User %s online status: %t", userID, isMember)
	return isMember, nil
}

// GetOnlineUsers gets all online users
func (u *UserSessionRepository) GetOnlineUsers(ctx context.Context) ([]string, error) {
	onlineUsers, err := u.redisRepo.SMembers(ctx, "online_users")
	if err != nil {
		u.logger.Errorf("Failed to get online users: %v", err)
		return nil, err
	}

	u.logger.Infof("Retrieved %d online users", len(onlineUsers))
	return onlineUsers, nil
}
