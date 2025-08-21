package miniappmerchant

import (
	"encoding/json"
	"fmt"
	"net/http"

	miniapp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app_merchant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type HttpStore struct {
	Application miniapp_application.MiniAppMerchantApplication
	logger      shared_utils.Logger
}

func NewMiniAppMerchantAdapter(app miniapp_application.MiniAppMerchantApplication, logger shared_utils.Logger) *HttpStore {
	return &HttpStore{Application: app, logger: logger}
}

func (h *HttpStore) extractMaker(r *http.Request) (*entities.User, error) {
	u := ctx_util.ExtractUserContext(r)
	if u.IsIncomplete() {
		return nil, fmt.Errorf("unauthorized user context")
	}
	return &entities.User{
		UserCode:    u.UserID,
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
		Department:  u.Department,
	}, nil
}

func (h *HttpStore) CreateMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	var reqDTO MiniAppMerchantDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := reqDTO.Validate(true); err != nil {
		common_util.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.Unauthorized, http.StatusUnauthorized, nil)
		return
	}

	merchantReq := ToMiniAppMerchantDomainFromUpdateDTO(&reqDTO)
	err = h.Application.CreateMerchant(r.Context(), *merchantReq, *maker)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.BaseResponseMaker(nil, w, "Create MiniAppMerchant request successfully created", 201)
}

func (h *HttpStore) UpdateMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok || id == "" {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	var reqDTO MiniAppMerchantDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if reqDTO.IsEmpty() {
		common_util.SendErrorResponse(w, common_util.NoDataProvidedForUpdate, http.StatusBadRequest, nil)
		return
	}

	if err := reqDTO.Validate(false); err != nil {
		common_util.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.Unauthorized, http.StatusUnauthorized, nil)
		return
	}

	merchantReq := ToMiniAppMerchantDomainFromUpdateDTO(&reqDTO)
	err = h.Application.UpdateMerchant(r.Context(), id, *merchantReq, *maker)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Update request successfully created")
}

func (h *HttpStore) DeleteMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok || id == "" {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.Unauthorized, http.StatusUnauthorized, nil)
		return
	}

	err = h.Application.DeleteMerchant(r.Context(), id, *maker)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Delete request successfully created")
}

func (h *HttpStore) Enable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok || id == "" {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}
	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.Unauthorized, http.StatusUnauthorized, nil)
		return
	}

	err = h.Application.EnableDisableMerchant(r.Context(), id, *maker, true)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Enable request successfully created")
}

func (h *HttpStore) Disable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok || id == "" {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}
	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.Unauthorized, http.StatusUnauthorized, nil)
		return
	}

	err = h.Application.EnableDisableMerchant(r.Context(), id, *maker, false)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Disable request successfully created")
}

func (h *HttpStore) GetMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok || id == "" {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	result, err := h.Application.FetchMerchantByID(r.Context(), id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	response := ToMiniAppMerchantResponseDTO(result)
	common_util.WriteSuccessResponse(w, response, "Fetch successful")
}

func (h *HttpStore) GetAllMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	filter := common_util.ExtractFilterParams(r)

	result, err := h.Application.FetchMerchants(r.Context(), filter)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	var response []*MiniAppMerchantResponseDTO
	for _, doc := range result.Data {
		response = append(response, ToMiniAppMerchantResponseDTO(doc))
	}

	res := common_util.PaginatedResponse[[]*MiniAppMerchantResponseDTO]{
		Data: response,
		Meta: result.Meta,
	}

	common_util.WriteSuccessResponse(w, res, "List successful")
}
