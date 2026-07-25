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

type selfActivationKYCService struct {
	repo           storage.SelfActivationKYCRepository
	cpsService     service.CPSActionService
	customerSvc    service.CustomerService
	accountService account_lookup.Account
	cpsUserRepo    storage.CpsUserRepository
	coreio         coreio.CBECoreAPIInterface
	tokenProvider  service.TokenProviderService
	logger         utils.Logger
	bucketName     string
	cfg            *config.VaultConfig
	smsService     *lib.NotificationStore
	minioClient    *s3.Client
}

func NewSelfActivationKYCService(repo storage.SelfActivationKYCRepository,
	cpsService service.CPSActionService,
	cpsUserRepo storage.CpsUserRepository,
	customerSvc service.CustomerService,
	accountLookUpService account_lookup.Account,
	coreio coreio.CBECoreAPIInterface,
	tokenProvider service.TokenProviderService,
	logger utils.Logger,
	bucketName string,
	cfg *config.VaultConfig,
	smsService *lib.NotificationStore,
	minioClient *s3.Client,
) service.SelfActivateKYCService {
	return &selfActivationKYCService{
		repo:           repo,
		cpsService:     cpsService,
		cpsUserRepo:    cpsUserRepo,
		customerSvc:    customerSvc,
		accountService: accountLookUpService,
		coreio:         coreio,
		tokenProvider:  tokenProvider,
		logger:         logger,
		bucketName:     bucketName,
		cfg:            cfg,
		smsService:     smsService,
		minioClient:    minioClient,
	}
}

func (s *selfActivationKYCService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]dto.CustomerKYCResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "SelfActivationKYC", "FindAllWithPagination")
	defer span.End()

	result, err := s.repo.FindAllWithPaginationSA(ctx, *filterParam)
	if err != nil {
		return nil, err
	}

	return core.MapSelfActivationUserToResponsePaginated(result), nil
}

func (s *selfActivationKYCService) FindByID(ctx context.Context, id string) (*dto.CustomerKYCResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "SelfActivationKYC", "FindByID")
	defer span.End()

	result, err := s.repo.FindByIDSA(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedResponse := core.MapSelfActivationUserToResponse(result)

	// userData, err := s.customerSvc.SearchCustomerByCIF(ctx, result.CustomerNumber)
	// if err != nil {
	// 	return nil, err
	// }

	// chosenAccounts := make(map[string]struct{}, len(result.ChosenAccounts))
	// for _, account := range result.ChosenAccounts {
	// 	chosenAccounts[account] = struct{}{}
	// }

	// for _, user := range userData {
	// 	if _, ok := chosenAccounts[user.AccountName]; ok {
	// 		mappedResponse.LinkedAccount = append(mappedResponse.LinkedAccount, dto.LinkedAccounts{
	// 			AccountHolderName:  user.CustomerName,
	// 			PhoneNumber:        user.Phone,
	// 			AccountType:        user.AccountType,
	// 			ProductCode:        "-",
	// 			LinkedChannel:      "-",
	// 			Currency:           user.Currency,
	// 			BranchName:         user.BranchName,
	// 			BranchCode:         user.BranchCode,
	// 			InActive:           user.InactiveFlag,
	// 			PostingRestriction: "-",
	// 			RestrictionType:    user.RestrictionType,
	// 			AndorAccount:       nil,
	// 		})
	// 	}
	// }

	return mappedResponse, nil
}

