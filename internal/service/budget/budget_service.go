package budget

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/budget/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"mime/multipart"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BudgetService struct {
	repo       storage.BudgetRepository
	cpsService service.CPSActionService
	bucketName string
	minio      config.MinioClientInterface
	logger     utils.Logger
	cfg        *config.VaultConfig
}

// Constructor
func NewBudgetService(
	repo storage.BudgetRepository,
	cpsService service.CPSActionService,
	bucketName string,
	minio config.MinioClientInterface,
	cfg *config.VaultConfig,
	logger utils.Logger,
) service.BudgetService {
	return &BudgetService{
		repo:       repo,
		cpsService: cpsService,
		bucketName: bucketName,
		minio:      minio,
		logger:     logger,
		cfg:        cfg,
	}
}

func (b *BudgetService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	err := b.repo.AuthorizeCPSAction(ctx, cpsAction)
	if err != nil {
		return nil, err
	}
	return cpsAction, nil
}

func (b *BudgetService) CreateBudgetIcon(ctx context.Context, fileHeader *multipart.FileHeader, file *multipart.File) error {
	makerUser := local_util.ExtractUserFromContext(ctx)
	cpsAction := model.CPSAction{
		ActionCode:    utils.RandomGenerator(20),
		CurrentAction: make(map[string]interface{}),
	}
	_, ok := cpsAction.CurrentAction.(map[string]interface{})
	if !ok {
		b.logger.Errorf("invalid current action type")
		return nil
	}

	saveObj, err := core.SaveIconToMinio(ctx, fileHeader, b.bucketName, b.minio, b.logger)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	iconURL := "https://" + b.cfg.MinioEndPoint + "/" + saveObj.Bucket + "/" + saveObj.Key
	cpsAction.CurrentAction = model.Icon{
		Icon:      iconURL,
		CreatedAt: time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, cpsAction.CurrentAction, string(constants.RequestCreateBudgetIcon), constants.CREATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		return err
	}

	return nil
}
func (b *BudgetService) BudgetFetchIcons(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Icon], error) {
	icons, err := b.repo.FetchIcons(ctx, filterParams)
	if err != nil {
		b.logger.Errorf("failed to fetch budget icons: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return icons, nil
}

func (b *BudgetService) BudgetUpdateIcon(ctx context.Context, id string, fileHeader *multipart.FileHeader, file *multipart.File) error {
	makerUser := local_util.ExtractUserFromContext(ctx)
	cpsAction := model.CPSAction{
		ActionCode:    utils.RandomGenerator(20),
		CurrentAction: make(map[string]interface{}),
	}
	_, ok := cpsAction.CurrentAction.(map[string]interface{})
	if !ok {
		b.logger.Errorf("invalid current action type")
		return nil
	}

	saveObj, err := core.SaveIconToMinio(ctx, fileHeader, b.bucketName, b.minio, b.logger)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	iconURL := "https://" + b.cfg.MinioEndPoint + "/" + saveObj.Bucket + "/" + saveObj.Key
	cpsAction.CurrentAction = model.Icon{
		Icon:      iconURL,
		CreatedAt: time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, cpsAction.CurrentAction, string(constants.RequestUpdateBudgetIcon), constants.CREATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		return err
	}

	return nil
}

func (b *BudgetService) BudgetCreateColor(ctx context.Context, color *model.Color) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	if color.Color == "" {
		b.logger.Errorf("color cannot be empty")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}

	color.CreatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, color, string(constants.RequestCreateBudgetColor), constants.CREATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		return err
	}

	return nil
}

func (b *BudgetService) BudgetFetchColors(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Color], error) {
	colors, err := b.repo.FetchColors(ctx, filterParams)
	if err != nil {
		b.logger.Errorf("failed to fetch budget colors: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return colors, nil
}

func (b *BudgetService) BudgetUpdateColor(ctx context.Context, id string, color *model.Color) error {
	makerUser := local_util.ExtractUserFromContext(ctx)

	if color.Color == "" {
		b.logger.Errorf("color cannot be empty")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}

	color.UpdatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, color, string(constants.RequestUpdateBudgetColor), constants.UPDATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		return err
	}

	return nil
}

func (b *BudgetService) BudgetCheckerApproval(ctx context.Context, actioncode string) error {
	action := &model.CPSAction{
		ActionCode: actioncode,
	}

	if err := b.cpsService.ApproveCPSAction(ctx, action); err != nil {
		b.logger.Errorf("failed to approve budget action: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}
