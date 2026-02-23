package budget_category

import (
	"cbe-super-app-cps-action/internal/constants"
	budget_category_dto "cbe-super-app-cps-action/internal/constants/dto/budget_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"path"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	config "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Budget Category", "Authorize")
	defer span.End()

	b.logger.Infof("[BudgetCatSvc][Authorize] action: %s", action.RequestAction)
	var err error
	budgetCategory, marshal_err := local_util.JsonUnmarshal[model.BudgetCategory](action.CurrentAction)
	if marshal_err != nil || budgetCategory == nil {
		span.AddEvent("[Authorize] failed to unmarshal current action", trace.WithAttributes(
			attribute.String("error", marshal_err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		b.logger.Errorf("[BudgetCatSvc][Authorize] unmarshal err: %v", marshal_err)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	switch action.RequestAction {
	case string(constants.RequestCreateBudgetCategory):
		err = b.budgetCategoryRepo.CreateBudgetCategory(ctx, budgetCategory)
		if err != nil {
			span.AddEvent("[Authorize] failed to create budget category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			b.logger.Errorf("[BudgetCatSvc][Authorize] create err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BudgetCatSvc][Authorize] created")
	case string(constants.RequestUpdateBudgetCategory):
		err = b.budgetCategoryRepo.UpdateBudgetCategory(ctx, action.UniqueId, budgetCategory)
		if err != nil {
			span.AddEvent("[Authorize] failed to update budget category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			b.logger.Errorf("[BudgetCatSvc][Authorize] update err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BudgetCatSvc][Authorize] updated id: %s", action.UniqueId)
	case string(constants.RequestDisableBudgetCategory):
		err = b.budgetCategoryRepo.UpdateBudgetCategory(ctx, action.UniqueId, budgetCategory)
		if err != nil {
			span.AddEvent("[Authorize] failed to disable budget category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			b.logger.Errorf("[BudgetCatSvc][Authorize] disable err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BudgetCatSvc][Authorize] disabled id: %s", action.UniqueId)
	case string(constants.RequestEnableBudgetCategory):
		err = b.budgetCategoryRepo.UpdateBudgetCategory(ctx, action.UniqueId, budgetCategory)
		if err != nil {
			span.AddEvent("[Authorize] failed to enable budget category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			b.logger.Errorf("[BudgetCatSvc][Authorize] enable err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BudgetCatSvc][Authorize] enabled id: %s", action.UniqueId)
	case string(constants.RequestDeleteBudgetCategory):
		err = b.budgetCategoryRepo.DeleteBudgetCategory(ctx, action.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] failed to delete budget category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			b.logger.Errorf("[BudgetCatSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BudgetCatSvc][Authorize] deleted id: %s", action.UniqueId)
	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", action.RequestAction)))
		b.logger.Errorf("[BudgetCatSvc][Authorize] unsupported: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	action.ActionStatus = (string)(constants.Approved)
	b.logger.Infof("[BudgetCatSvc][Authorize] done: %s", action.RequestAction)
	return action, nil
}

func (b *BudgetCategoryService) CreateBudgetCategory(ctx context.Context, req budget_category_dto.CreateBudgetRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateBudgetCategory", "Budget Category", "CreateBudgetCategory")
	defer span.End()

	makerUser := local_util.ExtractUserFromContext(ctx)

	iconURL := ""
	if req.Icon != nil {
		url, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, req.Icon, string(constants.BudgetCategoryFolderName), *b.cfg, "", b.logger)
		if err != nil {
			span.AddEvent("[CreateBudgetCategory] failed to upload icon", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("name", req.Name),
			))
			b.logger.Errorf("[BudgetCatSvc][Create] upload icon err: %v", err)
			return err
		}
		iconURL = url
	}
	isDuplicate, err := b.budgetCategoryRepo.FindByName(ctx, req.Name)
	if err != nil {
		span.AddEvent("[CreateBudgetCategory] failed to check for duplicate", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("name", req.Name),
		))
		b.logger.Errorf("[BudgetCatSvc][Create] dup check err: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	if isDuplicate != nil {
		span.AddEvent("[CreateBudgetCategory] duplicate budget category found", trace.WithAttributes(attribute.String("name", req.Name)))
		b.logger.Errorf("[BudgetCatSvc][Create] name exists")
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
		span.AddEvent("[CreateBudgetCategory] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("name", req.Name),
		))
		b.logger.Errorf("[BudgetCatSvc][Create] cps action err: %v", err)
		return err
	}

	b.logger.Infof("[BudgetCatSvc][Create] request created")
	return nil
}

func (b *BudgetCategoryService) FetchBudgetCategory(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]budget_category_dto.BudgetCategoryResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchBudgetCategory", "Budget Category", "FetchBudgetCategory")
	defer span.End()

	budgetCategories, err := b.budgetCategoryRepo.FindAllBudgetCategories(ctx, filterParams)
	if err != nil {
		span.AddEvent("[FetchBudgetCategory] failed to fetch budget categories", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		b.logger.Errorf("[BudgetCatSvc][FetchAll] err: %v", err)
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

	b.logger.Infof("[BudgetCatSvc][FetchAll] count: %d", len(responseData))
	return &types.PaginatedResponse[[]budget_category_dto.BudgetCategoryResponse]{
		Data: responseData,
		Meta: budgetCategories.Meta,
	}, nil
}

func (b *BudgetCategoryService) FetchBudgetCategoryByID(ctx context.Context, id string) (*budget_category_dto.BudgetCategoryResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchBudgetCategoryByID", "Budget Category", "FetchBudgetCategoryByID")
	defer span.End()

	budgetCategory, err := b.budgetCategoryRepo.FindBudgetCategoryByID(ctx, id)
	if err != nil {
		span.AddEvent("[FetchBudgetCategoryByID] failed to fetch budget category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[BudgetCatSvc][FetchByID] err: %v", err)
		return nil, err
	}
	b.logger.Infof("[BudgetCatSvc][FetchByID] found id: %s", id)

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
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateBudgetCategory", "Budget Category", "UpdateBudgetCategory")
	defer span.End()

	b.logger.Infof("[BudgetCatSvc][Update] id: %s", id)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existingBudgetCategory, err := b.budgetCategoryRepo.FindBudgetCategoryByID(ctx, id)
	if err != nil {
		span.AddEvent("[UpdateBudgetCategory] failed to find budget category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[BudgetCatSvc][Update] find err: %v", err)
		return err
	}

	newBudgetCategory := *existingBudgetCategory

	if req.Name != "" {

		isDuplicate, err := b.budgetCategoryRepo.FindByName(ctx, req.Name)
		if err != nil {
			span.AddEvent("[UpdateBudgetCategory] failed to check for duplicate", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			b.logger.Errorf("[BudgetCatSvc][Update] dup check err: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if isDuplicate != nil && isDuplicate.ID.Hex() != id {
			span.AddEvent("[UpdateBudgetCategory] duplicate budget category found", trace.WithAttributes(attribute.String("name", req.Name)))
			b.logger.Errorf("[BudgetCatSvc][Update] name exists")
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
			string(constants.BudgetCategoryFolderName),
			*b.cfg,
			objectkey,
			b.logger,
		)
		if err != nil {
			span.AddEvent("[UpdateBudgetCategory] failed to upload icon", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			b.logger.Errorf("[BudgetCatSvc][Update] upload icon err: %v", err)
			return err
		}
		newBudgetCategory.Icon = url
	}
	if newBudgetCategory == *existingBudgetCategory {
		span.AddEvent("[UpdateBudgetCategory] no changes to update", trace.WithAttributes(attribute.String("id", id)))
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
		span.AddEvent("[UpdateBudgetCategory] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}

func (b *BudgetCategoryService) DeleteBudgetCategory(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteBudgetCategory", "Budget Category", "DeleteBudgetCategory")
	defer span.End()

	b.logger.Infof("[BudgetCatSvc][Delete] id: %s", id)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existingBudgetCategory, err := b.budgetCategoryRepo.FindBudgetCategoryByID(ctx, id)
	if err != nil {
		span.AddEvent("[DeleteBudgetCategory] failed to find budget category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[BudgetCatSvc][Delete] find err: %v", err)
		return err
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, existingBudgetCategory, existingBudgetCategory, string(constants.RequestDeleteBudgetCategory), constants.DELETE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		span.AddEvent("[DeleteBudgetCategory] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[BudgetCatSvc][Delete] cps action err: %v", err)
		return err
	}
	b.logger.Infof("[BudgetCatSvc][Delete] request created id: %s", id)
	return nil
}

func (b *BudgetCategoryService) EnableOrDisableBudgetCategory(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableBudgetCategory", "Budget Category", "EnableOrDisableBudgetCategory")
	defer span.End()

	b.logger.Infof("[BudgetCatSvc][EnableDisable] id: %s enable: %v", id, enable)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existingBudgetCategory, err := b.budgetCategoryRepo.FindBudgetCategoryByID(ctx, id)
	if err != nil {
		span.AddEvent("[EnableOrDisableBudgetCategory] failed to find budget category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[BudgetCatSvc][EnableDisable] find err: %v", err)
		return err
	}

	if enable && enable == existingBudgetCategory.Enabled {
		span.AddEvent("[EnableOrDisableBudgetCategory] budget category already enabled", trace.WithAttributes(attribute.String("id", id)))
		b.logger.Errorf("[BudgetCatSvc][EnableDisable] already enabled")
		return errors.New(localization.ErrorBudgetCategoryAlreadyEnabled.Code)
	}
	if !enable && enable == existingBudgetCategory.Enabled {
		span.AddEvent("[EnableOrDisableBudgetCategory] budget category already disabled", trace.WithAttributes(attribute.String("id", id)))
		b.logger.Errorf("[BudgetCatSvc][EnableDisable] already disabled")
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
		span.AddEvent("[EnableOrDisableBudgetCategory] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}
