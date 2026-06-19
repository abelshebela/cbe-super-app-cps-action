package customer

import (
	"cbe-super-app-cps-action/internal/constants"
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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	coreio "github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type customerKYCService struct {
	repo           storage.CustomerKYCRepository
	cpsService     service.CPSActionService
	accountService account_lookup.Account
	cpsUserRepo    storage.CpsUserRepository
	coreio         coreio.CBECoreAPIInterface
	tokenProvider  service.TokenProviderService
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
	coreio coreio.CBECoreAPIInterface,
	tokenProvider service.TokenProviderService,
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
		coreio:         coreio,
		tokenProvider:  tokenProvider,
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

	// If the kyc is in review status, get the start, and expiry time
	if result.KYCStatus == imodel.KYCStatusInReview {
		review, err := s.repo.FindKycInReview(ctx, id)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		if review != nil {
			mappedResponse.KYCReviewStartedAt = &review.StartedAt
			mappedResponse.KYCReviewExpiresAt = &review.ExpiresAt
			mappedResponse.Reviewer = &review.Reviewer
		}
	}

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

func (s *customerKYCService) StartKycReview(ctx context.Context, id string) (*imodel.StartedKycReview, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

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

	kyc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] find err: %v", err)
		return nil, err
	}

	if kyc.KYCStatus != imodel.KYCStatusPending {
		log.Warnf("[CustKycSvc][StartKycReview] kyc status is not pending id: %s, status: %s", id, kyc.KYCStatus)
		return nil, errors.New("KYC review can only be started for KYC requests with pending status")
	}

	existingReview, err := s.repo.FindKycInReview(ctx, id)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		log.Errorf("[CustKycSvc][StartKycReview] failed to check existing review: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if existingReview != nil && existingReview.IsActive && existingReview.ExpiresAt.After(time.Now()) {
		log.Warnf("[CustKycSvc][StartKycReview] active review already exists for kyc id: %s", id)
		return nil, errors.New("An active review already exists for this KYC request")
	}

	if existingReview != nil && existingReview.ExpiresAt.Before(time.Now()) {
		log.Warnf("[CustKycSvc][StartKycReview] review already exists but expired for kyc id: %s", id)
		return nil, errors.New("An expired review already exists. Please pick the review to restart the review process.")
	}

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

	now := time.Now()
	if kycInReview.ExpiresAt.After(now) && kycInReview.Reviewer.ID != userID {
		log.Warnf("[CustKycSvc][PickKycReview] user %s is not the current reviewer for kyc id: %s", makerUser.UserCode, id)
		return errors.New("This KYC review is currently assigned to another reviewer")
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

		var firstName, middleName, lastName string

		nameParts := strings.Fields(strings.TrimSpace(userData.KYCData.FullName))

		switch len(nameParts) {
		case 1:
			firstName = nameParts[0]
		case 2:
			firstName = nameParts[0]
			lastName = nameParts[1]
		default:
			firstName = nameParts[0]
			middleName = nameParts[1]
			lastName = strings.Join(nameParts[2:], " ")
		}

		data := coreio.CreateCustomerParam{
			FirstName:  firstName,
			MiddleName: middleName,
			LastName:   lastName,

			PhoneNumber: userData.KYCData.PhoneNumber,

			Address: strings.Join([]string{
				"Region: " + userData.KYCData.Address.Region,
				"Zone: " + userData.KYCData.Address.Zone,
				"Kebele: " + userData.KYCData.Address.Kebele,
				"Woreda: " + userData.KYCData.Address.Woreda,
			}, ", "),

			PostalCode:     constants.Empty,
			ISOCountryCode: userData.KYCData.Country,

			AccountOffice: constants.Empty,
			Industry:      constants.Empty,

			ISONationalityCode: userData.KYCData.Nationality,
			ISOResidentCode:    userData.KYCData.Country,

			UniqueID:   userData.KYCData.OriginID,
			IssuesBy:   constants.Empty,
			IssuedDate: constants.Empty,
			ExpiryDate: constants.Empty,

			Gender:      userData.KYCData.Gender,
			DateOfBirth: userData.KYCData.BirthDate.Format("2006-01-02"),

			MaritalStatus: userData.KYCData.MaritalStatus,
			Email:         userData.KYCData.Email,

			EmploymentStatus: userData.KYCData.EmployementStatus,
			Occupation:       userData.KYCData.Occupation,

			EmployerName:     constants.Empty,
			EmployerAddress:  constants.Empty,
			EmployerBusiness: userData.KYCData.SourceOfIncome,

			CustomerCurrency:  userData.KYCData.Currency,
			Salary:            userData.KYCData.MonthlyIncome,
			AnnualBonus:       constants.Empty,
			NetMonthlyIncome:  userData.KYCData.MonthlyIncome,
			NetMonthlyExpence: constants.Empty,

			TinNumber:     userData.KYCData.USTIN,
			MotherName:    userData.KYCData.MothersName,
			CustomerGroup: string(constants.MASS),
		}

		s.logger.Infof("[CustKycSvc][Authorize] creating core account for customer: %s", data.UniqueID)
		userAccount, err := core.CreateAccountToCore(ctx, data, s.accountService, s.coreio, s.logger)
		if err != nil {
			log.Errorf("[CustKycSvc][Authorize] core account creation failed: %v", err)
			return nil, err
		}

		log.Infof("[CustKycSvc][Authorize] core account created successfully for customer: %s, account number: %s", data.UniqueID, userAccount)
		token, err := s.tokenProvider.GetToken(ctx)
		if err != nil {
			log.Errorf("[CustKycSvc][Authorize] token fetch: %v", err)
			return nil, err
		}

		log.Infof("[CustKycSvc][Authorize] token fetched successfully for account creation %v", token)
		accountResult, err := s.coreio.AccountCreation(ctx, coreio.AccountCreationParam{
			Username:       s.cfg.CbeCoreUsername,
			Password:       s.cfg.CbeCorePassword,
			CustomerNumber: userAccount.Detail.Customer,
			Category:       userData.KYCData.AccountType,
			Currency:       userData.KYCData.Currency,
			Url:            s.cfg.AccountOpeningURL,
			Header: map[string]string{
				"Authorization": "Bearer " + token,
			},
		})

		if err != nil {
			log.Errorf("[CustKycSvc][Authorize] account creation: %v", err)
			return nil, err
		}

		log.Infof("[CustKycSvc][Authorize] account creation response *****************: %v", accountResult)
		if !accountResult.Success {
			log.Errorf("[CustKycSvc][Authorize] account creation failed: %v", accountResult.Messages)
			return nil, fmt.Errorf("account creation failed: %s", strings.Join(accountResult.Messages, ", "))
		}
		log.Infof("[CustKycSvc][Authorize] account created for customer=%s", userAccount.Detail.Customer)

		if err = s.repo.CreateUser(ctx, userAccount, *userData); err != nil {
			return nil, err
		}

		if err = s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, string(constants.KYCStatusApproved), "", true); err != nil {
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

		rejectionReason := strings.TrimSpace(userData.KYCRejectReason)
		if rejectionReason == "" {
			return nil, errors.New("rejection reason is required")
		}

		if err = s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, string(constants.KYCStatusRejected), rejectionReason, false); err != nil {
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

	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}

	return cpsAction, nil
}
