package kyc_verifier

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/kyc_verifier/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type KYCVerifier struct {
	Repo           storage.KYCVerifierRepository
	cpsService     service.CPSActionService
	AccountService account_lookup.Account
	logger         utils.Logger
}

func NewKYCVerifierService(client *mongo.Client, repo storage.KYCVerifierRepository, accountLookUpService account_lookup.Account, cpsAction service.CPSActionService, logger utils.Logger) service.KYCVerifierService {
	return &KYCVerifier{
		Repo:           repo,
		cpsService:     cpsAction,
		AccountService: accountLookUpService,
		logger:         logger,
	}
}

func (s *KYCVerifier) FetchKYCList(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dto.KYCVerifierResponse], error) {
	return s.Repo.FindAllWithPaginationPopulated(ctx, *filterParams)
}

func (s *KYCVerifier) FetchKYCByID(ctx context.Context, id string) (*dto.KYCVerifierResponse, error) {
	return s.Repo.FindByIDPopulated(ctx, id)
}

func (s *KYCVerifier) UpdateKYC(ctx context.Context, id string, req dto.UpdateKYCRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	prev, err := s.Repo.FindByID(ctx, id)
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

	prev, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := core.ValidateApprove(req); err != nil {
		return err
	}
	// Build the fully updated KYC document to include in CurrentAction
	updated := *prev
	updated.KYCApproved = req.Approve
	if !req.Approve {
		updated.KYCRejectReason = req.Reason
	} else {
		updated.KYCRejectReason = ""
	}
	updated.KYCActivityBy = map[string]any{"admin_id": makerData.UserID, "full_name": makerData.FullName, "role": "kyc_verifier"}
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
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		if err := s.Repo.Update(ctx, action.UniqueId, updated); err != nil {
			return nil, err
		}

	case string(constants.RequestApproveKYC):
		// CurrentAction contains the fully-updated CustomerKYC
		updated, err := local_util.JsonUnmarshal[model.CustomerKYC](action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		if err := s.Repo.Update(ctx, action.UniqueId, updated); err != nil {
			return nil, err
		}
	default:
		s.logger.Errorf("Unsupported action requested: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
	return action, nil
}
