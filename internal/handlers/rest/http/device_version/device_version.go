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
)

type deviceVersionAdapter struct {
	svc    service.DeviceVersionServiceSrv
	logger utils.Logger
}

func InitDeviceVersionAdapter(s service.DeviceVersionServiceSrv, logger utils.Logger) deviceversion.DeviceVersionHandler {
	return &deviceVersionAdapter{svc: s, logger: logger}
}

// CreateDeviceVersion
func (h *deviceVersionAdapter) CreateDeviceVersion(w http.ResponseWriter, r *http.Request) {
	var req dvdto.CreateDeviceVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[CreateDeviceVersion] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	req.Clean()
	if err := req.Validate(); err != nil {
		h.logger.Errorf("[CreateDeviceVersion] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	req.Platform = strings.ToUpper(req.Platform)
	if err := h.svc.CreateDeviceVersion(r.Context(), req); err != nil {
		h.logger.Errorf("[CreateDeviceVersion] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_CREATE_REQUEST_CREATED", StatusCode: localization.StatusOK, Message: "Device version create request submitted", Type: "success"}, nil)
}

// UpdateDeviceVersion
func (h *deviceVersionAdapter) UpdateDeviceVersion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	var req dvdto.UpdateDeviceVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[UpdateDeviceVersion] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}
	req.ID = id
	req.Clean()
	if err := req.Validate(); err != nil {
		h.logger.Errorf("[UpdateDeviceVersion] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	if err := h.svc.UpdateDeviceVersion(r.Context(), id, req); err != nil {
		h.logger.Errorf("[UpdateDeviceVersion] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_UPDATE_REQUEST_CREATED", StatusCode: localization.StatusOK, Message: "Device version update request submitted", Type: "success"}, nil)
}

// GetAllDeviceVersions
func (h *deviceVersionAdapter) GetAllDeviceVersions(w http.ResponseWriter, r *http.Request) {
	filter := common_utils.ExtractFilterParams(r)
	res, err := h.svc.GetAllDeviceVersions(r.Context(), filter)
	if err != nil {
		h.logger.Errorf("[GetAllDeviceVersions] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSIONS_FETCHED", StatusCode: localization.StatusOK, Message: "Device versions fetched", Type: "success"}, res)
}

// GetDeviceVersionByID
func (h *deviceVersionAdapter) GetDeviceVersionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	res, err := h.svc.GetDeviceVersionByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[GetDeviceVersionByID] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_FETCHED", StatusCode: localization.StatusOK, Message: "Device version fetched", Type: "success"}, res)
}

// Enable
func (h *deviceVersionAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	if err := h.svc.EnableDisableDeviceVersion(r.Context(), id, true); err != nil {
		h.logger.Errorf("[Enable] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_ENABLE_REQUEST_CREATED", StatusCode: localization.StatusOK, Message: "Device version enable request submitted", Type: "success"}, nil)
}

// Disable
func (h *deviceVersionAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	if err := h.svc.EnableDisableDeviceVersion(r.Context(), id, false); err != nil {
		h.logger.Errorf("[Disable] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.ResponseCode{Code: "SUCCESS_DEVICE_VERSION_DISABLE_REQUEST_CREATED", StatusCode: localization.StatusOK, Message: "Device version disable request submitted", Type: "success"}, nil)
}
