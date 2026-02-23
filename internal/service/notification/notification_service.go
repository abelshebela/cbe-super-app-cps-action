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
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helper "cbe-super-app-cps-action/internal/service/notification/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// notificationService implements NotificationService
type notificationService struct {
	repo       storage.NotificationRepository
	logger     shared_utils.Logger
	cpsService service.CPSActionService
	store      *lib.NotificationStore
}

func InitNotificationService(repo storage.NotificationRepository, logger shared_utils.Logger, cpsService service.CPSActionService, store *lib.NotificationStore) service.NotificationService {
	return &notificationService{
		repo:       repo,
		logger:     logger,
		cpsService: cpsService,
		store:      store,
	}
}

func (s *notificationService) CreateNotification(ctx context.Context, req notify.NotificationRequest) (*notify.NotificationResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateNotification", "Notification", "CreateNotification")
	defer span.End()

	s.logger.Infof("[NotifSvc][Create] creating")
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		s.logger.Errorf("[NotifSvc][Create] incomplete user")
		span.AddEvent("Incomplete user context", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
		))
		return nil, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	entity := helper.BuildCreateNotification(req)

	cpsAction := lib.CpsModelBuilder("", maker, nil, entity, string(constants.RequestCreatePublicNotification), constants.CREATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[NotifSvc][Create] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	s.logger.Infof("[NotifSvc][Create] request created")
	return notify.MapNotificationToResponse(&entity), nil
}

