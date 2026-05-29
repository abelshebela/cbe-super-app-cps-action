package bankService

import (
	"cbe-super-app-cps-action/internal/constants"
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type BankService struct {
	cpsService  service.CPSActionService
	logger      utils.Logger
	repo        storage.BankRepository
	oracleRepo  storage.BankOracleRepository
	cfg         *config.VaultConfig
	minio       *s3.Client
	minioPubUrl string
	bucketName  string
}

func NewBankService(logger utils.Logger, repo storage.BankRepository, oracleRepo storage.BankOracleRepository, cpsService service.CPSActionService, minio *s3.Client, minioPubUrl string, cfg *config.VaultConfig, bucketName string) service.BankService {
	return &BankService{
		logger:      logger,
		repo:        repo,
		oracleRepo:  oracleRepo,
		cpsService:  cpsService,
		cfg:         cfg,
		minio:       minio,
		minioPubUrl: minioPubUrl,
		bucketName:  bucketName,
	}
}

func (b *BankService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Bank", "Authorize")
	defer span.End()

	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("failed to marshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		log.Errorf("[BankSvc][Authorize] marshal err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		span.AddEvent("failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		log.Errorf("[BankSvc][Authorize] unmarshal err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// actionData := bank_core.Bank_mapper(actionMap.(map[string]interface{}))
	// if cpsAction.UniqueId != "" {
	// 	objID, err := bson.ObjectIDFromHex(cpsAction.UniqueId)
	// 	if err != nil {
	// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// 	}
	// 	actionData.ID = objID
	// }
	actionData := bank_core.Bank_oracle_mapper(actionMap.(map[string]interface{}))

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateBank):
		actionData.CreateAt = time.Now().String()
		err := b.oracleRepo.Create(ctx, &actionData)
		if err != nil {
			span.AddEvent("[Authorize] bank create action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			log.Errorf("[BankSvc][Authorize] create err: %v", err)
			return nil, err
		}
		log.Infof("[BankSvc][Authorize] created")
	case string(constants.RequestDeleteBank):
		err := b.oracleRepo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] bank delete action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			log.Errorf("[BankSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		log.Infof("[BankSvc][Authorize] deleted id: %s", cpsAction.UniqueId)

	case string(constants.RequestDisableBank), string(constants.RequestEnableBank):
		actionData.UpdateAt = time.Now().String()
		var err error
		if actionData.IsEnabled == 1 {
			err = b.oracleRepo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		} else {
			err = b.oracleRepo.EnableOrDisable(ctx, cpsAction.UniqueId, false)

		}
		if err != nil {
			span.AddEvent("[Authorize] bank enable/disable action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			log.Errorf("[BankSvc][Authorize] enable/disable err: %v", err)
			return nil, err
		}
		log.Infof("[BankSvc][Authorize] enable/disable done id: %s", cpsAction.UniqueId)
	case string(constants.RequestUpdateBankLogo):
		err := b.oracleRepo.Update(ctx, cpsAction.UniqueId, &actionData)
		if err != nil {
			span.AddEvent("[Authorize] bank update logo action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			log.Errorf("[BankSvc][Authorize] update logo err: %v", err)
			return nil, err
		}
		log.Infof("[BankSvc][Authorize] logo updated id: %s", cpsAction.UniqueId)
	case string(constants.RequestUpdateBank):
		actionData.UpdateAt = time.Now().String()
		err := b.oracleRepo.Update(ctx, cpsAction.UniqueId, &actionData)
		if err != nil {
			span.AddEvent("[Authorize] bank update action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			log.Errorf("[BankSvc][Authorize] update err: %v", err)
			return nil, err
		}
		log.Infof("[BankSvc][Authorize] updated id: %s", cpsAction.UniqueId)
	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
		log.Errorf("[BankSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		return nil, fmt.Errorf("%s", localization.MsgBankInvalidRequestAction)
	}
	log.Infof("[BankSvc][Authorize] done: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

func (b *BankService) CreateOneBank(ctx context.Context, bank_request bank_dto.CreateBankRequest) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateOneBank", "Bank", "CreateOneBank")
	defer span.End()

	log.Infof("[BankSvc][CreateOneBank] creating")
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[CreateOneBank] incomplete user data")
		log.Errorf("[BankSvc][CreateOneBank] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	URL, err := lib.UploadFileToMinio(ctx, b.minio, b.bucketName, bank_request.Logo, string(constants.BankFolderName), *b.cfg, "", b.logger)
	if err != nil {
		span.AddEvent("[CreateOneBank] failed to upload logo", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("bank_name", bank_request.Name),
		))
		log.Errorf("[BankSvc][CreateOneBank] upload logo err: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	bank := imodel.BankOracle{
		BankName:      bank_request.Name,
		Logo:          URL,
		BICCode:       bank_request.BICCode,
		IsEnabled:     0,
		AccountLength: bank_request.AccountLength,
	}
	if bank_request.HasAlphaNumeric != nil {
		if *bank_request.HasAlphaNumeric {
			bank.HasAlphaNumeric = 1
		} else {
			bank.HasAlphaNumeric = 0
		}
	}

	if bank_request.IsCBE != nil {
		if *bank_request.IsCBE {
			bank.IS_CBE = 1
		} else {
			bank.IS_CBE = 0
		}
	}

	result, err := b.oracleRepo.FindByNameOrBIC(ctx, bank_request.BICCode, "")
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code {
			log.Errorf("[BankSvc][CreateOneBank] dup check err: %v", err)
			return err
		}
	}
	if result != nil && result.ID != "" {
		log.Infof("[BankSvc][CreateOneBank] found existing bank name: %s, bic: %s, id: %s", result.BankName, result.BICCode, result.ID)
		if bank_request.BICCode != "" && result.BICCode != "" && strings.EqualFold(result.BICCode, bank_request.BICCode) {
			log.Errorf("[BankSvc][CreateOneBank] BIC exists")
			return fmt.Errorf("%s", localization.ErrorBankWithBICAlreadyExists.Code)
		}
	}

	result, err = b.oracleRepo.FindByNameOrBIC(ctx, "", bank_request.Name)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code {
			log.Errorf("[BankSvc][CreateOneBank] dup check err: %v", err)
			return err
		}
	}
	if result != nil && result.ID != "" {
		log.Infof("[BankSvc][CreateOneBank] found existing bank name: %s, bic: %s, id: %s", result.BankName, result.BICCode, result.ID)

		if bank_request.Name != "" && result.BankName != "" && strings.EqualFold(result.BankName, bank_request.Name) {
			log.Errorf("[BankSvc][CreateOneBank] name exists")
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
		log.Errorf("[BankSvc][CreateOneBank] cps action err: %v", err)
		return err
	}
	log.Infof("[BankSvc][CreateOneBank] request created")
	return nil
}

func (b *BankService) DeleteOneBank(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteOneBank", "Bank", "DeleteOneBank")
	defer span.End()

	log.Infof("[BankSvc][DeleteOneBank] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[DeleteOneBank] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[BankSvc][DeleteOneBank] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	bank, err := b.oracleRepo.FindByID(ctx, id)

	if err != nil {
		span.AddEvent("[DeleteOneBank] failed to find bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		log.Errorf("[BankSvc][DeleteOneBank] find err: %v", err)
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
		log.Errorf("[BankSvc][DeleteOneBank] cps action err: %v", err)
		return err
	}
	log.Infof("[BankSvc][DeleteOneBank] request created id: %s", id)
	return nil
}

func (b *BankService) EnableOrDisableBank(ctx context.Context, id string, enableDisable bool) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableBank", "Bank", "EnableOrDisableBank")
	defer span.End()

	log.Infof("[BankSvc][EnableOrDisable] id: %s enabled: %v", id, enableDisable)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[EnableOrDisableBank] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[BankSvc][EnableOrDisable] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	bank, err := b.oracleRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[EnableOrDisableBank] failed to find bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		log.Errorf("[BankSvc][EnableOrDisable] find err: %v", err)
		return err
	}

	if bank.IsEnabled == 1 && enableDisable {
		span.AddEvent("[EnableOrDisableBank] bank already enabled", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[BankSvc][EnableOrDisable] already enabled")
		return fmt.Errorf("%s", localization.ErrorBankAlreadyEnabled.Code)
	} else if bank.IsEnabled == 0 && !enableDisable {
		span.AddEvent("[EnableOrDisableBank] bank already disabled", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[BankSvc][EnableOrDisable] already disabled")
		return fmt.Errorf("%s", localization.ErrorBankAlreadyDisabled.Code)
	}

	newBankData := *bank
	// newBankData.Enabled = enableDisable

	var enable string
	if enableDisable {
		newBankData.IsEnabled = 1
		enable = string(constants.RequestEnableBank)
	} else {
		newBankData.IsEnabled = 0
		enable = string(constants.RequestDisableBank)
	}

	// var enable string
	// if enableDisable {
	// 	enable = string(constants.RequestEnableBank)
	// } else {
	// 	enable = string(constants.RequestDisableBank)
	// }
	action := lib.CpsModelBuilder(id, makerData, bank, newBankData, enable, constants.UPDATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		span.AddEvent("[EnableOrDisableBank] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		log.Errorf("[BankSvc][EnableOrDisable] cps action err: %v", err)
		return err
	}
	log.Infof("[BankSvc][EnableOrDisable] request created id: %s", id)
	return nil
}

func (b *BankService) GetAllBank(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.BankOracle], error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllBank", "Bank", "GetAllBank")
	defer span.End()

	result, err := b.oracleRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("[GetAllBank] failed to fetch banks", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		log.Errorf("[BankSvc][GetAllBank] fetch err: %v", err)
		return nil, err
	}
	log.Infof("[BankSvc][GetAllBank] count: %d", len(result.Data))
	return result, nil
}

func (b *BankService) GetOneBank(ctx context.Context, id string) (*imodel.BankOracle, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetOneBank", "Bank", "GetOneBank")
	defer span.End()

	result, err := b.oracleRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[GetOneBank] failed to fetch bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		log.Errorf("[BankSvc][GetOneBank] fetch err: %v", err)
		return nil, err
	}
	log.Infof("[BankSvc][GetOneBank] found id: %s", id)
	return result, nil
}

// GetOneBankByBIC implements [service.BankService].
func (b *BankService) GetOneBankByBIC(ctx context.Context, bicCode string) (*imodel.BankOracle, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetOneBankByBIC", "Bank", "GetOneBankByBIC")
	defer span.End()
	log.Infof("[BankSvc][GetOneBankByBIC] bic: %s", bicCode)
	result, err := b.oracleRepo.FindByBIC(ctx, bicCode)
	if err != nil {
		span.AddEvent("[GetOneBankByBIC] failed to fetch bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("bic_code", bicCode),
		))
		log.Errorf("[BankSvc][GetOneBankByBIC] fetch err: %v", err)
		return nil, err
	}
	return result, nil
}
func (b *BankService) UpdateLogo(ctx context.Context, id string, logo bank_dto.UpdateLogo) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateLogo", "Bank", "UpdateLogo")
	defer span.End()

	log.Infof("[BankSvc][UpdateLogo] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[UpdateLogo] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[BankSvc][UpdateLogo] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	bank, err := b.oracleRepo.FindByID(ctx, id)

	if err != nil {
		span.AddEvent("[UpdateLogo] failed to find bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		log.Errorf("[BankSvc][UpdateLogo] find err: %v", err)
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
		log.Errorf("[BankSvc][UpdateLogo] upload err: %v", err)
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
		log.Errorf("[BankSvc][UpdateLogo] cps action err: %v", err)
		return err
	}
	log.Infof("[BankSvc][UpdateLogo] request created id: %s", id)
	return nil
}

func (b *BankService) UpdateOneBank(ctx context.Context, id string, bank_request bank_dto.UpdateBankRequest) error {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateOneBank", "Bank", "UpdateOneBank")
	defer span.End()

	log.Infof("[BankSvc][UpdateOneBank] id: %s, body: %+v", id, bank_request)

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("[UpdateOneBank] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[BankSvc][UpdateOneBank] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	bank, err := b.oracleRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[UpdateOneBank] failed to find bank", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		log.Errorf("[BankSvc][UpdateOneBank] find err: %v", err)
		return err
	}

	updatedBank := *bank

	if bank_request.BICCode != "" {
		updatedBank.BICCode = bank_request.BICCode
	}
	if bank_request.Name != "" {
		updatedBank.BankName = bank_request.Name
	}
	if bank_request.AccountLength != 0 {
		updatedBank.AccountLength = bank_request.AccountLength
	}
	if bank_request.HasAlphaNumeric != nil {
		if *bank_request.HasAlphaNumeric {
			updatedBank.HasAlphaNumeric = 1
		} else {
			updatedBank.HasAlphaNumeric = 0
		}
	}
	if bank_request.IsCBE != nil {
		log.Infof("[BankSvc][UpdateOneBank] IsCBE: %v", *bank_request.IsCBE)
		if *bank_request.IsCBE {
			updatedBank.IS_CBE = 1
			log.Infof("[BankSvc][UpdateOneBank] Set IS_CBE to 1 for id: %s", id)
		} else {
			updatedBank.IS_CBE = 0
			log.Infof("[BankSvc][UpdateOneBank] Set IS_CBE to 0 for id: %s", id)
		}
	} else {
		log.Infof("[BankSvc][UpdateOneBank] IsCBE is nil, skipping update for id: %s", id)
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
			log.Errorf("[BankSvc][UpdateOneBank] upload err: %v", err)
			return errors.New(localization.ErrorFileUploadFailed.Code)
		}

		updatedBank.Logo = URL
	}

	result, err := b.oracleRepo.FindByNameOrBIC(ctx, bank_request.BICCode, "")
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code {
			log.Errorf("[BankSvc][UpdateOneBank] dup check err: %v", err)
			return err
		}
	}
	if result != nil && result.ID != "" {
		log.Infof("[BankSvc][UpdateOneBank] found existing bank name: %s, bic: %s, id: %s", result.BankName, result.BICCode, result.ID)
		if bank_request.BICCode != "" && result.BICCode != "" && strings.EqualFold(result.BICCode, bank_request.BICCode) && result.ID != id {
			log.Errorf("[BankSvc][UpdateOneBank] BIC exists")
			return fmt.Errorf("%s", localization.ErrorBankWithBICAlreadyExists.Code)
		}
	}
	result, err = b.oracleRepo.FindByNameOrBIC(ctx, "", bank_request.Name)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code {
			log.Errorf("[BankSvc][UpdateOneBank] dup check err: %v", err)
			return err
		}
	}
	if result != nil && result.ID != "" {
		log.Infof("[BankSvc][UpdateOneBank] found existing bank name: %s, bic: %s, id: %s", result.BankName, result.BICCode, result.ID)

		if bank_request.Name != "" && result.BankName != "" && strings.EqualFold(result.BankName, bank_request.Name) && result.ID != id {
			log.Errorf("[BankSvc][UpdateOneBank] name exists")
			return fmt.Errorf("%s", localization.ErrorBankWithNameAlreadyExists.Code)
		}
	}

	action := lib.CpsModelBuilder(id, makerData, bank, updatedBank, string(constants.RequestUpdateBank), constants.UPDATE)

	err = b.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		span.AddEvent("[UpdateOneBank] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		log.Errorf("[BankSvc][UpdateOneBank] cps action err: %v", err)
		return err
	}
	log.Infof("[BankSvc][UpdateOneBank] request created id: %s", id)
	return nil
}
