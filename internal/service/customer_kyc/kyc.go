package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	accountLookupDto "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	"cbe-super-app-cps-action/internal/constants/lib"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/customer_kyc/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerKYCService struct {
	repo           storage.CustomerKYCRepository
	cpsService     service.CPSActionService
	accountService account_lookup.Account
	logger         utils.Logger
	minio          *s3.Client
	bucketName     string
	cfg            *config.VaultConfig
	minioEndPoint  string
}

func NewCustomerKYCService(repo storage.CustomerKYCRepository,
	cpsService service.CPSActionService,
	accountLookUpService account_lookup.Account,
	logger utils.Logger,
	minio *s3.Client,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
) service.CustomerKYCService {
	return &customerKYCService{
		repo:           repo,
		cpsService:     cpsService,
		accountService: accountLookUpService,
		logger:         logger,
		minio:          minio,
		bucketName:     bucketName,
		cfg:            cfg,
		minioEndPoint:  minioEndPoint,
	}
}

func (s *customerKYCService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]dto.CustomerKYCResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "CustomerKYC", "FindAllWithPagination")
	defer span.End()

	result, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		return nil, err
	}

	mappedResponse := core.MapCustomerKYCToResponsePaginated(result)

	return mappedResponse, nil
}

func (s *customerKYCService) FindByID(ctx context.Context, id string) (*dto.CustomerKYCResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "CustomerKYC", "FindByID")
	defer span.End()

	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	mappedResponse := core.MapCustomerKYCToResponse(result)

	return mappedResponse, nil
}

func (s *customerKYCService) EnableOrDisable(ctx context.Context, id, reason string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	if local_util.IsIncomplete(makerUser) {
		log.Errorf("[CustKycSvc][EnableDisable] incomplete maker user")
		return errors.New(constants.IncompleteUserInfo)
	}

	userReq, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustKycSvc][EnableDisable] find err: %v", err)
		return err
	}

	if enable && userReq.KYCStatus == imodel.KYCStatusApproved {
		log.Warnf("[CustKycSvc][EnableDisable] already approved id: %s", id)
		return errors.New("Customer KYC is already approved")
	}

	if !enable && userReq.KYCStatus == imodel.KYCStatusRejected {
		log.Warnf("[CustKycSvc][EnableDisable] already rejected id: %s", id)
		return errors.New("Customer KYC is already rejected")
	}

	var action constants.RequestAction
	if enable {
		action = constants.RequestApproveCustomerKYC
	} else {
		action = constants.RequestRejectCustomerKYC
	}

	newReq := *userReq
	newReq.Enabled = enable
	if enable {
		newReq.KYCStatus = imodel.KYCStatusApproved
	} else {
		newReq.KYCStatus = imodel.KYCStatusRejected
		newReq.KYCRejectReason = reason
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, userReq, newReq, string(action), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustKycSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *customerKYCService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "CustomerKYC", "Authorize")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[CustKycSvc][Authorize] action: %s", cpsAction.RequestAction)

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestApproveCustomerKYC):
		userData, err := local_util.JsonUnmarshal[imodel.CustomerKYC](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}

		data := accountLookupDto.CreateAccountRequest{
			CustomerName:      userData.KYCData.FullName,
			Gender:            constants.Gender(userData.KYCData.Gender),
			PhoneNumber:       userData.KYCData.PhoneNumber,
			AccountType:       userData.KYCData.AccountType,
			AccountBranchType: "",
			Picture:           userData.KYCData.SelfiePhoto,
		}
		userAccount, err := core.CreateAccountToCore(ctx, data, s.accountService, s.logger)
		if err != nil {
			log.Errorf("[CustKycSvc][Authorize] core account creation failed: %v", err)
			return nil, err
		}

		if err = s.repo.CreateUser(ctx, userAccount, *userData); err != nil {
			return nil, err
		}

		if err = s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, string(constants.KYCStatusApproved), "", true); err != nil {
			return nil, err
		}

	case string(constants.RequestRejectCustomerKYC):
		userData, err := local_util.JsonUnmarshal[imodel.CustomerKYC](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}

		rejectionReason := strings.TrimSpace(userData.KYCRejectReason)
		if rejectionReason == "" {
			return nil, errors.New("rejection reason is required")
		}

		if err = s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, string(constants.KYCStatusRejected), rejectionReason, false); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}

	return cpsAction, nil
}
