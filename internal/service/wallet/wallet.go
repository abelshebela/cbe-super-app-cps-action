package wallet

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/wallet/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	walletCatch "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/catch/wallet"


	"path"

	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_model "cbe-super-app-cps-action/internal/constants/model"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type walletService struct {
	repo        storage.WalletOracleRepository
	cpsService  service.CPSActionService
	serviceRepo storage.ServicesRepository
	logger      utils.Logger
	minio       *s3.Client
	bucketName  string
	minioPubUrl string
	cfg         *config.VaultConfig
	walletCatch  walletCatch.WalletCatch
}

func NewWalletService(repo storage.WalletOracleRepository, cps service.CPSActionService, serviceRepo storage.ServicesRepository, minio *s3.Client, minioPubUrl string, bucketName string, cfg *config.VaultConfig,walletCatch walletCatch.WalletCatch, logger utils.Logger) service.WalletService {
	return &walletService{
		repo:        repo,
		cpsService:  cps,
		serviceRepo: serviceRepo,
		logger:      logger,
		minio:       minio,
		bucketName:  bucketName,
		cfg:         cfg,
		walletCatch: walletCatch,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, req walletDto.WalletRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateWallet", "walletService", "walletService")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[WalletSvc][Create] name: %s", req.Name)

	// Check name, code, and service_id independently (Find uses OR only when multiple args are set; we run ordered checks for clear errors).
	if strings.TrimSpace(req.Name) != "" {
		exist, err := s.repo.Find(ctx, "", req.Name, "")
		if err != nil {
			span.AddEvent("Repo find error", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		if exist != nil {
			span.AddEvent("Wallet name already exists", trace.WithAttributes(attribute.String("name", req.Name)))
			return errors.New(localization.ErrorWalletNameAlreadyExists.Code)
		}
	}
	if strings.TrimSpace(req.UniqueCode) != "" {
		exist, err := s.repo.Find(ctx, req.UniqueCode, "", "")
		if err != nil {
			span.AddEvent("Repo find error", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		if exist != nil {
			span.AddEvent("Wallet code already exists", trace.WithAttributes(attribute.String("unique_code", req.UniqueCode)))
			return errors.New(localization.ErrorWalletCodeAlreadyExists.Code)
		}
	}
	services, err := core.CheckServices(req.Self, req.Other, req.Agent, req.SelfServiceID, req.OtherServiceID, req.AgentServiceID, s.serviceRepo, ctx, s.logger)
	if err != nil {
		return err
	}

	if err := core.CheckServiceIDInWalletService(ctx, req.SelfServiceID, req.OtherServiceID, req.AgentServiceID, s.repo, nil, s.logger); err != nil {
		return err
	}

	// code := strings.ToUpper(strings.TrimSpace(req.UniqueCode))

	URL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, string(constants.WalletFolderName), *s.cfg, "", s.logger)
	if err != nil {
		span.AddEvent("UploadFileToMinio error", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	wallet := core.ToCreateWalletDoc(req, URL, services, nil)
	wallet.Enabled = false

	//here since the unique id is nil 000.. use other unique id like the code
	// if err := core.HandleCPSAction(ctx, s.cpsService, wallet.ID.Hex(), constants.RequestCreateWallet, wallet, nil, constants.ActionCreate); err != nil {
	if err := core.HandleCPSAction(ctx, s.cpsService, "", constants.RequestCreateWallet, wallet, nil, constants.ActionCreate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("unique_code", wallet.UniqueCode)))
		log.Errorf("[WalletSvc][Create] cps action err: %v", err)
		return err
	}

	span.AddEvent("Wallet created", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
	return nil
}

func (s *walletService) UpdateWallet(ctx context.Context, id string, req walletDto.WalletRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateWallet", "walletService", "walletService")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[WalletSvc][Update] id: %s", id)

	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return err
	}

	if req.Name != "" {
		exist, err := s.repo.Find(ctx, "", req.Name, "")
		if err != nil {
			span.AddEvent("Repo find error", trace.WithAttributes(attribute.String("error", err.Error())))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if exist != nil && !strings.EqualFold(strings.TrimSpace(exist.ID), strings.TrimSpace(id)) {
			span.AddEvent("Wallet name already exists", trace.WithAttributes(attribute.String("name", req.Name)))
			log.Infof("[WalletSvc][Update] wallet name already exists name: %s req.id: %s found.id: %s", req.Name, id, exist.ID)
			return errors.New(localization.ErrorWalletNameAlreadyExists.Code)
		}
	}
	if req.UniqueCode != "" {
		exist, err := s.repo.Find(ctx, req.UniqueCode, "", "")
		if err != nil {
			span.AddEvent("Repo find error", trace.WithAttributes(attribute.String("error", err.Error())))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if exist != nil && !strings.EqualFold(strings.TrimSpace(exist.ID), strings.TrimSpace(id)) {
			span.AddEvent("Wallet code already exists", trace.WithAttributes(attribute.String("unique_code", req.UniqueCode)))
			log.Infof("[WalletSvc][Update] wallet code already exists unique_code: %s req.id: %s found.id: %s", req.UniqueCode, id, exist.ID)
			return errors.New(localization.ErrorWalletCodeAlreadyExists.Code)
		}
	}

	var avatarURL string
	if req.Avatar != nil {
		var objectkey string
		if prevWallet.Avatar != "" {
			objectkey = path.Base(prevWallet.Avatar)
		}

		avatarURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, string(constants.WalletFolderName), *s.cfg, objectkey, s.logger)
		if err != nil {
			span.AddEvent("UploadFileToMinio error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	} else {
		avatarURL = prevWallet.Avatar
	}
	services, err := core.CheckServices(req.Self, req.Other, req.Agent, req.SelfServiceID, req.OtherServiceID, req.AgentServiceID, s.serviceRepo, ctx, s.logger)
	if err != nil {
		return err
	}

	if err := core.CheckServiceIDInWalletService(ctx, req.SelfServiceID, req.OtherServiceID, req.AgentServiceID, s.repo, prevWallet, s.logger); err != nil {
		return err
	}

	wallet := core.ToCreateWalletDoc(req, avatarURL, services, prevWallet)

	UpdateWallet, change_count := core.ToUpdateWalletDoc(*prevWallet, *wallet)
	if avatarURL != prevWallet.Avatar {
		change_count++
	}
	UpdateWallet.Avatar = avatarURL

	if change_count == 0 {
		span.AddEvent("No changes detected", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorNoChangesDetected.Code)
	}
	// Convert struct to map with json tags as keys
	b, _ := json.Marshal(UpdateWallet)
	var walletMap map[string]interface{}
	json.Unmarshal(b, &walletMap)

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestUpdateWallet, walletMap, *prevWallet, constants.ActionUpdate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("unique_code", UpdateWallet.UniqueCode)))
		log.Errorf("[WalletSvc][Update] cps action err: %v", err)
		return err
	}

	span.AddEvent("Wallet updated", trace.WithAttributes(attribute.String("id", id)))
	return nil
}

func (s *walletService) DeleteWallet(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteWallet", "walletService", "walletService")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorWalletNotFound.Code)
	}

	now := time.Now()
	deletedWallet := *prevWallet
	deletedWallet.IsDeleted = true
	deletedWallet.DeletedAt = &now

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestDeleteWallet, deletedWallet, *prevWallet, constants.ActionDelete); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("unique_code", deletedWallet.UniqueCode)))
		log.Errorf("[WalletSvc][Delete] cps action err: %v", err)
		return err
	}

	span.AddEvent("Wallet deleted", trace.WithAttributes(attribute.String("id", id)))
	return nil
}

