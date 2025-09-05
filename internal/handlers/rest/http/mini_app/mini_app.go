package miniapphandler

import (
	miniappdto "cbe-super-app-cps-action/internal/constants/dto/mini_app"
	miniAppInbound "cbe-super-app-cps-action/internal/constants/interfaces/mini_app"
	"cbe-super-app-cps-action/internal/constants/localization"
	miniappcore "cbe-super-app-cps-action/internal/handlers/rest/http/mini_app/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
)

type HttpStore struct {
	miniAppSrvc service.MiniAppService
	logger      utils.Logger
}

func InitMiniAppAdapter(app service.MiniAppService, logger utils.Logger) miniAppInbound.MiniAppInbound {
	return &HttpStore{
		miniAppSrvc: app,
		logger:      logger,
	}
}

func (h *HttpStore) CreateMiniApp(w http.ResponseWriter, r *http.Request) {
	req, err := miniappcore.ParseMiniAppRequestFromMultipartForm(r, true)
	if err != nil {
		h.logger.Errorf("failed to parse mini app request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(true); err != nil {
		h.logger.Errorf("mini app request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	dto, err := miniappcore.ToMiniAppCreateRequest(req, true)
	if err != nil {
		h.logger.Errorf("failed to convert request to DTO: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := h.miniAppSrvc.CreateMiniApp(r.Context(), dto); err != nil {
		h.logger.Errorf("failed to create mini app: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app creation request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantCreateRequestCreated, nil)
}

func (h *HttpStore) UpdateMiniApp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	req, err := miniappcore.ParseMiniAppRequestFromMultipartForm(r, false)
	if err != nil {
		h.logger.Errorf("failed to parse mini app update request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	if err := req.Validate(false); err != nil {
		h.logger.Errorf("mini app request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	dtoForUpdate, err := miniappcore.ToMiniAppCreateRequest(req, false)
	if err != nil {
		h.logger.Errorf("failed to convert update request to DTO: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	if miniappdto.IsEmptyUpdate(dtoForUpdate) {
		h.logger.Warnf("no data provided for mini app update, mini app ID: %s", id)
		localization.SendErrorResponse(w, localization.ErrorUpdateMiniAppEmptyPayload, nil, nil)
		return
	}

	dto, err := miniappcore.ToMiniAppCreateRequest(req, false)
	if err != nil {
		h.logger.Errorf("failed to convert update request to DTO: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	dto.ID = id

	if err := h.miniAppSrvc.UpdateMiniApp(r.Context(), dto); err != nil {
		h.logger.Errorf("failed to update mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app update request submitted successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessUpdateMiniAppRequestCreated, nil)
}

func (h *HttpStore) DeleteMiniApp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for delete")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	if err := h.miniAppSrvc.DeleteMiniApp(r.Context(), id); err != nil {
		h.logger.Errorf("failed to delete mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app delete request submitted successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppDeletedRequestCreated, nil)
}

func (h *HttpStore) EnableMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for enable")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	if err := h.miniAppSrvc.EnableDisableMiniAppByID(r.Context(), id, true); err != nil {
		h.logger.Errorf("failed to enable mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app enable request submitted successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppEnableRequestSubmitted, nil)
}

func (h *HttpStore) DisableMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for disable")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	if err := h.miniAppSrvc.EnableDisableMiniAppByID(r.Context(), id, false); err != nil {
		h.logger.Errorf("failed to disable mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app disable request submitted successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppDisableRequestSubmitted, nil)
}

func (h *HttpStore) DetailMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for fetch")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	miniApp, err := h.miniAppSrvc.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("failed to fetch mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app fetched successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppFetchedByID, miniappcore.ToMiniAppResponse(miniApp))
}

func (h *HttpStore) ListMiniApp(w http.ResponseWriter, r *http.Request) {
	filter := local_util.ExtractFilterParams(r)
	h.logger.Infof("fetching mini apps with filter: %+v", filter)

	list, err := h.miniAppSrvc.ListMiniApp(r.Context(), *filter)
	if err != nil {
		h.logger.Errorf("failed to fetch mini apps: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini apps fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessMiniAppsRetrieved, list)
}
