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
	var reqDTO CreateMiniAppMerchantDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidInput, 0, nil)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		common_util.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	miniApp := ToMiniAppMerchantDomainFromCreateDTO(&reqDTO)
	cpsReq := entities.CreateCPSAction{
		MakerUser:  *maker,
		ActionData: miniApp,
	}

	_, err = h.Application.CreateOne(r.Context(), cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Create MiniAppMerchant request submitted")
}

func (h *HttpStore) UpdateMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var reqDTO UpdateMiniAppMerchantDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := reqDTO.Validate(); err != nil {
		common_util.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.Unauthorized, http.StatusUnauthorized, nil)
		return
	}

	miniApp := ToMiniAppMerchantDomainFromUpdateDTO(&reqDTO, id)
	cpsReq := entities.CreateCPSAction{
		MakerUser:  *maker,
		ActionData: miniApp,
	}

	_, err = h.Application.UpdateOne(r.Context(), cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Update request submitted")
}

func (h *HttpStore) DeleteMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	_, err = h.Application.DeleteOne(r.Context(), id, entities.CreateCPSAction{MakerUser: *maker})
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Delete request submitted")
}

func (h *HttpStore) Enable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	_, err = h.Application.EnableOne(r.Context(), id, entities.CreateCPSAction{MakerUser: *maker})
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Enable request submitted")
}

func (h *HttpStore) Disable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	maker, err := h.extractMaker(r)
	if err != nil {
		common_util.SendErrorResponse(w, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	_, err = h.Application.DisableOne(r.Context(), id, entities.CreateCPSAction{MakerUser: *maker})
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Disable request submitted")
}

func (h *HttpStore) GetMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	if id == "" {
		common_util.SendErrorResponse(w, "missing id", http.StatusBadRequest, nil)
		return
	}

	result, err := h.Application.Detail(r.Context(), id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, result, "Fetch successful")
}

func (h *HttpStore) GetAllMiniAppMerchant(w http.ResponseWriter, r *http.Request) {
	filter := common_util.ExtractFilterParams(r)

	result, err := h.Application.List(r.Context(), filter)
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
