package ussd_merchant_service

import (
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	ussd_merchant_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	localization "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service/ussd_merchant/core"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ussdMerchantService struct {
	repo        storage.UssdMerchantRepository
	serviceRepo storage.ServicesRepository
	cfg         config.VaultConfig
	cpsService  service.CPSActionService
	bucketName  string
	minio       *s3.Client
	logger      utils.Logger
}

func NewUssdMerchantService(repo storage.UssdMerchantRepository, serviceRepo storage.ServicesRepository, minio *s3.Client, cpsService service.CPSActionService, bucketName string, cfg config.VaultConfig, logger utils.Logger) service.UssdMerchantService {
	return &ussdMerchantService{
		repo:        repo,
		cfg:         cfg,
		serviceRepo: serviceRepo,
		cpsService:  cpsService,
		bucketName:  bucketName,
		minio:       minio,
		logger:      logger,
	}
}

func (s *ussdMerchantService) CreateUssdMerchant(ctx context.Context, req ussd_merchant_dto.CreateUssdMerchantRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateUssdMerchant", "UssdMerchant", "Create")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[UssdMerchSvc][Create] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByOr(ctx, req.PhoneNumber, req.Email, req.AccountNumber)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		log.Errorf("[UssdMerchSvc][Create] exist check err: %v", err)
		return err
	}

	if req.Service != "" {
		log.Infof("[UssdMerchSvc][Create] validating service")
		_, err := s.serviceRepo.FindByID(ctx, req.Service)
		if err != nil {
			log.Errorf("[UssdMerchSvc][Create] invalid service: %v", err)
			return errors.New(localization.ErrorServiceNotFound.Code)
		}
		log.Infof("[UssdMerchSvc][Create] service valid")
	}

	if err := core.ExistingIdentifier(&existing, req); err != nil {
		log.Infof("[UssdMerchSvc][Create] data exists: %v", err)
		return err
	}

	URL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Logo, string(constants.UssdMerchantFolderName), s.cfg, "", s.logger)
	if err != nil {
		log.Errorf("[UssdMerchSvc][Create] upload logo err: %v", err)
		return errors.New(localization.ErrorFileUploadFailed.Code)
	}

	ussdMerchant := core.UssdMerchant(req)

	ussdMerchant.Logo = URL
	core.CreateCredentials(&ussdMerchant, s.cfg)

	cspActionModel := lib.CpsModelBuilder(constants.Empty, makerData, nil, ussdMerchant, constants.RequestCreateUssdMerchant, constants.CREATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cspActionModel); err != nil {
		log.Errorf("[UssdMerchSvc][Create] cps action err: %v", err)
		return err
	}

	log.Infof("[UssdMerchSvc][Create] done")
	return nil
}

func (s *ussdMerchantService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "UssdMerchant", "List")
	defer span.End()

	data, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		log.Errorf("[UssdMerchSvc][FindAll] err: %v", err)
		return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{}, err
	}

	return data, nil
}

func (s *ussdMerchantService) GetUssdMerchantByID(ctx context.Context, id string) (ussd_merchant_dto.UssdMerchantResponse, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetUssdMerchantByID", "UssdMerchant", "Get")
	defer span.End()

	data, err := s.repo.FindById(ctx, id)
	if err != nil {
		log.Errorf("[UssdMerchSvc][GetByID] err: %v", err)
		return ussd_merchant_dto.UssdMerchantResponse{}, err
	}

	service, err := s.serviceRepo.FindByID(ctx, data.Service)
	if err != nil {
		s.logger.Errorf("[GetUssdMerchantByID][service] error while getting service for uud merchant err :%v", err)
		return ussd_merchant_dto.UssdMerchantResponse{}, localization.ErrorUnexpectedError
	}

	data.ServiceName = service.ServiceName
	return data, nil
}