func (s *walletService) EnableOrDisableWallet(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableWallet", "walletService", "walletService")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorWalletNotFound.Code)
	}

	if enable && prevWallet.Enabled {
		span.AddEvent("Wallet already enabled", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorWalletAlreadyEnabled.Code)
	}
	if !enable && !prevWallet.Enabled {
		span.AddEvent("Wallet already disabled", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorWalletAlreadyDisabled.Code)
	}

	updatedWallet := *prevWallet
	updatedWallet.Enabled = enable
	updatedWallet.LastModifiedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableWallet
	} else {
		action = constants.RequestDisableWallet
	}

	if err := core.HandleCPSAction(ctx, s.cpsService, id, action, updatedWallet, *prevWallet, constants.ActionUpdate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("unique_code", updatedWallet.UniqueCode)))
		log.Errorf("[WalletSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	span.AddEvent("Wallet enable/disable updated", trace.WithAttributes(attribute.String("id", id), attribute.Bool("enabled", enable)))
	return nil
}

func (s *walletService) EnableOrDisableWalletService(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableWalletService", "walletService", "walletService")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	prevWalletService, err := s.repo.FindWalletServiceByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorWalletServiceNotFound.Code)
	}

	if enable && prevWalletService.IsEnabled == 1 {
		span.AddEvent("Wallet service already enabled", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorWalletServiceAlreadyEnabled.Code)
	}
	if !enable && prevWalletService.IsEnabled == 0 {
		span.AddEvent("Wallet service already disabled", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorWalletServiceAlreadyDisabled.Code)
	}

	updatedWalletService := *prevWalletService
	if enable {
		updatedWalletService.IsEnabled = 1
	} else {
		updatedWalletService.IsEnabled = 0

	}
	updatedWalletService.LastModifiedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableWalletService
	} else {
		action = constants.RequestDisableWalletService
	}

	if err := core.HandleCPSAction(ctx, s.cpsService, id, action, updatedWalletService, *prevWalletService, constants.ActionUpdate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("service_id", updatedWalletService.ID)))
		log.Errorf("[WalletSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	span.AddEvent("Wallet service enable/disable updated", trace.WithAttributes(attribute.String("id", id), attribute.Bool("enabled", enable)))
	return nil
}

