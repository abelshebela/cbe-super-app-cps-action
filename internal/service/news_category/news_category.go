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
	n.logger.Infof("[Authorize] authorizing news category action: %s", cpsAction.RequestAction)

	actionData, err := local_util.JsonUnmarshal[model.NewsCategoryCPSAction](cpsAction.CurrentAction)
	if err != nil {
		n.logger.Errorf("[Authorize] failed to unmarshal CurrentAction: %v", err)
		return nil, fmt.Errorf("failed to unmarshal CurrentAction to NewsCategoryCPSAction: %v", err)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateNewsCategory):
		n.logger.Infof("[Authorize] creating news category")
		err := n.repo.Create(ctx, actionData.CategoryNameList)
		if err != nil {
			n.logger.Errorf("[Authorize] create news category action failed: %v", err)
			return nil, err
		}
		n.logger.Infof("[Authorize] news category created successfully")
	case string(constants.RequestUpdateNewsCategory):
		n.logger.Infof("[Authorize] updating news category id=%s", actionData.ID.Hex())
		err := n.repo.Update(ctx, actionData.ID.Hex(), actionData.CategoryName)
		if err != nil {
			n.logger.Errorf("[Authorize] update news category action failed: %v", err)
			return nil, err
		}
		n.logger.Infof("[Authorize] news category updated successfully for id: %s", actionData.ID.Hex())
		return cpsAction, nil

	case string(constants.RequestDeleteNewsCategory):
		n.logger.Infof("[Authorize] deleting news category id=%s", actionData.ID.Hex())
		err := n.repo.Delete(ctx, actionData.ID.Hex())
		if err != nil {
			n.logger.Errorf("[Authorize] delete news category action failed: %v", err)
			return nil, err
		}
		n.logger.Infof("[Authorize] news category deleted successfully for id: %s", actionData.ID.Hex())
		return cpsAction, nil

	default:
		n.logger.Errorf("[Authorize] invalid action: %s", cpsAction.RequestAction)
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	n.logger.Infof("[Authorize] news category action authorized successfully: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

// CreateNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) CreateNewsCategory(ctx context.Context, categoryName []string) error {
	n.logger.Infof("[CreateNewsCategory] creating %d news categories", len(categoryName))
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		n.logger.Errorf("[CreateNewsCategory] incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}
	news_cat, err := n.repo.FindByNames(ctx, categoryName)
	code, _ := local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		if err != nil {
			n.logger.Errorf("[CreateNewsCategory] FindByNames repository error: %v", err)
			return err
		}
	}

	if news_cat != nil {
		n.logger.Errorf("[CreateNewsCategory] news category with name already exists")
		return fmt.Errorf("%s", localization.ErrorNewsCategoryWithNameAlreadyExists.Code)
	}
	cpsAction := lib.CpsModelBuilder("", makerData, nil, model.NewsCategoryCPSAction{CategoryNameList: categoryName}, string(constants.RequestCreateNewsCategory), constants.CREATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("[CreateNewsCategory] failed to create CPS action: %v", err)
		return err
	}
	n.logger.Infof("[CreateNewsCategory] CPS action created successfully for %d categories", len(categoryName))

	return nil
}

// DeleteNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) DeleteNewsCategory(ctx context.Context, id string) error {
	n.logger.Infof("[DeleteNewsCategory] deleting news category for id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		n.logger.Errorf("[DeleteNewsCategory] incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[DeleteNewsCategory] resource not found: %s", id)
		return fmt.Errorf("%s", localization.ErrorResourceNotFound.Code)
	} else if err != nil {
		n.logger.Errorf("[DeleteNewsCategory] repository Get error: %v", err)
		return err
	}

	updated_news_category := *news_category
	updated_news_category.IsDeleted = true
	updated_news_category.DeletedAt = time.Now()
	updated_news_category.LastModifiedAt = time.Now()

	cpsAction := lib.CpsModelBuilder(id, makerData, news_category, updated_news_category, string(constants.RequestDeleteNewsCategory), constants.DELETE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("[DeleteNewsCategory] failed to create CPS action: %v", err)
		return err
	}
	n.logger.Infof("[DeleteNewsCategory] CPS action created successfully for id: %s", id)

	return nil
}

// FindAllWithPagination implements service.NewsCategoryService.
func (n *NewsCategoryService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.NewsCategory], error) {
	categories, err := n.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		n.logger.Errorf("[FindAllWithPagination] failed to fetch news categories: %v", err)
		return nil, err
	}
	n.logger.Infof("[FindAllWithPagination] retrieved %d news categories", len(categories.Data))
	return categories, nil
}

// GetNewsCategoryByID implements service.NewsCategoryService.
func (n *NewsCategoryService) GetNewsCategoryByID(ctx context.Context, id string) (*model.NewsCategory, error) {
	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[GetNewsCategoryByID] resource not found: %s", id)
		return nil, fmt.Errorf("%s", code)
	} else if err != nil {
		n.logger.Errorf("[GetNewsCategoryByID] failed to get news category: %v", err)
		return nil, err
	}
	n.logger.Infof("[GetNewsCategoryByID] news category retrieved successfully for id: %s", id)
	return news_category, nil
}

// UpdateNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) UpdateNewsCategory(ctx context.Context, id string, categoryName string) error {
	n.logger.Infof("[UpdateNewsCategory] updating news category id=%s", id)
	makerData := local_util.ExtractUserFromContext(ctx)

	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("[UpdateNewsCategory] resource not found: %s", id)
		return fmt.Errorf("%s", code)
	} else if err != nil {
		n.logger.Errorf("[UpdateNewsCategory] failed to get news category: %v", err)
		return err
	}
	news_cat, err := n.repo.FindByNames(ctx, []string{categoryName})
	code, _ = local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		if err != nil {
			n.logger.Errorf("[UpdateNewsCategory] FindByNames repository error: %v", err)
			return err
		}
	}

	if news_cat != nil {
		n.logger.Errorf("[UpdateNewsCategory] news category with name already exists")
		return fmt.Errorf("%s", localization.ErrorNewsCategoryWithNameAlreadyExists.Code)
	}

	updated_news_category := *news_category
	updated_news_category.CategoryName = categoryName
	updated_news_category.LastModifiedAt = time.Now()

	cpsAction := lib.CpsModelBuilder(id, makerData, news_category, updated_news_category, string(constants.RequestUpdateNewsCategory), constants.UPDATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("[UpdateNewsCategory] failed to create CPS action: %v", err)
		return err
	}
	n.logger.Infof("[UpdateNewsCategory] CPS action created successfully for id: %s", id)

	return nil
}
