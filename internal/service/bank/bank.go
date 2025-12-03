package bankService

import (
	"cbe-super-app-cps-action/internal/constants"
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	"cbe-super-app-cps-action/internal/constants/localization"
	bank_core "cbe-super-app-cps-action/internal/service/bank/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BankService struct {
	cpsService  service.CPSActionService
	logger      utils.Logger
	repo        storage.BankRepository
	cfg         *config.VaultConfig
	minio       *s3.Client
	minioPubUrl string
	bucketName  string
}

func NewBankService(logger utils.Logger, repo storage.BankRepository, cpsService service.CPSActionService, minio *s3.Client, minioPubUrl string, cfg *config.VaultConfig, bucketName string) service.BankService {
	return &BankService{
		logger:      logger,
		repo:        repo,
		cpsService:  cpsService,
		cfg:         cfg,
		minio:       minio,
		minioPubUrl: minioPubUrl,
		bucketName:  bucketName,
	}
}

func (b *BankService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		b.logger.Errorf("failed to marshal CurrentAction: %v\n", err)
		return nil, fmt.Errorf("failed to marshal CurrentAction: %v", err)
	}

	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		b.logger.Errorf("failed to unmarshal CurrentAction: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal to interface{}: %v", err)
	}

	actionData := bank_core.Bank_mapper(actionMap.(map[string]interface{}))
	if cpsAction.UniqueId != "" {
		objID, err := bson.ObjectIDFromHex(cpsAction.UniqueId)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		actionData.ID = objID
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
		err := b.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			b.logger.Errorf("Bank Delete action  failed", "error", err)
			return nil, err
		}

	case string(constants.RequestEnableDisableBank):
		err := b.repo.EnableOrDisable(ctx, cpsAction.UniqueId, actionData.Enabled)

		if err != nil {
			b.logger.Errorf("Bank Enable Disable action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestUpdateBankLogo):
		err := b.repo.Update(ctx, cpsAction.UniqueId, &actionData)
		if err != nil {
			b.logger.Errorf("Bank update Logo action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestUpdateBank):
		err := b.repo.Update(ctx, cpsAction.UniqueId, &actionData)
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

	URL, err := lib.UploadFileToMinio(ctx, b.minio, "", bank_request.Logo, b.bucketName, *b.cfg, "", b.logger)

	if err != nil {
		b.logger.Errorf("UploadFileToMinio failed", "error", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	bank := model.Bank{
		Name:    bank_request.Name,
		BIC:     bank_request.BIC,
		Code:    bank_request.Code,
		Logo:    URL,
		Enabled: true,
	}

	result, err := b.repo.FindByNameOrBICOrCode(ctx, bank_request.BIC, bank_request.Code, bank_request.Name)

	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code {
			return err
		}
	}

	if result != nil {
		if bank_request.BIC != "" && strings.EqualFold(result.BIC, bank_request.BIC) {
			return fmt.Errorf("%s", localization.ErrorBankWithBICAlreadyExists.Code)
		}
		if bank_request.Code != "" && strings.EqualFold(result.Code, bank_request.Code) {
			return fmt.Errorf("%s", localization.ErrorBankWithCodeAlreadyExists.Code)
		}
		if bank_request.Name != "" && strings.EqualFold(result.Name, bank_request.Name) {
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

	enable := string(constants.RequestEnableDisableBank)
	// if !enableDisable {
	// 	enable = string(constants.RequestDisableBank)
	// }

	action := lib.CpsModelBuilder(id, makerData, bank, newBankData, enable, constants.UPDATE)

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

	var objectkey string
	if bank.Logo != "" {
		objectkey = path.Base(bank.Logo)
	}

	URL, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, logo.Logo, b.bucketName, *b.cfg, objectkey, b.logger)

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
	var logoUrl string
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	bank, err := b.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	updatedBank := *bank

	if bank_request.BIC != "" {
		updatedBank.BIC = bank_request.BIC
	}
	if bank_request.Code != "" {
		updatedBank.Code = bank_request.Code
	}
	if bank_request.Name != "" {
		updatedBank.Name = bank_request.Name
	}

	if bank_request.Logo != nil {
		var objectkey string
		if bank.Logo != "" {
			objectkey = path.Base(bank.Logo)
		}

		URL, err := lib.UploadFileToMinio(
			ctx,
			b.minio,
			b.bucketName,
			bank_request.Logo,
			b.bucketName,
			*b.cfg,
			objectkey,
			b.logger,
		)
		if err != nil {
			b.logger.Errorf("UploadFileToMinio failed", "error", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}

		logoUrl = URL
	}

	result, err := b.repo.FindByNameOrBICOrCode(ctx, bank_request.BIC, bank_request.Code, bank_request.Name)

	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code {
			return err
		}
	}

	if result != nil {
		if bank_request.BIC != "" && strings.EqualFold(result.BIC, bank_request.BIC) {
			return fmt.Errorf("%s", localization.ErrorBankWithBICAlreadyExists.Code)
		}
		if bank_request.Code != "" && strings.EqualFold(result.Code, bank_request.Code) {
			return fmt.Errorf("%s", localization.ErrorBankWithCodeAlreadyExists.Code)
		}
		if bank_request.Name != "" && strings.EqualFold(result.Name, bank_request.Name) {
			return fmt.Errorf("%s", localization.ErrorBankWithNameAlreadyExists.Code)
		}
	}
	updatedBank.Logo = logoUrl
	action := lib.CpsModelBuilder(id, makerData, bank, updatedBank, string(constants.RequestUpdateBank), constants.UPDATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}
