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

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerKYCService struct {
	repo           storage.CustomerKYCRepository
	userRepo       storage.CustomerKYCRepository
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

// func (s *customerKYCService) Create(ctx context.Context, req dto.CreateCustomerKYCRequest) error {
// 	ctx, span := local_util.TraceLogger(ctx, "service", "CreateCustomerKYC", "CustomerKYC", "Create")
// 	defer span.End()

// 	makerData := local_util.ExtractUserFromContext(ctx)
// 	if local_util.IsIncomplete(makerData) {
// 		s.logger.Errorf("[CustKycSvc][Create] incomplete user")
// 		return errors.New(constants.IncompleteUserInfo)
// 	}

// 	var (
// 		idCardFront string
// 		idCardBack  string
// 		video       string
// 		err         error
// 	)

// 	if req.LivenessCheck.IDCardFront != nil {
// 		idCardFront, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.LivenessCheck.IDCardFront, string(constants.CustomerKYCFolderName), *s.cfg, "", s.logger)
// 		if err != nil {
// 			span.AddEvent("Failed to files", trace.WithAttributes(
// 				attribute.String("error", err.Error()),
// 			))
// 			return err
// 		}
// 	}

// 	if req.LivenessCheck.IDCardBack != nil {
// 		idCardBack, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.LivenessCheck.IDCardBack, string(constants.CustomerKYCFolderName), *s.cfg, "", s.logger)
// 		if err != nil {
// 			span.AddEvent("Failed to files", trace.WithAttributes(
// 				attribute.String("error", err.Error()),
// 			))
// 			return err
// 		}
// 	}

// 	if req.LivenessCheck.LivenessCheckVideo != nil {
// 		video, err = lib.UploadVideoToMinio(ctx, s.minio, s.bucketName, req.LivenessCheck.LivenessCheckVideo, string(constants.CustomerKYCFolderName), *s.cfg, "", s.logger)
// 		if err != nil {
// 			span.AddEvent("Failed to files", trace.WithAttributes(
// 				attribute.String("error", err.Error()),
// 			))
// 			return err
// 		}
// 	}

// 	kyc := &imodel.CustomerKYC{
// 		CustomerCode: local_util.GenerateCustomerCode(),
// 		AccountType:  constants.AccountType(req.AccountType),
// 		CustomerName: imodel.CustomerInfo{
// 			FirstName:   req.CustomerName.FirstName,
// 			MiddleName:  req.CustomerName.MiddleName,
// 			LastName:    req.CustomerName.LastName,
// 			PhoneNumber: req.CustomerName.PhoneNumber,
// 			Email:       req.CustomerName.Email,
// 			DateOfBirth: req.CustomerName.DateOfBirth,
// 			Gender:      req.CustomerName.Gender,
// 			MotherName:  req.CustomerName.MotherName,
// 		},
// 		Address: imodel.Address{
// 			Country:     req.Address.Country,
// 			Region:      req.Address.Region,
// 			City:        req.Address.City,
// 			SubCity:     req.Address.SubCity,
// 			Wereda:      req.Address.Wereda,
// 			Kebele:      req.Address.Kebele,
// 			HouseNumber: req.Address.HouseNumber,
// 		},
// 		Nationality:          req.Nationality,
// 		MaritalStatus:        constants.MaritalStatus(req.MaritalStatus),
// 		CustomerStatus:       constants.CustomerPending,
// 		EmploymentStatus:     constants.EmploymentStatus(req.EmploymentStatus),
// 		Occupation:           req.Occupation,
// 		AverageMonthlyIncome: req.AverageMonthlyIncome,
// 		EducationStatus:      req.EducationStatus,
// 		SourceOfFund:         req.SourceOfFund,
// 		KYCStatus:            "PENDING",
// 		MoneyLaunderingFree:  true,
// 		TermsAndConditions:   req.TermsAndConditions,
// 		LivenessCheck: imodel.LivenessCheck{
// 			IDCardFront:        idCardFront,
// 			IDCardBack:         idCardBack,
// 			LivenessCheckVideo: video,
// 		},
// 		VerificationResult: imodel.VerificationResult{
// 			FaceMatchScore:             req.VerificationResult.FaceMatchScore,
// 			LivenessResult:             req.VerificationResult.LivenessResult,
// 			DocumentAuthenticityResult: req.VerificationResult.DocumentAuthenticityResult,
// 		},
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	action := lib.CpsModelBuilder("", makerData, nil, kyc, string(constants.RequestCreateCustomerKYC), constants.CREATE)

// 	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
// 		s.logger.Errorf("[CustKycSvc][Create] cps action err: %v", err)
// 		return err
// 	}

// 	return nil
// }

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
	makerUser := local_util.ExtractUserFromContext(ctx)

	userReq, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[CustKycSvc][EnableDisable] find err: %v", err)
		return err
	}

	if userReq.Enabled && enable && userReq.KYCStatus == imodel.KYCStatusApproved {
		s.logger.Warnf("[CustKycSvc][EnableDisable] already in state id: %s, enabled: %v", id, enable)
		return errors.New("Customer KYC is already approved")
	}

	if !userReq.Enabled && !enable && userReq.KYCStatus == imodel.KYCStatusRejected {
		s.logger.Warnf("[CustKycSvc][EnableDisable] already in state id: %s, disabled: %v", id, enable)
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
	if !enable {
		newReq.KYCRejectReason = reason
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, userReq, newReq, string(action), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		s.logger.Errorf("[CustKycSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	return nil
}

// func (s *customerKYCService) Delete(ctx context.Context, id string) error {
// 	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteKYC", "CustomerKYC", "Delete")
// 	defer span.End()

// 	s.logger.Infof("[CustKycSvc][Delete] id: %s", id)
// 	makerData := local_util.ExtractUserFromContext(ctx)
// 	if local_util.IsIncomplete(makerData) {
// 		s.logger.Errorf("[CustKycSvc][Delete] incomplete user")
// 		return errors.New(constants.IncompleteUserInfo)
// 	}

// 	kyc, err := s.repo.FindByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	action := lib.CpsModelBuilder(id, makerData, kyc, nil, string(constants.RequestDeleteCustomerKYC), constants.DELETE)

// 	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
// 		s.logger.Errorf("[CustKycSvc][Delete] cps action err: %v", err)
// 		return err
// 	}

// 	return nil
// }

// func (s *customerKYCService) UpdateKYCStatus(ctx context.Context, id string, status string) error {
// 	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateKYCStatus", "CustomerKYC", "UpdateKYCStatus")
// 	defer span.End()

// 	s.logger.Infof("[CustKycSvc][UpdateStatus] id: %s status: %s", id, status)
// 	makerData := local_util.ExtractUserFromContext(ctx)
// 	if local_util.IsIncomplete(makerData) {
// 		s.logger.Errorf("[CustKycSvc][UpdateStatus] incomplete user")
// 		return errors.New(constants.IncompleteUserInfo)
// 	}

// 	kyc, err := s.repo.FindByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	currentAction := *kyc
// 	// currentAction.KYCStatus = status

// 	action := lib.CpsModelBuilder(id, makerData, kyc, currentAction, string(constants.RequestUpdateCustomerKYC), constants.UPDATE)

// 	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
// 		s.logger.Errorf("[CustKycSvc][UpdateStatus] cps action err: %v", err)
// 		return err
// 	}

// 	return nil
// }

func (s *customerKYCService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "CustomerKYC", "Authorize")
	defer span.End()
	s.logger.Infof("[CustKycSvc][Authorize] action: %s", cpsAction.RequestAction)

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestApproveCustomerKYC):
		userData, err := local_util.JsonUnmarshal[imodel.CustomerKYC](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}
		// Create the account in core
		data := accountLookupDto.CreateAccountRequest{
			CustomerName:      userData.KYCData.FullName,
			Gender:            constants.Gender(userData.KYCData.Gender),
			PhoneNumber:       userData.KYCData.PhoneNumber,
			AccountType:       "",
			AccountBranchType: "",
			Picture:           userData.KYCData.SelfiePhoto,
		}
		userAccount, err := core.CreateAccountToCore(ctx, data, s.accountService, s.logger)
		if err != nil {
			s.logger.Errorf("Core account creation failed: %v", err)
			return nil, err
		}

		// Create the account in user table
		err = s.repo.CreateUser(ctx, userAccount, *userData)
		if err != nil {
			return nil, err
		}
		// Update customer status
		err = s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, string(constants.KYCStatusApproved), "", true)
		if err != nil {
			return nil, err
		}
	case string(constants.RequestRejectCustomerKYC):
		userData, err := local_util.JsonUnmarshal[imodel.CustomerKYC](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}

		rejectionReason := userData.KYCRejectReason
		err = s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, string(constants.KYCStatusRejected), rejectionReason, false)
		if err != nil {
			return nil, err
		}

	// case string(constants.RequestCreateCustomerKYC):
	// 	data, err := local_util.JsonUnmarshal[imodel.CustomerKYC](cpsAction.CurrentAction)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if err := s.repo.Create(ctx, data); err != nil {
	// 		return nil, err
	// 	}
	// case string(constants.RequestDeleteCustomerKYC):
	// 	if err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
	// 		return nil, err
	// 	}
	// case string(constants.RequestUpdateCustomerKYC):
	// 	data, err := local_util.JsonUnmarshal[dto.UpdateKYCStatusRequest](cpsAction.CurrentAction)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if err := s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, data.KYCStatus); err != nil {
	// 		return nil, err
	// 	}
	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}

	return cpsAction, nil
}
