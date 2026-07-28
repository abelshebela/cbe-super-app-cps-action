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
	"encoding/csv"
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
		smsService:     smsService,
	}
}

func (r *customerKYCService) GetKYCReviewPeriod() time.Duration {
	period := r.cfg.KYCReviewPeriodInDays

	switch r.cfg.GoEnv {
	case "dev", "qa", "uat":
		return time.Duration(period) * time.Minute
	default:
		return time.Duration(period) * 24 * time.Hour
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
	return core.MapCustomerKYCToResponse(result), nil
}

func (s *customerKYCService) EnableOrDisable(ctx context.Context, id, reason string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	if local_util.IsIncomplete(makerUser) {
		log.Errorf("[CustKycSvc][EnableDisable] incomplete maker user")
		return errors.New(constants.IncompleteUserInfo)
	}

	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustKycSvc][EnableDisable] find err: %v", err)
		return err
	}

	if string(result.Review.Status) == string(imodel.KYCStatusPending) {
		log.Warnf("[CustKycSvc][EnableDisable] kyc status is pending id: %s, enabled: %v", id, enable)
		return errors.New("Start KYC review before approving or rejecting the KYC request")
	}

	if string(result.Review.Status) == string(imodel.KYCStatusInReview) && result.Review.ReviewExpiresAt.Before(time.Now()) {
		log.Warnf("[CustKycSvc][EnableDisable] kyc review expired id: %s, enabled: %v", id, enable)
		return errors.New("The KYC review period has expired.")
	}

	if result.Enabled && enable && string(result.Review.Status) == string(imodel.KYCStatusApproved) {
		log.Warnf("[CustKycSvc][EnableDisable] already in state id: %s, enabled: %v", id, enable)
		return errors.New("Customer KYC is already approved")
	}

	if !enable && string(result.Review.Status) == string(imodel.KYCStatusRejected) {
		log.Warnf("[CustKycSvc][EnableDisable] already rejected id: %s", id)
		return errors.New("Customer KYC is already rejected")
	}

	var action constants.RequestAction
	if enable {
		action = constants.RequestApproveCustomerKYC
	} else {
		action = constants.RequestRejectCustomerKYC
	}

	newReq := *result
	newReq.Enabled = enable
	if enable {
		newReq.Review.Status = imodel.KYCStatusApproved
	} else {
		newReq.Review.Status = imodel.KYCStatusRejected
		newReq.Review.Decision.RejectionReason = reason
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, result, newReq, string(action), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustKycSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *customerKYCService) StartKycReview(ctx context.Context, id string) (*dto.CustomerKYCResponse, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

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

	if kyc.Review == nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] kyc.review is null or status is empty: %v", err)
		return nil, fmt.Errorf("KYC status not found")
	}

	if string(kyc.Review.Status) != string(imodel.KYCStatusPending) {
		log.Warnf("[CustKycSvc][StartKycReview] kyc status is not pending id: %s, status: %s", id, kyc.Review.Status)
		return nil, errors.New("KYC review can only be started for KYC requests with pending status")
	}

	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] failed to fetch user info: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	now := time.Now()
	expiresAt := time.Now().Add(s.GetKYCReviewPeriod())
	newReq := &imodel.CustomerKYC{
		Review: &imodel.KYCReview{
			Status:          imodel.ReviewStatus(imodel.KYCStatusInReview),
			ReviewStartedAt: &now,
			ReviewExpiresAt: &expiresAt,
			Decision: &imodel.KYCReviewDecision{
				Reviewer: &imodel.UserInfo{
					ID:          userID,
					UserCode:    cpsUser.UserCode,
					FullName:    cpsUser.FullName,
					Email:       cpsUser.Email,
					PhoneNumber: cpsUser.PhoneNumber,
				},
			},
		},
		LastModifiedAt: time.Now(),
	}

	result, err := s.repo.UpdateKYC(ctx, id, newReq)
	if err != nil {
		return nil, err
	}

	return core.MapCustomerKYCToResponse(result), nil
}

