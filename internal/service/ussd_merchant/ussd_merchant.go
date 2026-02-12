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
		s.logger.Errorf("[CreateUssdMerchant] Error while extracting user info from context it's incomplete ")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByOr(ctx, req.PhoneNumber, req.Email, req.AccountNumber)
	if err != nil {
		s.logger.Errorf("[CreateUssdMerchant] error whil checking existing information error: %v", err)
		return err
	}

	if req.Service != "" {
		s.logger.Infof("[CreateUssdMerchant] service is not empty, checking if it's valid")
		_, err := s.serviceRepo.FindByID(ctx, req.Service)
		if err != nil {
			s.logger.Errorf("[CreateUssdMerchant] service is not valid, error: %v", err)
			return errors.New(localization.ErrorServiceNotFound.Code)
		}
		s.logger.Infof("[CreateUssdMerchant] service is valid")
	}

	if err := core.ExistingIdentifier(&existing, req); err != nil {
		s.logger.Infof("[CreateUssdMerchant] the entered data is already existed error: %v", err)
		return err
	}

	URL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Logo, string(constants.UssdMerchantFolderName), s.cfg, "", s.logger)
	if err != nil {
		s.logger.Errorf("[CreateOneBank] failed to upload logo: %v", err)
		return errors.New(localization.ErrorFileUploadFailed.Code)
	}

	ussdMerchant := core.UssdMerchant(req)

	ussdMerchant.Logo = URL
	core.CreateCredentials(&ussdMerchant, s.cfg)

	cspActionModel := lib.CpsModelBuilder(constants.Empty, makerData, nil, ussdMerchant, constants.RequestCreateUssdMerchant, constants.CREATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cspActionModel); err != nil {
		s.logger.Errorf("[CreateUssdMerchant] failed to create CPS action: %v", err)
		return err
	}

	s.logger.Infof("[CreateUssdMerchant] CreateUssdMerchant Successful")
	return nil
}

func (s *ussdMerchantService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "UssdMerchant", "List")
	defer span.End()

	data, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] error whil listing data: %v", err)
		return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{}, err
	}

	return data, nil
}

func (s *ussdMerchantService) GetUssdMerchantByID(ctx context.Context, id string) (ussd_merchant_dto.UssdMerchantResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUssdMerchantByID", "UssdMerchant", "Get")
	defer span.End()

	data, err := s.repo.FindById(ctx, id)
	if err != nil {
		s.logger.Errorf("[GetUssdMerchantByID] error while fetching merchant details: %v", err)
		return ussd_merchant_dto.UssdMerchantResponse{}, err
	}

	return data, nil
}

func (s *ussdMerchantService) EnableUssdMerchant(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableUssdMerchant", "UssdMerchant", "Enable")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[EnabledUssdMerchant] Error while extracting user info from context it's incomplete ")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevData, err := s.repo.FindById(ctx, id)
	if err != nil {
		s.logger.Errorf("[EnabledUssdMerchant] error whil checking existing information error: %v", err)
		return err
	}

	if prevData.Enabled {
		s.logger.Errorf("[EnabledUssdMerchant] current merchant already enabled")
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	currentData := prevData

	currentData.Enabled = true

	cpsActionModel := lib.CpsModelBuilder(id, makerData, prevData, currentData, constants.RequestEnableUssdMerchant, constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[EnabledUssdMerchant] error whil creating CPS action: %v", err)
		return err
	}

	return nil
}

