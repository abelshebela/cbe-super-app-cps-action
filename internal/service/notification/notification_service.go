package notification

import (
	"context"
	"fmt"

	"cbe-super-app-cps-action/internal/constants"
	notify "cbe-super-app-cps-action/internal/constants/dto/notification"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helper "cbe-super-app-cps-action/internal/service/notification/core"
	"cbe-super-app-cps-action/internal/storage"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// notificationService implements NotificationService
type notificationService struct {
	repo       storage.NotificationRepository
	logger     shared_utils.Logger
	cpsService service.CPSActionService
}

func InitNotificationService(repo storage.NotificationRepository, logger shared_utils.Logger, cpsService service.CPSActionService) service.NotificationService {
	return &notificationService{
		repo:       repo,
		logger:     logger,
		cpsService: cpsService,
	}
}

func (s *notificationService) CreateNotification(ctx context.Context, req notify.NotificationRequest) (*notify.NotificationResponse, error) {
	s.logger.Infof("Creating notification request")

	// uniqueness check
	exist, err := s.repo.NotificationExists(ctx, req.NotificationType, constants.NotificationFor(req.For), nil)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("NOTIFICATION_ALREADY_EXISTS")
	}

	entity := helper.BuildCreateNotification(req)
	return notify.MapNotificationToResponse(&entity), nil
}

func (s *notificationService) UpdateNotification(ctx context.Context, id string, req notify.NotificationRequest) (*notify.NotificationResponse, error) {
	s.logger.Infof("Updating notification request: %s", id)
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cur := helper.BuildUpdateNotification(prev, req)
	return notify.MapNotificationToResponse(&cur), nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, id string) error {
	s.logger.Infof("Deleting notification request: %s", id)
	if id == "" {
		return fmt.Errorf("INVALID_ID")
	}
	return nil
}

func (s *notificationService) EnableDisableNotification(ctx context.Context, id string) error {
	s.logger.Infof("Enable/Disable notification request: %s", id)
	if id == "" {
		return fmt.Errorf("INVALID_ID")
	}
	return nil
}

func (s *notificationService) FetchNotificationByID(ctx context.Context, id string) (*notify.NotificationResponse, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return notify.MapNotificationToResponse(entity), nil
}

func (s *notificationService) FetchNotifications(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*notify.NotificationResponse], error) {
	if filterParam == nil {
		f := types.Filter{}
		filterParam = &f
	}
	entities, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		return nil, err
	}
	resp := make([]*notify.NotificationResponse, 0, len(entities.Data))
	for _, e := range entities.Data {
		resp = append(resp, notify.MapNotificationToResponse(e))
	}
	return &types.PaginatedResponse[[]*notify.NotificationResponse]{
		Data: resp,
		Meta: entities.Meta,
	}, nil
}

func (s *notificationService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	requested := action.RequestAction

	notif, err := func() (model.Notification, error) {
		return helper.BindNotificationFromAction(action.CurrentAction)
	}()
	if err != nil {
		return nil, err
	}

	switch requested {
	case string(constants.RequestCreatePublicNotification):
		if err := s.repo.Create(ctx, &notif); err != nil {
			return nil, err
		}
	case string(constants.RequestUpdatePublicNotification):
		if err := s.repo.Update(ctx, notif.ID.Hex(), &notif); err != nil {
			return nil, err
		}
	case string(constants.RequestDeleteNotification):
		if err := s.repo.Delete(ctx, notif.ID.Hex()); err != nil {
			return nil, err
		}
	case string(constants.RequestEnableNotification):
		if _, err := s.repo.EnableDisableNotification(ctx, notif.ID.Hex(), true); err != nil {
			return nil, err
		}
	case string(constants.RequestDisableNotification):
		if _, err := s.repo.EnableDisableNotification(ctx, notif.ID.Hex(), false); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("UNSUPPORTED_ACTION")
	}

	action.CurrentAction = notif
	return action, nil
}
