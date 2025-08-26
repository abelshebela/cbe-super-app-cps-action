package bankService

import (
	"cbe-super-app-cps-action/internal/constants"
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type BankService struct {
	cpsService service.CPSActionService
	logger     utils.Logger
	repo       storage.BankRepository
	cfg        *config.VaultConfig
	minio      config.MinioClientInterface
	bucketName string
}

func NewBankService(logger utils.Logger, repo storage.BankRepository, cpsService service.CPSActionService, minio config.MinioClientInterface, cfg *config.VaultConfig, bucketName string) service.BankService {
	return &BankService{
		logger:     logger,
		repo:       repo,
		cpsService: cpsService,
		cfg:        cfg,
		minio:      minio,
		bucketName: bucketName,
	}
}

func (b *BankService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	var actionData model.Bank
	raw, _ := bson.Marshal(cpsAction.CurrentAction)
	if err := bson.Unmarshal(raw, &actionData); err != nil {
		b.logger.Errorf("failed to unmarshal action data for authorization, action_code: %s", (cpsAction.ActionCode))
		return nil, fmt.Errorf("%s", localization.MsgBankInvalidRequestAction)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateBank):
		actionData.CreatedAt = time.Now()
		err := b.repo.Create(ctx, &actionData)
		if err != nil {
			b.logger.Errorf("Bank Create action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestDeleteBank):
		err := b.repo.Delete(ctx, actionData.ID)
		if err != nil {
			b.logger.Errorf("Bank Delete action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestEnableDisableBank):
		err := b.repo.EnableOrDisable(ctx, actionData.ID, actionData.Enabled)
		if err != nil {
			b.logger.Errorf("Bank Enable Disable action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestUpdateBankLogo):
		err := b.repo.Update(ctx, actionData.ID, &actionData)
		if err != nil {
			b.logger.Errorf("Bank update Logo action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestUpdateBank):
		err := b.repo.Update(ctx, actionData.ID, &actionData)
		if err != nil {
			b.logger.Errorf("Bank update action failed", "error", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s", localization.MsgBankInvalidRequestAction)
	}
	return cpsAction, nil
}

func (b *BankService) CreateOneBank(ctx context.Context, bank_request bank_dto.CreateBankRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		b.logger.Errorf("Create Bank failed incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	URL, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, bank_request.Logo, b.bucketName, b.cfg.MinioEndPoint, b.logger)

	if err != nil {
		b.logger.Errorf("UploadFileToMinio failed", "error", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	bank := model.Bank{
		Name: bank_request.Name,
		BIC:  bank_request.BIC,
		Code: bank_request.Code,
		Logo: URL,
	}

	searchParams := ""

	if bank_request.BIC != "" {
		searchParams += fmt.Sprintf("bic=%s&", bank_request.BIC)
	}
	if bank_request.Code != "" {
		searchParams += fmt.Sprintf("code=%s&", bank_request.Code)
	}
	if bank_request.Name != "" {
		searchParams += fmt.Sprintf("name=%s", bank_request.Name)
	}

	result, err := b.repo.FindByNameOrBICOrCode(ctx, bank_request.BIC, bank_request.Code, bank_request.Name)
	if err != nil && err.Error() != localization. {
		return err
	}

	if result != nil {
		if bank_request.BIC != "" && result.BIC == bank_request.BIC {
			return fmt.Errorf("%s", localization.ErrorBankWithBICAlreadyExists.Code)
		}
		if bank_request.Code != "" && result.Code == bank_request.Code {
			return fmt.Errorf("%s", localization.ErrorBankWithCodeAlreadyExists.Code)
		}
		if bank_request.Name != "" && result.Name == bank_request.Name {
			return fmt.Errorf("%s", localization.ErrorBankWithNameAlreadyExists.Code)
		}
	}

	action := lib.CpsModelBuilder("", makerData, nil, bank, string(constants.RequestCreateBank), constants.CREATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}

func (b *BankService) DeleteOneBank(ctx context.Context, id string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		b.logger.Errorf("Delete Bank failed incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}
	bank, err := b.repo.FindByID(ctx, id)

	if err != nil {
		return err
	}

	newBankData := *bank

	newBankData.IsDeleted = true
	action := lib.CpsModelBuilder(id, makerData, bank, newBankData, string(constants.RequestDeleteBank), constants.DELETE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}

func (b *BankService) EnableOrDisableBank(ctx context.Context, id string, enableDisable bool) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		b.logger.Errorf("Create Bank failed incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	bank, err := b.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if bank.Enabled && enableDisable {
		return fmt.Errorf("%s", localization.ErrorBankAlreadyEnabled.Code)
	} else if !bank.Enabled && !enableDisable {
		return fmt.Errorf("%s", localization.ErrorBankAlreadyDisabled.Code)
	}

	newBankData := *bank
	newBankData.Enabled = enableDisable

	action := lib.CpsModelBuilder(id, makerData, bank, newBankData, string(constants.RequestEnableDisableBank), constants.UPDATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}

func (b *BankService) GetAllBank(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Bank], error) {
	return b.repo.FindAllWithPagination(ctx, *filterParams)
}

func (b *BankService) GetOneBank(ctx context.Context, id string) (*model.Bank, error) {
	return b.repo.FindByID(ctx, id)
}

func (b *BankService) UpdateLogo(ctx context.Context, id string, logo bank_dto.UpdateLogo) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	bank, err := b.repo.FindByID(ctx, id)

	if err != nil {
		return err
	}

	URL, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, logo.Logo, b.bucketName, b.cfg.MinioEndPoint, b.logger)

	if err != nil {
		b.logger.Errorf("UploadFileToMinio failed", "error", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	newBankData := *bank

	newBankData.Logo = URL

	action := lib.CpsModelBuilder(id, makerData, bank, newBankData, string(constants.RequestUpdateBankLogo), constants.UPDATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}

func (b *BankService) UpdateOneBank(ctx context.Context, id string, bank_request bank_dto.UpdateBankRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	bank, err := b.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	updatedBank := *bank
	searchParams := ""

	if bank_request.BIC != "" {
		updatedBank.BIC = bank_request.BIC
		searchParams += fmt.Sprintf("bic=%s&", bank_request.BIC)
	}
	if bank_request.Code != "" {
		updatedBank.Code = bank_request.Code
		searchParams += fmt.Sprintf("code=%s&", bank_request.Code)
	}
	if bank_request.Name != "" {
		updatedBank.Name = bank_request.Name
		searchParams += fmt.Sprintf("name=%s", bank_request.Name)
	}

	result, err := b.repo.FindByNameOrBICOrCode(ctx, bank_request.BIC, bank_request.Code, bank_request.Name)
	if err != nil && err.Error() != localization.MsgFileNotFound {
		return err
	}

	if result != nil {
		if bank_request.BIC != "" && result.BIC == bank_request.BIC {
			return fmt.Errorf("%s", localization.ErrorBankWithBICAlreadyExists.Code)
		}
		if bank_request.Code != "" && result.Code == bank_request.Code {
			return fmt.Errorf("%s", localization.ErrorBankWithCodeAlreadyExists.Code)
		}
		if bank_request.Name != "" && result.Name == bank_request.Name {
			return fmt.Errorf("%s", localization.ErrorBankWithNameAlreadyExists.Code)
		}
	}

	action := lib.CpsModelBuilder(id, makerData, bank, updatedBank, string(constants.RequestUpdateBank), constants.UPDATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}
