package budget_category

import (
	"cbe-super-app-cps-action/internal/constants"
	budget_category_dto "cbe-super-app-cps-action/internal/constants/dto/budget_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	config "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BudgetCategoryService struct {
	budgetCategoryRepo storage.BudgetCategoryRepository
	cpsService         service.CPSActionService
	logger             utils.Logger
	minio              *s3.Client
	bucketName         string
	cfg                *config.VaultConfig
	minioEndPoint      string
}

func NewBudgetCategoryService(
	budgetCategoryRepo storage.BudgetCategoryRepository,
	cpsService service.CPSActionService,
	logger utils.Logger,
	minio *s3.Client,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
) service.BudgetCategoryService {
	return &BudgetCategoryService{
		budgetCategoryRepo: budgetCategoryRepo,
		cpsService:         cpsService,
		logger:             logger,
		minio:              minio,
		bucketName:         bucketName,
		cfg:                cfg,
		minioEndPoint:      minioEndPoint,
	}
}

func (b *BudgetCategoryService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	b.logger.Infof("[Authorize] authorizing budget category action: %s", action.RequestAction)
	var err error
	budgetCategory, marshal_err := local_util.JsonUnmarshal[model.BudgetCategory](action.CurrentAction)
	if marshal_err != nil || budgetCategory == nil {
		b.logger.Errorf("[Authorize] failed to unmarshal current action: %v", marshal_err)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	switch action.RequestAction {
	case string(constants.RequestCreateBudgetCategory):
		err = b.budgetCategoryRepo.CreateBudgetCategory(ctx, budgetCategory)
		if err != nil {
			b.logger.Errorf("[Authorize] failed to create budget category: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] budget category created successfully")
	case string(constants.RequestUpdateBudgetCategory):
		err = b.budgetCategoryRepo.UpdateBudgetCategory(ctx, action.UniqueId, budgetCategory)
		if err != nil {
			b.logger.Errorf("[Authorize] failed to update budget category: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] budget category updated successfully for id: %s", action.UniqueId)
	case string(constants.RequestDisableBudgetCategory):
		err = b.budgetCategoryRepo.UpdateBudgetCategory(ctx, action.UniqueId, budgetCategory)
		if err != nil {
			b.logger.Errorf("[Authorize] failed to disable budget category: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] budget category disabled successfully for id: %s", action.UniqueId)
	case string(constants.RequestEnableBudgetCategory):
		err = b.budgetCategoryRepo.UpdateBudgetCategory(ctx, action.UniqueId, budgetCategory)
		if err != nil {
			b.logger.Errorf("[Authorize] failed to enable budget category: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] budget category enabled successfully for id: %s", action.UniqueId)
	case string(constants.RequestDeleteBudgetCategory):
		err = b.budgetCategoryRepo.DeleteBudgetCategory(ctx, action.UniqueId)
		if err != nil {
			b.logger.Errorf("[Authorize] failed to delete budget category: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] budget category deleted successfully for id: %s", action.UniqueId)
	default:
		b.logger.Errorf("[Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	action.ActionStatus = (string)(constants.Approved)
	b.logger.Infof("[Authorize] budget category action authorized successfully: %s", action.RequestAction)
	return action, nil
}

func (b *BudgetCategoryService) CreateBudgetCategory(ctx context.Context, req budget_category_dto.CreateBudgetRequest) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	iconURL := ""
	if req.Icon != nil {
		url, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, req.Icon, string(constants.BudgetCategoryIcon), *b.cfg, "", b.logger)
		if err != nil {
			b.logger.Errorf("[CreateBudgetCategory] failed to upload icon: %v", err)
			return err
		}
		iconURL = url
	}
	isDuplicate, err := b.budgetCategoryRepo.FindByName(ctx, req.Name)
	if err != nil {
		b.logger.Errorf("[CreateBudgetCategory] failed to check for duplicate budget category name: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	if isDuplicate != nil {
		b.logger.Errorf("[CreateBudgetCategory] duplicate budget category found")
		return errors.New(localization.ErrorBudgetCategoryNameAlreadyExists.Code)
	}

	budgetCategory := &model.BudgetCategory{
		Name:      req.Name,
		Color:     req.Color,
		Icon:      iconURL,
		Enabled:   true,
		IsDeleted: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, budgetCategory, string(constants.RequestCreateBudgetCategory), constants.CREATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		b.logger.Errorf("[CreateBudgetCategory] failed to create CPS action: %v", err)
		return err
	}

	b.logger.Infof("[CreateBudgetCategory] budget category creation request created successfully")
	return nil
}

func (b *BudgetCategoryService) FetchBudgetCategory(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]budget_category_dto.BudgetCategoryResponse], error) {
	budgetCategories, err := b.budgetCategoryRepo.FindAllBudgetCategories(ctx, filterParams)
	if err != nil {
		b.logger.Errorf("[FetchBudgetCategory] failed to fetch budget categories: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var responseData []budget_category_dto.BudgetCategoryResponse
	for _, budgetCategory := range budgetCategories.Data {
		responseData = append(responseData, budget_category_dto.BudgetCategoryResponse{
			ID:        budgetCategory.ID.Hex(),
			Name:      budgetCategory.Name,
			Color:     budgetCategory.Color,
			Icon:      budgetCategory.Icon,
			Enabled:   budgetCategory.Enabled,
			CreatedAt: budgetCategory.CreatedAt.Format(time.RFC3339),
			UpdatedAt: budgetCategory.UpdatedAt.Format(time.RFC3339),
		})
	}

	b.logger.Infof("[FetchBudgetCategory] retrieved %d budget categories", len(responseData))
	return &types.PaginatedResponse[[]budget_category_dto.BudgetCategoryResponse]{
		Data: responseData,
		Meta: budgetCategories.Meta,
	}, nil
}

func (b *BudgetCategoryService) FetchBudgetCategoryByID(ctx context.Context, id string) (*budget_category_dto.BudgetCategoryResponse, error) {
	budgetCategory, err := b.budgetCategoryRepo.FindBudgetCategoryByID(ctx, id)
	if err != nil {
		b.logger.Errorf("[FetchBudgetCategoryByID] failed to fetch budget category: %v", err)
		return nil, err
	}
	b.logger.Infof("[FetchBudgetCategoryByID] budget category retrieved successfully for id: %s", id)

	return &budget_category_dto.BudgetCategoryResponse{
		ID:        budgetCategory.ID.Hex(),
		Name:      budgetCategory.Name,
		Color:     budgetCategory.Color,
		Icon:      budgetCategory.Icon,
		Enabled:   budgetCategory.Enabled,
		CreatedAt: budgetCategory.CreatedAt.Format(time.RFC3339),
		UpdatedAt: budgetCategory.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (b *BudgetCategoryService) UpdateBudgetCategory(ctx context.Context, id string, req budget_category_dto.UpdateBudgetRequest) error {
	b.logger.Infof("[UpdateBudgetCategory] updating budget category for id: %s", id)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existingBudgetCategory, err := b.budgetCategoryRepo.FindBudgetCategoryByID(ctx, id)
	if err != nil {
		b.logger.Errorf("[UpdateBudgetCategory] failed to find budget category: %v", err)
		return err
	}

	newBudgetCategory := *existingBudgetCategory

	if req.Name != "" {

		isDuplicate, err := b.budgetCategoryRepo.FindByName(ctx, req.Name)
		if err != nil {
			b.logger.Errorf("[UpdateBudgetCategory] failed to check for duplicate budget category name: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if isDuplicate != nil && isDuplicate.ID.Hex() != id {
			b.logger.Errorf("[UpdateBudgetCategory] duplicate budget category found")
			return errors.New(localization.ErrorBudgetCategoryNameAlreadyExists.Code)
		}

		newBudgetCategory.Name = req.Name
	}
	if req.Color != "" {
		newBudgetCategory.Color = req.Color
	}
	if req.Icon != nil {
		var objectkey string
		if existingBudgetCategory.Icon != "" {
			objectkey = path.Base(existingBudgetCategory.Icon)
		}

		url, err := lib.UploadFileToMinio(
			ctx,
			b.minio,
			b.bucketName,
			req.Icon,
			string(constants.BudgetCategoryIcon),
			*b.cfg,
			objectkey,
			b.logger,
		)
		if err != nil {
			b.logger.Errorf("[UpdateBudgetCategory] failed to upload icon: %v", err)
			return err
		}
		newBudgetCategory.Icon = url
	}
	if newBudgetCategory == *existingBudgetCategory {
		return errors.New(localization.ErrorNoChangesToUpdate.Code)
	}

	newBudgetCategory.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(
		id,
		makerUser,
		existingBudgetCategory,
		&newBudgetCategory,
		string(constants.RequestUpdateBudgetCategory),
		constants.UPDATE,
	)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		return err
	}

	return nil
}

func (b *BudgetCategoryService) DeleteBudgetCategory(ctx context.Context, id string) error {
	b.logger.Infof("[DeleteBudgetCategory] deleting budget category for id: %s", id)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existingBudgetCategory, err := b.budgetCategoryRepo.FindBudgetCategoryByID(ctx, id)
	if err != nil {
		b.logger.Errorf("[DeleteBudgetCategory] failed to find budget category: %v", err)
		return err
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existingBudgetCategory, existingBudgetCategory, string(constants.RequestDeleteBudgetCategory), constants.DELETE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		b.logger.Errorf("[DeleteBudgetCategory] failed to create CPS action: %v", err)
		return err
	}
	b.logger.Infof("[DeleteBudgetCategory] budget category deletion request created successfully for id: %s", id)
	return nil
}

func (b *BudgetCategoryService) EnableOrDisableBudgetCategory(ctx context.Context, id string, enable bool) error {
	b.logger.Infof("[EnableOrDisableBudgetCategory] processing budget category enable/disable for id: %s, enabled: %v", id, enable)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existingBudgetCategory, err := b.budgetCategoryRepo.FindBudgetCategoryByID(ctx, id)
	if err != nil {
		b.logger.Errorf("[EnableOrDisableBudgetCategory] failed to find budget category: %v", err)
		return err
	}

	if enable && enable == existingBudgetCategory.Enabled {
		b.logger.Errorf("[EnableOrDisableBudgetCategory] budget category already enabled")
		return errors.New(localization.ErrorBudgetCategoryAlreadyEnabled.Code)
	}
	if !enable && enable == existingBudgetCategory.Enabled {
		b.logger.Errorf("[EnableOrDisableBudgetCategory] budget category already disabled")
		return errors.New(localization.ErrorBudgetCategoryAlreadyDisabled.Code)
	}

	update := *existingBudgetCategory
	update.Enabled = enable
	update.UpdatedAt = time.Now()

	var requestType string
	if enable {
		requestType = string(constants.RequestEnableBudgetCategory)
	} else {
		requestType = string(constants.RequestDisableBudgetCategory)
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existingBudgetCategory, update, requestType, constants.UPDATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		return err
	}

	return nil
}
