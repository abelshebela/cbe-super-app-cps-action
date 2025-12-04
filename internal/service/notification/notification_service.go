package notification

import (
	"context"
	"errors"
	"fmt"
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
	s.logger.Infof("Creating notification request: %+v", req)
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		s.logger.Errorf("incomplete user context for notification action | context = %v", maker)
		return nil, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	entity := helper.BuildCreateNotification(req)

	cpsAction := lib.CpsModelBuilder("", maker, nil, entity, string(constants.RequestCreatePublicNotification), constants.CREATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("failed to create CPS action for notification | action=%s | err=%v", constants.RequestCreatePublicNotification, err)
		return nil, err
	}

	return notify.MapNotificationToResponse(&entity), nil
}

func (s *notificationService) UpdateNotification(ctx context.Context, id string, req notify.NotificationRequest) (*notify.NotificationResponse, error) {
	if s.logger != nil {
		s.logger.Infof("Updating notification request: %s", id)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cur := helper.BuildUpdateNotification(prev, req)

	makerdata := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, makerdata, prev, cur, string(constants.RequestUpdatePublicNotification), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for notification | action=%s | err=%v", constants.RequestUpdatePublicNotification, err)
		}
		return nil, err
	}

	return notify.MapNotificationToResponse(&cur), nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, id string) error {
	if s.logger != nil {
		s.logger.Infof("Deleting notification request: %s", id)
	}
	if id == "" {
		return fmt.Errorf("INVALID_ID")
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for notification action | context = %v", maker)
		}
		return fmt.Errorf("INCOMPLETE_USER_INFO")
	}
	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, nil, string(constants.RequestDeleteNotification), string(constants.DELETE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)

}

func (s *notificationService) EnableNotification(ctx context.Context, id string) error {
	if s.logger != nil {
		s.logger.Infof("Enabling notification request: %s", id)
	}
	if id == "" {
		return fmt.Errorf("INVALID_ID")
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		return err
	}
	if prev.Enabled {
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}
	updated := prev
	updated.Enabled = true

	maker := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, maker, prev, &updated, string(constants.RequestEnableNotification), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *notificationService) DisableNotification(ctx context.Context, id string) error {
	if s.logger != nil {
		s.logger.Infof("Disabling notification request: %s", id)
	}
	if id == "" {
		return fmt.Errorf("INVALID_ID")
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		return err
	}
	if !prev.Enabled {
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}
	updated := prev
	updated.Enabled = false

	maker := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, maker, prev, &updated, string(constants.RequestDisableNotification), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *notificationService) FetchNotificationByID(ctx context.Context, id string) (*notify.NotificationResponse, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return notify.MapNotificationToResponse(entity), nil
}

func (s *notificationService) FetchNotifications(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*notify.NotificationResponse], error) {
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
	fmt.Println("//// here to authorize..... ", action)
	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime

	switch action.RequestAction {
	case string(constants.RequestCreatePublicNotification):
		notif, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if notif.ID.IsZero() {
			notif.ID = bson.NewObjectID()
		}
		if notif.Title == "" || notif.NotificationType == "" || string(notif.For) == "" {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Create(ctx, &notif); err != nil {
			fmt.Println("createiting there issue on this>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>><<<<<<<<<<<<<<<<<<<<<<")
			return nil, err
		}
		return action, nil

	case string(constants.RequestUpdatePublicNotification):
		notif, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		id := action.UniqueId
		if !notif.ID.IsZero() {
			id = notif.ID.Hex()
		}
		// // Uniqueness check to avoid duplicate type+for on update
		// exists, err := s.repo.NotificationExists(ctx, notif.NotificationType, notif.For, &id)
		// if err != nil {
		// 	return nil, err
		// }
		// if exists {
		// 	return nil, errors.New(localization.ErrorNotificationAlreadyExists.Code)
		// }
		if err := s.repo.Update(ctx, id, &notif); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.RequestDeleteNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Delete(ctx, action.UniqueId); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.RequestEnableNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if _, err := s.repo.EnableDisableNotification(ctx, action.UniqueId, true); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.RequestDisableNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if _, err := s.repo.EnableDisableNotification(ctx, action.UniqueId, false); err != nil {
			return nil, err
		}
		return action, nil
	}

	return nil, errors.New(localization.ErrorInvalidRequest.Code)
}
