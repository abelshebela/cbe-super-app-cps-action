package notification

import (
	"context"
	"fmt"
	"time"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	inappnotification "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/in-app-notification"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, notification NotificationRequest) (*Notification, error)
	UpdateNotification(ctx context.Context, id string, notification NotificationRequest) (*Notification, *Notification, error)
	DeleteNotification(ctx context.Context, id string) (*Notification, *Notification, error)
	EnableDisableNotification(ctx context.Context, id string, enable bool) (*Notification, *Notification, error)
	FetchNotificationByID(ctx context.Context, id string) (*Notification, error)
	FetchNotifications(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*Notification], error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

// Service implements NotificationService
type Service struct {
	Repository                  NotificationRepository
	InAppNotificationRepository inappnotification.InAppNotificationRepository
	logger                      shared_utils.Logger
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(repository NotificationRepository, inAppNotificationRepository inappnotification.InAppNotificationRepository, logger shared_utils.Logger) NotificationService {
	return &Service{
		Repository:                  repository,
		InAppNotificationRepository: inAppNotificationRepository,
		logger:                      logger,
	}
}

func (s *Service) CreateNotification(ctx context.Context, notification NotificationRequest) (*Notification, error) {
	s.logger.Infof("Creating notification", "type", notification.NotificationType)

	// Validate notification type and recipient
	exist, err := s.Repository.NotificationExists(ctx, notification.NotificationType, NotificationFor(notification.For), nil)
	if err != nil {
		s.logger.Errorf("Failed to check if notification exists", "error", err)
		return nil, err
	}

	if exist {
		s.logger.Errorf("Notification already exists", "type", notification.NotificationType, "for", notification.For)
		return nil, fmt.Errorf("NOTIFICATION_ALREADY_EXISTS")
	}

	result := Notification{
		NotificationType: notification.NotificationType,
		NotificationBody: notification.NotificationBody,
		IsPublic:         notification.IsPublic,
		For:              NotificationFor(notification.For),
		CreatedBy:        notification.CreatedBy,
		Title:            notification.Title,
		Status:           StatusPending,
		Seen:             false,
		Enabled:          false,
		IsDeleted:        false,
		CreatedAt:        time.Now(),
		LastModified:     time.Now(),
	}

	s.logger.Infof("Notification constructed successfully", "id", result.ID)
	return &result, nil
}

// UpdateNotification updates an existing notification without persisting
func (s *Service) UpdateNotification(ctx context.Context, id string, notification NotificationRequest) (*Notification, *Notification, error) {
	s.logger.Infof("Updating notification", "id", id)

	// Fetch previous notification
	prevNotification, err := s.Repository.FetchNotificationByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch notification", "id", id, "error", err)
		return nil, nil, err
	}

	// Check if notification type and recipient combination already exists
	if notification.NotificationType != "" && notification.For != "" {
		exist, err := s.Repository.NotificationExists(ctx, notification.NotificationType, NotificationFor(notification.For), &id)
		if err != nil {
			s.logger.Errorf("Failed to check if notification exists", "error", err)
			return nil, nil, err
		}
		if exist {
			s.logger.Errorf("Notification already exists", "type", notification.NotificationType, "for", notification.For)
			return nil, nil, fmt.Errorf("NOTIFICATION_ALREADY_EXISTS")
		}
	}

	curNotification := Notification{
		ID:               prevNotification.ID,
		NotificationType: nonEmptyString(notification.NotificationType, prevNotification.NotificationType),
		NotificationBody: nonEmptyString(notification.NotificationBody, prevNotification.NotificationBody),
		Title:            nonEmptyString(notification.Title, prevNotification.Title),
		IsPublic:         notification.IsPublic || prevNotification.IsPublic,
		For:              nonEmptyNotificationFor(notification.For, prevNotification.For),
		CreatedBy:        nonEmptyString(notification.CreatedBy, prevNotification.CreatedBy),
		Seen:             prevNotification.Seen,
		Enabled:          prevNotification.Enabled,
		IsDeleted:        prevNotification.IsDeleted,
		CreatedAt:        prevNotification.CreatedAt,
		LastModified:     time.Now(),
	}

	s.logger.Infof("Notification updated successfully", "id", id)
	return &curNotification, prevNotification, nil
}

// DeleteNotification soft-deletes a notification without persisting
func (s *Service) DeleteNotification(ctx context.Context, id string) (*Notification, *Notification, error) {
	s.logger.Infof("Deleting notification", "id", id)

	// Fetch previous notification
	prevNotification, err := s.Repository.FetchNotificationByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch notification", "id", id, "error", err)
		return nil, nil, err
	}

	curNotification := generateNotification(*prevNotification)
	curNotification.IsDeleted = true
	curNotification.DeletedAt = time.Now()
	curNotification.LastModified = time.Now()

	s.logger.Infof("Notification marked for deletion", "id", id)
	return curNotification, prevNotification, nil
}

