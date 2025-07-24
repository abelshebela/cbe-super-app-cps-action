package miniapphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	miniapp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/miniapp"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
)

type HttpStore struct {
	Application miniapp_application.ApplicationAbstracts
	logger      util.Logger
}

func NewMiniAppAdapter(app miniapp_application.ApplicationAbstracts, logger util.Logger) Inbound.MiniAppInbound {
	return &HttpStore{
		Application: app,
		logger:      logger,
	}
}

func (h *HttpStore) createUser(r *http.Request) (*entities.User, error) {
	ctx_extract := ctx_util.ExtractUserContext(r)
	if ctx_extract.IsIncomplete() {
		return nil, fmt.Errorf("Unauthrozed")
	}
	maker := &entities.User{
		UserCode:    ctx_extract.UserID,
		FullName:    ctx_extract.FullName,
		PhoneNumber: ctx_extract.PhoneNumber,
		Department:  ctx_extract.Department,
	}

	return maker, nil
}

func (h *HttpStore) MakerCreateMiniApp(w http.ResponseWriter, r *http.Request) {
	var req dto.MiniAppCreateRequest
	makerUser := contexts.ExtractUserContext(r)
	if makerUser.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}

	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "app_icon", 10<<20)
	if err != nil {
		h.logger.Errorf("error parsing file: %v", err)
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	defer file.Close()

	isEventMiniApp, err := strconv.ParseBool(r.FormValue("is_event_mini_app"))
	if err != nil {
		isEventMiniApp = false
	}

	isThreeClick, err := strconv.ParseBool(r.FormValue("is_three_click"))
	if err != nil {
		isThreeClick = false
	}

	enabled, err := strconv.ParseBool(r.FormValue("enabled"))
	if err != nil {
		enabled = false
	}

	req.AppIcon = fileHeader
	req.AppName = r.FormValue("app_name")
	req.CommisonGLAccount = r.FormValue("commison_gl_account")
	req.MerchantID = r.FormValue("merchant_id")
	req.IsEventMiniApp = isEventMiniApp
	req.IsThreeClick = isThreeClick
	req.Enabled = enabled

	var appType dto.AppType
	_ = json.Unmarshal([]byte(r.FormValue("app_type")), &appType)

	var productCodes []dto.ProductCode
	_ = json.Unmarshal([]byte(r.FormValue("product_code")), &productCodes)

	var credentials []dto.CredentialInformation
	_ = json.Unmarshal([]byte(r.FormValue("credential")), &credentials)

	req.ProductCode = productCodes
	req.Credential = credentials
	req.AppType = appType

	fmt.Println(req)

	if err := req.Validate(); err != nil {
		common_util.SendErrorResponse(w, err, 0, nil)
		return
	}
	var maker entities.User

	maker.FullName = makerUser.FullName
	maker.UserCode = makerUser.UserCode
	maker.PhoneNumber = makerUser.PhoneNumber
	maker.Department = makerUser.Department

	response, err := h.Application.MakerCreateMiniApp(r.Context(), &req, maker)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to create mini app")
		return
	}

	utils.WriteSuccessResponse(w, response, "mini App request successfully created")
}

func (h *HttpStore) MakerUpdateMiniApp(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var req dto.MiniAppCreateRequest

	// Try to parse optional file
	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "app_icon", 10<<20)
	if err != nil && err.Error() != common_util.ErrMissingFile {
		h.logger.Errorf("error parsing file: %v", err)
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	if file != nil {
		defer file.Close()
		req.AppIcon = fileHeader
	}

	// Read values from form
	req.ID = id
	req.AppName = r.FormValue("app_name")
	req.CommisonGLAccount = r.FormValue("commison_gl_account")
	req.MerchantID = r.FormValue("merchant_id")

	// Parse optional bools
	if val := r.FormValue("is_event_mini_app"); val != "" {
		if parsed, err := strconv.ParseBool(val); err == nil {
			req.IsEventMiniApp = parsed
		}
	}
	if val := r.FormValue("is_three_click"); val != "" {
		if parsed, err := strconv.ParseBool(val); err == nil {
			req.IsThreeClick = parsed
		}
	}
	if val := r.FormValue("enabled"); val != "" {
		if parsed, err := strconv.ParseBool(val); err == nil {
			req.Enabled = parsed
		}
	}

	// Parse complex fields if provided
	if val := r.FormValue("app_type"); val != "" {
		var appType dto.AppType
		if err := json.Unmarshal([]byte(val), &appType); err == nil {
			req.AppType = appType
		}
	}
	if val := r.FormValue("product_code"); val != "" {
		var productCodes []dto.ProductCode
		if err := json.Unmarshal([]byte(val), &productCodes); err == nil {
			req.ProductCode = productCodes
		}
	}
	if val := r.FormValue("credential"); val != "" {
		var credentials []dto.CredentialInformation
		if err := json.Unmarshal([]byte(val), &credentials); err == nil {
			req.Credential = credentials
		}
	}

	// Check if at least one field was provided
	if req.AppName == "" && req.CommisonGLAccount == "" && req.MerchantID == "" &&
		len(req.ProductCode) == 0 && len(req.Credential) == 0 &&
		req.AppIcon == nil && req.AppType.UAT == "" && req.AppType.Production == "" &&
		req.AppType.Test == "" && req.AppType.Dev == "" &&
		!req.IsEventMiniApp && !req.IsThreeClick && !req.Enabled {
		common_util.SendErrorResponse(w, "NO_DATA_PROVIDED_FOR_UPDATE", http.StatusBadRequest, nil)
		return
	}

	// Extract user
	makerUser := contexts.ExtractUserContext(r)
	if makerUser.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}

	maker := entities.User{
		UserCode:    makerUser.UserCode,
		FullName:    makerUser.FullName,
		PhoneNumber: makerUser.PhoneNumber,
		Department:  makerUser.Department,
	}

	// Call domain application logic
	response, err := h.Application.MakerUpdateMiniApp(r.Context(), &req, maker)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to update mini app")
		return
	}

	utils.WriteSuccessResponse(w, response, "Mini App update request created successfully")
}

func (h *HttpStore) MakerDeleteMiniApp(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	maker, UserErr := h.createUser(r)
	if UserErr != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ActionCode, err := h.Application.MakerDeleteMiniApp(r.Context(), *maker, id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to delete mini app")
		return
	}
	utils.WriteSuccessResponse(w, ActionCode, "successful")
}

func (h *HttpStore) ListMiniApp(w http.ResponseWriter, r *http.Request) {
	filterParam := common_util.ExtractFilterParams(r)

	list, err := h.Application.ListMiniApp(r.Context(), filterParam)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to list mini apps")
		return
	}
	utils.WriteSuccessResponse(w, list, "successful")
}

func (h *HttpStore) DetailMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Missing mini app ID")
		return
	}
	detail, err := h.Application.DetailMiniAppByID(r.Context(), id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get mini app detail")
		return
	}
	utils.WriteSuccessResponse(w, detail, "successful feached miniapp by ID")
}
