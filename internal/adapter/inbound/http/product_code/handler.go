package productcode

import (
	"encoding/json"
	"net/http"

	productcode "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/product_code"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/product_code"
	context "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// ProductCodeHTTPStore handles HTTP requests for product code operations
type ProductCodeHTTPStore struct {
	Application productcode.Application
	logger      shared.Logger
}

// NewProductCodeHTTPHandler initializes a new ProductCodeHTTPStore
func NewProductCodeHTTPHandler(app productcode.Application, logger shared.Logger) outbound.ProductCodeHandler {
	return &ProductCodeHTTPStore{
		Application: app,
		logger:      logger,
	}
}

// extractUserAndMaker extracts user context and creates a maker, sending an error response if incomplete
func (h *ProductCodeHTTPStore) extractUserAndMaker(w http.ResponseWriter, r *http.Request) (cps_entities.User, bool) {
	userContext := context.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		h.logger.Errorf("[productcode.extractUserAndMaker] incomplete user context")
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return cps_entities.User{}, false
	}
	maker := cps_entities.User{
		UserCode:    userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}
	return maker, true
}

// extractID extracts and validates the ID parameter, sending an error response if invalid
func (h *ProductCodeHTTPStore) extractID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id, ok := utils.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("[productcode.extractID] missing or invalid parameter 'id'")
		utils.SendErrorResponse(w, utils.InvalidInputParameters, 0, nil)
		return "", false
	}
	return id, true
}

// sendResponse sends success or error responses
func (h *ProductCodeHTTPStore) sendResponse(w http.ResponseWriter, err error, successMessage string, data interface{}) {
	if err != nil {
		h.logger.Errorf("[productcode.sendResponse] error: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, data, successMessage)
}

// UpdateProductCode handles updating a product code
func (h *ProductCodeHTTPStore) UpdateProductCode(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	var req ProductCodeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[productcode.parseAndValidateProductCodeRequest] failed to parse JSON: %v", err)
		utils.SendErrorResponse(w, utils.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	domainReq := ToDomainProductCodeRequest(req)
	if req.isEmpty() {
		h.logger.Errorf("[productcode.UpdateProductCode] no data provided for update, id: %s", id)
		utils.SendErrorResponse(w, utils.NoDataProvidedForUpdate, http.StatusBadRequest, nil)
		return
	}

	err := h.Application.Update(r.Context(), id, domainReq, maker)
	h.sendResponse(w, err, "Product code update request submitted successfully", nil)
}

// FetchProductCodeByID retrieves a single product code
func (h *ProductCodeHTTPStore) FetchProductCodeByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	data, err := h.Application.FetchByID(r.Context(), id)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	res := ToProductCodeResponse(*data)
	h.sendResponse(w, nil, "Product code successfully retrieved", res)
}

// FetchProductCodes retrieves all product codes
func (h *ProductCodeHTTPStore) FetchProductCodes(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractFilterParams(r)
	list, err := h.Application.FetchAll(r.Context(), filterParams)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	docs := ToProductCodeResponses(list.Data)
	res := utils.PaginatedResponse[[]*ProductCodeResponse]{
		Data: docs,
		Meta: list.Meta,
	}
	h.sendResponse(w, nil, "Product codes successfully retrieved", res)
}