func (s *selfActivationKYCService) EnableOrDisable(ctx context.Context, id, reason string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	if local_util.IsIncomplete(makerUser) {
		log.Errorf("[CustKycSvc][EnableDisable] incomplete maker user")
		return errors.New(constants.IncompleteUserInfo)
	}

	result, err := s.repo.FindByIDSA(ctx, id)
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
		return errors.New("Self activation KYC is already approved")
	}

	if !enable && string(result.Review.Status) == string(imodel.KYCStatusRejected) {
		log.Warnf("[CustKycSvc][EnableDisable] already rejected id: %s", id)
		return errors.New("Self activation KYC is already rejected")
	}

	var action constants.RequestAction
	if enable {
		action = constants.RequestApproveSelfActivateKyc
	} else {
		action = constants.RequestRejectSelfActivateKyc
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

func (s *selfActivationKYCService) StartKycReview(ctx context.Context, id string) (*dto.CustomerKYCResponse, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	userID, err := bson.ObjectIDFromHex(makerUser.UserID)
	if err != nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	kyc, err := s.repo.FindByIDSA(ctx, id)
	if err != nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] find err: %v", err)
		return nil, err
	}

	if kyc.Review == nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] kyc.review is null or status is empty: %v", err)
		return nil, fmt.Errorf("KYC status not found")
	}

	if string(kyc.Review.Status) != string(imodel.KYCStatusPending) {
		log.Warnf("[SelfActivationKYC][StartKycReview] kyc status is not pending id: %s, status: %s", id, kyc.Review.Status)
		return nil, errors.New("KYC review can only be started for KYC requests with pending status")
	}

	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] failed to fetch user info: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	now := time.Now()
	expiresAt := time.Now().Add(30 * time.Minute)
	newReq := &imodel.SelfActivationUser{
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

	result, err := s.repo.UpdateKYCSA(ctx, id, newReq)
	if err != nil {
		return nil, err
	}

	return core.MapSelfActivationUserToResponse(result), nil
}

