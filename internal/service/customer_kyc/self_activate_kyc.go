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

	if result.KYC.KYCStatus == string(imodel.KYCStatusInReview) {
		review, err := s.repo.FindKycInReviewSA(ctx, id)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		if review != nil {
			mappedResponse.KYCReviewStartedAt = review.StartedAt
			mappedResponse.KYCReviewExpiresAt = review.ExpiresAt
			mappedResponse.Reviewer = &review.Reviewer
		}
	}

	return mappedResponse, nil
}

func (s *selfActivationKYCService) EnableOrDisable(ctx context.Context, id, reason string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	if local_util.IsIncomplete(makerUser) {
		log.Errorf("[CustKycSvc][EnableDisable] incomplete maker user")
		return errors.New(constants.IncompleteUserInfo)
	}

	userReq, err := s.repo.FindByIDSA(ctx, id)
	if err != nil {
		log.Errorf("[CustKycSvc][EnableDisable] find err: %v", err)
		return err
	}

	if userReq.KYC.KYCStatus == string(imodel.KYCStatusPending) {
		log.Warnf("[CustKycSvc][EnableDisable] kyc status is pending id: %s, enabled: %v", id, enable)
		return errors.New("Start KYC review before approving or rejecting the KYC request")
	}

	if userReq.KYC.KYCStatus == string(imodel.KYCStatusInReview) {
		review, err := s.repo.FindKycInReviewSA(ctx, id)
		if err != nil {
			log.Errorf("[CustKycSvc][EnableDisable] failed to find kyc review: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		if review != nil && review.ExpiresAt.Before(time.Now()) {
			log.Warnf("[CustKycSvc][EnableDisable] kyc review expired id: %s, enabled: %v", id, enable)
			return errors.New("The KYC review period has expired.")
		}
	}

	if userReq.Enabled && enable && userReq.KYC.KYCStatus == string(imodel.KYCStatusApproved) {
		log.Warnf("[CustKycSvc][EnableDisable] already in state id: %s, enabled: %v", id, enable)
		return errors.New("Customer KYC is already approved")
	}

	if !enable && userReq.KYC.KYCStatus == string(imodel.KYCStatusRejected) {
		log.Warnf("[CustKycSvc][EnableDisable] already rejected id: %s", id)
		return errors.New("Customer KYC is already rejected")
	}

	var action constants.RequestAction
	if enable {
		action = constants.RequestApproveSelfActivateKyc
	} else {
		action = constants.RequestRejectSelfActivateKyc
	}

	newReq := *userReq
	newReq.Enabled = enable
	if enable {
		newReq.KYC.KYCStatus = string(imodel.KYCStatusApproved)
	} else {
		newReq.KYC.KYCStatus = string(imodel.KYCStatusRejected)
		newReq.KYCRejectReason = reason
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, userReq, newReq, string(action), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustKycSvc][EnableDisable] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *selfActivationKYCService) StartKycReview(ctx context.Context, id string) (*imodel.StartedKycReview, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	kycID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

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

	if kyc.KYC.KYCStatus != string(imodel.KYCStatusPending) {
		log.Warnf("[SelfActivationKYC][StartKycReview] kyc status is not pending id: %s, status: %s", id, kyc.KYC.KYCStatus)
		return nil, errors.New("KYC review can only be started for KYC requests with pending status")
	}

	// existingReview, err := s.repo.FindKycInReviewSA(ctx, id)
	// if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
	// 	log.Errorf("[SelfActivationKYC][StartKycReview] failed to check existing review: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }
	// if existingReview != nil &&
	// 	existingReview.ReviewStatus == string(imodel.KYCStatusInReview) {

	// 	now := time.Now()

	// 	switch {
	// 	case existingReview.ExpiresAt == nil || existingReview.ExpiresAt.After(now):
	// 		log.Warnf("[SelfActivationKYC][StartKycReview] active review already exists for kyc id: %s", id)
	// 		return nil, errors.New("An active review already exists for this KYC request")

	// 	default:
	// 		log.Warnf("[SelfActivationKYC][StartKycReview] review already exists but expired for kyc id: %s", id)
	// 		return nil, errors.New("An expired review already exists. Please pick the review to restart the review process.")
	// 	}
	// }

	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] failed to fetch user info: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	now := time.Now()
	expiresAt := time.Now().Add(30 * time.Minute)
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
		StartedAt:      &now,
		ExpiresAt:      &expiresAt,
		IsActive:       true,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	newReview, err := s.repo.StartKycReviewSA(ctx, newReq)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateKYCStatusSA(ctx, id, string(imodel.KYCStatusInReview), "", false)
	if err != nil {
		return nil, err
	}

	return newReview, nil
}

