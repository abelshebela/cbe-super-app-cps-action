package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	cust_kyc_dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
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

	coreio "github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type selfActivationKYCService struct {
	repo           storage.CustomerKYCRepository
	cpsService     service.CPSActionService
	accountService account_lookup.Account
	cpsUserRepo    storage.CpsUserRepository
	coreio         coreio.CBECoreAPIInterface
	tokenProvider  service.TokenProviderService
	logger         utils.Logger
	bucketName     string
	cfg            *config.VaultConfig
	smsService     *lib.NotificationStore
}

func NewSelfActivationKYCService(repo storage.CustomerKYCRepository,
	cpsService service.CPSActionService,
	cpsUserRepo storage.CpsUserRepository,
	accountLookUpService account_lookup.Account,
	coreio coreio.CBECoreAPIInterface,
	tokenProvider service.TokenProviderService,
	logger utils.Logger,
	bucketName string,
	cfg *config.VaultConfig,
	smsService *lib.NotificationStore,
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
	}
}

func (s *selfActivationKYCService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]cust_kyc_dto.CustomerKYCResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "SelfActivationKYC", "FindAllWithPagination")
	defer span.End()

	result, err := s.repo.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		return nil, err
	}

	return core.MapCustomerKYCToResponsePaginated(result), nil
}

func (s *selfActivationKYCService) FindByID(ctx context.Context, id string) (*cust_kyc_dto.CustomerKYCResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindByID", "SelfActivationKYC", "FindByID")
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

func (s *selfActivationKYCService) EnableOrDisable(ctx context.Context, id, reason string, enable bool) error {
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
		action = constants.RequestApproveSelfActivateKyc
	} else {
		action = constants.RequestRejectSelfActivateKyc
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

	kyc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] find err: %v", err)
		return nil, err
	}

	if kyc.KYCStatus != imodel.KYCStatusPending {
		log.Warnf("[SelfActivationKYC][StartKycReview] kyc status is not pending id: %s, status: %s", id, kyc.KYCStatus)
		return nil, errors.New("KYC review can only be started for KYC requests with pending status")
	}

	existingReview, err := s.repo.FindKycInReview(ctx, id)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		log.Errorf("[SelfActivationKYC][StartKycReview] failed to check existing review: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if existingReview != nil && existingReview.IsActive && existingReview.ExpiresAt.After(time.Now()) {
		log.Warnf("[SelfActivationKYC][StartKycReview] active review already exists for kyc id: %s", id)
		return nil, errors.New("An active review already exists for this KYC request")
	}

	if existingReview != nil && existingReview.ExpiresAt.Before(time.Now()) {
		log.Warnf("[SelfActivationKYC][StartKycReview] review already exists but expired for kyc id: %s", id)
		return nil, errors.New("An expired review already exists. Please pick the review to restart the review process.")
	}

	cpsUser, err := s.cpsUserRepo.FindByID(ctx, makerUser.UserCode)
	if err != nil {
		log.Errorf("[SelfActivationKYC][StartKycReview] failed to fetch user info: %v", err)
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

func (s *selfActivationKYCService) PickKycReview(ctx context.Context, id string, reason string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	kycInReview, err := s.repo.FindKycInReview(ctx, id)
	if err != nil {
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New("No started  KYC review found.")
		}
		log.Errorf("[SelfActivationKYC][PickKycReview] find err: %v", err)
		return err
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

	cpsActionData := lib.CpsModelBuilder(id, makerUser, nil, newReq, string(constants.RequestPickSelfActivateKycReview), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SelfActivationKYC][PickKycReview] cps action err: %v", err)
		return err
	}

	return nil
}

func (s *selfActivationKYCService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "SelfActivationKYC", "Authorize")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[SelfActivationKYC][Authorize] action: %s", cpsAction.RequestAction)

	switch string(cpsAction.RequestAction) {
	// case string(constants.RequestApproveSelfActivateKyc):

	case string(constants.RequestRejectSelfActivateKyc):
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
			log.Errorf("[SelfActivationKYC][Authorize] failed to check existing review: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		if existingReview != nil {
			existingReview.ReviewStatus = string(constants.KYCStatusRejected)
			_, err = s.repo.UpdateKycReview(ctx, cpsAction.UniqueId, existingReview)
			if err != nil {
				log.Errorf("[SelfActivationKYC][Authorize] failed to update review status: %v", err)
				return nil, errors.New(localization.ErrorUnexpectedError.Code)
			}
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
					s.logger.Errorf("[SelfActivationKYC][Authorize] SMS send failed for phone %s: %v", phone, err)
				}
			}()
		}

	case string(constants.RequestPickSelfActivateKycReview):
		reviewData, err := local_util.JsonUnmarshal[imodel.StartedKycReview](cpsAction.CurrentAction)
		if err != nil {
			return nil, err
		}

		existingReview, err := s.repo.FindKycInReview(ctx, cpsAction.UniqueId)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			log.Errorf("[SelfActivationKYC][Authorize] failed to check existing review: %v", err)
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
			log.Errorf("[SelfActivationKYC][Authorize] failed to update review expiration: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	default:
		return nil, fmt.Errorf("unsupported action: %s", cpsAction.RequestAction)
	}

	return cpsAction, nil
}
