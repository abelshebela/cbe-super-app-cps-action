package productcode

import (
	"encoding/json"
	"fmt"
	"net/http"

	dto "cbe-super-app-cps-action/internal/constants/dto/productcode"
	"cbe-super-app-cps-action/internal/constants/errors"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/pkgs/utils"

	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// ProductCodeAdapter handles HTTP requests for product code operations
type ProductCodeAdapter struct {
	productCodeApplication service.ProductCodeService
	logger                 shared.Logger
}

// NewProductCodeHTTPHandler initializes a new ProductCodeAdapter
func InitProductcodeAdapter(service service.ProductCodeService, logger shared.Logger) *ProductCodeAdapter {
	return &ProductCodeAdapter{
		productCodeApplication: service,
		logger:                 logger,
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
		localization.SendErrorByCodeResponse(w, errors.ErrBadRequest.Error())
		return
	}

	err = req.Validate()
	if err != nil {
		h.logger.Errorf("[productcode.UpdateProductCode] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	domainReq := dto.ToDomainProductCodeRequest(req)

	// For now, we'll skip the maker parameter until user context is implemented
	old, new, err := h.productCodeApplication.UpdateProductCode(r.Context(), domainReq)
	fmt.Println("the error:", err)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{
		Type: "success",
		StatusCode: 200,
	}, map[string]*model.ProductCode{
		"old": old,
		"new": new,
	})
}

// FetchProductCodeByID retrieves a single product code
func (h *ProductCodeAdapter) FetchProductCodeByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractID(w, r)
	if err != nil {
		return
	}

	data, err := h.productCodeApplication.FetchProductCodeByID(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	res := dto.ToProductCodeResponse(*data)
	localization.SendSuccessResponse(w, localization.ResponseCode{
		Type: "success",
		StatusCode: 200,
	}, res)
}

// FetchProductCodes retrieves all product codes
func (h *ProductCodeAdapter) FetchProductCodes(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractFilterParams(r)
	list, err := h.productCodeApplication.FetchAllProductCodes(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	docs := dto.ToProductCodeResponses(list.Data)
	res := types.PaginatedResponse[[]*dto.ProductCodeResponse]{
		Data: docs,
		Meta: list.Meta,
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{
		Type: "success",
		StatusCode: 200,
	}, res)
}
