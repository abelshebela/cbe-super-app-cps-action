package newscategory_service

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"fmt"
	"time"

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
	n.logger.Infof("Authorize: starting for action id=%s request=%s", cpsAction.ID, string(cpsAction.RequestAction))

	actionData, err := local_util.JsonUnmarshal[model.NewsCategoryCPSAction](cpsAction.CurrentAction)
	if err != nil {
		n.logger.Errorf("Authorize: failed to unmarshal CurrentAction: %v", err)
		return nil, fmt.Errorf("failed to unmarshal CurrentAction to NewsCategoryCPSAction: %v", err)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateNewsCategory):
		n.logger.Infof("Authorize: handling CreateNewsCategory for categories=%v", actionData.CategoryNameList)
		err := n.repo.Create(ctx, actionData.CategoryNameList)
		if err != nil {
			n.logger.Errorf("Authorize: Create repository call failed: %v", err)
			return nil, err
		}
		n.logger.Infof("Authorize: CreateNewsCategory repository call succeeded")
	case string(constants.RequestUpdateNewsCategory):
		n.logger.Infof("Authorize: handling UpdateNewsCategory id=%s newName=%v", actionData.ID.Hex(), actionData.CategoryName)
		err := n.repo.Update(ctx, actionData.ID.Hex(), actionData.CategoryName)
		if err != nil {
			n.logger.Errorf("Authorize: Update repository call failed: %v", err)
		} else {
			n.logger.Infof("Authorize: UpdateNewsCategory repository call succeeded id=%s", actionData.ID.Hex())
		}
		return nil, err

	case string(constants.RequestDeleteNewsCategory):
		n.logger.Infof("Authorize: handling DeleteNewsCategory id=%s", actionData.ID.Hex())
		err := n.repo.Delete(ctx, actionData.ID.Hex())
		if err != nil {
			n.logger.Errorf("Authorize: Delete repository call failed: %v", err)
		} else {
			n.logger.Infof("Authorize: DeleteNewsCategory repository call succeeded id=%s", actionData.ID.Hex())
		}
		return nil, err

	default:
		n.logger.Errorf("Authorize: invalid action %s", cpsAction.RequestAction)
		return nil, fmt.Errorf("%s", localization.MsgInvalidAction)
	}
	n.logger.Infof("Authorize: completed for action id=%s", cpsAction.ID)
	return cpsAction, nil
}

// CreateNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) CreateNewsCategory(ctx context.Context, categoryName []string) error {
	n.logger.Infof("CreateNewsCategory: starting for categories=%v", categoryName)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		n.logger.Errorf("CreateNewsCategory: incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}
	news_cat, err := n.repo.FindByNames(ctx, categoryName)
	code, _ := local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		if err != nil {
			n.logger.Errorf("CreateNewsCategory: FindByNames repository error: %v", err)
			return err
		}
	}

	if news_cat != nil {
		n.logger.Errorf("CreateNewsCategory: news category with name already exists: %v", categoryName)
		return fmt.Errorf("%s", localization.ErrorNewsCategoryWithNameAlreadyExists.Code)
	}
	cpsAction := lib.CpsModelBuilder("", makerData, nil, model.NewsCategoryCPSAction{CategoryNameList: categoryName}, string(constants.RequestCreateNewsCategory), constants.CREATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("CreateNewsCategory: failed to create CPS action: %v", err)
		return err
	}
	n.logger.Infof("CreateNewsCategory: CPS action created successfully for categories=%v", categoryName)

	return nil
}

// DeleteNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) DeleteNewsCategory(ctx context.Context, id string) error {
	n.logger.Infof("DeleteNewsCategory: starting for id=%s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		n.logger.Errorf("DeleteNewsCategory: incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("DeleteNewsCategory: resource not found id=%s", id)
		return fmt.Errorf("%s", localization.ErrorResourceNotFound.Code)
	} else if err != nil {
		n.logger.Errorf("DeleteNewsCategory: repository Get error for id=%s: %v", id, err)
		return err
	}

	updated_news_category := *news_category
	updated_news_category.IsDeleted = true
	updated_news_category.DeletedAt = time.Now()
	updated_news_category.LastModifiedAt = time.Now()

	cpsAction := lib.CpsModelBuilder(id, makerData, news_category, updated_news_category, string(constants.RequestDeleteNewsCategory), constants.DELETE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("DeleteNewsCategory: failed to create CPS action for id=%s: %v", id, err)
		return err
	}
	n.logger.Infof("DeleteNewsCategory: CPS action created successfully for id=%s", id)

	return nil
}

// FindAllWithPagination implements service.NewsCategoryService.
func (n *NewsCategoryService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.NewsCategory], error) {
	n.logger.Infof("FindAllWithPagination: starting with filter=%+v", filter)
	categories, err := n.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		n.logger.Errorf("FindAllWithPagination: repository error: %v", err)
		return nil, err
	}
	n.logger.Infof("FindAllWithPagination: found %d categories", len(categories.Data))
	return categories, nil
}

// GetNewsCategoryByID implements service.NewsCategoryService.
func (n *NewsCategoryService) GetNewsCategoryByID(ctx context.Context, id string) (*model.NewsCategory, error) {
	n.logger.Infof("GetNewsCategoryByID: starting for id=%s", id)
	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("GetNewsCategoryByID: resource not found id=%s", id)
		return nil, fmt.Errorf("%s", code)
	} else if err != nil {
		n.logger.Errorf("GetNewsCategoryByID: repository Get error id=%s: %v", id, err)
		return nil, err
	}
	n.logger.Infof("GetNewsCategoryByID: success id=%s", id)
	return news_category, nil
}

// UpdateNewsCategory implements service.NewsCategoryService.
func (n *NewsCategoryService) UpdateNewsCategory(ctx context.Context, id string, categoryName string) error {
	n.logger.Infof("UpdateNewsCategory: starting id=%s newName=%s", id, categoryName)
	makerData := local_util.ExtractUserFromContext(ctx)

	news_category, err := n.repo.Get(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		n.logger.Errorf("UpdateNewsCategory: resource not found id=%s", id)
		return fmt.Errorf("%s", code)
	} else if err != nil {
		n.logger.Errorf("UpdateNewsCategory: repository Get error id=%s: %v", id, err)
		return err
	}
	news_cat, err := n.repo.FindByNames(ctx, []string{categoryName})
	code, _ = local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code {
		if err != nil {
			n.logger.Errorf("UpdateNewsCategory: FindByNames repository error: %v", err)
			return err
		}
	}

	if news_cat != nil {
		n.logger.Errorf("CreateNewsCategory: news category with name already exists: %v", categoryName)
		return fmt.Errorf("%s", localization.ErrorNewsCategoryWithNameAlreadyExists.Code)
	}

	updated_news_category := *news_category
	updated_news_category.CategoryName = categoryName
	updated_news_category.LastModifiedAt = time.Now()

	cpsAction := lib.CpsModelBuilder(id, makerData, news_category, updated_news_category, string(constants.RequestUpdateNewsCategory), constants.UPDATE)

	if err := n.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		n.logger.Errorf("UpdateNewsCategory: failed to create CPS action for id=%s: %v", id, err)
		return err
	}
	n.logger.Infof("UpdateNewsCategory: CPS action created successfully for id=%s", id)

	return nil
}