func (s *customerKYCService) PickKycReview(ctx context.Context, id string, reason string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	kycInReview, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New("No started  KYC review found.")
		}
		log.Errorf("[CustKycSvc][PickKycReview] find err: %v", err)
		return err
	}

	if string(kycInReview.Review.Status) != string(imodel.KYCStatusInReview) {
		return fmt.Errorf("This KYC is already been: %s", kycInReview.ReviewStatus)
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
	if kycInReview.Review.ReviewExpiresAt.After(now) && kycInReview.Review.Decision.Reviewer.ID != userID {
		log.Warnf("[CustKycSvc][PickKycReview] user %s is not the current reviewer for kyc id: %s", makerUser.UserCode, id)
		return errors.New("This KYC review is currently assigned to another reviewer")
	} else if kycInReview.Review.ReviewExpiresAt.After(now) {
		return errors.New("You already picked the review")
	}

	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[CustKycSvc][StartKycReview] failed to fetch user info: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	expiresAt := now.Add(s.GetKYCReviewPeriod())
	assignment := imodel.KYCReviewAssignment{
		PickedAt: &now,
		PickedBy: &imodel.UserInfo{
			ID:          userID,
			UserCode:    cpsUser.UserCode,
			FullName:    cpsUser.FullName,
			Email:       cpsUser.Email,
			PhoneNumber: cpsUser.PhoneNumber,
		},
		PickReason: reason,
	}

	kycInReview.Review.Assignments = append(
		[]imodel.KYCReviewAssignment{assignment},
		kycInReview.Review.Assignments...,
	)

	kycInReview.Review.PickCount++
	kycInReview.Review.ReviewStartedAt = &now
	kycInReview.Review.ReviewExpiresAt = &expiresAt
	kycInReview.Review.Decision.Reviewer = &imodel.UserInfo{
		ID:          userID,
		UserCode:    cpsUser.UserCode,
		FullName:    cpsUser.FullName,
		Email:       cpsUser.Email,
		PhoneNumber: cpsUser.PhoneNumber,
	}

	_, err = s.repo.UpdateKYC(ctx, id, kycInReview)
	if err != nil {
		log.Errorf("[CustKycSvc][PickKycReview] failed to update review expiration: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (s *customerKYCService) ExportKYCOnboarding(ctx context.Context, from, to time.Time, fileType, status, customerName string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if s.minio == nil {
		log.Errorf("[ExportKYCOnboarding] minio client is not configured")
		return "", errors.New(localization.CpsUserDataExportedError.Code)
	}

	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType != string(lib.FileTypeCSV) && fileType != string(lib.FileTypePDF) {
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}

	data, err := s.repo.FindKYCOnboardingForExport(ctx, from, to, status, customerName)
	if err != nil {
		log.Errorf("[ExportKYCOnboarding] failed to fetch onboarding requests: %v", err)
		return "", err
	}

	headers := []string{
		"Customer Name",
		"Phone Number",
		"Gender",
		"Date of Birth",
		"Region",
		"Registration Date & Time",
		"Rejection Reason",
		"Customer Status",
		"KYC Status",
	}

	ext := "csv"
	if fileType == string(lib.FileTypePDF) {
		ext = "pdf"
	}

	objectName := fmt.Sprintf("onboarding_kyc_requests_%s_to_%s_%d.%s", from.Format("20060102"), to.Format("20060102"), time.Now().Unix(), ext)

	if fileType == string(lib.FileTypePDF) {
		rows := make([][]string, 0, len(data))
		for _, item := range data {
			rows = append(rows, core.BuildRow(item))
		}
		url, exportErr := lib.ExportPDFAndUpload(ctx, s.minio, s.bucketName, *s.cfg, objectName, headers, rows, lib.PDFExportOptions{PageSize: "A4"}, nil, s.logger)
		if exportErr != nil {
			log.Errorf("[CPSUser] pdf export failed: %v", exportErr)
			return "", fmt.Errorf("failed to export onboarding kyc data")
		}

		return url, nil
	}

	url, exportErr := lib.ExportCSVAndUpload(ctx, s.minio, s.bucketName, *s.cfg, objectName, headers, func(writer *csv.Writer) error {
		for _, item := range data {
			if err := writer.Write(core.BuildRow(item)); err != nil {
				return err
			}
		}
		return nil
	}, s.logger)
	if exportErr != nil {
		log.Errorf("[CPSUser] csv export failed: %v", exportErr)
		return "", fmt.Errorf("failed to export onboarding kyc data")
	}

	return url, nil
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

		// var firstName, middleName, lastName string

		// nameParts := strings.Fields(strings.TrimSpace(userData.KYCData.FullName))

		// switch len(nameParts) {
		// case 1:
		// 	firstName = nameParts[0]
		// case 2:
		// 	firstName = nameParts[0]
		// 	lastName = nameParts[1]
		// default:
		// 	firstName = nameParts[0]
		// 	middleName = nameParts[1]
		// 	lastName = strings.Join(nameParts[2:], " ")
		// }

		// token, err := s.tokenProvider.GetToken(ctx)
		// if err != nil {
		// 	log.Errorf("[CustKycSvc][Authorize] token fetch: %v", err)
		// 	return nil, err
		// }

		// data := coreio.CreateCustomerParam{
		// 	FirstName:  strings.ToUpper(strings.TrimSpace(firstName)),
		// 	MiddleName: strings.ToUpper(strings.TrimSpace(middleName)),
		// 	LastName:   strings.ToUpper(strings.TrimSpace(lastName)),

		// 	PhoneNumber: userData.KYCData.PhoneNumber,

		// 	Address: strings.TrimSpace(userData.KYCData.Address.Woreda),

		// 	PostalCode:     constants.Empty,
		// 	ISOCountryCode: "ET",

		// 	AccountOffice: s.cfg.CentralKYCBranchCode,
		// 	Industry:      constants.Empty,

		// 	ISONationalityCode: "ET",
		// 	ISOResidentCode:    "ET",

		// 	UniqueID: t24LegalID(local_util.NonEmptyString(userData.KYCData.OriginID, userData.KYCData.Sub)),
		// 	IssuesBy: strings.ToUpper(string(userData.KYCData.Vendor)),
		// 	// IssuedDate: t24IssuedDate(userData.KYCData.IssuedDate),
		// 	ExpiryDate: constants.Empty,

		// 	Gender:      strings.ToUpper(strings.TrimSpace(userData.KYCData.Gender)),
		// 	DateOfBirth: userData.KYCData.BirthDate.Format("20060102"),

		// 	MaritalStatus: strings.ToUpper(strings.TrimSpace(userData.KYCData.MaritalStatus)),
		// 	Email:         userData.KYCData.Email,

		// 	EmploymentStatus: t24EmploymentStatus(userData.KYCData.EmployementStatus),
		// 	Occupation:       strings.ToUpper(strings.TrimSpace(userData.KYCData.Occupation)),

		// 	EmployerName:     constants.Empty,
		// 	EmployerAddress:  constants.Empty,
		// 	EmployerBusiness: userData.KYCData.SourceOfIncome,

		// 	CustomerCurrency:  t24Currency(userData.KYCData.Currency),
		// 	Salary:            t24Amount(userData.KYCData.MonthlyIncome),
		// 	AnnualBonus:       constants.Empty,
		// 	NetMonthlyIncome:  t24Amount(userData.KYCData.MonthlyIncome),
		// 	NetMonthlyExpence: constants.Empty,

		// 	TinNumber:     userData.KYCData.USTIN,
		// 	MotherName:    strings.ToUpper(strings.TrimSpace(userData.KYCData.MothersName)),
		// 	CustomerGroup: t24CustomerGroup(userData.KYCData.SubAccountType),
		// 	NationalId:    userData.KYCData.Sub,

		// 	Url: s.cfg.CustomerCreateURL,
		// 	Header: map[string]string{
		// 		"Authorization": "Bearer " + token,
		// 	},
		// }

		// s.logger.Infof("[CustKycSvc][Authorize] core call params — uniqueID=%s employmentStatus=%s salary=%s nationality=%s country=%s", data.UniqueID, data.EmploymentStatus, data.Salary, data.ISONationalityCode, data.ISOCountryCode)
		// userAccount, err := core.CreateAccountToCore(ctx, data, s.accountService, s.coreio, s.cfg, s.logger)
		// if err != nil {
		// 	log.Errorf("[CustKycSvc][Authorize] core account creation failed: %v", err)
		// 	if strings.Contains(err.Error(), "customer creation failed") {
		// 		return nil, errors.New(localization.ErrorCustomerCreationOnCoreFailed.Code)
		// 	}
		// 	return nil, err
		// }

		// if err = s.repo.CreateUser(ctx, userAccount, *userData); err != nil {
		// 	return nil, err
		// }

		err = s.repo.ApproveOrReject(ctx, cpsAction.UniqueId, string(imodel.KYCStatusApproved), "", true)
		if err != nil {
			return nil, err
		}

		if s.smsService != nil {
			phone := userData.KYCData.PhoneNumber
			name := userData.KYCData.FullName
			go func() {
				msg := fmt.Sprintf(
					"Dear %s, Congratulations! Your application for opening a new account and superapp activation is successful, Please Check and continue your account creation.",
					name,
				)
				if err := s.smsService.PublishSMSMessage(context.Background(), types.SMSKafkaMessage{
					Recipient:   phone,
					MessageBody: msg,
				}); err != nil {
					s.logger.Errorf("[CustKycSvc][Authorize] SMS send failed for phone %s: %v", phone, err)
				}
			}()
		}

	case string(constants.RequestRejectCustomerKYC):
		userData, err := local_util.JsonUnmarshal[imodel.CustomerKYC](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}

		rejectionReason := strings.TrimSpace(userData.Review.Decision.RejectionReason)
		if rejectionReason == "" {
			return nil, errors.New("rejection reason is required")
		}

		err = s.repo.ApproveOrReject(ctx, cpsAction.UniqueId, string(imodel.KYCStatusRejected), rejectionReason, false)
		if err != nil {
			return nil, err
		}

		if s.smsService != nil {
			phone := userData.KYCData.PhoneNumber
			name := userData.KYCData.FullName
			reason := rejectionReason
			go func() {
				msg := fmt.Sprintf(
					"Dear %s, Your application for opening a new account and CBE superapp activation is Rejected Due to %s, please correct and apply again. Thank You",
					name, reason,
				)
				if err := s.smsService.PublishSMSMessage(context.Background(), types.SMSKafkaMessage{
					Recipient:   phone,
					MessageBody: msg,
				}); err != nil {
					s.logger.Errorf("[CustKycSvc][Authorize] SMS send failed for phone %s: %v", phone, err)
				}
			}()
		}

	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}

	return cpsAction, nil
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

func t24IssuedDate(issuedDate string) string {
	if strings.TrimSpace(issuedDate) != "" {
		return issuedDate
	}
	return "20210101"
}

func t24LegalID(id string) string {
	const maxLen = 35
	if len(id) <= maxLen {
		return id
	}
	return id[:maxLen]
}
