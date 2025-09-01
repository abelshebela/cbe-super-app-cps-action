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


type ProductCodeAdapter struct {
	productCodeApplication service.ProductCodeService
	logger                 shared.Logger
}

func InitProductcodeAdapter(service service.ProductCodeService, logger shared.Logger) ProductCodeAdapter {
	return ProductCodeAdapter{
		productCodeApplication: service,
		logger:                 logger,
	}
}


// @Summary Update a product code
// @Description Updates an existing product code by its ID. This is a pending action that requires approval.
// @Tags Product Code
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product Code ID"
// @Param productCode body dto.UpdateProductCodeRequest true "Update Product Code Request"
// @Success 200 {object} map[string]model.ProductCode
// @Failure 400 {object} localization.ResponseCode
// @Failure 401 {object} localization.ResponseCode
// @Failure 403 {object} localization.ResponseCode
// @Failure 404 {object} localization.ResponseCode
// @Router /productcodes/{id} [patch]
func (h *ProductCodeAdapter) UpdateProductCode(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractID(w, r)
	if err != nil {
		h.logger.Errorf("[productcode.UpdateProductCode] failed to extract ID: %v", err)
		localization.GetResponseCodeByCode(err.Error())
		return
	}

	var req dto.UpdateProductCodeRequest
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[productcode.UpdateProductCode] failed to parse JSON: %v", err)
		localization.GetResponseCodeByCode(err.Error())
		return
	}

	err = req.Validate()
	if err != nil {
		h.logger.Errorf("[productcode.UpdateProductCode] validation error: %v", err)
		localization.GetResponseCodeByCode(err.Error())
		return
	}

	domainReq := dto.ToDomainProductCodeRequest(req)

	old, new, err := h.productCodeApplication.UpdateProductCode(r.Context(), domainReq)
	if err != nil {
		localization.GetResponseCodeByCode(err.Error())
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


// @Summary Fetch a product code by ID
// @Description Retrieves a single product code by its unique ID.
// @Tags Product Code
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product Code ID"
// @Success 200 {object} dto.ProductCodeResponse
// @Failure 400 {object} localization.ResponseCode
// @Failure 401 {object} localization.ResponseCode
// @Failure 403 {object} localization.ResponseCode
// @Failure 404 {object} localization.ResponseCode
// @Router /productcodes/{id} [get]
func (h *ProductCodeAdapter) FetchProductCodeByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractID(w, r)
	if err != nil {
		return
	}

	data, err := h.productCodeApplication.FetchProductCodeByID(r.Context(), id)
	if err != nil {
		localization.GetResponseCodeByCode(err.Error())
		return
	}
	res := dto.ToProductCodeResponse(*data)
	localization.SendSuccessResponse(w, localization.ResponseCode{
		Type:       "success",
		StatusCode: 200,
	}, res)
}


type ProductCodePaginatedResponse types.PaginatedResponse[[]*dto.ProductCodeResponse]

// @Summary Fetch all product codes
// @Description Retrieves a paginated list of product codes. Search using service_name.Filter using {_id,service_name,created_at,...}
// @Tags Product Code
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param search query string false "Search term"
// @Param filter query string false "filter term"
// @Success 200 {object} ProductCodePaginatedResponse
// @Failure 400 {object} localization.ResponseCode
// @Failure 401 {object} localization.ResponseCode
// @Failure 403 {object} localization.ResponseCode
// @Router /productcodes [get]
func (h *ProductCodeAdapter) FetchProductCodes(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractFilterParams(r)
	list, err := h.productCodeApplication.FetchAllProductCodes(r.Context(), filterParams)
	if err != nil {
		localization.GetResponseCodeByCode(err.Error())
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