func (s *selfActivationKYCService) PickKycReview(ctx context.Context, id string, reason string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	kycInReview, err := s.repo.FindByIDSA(ctx, id)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New("No started  KYC review found.")
		}
		log.Errorf("[SelfActivationKYC][PickKycReview] find err: %v", err)
		return err
	}

	if string(kycInReview.Review.Status) != string(imodel.KYCStatusInReview) {
		return fmt.Errorf("This KYC is already been: %s", kycInReview.Review.Status)
	}

	if kycInReview == nil {
		log.Warnf("[SelfActivationKYC][PickKycReview] no active review found for kyc id: %s", id)
		return errors.New("No active review found for this KYC request")
	}

	userID, err := bson.ObjectIDFromHex(makerUser.UserID)
	if err != nil {
		log.Errorf("[SelfActivationKYC][PickKycReview] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	now := time.Now()
	if kycInReview.Review.ReviewExpiresAt.After(now) && kycInReview.Review.Decision.Reviewer.ID != userID {
		log.Warnf("[SelfActivationKYC][PickKycReview] user %s is not the current reviewer for kyc id: %s", makerUser.UserCode, id)
		return errors.New("This KYC review is currently assigned to another reviewer")
	} else if kycInReview.Review.ReviewExpiresAt.After(now) {
		return errors.New("You already picked the review")
	}

	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[SelfActivationKYC][PickKycReview] failed to fetch user info: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	expiresAt := now.Add(30 * time.Minute)
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

	_, err = s.repo.UpdateKYCSA(ctx, id, kycInReview)
	if err != nil {
		log.Errorf("[SelfActivationKYC][PickKycReview] failed to update review expiration: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// cpsActionData := lib.CpsModelBuilder(id, makerUser, nil, newReq, string(constants.RequestPickSelfActivateKycReview), constants.CREATE)

	// if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
	// 	log.Errorf("[SelfActivationKYC][PickKycReview] cps action err: %v", err)
	// 	return err
	// }

	return nil
}

func (s *selfActivationKYCService) ExportUserSelfActivation(ctx context.Context, from, to time.Time, fileType, status, customerName string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if s.minioClient == nil {
		log.Errorf("[ExportUserSelfActivation] minio client is not configured")
		return "", errors.New(localization.CpsUserDataExportedError.Code)
	}

	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType != string(lib.FileTypeCSV) && fileType != string(lib.FileTypePDF) {
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}

	data, err := s.repo.FindForExport(ctx, from, to, status, customerName)
	if err != nil {
		log.Errorf("[ExportUserSelfActivation] failed to fetch self activation requests: %v", err)
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

	objectName := fmt.Sprintf("self_activation_requests%s_to_%s_%d.%s", from.Format("20060102"), to.Format("20060102"), time.Now().Unix(), ext)

	if fileType == string(lib.FileTypePDF) {
		rows := make([][]string, 0, len(data))
		for _, item := range data {
			rows = append(rows, core.BuildRow(item))
		}
		url, exportErr := lib.ExportPDFAndUpload(ctx, s.minioClient, s.bucketName, *s.cfg, objectName, headers, rows, lib.PDFExportOptions{PageSize: "A4"}, nil, s.logger)
		if exportErr != nil {
			log.Errorf("[CPSUser] pdf export failed: %v", exportErr)
			return "", fmt.Errorf("failed to export self-activation kyc data")
		}

		return url, nil
	}

	url, exportErr := lib.ExportCSVAndUpload(ctx, s.minioClient, s.bucketName, *s.cfg, objectName, headers, func(writer *csv.Writer) error {
		for _, item := range data {
			if err := writer.Write(core.BuildRow(item)); err != nil {
				return err
			}
		}
		return nil
	}, s.logger)
	if exportErr != nil {
		log.Errorf("[CPSUser] csv export failed: %v", exportErr)
		return "", fmt.Errorf("failed to export self-activation kyc data")
	}

	return url, nil
}

func (s *selfActivationKYCService) GetUsersActionLog(ctx context.Context, customerNumber string) (*types.PaginatedResponse[[]imodel.SelfActivationActionLog], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUsersActionLog", "SelfActivationKYC", "GetUsersActionLog")
	defer span.End()

	return s.repo.GetUsersActionLog(ctx, customerNumber)
}

func (s *selfActivationKYCService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "SelfActivationKYC", "Authorize")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[SelfActivationKYC][Authorize] action: %s", cpsAction.RequestAction)

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestApproveSelfActivateKyc):
		userData, err := local_util.JsonUnmarshal[imodel.SelfActivationUser](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}

		exists, err := s.repo.CheckIfUserOrAccountExists(ctx, userData)
		if exists {
			s.logger.Errorf("This user exists and has linked account: %v", err)

			err := s.repo.ApproveOrRejectSA(ctx, cpsAction.UniqueId, string(imodel.KYCStatusCancelled), "", false)
			if err != nil {
				return nil, err
			}

			return nil, fmt.Errorf("The user has already been approved through a branch. This request has been cancelled.")
		} else {

			err := s.repo.ApproveOrRejectSA(ctx, cpsAction.UniqueId, string(imodel.KYCStatusApproved), "", true)
			if err != nil {
				return nil, err
			}

			if s.smsService != nil {
				phone := userData.PhoneNumber
				name := userData.Name
				go func() {
					msg := fmt.Sprintf(
						"Dear %s, Congratulations! Your application for superapp activation is successful. Welcome to CBE Super App!",
						name,
					)
					if err := s.smsService.PublishSMSMessage(context.Background(), types.SMSKafkaMessage{
						Recipient:   phone,
						MessageBody: msg,
					}); err != nil {
						s.logger.Errorf("[SelfActivationKYC][Authorize] SMS send failed for phone %s: %v", phone, err)
					}
				}()
			}
		}

	case string(constants.RequestRejectSelfActivateKyc):
		userData, err := local_util.JsonUnmarshal[imodel.SelfActivationUser](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}

		rejectionReason := strings.TrimSpace(userData.Review.Decision.RejectionReason)
		if rejectionReason == "" {
			return nil, errors.New("rejection reason is required")
		}

		err = s.repo.ApproveOrRejectSA(ctx, cpsAction.UniqueId, string(imodel.KYCStatusRejected), rejectionReason, false)
		if err != nil {
			return nil, err
		}

		if s.smsService != nil {
			phone := userData.PhoneNumber
			name := userData.Name
			reason := rejectionReason
			go func() {
				msg := fmt.Sprintf(
					"Dear %s, Your application for CBE superapp activation is Rejected Due to %s, please correct and apply again. Thank You",
					name, reason,
				)
				if err := s.smsService.PublishSMSMessage(context.Background(), types.SMSKafkaMessage{
					Recipient:   phone,
					MessageBody: msg,
				}); err != nil {
					s.logger.Errorf("[SelfActivationKYC][Authorize] SMS send failed for phone %s: %v", phone, err)
				}
			}()
		}

	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}

	return cpsAction, nil
}
