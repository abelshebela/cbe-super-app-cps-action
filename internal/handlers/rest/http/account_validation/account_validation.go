package accountvalidation

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/account_validation"
	accountvalidationInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_validation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accountValidationAdapter struct {
	accountValidationService service.AccountValidationService
	logger                   utils.Logger
}

func NewHttpAccountValidation(accountValidationService service.AccountValidationService, logger utils.Logger) accountvalidationInterface.AccountValidation {
	return &accountValidationAdapter{
		accountValidationService: accountValidationService,
		logger:                   logger,
	}
}
func (h *accountValidationAdapter) FetchAccountValidation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // get from path instead of query
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequest.Code)
		return
	}

	resp, err := h.accountValidationService.GetAccountValidation(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data, err := StructToMap(resp)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessValidationRuleFetched, data)
}

func (h *accountValidationAdapter) UpdateAccountValidation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// ─── Parse Request Body ───────────────────────────────────────────────
	var req dto.ValidationRuleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorResponse(w, localization.ErrorValidationRuleConvertIDFailed, nil, nil)
		return
	}

	// ─── Validation Checks ───────────────────────────────────────────────
	if req.MinLength > req.MaxLength {
		localization.SendErrorResponse(w, localization.ErrorValidationRuleMinMaxLengthMismatch, nil, nil)
		return
	}

	if err := req.Validate(); err != nil {
		// assuming Validate returns a custom error with code
		localization.SendErrorResponse(w, localization.ErrorValidationFailed, nil, nil)
		return
	}

	// ─── Extract User Context ────────────────────────────────────────────
	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	maker := dto.User{
		ID:          userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}

	// ─── Map DTO To Model ────────────────────────────────────────────────
	rule := dto.ToModel(req)

	if err := h.accountValidationService.UpdateAccountValidation(r.Context(), id, rule); err != nil {
		localization.SendErrorResponse(w, localization.ErrorValidationRuleUpdateFailed, nil, nil)
		return
	}

	// ─── Success Response ────────────────────────────────────────────────
	localization.SendSuccessResponse(w, localization.SuccessValidationRuleApproved, map[string]interface{}{
		"maker": maker, // optional, if you want to include info
	})
}

func (s *accountValidationAdapter) FetchAllAccountValidation(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	accountValidation, err := s.accountValidationService.GetAllAccountValidation(ctx, *filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessValidationRuleFetched, accountValidation)
}
