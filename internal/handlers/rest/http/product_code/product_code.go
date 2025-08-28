package productcode

import (
	"encoding/json"
	"net/http"

	dto "cbe-super-app-cps-action/internal/constants/dto/productcode"
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
		h.logger.Errorf("[productcode.UpdateProductCode] failed to extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorToResponseCode(err.Error(), int(http.StatusBadRequest), "failed to extract ID").Code)
		return
	}

	var req dto.UpdateProductCodeRequest
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[productcode.UpdateProductCode] failed to parse JSON: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorToResponseCode(err.Error(), int(http.StatusBadRequest), "failed to parse JSON").Code)
		return
	}

	err = req.Validate()
	if err != nil {
		h.logger.Errorf("[productcode.UpdateProductCode] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorToResponseCode(err.Error(), int(http.StatusBadRequest), "validation error").Code)
		return
	}

	domainReq := dto.ToDomainProductCodeRequest(req)

	old, new, err := h.productCodeApplication.UpdateProductCode(r.Context(), domainReq)
	if err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorToResponseCode(err.Error(), int(http.StatusBadRequest), "failed to make update product code request").Code)
		return
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{
		Type:       "success",
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
		localization.SendErrorByCodeResponse(w, localization.ErrorToResponseCode(err.Error(), int(http.StatusBadRequest), "failed to fetch product by id").Code)
		return
	}
	res := dto.ToProductCodeResponse(*data)
	localization.SendSuccessResponse(w, localization.ResponseCode{
		Type:       "success",
		StatusCode: 200,
	}, res)
}

// FetchProductCodes retrieves all product codes
func (h *ProductCodeAdapter) FetchProductCodes(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractFilterParams(r)
	list, err := h.productCodeApplication.FetchAllProductCodes(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorToResponseCode(err.Error(), int(http.StatusBadRequest), "failed to fetch all product codes").Code)
		return
	}
	docs := dto.ToProductCodeResponses(list.Data)
	res := types.PaginatedResponse[[]*dto.ProductCodeResponse]{
		Data: docs,
		Meta: list.Meta,
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{
		Type:       "success",
		StatusCode: 200,
	}, res)
}
