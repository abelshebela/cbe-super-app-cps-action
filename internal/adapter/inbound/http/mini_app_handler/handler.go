package miniapphandler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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

func (h *HttpStore) getValues(r *http.Request, isCreate bool) (*MiniAppRequest, error) {
	var req MiniAppRequest

	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "app_icon", 2<<20)
	if err != nil && err.Error() != common_util.ErrMissingFile && !isCreate {
		h.logger.Errorf("error parsing file: %v", err)
		return nil, err
	}
	if file != nil {
		defer file.Close()
	}

	get := func(key string) string {
		return strings.TrimSpace(r.FormValue(key))
	}

	parseBool := func(key string) (bool, error) {
		value := strings.ToLower(get(key))
		switch value {
		case "true":
			return true, nil
		case "false", "":
			return false, nil
		default:
			return false, errors.New("invalid boolean value for " + key)
		}
	}

	req.AppName = get("app_name")
	req.CommissionGLAccount = get("commission_gl_account")
	req.MerchantID = get("merchant_id")
	req.AppType = get("app_type")
	req.URL = get("url")
	req.MPAASID = get("mpaas_id")
	req.AppViewType = get("app_view_type")

	req.IFBProductCode = get("ifb_product_code")
	req.IFBVATCode = get("ifb_vat_code")
	req.IFBServiceFeeCode = get("ifb_service_fee_code")
	req.CBProductCode = get("cb_product_code")
	req.CBVATCode = get("cb_vat_code")
	req.CBServiceFeeCode = get("cb_service_fee_code")

	isEventMiniApp, err := parseBool("is_event_mini_app")
	if err != nil {
		return nil, err
	}
	req.IsEventMiniApp = isEventMiniApp

	isThreeClick, err := parseBool("is_three_click")
	if err != nil {
		return nil, err
	}
	req.IsThreeClick = isThreeClick
	req.AppIcon = fileHeader

	return &req, nil
}

func (h *HttpStore) MakerCreateMiniApp(w http.ResponseWriter, r *http.Request) {
	makerUser := contexts.ExtractUserContext(r)
	if makerUser.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, http.StatusBadRequest, nil)
		return
	}

	req, err := h.getValues(r, true)
	if err != nil {
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(true); err != nil {
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	var maker entities.User
	maker.FullName = makerUser.FullName
	maker.UserCode = makerUser.UserCode
	maker.PhoneNumber = makerUser.PhoneNumber
	maker.Department = makerUser.Department

	dto, err := req.ToMiniAppCreateRequest(true)
	if err != nil {
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	_, err = h.Application.MakerCreateMiniApp(r.Context(), dto, maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Mini App request successfully created")
}

func (h *HttpStore) MakerUpdateMiniApp(w http.ResponseWriter, r *http.Request) {

	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
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

	req, err := h.getValues(r, false)

	if err != nil {
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(false); err != nil {
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	dto, err := req.ToMiniAppCreateRequest(false)
	if err != nil {
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	dto.ID = id
	if IsEmptyUpdate(dto) {
		utils.SendErrorResponse(w, common_util.NoDataProvidedForUpdate, http.StatusBadRequest, nil)
		return
	}

	// Call domain application logic
	_, err = h.Application.MakerUpdateMiniApp(r.Context(), dto, maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Mini App update request created successfully")
}

func (h *HttpStore) MakerDeleteMiniApp(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	maker, UserErr := h.createUser(r)
	if UserErr != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := h.Application.MakerDeleteMiniApp(r.Context(), *maker, id)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, nil, "successful")
}

func (h *HttpStore) ListMiniApp(w http.ResponseWriter, r *http.Request) {
	filterParam := common_util.ExtractFilterParams(r)

	list, err := h.Application.ListMiniApp(r.Context(), filterParam)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	var docs []MiniAppResponse
	for _, doc := range list.Data {
		res := ToMiniAppResponse(doc)
		docs = append(docs, res)
	}

	res := common_util.PaginatedResponse[*[]MiniAppResponse]{
		Data: &docs,
		Meta: list.Meta,
	}
	utils.WriteSuccessResponse(w, res, "successful")
}

func (h *HttpStore) DetailMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.SendErrorResponse(w, common_util.InvalidID, 0, nil)
		return
	}
	detail, err := h.Application.DetailMiniAppByID(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, detail, "successful feached miniapp by ID")
}

func (h *HttpStore) EnableMiniAppByID(w http.ResponseWriter, r *http.Request) {
	h.enableDisableMiniApp(w, r, true)
}

func (h *HttpStore) DisableMiniAppByID(w http.ResponseWriter, r *http.Request) {
	h.enableDisableMiniApp(w, r, false)
}

func (h *HttpStore) enableDisableMiniApp(w http.ResponseWriter, r *http.Request, enable bool) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

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

	_, err := h.Application.EnableDisableMiniAppByID(r.Context(), id, enable, maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	action := "disabled"
	if enable {
		action = "enabled"
	}

	utils.WriteSuccessResponse(w, nil, fmt.Sprintf("Mini App request to be %s successfully created", action))
}
