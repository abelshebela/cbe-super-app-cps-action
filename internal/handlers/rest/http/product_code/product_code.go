package productcode

import (
	"encoding/json"

	"net/http"

	product_code_dto "cbe-super-app-cps-action/internal/constants/dto/productcode"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/pkgs/utils"

	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

// Update Product Code
//
//	@Summary		Update a product code
//	@Description	Updates an existing product code by its ID. This is a pending action that requires approval.
//	@Tags			ProductCode
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string									true	"Product Code ID"
//	@Param			productCode			body		productcode.UpdateProductCodeRequest	true	"Update Product Code Request"
//	@Success		200					{object}	localization.StandardResponse{data=map[string]model.ProductCode}
//	@Failure		400,401,403,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/product-codes/{id} [patch]
func (h *ProductCodeAdapter) UpdateProductCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.TraceLogger(r.Context(), "handler", "productCode", "ProductCodeAdapter", "UpdateProductCode")
	defer span.End()
	id, err := utils.ExtractID(w, r)
	if err != nil {
		span.AddEvent("Failed to extract ID", trace.WithAttributes(attribute.String("error", err.Error())))
		h.logger.Errorf("[productcode.UpdateProductCode] failed to extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	var req product_code_dto.UpdateProductCodeRequest
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.AddEvent("Failed to parse JSON", trace.WithAttributes(attribute.String("error", err.Error())))
		h.logger.Errorf("[productcode.UpdateProductCode] failed to parse JSON: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	err = req.Validate()
	if err != nil {
		span.AddEvent("Validation error", trace.WithAttributes(attribute.String("error", err.Error())))
		h.logger.Errorf("[productcode.UpdateProductCode] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	domainReq := product_code_dto.ToDomainProductCodeRequest(req)

	old, new, err := h.productCodeApplication.UpdateProductCode(ctx, domainReq)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		h.logger.Errorf("[UpdateProductCode] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("ProductCode updated", trace.WithAttributes(attribute.String("id", id)))
	h.logger.Infof("[UpdateProductCode] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessProductCodeUpdated, map[string]*model.ProductCode{
		"old": old,
		"new": new,
	})
}

// Get Product Code by ID
//
//	@Summary		Fetch a product code by ID
//	@Description	Retrieves a single product code by its unique ID.
//	@Tags			ProductCode
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id					path		string	true	"Product Code ID"
//	@Success		200					{object}	localization.StandardResponse{data=productcode.ProductCodeResponse}
//	@Failure		400,401,403,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/product-codes/{id} [get]
func (h *ProductCodeAdapter) FetchProductCodeByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.TraceLogger(r.Context(), "handler", "productCode", "ProductCodeAdapter", "FetchProductCodeByID")
	defer span.End()
	id, err := utils.ExtractID(w, r)
	if err != nil {
		span.AddEvent("Failed to extract ID", trace.WithAttributes(attribute.String("error", err.Error())))
		return
	}

	data, err := h.productCodeApplication.FetchProductCodeByID(ctx, id)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		h.logger.Errorf("[FetchProductCodeByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("ProductCode retrieved", trace.WithAttributes(attribute.String("id", id)))
	h.logger.Infof("[FetchProductCodeByID] product code retrieved successfully for id: %s", id)
	res := product_code_dto.ToProductCodeResponse(*data)
	localization.SendSuccessResponse(w, localization.SuccessProductCodeFetched, res)
}

type ProductCodePaginatedResponse types.PaginatedResponse[[]*product_code_dto.ProductCodeResponse]

// Get All Product Codes
//
//	@Summary		Fetch all product codes
//	@Description	Retrieves a paginated list of product codes. Search using service_name.Filter using {_id,service_name,created_at,...}
//	@Tags			ProductCode
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page			query		int		false	"Page number"
//	@Param			per_page		query		int		false	"Items per page"
//	@Param			search			query		string	false	"Search term"
//	@Param			filter			query		string	false	"filter term"
//	@Success		200				{object}	localization.StandardResponse{data=ProductCodePaginatedResponse}
//	@Failure		400,401,403,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/productcodes [get]
func (h *ProductCodeAdapter) FetchProductCodes(w http.ResponseWriter, r *http.Request) {
	ctx, span := utils.TraceLogger(r.Context(), "handler", "productCode", "ProductCodeAdapter", "FetchProductCodes")
	defer span.End()
	filterParams := utils.ExtractFilterParams(r)
	list, err := h.productCodeApplication.FetchAllProductCodes(ctx, filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		h.logger.Errorf("[FetchProductCodes] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("ProductCodes retrieved", trace.WithAttributes(attribute.Int("count", len(list.Data))))
	h.logger.Infof("[FetchProductCodes] retrieved %d product codes", len(list.Data))
	docs := product_code_dto.ToProductCodeResponses(list.Data)
	res := types.PaginatedResponse[[]*product_code_dto.ProductCodeResponse]{
		Data: docs,
		Meta: list.Meta,
	}
	localization.SendSuccessResponse(w, localization.SuccessProductCodesFetched, res)
}
