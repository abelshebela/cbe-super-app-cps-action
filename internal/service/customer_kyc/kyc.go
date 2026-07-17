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
	smsService     *lib.NotificationStore
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
	smsService *lib.NotificationStore,
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
		smsService:     smsService,
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
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New("No started  KYC review found.")
		}
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
		return s.authorizeApprove(ctx, cpsAction)
	case string(constants.RequestRejectCustomerKYC):
		return s.authorizeReject(ctx, cpsAction)
	case string(constants.RequestPickKycReview):
		return s.authorizePickReview(ctx, cpsAction)
	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}
}

func (s *customerKYCService) authorizeApprove(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	userData, err := local_util.JsonUnmarshal[imodel.CustomerKYC](cpsAction.CurrentAction)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenProvider.GetToken(ctx)
	if err != nil {
		log.Errorf("[CustKycSvc][authorizeApprove] token fetch: %v", err)
		return nil, err
	}

	data := s.buildT24CustomerParam(userData, token)
	log.Infof("[CustKycSvc][authorizeApprove] core call params — uniqueID=%s employmentStatus=%s salary=%s nationality=%s country=%s",
		data.UniqueID, data.EmploymentStatus, data.Salary, data.ISONationalityCode, data.ISOCountryCode)

	userAccount, err := core.CreateAccountToCore(ctx, data, s.accountService, s.coreio, s.cfg, s.logger)
	if err != nil {
		log.Errorf("[CustKycSvc][authorizeApprove] core account creation failed: %v", err)
		if strings.Contains(err.Error(), "customer creation failed") {
			return nil, errors.New(localization.ErrorCustomerCreationOnCoreFailed.Code)
		}
		return nil, err
	}

	if err = s.repo.CreateUser(ctx, userAccount, *userData); err != nil {
		return nil, err
	}

	if err = s.repo.UpdateKYCStatus(ctx, cpsAction.UniqueId, string(constants.KYCStatusApproved), "", true); err != nil {
		return nil, err
	}

	review, err := s.findAndValidateReview(ctx, cpsAction.UniqueId, log)
	if err != nil {
		return nil, err
	}

	review.ReviewStatus = string(constants.KYCStatusApproved)
	if _, err = s.repo.UpdateKycReview(ctx, cpsAction.UniqueId, review); err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.sendApprovalSMS(userData.KYCData.PhoneNumber, userData.KYCData.FullName,
		userAccount.AccountCreationDetail.Detail.AccountNumber)

	return cpsAction, nil
}

