package newstag_service

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type newsTagService struct {
	repo       storage.NewsTagRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewNewsTagService(newsTagRepository storage.NewsTagRepository, cpsService service.CPSActionService, logger utils.Logger) service.NewsTagService {
	return &newsTagService{
		repo:       newsTagRepository,
		cpsService: cpsService,
		logger:     logger,
	}
}

// Authorize implements service.NewsTagService.
func (n *newsTagService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "NewsTag", "Authorize")
	defer span.End()

	log.Infof("[NewsTagSvc][Authorize] action: %s", cpsAction.RequestAction)

	actionData, err := local_util.JsonUnmarshal[model.NewsTagCPSAction](cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[NewsTagSvc][Authorize] unmarshal err: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, fmt.Errorf("failed to unmarshal CurrentAction to NewsTagCPSAction: %v", err)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateNewsTag):
		log.Infof("[NewsTagSvc][Authorize] creating: %v", actionData.TagName)
		err := n.repo.Create(ctx, actionData.TagNameList)
		if err != nil {
			log.Errorf("[NewsTagSvc][Authorize] create err: %v", err)
			span.AddEvent("Failed to create news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		log.Infof("[NewsTagSvc][Authorize] created")
	case string(constants.RequestUpdateNewsTag):
		log.Infof("[NewsTagSvc][Authorize] updating id: %s", actionData.ID.Hex())
		err := n.repo.Update(ctx, actionData.ID.Hex(), actionData.TagName)
		if err != nil {
			log.Errorf("[NewsTagSvc][Authorize] update err: %v", err)
			span.AddEvent("Failed to update news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		log.Infof("[NewsTagSvc][Authorize] updated id: %s", actionData.ID.Hex())
		return cpsAction, nil

	case string(constants.RequestDeleteNewsTag):
		log.Infof("[NewsTagSvc][Authorize] deleting id: %s", actionData.ID.Hex())
		err := n.repo.Delete(ctx, actionData.ID.Hex())
		if err != nil {
			log.Errorf("[NewsTagSvc][Authorize] delete err: %v", err)
			span.AddEvent("Failed to delete news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		log.Infof("[NewsTagSvc][Authorize] deleted id: %s", actionData.ID.Hex())
		return cpsAction, nil

	default:
		log.Errorf("[NewsTagSvc][Authorize] invalid: %s", cpsAction.RequestAction)
		span.AddEvent("Invalid action", trace.WithAttributes(
			attribute.String("error", localization.MsgInvalidAction),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	log.Infof("[NewsTagSvc][Authorize] done: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

// CreateNewsTags implements service.NewsTagService.
func (n *newsTagService) CreateNewsTags(ctx context.Context, tagName []string) error {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateNewsTags", "NewsTag", "CreateNewsTags")
	defer span.End()

	log.Infof("[NewsTagSvc][Create] count: %d", len(tagName))
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[NewsTagSvc][Create] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}
	news_tag, err := n.repo.FindByNames(ctx, tagName)
	code, _ := local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		log.Errorf("[NewsTagSvc][Create] find names err: %v", err)
		span.AddEvent("Failed to find news tags", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	if news_tag != nil {
		log.Errorf("[NewsTagSvc][Create] name exists")
		span.AddEvent("News tag already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorNewsTagWithNameAlreadyExists.Code),
		))
		return fmt.Errorf("%s", localization.ErrorNewsTagWithNameAlreadyExists.Code)
	}
	cpsAction := lib.CpsModelBuilder("", makerData, nil, model.NewsTagCPSAction{TagNameList: tagName}, string(constants.RequestCreateNewsTag), constants.CREATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Errorf("[NewsTagSvc][Create] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	log.Infof("[NewsTagSvc][Create] request created count: %d", len(tagName))
	return nil
}

// DeleteNewsTag implements service.NewsTagService.
func (n *newsTagService) DeleteNewsTag(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteNewsTag", "NewsTag", "DeleteNewsTag")
	defer span.End()

	log.Infof("[NewsTagSvc][Delete] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[NewsTagSvc][Delete] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	news_tag, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		log.Errorf("[NewsTagSvc][Delete] not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorResourceNotFound.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorResourceNotFound.Code)
	}
	if err != nil {
		log.Errorf("[NewsTagSvc][Delete] get err: %v", err)
		span.AddEvent("Failed to get news tag", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	updated_news_tag := *news_tag
	updated_news_tag.IsDeleted = true
	updated_news_tag.DeletedAt = time.Now()
	updated_news_tag.LastModifiedAt = time.Now()

	cpsAction := lib.CpsModelBuilder(id, makerData, news_tag, updated_news_tag, string(constants.RequestDeleteNewsTag), constants.DELETE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Errorf("[NewsTagSvc][Delete] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	log.Infof("[NewsTagSvc][Delete] request created id: %s", id)
	return nil
}

// FindAllWithPagination implements service.NewsTagService.
func (n *newsTagService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]model.NewsTag], error) {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "NewsTag", "FindAllWithPagination")
	defer span.End()

	tag, err := n.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		log.Errorf("[NewsTagSvc][FindAll] fetch err: %v", err)
		span.AddEvent("Failed to fetch news tags", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	log.Infof("[NewsTagSvc][FindAll] count: %d", len(tag.Data))
	return tag, nil
}

// GetNewsTagByID implements service.NewsTagService.
func (n *newsTagService) GetNewsTagByID(ctx context.Context, id string) (*model.NewsTag, error) {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetNewsTagByID", "NewsTag", "GetNewsTagByID")
	defer span.End()

	news_tag, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		log.Errorf("[NewsTagSvc][GetByID] not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return nil, fmt.Errorf("%s", code)
	} else if err != nil {
		log.Errorf("[NewsTagSvc][GetByID] get err: %v", err)
		span.AddEvent("Failed to get news tag", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	log.Infof("[NewsTagSvc][GetByID] found id: %s", id)
	return news_tag, nil
}

// UpdateNewsTag implements service.NewsTagService.
func (n *newsTagService) UpdateNewsTag(ctx context.Context, id string, tagName string) error {
	log := local_util.LoggerFromCtx(ctx, n.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateNewsTag", "NewsTag", "UpdateNewsTag")
	defer span.End()

	log.Infof("[NewsTagSvc][Update] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	news_tag, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		log.Errorf("[NewsTagSvc][Update] not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", code)
	} else if err != nil {
		log.Errorf("[NewsTagSvc][Update] get err: %v", err)
		span.AddEvent("Failed to get news tag", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	news_cat, err := n.repo.FindByNames(ctx, []string{tagName})
	code, _ = local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		if err != nil {
			log.Errorf("[NewsTagSvc][Update] find names err: %v", err)
			span.AddEvent("Failed to find news tags", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	if news_cat != nil {
		log.Errorf("[NewsTagSvc][Update] name exists")
		span.AddEvent("News tag already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorNewsCategoryWithNameAlreadyExists.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorNewsCategoryWithNameAlreadyExists.Code)
	}

	updated_news_tag := *news_tag
	updated_news_tag.TagName = tagName
	updated_news_tag.LastModifiedAt = time.Now()

	cpsAction := lib.CpsModelBuilder(id, makerData, news_tag, updated_news_tag, string(constants.RequestUpdateNewsTag), constants.UPDATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Errorf("[NewsTagSvc][Update] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	log.Infof("[NewsTagSvc][Update] request created id: %s", id)
	return nil
}