func (s *ussdMerchantService) DisableUssdMerchant(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableUssdMerchant", "UssdMerchant", "Enable")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[DisableUssdMerchant] Error while extracting user info from context it's incomplete ")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevData, err := s.repo.FindById(ctx, id)
	if err != nil {
		s.logger.Errorf("[DisableUssdMerchant] error whil checking existing information error: %v", err)
		return err
	}

	if !prevData.Enabled {
		s.logger.Errorf("[DisableUssdMerchant] current merchant already disabled")
		return errors.New(localization.ErrorAlreadyDisabled.Code)

	}
	currentData := prevData

	currentData.Enabled = false
	cpsActionModel := lib.CpsModelBuilder(id, makerData, prevData, currentData, constants.RequestDisableUssdMerchant, constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("[CreateUssdMerchant] error whil creating CPS action: %v", err)
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
		s.logger.Errorf("[UpdateUssdMerchantService] Error while extracting user info from context it's incomplete ")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prevMerchant, err := s.repo.FindById(ctx, id)
	if err != nil {
		s.logger.Errorf("[UpdateUssdMerchantService] error whil checking existing information error: %v", err)
		return err
	}

	if req.PhoneNumber != "" || req.Email != "" || req.AccountNumber != "" {
		existing, err := s.repo.FindByOr(ctx, req.PhoneNumber, req.Email, req.AccountNumber)
		if err != nil {
			s.logger.Errorf("[UpdateUssdMerchantService] error whil checking existing information error: %v", err)
			return err
		}
		if err := core.ExistingIdentifierForUpdate(existing, id, req); err != nil {
			s.logger.Infof("[UpdateUssdMerchantService] the entered data is already existed error: %v", err)
			return err
		}
	}

	if req.Service != "" && !strings.EqualFold(req.Service, prevMerchant.Service) {
		s.logger.Infof("[UpdateUssdMerchantService] service is not empty, checking if it's valid")
		_, err := s.serviceRepo.FindByID(ctx, req.Service)
		if err != nil {
			s.logger.Errorf("[UpdateUssdMerchantService] service is not valid, error: %v", err)
			return errors.New(localization.ErrorServiceNotFound.Code)
		}
		s.logger.Infof("[UpdateUssdMerchantService] service is valid")
	}

	if err := core.ExistingIdentifierForUpdate(existing, id, req); err != nil {
		s.logger.Infof("[UpdateUssdMerchantService] the entered data is already existed error: %v", err)
		return err
	}

	if req.Logo != nil {
		URL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Logo, string(constants.BankFolderName), s.cfg, "", s.logger)
		if err != nil {
			s.logger.Errorf("[UpdateUssdMerchantService] failed to upload logo: %v", err)
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
		s.logger.Errorf("[UpdateUssdMerchantService] failed to create CPS action: %v", err)
		return err
	}

	s.logger.Infof("UpdateUssdMerchantService seccussfuly implemented")
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
		s.logger.Errorf("[AuthorizeUssdMerchant] failed to unmarshal current action: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch action {
	case string(constants.RequestCreateUssdMerchant):
		if err := s.repo.Create(ctx, *curMerchant); err != nil {
			s.logger.Errorf("[AuthorizeUssdMerchant] failed to create ussd merchant: %v", err)
			return nil, errors.New(localization.ErrorUnhandledServer.Code)
		}
	case string(constants.RequestUpdateUssdMerchant):
		update := core.ModelToBson(curMerchant)
		if err := s.repo.Update(ctx, cpsAction.UniqueId, update); err != nil {
			s.logger.Errorf("[AuthorizeUssdMerchant] failed to update ussd merchant: %v", err)
			return nil, errors.New(localization.ErrorUnhandledServer.Code)
		}
	case string(constants.RequestEnableUssdMerchant):
		if err := s.repo.Update(ctx, cpsAction.UniqueId, bson.M{"enabled": true}); err != nil {
			s.logger.Errorf("[AuthorizeUssdMerchant] failed to enable ussd merchant: %v", err)
			return nil, err
		}
	case string(constants.RequestDisableUssdMerchant):
		if err := s.repo.Update(ctx, cpsAction.UniqueId, bson.M{"enabled": false}); err != nil {
			s.logger.Errorf("[AuthorizeUssdMerchant] failed to disable ussd merchant: %v", err)
			return nil, err
		}
	default:
		s.logger.Errorf("[AuthorizeUssdMerchant] invalid action: %s", action)
		return nil, errors.New(localization.ErrorInvalidAction.Code)
	}

	return cpsAction, nil
}
