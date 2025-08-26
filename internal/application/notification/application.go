package notification_application

import (
	"context"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	notification_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// NotificationApplicationAbstracts defines the interface for notification application logic
type NotificationApplicationAbstracts interface {
	CreateNotification(ctx context.Context, notification dto.NotificationRequest, maker cps_entities.User) error
	UpdateNotification(ctx context.Context, id string, notification dto.NotificationRequest, maker cps_entities.User) error
	DeleteNotification(ctx context.Context, id string, maker cps_entities.User) error
	EnableDisableNotification(ctx context.Context, id string, maker cps_entities.User, enable bool) error

	FetchNotificationByID(ctx context.Context, id string) (*notification_entity.Notification, error)
	FetchNotifications(ctx context.Context, filterParam *constant.MongoFilter) (*common_util.PaginatedResponse[[]*notification_entity.Notification], error)
}

// NotificationApplication implements NotificationApplicationAbstracts
type NotificationApplication struct {
	service    domain.NotificationService
	cpsService cps_service.CPSActionService
	logger     utils.Logger
}

// NewNotificationApplication creates a new notification application instance
func NewNotificationApplication(
	service domain.NotificationService,
	cpsService cps_service.CPSActionService,
	logger utils.Logger,
) NotificationApplicationAbstracts {
	return &NotificationApplication{
		service:    service,
		cpsService: cpsService,
		logger:     logger,
	}
}

// handleCPSAction encapsulates the common CPS action logic
func (a *NotificationApplication) handleCPSAction(
	ctx context.Context,
	maker cps_entities.User,
	requestAction cps_const.RequestAction,
	curData, prevData interface{},
	actionType cps_const.ActionType,
) error {
	cpsAction := a.cpsService.BuildCPSAction(ctx, cps_entities.CreateCPSRequest{
		User:          maker,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    actionType,
	})

	_, err := a.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("Failed to create CPS action for %s: %v", requestAction, err)
		return err
	}
	return nil
}

// CreateNotification creates a new notification
func (a *NotificationApplication) CreateNotification(ctx context.Context, notification dto.NotificationRequest, maker cps_entities.User) error {

	res, err := a.service.CreateNotification(ctx, notification)
	if err != nil {
		a.logger.Errorf("Failed to create notification: %v", err)
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestCreateNotification, res, nil, cps_const.ActionCreate)
}

// UpdateNotification updates an existing notification
func (a *NotificationApplication) UpdateNotification(ctx context.Context, id string, notification dto.NotificationRequest, maker cps_entities.User) error {
	curAction, prevAction, err := a.service.UpdateNotification(ctx, id, notification)
	if err != nil {
		a.logger.Errorf("Failed to update notification ID %s: %v", id, err)
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestUpdateNotification, curAction, prevAction, cps_const.ActionUpdate)
}

// DeleteNotification deletes a notification
func (a *NotificationApplication) DeleteNotification(ctx context.Context, id string, maker cps_entities.User) error {
	curAction, prevAction, err := a.service.DeleteNotification(ctx, id)
	if err != nil {
		a.logger.Errorf("Failed to delete notification ID %s: %v", id, err)
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestDeleteNotification, curAction, prevAction, cps_const.ActionDelete)
}

// EnableDisableNotification enables or disables a notification
func (a *NotificationApplication) EnableDisableNotification(ctx context.Context, id string, maker cps_entities.User, enable bool) error {
	curAction, prevAction, err := a.service.EnableDisableNotification(ctx, id, enable)
	if err != nil {
		a.logger.Errorf("Failed to enable/disable notification ID %s: %v", id, err)
		return err
	}

	var action cps_const.RequestAction
	if enable {
		action = cps_const.RequestEnableNotification
	} else {
		action = cps_const.RequestDisableNotification
	}

	return a.handleCPSAction(ctx, maker, action, curAction, prevAction, cps_const.ActionUpdate)
}

// FetchNotificationByID fetches a notification by its ID
func (a *NotificationApplication) FetchNotificationByID(ctx context.Context, id string) (*notification_entity.Notification, error) {
	return a.service.FetchNotificationByID(ctx, id)
}

// FetchNotifications fetches notifications with pagination and filtering
func (a *NotificationApplication) FetchNotifications(ctx context.Context, filterParam *constant.MongoFilter) (*common_util.PaginatedResponse[[]*notification_entity.Notification], error) {
	return a.service.FetchNotifications(ctx, filterParam)
}