// EnableDisableNotification enables or disables a notification without persisting
func (s *Service) EnableDisableNotification(ctx context.Context, id string, enable bool) (*Notification, *Notification, error) {
	s.logger.Infof("EnableDisable notification", "id", id, "enable", enable)

	// Fetch previous notification
	prevNotification, err := s.Repository.FetchNotificationByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch notification", "id", id, "error", err)
		return nil, nil, err
	}

	if enable && prevNotification.Enabled {
		s.logger.Errorf("Notification already enabled", "id", id)
		return nil, nil, fmt.Errorf(common_util.ErrAlreadyEnabled)
	}

	if !enable && !prevNotification.Enabled {
		s.logger.Errorf("Notification already disabled", "id", id)
		return nil, nil, fmt.Errorf(common_util.ErrAlreadyDisabled)
	}

	curNotification := generateNotification(*prevNotification)
	curNotification.Enabled = enable
	curNotification.LastModified = time.Now()

	s.logger.Infof("Notification enable/disable prepared", "id", id, "enable", enable)
	return curNotification, prevNotification, nil
}

// Authorize handles persistence for notification actions
func (s *Service) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	requestedAction := action.RequestAction

	var notification *Notification
	var err error

	bindErr := common_util.BindAction(action.CurrentAction, &notification)
	if bindErr != nil {
		s.logger.Errorf("Failed to bind current action to notification: %v", bindErr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	switch requestedAction {
	case cps_const.RequestCreateNotification:
		err = s.Repository.RunInTransaction(ctx, func(ctx context.Context) error {
			notification.Status = StatusDelivered
			notification, err = s.Repository.CreateNotification(ctx, notification)
			if err != nil {
				s.logger.Errorf("Failed to create notification in repository", "error", err)
				return err
			}

			inAppNotification := inappnotification.InAppNotification{
				NotificationType: notification.NotificationType,
				NotificationBody: notification.NotificationBody,
				IsPublic:         notification.IsPublic,
				For:              string(notification.For),
				CreatedBy:        notification.CreatedBy,
			}

			err = s.InAppNotificationRepository.CreateInAppNotification(ctx, &inAppNotification)
			if err != nil {
				s.logger.Errorf("Failed to create in-app notification in repository", "error", err)
				return err
			}
			return nil
		})
		if err != nil {
			s.logger.Errorf("Failed to create notification in repository", "error", err)
			return nil, err
		}

	case cps_const.RequestUpdateNotification:
		notification, err = s.Repository.UpdateNotification(ctx, notification)
		if err != nil {
			s.logger.Errorf("Failed to update notification in repository", "error", err)
			return nil, err
		}
	case cps_const.RequestDeleteNotification:
		err = s.Repository.RunInTransaction(ctx, func(ctx context.Context) error {
			_, err := s.Repository.DeleteNotification(ctx, notification.ID)
			if err != nil {
				s.logger.Errorf("Failed to delete notification in repository", "id", notification.ID, "error", err)
				return err
			}

			err = s.InAppNotificationRepository.DeleteInAppNotification(ctx, notification.ID)
			if err != nil {
				s.logger.Errorf("Failed to delete in-app notification in repository", "id", notification.ID, "error", err)
				return err
			}
			return nil
		})
	case cps_const.RequestEnableNotification:
		err = s.Repository.RunInTransaction(ctx, func(ctx context.Context) error {

			notification, err = s.Repository.EnableDisableNotification(ctx, notification.ID, true)
			if err != nil {
				s.logger.Errorf("Failed to enable notification in repository", "id", notification.ID, "error", err)
				return err
			}

			err = s.InAppNotificationRepository.EnableDisableInAppNotification(ctx, notification.ID, true)
			if err != nil {
				s.logger.Errorf("Failed to enable in-app notification in repository", "id", notification.ID, "error", err)
				return err
			}
			return nil
		})
	case cps_const.RequestDisableNotification:
		err = s.Repository.RunInTransaction(ctx, func(ctx context.Context) error {
			notification, err = s.Repository.EnableDisableNotification(ctx, notification.ID, false)
			if err != nil {
				s.logger.Errorf("Failed to disable notification in repository", "id", notification.ID, "error", err)
				return err
			}

			err = s.InAppNotificationRepository.EnableDisableInAppNotification(ctx, notification.ID, false)
			if err != nil {
				s.logger.Errorf("Failed to disable in-app notification in repository", "id", notification.ID, "error", err)
				return err
			}
			return nil
		})
	default:
		s.logger.Errorf("Unsupported action requested", "action", requestedAction)
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	action.CurrentAction = notification
	s.logger.Infof("Authorization completed for action", "action", requestedAction, "id", notification.ID)
	return action, nil
}

// FetchNotificationByID fetches a notification by its ID
func (s *Service) FetchNotificationByID(ctx context.Context, id string) (*Notification, error) {
	s.logger.Infof("Fetching notification by ID", "id", id)
	return s.Repository.FetchNotificationByID(ctx, id)
}

// FetchNotifications fetches notifications with pagination and filtering
func (s *Service) FetchNotifications(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*Notification], error) {
	s.logger.Infof("Fetching notifications with filter", "filter", filterParam)
	return s.Repository.FetchNotifications(ctx, filterParam)
}

func generateNotification(notification Notification) *Notification {
	return &Notification{
		ID:                notification.ID,
		NotificationType:  notification.NotificationType,
		NotificationBody:  notification.NotificationBody,
		IsPublic:          notification.IsPublic,
		For:               notification.For,
		CreatedBy:         notification.CreatedBy,
		NotificationParts: notification.NotificationParts,
		Seen:              notification.Seen,
		Enabled:           notification.Enabled,
		IsDeleted:         notification.IsDeleted,
		CreatedAt:         notification.CreatedAt,
		LastModified:      notification.LastModified,
		DeletedAt:         notification.DeletedAt,
	}
}
