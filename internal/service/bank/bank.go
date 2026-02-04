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
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

// GetOneBankByBIC implements [service.BankService].
func (b *BankService) GetOneBankByBIC(ctx context.Context, bicCode string) (*model.Bank, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetOneBankByBIC", "Bank", "GetOneBankByBIC")
	defer span.End()
	b.logger.Infof("[GetOneBankByBIC] fetching bank for BIC: %s", bicCode)
	result, err := b.repo.FindByBIC(ctx, bicCode)
	if err != nil {
		span.AddEvent("[GetOneBankByBIC] failed to fetch bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("bic_code", bicCode),
		))
		b.logger.Errorf("[GetOneBankByBIC] failed to fetch bank: %v", err)
		return nil, err
	}
	return result, nil
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Bank", "Authorize")
	defer span.End()

	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("failed to marshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		b.logger.Errorf("failed to marshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		span.AddEvent("failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		b.logger.Errorf("failed to unmarshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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
			span.AddEvent("[Authorize] bank create action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[Authorize] bank create action failed: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] bank created successfully")
	case string(constants.RequestDeleteBank):
		err := b.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] bank delete action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[Authorize] bank delete action failed: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] bank deleted successfully for id: %s", cpsAction.UniqueId)

	case string(constants.RequestEnableDisableBank):
		err := b.repo.EnableOrDisable(ctx, cpsAction.UniqueId, actionData.Enabled)

		if err != nil {
			span.AddEvent("[Authorize] bank enable/disable action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[Authorize] bank enable/disable action failed: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] bank enable/disable action completed successfully for id: %s", cpsAction.UniqueId)
	case string(constants.RequestUpdateBankLogo):
		err := b.repo.Update(ctx, cpsAction.UniqueId, &actionData)
		if err != nil {
			span.AddEvent("[Authorize] bank update logo action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[Authorize] bank update logo action failed: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] bank logo updated successfully for id: %s", cpsAction.UniqueId)
	case string(constants.RequestUpdateBank):
		err := b.repo.Update(ctx, cpsAction.UniqueId, &actionData)
		if err != nil {
			span.AddEvent("[Authorize] bank update action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[Authorize] bank update action failed: %v", err)
			return nil, err
		}
		b.logger.Infof("[Authorize] bank updated successfully for id: %s", cpsAction.UniqueId)
	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
		b.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, fmt.Errorf("%s", localization.MsgBankInvalidRequestAction)
	}
	b.logger.Infof("[Authorize] bank action authorized successfully: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

func (b *BankService) CreateOneBank(ctx context.Context, bank_request bank_dto.CreateBankRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateOneBank", "Bank", "CreateOneBank")
	defer span.End()

	b.logger.Infof("[CreateOneBank] creating bank")
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[CreateOneBank] incomplete user data")
		b.logger.Errorf("[CreateOneBank] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	URL, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, bank_request.Logo, string(constants.BankFolderName), *b.cfg, "", b.logger)
	if err != nil {
		span.AddEvent("[CreateOneBank] failed to upload logo", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("bank_name", bank_request.Name),
		))
		b.logger.Errorf("[CreateOneBank] failed to upload logo: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	bank := model.Bank{
		Name:    bank_request.Name,
		BICCode: bank_request.BICCode,
		Type:    string(bank_request.Type),
		Logo:    URL,
		Enabled: false,
	}

	result, err := b.repo.FindByNameOrBIC(ctx, bank_request.BICCode, bank_request.Name)

	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code {
			b.logger.Errorf("[CreateOneBank] failed to check for duplicate bank: %v", err)
			return err
		}
	}

	if result != nil {
		if bank_request.BICCode != "" && result.BICCode != "" && result.BICCode == bank_request.BICCode {
			b.logger.Errorf("[CreateOneBank] bank with BIC already exists")
			return fmt.Errorf("%s", localization.ErrorBankWithBICAlreadyExists.Code)
		}
		if bank_request.Name != "" && result.Name != "" && strings.ToLower(result.Name) == strings.ToLower(bank_request.Name) {
			b.logger.Errorf("[CreateOneBank] bank with name already exists")
			return fmt.Errorf("%s", localization.ErrorBankWithNameAlreadyExists.Code)
		}
	}

	action := lib.CpsModelBuilder("", makerData, nil, bank, string(constants.RequestCreateBank), constants.CREATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		span.AddEvent("[CreateOneBank] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("bank_name", bank_request.Name),
		))
		b.logger.Errorf("[CreateOneBank] failed to create CPS action: %v", err)
		return err
	}
	b.logger.Infof("[CreateOneBank] bank creation request created successfully")
	return nil
}

func (b *BankService) DeleteOneBank(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteOneBank", "Bank", "DeleteOneBank")
	defer span.End()

	b.logger.Infof("[DeleteOneBank] deleting bank for id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[DeleteOneBank] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		b.logger.Errorf("[DeleteOneBank] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	bank, err := b.repo.FindByID(ctx, id)

	if err != nil {
		span.AddEvent("[DeleteOneBank] failed to find bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[DeleteOneBank] failed to find bank: %v", err)
		return err
	}

	newBankData := *bank

	newBankData.IsDeleted = true
	action := lib.CpsModelBuilder(id, makerData, bank, newBankData, string(constants.RequestDeleteBank), constants.DELETE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		span.AddEvent("[DeleteOneBank] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[DeleteOneBank] failed to create CPS action: %v", err)
		return err
	}
	b.logger.Infof("[DeleteOneBank] bank deletion request created successfully for id: %s", id)
	return nil
}

func (b *BankService) EnableOrDisableBank(ctx context.Context, id string, enableDisable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableBank", "Bank", "EnableOrDisableBank")
	defer span.End()

	b.logger.Infof("[EnableOrDisableBank] processing bank enable/disable for id: %s, enabled: %v", id, enableDisable)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[EnableOrDisableBank] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		b.logger.Errorf("[EnableOrDisableBank] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	bank, err := b.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[EnableOrDisableBank] failed to find bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[EnableOrDisableBank] failed to find bank: %v", err)
		return err
	}

	if bank.Enabled && enableDisable {
		span.AddEvent("[EnableOrDisableBank] bank already enabled", trace.WithAttributes(attribute.String("id", id)))
		b.logger.Errorf("[EnableOrDisableBank] bank already enabled")
		return fmt.Errorf("%s", localization.ErrorBankAlreadyEnabled.Code)
	} else if !bank.Enabled && !enableDisable {
		span.AddEvent("[EnableOrDisableBank] bank already disabled", trace.WithAttributes(attribute.String("id", id)))
		b.logger.Errorf("[EnableOrDisableBank] bank already disabled")
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
		span.AddEvent("[EnableOrDisableBank] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[EnableOrDisableBank] failed to create CPS action: %v", err)
		return err
	}
	b.logger.Infof("[EnableOrDisableBank] bank enable/disable request created successfully for id: %s", id)
	return nil
}

func (b *BankService) GetAllBank(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]model.Bank], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllBank", "Bank", "GetAllBank")
	defer span.End()

	result, err := b.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("[GetAllBank] failed to fetch banks", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		b.logger.Errorf("[GetAllBank] failed to fetch banks: %v", err)
		return nil, err
	}
	b.logger.Infof("[GetAllBank] retrieved %d banks", len(result.Data))
	return result, nil
}

