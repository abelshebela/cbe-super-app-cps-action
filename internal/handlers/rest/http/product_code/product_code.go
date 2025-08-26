package productcode

import (
	"encoding/json"
	"net/http"

	dto "cbe-super-app-cps-action/internal/constants/dto/productcode"
	"cbe-super-app-cps-action/internal/constants/errors"
	"cbe-super-app-cps-action/internal/constants/response"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/pkgs/utils"

	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// ProductCodeAdapter handles HTTP requests for product code operations
type ProductCodeAdapter struct {
	service service.ProductCodeService
	logger  shared.Logger
}

// NewProductCodeHTTPHandler initializes a new ProductCodeAdapter
func InitProductcodeAdapter(service service.ProductCodeService, logger shared.Logger) *ProductCodeAdapter {
	return &ProductCodeAdapter{
		service: service,
		logger:  logger,
	}
}

// UpdateProductCode handles updating a product code
func (h *ProductCodeAdapter) UpdateProductCode(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractID(w, r)
	if err != nil {
		return
	}

	var req dto.UpdateProductCodeRequest
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[productcode.UpdateProductCode] failed to parse JSON: %v", err)
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	err = req.Validate() // Note: Changed to exported method
	if err != nil {
		h.logger.Errorf("[productcode.UpdateProductCode] validation error: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	// TODO: Implement user context extraction when available
	// maker, ok := h.extractUserAndMaker(w, r)
	// if !ok {
	//     return
	// }

	domainReq := dto.ToDomainProductCodeRequest(req)

	// For now, we'll skip the maker parameter until user context is implemented
	_, existing, err := h.service.UpdateProductCode(r.Context(), domainReq)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Product code update request submitted successfully", existing, nil)
}

// FetchProductCodeByID retrieves a single product code
func (h *ProductCodeAdapter) FetchProductCodeByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractID(w, r)
	if err != nil {
		return
	}

	data, err := h.service.FetchProductCodeByID(r.Context(), id)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}
	res := dto.ToProductCodeResponse(*data)
	response.SendSuccessResponse(w, http.StatusOK, "Product code successfully retrieved", res, nil)
}

// FetchProductCodes retrieves all product codes
func (h *ProductCodeAdapter) FetchProductCodes(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractFilterParams(r)
	list, err := h.service.FetchAllProductCodes(r.Context(), filterParams)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	docs := dto.ToProductCodeResponses(list.Data)
	res := types.PaginatedResponse[[]*dto.ProductCodeResponse]{
		Data: docs,
		Meta: list.Meta,
	}
	response.SendSuccessResponse(w, http.StatusOK, "Product codes successfully retrieved", res, nil)
}

// TODO: Implement user context extraction when the utility is available
// func (h *ProductCodeAdapter) extractUserAndMaker(w http.ResponseWriter, r *http.Request) (cps_entities.User, bool) {
//     userContext := context.ExtractUserContext(r)
//     if userContext.IsIncomplete() {
//         h.logger.Errorf("[productcode.extractUserAndMaker] incomplete user context")
//         response.SendErrorResponse(w, errors.ErrUnauthorized)
//         return cps_entities.User{}, false
//     }
//     maker := cps_entities.User{
//         UserCode:    userContext.UserID,
//         FullName:    userContext.FullName,
//         PhoneNumber: userContext.PhoneNumber,
//         Department:  userContext.Department,
//     }
//     return maker, true
// }
