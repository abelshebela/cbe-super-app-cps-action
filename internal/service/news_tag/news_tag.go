package newstag_service

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "NewsTag", "Authorize")
	defer span.End()

	n.logger.Infof("[Authorize] authorizing news tag action: %s", cpsAction.RequestAction)

	actionData, err := local_util.JsonUnmarshal[model.NewsTagCPSAction](cpsAction.CurrentAction)
	if err != nil {
		n.logger.Errorf("[Authorize] failed to unmarshal CurrentAction: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, fmt.Errorf("failed to unmarshal CurrentAction to NewsTagCPSAction: %v", err)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateNewsTag):
		n.logger.Infof("[Authorize] creating news tag: %v", actionData.TagName)
		err := n.repo.Create(ctx, actionData.TagNameList)
		if err != nil {
			n.logger.Errorf("[Authorize] create news tag action failed: %v", err)
			span.AddEvent("Failed to create news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		n.logger.Infof("[Authorize] news tag created successfully")
	case string(constants.RequestUpdateNewsTag):
		n.logger.Infof("[Authorize] updating news tag id=%s to name=%s", actionData.ID.Hex(), actionData.TagName)
		err := n.repo.Update(ctx, actionData.ID.Hex(), actionData.TagName)
		if err != nil {
			n.logger.Errorf("[Authorize] update news tag action failed: %v", err)
			span.AddEvent("Failed to update news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		n.logger.Infof("[Authorize] news tag updated successfully for id: %s", actionData.ID.Hex())
		return cpsAction, nil

	case string(constants.RequestDeleteNewsTag):
		n.logger.Infof("[Authorize] deleting news tag id=%s", actionData.ID.Hex())
		err := n.repo.Delete(ctx, actionData.ID.Hex())
		if err != nil {
			n.logger.Errorf("[Authorize] delete news tag action failed: %v", err)
			span.AddEvent("Failed to delete news tag", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		n.logger.Infof("[Authorize] news tag deleted successfully for id: %s", actionData.ID.Hex())
		return cpsAction, nil

	default:
		n.logger.Errorf("[Authorize] invalid action: %s", cpsAction.RequestAction)
		span.AddEvent("Invalid action", trace.WithAttributes(
			attribute.String("error", localization.MsgInvalidAction),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	n.logger.Infof("[Authorize] news tag action authorized successfully: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

// CreateNewsTags implements service.NewsTagService.
func (n *newsTagService) CreateNewsTags(ctx context.Context, tagName []string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateNewsTags", "NewsTag", "CreateNewsTags")
	defer span.End()

	n.logger.Infof("[CreateNewsTags] creating news tags: %d tags", len(tagName))
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		n.logger.Errorf("[CreateNewsTags] incomplete user data")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}
	news_tag, err := n.repo.FindByNames(ctx, tagName)
	code, _ := local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[CreateNewsTags] error checking existing tags: %v", err)
		span.AddEvent("Failed to find news tags", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	if news_tag != nil {
		n.logger.Errorf("[CreateNewsTags] tag already exists")
		span.AddEvent("News tag already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorNewsTagWithNameAlreadyExists.Code),
		))
		return fmt.Errorf("%s", localization.ErrorNewsTagWithNameAlreadyExists.Code)
	}
	cpsAction := lib.CpsModelBuilder("", makerData, nil, model.NewsTagCPSAction{TagNameList: tagName}, string(constants.RequestCreateNewsTag), constants.CREATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("[CreateNewsTags] failed to create CPS action: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	n.logger.Infof("[CreateNewsTags] CPS action created successfully for %d tags", len(tagName))
	return nil
}

// DeleteNewsTag implements service.NewsTagService.
func (n *newsTagService) DeleteNewsTag(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteNewsTag", "NewsTag", "DeleteNewsTag")
	defer span.End()

	n.logger.Infof("[DeleteNewsTag] deleting news tag for id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		n.logger.Errorf("[DeleteNewsTag] incomplete user data")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	news_tag, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[DeleteNewsTag] resource not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorResourceNotFound.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorResourceNotFound.Code)
	}
	if err != nil {
		n.logger.Errorf("[DeleteNewsTag] failed to get news tag: %v", err)
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
		n.logger.Errorf("[DeleteNewsTag] failed to create CPS action: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	n.logger.Infof("[DeleteNewsTag] CPS action created successfully for id: %s", id)
	return nil
}

// FindAllWithPagination implements service.NewsTagService.
func (n *newsTagService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.NewsTag], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "NewsTag", "FindAllWithPagination")
	defer span.End()

	tag, err := n.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		n.logger.Errorf("[FindAllWithPagination] failed to fetch news tags: %v", err)
		span.AddEvent("Failed to fetch news tags", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	n.logger.Infof("[FindAllWithPagination] retrieved %d news tags", len(tag.Data))
	return tag, nil
}

// GetNewsTagByID implements service.NewsTagService.
func (n *newsTagService) GetNewsTagByID(ctx context.Context, id string) (*model.NewsTag, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetNewsTagByID", "NewsTag", "GetNewsTagByID")
	defer span.End()

	news_tag, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[GetNewsTagByID] resource not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return nil, fmt.Errorf("%s", code)
	} else if err != nil {
		n.logger.Errorf("[GetNewsTagByID] failed to get news tag: %v", err)
		span.AddEvent("Failed to get news tag", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	n.logger.Infof("[GetNewsTagByID] news tag retrieved successfully for id: %s", id)
	return news_tag, nil
}

// UpdateNewsTag implements service.NewsTagService.
func (n *newsTagService) UpdateNewsTag(ctx context.Context, id string, tagName string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateNewsTag", "NewsTag", "UpdateNewsTag")
	defer span.End()

	n.logger.Infof("[UpdateNewsTag] updating news tag id=%s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	news_tag, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[UpdateNewsTag] resource not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", code)
	} else if err != nil {
		n.logger.Errorf("[UpdateNewsTag] failed to get news tag: %v", err)
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
			n.logger.Errorf("[UpdateNewsTag] FindByNames repository error: %v", err)
			span.AddEvent("Failed to find news tags", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	if news_cat != nil {
		n.logger.Errorf("[UpdateNewsTag] news tag with name already exists")
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
		n.logger.Errorf("[UpdateNewsTag] failed to create CPS action: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	n.logger.Infof("[UpdateNewsTag] CPS action created successfully for id: %s", id)
	return nil
}
