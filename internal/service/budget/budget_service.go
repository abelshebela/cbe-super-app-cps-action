package budget

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"mime/multipart"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BudgetService struct {
	colorRepo   storage.ColorRepository
	iconRepo    storage.IconRepository
	cpsService  service.CPSActionService
	bucketName  string
	minio       config.MinioClientInterface
	minioPubUrl string
	logger      utils.Logger
	cfg         *config.VaultConfig
}

// Constructor
func NewBudgetService(
	iconRepo storage.IconRepository,
	colorRepo storage.ColorRepository,
	cpsService service.CPSActionService,
	bucketName string,
	minio config.MinioClientInterface,
	minioPubUrl string,
	cfg *config.VaultConfig,
	logger utils.Logger,
) service.BudgetService {
	return &BudgetService{
		iconRepo:    iconRepo,
		colorRepo:   colorRepo,
		cpsService:  cpsService,
		bucketName:  bucketName,
		minio:       minio,
		minioPubUrl: minioPubUrl,
		logger:      logger,
		cfg:         cfg,
	}
}

func (b *BudgetService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	color, ok := action.CurrentAction.(*model.Color)
	if !ok || color == nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	icon, ok := action.CurrentAction.(*model.Icon)
	if !ok || icon == nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	var err error
	switch action.RequestAction {
	case string(constants.RequestCreateBudgetColor):
		err = b.colorRepo.Create(ctx, color)
	case string(constants.RequestUpdateBudgetColor):
		err = b.colorRepo.Update(ctx, color.ID.Hex(), color)
	case string(constants.RequestCreateBudgetIcon):
		err = b.iconRepo.Create(ctx, icon)
	case string(constants.RequestUpdateBudgetIcon):
		err = b.iconRepo.Update(ctx, icon.ID.Hex(), icon)
	default:
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		return nil, err
	}

	action.CurrentAction = color
	action.CurrentAction = icon
	return action, nil
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

	iconURL, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, fileHeader, "budget-icons", b.minioPubUrl, b.logger)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

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
	icons, err := b.iconRepo.FindAllWithPagination(ctx, filterParams)
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

	iconURL, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, fileHeader, "budget-icons", b.minioPubUrl, b.logger)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsAction.CurrentAction = model.Icon{
		Icon:      iconURL,
		CreatedAt: time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, cpsAction.CurrentAction, string(constants.RequestUpdateBudgetIcon), constants.UPDATE)

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

	exists, err := b.colorRepo.FindByID(ctx, color.ID.Hex())
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			return err
		}
	}
	if exists != nil {
		return errors.New(localization.ErrorDuplicateColorExists.Code)
	}

	color.CreatedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, color, string(constants.RequestCreateBudgetColor), constants.CREATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		return err
	}

	return nil
}

func (b *BudgetService) BudgetFetchColors(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Color], error) {
	colors, err := b.colorRepo.FindAllWithPagination(ctx, filterParams)
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

	exists, err := b.colorRepo.Find(ctx, bson.M{"color": color.Color, "is_deleted": false})
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			return err
		}
	}

	if exists != nil {
		return errors.New(localization.ErrorDuplicateColorExists.Code)
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
