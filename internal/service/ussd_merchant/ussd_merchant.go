package ussd_merchant_service

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/v2/bson"

	"cbe-super-app-cps-action/internal/constants"
	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/ussd_merchant/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

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
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateUssdMerchant", "UssdMerchant", "Create")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[UssdMerchSvc][Create] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByOr(ctx, req.PhoneNumber, req.Email, req.AccountNumber)
	if err != nil {
		s.logger.Errorf("[UssdMerchSvc][Create] exist check err: %v", err)
		return err
	}

	if req.Service != "" {
		s.logger.Infof("[UssdMerchSvc][Create] validating service")
		_, err := s.serviceRepo.FindByID(ctx, req.Service)
		if err != nil {
			s.logger.Errorf("[UssdMerchSvc][Create] invalid service: %v", err)
			return errors.New(localization.ErrorServiceNotFound.Code)
		}
		s.logger.Infof("[UssdMerchSvc][Create] service valid")
	}

	if err := core.ExistingIdentifier(&existing, req); err != nil {
		s.logger.Infof("[UssdMerchSvc][Create] data exists: %v", err)
		return err
	}

	URL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Logo, string(constants.UssdMerchantFolderName), s.cfg, "", s.logger)
	if err != nil {
		s.logger.Errorf("[UssdMerchSvc][Create] upload logo err: %v", err)
		return errors.New(localization.ErrorFileUploadFailed.Code)
	}

	ussdMerchant := core.UssdMerchant(req)

	ussdMerchant.Logo = URL
	core.CreateCredentials(&ussdMerchant, s.cfg)

	cspActionModel := lib.CpsModelBuilder(constants.Empty, makerData, nil, ussdMerchant, constants.RequestCreateUssdMerchant, constants.CREATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cspActionModel); err != nil {
		s.logger.Errorf("[UssdMerchSvc][Create] cps action err: %v", err)
		return err
	}

	s.logger.Infof("[UssdMerchSvc][Create] done")
	return nil
}

func (s *ussdMerchantService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "UssdMerchant", "List")
	defer span.End()

	data, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		s.logger.Errorf("[UssdMerchSvc][FindAll] err: %v", err)
		return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{}, err
	}

	return data, nil
}

func (s *ussdMerchantService) GetUssdMerchantByID(ctx context.Context, id string) (ussd_merchant_dto.UssdMerchantResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUssdMerchantByID", "UssdMerchant", "Get")
	defer span.End()

	data, err := s.repo.FindById(ctx, id)
	if err != nil {
		s.logger.Errorf("[UssdMerchSvc][GetByID] err: %v", err)
		return ussd_merchant_dto.UssdMerchantResponse{}, err
	}

	return data, nil
}

