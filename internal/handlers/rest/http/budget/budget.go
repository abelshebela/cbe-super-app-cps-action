package budget

import (
	budget_dto "cbe-super-app-cps-action/internal/constants/dto/budget"
	"cbe-super-app-cps-action/internal/constants/interfaces/budget"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/handlers/rest/http/budget/core"
	"cbe-super-app-cps-action/internal/service"
	common_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type budgetAdapter struct {
	budgetApplication service.BudgetService
	logger            utils.Logger
}

func InitBudgetAdapter(budgetApplication service.BudgetService, logger utils.Logger) budget.BudgetPortHandler {
	return &budgetAdapter{
		logger:            logger,
		budgetApplication: budgetApplication,
	}
}

func (b *budgetAdapter) CreateBudgetIcon(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := core.ParseMultipartFormFile(r, "icons_image", 10<<20)
	if err != nil {
		b.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := core.FileValidator(w, file, *fileHeader, b.logger); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err = b.budgetApplication.CreateBudgetIcon(r.Context(), fileHeader, &file)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetIconRequestSubmittedForApproval, nil)
}
func (b *budgetAdapter) BudgetFetchIcons(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	icons, err := b.budgetApplication.BudgetFetchIcons(r.Context(), filterParams)
	if err != nil {
		b.logger.Errorf("failed to fetch budget icons: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetIconsFetched, icons)
}

func (b *budgetAdapter) BudgetUpdateIcon(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	file, fileHeader, err := core.ParseMultipartFormFile(r, "icons_image", 10<<20)
	if err != nil {
		b.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorMissingOrInvalidImage, nil, nil)
		return
	}

	if err := core.FileValidator(w, file, *fileHeader, b.logger); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err = b.budgetApplication.BudgetUpdateIcon(r.Context(), id, fileHeader, &file)
	if err != nil {
		b.logger.Errorf("failed to update budger icon: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())

		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetIconRequestSubmittedForApproval, nil)
}
func (b *budgetAdapter) BudgetCreateColor(w http.ResponseWriter, r *http.Request) {
	var req budget_dto.BudgetCreateColor

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		b.logger.Errorf("error decoding request body: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidJSONPayload, nil, nil)
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	color := &model.Color{
		Color: req.Color,
	}
	err := b.budgetApplication.BudgetCreateColor(r.Context(), color)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetColorRequestSubmittedForApproval, nil)

}
func (b *budgetAdapter) BudgetFetchColors(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	if filterParams.Page < 1 {
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	colors, err := b.budgetApplication.BudgetFetchColors(r.Context(), filterParams)
	if err != nil {
		b.logger.Errorf("failed to fetch budget colors: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetColorsFetched, colors)
}
func (b *budgetAdapter) BudgetUpdateColor(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	var req budget_dto.UpdateColorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidJSONPayload, nil, nil)
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	color := &model.Color{
		Color: req.Color,
	}

	err := b.budgetApplication.BudgetUpdateColor(r.Context(), id, color)
	if err != nil {
		b.logger.Errorf("failed to update budget color: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetColorUpdateSubmittedForApproval, nil)
}

func (b *budgetAdapter) BudgetCheckerApproval(w http.ResponseWriter, r *http.Request) {
	actionCode, ok := common_util.GetParam(r, "action_code")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'action_code'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err := b.budgetApplication.BudgetCheckerApproval(r.Context(), actionCode)
	if err != nil {
		b.logger.Errorf("failed to approve action: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCheckerActionApproved, nil)
}
