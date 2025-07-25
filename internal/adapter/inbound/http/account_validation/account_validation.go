package accountvalidation_inbound

import (
	"encoding/json"
	"fmt"

	// "fmt"
	"net/http"

	accountvalidation_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_validation"

	accountvalidation "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_validation"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	common_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type HttpStore struct {
	Application accountvalidation.ApplicationAbstracts
	logger      common_utils.Logger
}

func NewHttpAccountValidation(app accountvalidation.ApplicationAbstracts, logger common_utils.Logger) inbound.Inbound {
	return &HttpStore{
		Application: app,
		logger:      logger,
	}
}

func (h *HttpStore) FetchAccountValidation(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		utils.SendErrorResponse(w, utils.InvalidInputParameters, 0, nil)
		return
	}

	resp, err := h.Application.GetAccountValidation(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := utils.StructToMap(resp.Validation)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 500, nil)
	}
	utils.BaseResponseMaker(data, w, "Account fetched successfully", 200)
}
func (h *HttpStore) FetchAllAccountValidation(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		common_util.SendErrorResponse(w, "INVALID_INPUT_PARAMETERS", 0, nil)
		return
	}

	resp, err := h.Application.GetAllAccountValidation(r.Context(), common_util.Filter(*filterParams))
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.BaseResponseMaker(resp, w, "Account fetched successfully", 200)
}

func (h *HttpStore) UpdateAccountValidationMaker(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		utils.SendErrorResponse(w, utils.InvalidInputParameters, 0, nil)
		return
	}

	var req accountvalidation_app.ValidationRuleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}

	fmt.Println("**************************88")
	fmt.Println(req.MaxLength < req.MinLength)
	fmt.Println("**************************88")

	if req.MinLength > req.MaxLength {
		utils.SendErrorResponse(w, "MAX_NOT_BE_LESS", 0, nil)
		return
	}
	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}

	maker := domain.User{
		ID:          userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}

	_, err := h.Application.UpdateAccountValidationRequest(r.Context(), id, accountvalidation_app.ToDomainValidationRule(req), maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	message := "Update request submitted for approval"
	utils.BaseResponseMaker(map[string]interface{}{}, w, message, 200)
}

func (h *HttpStore) UpdateAccountValidationChecker(w http.ResponseWriter, r *http.Request) {
	actionCode, ok := common_util.GetParam(r, "action_code")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'action_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	var req accountvalidation_app.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}

	req.ActionCode = actionCode
	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}

	checker := domain.User{
		ID:          userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}

	decison := common_util.Decison(req.Decison)
	err := h.Application.UpdateAccountValidation(
		r.Context(),
		req.ActionCode,
		decison,
		checker,
		func() string {
			if !req.Decison {
				return req.RejectedReason
			}
			return ""
		}(),
	)
	if err != nil {
		errMsg := err.Error()
		utils.SendErrorResponse(w, errMsg, 0, nil)
		return
	}

	action := "approved"
	if !req.Decison {
		action = "rejected"
	}
	data, err := utils.StructToMap(map[string]interface{}{
		"action": action,
	})
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}
	message := "update request " + action + " successfully Approved"
	utils.BaseResponseMaker(data, w, message, 200)
}