func (s *ussdMerchantService) EnableUssdMerchant(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableUssdMerchant", "UssdMerchant", "Enable")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[UssdMerchSvc][Enable] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevData, err := s.repo.FindById(ctx, id)
	if err != nil {
		s.logger.Errorf("[UssdMerchSvc][Enable] find err: %v", err)
		return err
	}

	if prevData.Enabled {
		s.logger.Errorf("[UssdMerchSvc][Enable] already enabled")
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	currentData := prevData

	currentData.Enabled = true

	cpsActionModel := lib.CpsModelBuilder(id, makerData, prevData, currentData, constants.RequestEnableUssdMerchant, constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[UssdMerchSvc][Enable] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *ussdMerchantService) DisableUssdMerchant(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableUssdMerchant", "UssdMerchant", "Enable")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[UssdMerchSvc][Disable] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevData, err := s.repo.FindById(ctx, id)
	if err != nil {
		s.logger.Errorf("[UssdMerchSvc][Disable] find err: %v", err)
		return err
	}

	if !prevData.Enabled {
		s.logger.Errorf("[UssdMerchSvc][Disable] already disabled")
		return errors.New(localization.ErrorAlreadyDisabled.Code)

	}
	currentData := prevData

	currentData.Enabled = false
	cpsActionModel := lib.CpsModelBuilder(id, makerData, prevData, currentData, constants.RequestDisableUssdMerchant, constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[UssdMerchSvc][Disable] cps action err: %v", err)
		return err
	}

	return nil
}
func (s *ussdMerchantService) UpdateUssdMerchant(ctx context.Context, id string, req ussd_merchant_dto.UpdateUssdMerchantRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateUssdMerchantService", "UssdMerchant", "Update")
	defer span.End()
	var URL string
	var existing imodel.UssdMerchant

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[UssdMerchSvc][Update] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevMerchant, err := s.repo.FindById(ctx, id)
	if err != nil {
		s.logger.Errorf("[UssdMerchSvc][Update] find err: %v", err)
		return err
	}

	if req.PhoneNumber != "" || req.Email != "" || req.AccountNumber != "" {
		existing, err := s.repo.FindByOr(ctx, req.PhoneNumber, req.Email, req.AccountNumber)
		if err != nil {
			s.logger.Errorf("[UssdMerchSvc][Update] exist check err: %v", err)
			return err
		}
		if err := core.ExistingIdentifierForUpdate(existing, id, req); err != nil {
			s.logger.Infof("[UssdMerchSvc][Update] data exists: %v", err)
			return err
		}
	}

	if req.Service != "" && !strings.EqualFold(req.Service, prevMerchant.Service) {
		s.logger.Infof("[UssdMerchSvc][Update] validating service")
		_, err := s.serviceRepo.FindByID(ctx, req.Service)
		if err != nil {
			s.logger.Errorf("[UssdMerchSvc][Update] invalid service: %v", err)
			return errors.New(localization.ErrorServiceNotFound.Code)
		}
		s.logger.Infof("[UssdMerchSvc][Update] service valid")
	}

	if err := core.ExistingIdentifierForUpdate(existing, id, req); err != nil {
		s.logger.Infof("[UssdMerchSvc][Update] data exists: %v", err)
		return err
	}

	if req.Logo != nil {
		URL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Logo, string(constants.BankFolderName), s.cfg, "", s.logger)
		if err != nil {
			s.logger.Errorf("[UssdMerchSvc][Update] upload logo err: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	}

	prevData := prevMerchant
	currData := prevMerchant

	ussdMerchant := core.UssdMerchantUpdate(req, &currData)
	if URL != "" {
		ussdMerchant.Logo = URL
	}

	cspActionModel := lib.CpsModelBuilder(id, makerData, prevData, ussdMerchant, constants.RequestUpdateUssdMerchant, constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cspActionModel); err != nil {
		s.logger.Errorf("[UssdMerchSvc][Update] cps action err: %v", err)
		return err
	}

	s.logger.Infof("[UssdMerchSvc][Update] done")
	return nil
}

func (s *ussdMerchantService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "AuthorizeUssdMerchant", "UssdMerchant", "Authorize")
	defer span.End()
	var curMerchant *imodel.UssdMerchant
	var err error

	action := cpsAction.RequestAction
	curMerchant, err = local_util.JsonUnmarshal[imodel.UssdMerchant](cpsAction.CurrentAction)
	if err != nil {
		s.logger.Errorf("[UssdMerchSvc][Authorize] unmarshal err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch action {
	case string(constants.RequestCreateUssdMerchant):
		if err := s.repo.Create(ctx, *curMerchant); err != nil {
			s.logger.Errorf("[UssdMerchSvc][Authorize] create err: %v", err)
			return nil, errors.New(localization.ErrorUnhandledServer.Code)
		}
	case string(constants.RequestUpdateUssdMerchant):
		update := core.ModelToBson(curMerchant)
		if err := s.repo.Update(ctx, cpsAction.UniqueId, update); err != nil {
			s.logger.Errorf("[UssdMerchSvc][Authorize] update err: %v", err)
			return nil, errors.New(localization.ErrorUnhandledServer.Code)
		}
	case string(constants.RequestEnableUssdMerchant):
		if err := s.repo.Update(ctx, cpsAction.UniqueId, bson.M{"enabled": true}); err != nil {
			s.logger.Errorf("[UssdMerchSvc][Authorize] enable err: %v", err)
			return nil, err
		}
	case string(constants.RequestDisableUssdMerchant):
		if err := s.repo.Update(ctx, cpsAction.UniqueId, bson.M{"enabled": false}); err != nil {
			s.logger.Errorf("[UssdMerchSvc][Authorize] disable err: %v", err)
			return nil, err
		}
	default:
		s.logger.Errorf("[UssdMerchSvc][Authorize] invalid: %s", action)
		return nil, errors.New(localization.ErrorInvalidAction.Code)
	}

	return cpsAction, nil
}
