package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	"cbe-super-app-cps-action/internal/constants/lib"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type customerKYCService struct {
	repo          storage.CustomerKYCRepository
	cpsService    service.CPSActionService
	logger        utils.Logger
	minio         *s3.Client
	bucketName    string
	cfg           *config.VaultConfig
	minioEndPoint string
}

func NewCustomerKYCService(repo storage.CustomerKYCRepository,
	cpsService service.CPSActionService,
	logger utils.Logger,
	minio *s3.Client,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
) service.CustomerKYCService {
	return &customerKYCService{
		repo:          repo,
		cpsService:    cpsService,
		logger:        logger,
		minio:         minio,
		bucketName:    bucketName,
		cfg:           cfg,
		minioEndPoint: minioEndPoint,
	}
}

func (s *customerKYCService) Create(ctx context.Context, req dto.CreateCustomerKYCRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateCustomerKYC", "CustomerKYC", "Create")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[CustomerKYCService.Create] incomplete user data")
		return errors.New(constants.IncompleteUserInfo)
	}

	var (
		idCardFront string
		idCardBack  string
		video       string
		err         error
	)

	if req.LivenessCheck.IDCardFront != nil {
		idCardFront, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.LivenessCheck.IDCardFront, string(constants.CustomerKYCFolderName), *s.cfg, "", s.logger)
		if err != nil {
			span.AddEvent("Failed to files", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			return err
		}
	}

	if req.LivenessCheck.IDCardBack != nil {
		idCardBack, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.LivenessCheck.IDCardBack, string(constants.CustomerKYCFolderName), *s.cfg, "", s.logger)
		if err != nil {
			span.AddEvent("Failed to files", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			return err
		}
	}

	if req.LivenessCheck.LivenessCheckVideo != nil {
		video, err = lib.UploadVideoToMinio(ctx, s.minio, s.bucketName, req.LivenessCheck.LivenessCheckVideo, string(constants.CustomerKYCFolderName), *s.cfg, "", s.logger)
		if err != nil {
			span.AddEvent("Failed to files", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			return err
		}
	}

	kyc := &imodel.CustomerKYC{
		CustomerCode: local_util.GenerateCustomerCode(),
		AccountType:  imodel.AccountType(req.AccountType),
		CustomerName: imodel.CustomerInfo{
			FirstName:   req.CustomerName.FirstName,
			MiddleName:  req.CustomerName.MiddleName,
			LastName:    req.CustomerName.LastName,
			PhoneNumber: req.CustomerName.PhoneNumber,
			Email:       req.CustomerName.Email,
			DateOfBirth: req.CustomerName.DateOfBirth,
			Gender:      req.CustomerName.Gender,
			MotherName:  req.CustomerName.MotherName,
		},
		Address: imodel.Address{
			Country:     req.Address.Country,
			Region:      req.Address.Region,
			City:        req.Address.City,
			SubCity:     req.Address.SubCity,
			Wereda:      req.Address.Wereda,
			Kebele:      req.Address.Kebele,
			HouseNumber: req.Address.HouseNumber,
		},
		Nationality:          req.Nationality,
		MaritalStatus:        imodel.MaritalStatus(req.MaritalStatus),
		CustomerStatus:       imodel.CustomerPending,
		EmploymentStatus:     imodel.EmploymentStatus(req.EmploymentStatus),
		Occupation:           req.Occupation,
		AverageMonthlyIncome: req.AverageMonthlyIncome,
		EducationStatus:      req.EducationStatus,
		SourceOfFund:         req.SourceOfFund,
		KYCStatus:            "PENDING",
		MoneyLaunderingFree:  true,
		TermsAndConditions:   req.TermsAndConditions,
		LivenessCheck: imodel.LivenessCheck{
			IDCardFront:        idCardFront,
			IDCardBack:         idCardBack,
			LivenessCheckVideo: video,
		},
		VerificationResult: imodel.VerificationResult{
			FaceMatchScore:             req.VerificationResult.FaceMatchScore,
			LivenessResult:             req.VerificationResult.LivenessResult,
			DocumentAuthenticityResult: req.VerificationResult.DocumentAuthenticityResult,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	action := lib.CpsModelBuilder("", makerData, nil, kyc, string(constants.RequestCreateCustomerKYC), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		s.logger.Errorf("[CustomerKYCService.Create] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (s *customerKYCService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.CustomerKYC], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "CustomerKYC", "FindAllWithPagination")
	defer span.End()

	return s.repo.FindAllWithPagination(ctx, *filterParam)
}

func (s *customerKYCService) FindByID(ctx context.Context, id string) (*imodel.CustomerKYC, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "CustomerKYC", "FindByID")
	defer span.End()

	return s.repo.FindByID(ctx, id)
}

func (s *customerKYCService) Delete(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteKYC", "CustomerKYC", "Delete")
	defer span.End()

	s.logger.Infof("[CustomerKYCService.Delete] deleting kyc request: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[CustomerKYCService.Delete] incomplete user data")
		return errors.New(constants.IncompleteUserInfo)
	}

	kyc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	action := lib.CpsModelBuilder(id, makerData, kyc, nil, string(constants.RequestDeleteCustomerKYC), constants.DELETE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		s.logger.Errorf("[CustomerKYCService.Delete] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (s *customerKYCService) UpdateKYCStatus(ctx context.Context, id string, status string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateKYCStatus", "CustomerKYC", "UpdateKYCStatus")
	defer span.End()

	s.logger.Infof("[CustomerKYCService.UpdateKYCStatus] updating kyc status for request: %s to %s", id, status)
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		s.logger.Errorf("[CustomerKYCService.UpdateKYCStatus] incomplete user data")
		return errors.New(constants.IncompleteUserInfo)
	}

	kyc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	currentAction := *kyc
	currentAction.KYCStatus = status

	action := lib.CpsModelBuilder(id, makerData, kyc, currentAction, string(constants.RequestUpdateCustomerKYC), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		s.logger.Errorf("[CustomerKYCService.UpdateKYCStatus] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

func (s *customerKYCService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "CustomerKYC", "Authorize")
	defer span.End()
	s.logger.Infof("[CustomerKYCService.Authorize] authorizing kyc action: %s", cpsAction.RequestAction)

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateCustomerKYC):
		data, err := local_util.JsonUnmarshal[imodel.CustomerKYC](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}
		if err := s.repo.Create(ctx, data); err != nil {
			return nil, err
		}
	case string(constants.RequestDeleteCustomerKYC):
		if err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			return nil, err
		}
	case string(constants.RequestUpdateCustomerKYC):
		data, err := local_util.JsonUnmarshal[dto.UpdateKYCStatusRequest](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}
		if err := s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, data.KYCStatus); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}

	return cpsAction, nil
}