func (s *selfActivationKYCService) PickKycReview(ctx context.Context, id string, reason string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	kycInReview, err := s.repo.FindKycInReviewSA(ctx, id)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New("No started  KYC review found.")
		}
		log.Errorf("[SelfActivationKYC][PickKycReview] find err: %v", err)
		return err
	}

	if kycInReview.ReviewStatus != string(imodel.KYCStatusInReview) {
		return fmt.Errorf("This KYC is already been: %s", kycInReview.ReviewStatus)
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
	if kycInReview.ExpiresAt.After(now) && kycInReview.Reviewer.ID != userID {
		log.Warnf("[SelfActivationKYC][PickKycReview] user %s is not the current reviewer for kyc id: %s", makerUser.UserCode, id)
		return errors.New("This KYC review is currently assigned to another reviewer")
	}

	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[SelfActivationKYC][PickKycReview] failed to fetch user info: %v", err)
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

	// cpsActionData := lib.CpsModelBuilder(id, makerUser, nil, newReq, string(constants.RequestPickSelfActivateKycReview), constants.CREATE)

	// if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
	// 	log.Errorf("[SelfActivationKYC][PickKycReview] cps action err: %v", err)
	// 	return err
	// }

	expiresAt := now.Add(30 * time.Minute)

	kycInReview.PickedAt = &now
	kycInReview.StartedAt = &now
	kycInReview.ExpiresAt = &expiresAt
	kycInReview.Reviewer = *newReq.PickedBy
	kycInReview.PickedBy = newReq.PickedBy
	kycInReview.PickReason = newReq.PickReason
	kycInReview.PickCount += 1

	_, err = s.repo.UpdateKycReviewSA(ctx, id, kycInReview)
	if err != nil {
		log.Errorf("[SelfActivationKYC][PickKycReview] failed to update review expiration: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (s *selfActivationKYCService) ExportUserSelfActivation(ctx context.Context, from, to time.Time, fileType, customerName string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if s.minioClient == nil {
		log.Errorf("[ExportUserSelfActivation] minio client is not configured")
		return "", errors.New(localization.CpsUserDataExportedError.Code)
	}

	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType != string(lib.FileTypeCSV) && fileType != string(lib.FileTypePDF) {
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}

	data, err := s.repo.FindForExport(ctx, from, to, customerName)
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

		rejectionReason := strings.TrimSpace(userData.KYCRejectReason)
		if exists {
			s.logger.Errorf("This user exists and has linked account: %v", err)

			if err = s.repo.UpdateKYCStatusSA(ctx, cpsAction.UniqueId, string(constants.KYCStatusCancelled), rejectionReason, false); err != nil {
				return nil, err
			}

			existingReview, err := s.repo.FindKycInReviewSA(ctx, cpsAction.UniqueId)
			if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
				log.Errorf("[SelfActivationKYC][Authorize] failed to check existing review: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}

			if existingReview != nil {
				existingReview.ReviewStatus = string(constants.KYCStatusCancelled)
				_, err = s.repo.UpdateKycReviewSA(ctx, cpsAction.UniqueId, existingReview)
				if err != nil {
					log.Errorf("[SelfActivationKYC][Authorize] failed to update review status: %v", err)
					return nil, errors.New(localization.ErrorUnexpectedError.Code)
				}
			}

			return nil, fmt.Errorf("The user has already been approved through a branch. This request has been cancelled.")
		} else {

			if err = s.repo.UpdateKYCStatusSA(ctx, cpsAction.UniqueId, string(constants.KYCStatusApproved), rejectionReason, false); err != nil {
				return nil, err
			}

			existingReview, err := s.repo.FindKycInReviewSA(ctx, cpsAction.UniqueId)
			if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
				log.Errorf("[SelfActivationKYC][Authorize] failed to check existing review: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
			if existingReview != nil {
				existingReview.ReviewStatus = string(constants.KYCStatusApproved)
				_, err = s.repo.UpdateKycReviewSA(ctx, cpsAction.UniqueId, existingReview)
				if err != nil {
					log.Errorf("[SelfActivationKYC][Authorize] failed to update review status: %v", err)
					return nil, errors.New(localization.ErrorUnexpectedError.Code)
				}
			}

			if s.smsService != nil {
				phone := userData.PhoneNumber
				name := userData.Name
				// reason := rejectionReason
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

		rejectionReason := strings.TrimSpace(userData.KYCRejectReason)
		if rejectionReason == "" {
			return nil, errors.New("rejection reason is required")
		}

		if err = s.repo.UpdateKYCStatusSA(ctx, cpsAction.UniqueId, string(constants.KYCStatusRejected), rejectionReason, false); err != nil {
			return nil, err
		}

		existingReview, err := s.repo.FindKycInReviewSA(ctx, cpsAction.UniqueId)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			log.Errorf("[SelfActivationKYC][Authorize] failed to check existing review: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		if existingReview != nil {
			existingReview.ReviewStatus = string(constants.KYCStatusRejected)
			_, err = s.repo.UpdateKycReviewSA(ctx, cpsAction.UniqueId, existingReview)
			if err != nil {
				log.Errorf("[SelfActivationKYC][Authorize] failed to update review status: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
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

		// case string(constants.RequestPickSelfActivateKycReview):
		// 	reviewData, err := local_util.JsonUnmarshal[imodel.StartedKycReview](cpsAction.CurrentAction)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	existingReview,
		// 		err := s.repo.FindKycInReviewSA(ctx, cpsAction.UniqueId)
		// 	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		// 		log.Errorf("[SelfActivationKYC][Authorize] failed to check existing review: %v", err)
		// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
		// 	}
		// 	now := time.Now()
		// 	expiresAt := now.Add(30 * time.Minute)

		// 	existingReview.PickedAt = &now
		// 	existingReview.StartedAt = &now
		// 	existingReview.ExpiresAt = &expiresAt
		// 	existingReview.Reviewer = *reviewData.PickedBy
		// 	existingReview.PickedBy = reviewData.PickedBy
		// 	existingReview.PickReason = reviewData.PickReason
		// 	existingReview.PickCount += 1

		// 	_, err = s.repo.UpdateKycReviewSA(ctx, cpsAction.UniqueId, existingReview)
		// 	if err != nil {
		// 		log.Errorf("[SelfActivationKYC][Authorize] failed to update review expiration: %v", err)
		// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
		// 	}
		// 	return cpsAction, nil

	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}

	return cpsAction, nil
}
