package bulk_service

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/bulk_service"
	"cbe-super-app-cps-action/internal/constants/interfaces/bulk_service"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	// "fmt"

	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/localization"
	util "cbe-super-app-cps-action/pkgs/utils"

	"go.opentelemetry.io/otel/attribute"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type bulk_services_paginated_resp *types.PaginatedResponse[[]*model.APPAccessList]

type bulk_serviceAdapter struct {
	bulkService service.BulkService
	logger      utils.Logger
}

func InitBulkServiceAdapter(bulk_service service.BulkService, logger utils.Logger) bulk_service.BulkServiceHandler {
	return &bulk_serviceAdapter{
		logger:      logger,
		bulkService: bulk_service,
	}
}

// GetAllBulkServices retrieves all bulk services with pagination
//
//	@Summary		Get all bulk services
//	@Description	Retrieves a paginated list of all bulk services with optional filtering
//	@Tags			Bulk Services
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																	false	"Page number"		default(1)
//	@Param			per_page	query		int																	false	"Items per page"	default(10)
//	@Param			search		query		string																false	"Search term access_list_name and key"
//	@Param			ussd_enabled		query		string														false	"Search term"
//	@Param			enabled		query		string																false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=bulk_services_paginated_resp}	"Bulk services retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}								"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}								"Internal server error"
//	@Security		BearerAuth
//	@Router			/bulk_services [get]
func (h *bulk_serviceAdapter) GetAllBulkServices(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "", "getAllBulkServices", "handler", "bulkService")
	defer span.End()
	filter_params := util.ExtractFilterParams(r)
	bulk_services, err := h.bulkService.GetAllBulkServices(ctx, filter_params)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAllBulkServices] service error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.UnableToFetchBulkService.Code)
		return
	}
	span.SetAttributes(attribute.Int("bulk_service.count", len(bulk_services.Data)))
	h.logger.Infof("[GetAllBulkServices] retrieved %d bulk services", len(bulk_services.Data))
	localization.SendSuccessResponse(w, localization.BulkServiceFetchSuccessfully, bulk_services)

}

// EnableBulkService enables one or more bulk services
//
//	@Summary		Enable bulk services
//	@Description	Enables one or more bulk services by their keys
//	@Tags			Bulk Services
//	@Accept			json
//	@Produce		json
//	@Param			request	body		bulk_service.BulkServiceDTO				true	"Bulk service enable request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Bulk services enabled successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/bulk_services/enable [post]
func (h *bulk_serviceAdapter) EnableBulkService(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "", "enableBulkService", "handler", "bulkService")
	defer span.End()
	var req dto.BulkServiceDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[EnableBulkService] failed to decode request: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	span.SetAttributes(attribute.Int("bulk_service.keys_count", len(req.Keys)))
	err := h.bulkService.EnableBulkService(ctx, req.Keys)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[EnableBulkService] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[EnableBulkService] request sent successfully for %d service keys", len(req.Keys))
	localization.SendSuccessResponse(w, localization.BulkServiceEnableRequestSuccess, nil)
}

// DisableBulkService disables one or more bulk services
//
//	@Summary		Disable bulk services
//	@Description	Disables one or more bulk services by their keys
//	@Tags			Bulk Services
//	@Accept			json
//	@Produce		json
//	@Param			request	body		bulk_service.BulkServiceDTO				true	"Bulk service disable request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Bulk services disabled successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/bulk_services/disable [post]
func (h *bulk_serviceAdapter) DisableBulkService(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "", "disableBulkService", "handler", "bulkService")
	defer span.End()
	var req dto.BulkServiceDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DisableBulkService] failed to decode request: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	span.SetAttributes(attribute.Int("bulk_service.keys_count", len(req.Keys)))
	err := h.bulkService.DisableBulkService(ctx, req.Keys)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DisableBulkService] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[DisableBulkService] request sent successfully for %d service keys", len(req.Keys))
	localization.SendSuccessResponse(w, localization.BulkServiceDisableRequestSuccess, nil)
}