func (s *notificationService) UpdateNotification(ctx context.Context, id string, req notify.NotificationRequest) (*notify.NotificationResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateNotification", "Notification", "UpdateNotification")
	defer span.End()

	s.logger.Infof("[NotifSvc][Update] id: %s", id)
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[NotifSvc][Update] find err: %v", err)
		span.AddEvent("Failed to find notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	cur := helper.BuildUpdateNotification(prev, req)

	makerdata := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, makerdata, prev, cur, string(constants.RequestUpdatePublicNotification), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[NotifSvc][Update] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}

	s.logger.Infof("[NotifSvc][Update] request created id: %s", id)
	return notify.MapNotificationToResponse(&cur), nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteNotification", "Notification", "DeleteNotification")
	defer span.End()

	s.logger.Infof("[NotifSvc][Delete] id: %s", id)
	if id == "" {
		s.logger.Errorf("[NotifSvc][Delete] invalid id")
		span.AddEvent("Invalid id provided", trace.WithAttributes(
			attribute.String("error", "INVALID_ID"),
		))
		return fmt.Errorf("INVALID_ID")
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[NotifSvc][Delete] find err: %v", err)
		span.AddEvent("Failed to find notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		s.logger.Errorf("[NotifSvc][Delete] incomplete user")
		span.AddEvent("Incomplete user context", trace.WithAttributes(
			attribute.String("error", "INCOMPLETE_USER_INFO"),
			attribute.String("id", id),
		))
		return fmt.Errorf("INCOMPLETE_USER_INFO")
	}
	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, nil, string(constants.RequestDeleteNotification), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[NotifSvc][Delete] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	s.logger.Infof("[NotifSvc][Delete] request created id: %s", id)
	return nil

}

func (s *notificationService) EnableNotification(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableNotification", "Notification", "EnableNotification")
	defer span.End()

	s.logger.Infof("[NotifSvc][Enable] id: %s", id)
	if id == "" {
		s.logger.Errorf("[NotifSvc][Enable] invalid id")
		span.AddEvent("Invalid id provided", trace.WithAttributes(
			attribute.String("error", "INVALID_ID"),
		))
		return fmt.Errorf("INVALID_ID")
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[NotifSvc][Enable] find err: %v", err)
		span.AddEvent("Failed to find notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if prev.Enabled {
		s.logger.Errorf("[NotifSvc][Enable] already enabled: %s", id)
		span.AddEvent("Notification already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorUserAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}
	updated := prev
	updated.Enabled = true

	maker := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, maker, prev, &updated, string(constants.RequestEnableNotification), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[NotifSvc][Enable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	s.logger.Infof("[NotifSvc][Enable] request created id: %s", id)
	return nil
}

func (s *notificationService) DisableNotification(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableNotification", "Notification", "DisableNotification")
	defer span.End()

	s.logger.Infof("[NotifSvc][Disable] id: %s", id)
	if id == "" {
		s.logger.Errorf("[NotifSvc][Disable] invalid id")
		span.AddEvent("Invalid id provided", trace.WithAttributes(
			attribute.String("error", "INVALID_ID"),
		))
		return fmt.Errorf("INVALID_ID")
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[NotifSvc][Disable] find err: %v", err)
		span.AddEvent("Failed to find notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if !prev.Enabled {
		s.logger.Errorf("[NotifSvc][Disable] already disabled: %s", id)
		span.AddEvent("Notification already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorUserAlreadyDisabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}
	updated := prev
	updated.Enabled = false

	maker := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, maker, prev, &updated, string(constants.RequestDisableNotification), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[NotifSvc][Disable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	s.logger.Infof("[NotifSvc][Disable] request created id: %s", id)
	return nil
}

func (s *notificationService) FetchNotificationByID(ctx context.Context, id string) (*notify.NotificationResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchNotificationByID", "Notification", "FetchNotificationByID")
	defer span.End()

	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[NotifSvc][FetchByID] fetch err: %v", err)
		span.AddEvent("Failed to fetch notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	s.logger.Infof("[NotifSvc][FetchByID] found id: %s", id)
	return notify.MapNotificationToResponse(entity), nil
}

func (s *notificationService) FetchNotifications(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]model.Notification], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchNotifications", "Notification", "FetchNotifications")
	defer span.End()

	entities, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		s.logger.Errorf("[NotifSvc][FetchAll] fetch err: %v", err)
		span.AddEvent("Failed to fetch notifications", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	s.logger.Infof("[NotifSvc][FetchAll] done")
	return entities, nil
}

func (s *notificationService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Notification", "Authorize")
	defer span.End()

	s.logger.Infof("[NotifSvc][Authorize] action: %s", action.RequestAction)
	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime

	switch action.RequestAction {
	case string(constants.RequestCreatePublicNotification):
		notif, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] bind err: %v", err)
			span.AddEvent("Failed to bind notification from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if notif.ID.IsZero() {
			notif.ID = bson.NewObjectID()
		}
		if notif.Title == "" || notif.NotificationType == "" || string(notif.For) == "" {
			s.logger.Errorf("[NotifSvc][Authorize] invalid data")
			span.AddEvent("Invalid notification data", trace.WithAttributes(
				attribute.String("error", localization.ErrorInvalidRequest.Code),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Create(ctx, &notif); err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] create err: %v", err)
			span.AddEvent("Failed to create notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		s.logger.Infof("[NotifSvc][Authorize] created")

		// Publish in-app broadcast to Kafka when a public notification is approved
		if s.store != nil && notif.IsPublic {
			// default 30-day expiry; adjust as needed or extend DTO to accept expiry
			exp := time.Now().Add(30 * 24 * time.Hour)
			payload := types.InAppBroadcastMessage{
				Title:     notif.Title,
				Message:   notif.NotificationBody,
				Type:      "inapp",
				ExpiresAt: exp,
			}
			if err := s.store.PublishInAppBroadcast(ctx, payload); err != nil {
				s.logger.Errorf("[NotifSvc][Authorize] broadcast err: %v", err)
			} else {
				s.logger.Infof("[NotifSvc][Authorize] broadcast published")
			}
		}
		return action, nil

	case string(constants.RequestUpdatePublicNotification):
		notif, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] bind err: %v", err)
			span.AddEvent("Failed to bind notification from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		id := action.UniqueId
		if !notif.ID.IsZero() {
			id = notif.ID.Hex()
		}
		if err := s.repo.Update(ctx, id, &notif); err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] update err: %v", err)
			span.AddEvent("Failed to update notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		s.logger.Infof("[NotifSvc][Authorize] updated id: %s", id)
		return action, nil

	case string(constants.RequestDeleteNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] bind err: %v", err)
			span.AddEvent("Failed to bind notification from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Delete(ctx, action.UniqueId); err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] delete err: %v", err)
			span.AddEvent("Failed to delete notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		s.logger.Infof("[NotifSvc][Authorize] deleted id: %s", action.UniqueId)
		return action, nil

	case string(constants.RequestEnableNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] bind err: %v", err)
			span.AddEvent("Failed to bind notification from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if _, err := s.repo.EnableDisableNotification(ctx, action.UniqueId, true); err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] enable err: %v", err)
			span.AddEvent("Failed to enable notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		s.logger.Infof("[NotifSvc][Authorize] enabled id: %s", action.UniqueId)
		return action, nil

	case string(constants.RequestDisableNotification):
		_, err := helper.BindNotificationFromAction(action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] bind err: %v", err)
			span.AddEvent("Failed to bind notification from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if _, err := s.repo.EnableDisableNotification(ctx, action.UniqueId, false); err != nil {
			s.logger.Errorf("[NotifSvc][Authorize] disable err: %v", err)
			span.AddEvent("Failed to disable notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		s.logger.Infof("[NotifSvc][Authorize] disabled id: %s", action.UniqueId)
		return action, nil
	default:
		s.logger.Errorf("[NotifSvc][Authorize] unsupported: %s", action.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(action.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}
