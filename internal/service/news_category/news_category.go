package newscategory_service

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type NewsCategoryService struct {
	repo       storage.NewsCategoryRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewNewsCategoryService(repo storage.NewsCategoryRepository, cpsService service.CPSActionService, logger utils.Logger) service.NewsCategoryService {
	return &NewsCategoryService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

// Authorize implements service.NewsCategoryService.
func (n *NewsCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "NewsCategory", "Authorize")
	defer span.End()

	n.logger.Infof("[NewsCatSvc][Authorize] action: %s", cpsAction.RequestAction)

	actionData, err := local_util.JsonUnmarshal[model.NewsCategoryCPSAction](cpsAction.CurrentAction)
	if err != nil {
		n.logger.Errorf("[NewsCatSvc][Authorize] unmarshal err: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, fmt.Errorf("failed to unmarshal CurrentAction to NewsCategoryCPSAction: %v", err)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateNewsCategory):
		n.logger.Infof("[NewsCatSvc][Authorize] creating")
		err := n.repo.Create(ctx, actionData.CategoryNameList)
		if err != nil {
			n.logger.Errorf("[NewsCatSvc][Authorize] create err: %v", err)
			span.AddEvent("Failed to create news category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		n.logger.Infof("[NewsCatSvc][Authorize] created")
	case string(constants.RequestUpdateNewsCategory):
		n.logger.Infof("[NewsCatSvc][Authorize] updating id: %s", actionData.ID.Hex())
		err := n.repo.Update(ctx, actionData.ID.Hex(), actionData.CategoryName)
		if err != nil {
			n.logger.Errorf("[NewsCatSvc][Authorize] update err: %v", err)
			span.AddEvent("Failed to update news category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		n.logger.Infof("[NewsCatSvc][Authorize] updated id: %s", actionData.ID.Hex())
		return cpsAction, nil

	case string(constants.RequestDeleteNewsCategory):
		n.logger.Infof("[NewsCatSvc][Authorize] deleting id: %s", actionData.ID.Hex())
		err := n.repo.Delete(ctx, actionData.ID.Hex())
		if err != nil {
			n.logger.Errorf("[NewsCatSvc][Authorize] delete err: %v", err)
			span.AddEvent("Failed to delete news category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}
		n.logger.Infof("[NewsCatSvc][Authorize] deleted id: %s", actionData.ID.Hex())
		return cpsAction, nil

	default:
		n.logger.Errorf("[NewsCatSvc][Authorize] invalid: %s", cpsAction.RequestAction)
		span.AddEvent("Invalid action", trace.WithAttributes(
			attribute.String("error", localization.MsgInvalidAction),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	n.logger.Infof("[NewsCatSvc][Authorize] done: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

// CreateNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) CreateNewsCategory(ctx context.Context, categoryName []string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateNewsCategory", "NewsCategory", "CreateNewsCategory")
	defer span.End()

	n.logger.Infof("[NewsCatSvc][Create] count: %d", len(categoryName))
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		n.logger.Errorf("[NewsCatSvc][Create] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}
	news_cat, err := n.repo.FindByNames(ctx, categoryName)
	code, _ := local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		if err != nil {
			n.logger.Errorf("[NewsCatSvc][Create] find names err: %v", err)
			span.AddEvent("Failed to find news categories", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			return err
		}
	}

	if news_cat != nil {
		n.logger.Errorf("[NewsCatSvc][Create] name exists")
		span.AddEvent("News category already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorNewsCategoryWithNameAlreadyExists.Code),
		))
		return fmt.Errorf("%s", localization.ErrorNewsCategoryWithNameAlreadyExists.Code)
	}
	cpsAction := lib.CpsModelBuilder("", makerData, nil, model.NewsCategoryCPSAction{CategoryNameList: categoryName}, string(constants.RequestCreateNewsCategory), constants.CREATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("[NewsCatSvc][Create] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	n.logger.Infof("[NewsCatSvc][Create] request created count: %d", len(categoryName))

	return nil
}

// DeleteNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) DeleteNewsCategory(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteNewsCategory", "NewsCategory", "DeleteNewsCategory")
	defer span.End()

	n.logger.Infof("[NewsCatSvc][Delete] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		n.logger.Errorf("[NewsCatSvc][Delete] incomplete user")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[NewsCatSvc][Delete] not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorResourceNotFound.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorResourceNotFound.Code)
	} else if err != nil {
		n.logger.Errorf("[NewsCatSvc][Delete] get err: %v", err)
		span.AddEvent("Failed to get news category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	updated_news_category := *news_category
	updated_news_category.IsDeleted = true
	updated_news_category.DeletedAt = time.Now()
	updated_news_category.LastModifiedAt = time.Now()

	cpsAction := lib.CpsModelBuilder(id, makerData, news_category, updated_news_category, string(constants.RequestDeleteNewsCategory), constants.DELETE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("[NewsCatSvc][Delete] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	n.logger.Infof("[NewsCatSvc][Delete] request created id: %s", id)

	return nil
}

// FindAllWithPagination implements service.NewsCategoryService.
func (n *NewsCategoryService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]model.NewsCategory], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "NewsCategory", "FindAllWithPagination")
	defer span.End()

	categories, err := n.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		n.logger.Errorf("[NewsCatSvc][FindAll] fetch err: %v", err)
		span.AddEvent("Failed to fetch news categories", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	n.logger.Infof("[NewsCatSvc][FindAll] count: %d", len(categories.Data))
	return categories, nil
}

// GetNewsCategoryByID implements service.NewsCategoryService.
func (n *NewsCategoryService) GetNewsCategoryByID(ctx context.Context, id string) (*model.NewsCategory, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetNewsCategoryByID", "NewsCategory", "GetNewsCategoryByID")
	defer span.End()

	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[NewsCatSvc][GetByID] not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return nil, fmt.Errorf("%s", code)
	} else if err != nil {
		n.logger.Errorf("[NewsCatSvc][GetByID] get err: %v", err)
		span.AddEvent("Failed to get news category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	n.logger.Infof("[NewsCatSvc][GetByID] found id: %s", id)
	return news_category, nil
}

// UpdateNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) UpdateNewsCategory(ctx context.Context, id string, categoryName string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateNewsCategory", "NewsCategory", "UpdateNewsCategory")
	defer span.End()

	n.logger.Infof("[NewsCatSvc][Update] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)

	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[NewsCatSvc][Update] not found: %s", id)
		span.AddEvent("Resource not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", code)
	} else if err != nil {
		n.logger.Errorf("[NewsCatSvc][Update] get err: %v", err)
		span.AddEvent("Failed to get news category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	news_cat, err := n.repo.FindByNames(ctx, []string{categoryName})
	code, _ = local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		if err != nil {
			n.logger.Errorf("[NewsCatSvc][Update] find names err: %v", err)
			span.AddEvent("Failed to find news categories", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
	}

	if news_cat != nil {
		n.logger.Errorf("[NewsCatSvc][Update] name exists")
		span.AddEvent("News category already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorNewsCategoryWithNameAlreadyExists.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorNewsCategoryWithNameAlreadyExists.Code)
	}

	updated_news_category := *news_category
	updated_news_category.CategoryName = categoryName
	updated_news_category.LastModifiedAt = time.Now()

	cpsAction := lib.CpsModelBuilder(id, makerData, news_category, updated_news_category, string(constants.RequestUpdateNewsCategory), constants.UPDATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("[NewsCatSvc][Update] cps action err: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	n.logger.Infof("[NewsCatSvc][Update] request created id: %s", id)

	return nil
}
