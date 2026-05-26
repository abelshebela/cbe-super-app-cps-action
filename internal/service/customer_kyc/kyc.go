package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	accountLookupDto "cbe-super-app-cps-action/internal/constants/dto/account_lookup"
	dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
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
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type customerKYCService struct {
	repo           storage.CustomerKYCRepository
	userRepo       storage.CustomerKYCRepository
	cpsService     service.CPSActionService
	accountService account_lookup.Account
	cpsUserRepo    storage.CpsUserRepository
	logger         utils.Logger
	minio          *s3.Client
	bucketName     string
	cfg            *config.VaultConfig
	minioEndPoint  string
}

func NewCustomerKYCService(repo storage.CustomerKYCRepository,
	cpsService service.CPSActionService,
	cpsUserRepo storage.CpsUserRepository,
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
		cpsUserRepo:    cpsUserRepo,
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
// 		log.Errorf("[CustKycSvc][Create] incomplete user")
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
// 		log.Errorf("[CustKycSvc][Create] cps action err: %v", err)
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
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	userReq, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustKycSvc][EnableDisable] find err: %v", err)
		return err
	}

	if userReq.KYCStatus == imodel.KYCStatusPending {
		log.Warnf("[CustKycSvc][EnableDisable] kyc status is pending id: %s, enabled: %v", id, enable)
		return errors.New("Start KYC review before approving or rejecting the KYC request")
	}

	// Check if review time is not expired before allowing approval or rejection of the KYC request
	if userReq.KYCStatus == imodel.KYCStatusInReview {
		review, err := s.repo.FindKycInReview(ctx, id)
		if err != nil {
			log.Errorf("[CustKycSvc][EnableDisable] failed to find kyc review: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		if review != nil && review.ExpiresAt.Before(time.Now()) {
			log.Warnf("[CustKycSvc][EnableDisable] kyc review expired id: %s, enabled: %v", id, enable)
			return errors.New("The KYC review period has expired.")
		}
	}

	if userReq.Enabled && enable && userReq.KYCStatus == imodel.KYCStatusApproved {
		log.Warnf("[CustKycSvc][EnableDisable] already in state id: %s, enabled: %v", id, enable)
		return errors.New("Customer KYC is already approved")
	}

	if !userReq.Enabled && !enable && userReq.KYCStatus == imodel.KYCStatusRejected {
		log.Warnf("[CustKycSvc][EnableDisable] already in state id: %s, disabled: %v", id, enable)
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
		log.Errorf("[CustKycSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	return nil
}

// func (s *customerKYCService) Delete(ctx context.Context, id string) error {

// 	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteKYC", "CustomerKYC", "Delete")
// 	defer span.End()

// 	log.Infof("[CustKycSvc][Delete] id: %s", id)
// 	makerData := local_util.ExtractUserFromContext(ctx)
// 	if local_util.IsIncomplete(makerData) {
// 		log.Errorf("[CustKycSvc][Delete] incomplete user")
// 		return errors.New(constants.IncompleteUserInfo)
// 	}

// 	kyc, err := s.repo.FindByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	action := lib.CpsModelBuilder(id, makerData, kyc, nil, string(constants.RequestDeleteCustomerKYC), constants.DELETE)

// 	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
// 		log.Errorf("[CustKycSvc][Delete] cps action err: %v", err)
// 		return err
// 	}

// 	return nil
// }

// func (s *customerKYCService) UpdateKYCStatus(ctx context.Context, id string, status string) error {

// 	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateKYCStatus", "CustomerKYC", "UpdateKYCStatus")
// 	defer span.End()

// 	log.Infof("[CustKycSvc][UpdateStatus] id: %s status: %s", id, status)
// 	makerData := local_util.ExtractUserFromContext(ctx)
// 	if local_util.IsIncomplete(makerData) {
// 		log.Errorf("[CustKycSvc][UpdateStatus] incomplete user")
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
// 		log.Errorf("[CustKycSvc][UpdateStatus] cps action err: %v", err)
// 		return err
// 	}

// 	return nil
// }

func (s *customerKYCService) StartKycReview(ctx context.Context, id string) (*imodel.StartedKycReview, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	kycRequest, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] find err: %v", err)
		return nil, err
	}

	if kycRequest.KYCStatus != imodel.KYCStatusPending {
		log.Warnf("[CustKycSvc][StartKycReview] kyc status is not pending id: %s, status: %s", id, kycRequest.KYCStatus)
		return nil, errors.New("This KYC request is under review or has already been reviewed.")
	}

	kycID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	userID, err := bson.ObjectIDFromHex(makerUser.UserID)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// prevent duplicate kyc action from being created if the review is already started
	existingReview, err := s.repo.FindKycInReview(ctx, id)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		log.Errorf("[CustKycSvc][StartKycReview] failed to check existing review: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if existingReview != nil {
		log.Warnf("[CustKycSvc][StartKycReview] review already started for kyc id: %s", id)
		return nil, errors.New("A review has already been started for this KYC request")
	}

	// Find the user who is starting the review and include their info in the review document
	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] failed to fetch user info: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	newReq := &imodel.StartedKycReview{
		KycID: kycID,
		Reviewer: imodel.UserInfo{
			ID:          userID,
			UserCode:    cpsUser.UserCode,
			FullName:    cpsUser.FullName,
			Email:       cpsUser.Email,
			PhoneNumber: cpsUser.PhoneNumber,
		},
		ReviewStatus:   string(imodel.KYCStatusInReview),
		StartedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(30 * time.Minute),
		IsActive:       true,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	newReview, err := s.repo.StartKycReview(ctx, newReq)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateKYCStatus(ctx, id, string(imodel.KYCStatusInReview), "", false)
	if err != nil {
		return nil, err
	}

	return newReview, nil
}

