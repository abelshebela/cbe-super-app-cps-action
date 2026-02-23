package kyc_verifier

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/kyc_verifier/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type KYCVerifier struct {
	repo              storage.KYCVerifierRepository
	userRepo          storage.UserRepository
	cpsService        service.CPSActionService
	accountService    account_lookup.Account
	linkedAccountRepo storage.LinkedAccountRepository
	logger            utils.Logger
	cfg               config.VaultConfig
}

func NewKYCVerifierService(client *mongo.Client, repo storage.KYCVerifierRepository, userRepo storage.UserRepository, accountLookUpService account_lookup.Account, cpsAction service.CPSActionService, linkedAccountRepo storage.LinkedAccountRepository, cfg config.VaultConfig, logger utils.Logger) service.KYCVerifierService {
	return &KYCVerifier{
		repo:              repo,
		userRepo:          userRepo,
		cpsService:        cpsAction,
		accountService:    accountLookUpService,
		linkedAccountRepo: linkedAccountRepo,
		logger:            logger,
		cfg:               cfg,
	}
}

func (s *KYCVerifier) FetchKYCList(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]dto.KYCVerifierResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchKYCList", "KYCVerifier", "FetchKYCList")
	defer span.End()

	result, err := s.repo.FindAllWithPaginationPopulated(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch KYC list", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}

func (s *KYCVerifier) FetchKYCByID(ctx context.Context, id string) (*dto.KYCVerifierResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchKYCByID", "KYCVerifier", "FetchKYCByID")
	defer span.End()

	result, err := s.repo.FindByIDPopulated(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch KYC", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return result, nil
}

func (s *KYCVerifier) UpdateKYC(ctx context.Context, id string, req dto.UpdateKYCRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateKYC", "KYCVerifier", "UpdateKYC")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find KYC", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if err := core.ValidateUpdate(req); err != nil {
		span.AddEvent("Validation failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	update := core.MapUpdateToModel(prev, req)
	cpsAction := lib.CpsModelBuilder(id, makerData, prev, update, string(constants.RequestUpdateKYC), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	return nil
}

func (s *KYCVerifier) ApproveKYC(ctx context.Context, id string, req dto.ApproveKYCRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "ApproveKYC", "KYCVerifier", "ApproveKYC")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find KYC", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if err := core.ValidateApprove(req); err != nil {
		span.AddEvent("Validation failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	// 	type KYCStatus string

	// const (
	// 	KYCStatusPending  KYCStatus = "PENDING"
	// 	KYCStatusApproved KYCStatus = "APPROVED"
	// 	KYCStatusRejected KYCStatus = "REJECTED"
	// )

	// Build the fully updated KYC document to include in CurrentAction
	updated := *prev
	updated.KYCApproved = req.Approve
	if req.Approve {
		updated.KYCStatus = shared_constant.KYCStatusApproved
	} else {
		updated.KYCStatus = shared_constant.KYCStatusRejected
	}
	updated.KYCActivityBy = map[string]any{"admin_id": makerData.UserID, "full_name": makerData.FullName, "role": "kyc_verifier"}
	if !req.Approve {
		updated.KYCRejectReason = req.Reason
	}
	cpsAction := lib.CpsModelBuilder(id, makerData, prev, &updated, string(constants.RequestApproveKYC), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	return nil
}

func (s *KYCVerifier) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "KYCVerifier", "Authorize")
	defer span.End()

	if action.ActionStatus != constants.Approved {
		s.logger.Errorf("[KycVerifSvc][Authorize] invalid status")
		span.AddEvent("CPS action status invalid", trace.WithAttributes(
			attribute.String("error", localization.ErrorCPSActionStatusInvalid.Code),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestUpdateKYC):
		updated, err := local_util.JsonUnmarshal[model.CustomerKYC](action.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		lib.GoRoutinBaker(types.BakerOptions{}, func() {
			if err := s.repo.Update(ctx, action.UniqueId, updated); err != nil {
				s.logger.Errorf("[KycVerifSvc][Authorize] update kyc job err: %v", err)
				span.AddEvent("Failed to update KYC in job process", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", action.UniqueId),
				))
			}
		}, func() {
			bgCtx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.BgJobTimeout))
			defer cancel()

			user, err := s.userRepo.FindById(bgCtx, action.UniqueId)
			if err != nil {
				s.logger.Errorf("[KycVerifSvc][Authorize] find user job err: %v", err)
				span.AddEvent("Failed to find user in job process", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", action.UniqueId),
				))
			}

			if err := core.AccountCreateAndLink(bgCtx, *action, action.UniqueId, *user, s.accountService, s.userRepo, s.linkedAccountRepo, s.logger); err != nil {
				s.logger.Errorf("[KycVerifSvc][Authorize] create account job err: %v", err)
				span.AddEvent("Failed to create account in job process", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", action.UniqueId),
				))
			}

		})

	case string(constants.RequestApproveKYC):
		// CurrentAction contains the fully-updated CustomerKYC
		updated, err := local_util.JsonUnmarshal[model.CustomerKYC](action.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
		if err := s.repo.Update(ctx, action.UniqueId, updated); err != nil {
			span.AddEvent("Failed to update KYC", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

		// update user
		if err := core.MapandUpdateuserFromKYC(ctx, s.userRepo, *updated, s.logger); err != nil {
			s.logger.Errorf("[KycVerifSvc][Authorize] update user err: %v", err)
			span.AddEvent("Failed to update user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
		}
	default:
		s.logger.Errorf("[KycVerifSvc][Authorize] unsupported action: %s", action.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(action.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
	return action, nil
}
