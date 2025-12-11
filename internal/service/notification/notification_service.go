package notification

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	notify "cbe-super-app-cps-action/internal/constants/dto/notification"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helper "cbe-super-app-cps-action/internal/service/notification/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	s.logger.Infof("[CreateNotification] creating notification request")
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		s.logger.Errorf("[CreateNotification] incomplete user context")
		return nil, errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	entity := helper.BuildCreateNotification(req)

	cpsAction := lib.CpsModelBuilder("", maker, nil, entity, string(constants.RequestCreatePublicNotification), constants.CREATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[CreateNotification] failed to create CPS action: %v", err)
		return nil, err
	}

	s.logger.Infof("[CreateNotification] notification creation request created successfully")
	return notify.MapNotificationToResponse(&entity), nil
}

func (s *notificationService) UpdateNotification(ctx context.Context, id string, req notify.NotificationRequest) (*notify.NotificationResponse, error) {
	s.logger.Infof("[UpdateNotification] updating notification for id: %s", id)
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[UpdateNotification] failed to find notification: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		return nil, err
	}
	cur := helper.BuildUpdateNotification(prev, req)

	makerdata := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, makerdata, prev, cur, string(constants.RequestUpdatePublicNotification), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[UpdateNotification] failed to create CPS action: %v", err)
		return nil, err
	}

	s.logger.Infof("[UpdateNotification] notification update request created successfully for id: %s", id)
	return notify.MapNotificationToResponse(&cur), nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, id string) error {
	s.logger.Infof("[DeleteNotification] deleting notification for id: %s", id)
	if id == "" {
		s.logger.Errorf("[DeleteNotification] invalid id provided")
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[DeleteNotification] failed to find notification: %v", err)
		return err
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		s.logger.Errorf("[DeleteNotification] incomplete user context")
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, nil, string(constants.RequestDeleteNotification), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[DeleteNotification] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[DeleteNotification] notification deletion request created successfully for id: %s", id)
	return nil

}

func (s *notificationService) EnableNotification(ctx context.Context, id string) error {
	s.logger.Infof("[EnableNotification] enabling notification for id: %s", id)
	if id == "" {
		s.logger.Errorf("[EnableNotification] invalid id provided")
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			s.logger.Errorf("[EnableNotification] notification not found: %s", id)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		s.logger.Errorf("[EnableNotification] failed to find notification: %v", err)
		return err
	}
	if prev.Enabled {
		s.logger.Errorf("[EnableNotification] notification already enabled: %s", id)
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	updated := prev
	updated.Enabled = true

	maker := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, maker, prev, &updated, string(constants.RequestEnableNotification), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[EnableNotification] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableNotification] notification enable request created successfully for id: %s", id)
	return nil
}

func (s *notificationService) DisableNotification(ctx context.Context, id string) error {
	s.logger.Infof("[DisableNotification] disabling notification for id: %s", id)
	if id == "" {
		s.logger.Errorf("[DisableNotification] invalid id provided")
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			s.logger.Errorf("[DisableNotification] notification not found: %s", id)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		s.logger.Errorf("[DisableNotification] failed to find notification: %v", err)
		return err
	}
	if !prev.Enabled {
		s.logger.Errorf("[DisableNotification] notification already disabled: %s", id)
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}
	updated := prev
	updated.Enabled = false

	maker := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, maker, prev, &updated, string(constants.RequestDisableNotification), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[DisableNotification] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[DisableNotification] notification disable request created successfully for id: %s", id)
	return nil
}

func (s *notificationService) FetchNotificationByID(ctx context.Context, id string) (*notify.NotificationResponse, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[FetchNotificationByID] failed to fetch notification: %v", err)
		return nil, err
	}
	s.logger.Infof("[FetchNotificationByID] notification retrieved successfully for id: %s", id)
	return notify.MapNotificationToResponse(entity), nil
}

func (s *notificationService) FetchNotifications(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*notify.NotificationResponse], error) {
	entities, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		s.logger.Errorf("[FetchNotifications] failed to fetch notifications: %v", err)
		return nil, err
	}
	resp := make([]*notify.NotificationResponse, 0, len(entities.Data))
	for _, e := range entities.Data {
		resp = append(resp, notify.MapNotificationToResponse(e))
	}
	s.logger.Infof("[FetchNotifications] retrieved %d notifications", len(resp))
	return &types.PaginatedResponse[[]*notify.NotificationResponse]{
		Data: resp,
		Meta: entities.Meta,
	}, nil
}

func (s *notificationService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[Authorize] authorizing notification action: %s", action.RequestAction)
	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime

	switch action.RequestAction {
	case string(constants.RequestCreatePublicNotification):
		notif, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to bind notification from action: %v", err)
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if notif.ID.IsZero() {
			notif.ID = bson.NewObjectID()
		}
		if notif.Title == "" || notif.NotificationType == "" || string(notif.For) == "" {
			s.logger.Errorf("[Authorize] invalid notification data")
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Create(ctx, &notif); err != nil {
			s.logger.Errorf("[Authorize] failed to create notification: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] notification created successfully")
		return action, nil

	case string(constants.RequestUpdatePublicNotification):
		notif, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to bind notification from action: %v", err)
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		id := action.UniqueId
		if !notif.ID.IsZero() {
			id = notif.ID.Hex()
		}
		if err := s.repo.Update(ctx, id, &notif); err != nil {
			s.logger.Errorf("[Authorize] failed to update notification: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] notification updated successfully for id: %s", id)
		return action, nil

	case string(constants.RequestDeleteNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to bind notification from action: %v", err)
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Delete(ctx, action.UniqueId); err != nil {
			s.logger.Errorf("[Authorize] failed to delete notification: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] notification deleted successfully for id: %s", action.UniqueId)
		return action, nil

	case string(constants.RequestEnableNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to bind notification from action: %v", err)
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if _, err := s.repo.EnableDisableNotification(ctx, action.UniqueId, true); err != nil {
			s.logger.Errorf("[Authorize] failed to enable notification: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] notification enabled successfully for id: %s", action.UniqueId)
		return action, nil

	case string(constants.RequestDisableNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to bind notification from action: %v", err)
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if _, err := s.repo.EnableDisableNotification(ctx, action.UniqueId, false); err != nil {
			s.logger.Errorf("[Authorize] failed to disable notification: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] notification disabled successfully for id: %s", action.UniqueId)
		return action, nil
	default:
		s.logger.Errorf("[Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}