func (s *customerKYCService) PickKycReview(ctx context.Context, id string, reason string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	kycInReview, err := s.repo.FindKycInReview(ctx, id)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] find err: %v", err)
		return err
	}

	if kycInReview == nil {
		log.Warnf("[CustKycSvc][PickKycReview] no active review found for kyc id: %s", id)
		return errors.New("No active review found for this KYC request")
	}

	userID, err := bson.ObjectIDFromHex(makerUser.UserID)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] failed to fetch user info: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	newReq := &imodel.StartedKycReview{
		PickedBy: &imodel.UserInfo{
			ID:          userID,
			UserCode:    cpsUser.UserCode,
			FullName:    cpsUser.FullName,
			Email:       cpsUser.Email,
			PhoneNumber: cpsUser.PhoneNumber,
		},
		PickReason: reason,
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, nil, newReq, string(constants.RequestPickKycReview), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] cps action err: %v", err)
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
			log.Errorf("Core account creation failed: %v", err)
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
		existingReview, err := s.repo.FindKycInReview(ctx, cpsAction.UniqueId)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		existingReview.ReviewStatus = string(constants.KYCStatusApproved)
		_, err = s.repo.UpdateKycReview(ctx, cpsAction.UniqueId, existingReview)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
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

		existingReview, err := s.repo.FindKycInReview(ctx, cpsAction.UniqueId)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			log.Errorf("[CustKycSvc][Authorize] failed to check existing review: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		existingReview.ReviewStatus = string(constants.KYCStatusRejected)
		_, err = s.repo.UpdateKycReview(ctx, cpsAction.UniqueId, existingReview)
		if err != nil {
			log.Errorf("[CustKycSvc][Authorize] failed to update review status: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
	case string(constants.RequestPickKycReview):
		reviewData, err := local_util.JsonUnmarshal[imodel.StartedKycReview](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}

		// Set a new expiration time for the review to be picked by another reviewer if the current reviewer fails to complete the review in time
		existingReview, err := s.repo.FindKycInReview(ctx, cpsAction.UniqueId)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			log.Errorf("[CustKycSvc][Authorize] failed to check existing review: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		now := time.Now()
		existingReview.PickedAt = &now
		existingReview.StartedAt = now
		existingReview.ExpiresAt = now.Add(30 * time.Minute)
		existingReview.Reviewer = *reviewData.PickedBy
		existingReview.PickedBy = reviewData.PickedBy
		existingReview.PickReason = reviewData.PickReason
		existingReview.PickCount += 1

		_, err = s.repo.UpdateKycReview(ctx, cpsAction.UniqueId, existingReview)
		if err != nil {
			log.Errorf("[CustKycSvc][Authorize] failed to update review expiration: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

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
