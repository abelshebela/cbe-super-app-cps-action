package notification

import (
	"context"
	"encoding/json"
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

	local_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

func decodeNotificationDocumentAction(current interface{}) (local_model.NotificationDocument, error) {
	out := local_model.NotificationDocument{}
	if current == nil {
		return out, nil
	}

	b, err := json.Marshal(current)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return out, err
	}

	return out, nil
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
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateNotification", "Notification", "CreateNotification")
	defer span.End()

	log.Infof("[NotifSvc][Create] creating")
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[NotifSvc][Create] incomplete user")
		span.AddEvent("Incomplete user context", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
		))
		return nil, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	entity := helper.BuildCreateNotification(req)

	cpsAction := lib.CpsModelBuilder("", maker, nil, entity, string(constants.RequestCreatePublicNotification), constants.CREATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Errorf("[NotifSvc][Create] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	log.Infof("[NotifSvc][Create] request created")
	return notify.MapNotificationToResponse(&entity), nil
}

func (s *notificationService) UpdateNotification(ctx context.Context, id string, req notify.NotificationRequest) (*notify.NotificationResponse, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateNotification", "Notification", "UpdateNotification")
	defer span.End()

	log.Infof("[NotifSvc][Update] id: %s", id)
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[NotifSvc][Update] find err: %v", err)
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
		log.Errorf("[NotifSvc][Update] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}

	log.Infof("[NotifSvc][Update] request created id: %s", id)
	return notify.MapNotificationToResponse(&cur), nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteNotification", "Notification", "DeleteNotification")
	defer span.End()

	log.Infof("[NotifSvc][Delete] id: %s", id)
	if id == "" {
		log.Errorf("[NotifSvc][Delete] invalid id")
		span.AddEvent("Invalid id provided", trace.WithAttributes(
			attribute.String("error", "INVALID_ID"),
		))
		return fmt.Errorf("INVALID_ID")
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[NotifSvc][Delete] find err: %v", err)
		span.AddEvent("Failed to find notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[NotifSvc][Delete] incomplete user")
		span.AddEvent("Incomplete user context", trace.WithAttributes(
			attribute.String("error", "INCOMPLETE_USER_INFO"),
			attribute.String("id", id),
		))
		return fmt.Errorf("INCOMPLETE_USER_INFO")
	}
	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, nil, string(constants.RequestDeleteNotification), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		log.Errorf("[NotifSvc][Delete] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	log.Infof("[NotifSvc][Delete] request created id: %s", id)
	return nil

}

func (s *notificationService) EnableNotification(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "EnableNotification", "Notification", "EnableNotification")
	defer span.End()

	log.Infof("[NotifSvc][Enable] id: %s", id)
	if id == "" {
		log.Errorf("[NotifSvc][Enable] invalid id")
		span.AddEvent("Invalid id provided", trace.WithAttributes(
			attribute.String("error", "INVALID_ID"),
		))
		return fmt.Errorf("INVALID_ID")
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[NotifSvc][Enable] find err: %v", err)
		span.AddEvent("Failed to find notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	enabled := true
	actionPayload := local_model.NotificationDocument{Enabled: &enabled}

	maker := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, maker, prev, actionPayload, string(constants.RequestEnableNotification), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		log.Errorf("[NotifSvc][Enable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	log.Infof("[NotifSvc][Enable] request created id: %s", id)
	return nil
}

func (s *notificationService) DisableNotification(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DisableNotification", "Notification", "DisableNotification")
	defer span.End()

	log.Infof("[NotifSvc][Disable] id: %s", id)
	if id == "" {
		log.Errorf("[NotifSvc][Disable] invalid id")
		span.AddEvent("Invalid id provided", trace.WithAttributes(
			attribute.String("error", "INVALID_ID"),
		))
		return fmt.Errorf("INVALID_ID")
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[NotifSvc][Disable] find err: %v", err)
		span.AddEvent("Failed to find notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	enabled := false
	actionPayload := local_model.NotificationDocument{Enabled: &enabled}

	maker := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, maker, prev, actionPayload, string(constants.RequestDisableNotification), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		log.Errorf("[NotifSvc][Disable] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	log.Infof("[NotifSvc][Disable] request created id: %s", id)
	return nil
}

func (s *notificationService) FetchNotificationByID(ctx context.Context, id string) (*notify.NotificationResponse, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "FetchNotificationByID", "Notification", "FetchNotificationByID")
	defer span.End()

	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[NotifSvc][FetchByID] fetch err: %v", err)
		span.AddEvent("Failed to fetch notification", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	log.Infof("[NotifSvc][FetchByID] found id: %s", id)
	return notify.MapNotificationToResponse(entity), nil
}

func (s *notificationService) FetchNotifications(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]local_model.NotificationDocument], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "FetchNotifications", "Notification", "FetchNotifications")
	defer span.End()

	entities, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		log.Errorf("[NotifSvc][FetchAll] fetch err: %v", err)
		span.AddEvent("Failed to fetch notifications", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	log.Infof("[NotifSvc][FetchAll] done")
	return entities, nil
}

func (s *notificationService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Notification", "Authorize")
	defer span.End()

	log.Infof("[NotifSvc][Authorize] action: %s", action.RequestAction)
	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime
	notif, err := decodeNotificationDocumentAction(action.CurrentAction)
	if err != nil {
		log.Errorf("[NotifSvc][Authorize] bind payload err: %v", err)
		span.AddEvent("Failed to bind notification payload", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, localization.ErrorUnexpectedError
	}

	switch action.RequestAction {
	case string(constants.RequestCreatePublicNotification):
		if notif.Title == "" || notif.Category == "" || notif.For == "" {
			log.Errorf("[NotifSvc][Authorize] invalid data")
			span.AddEvent("Invalid notification data", trace.WithAttributes(
				attribute.String("error", localization.ErrorInvalidRequest.Code),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Create(ctx, &notif); err != nil {
			log.Errorf("[NotifSvc][Authorize] create err: %v", err)
			span.AddEvent("Failed to create notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		log.Infof("[NotifSvc][Authorize] created")

		return action, nil

	case string(constants.RequestUpdatePublicNotification):
		id := action.UniqueId

		if err := s.repo.Update(ctx, id, &notif); err != nil {
			log.Errorf("[NotifSvc][Authorize] update err: %v", err)
			span.AddEvent("Failed to update notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		log.Infof("[NotifSvc][Authorize] updated id: %s", id)
		return action, nil

	case string(constants.RequestDeleteNotification):
		if err := s.repo.Delete(ctx, action.UniqueId); err != nil {
			log.Errorf("[NotifSvc][Authorize] delete err: %v", err)
			span.AddEvent("Failed to delete notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		log.Infof("[NotifSvc][Authorize] deleted id: %s", action.UniqueId)
		return action, nil

	case string(constants.RequestEnableNotification):
		if notif.Enabled == nil || !*notif.Enabled {
			log.Errorf("[NotifSvc][Authorize] invalid enable payload")
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if _, err := s.repo.EnableDisableNotification(ctx, action.UniqueId, true); err != nil {
			log.Errorf("[NotifSvc][Authorize] enable err: %v", err)
			span.AddEvent("Failed to enable notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		log.Infof("[NotifSvc][Authorize] enabled id: %s", action.UniqueId)
		return action, nil

	case string(constants.RequestDisableNotification):
		if notif.Enabled == nil || *notif.Enabled {
			log.Errorf("[NotifSvc][Authorize] invalid disable payload")
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if _, err := s.repo.EnableDisableNotification(ctx, action.UniqueId, false); err != nil {
			log.Errorf("[NotifSvc][Authorize] disable err: %v", err)
			span.AddEvent("Failed to disable notification", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		log.Infof("[NotifSvc][Authorize] disabled id: %s", action.UniqueId)
		return action, nil
	default:
		log.Errorf("[NotifSvc][Authorize] unsupported: %s", action.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(action.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}
