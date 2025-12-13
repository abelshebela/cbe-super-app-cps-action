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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (s *KYCVerifier) FetchKYCList(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dto.KYCVerifierResponse], error) {
	return s.repo.FindAllWithPaginationPopulated(ctx, *filterParams)
}

func (s *KYCVerifier) FetchKYCByID(ctx context.Context, id string) (*dto.KYCVerifierResponse, error) {
	return s.repo.FindByIDPopulated(ctx, id)
}

func (s *KYCVerifier) UpdateKYC(ctx context.Context, id string, req dto.UpdateKYCRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := core.ValidateUpdate(req); err != nil {
		return err
	}
	update := core.MapUpdateToModel(prev, req)
	cpsAction := lib.CpsModelBuilder(id, makerData, prev, update, string(constants.RequestUpdateKYC), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}
	return nil
}

func (s *KYCVerifier) ApproveKYC(ctx context.Context, id string, req dto.ApproveKYCRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := core.ValidateApprove(req); err != nil {
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
		return err
	}
	return nil
}

func (s *KYCVerifier) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	if action.ActionStatus != constants.Approved {
		s.logger.Errorf("Tried to authorize KYC action without approval")
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestUpdateKYC):
		updated, err := local_util.JsonUnmarshal[model.CustomerKYC](action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		lib.GoRoutinBaker(types.BakerOptions{}, func() {
			if err := s.repo.Update(ctx, action.UniqueId, updated); err != nil {
				s.logger.Errorf("[KYCVerifier] Error updating KYC in job proccess: %v", err)
			}
		}, func() {
			bgCtx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.BgJobTimeout))
			defer cancel()

			user, err := s.userRepo.FindById(bgCtx, action.UniqueId)
			if err != nil {
				s.logger.Errorf("[KYCVerifier] Error finding user in job proccess: %v", err)
			}

			if err := core.AccountCreateAndLink(bgCtx, *action, action.UniqueId, *user, s.accountService, s.userRepo, s.linkedAccountRepo, s.logger); err != nil {
				s.logger.Errorf("[KYCVerifier] Error creating account in job proccess: %v", err)
			}

		})

	case string(constants.RequestApproveKYC):
		// CurrentAction contains the fully-updated CustomerKYC
		updated, err := local_util.JsonUnmarshal[model.CustomerKYC](action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}
		if err := s.repo.Update(ctx, action.UniqueId, updated); err != nil {
			return nil, err
		}

		// update user
		if err := core.MapandUpdateuserFromKYC(ctx, s.userRepo, *updated, s.logger); err != nil {
			s.logger.Errorf("[KYCVerifier] Error updating user in job proccess: %v", err)
		}
	default:
		s.logger.Errorf("Unsupported action requested: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
	return action, nil
}