func (s *walletService) GetWallet(ctx context.Context, id string) (*local_model.WalletOracle, error) {
	// Same enrichment as list/gRPC: join services + access_lists for SERVICE_CODE / SERVICE_KEY.
	w, err := s.repo.FindByIDForGRPC(ctx, id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, errors.New(localization.ErrorWalletNotFound.Code)
	}
	return w, nil
}

func (s *walletService) GetAllWallet(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]local_model.WalletOracle], error) {
	return s.repo.FindAllWithPaginationForGRPC(ctx, filterParams)
}

// GetAllWalletForGRPC implements service.WalletService.
func (s *walletService) GetAllWalletForGRPC(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]local_model.WalletOracle], error) {
	return s.repo.FindAllWithPaginationForGRPC(ctx, filterParams)
}

// GetWalletForGRPC implements service.WalletService.
func (s *walletService) GetWalletForGRPC(ctx context.Context, id string) (*local_model.WalletOracle, error) {
	w, err := s.repo.FindByIDForGRPC(ctx, id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, errors.New(localization.ErrorWalletNotFound.Code)
	}
	return w, nil
}

func (s *walletService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "walletService", "walletService")
	defer span.End()

	b, err := json.Marshal(action.CurrentAction)
	if err != nil {
		return nil, err
	}

	var w map[string]interface{}
	err = json.Unmarshal(b, &w)
	if err != nil {
		span.AddEvent("JsonUnmarshal error", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}
	wallet := core.MapToModel(w, s.logger)

	

	switch action.RequestAction {
	case string(constants.RequestCreateWallet):
		span.AddEvent("Creating wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		if strings.TrimSpace(wallet.Name) != "" {
			if exist, e := s.repo.Find(ctx, "", wallet.Name, ""); e != nil {
				return nil, e
			} else if exist != nil {
				return nil, errors.New(localization.ErrorWalletNameAlreadyExists.Code)
			}
		}
		if strings.TrimSpace(wallet.UniqueCode) != "" {
			if exist, e := s.repo.Find(ctx, wallet.UniqueCode, "", ""); e != nil {
				return nil, e
			} else if exist != nil {
				return nil, errors.New(localization.ErrorWalletCodeAlreadyExists.Code)
			}
		}
		err = s.repo.Create(ctx, wallet)
		if err != nil{
			s.logger.Errorf("[wallet service authorizor create] error:%v",err)
			return  nil,err
		}
		enabledWallerService := buildEnableWalletServices(wallet)
		
		err := s.walletCatch.Set(ctx,walletCatch.WalletData{
				Name: wallet.Name,
				UniqueCode: wallet.UniqueCode,
				ServiceID: wallet.ID,
				Logo: wallet.Avatar,
				IsWalletEnabled: wallet.Enabled,
				ServicesType:enabledWallerService, 
				
			})
		if err != nil{
			s.logger.Errorf("[wallet service Authorizor create] unable to set on redis err:%v",err)
			return nil,err
		}
	case string(constants.RequestUpdateWallet):
		span.AddEvent("Updating wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.Update(ctx, action.UniqueId, wallet)
		if err != nil{
				s.logger.Errorf("[wallet service authorizor update] error:%v",err)
				return  nil,err
		}
		privUniqueCode,err:= s.getActionCodeFromPrivAction(action)
		if err != nil{
			return nil,err
		}
		enabledWallerService := buildEnableWalletServices(wallet)
		_,err = s.walletCatch.Update(ctx,privUniqueCode,walletCatch.WalletData{
				Name: wallet.Name,
				ServiceID: wallet.ID,
				UniqueCode: wallet.UniqueCode,
				Logo: wallet.Avatar,
				IsWalletEnabled: wallet.Enabled,
				ServicesType:enabledWallerService,
				
			})

			if err != nil{
				if err.Error() == walletCatch.ErrKeyNotFoundInCatch.Error(){
					_,walErr := s.walletCatch.Update(ctx,wallet.UniqueCode,walletCatch.WalletData{
							Name: wallet.Name,
							ServiceID: wallet.ID,
							UniqueCode: wallet.UniqueCode,
							Logo: wallet.Avatar,
							IsWalletEnabled: wallet.Enabled,
							ServicesType:enabledWallerService,
						})
					if walErr != nil{
					s.logger.Errorf("[wallet service Authorizor update] unable to set on redis err:%v",walErr)
					}
					return nil,walErr
					}
					
			s.logger.Errorf("[wallet service Authorizor update] unable to update on redis err:%v",err)
			return nil,err
		}

	case string(constants.RequestDeleteWallet):
		span.AddEvent("Deleting wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.Delete(ctx, action.UniqueId)
		if err != nil{
			s.logger.Errorf("[wallet service authorizor update] error:%v",err)

		}
		
		err := s.walletCatch.Delete(ctx,wallet.UniqueCode)
		if err != nil{
			s.logger.Errorf("[wallet service Authorizor update] unable to set on redis err:%v",err)
			return nil,err		
		}

	case string(constants.RequestEnableWallet):
		span.AddEvent("Enabling wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, true)

		if err != nil{
				s.logger.Errorf("[wallet service authorizor update] error:%v",err)
				return  nil,err
		}
      enabledWallerService := buildEnableWalletServices(wallet)
		err := s.walletCatch.Set(ctx,walletCatch.WalletData{
				Name: wallet.Name,
				ServiceID: wallet.ID,
				UniqueCode: wallet.UniqueCode,
				Logo: wallet.Avatar,
				IsWalletEnabled: true,
				ServicesType: enabledWallerService,
			})

			if err != nil{
			s.logger.Errorf("[wallet service Authorizor enable] unable to set on redis err:%v",err)
			return nil,err
		}
	case string(constants.RequestDisableWallet):
		span.AddEvent("Disabling wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, false)
		if err != nil{
				s.logger.Errorf("[wallet service authorizor update] error:%v",err)
				return  nil,err
		}
		enabledWallerService := buildEnableWalletServices(wallet)
		err := s.walletCatch.Set(ctx,walletCatch.WalletData{
				ServiceID: wallet.ID,
				Name: wallet.Name,
				UniqueCode: wallet.UniqueCode,
				Logo: wallet.Avatar,
				IsWalletEnabled: false,
				ServicesType: enabledWallerService,
			})

			if err != nil{
			s.logger.Errorf("[wallet service Authorizor update] unable to set on redis err:%v",err)
			return nil,err
		}
	case string(constants.RequestEnableWalletService):
		span.AddEvent("Enabling wallet service", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.EnableOrDisableService(ctx, action.UniqueId, true)
		if err != nil{
				s.logger.Errorf("[wallet service authorizor update] error:%v",err)
				return  nil,err
		}
		enabledWalletService := buildEnableWalletServices(wallet)
		_,err := s.walletCatch.Update(ctx,wallet.UniqueCode,walletCatch.WalletData{
				ServiceID: wallet.ID,
				Name: wallet.Name,
				UniqueCode: wallet.UniqueCode,
				Logo: wallet.Avatar,
				IsWalletEnabled:wallet.Enabled,
				ServicesType: enabledWalletService,
					
			})

			if err != nil{
			s.logger.Errorf("[wallet service Authorizor update] unable to set on redis err:%v",err)
			return nil,err
		}

	case string(constants.RequestDisableWalletService):
		span.AddEvent("Disabling wallet service", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.EnableOrDisableService(ctx, action.UniqueId, false)
		
		if err != nil{
				s.logger.Errorf("[wallet service authorizor update] error:%v",err)
				return  nil,err
		}

		enabledWalletService := buildEnableWalletServices(wallet)
		_,err := s.walletCatch.Update(ctx,wallet.UniqueCode,walletCatch.WalletData{
				ServiceID: wallet.ID,
				Name: wallet.Name,
				UniqueCode: wallet.UniqueCode,
				Logo: wallet.Avatar,
				IsWalletEnabled:wallet.Enabled,
				ServicesType: enabledWalletService,
				
			})

			if err != nil{
			s.logger.Errorf("[wallet service Authorizor update] unable to set on redis err:%v",err)
			return nil,err
			}

	default:
		span.AddEvent("Unsupported action", trace.WithAttributes(attribute.String("action", action.RequestAction)))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		return nil, err
	}

	action.CurrentAction = wallet
	return action, nil
}


func buildEnableWalletServices(wallet *local_model.WalletOracle)[]walletCatch.WalletServiceData{
var walletServices = []walletCatch.WalletServiceData{}
		if wallet.SelfServiceEnabled == 1{
			walletServices = append(walletServices, walletCatch.WalletServiceData{
				ServiceID: wallet.SelfServiceID,
				ServiceType: walletCatch.ServiceType(wallet.SelfServiceName),

			})
		}
		if wallet.OtherServiceEnabled == 1{
			walletServices = append(walletServices, walletCatch.WalletServiceData{
				ServiceID: wallet.OtherServiceID,
				ServiceType: walletCatch.ServiceType(wallet.OtherServiceName),

			})
		}
		if wallet.AgentServiceEnabled == 1{
			walletServices = append(walletServices, walletCatch.WalletServiceData{
				ServiceID: wallet.AgentServiceID,
				ServiceType: walletCatch.ServiceType(wallet.AgentServiceName),

			})
		}
		return  walletServices
}


func (s *walletService)getActionCodeFromPrivAction(action *model.CPSAction) (string,error) {
		b, err := json.Marshal(action.PreviousAction)
		if err != nil {
			return "", err
		}

		var w map[string]interface{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			return "", errors.New(localization.ErrorInvalidActionData.Code)
		}

		privWallet := core.MapToModel(w, s.logger)
		return  privWallet.UniqueCode,nil
	}