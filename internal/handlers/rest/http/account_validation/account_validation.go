package accountvalidation

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/account_validation"
	accountvalidationInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_validation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type paginated_resp *types.PaginatedResponse[[]*model.ValidationRule]
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

// FindById godoc
// @Summary Get account validation by ID
// @Description Retrieve a specific account validation rule by its ID.
// @Tags Account Validation
// @Accept json
// @Produce json
// @Param id path string true "Validation Rule ID"
// @Success 200 {object} localization.StandardResponse{data=accountvalidation.ValidationRuleDTO} "Validation rule retrieved"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /account_validation/{id} [get]
func (h *accountValidationAdapter) FindById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequest.Code)
		return
	}

	resp, err := h.accountValidationService.FindById(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	//note:we need to  convert model to dto

	// Send the struct directly instead of converting to map
	localization.SendSuccessResponse(w, localization.SuccessValidationRuleFetched, resp)
}

// Update godoc
// @Summary Update account validation rule
// @Description Update an existing account validation rule. MinLength must not exceed MaxLength.
// @Tags Account Validation
// @Accept json
// @Produce json
// @Param id path string true "Validation Rule ID"
// @Param request body accountvalidation.ValidationRuleDTO true "Validation rule data" example({"entity_type":"ACCOUNT","validation_for":"NUMBER","identifier":"ACCOUNT_NUMBER","min_length":5,"max_length":20,"enabled":true})
// @Success 200 {object} localization.StandardResponse{data=nil} "Validation rule updated"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Validation failed"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /account_validation/update/{id} [patch]
func (h *accountValidationAdapter) Update(w http.ResponseWriter, r *http.Request) {
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

	// ─── Map DTO To Model ────────────────────────────────────────────────
	rule := dto.ToModel(req)

	if err := h.accountValidationService.Update(r.Context(), id, rule); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	// ─── Success Response ────────────────────────────────────────────────
	localization.SendSuccessResponse(w, localization.SuccessValidationRuleApproved, localization.ResponseCode{})
}

// FindAllWithPagination godoc
// @Summary List account validation rules
// @Description Retrieve account validation rules with pagination and optional search.
// @Tags Account Validation
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1) example(1)
// @Param per_page query int false "Items per page" default(10) minimum(1) maximum(100) example(10)
// @Param search query string false "Search term" example("ACCOUNT_NUMBER")
// @Success 200 {object} localization.StandardResponse{data=paginated_resp} "Validation rules retrieved"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /account_validation [get]
func (s *accountValidationAdapter) FindAllWithPagination(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	accountValidation, err := s.accountValidationService.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessValidationRuleFetched, accountValidation)
}