func (b *BankService) GetOneBank(ctx context.Context, id string) (*model.Bank, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetOneBank", "Bank", "GetOneBank")
	defer span.End()

	result, err := b.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[GetOneBank] failed to fetch bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[GetOneBank] failed to fetch bank: %v", err)
		return nil, err
	}
	b.logger.Infof("[GetOneBank] bank retrieved successfully for id: %s", id)
	return result, nil
}

func (b *BankService) UpdateLogo(ctx context.Context, id string, logo bank_dto.UpdateLogo) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateLogo", "Bank", "UpdateLogo")
	defer span.End()

	b.logger.Infof("[UpdateLogo] updating bank logo for id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[UpdateLogo] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		b.logger.Errorf("[UpdateLogo] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	bank, err := b.repo.FindByID(ctx, id)

	if err != nil {
		span.AddEvent("[UpdateLogo] failed to find bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[UpdateLogo] failed to find bank: %v", err)
		return err
	}

	var objectkey string
	if bank.Logo != "" {
		objectkey = path.Base(bank.Logo)
	}

	URL, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, logo.Logo, string(constants.BankFolderName), *b.cfg, objectkey, b.logger)

	if err != nil {
		span.AddEvent("[UpdateLogo] failed to upload logo", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[UpdateLogo] failed to upload logo: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	newBankData := *bank

	newBankData.Logo = URL

	action := lib.CpsModelBuilder(id, makerData, bank, newBankData, string(constants.RequestUpdateBankLogo), constants.UPDATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		span.AddEvent("[UpdateLogo] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[UpdateLogo] failed to create CPS action: %v", err)
		return err
	}
	b.logger.Infof("[UpdateLogo] bank logo update request created successfully for id: %s", id)
	return nil
}

func (b *BankService) UpdateOneBank(ctx context.Context, id string, bank_request bank_dto.UpdateBankRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateOneBank", "Bank", "UpdateOneBank")
	defer span.End()

	b.logger.Infof("[UpdateOneBank] updating bank for id: %s", id)
	var logoUrl string
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[UpdateOneBank] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		b.logger.Errorf("[UpdateOneBank] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	bank, err := b.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[UpdateOneBank] failed to find bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[UpdateOneBank] failed to find bank: %v", err)
		return err
	}

	updatedBank := *bank

	if bank_request.BICCode != "" {
		updatedBank.BICCode = bank_request.BICCode
	}
	if bank_request.Name != "" {
		updatedBank.Name = bank_request.Name
	}
	if bank_request.Type != "" {
		updatedBank.Type = string(bank_request.Type)
	}

	logoUrl = bank.Logo
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
			string(constants.BankFolderName),
			*b.cfg,
			objectkey,
			b.logger,
		)
		if err != nil {
			span.AddEvent("[UpdateOneBank] UploadFileToMinio failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			b.logger.Errorf("UploadFileToMinio failed", "error", err)
			return errors.New(localization.ErrorFileUploadFailed.Code)
		}

		logoUrl = URL
	}

	result, err := b.repo.FindByNameOrBIC(ctx, bank_request.BICCode, bank_request.Name)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code {
			b.logger.Errorf("[UpdateOneBank] failed to check for duplicate bank: %v", err)
			return err
		}
	}
	if result != nil && result.ID.Hex() != id {
		if bank_request.BICCode != "" && result.BICCode != "" && result.BICCode == bank_request.BICCode && result.ID.Hex() != id {
			b.logger.Errorf("[UpdateOneBank] bank with BIC already exists")
			return fmt.Errorf("%s", localization.ErrorBankWithBICAlreadyExists.Code)
		}

		if bank_request.Name != "" && result.Name != "" && result.Name == bank_request.Name && result.ID.Hex() != id {
			b.logger.Errorf("[UpdateOneBank] bank with name already exists")
			return fmt.Errorf("%s", localization.ErrorBankWithNameAlreadyExists.Code)
		}
	}

	updatedBank.Logo = logoUrl
	action := lib.CpsModelBuilder(id, makerData, bank, updatedBank, string(constants.RequestUpdateBank), constants.UPDATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		span.AddEvent("[UpdateOneBank] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		b.logger.Errorf("[UpdateOneBank] failed to create CPS action: %v", err)
		return err
	}
	b.logger.Infof("[UpdateOneBank] bank update request created successfully for id: %s", id)
	return nil
}
