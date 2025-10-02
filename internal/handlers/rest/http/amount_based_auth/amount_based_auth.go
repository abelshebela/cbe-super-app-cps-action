package amount_based_auth

import (
	"encoding/json"
	"net/http"

	amountauthdto "cbe-super-app-cps-action/internal/constants/dto/amount_based_auth"
	amount_based "cbe-super-app-cps-action/internal/constants/interfaces/amount_based_auth"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	common_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type paginated_auth_tier_resp types.PaginatedResponse[[]*model.AuthTier]
type AmountBasedAuthHandler struct {
	Service service.AmountBasedAuthService
	logger  utils.Logger
}

func NewAmountBasedAuthHandler(service service.AmountBasedAuthService, logger utils.Logger) amount_based.AmountBasedAuthAdapter {
	return &AmountBasedAuthHandler{
		Service: service,
		logger:  logger,
	}
}

// GetAllAmountBasedAuth godoc
// @Summary List amount-based auth tiers
// @Description Fetch all amount-based authentication tiers with pagination.
// @Tags Amount-Based-Auth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1) example(1)
// @Param per_page query int false "Items per page" default(10) minimum(1) maximum(100) example(10)
// @Param search query string false "Search by method or range" example("PIN")
// @Success 200 {object} localization.StandardResponse{data=paginated_auth_tier_resp} "Fetched successfully"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /amount_based_auth [get]
func (a *AmountBasedAuthHandler) GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	customers, err := a.Service.FindAllWithPagination(r.Context(), *filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, customers)
}

// UpdateAmountBasedAuth godoc
// @Summary Update an amount-based auth tier
// @Description Update tier by ID and method. Validation depends on method: OPEN requires max_amount; PIN requires both min/max; OTP_PIN requires min_amount.
// @Tags Amount-Based-Auth
// @Accept json
// @Produce json
// @Param id path string true "Tier ID"
// @Param method path string true "Method" Enums(OPEN,PIN,OTP_PIN)
// @Param request body amountauthdto.UpdateAmountBasedAuthRequest true "Update payload" example({"min_amount":100,"max_amount":1000})
// @Success 200 {object} localization.StandardResponse{data=nil} "Update request sent"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid parameters or payload"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Tier not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Server error"
// @Security BearerAuth
// @Router /amount_based_auth/update/{id}/{method} [patch]
func (a *AmountBasedAuthHandler) UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	method, ok := common_util.GetParam(r, "method")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	id, ok := common_util.GetParam(r, "id")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var request amountauthdto.UpdateAmountBasedAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	// Validate the method parameter
	methodEnum := constants.Method(method)
	if methodEnum == constants.OPEN || methodEnum == constants.PIN || methodEnum == constants.OTPANDPIN {
	} else {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidMethod.Message)
		return
	}

	if !request.Validate(methodEnum) {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	if err := a.Service.UpdateAmountBasedAuth(r.Context(), id, methodEnum, request); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAmountBasedAuthRequestSent, nil)
}

// RejectAmountBasedAuth godoc
// @Summary Reject an amount-based auth action
// @Description Submit a rejection with CPS action payload.
// @Tags Amount-Based-Auth
// @Accept json
// @Produce json
// @Param id path string true "Action ID"
// @Param request body model.CPSAction true "CPS Action payload"
// @Success 200 {object} localization.StandardResponse{data=model.CPSAction} "Rejected"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Invalid input"
// @Security BearerAuth
// @Router /amount_based_auth/reject/{id} [patch]
func (a *AmountBasedAuthHandler) RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	_, ok := common_util.GetParam(r, "id")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var cpsReq model.CPSAction
	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	// For rejection, just return success since the actual rejection
	// would be handled by the CPS action system
	localization.SendSuccessResponse(w, localization.SuccessUserUpdated, cpsReq)
}
