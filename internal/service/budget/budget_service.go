package budget

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/budget/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BudgetService struct {
	repo   storage.AmountBasedAuthRepository
	logger utils.Logger
	cpsService service.CPSActionService
	bucketName string
	minio      config.MinioClientInterface
	logger     utils.Logger
	cfg        *config.VaultConfig
}

func NewBudgetService(repo storage.AmountBasedAuthRepository,cpsService service.CPSActionService,bucketName string,minio config.MinioClientInterface,cfg config.VaultConfig, logger utils.Logger) service.BudgetService {
	return &BudgetService{
		repo:   repo,
		logger: logger,
		cpsService: cpsService,
		bucketName: bucketName,
		minio: minio,
		cfg: &cfg,
	}
}

func (b *BudgetService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	b.logger.Infof("Budget service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}



func (b *BudgetService) CreateBudgetIcon(ctx context.Context, fileHeader *multipart.FileHeader, file *multipart.File) error{
	makerUser := local_util.ExtractUserFromContext(ctx)
	cpsAction.ActionCode = utils.RandomGenerator(20)
	iconData, ok := cpsAction.CurrentAction.(map[string]interface{})
	if !ok {
		b.logger.Errorf("invalid current action type")
		return nil, fmt.Errorf("invalid currect action")
	}

	saveObj,err := core.SaveIconToMinio(ctx,fileHeader,b.bucketName,&b.minio,b.logger)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	iconURL := "https://" + b.cfg.MinioEndPoint + "/" + saveObj.Bucket + "/" + saveObj.Key
	iconData := model.Icon{
		Icon: iconURL,
		CreatedAt: time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder(nil,makerUser,nil,iconData,string(constants.RequestCreateBudgetIcon),constants.Create)

	if err := b.cpsService.CreateCPSAction(ctx,&cpsActionData); err != nil {
		return err
	}
	
	return nil
}
BudgetFetchIcons(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Icon], error)
BudgetUpdateIcon(ctx context.Context, fileHeader *multipart.FileHeader, file *multipart.File) error
BudgetCreateColor(ctx context.Context, color *model.Color) error
BudgetFetchColors(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Color], error)
BudgetUpdateColor(ctx context.Context, color *model.Color) error
BudgetCheckerApproval(ctx context.Context)