func (s *customerKYCService) authorizeReject(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

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

	review, err := s.findAndValidateReview(ctx, cpsAction.UniqueId, log)
	if err != nil {
		return nil, err
	}

	review.ReviewStatus = string(constants.KYCStatusRejected)
	if _, err = s.repo.UpdateKycReview(ctx, cpsAction.UniqueId, review); err != nil {
		log.Errorf("[CustKycSvc][authorizeReject] failed to update review status: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.sendRejectionSMS(userData.KYCData.PhoneNumber, userData.KYCData.FullName, rejectionReason)

	return cpsAction, nil
}

func (s *customerKYCService) authorizePickReview(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	reviewData, err := local_util.JsonUnmarshal[imodel.StartedKycReview](cpsAction.CurrentAction)
	if err != nil {
		return nil, err
	}

	review, err := s.findAndValidateReview(ctx, cpsAction.UniqueId, log)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	review.PickedAt = &now
	review.StartedAt = now
	review.ExpiresAt = now.Add(30 * time.Minute)
	review.Reviewer = *reviewData.PickedBy
	review.PickedBy = reviewData.PickedBy
	review.PickReason = reviewData.PickReason
	review.PickCount += 1

	if _, err = s.repo.UpdateKycReview(ctx, cpsAction.UniqueId, review); err != nil {
		log.Errorf("[CustKycSvc][authorizePickReview] failed to update review expiration: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return cpsAction, nil
}

func (s *customerKYCService) findAndValidateReview(ctx context.Context, kycID string, log utils.Logger) (*imodel.StartedKycReview, error) {
	review, err := s.repo.FindKycInReview(ctx, kycID)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		log.Errorf("[CustKycSvc][findAndValidateReview] FindKycInReview failed for kyc id %s: %v", kycID, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if review == nil {
		log.Warnf("[CustKycSvc][findAndValidateReview] no active review found for kyc id: %s", kycID)
		return nil, errors.New(localization.ErrorNoExistingReview.Code)
	}
	return review, nil
}

func (s *customerKYCService) buildT24CustomerParam(userData *imodel.CustomerKYC, token string) coreio.CreateCustomerParam {
	first, middle, last := splitFullName(userData.KYCData.FullName)
	return coreio.CreateCustomerParam{
		FirstName:  strings.ToUpper(strings.TrimSpace(first)),
		MiddleName: strings.ToUpper(strings.TrimSpace(middle)),
		LastName:   strings.ToUpper(strings.TrimSpace(last)),

		PhoneNumber: userData.KYCData.PhoneNumber,
		Address:     strings.TrimSpace(userData.KYCData.Address.Woreda),

		PostalCode:     constants.Empty,
		ISOCountryCode: "ET",

		AccountOffice: s.cfg.CentralKYCBranchCode,
		Industry:      constants.Empty,

		ISONationalityCode: "ET",
		ISOResidentCode:    "ET",

		UniqueID: t24LegalID(local_util.NonEmptyString(userData.KYCData.OriginID, userData.KYCData.Sub)),
		IssuesBy: strings.ToUpper(string(userData.KYCData.Vendor)),
		// IssuedDate: t24IssuedDate(userData.KYCData.IssuedDate),
		ExpiryDate: constants.Empty,

		Gender:      strings.ToUpper(strings.TrimSpace(userData.KYCData.Gender)),
		DateOfBirth: userData.KYCData.BirthDate.Format("20060102"),

		MaritalStatus: strings.ToUpper(strings.TrimSpace(userData.KYCData.MaritalStatus)),
		Email:         userData.KYCData.Email,

		EmploymentStatus: t24EmploymentStatus(userData.KYCData.EmployementStatus),
		Occupation:       strings.ToUpper(strings.TrimSpace(userData.KYCData.Occupation)),

		EmployerName:     constants.Empty,
		EmployerAddress:  constants.Empty,
		EmployerBusiness: userData.KYCData.SourceOfIncome,

		CustomerCurrency:  t24Currency(userData.KYCData.Currency),
		Salary:            t24Amount(userData.KYCData.MonthlyIncome),
		AnnualBonus:       constants.Empty,
		NetMonthlyIncome:  t24Amount(userData.KYCData.MonthlyIncome),
		NetMonthlyExpence: constants.Empty,

		TinNumber:     userData.KYCData.USTIN,
		MotherName:    strings.ToUpper(strings.TrimSpace(userData.KYCData.MothersName)),
		CustomerGroup: t24CustomerGroup(userData.KYCData.SubAccountType),
		NationalId:    userData.KYCData.Sub,

		Url:    s.cfg.CustomerCreateURL,
		Header: map[string]string{"Authorization": "Bearer " + token},
	}
}

func (s *customerKYCService) sendApprovalSMS(phone, name, accountNumber string) {
	if s.smsService == nil {
		return
	}
	go func() {
		msg := fmt.Sprintf(
			"Dear %s, Congratulations! Your application for opening a new account and superapp activation is successful, Your new account number is %s. Welcome to CBE Super App!",
			name, accountNumber,
		)
		if err := s.smsService.PublishSMSMessage(context.Background(), types.SMSKafkaMessage{
			Recipient:   phone,
			MessageBody: msg,
		}); err != nil {
			s.logger.Errorf("[CustKycSvc] approval SMS send failed for phone %s: %v", phone, err)
		}
	}()
}

func (s *customerKYCService) sendRejectionSMS(phone, name, reason string) {
	if s.smsService == nil {
		return
	}
	go func() {
		msg := fmt.Sprintf(
			"Dear %s, Your application for opening a new account and CBE superapp activation is Rejected Due to %s, please correct and apply again. Thank You",
			name, reason,
		)
		if err := s.smsService.PublishSMSMessage(context.Background(), types.SMSKafkaMessage{
			Recipient:   phone,
			MessageBody: msg,
		}); err != nil {
			s.logger.Errorf("[CustKycSvc] rejection SMS send failed for phone %s: %v", phone, err)
		}
	}()
}

func splitFullName(fullName string) (first, middle, last string) {
	parts := strings.Fields(strings.TrimSpace(fullName))
	switch len(parts) {
	case 0:
		return "", "", ""
	case 1:
		return parts[0], "", ""
	case 2:
		return parts[0], "", parts[1]
	default:
		return parts[0], parts[1], strings.Join(parts[2:], " ")
	}
}

func t24CustomerGroup(subAccountType string) string {
	if g := strings.ToUpper(strings.TrimSpace(subAccountType)); g != "" {
		return g
	}
	return "RETAIL"
}

func t24Amount(v string) string {
	if strings.TrimSpace(v) == "" {
		return "0"
	}
	return v
}

func t24Currency(currency string) string {
	if strings.TrimSpace(currency) == "" {
		return "ETB"
	}
	return strings.ToUpper(strings.TrimSpace(currency))
}

func t24EmploymentStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "FULL_TIME", "FULLTIME", "FULL-TIME":
		return "EMPLOYED"
	case "SELF_EMPLOYED", "SELFEMPLOYED", "SELF-EMPLOYED":
		return "SELF-EMP"
	case "PART_TIME", "PARTTIME", "PART-TIME":
		return "EMPLOYED"
	case "UNEMPLOYED":
		return "UNEMPL"
	case "RETIRED":
		return "RETIRED"
	case "STUDENT":
		return "STUDENT"
	default:
		return strings.ToUpper(strings.TrimSpace(status))
	}
}

// t24IssuedDate returns the Fayda ID issuance date for T24's LEGAL.ISS.DATE field.
// If the stored IssuedDate is empty (older KYC records), falls back to 2 years before today
// so it is always in the past relative to T24's UAT/production business date.
func t24IssuedDate(issuedDate string) string {
	if strings.TrimSpace(issuedDate) != "" {
		return issuedDate
	}
	return "20210101"
}

// func t24IssuedDate(issuedDate string) string {
// 	if strings.TrimSpace(issuedDate) != "" {
// 		return issuedDate
// 	}
// 	return time.Now().AddDate(-2, 0, 0).Format("20060102")
// }

// t24LegalID truncates the ID to T24's LEGAL.ID max of 35 characters.
func t24LegalID(id string) string {
	const maxLen = 35
	if len(id) <= maxLen {
		return id
	}
	return id[:maxLen]
}