func (s *ussdMerchantService) EnableUssdMerchant(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "EnableUssdMerchant", "UssdMerchant", "Enable")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[UssdMerchSvc][Enable] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevData, err := s.repo.FindById(ctx, id)
	if err != nil {
		log.Errorf("[UssdMerchSvc][Enable] find err: %v", err)
		return err
	}

	if prevData.Enabled {
		log.Errorf("[UssdMerchSvc][Enable] already enabled")
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	currentData := prevData

	currentData.Enabled = true

	cpsActionModel := lib.CpsModelBuilder(id, makerData, prevData, currentData, constants.RequestEnableUssdMerchant, constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		log.Errorf("[UssdMerchSvc][Enable] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *ussdMerchantService) DisableUssdMerchant(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "EnableUssdMerchant", "UssdMerchant", "Enable")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[UssdMerchSvc][Disable] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevData, err := s.repo.FindById(ctx, id)
	if err != nil {
		log.Errorf("[UssdMerchSvc][Disable] find err: %v", err)
		return err
	}

	if !prevData.Enabled {
		log.Errorf("[UssdMerchSvc][Disable] already disabled")
		return errors.New(localization.ErrorAlreadyDisabled.Code)

	}
	currentData := prevData

	currentData.Enabled = false
	cpsActionModel := lib.CpsModelBuilder(id, makerData, prevData, currentData, constants.RequestDisableUssdMerchant, constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		log.Errorf("[UssdMerchSvc][Disable] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *ussdMerchantService) UpdateUssdMerchant(ctx context.Context, id string, req ussd_merchant_dto.UpdateUssdMerchantRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateUssdMerchantService", "UssdMerchant", "Update")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[UssdMerchSvc][Update] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevMerchant, err := s.repo.FindById(ctx, id)
	if err != nil {
		log.Errorf("[UssdMerchSvc][Update] find err: %v", err)
		return err
	}

	if err := s.checkIdentifierConflict(ctx, id, req); err != nil {
		log.Infof("[UssdMerchSvc][Update] identifier conflict: %v", err)
		return err
	}

	if err := s.validateServiceChange(ctx, req.Service, prevMerchant.Service); err != nil {
		log.Errorf("[UssdMerchSvc][Update] service validation err: %v", err)
		return err
	}

	logoURL, err := s.uploadLogoIfPresent(ctx, req)
	if err != nil {
		log.Errorf("[UssdMerchSvc][Update] upload logo err: %v", err)
		return err
	}

	ussdMerchant := core.UssdMerchantUpdate(req, &prevMerchant)
	if logoURL != "" {
		ussdMerchant.Logo = logoURL
	}

	cspActionModel := lib.CpsModelBuilder(id, makerData, prevMerchant, ussdMerchant, constants.RequestUpdateUssdMerchant, constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cspActionModel); err != nil {
		log.Errorf("[UssdMerchSvc][Update] cps action err: %v", err)
		return err
	}

	log.Infof("[UssdMerchSvc][Update] done")
	return nil
}

func (s *ussdMerchantService) checkIdentifierConflict(ctx context.Context, id string, req ussd_merchant_dto.UpdateUssdMerchantRequest) error {
	if req.PhoneNumber == "" && req.Email == "" && req.AccountNumber == "" {
		return nil
	}
	existing, err := s.repo.FindByOr(ctx, req.PhoneNumber, req.Email, req.AccountNumber)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return nil
		}
		return err
	}
	return core.ExistingIdentifierForUpdate(existing, id, req)
}

func (s *ussdMerchantService) validateServiceChange(ctx context.Context, newService, currentService string) error {
	if newService == "" || strings.EqualFold(newService, currentService) {
		return nil
	}
	_, err := s.serviceRepo.FindByID(ctx, newService)
	if err != nil {
		return errors.New(localization.ErrorServiceNotFound.Code)
	}
	return nil
}

func (s *ussdMerchantService) uploadLogoIfPresent(ctx context.Context, req ussd_merchant_dto.UpdateUssdMerchantRequest) (string, error) {
	if req.Logo == nil {
		return "", nil
	}
	url, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Logo, string(constants.BankFolderName), s.cfg, "", s.logger)
	if err != nil {
		return "", errors.New(localization.ErrorUnhandledServer.Code)
	}
	return url, nil
}

func (s *ussdMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "AuthorizeUssdMerchant", "UssdMerchant", "Authorize")
	defer span.End()
	var curMerchant *imodel.UssdMerchant
	var err error

	action := cpsAction.RequestAction
	curMerchant, err = local_util.JsonUnmarshal[imodel.UssdMerchant](cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[UssdMerchSvc][Authorize] unmarshal err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch action {
	case string(constants.RequestCreateUssdMerchant):
		if err := s.repo.Create(ctx, *curMerchant); err != nil {
			log.Errorf("[UssdMerchSvc][Authorize] create err: %v", err)
			return nil, errors.New(localization.ErrorUnhandledServer.Code)
		}
	case string(constants.RequestUpdateUssdMerchant):
		if err := s.repo.Update(ctx, cpsAction.UniqueId, curMerchant); err != nil {
			log.Errorf("[UssdMerchSvc][Authorize] update err: %v", err)
			return nil, errors.New(localization.ErrorUnhandledServer.Code)
		}
	case string(constants.RequestEnableUssdMerchant):
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			log.Errorf("[UssdMerchSvc][Authorize] enable err: %v", err)
			return nil, err
		}
	case string(constants.RequestDisableUssdMerchant):
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			log.Errorf("[UssdMerchSvc][Authorize] disable err: %v", err)
			return nil, err
		}
	case constants.RequestDeleteUssdMerchant:
		if err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			log.Errorf("[UssdMerchSvc][Authorize] delete err: %v", err)
			return nil, err
		}
	default:
		log.Errorf("[UssdMerchSvc][Authorize] invalid: %s", action)
		return nil, errors.New(localization.ErrorInvalidAction.Code)
	}

	return cpsAction, nil
}

func (s *ussdMerchantService) DeleteUssdMerchant(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteUssdMerchant", "UssdMerchantService", "UssdMerhantService")
	defer span.End()
	log.Infof("[ussdMercantSvc][Delete] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[UssdMerchSvc][Update] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevData, err := s.repo.FindById(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorUssdMerchantNotFound.Code)
	}

	now := time.Now()
	deletedUssdMerchant := prevData
	deletedUssdMerchant.IsDeleted = true
	deletedUssdMerchant.DeletedAt = now

	cspActionModel := lib.CpsModelBuilder(id, makerData, prevData, deletedUssdMerchant, constants.RequestDeleteUssdMerchant, constants.DELETE)
	if err := s.cpsService.CreateCPSAction(ctx, &cspActionModel); err != nil {
		log.Errorf("[UssdMerchSvc][Delete] cps action err: %v", err)
		return err
	}

	log.Infof("[UssdMerchSvc][Delete] done")
	return nil
}
