package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	servicesdto "cbe-super-app-cps-action/internal/constants/dto/services"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/services"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type (
	_servicesCreateReqForSwagger = servicesdto.CreateServiceRequest
	_servicesUpdateReqForSwagger = servicesdto.UpdateServiceRequest
)

type servicesAdapter struct {
	app    service.ServicesService
	logger utils.Logger
}

func InitServicesAdapter(app service.ServicesService, logger utils.Logger) inbound.ServicesHandler {
	return &servicesAdapter{app: app, logger: logger}
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}, logger utils.Logger) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		logger.Errorf("Failed to decode JSON request: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidJSONPayload, nil, nil)
		return false
	}
	return true
}

// Create godoc
//
//	@Summary		Create Service
//	@Description	Create a new service
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			request	body		servicesdto.CreateServiceRequest		true	"Service payload"
//	@Success		201		{object}	localization.StandardResponse{data=nil}	"CPS action created"
//	@Failure		400,500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/services [post]
func (a *servicesAdapter) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "createService", "handler", "services")
	defer span.End()

	var req servicesdto.CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("Failed to decode JSON request: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidJSONPayload, nil, nil)
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := a.app.Create(ctx, req); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		userCode, _ := r.Context().Value(constants.ContextKey("user_code")).(string)
		a.logger.Infof("[Create] request sent successfully for user_code: %s is_maker_only: %v", userCode, IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessServiceCreated, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceCreateRequestSubmitted, nil)
}

// Update godoc
//
//	@Summary		Update Service
//	@Description	Update an existing service
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string									true	"Service ID"
//	@Param			request		body		servicesdto.UpdateServiceRequest		true	"Service update payload"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS action created"
//	@Failure		400,404,500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/services/{id} [patch]
func (a *servicesAdapter) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "updateService", "handler", "services")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("service ID is required"))
		localization.SendErrorResponse(w, localization.ErrorInvalidID, nil, nil)
		return
	}

	var req servicesdto.UpdateServiceRequest
	if !decodeJSONBody(w, r, &req, a.logger) {
		span.RecordError(errors.New("invalid payload"))
		return
	}
	if err := req.Validate(); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("service.id", id))

	if err := a.app.Update(ctx, id, req); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		userCode, _ := r.Context().Value(constants.ContextKey("user_code")).(string)
		a.logger.Infof("[Update] request sent successfully for user_code: %s is_maker_only: %v", userCode, IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessServiceUpdated, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceUpdateRequestSubmitted, nil)
}

// Enable godoc
//
//	@Summary		Enable Service
//	@Description	Enable a service by ID
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string									true	"Service ID"
//	@Success		200				{object}	localization.StandardResponse{data=nil}	"CPS action created"
//	@Failure		400,404,409,500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/services/{id}/enable [patch]
func (a *servicesAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "enableService", "handler", "services")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("service ID is required for enable"))
		localization.SendErrorResponse(w, localization.ErrorInvalidID, nil, nil)
		return
	}

	if r.Header.Get("Authorization") == "" {
		span.RecordError(errors.New("missing Authorization header"))
		localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
		return
	}

	span.SetAttributes(attribute.String("service.id", id))
	if err := a.app.Enable(ctx, id); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		userCode, _ := r.Context().Value(constants.ContextKey("user_code")).(string)
		a.logger.Infof("[Enable] request sent successfully for user_code: %s is_maker_only: %v", userCode, IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessServiceEnabled, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceEnableRequestSubmitted, nil)
}

// Disable godoc
//
//	@Summary		Disable Service
//	@Description	Disable a service by ID
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string									true	"Service ID"
//	@Success		200				{object}	localization.StandardResponse{data=nil}	"CPS action created"
//	@Failure		400,404,409,500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/services/{id}/disable [patch]
func (a *servicesAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "disableService", "handler", "services")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("service ID is required for disable"))
		localization.SendErrorResponse(w, localization.ErrorInvalidID, nil, nil)
		return
	}

	if r.Header.Get("Authorization") == "" {
		span.RecordError(errors.New("missing Authorization header"))
		localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
		return
	}

	span.SetAttributes(attribute.String("service.id", id))
	if err := a.app.Disable(ctx, id); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		userCode, _ := r.Context().Value(constants.ContextKey("user_code")).(string)
		a.logger.Infof("[Disable] request sent successfully for user_code: %s is_maker_only: %v", userCode, IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessServiceDisabled, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceDisableRequestSubmitted, nil)
}

// GetAll godoc
//
//	@Summary		List Services
//	@Description	Retrieve services with pagination, filtering, and search.
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int		false	"Page number"		default(1)
//	@Param			per_page		query		int		false	"Items per page"	default(10)
//	@Param			service_name	query		string	false	"Filter by service_name"
//	@Param			service_code	query		string	false	"Filter by service_code"
//	@Param			service_type	query		string	false	"Filter by service_type"
//	@Param			enabled			query		bool	false	"Filter by enabled status"
//	@Param			search			query		string	false	"Search term (service_name, service_code, service_type)"
//	@Success		200				{object}	localization.StandardResponse{data=types.PaginatedResponse}
//	@Failure		500				{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/services [get]
func (a *servicesAdapter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getAllServices", "handler", "services")
	defer span.End()

	filter := local_util.ExtractFilterParams(r)
	list, err := a.app.GetAll(ctx, *filter)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("service.count", len(list.Data)))
	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, list)
}

// GetAllServiceList godoc
//
//	@Summary		List Service List
//	@Description	Retrieve service list with pagination, filtering, and search.
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int		false	"Page number"		default(1)
//	@Param			per_page		query		int		false	"Items per page"	default(10)
//	@Param			service_name	query		string	false	"Filter by service_name"
//	@Param			service_code	query		string	false	"Filter by service_code"
//	@Param			service_type	query		string	false	"Filter by service_type"
//	@Param			enabled			query		bool	false	"Filter by enabled status"
//	@Param			search			query		string	false	"Search term (service_name, service_code, service_type)"
//	@Success		200				{object}	localization.StandardResponse{data=types.PaginatedResponse}
//	@Failure		500				{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/services/list [get]
func (a *servicesAdapter) GetAllServiceList(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getAllServicesList", "handler", "servicesList")
	defer span.End()

	filter := local_util.ExtractFilterParams(r)
	list, err := a.app.GetAllServiceList(ctx, *filter)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("service_list.count", len(list.Data)))
	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, list)
}

// GetByID godoc
//
//	@Summary		Get Service by ID
//	@Description	Retrieve a service by its ID
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"Service ID"
//	@Success		200			{object}	localization.StandardResponse{data=entities.Services}
//	@Failure		400,404,500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/services/{id} [get]
func (a *servicesAdapter) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getServiceById", "handler", "services")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("service ID is required for get"))
		localization.SendErrorResponse(w, localization.ErrorInvalidID, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("service.id", id))
	item, err := a.app.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, item)
}
