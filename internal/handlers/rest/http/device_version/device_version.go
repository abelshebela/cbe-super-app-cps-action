package deviceversionhandler

import (
	dvdto "cbe-super-app-cps-action/internal/constants/dto/device_version"
	deviceversion "cbe-super-app-cps-action/internal/constants/interfaces/device_version"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type deviceVersionAdapter struct {
	svc    service.DeviceVersionServiceSrv
	logger utils.Logger
}

func InitDeviceVersionAdapter(s service.DeviceVersionServiceSrv, logger utils.Logger) deviceversion.DeviceVersionHandler {
	return &deviceVersionAdapter{svc: s, logger: logger}
}

// CreateDeviceVersion
//
//	@Summary		Create Device Version
//	@Description	Create a new device version with the provided information
//	@Tags			Device Version
//	@Accept			json
//	@Produce		json
//	@Param			device_version	body	dvdto.CreateDeviceVersionRequest	true	"Device version information"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Device version created successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/device_versions [post]
func (h *deviceVersionAdapter) CreateDeviceVersion(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "", "createDeviceVersion", "handler", "deviceVersion")
	defer span.End()
	var req dvdto.CreateDeviceVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateDeviceVersion] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}
	req.Clean()
	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateDeviceVersion] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	req.Platform = strings.ToUpper(req.Platform)
	span.SetAttributes(
		attribute.String("device_version.platform", req.Platform),
		attribute.String("device_version.version", req.LatestVersion),
	)
	if err := h.svc.CreateDeviceVersion(ctx, req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateDeviceVersion] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[CreateDeviceVersion] request sent successfully for platform: %s", req.Platform)
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_CREATE_REQUEST_CREATED", StatusCode: localization.StatusOK, Message: "Device version create request submitted", Type: "success"}, nil)
}

// UpdateDeviceVersion
//
//	@Summary		Update Device Version
//	@Description	Update an existing device version by its ID
//	@Tags			Device Version
//	@Accept			json
//	@Produce		json
//	@Param			id				path	string								true	"Device version ID"
//	@Param			device_version	body	dvdto.UpdateDeviceVersionRequest	true	"Device version update information"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Device version update request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/device_versions/{id} [patch]
func (h *deviceVersionAdapter) UpdateDeviceVersion(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "", "updateDeviceVersion", "handler", "deviceVersion")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	var req dvdto.UpdateDeviceVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateDeviceVersion] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}
	req.ID = id
	req.Clean()
	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateDeviceVersion] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	span.SetAttributes(
		attribute.String("device_version.id", id),
		attribute.String("device_version.platform", req.Platform),
		attribute.String("device_version.version", req.LatestVersion),
	)
	if err := h.svc.UpdateDeviceVersion(ctx, id, req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateDeviceVersion] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[UpdateDeviceVersion] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_UPDATE_REQUEST_CREATED", StatusCode: localization.StatusOK, Message: "Device version update request submitted", Type: "success"}, nil)
}

// GetAllDeviceVersions
//
//	@Summary		Get All Device Versions
//	@Description	Retrieve all device versions with pagination, filtering, and search. Searchable fields: latest_version, platform, enabled, force_update, created_at, last_modified, updated_at, updated_by, created_by.
//	@Tags			Device Version
//	@Accept			json
//	@Produce		json
//	@Param			page			query	int		false	"Page number"
//	@Param			per_page		query	int		false	"Items per page"
//	@Param			search			query	string	false	"Search term (searches latest_version,platform,enabled,updated_at,created_by,created_at,lastmodified_at)"
//	@Param			latest_version	query	string	false	"Filter by latest version"
//	@Param			platform		query	string	false	"Filter by platform (e.g., ANDROID, IOS)"
//	@Param			enabled			query	bool	false	"Filter by enabled status"
//	@Param			force_update	query	bool	false	"Filter by force update status"
//	@Param			created_at		query	string	false	"Filter by created at"
//	@Param			last_modified	query	string	false	"Filter by last modified"
//	@Param			updated_at		query	string	false	"Filter by updated at"
//	@Param			updated_by		query	string	false	"Filter by updated by"
//	@Param			created_by		query	string	false	"Filter by created by"
//	@Success		200	{object}	localization.StandardResponse{data=[]model.DeviceVersionControl}	"Device versions retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/device_versions [get]
func (h *deviceVersionAdapter) GetAllDeviceVersions(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "", "getAllDeviceVersions", "handler", "deviceVersion")
	defer span.End()
	filter := common_utils.ExtractFilterParams(r)
	res, err := h.svc.GetAllDeviceVersions(ctx, filter)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAllDeviceVersions] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("device_version.count", len(res.Data)))
	h.logger.Infof("[GetAllDeviceVersions] retrieved %d device versions", len(res.Data))
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSIONS_FETCHED", StatusCode: localization.StatusOK, Message: "Device versions fetched", Type: "success"}, res)
}

// GetDeviceVersionByID
//
//	@Summary		Get Device Version by ID
//	@Description	Retrieve a device version by its ID
//	@Tags			Device Version
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Device version ID"
//	@Success		200	{object}	localization.StandardResponse{data=model.DeviceVersionControl}	"Device version retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Device version not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/device_versions/{id} [get]
func (h *deviceVersionAdapter) GetDeviceVersionByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "", "getDeviceVersionById", "handler", "deviceVersion")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("device_version.id", id))
	res, err := h.svc.GetDeviceVersionByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetDeviceVersionByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[GetDeviceVersionByID] device version retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_FETCHED", StatusCode: localization.StatusOK, Message: "Device version fetched", Type: "success"}, res)
}

// Enable
//
//	@Summary		Enable Device Version
//	@Description	Enable a device version by its ID
//	@Tags			Device Version
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Device version ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Device version enable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Device version not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/device_versions/enable/{id} [patch]
func (h *deviceVersionAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "", "enableDeviceVersion", "handler", "deviceVersion")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("device_version.id", id))
	if err := h.svc.EnableDisableDeviceVersion(ctx, id, true); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[Enable] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_ENABLE_REQUEST_CREATED", StatusCode: localization.StatusOK, Message: "Device version enable request submitted", Type: "success"}, nil)
}

// Disable
//
//	@Summary		Disable Device Version
//	@Description	Disable a device version by its ID
//	@Tags			Device Version
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Device version ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Device version disable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Device version not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/device_version/disable/{id} [patch]
func (h *deviceVersionAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "", "disableDeviceVersion", "handler", "deviceVersion")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("device_version.id", id))
	if err := h.svc.EnableDisableDeviceVersion(ctx, id, false); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[Disable] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_DISABLE_REQUEST_CREATED", StatusCode: localization.StatusOK, Message: "Device version disable request submitted", Type: "success"}, nil)
}
