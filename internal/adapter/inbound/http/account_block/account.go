package accountblock_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/go-chi/chi"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AccountBlockHandler struct {
	service account_block.ApplicationService
	logger  utils.Logger
}

func NewAccountBlockHandler(service account_block.ApplicationService, logger utils.Logger) *AccountBlockHandler {
	return &AccountBlockHandler{service: service, logger: logger}
}

func (h *AccountBlockHandler) FilterSingleBranches(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Region   string `json:"region"`
        District string `json:"district"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "Invalid JSON body"}
        resp.SendJSON()
        return
    }

    req.Region = strings.TrimSpace(req.Region)
    req.District = strings.TrimSpace(req.District)

    if req.Region == "" || req.District == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "region and district are required"}
        resp.SendJSON()
        return
    }

    if len(req.Region) < 3 || len(req.District) < 3 {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "region and district must be at least 3 characters"}
        resp.SendJSON()
        return
    }

    branches, err := h.service.FilterSingleBranches(r.Context(), req.Region, req.District)
    if err != nil {
        h.logger.Errorf("FilterSingleBranches failed: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
        resp.SendJSON()
        return
    }

    if branches == nil {
        branches = []action.Branch{}
    }

    resp := common.Response[map[string]any]{
        ResponseWriter: w,
        Status:         http.StatusOK,
        Data:           map[string]any{"branches": branches},
    }
    resp.SendJSON()
}
func (h *AccountBlockHandler) DisableSingleBranch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BranchCode string `json:"branch_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.BranchCode) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "branch_code is required"}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
		resp.SendJSON()
		return
	}

	if strings.TrimSpace(department) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
		resp.SendJSON()
		return
	}

	branch, err := h.service.GetBranchByCode(r.Context(), req.BranchCode)
	if err != nil {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "Branch not found"}
		resp.SendJSON()
		return
	}

	maker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Department:  department,
	}

	err = h.service.DisableSingleBranch(r.Context(), branch, maker)
	if err != nil {
		h.logger.Errorf("DisableSingleBranch failed: %v", err)
		if strings.Contains(err.Error(), "Duplicate key error") || strings.Contains(err.Error(), "already exists") {
			resp := common.Response[any]{ResponseWriter: w, Status: http.StatusConflict, Data: "A pending disable action already exists for this branch"}
			resp.SendJSON()
			return
		}
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[map[string]string]{
		ResponseWriter: w,
		Status:         http.StatusCreated,
		Data:           map[string]string{"message": "Branch is disabled and CPS action created"},
	}
	resp.SendJSON()
}
func (h *AccountBlockHandler) ApproveSingleBranchDisable(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string  `json:"action_id"`
		Approve  bool    `json:"approve"`
		Reason   *string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "action_id is required"}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
		resp.SendJSON()
		return
	}

	if strings.TrimSpace(department) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveSingleBranchDisable requested by user: %s (%s)", userID, fullName)
	err := h.service.ApproveSingleBranchDisable(r.Context(), req.ActionID, req.Approve, req.Reason)
	if err != nil {
		h.logger.Errorf("ApproveSingleBranchDisable failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}

	resp := common.Response[map[string]string]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           map[string]string{"message": "Action processed successfully"},
	}
	resp.SendJSON()
}
func (h *AccountBlockHandler) FilterMultipleBranches(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Region   string `json:"region"`
		District string `json:"district"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "Invalid JSON body"}
		resp.SendJSON()
		return
	}
	req.Region = strings.TrimSpace(req.Region)
	req.District = strings.TrimSpace(req.District)
	if req.Region == "" || req.District == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "region and district are required"}
		resp.SendJSON()
		return
	}
	if len(req.Region) < 3 || len(req.District) < 3 {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "region and district must be at least 3 characters"}
		resp.SendJSON()
		return
	}
	branches, err := h.service.FilterMultipleBranches(r.Context(), req.Region, req.District)
	if err != nil {
		h.logger.Errorf("FilterMultipleBranches failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[map[string]any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           map[string]any{"branches": branches},
	}
	resp.SendJSON()
}

func (h *AccountBlockHandler) DisableMultipleBranches(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BranchCodes []string `json:"branch_codes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.BranchCodes) == 0 {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "branch_codes are required"}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" ||
		strings.TrimSpace(department) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
		resp.SendJSON()
		return
	}

	var branches []action.Branch
	for _, code := range req.BranchCodes {
		branch, err := h.service.GetBranchByCode(r.Context(), code)
		if err != nil {
			resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "Branch not found: " + code}
			resp.SendJSON()
			return
		}
		branches = append(branches, branch)
	}

	maker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Department:  department,
	}

	err := h.service.DisableMultipleBranches(r.Context(), branches, maker)
	if err != nil {
		h.logger.Errorf("DisableMultipleBranches failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
		resp.SendJSON()
		return
	}

	resp := common.Response[map[string]string]{
		ResponseWriter: w,
		Status:         http.StatusCreated,
		Data:           map[string]string{"message": "Branches are disabled and CPS action created"},
	}
	resp.SendJSON()
}
func (h *AccountBlockHandler) ApproveBulkBranchesDisable(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string  `json:"action_id"`
		Approve  bool    `json:"approve"`
		Reason   *string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "action_id is required"}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
		resp.SendJSON()
		return
	}

	if strings.TrimSpace(department) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveBulkBranchesDisable requested by user: %s (%s)", userID, fullName)
	err := h.service.ApproveBulkBranchesDisable(r.Context(), req.ActionID, req.Approve, req.Reason)
	if err != nil {
		h.logger.Errorf("ApproveBulkBranchesDisable failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}

	resp := common.Response[map[string]string]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           map[string]string{"message": "Action processed successfully"},
	}
	resp.SendJSON()
}

func (h *AccountBlockHandler) BlockRegion(w http.ResponseWriter, r *http.Request) {
    var req struct {
        RegionCode string `json:"region_code"`
        RegionName string `json:"region_name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
        strings.TrimSpace(req.RegionCode) == "" || strings.TrimSpace(req.RegionName) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "region_code and name are required"}
        resp.SendJSON()
        return
    }
    userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
    fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
    phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
    department, _ := r.Context().Value(constant.ContextKey("department")).(string)
    if strings.TrimSpace(department) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
        resp.SendJSON()
        return
    }

    maker := action.User{
        UserID:      userID,
        FullName:    fullName,
        PhoneNumber: phoneNumber,
        Department:  department,
    }

    cpsAction := action.CPSAction{
        MakerID:          maker.UserID,
        MakerName:        maker.FullName,
        MakerPhoneNumber: maker.PhoneNumber,
        Department:       maker.Department,
    }

    err := h.service.BlockRegion(r.Context(), req.RegionCode, cpsAction)
    if err != nil {
        h.logger.Errorf("BlockRegion failed: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
        resp.SendJSON()
        return
    }

    resp := common.Response[map[string]string]{
        ResponseWriter: w,
        Status:         http.StatusCreated,
        Data:           map[string]string{"message": "Region block action created"},
    }
    resp.SendJSON()
}

func (h *AccountBlockHandler) UpdateRegion(w http.ResponseWriter, r *http.Request) {
    var req action.Region
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.RegionCode) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "region code is required"}
        resp.SendJSON()
        return
    }

    userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
    fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
    phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
    department, _ := r.Context().Value(constant.ContextKey("department")).(string)

    if strings.TrimSpace(userID) == "" ||
        strings.TrimSpace(fullName) == "" ||
        strings.TrimSpace(phoneNumber) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
        resp.SendJSON()
        return
    }
    if strings.TrimSpace(department) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
        resp.SendJSON()
        return
    }

    h.logger.Infof("UpdateRegion requested by user: %s (%s)", userID, fullName)

    err := h.service.UpdateRegion(r.Context(), req)
    if err != nil {
        h.logger.Errorf("UpdateRegion failed: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
        resp.SendJSON()
        return
    }
    resp := common.Response[map[string]string]{
        ResponseWriter: w,
        Status:         http.StatusOK,
        Data:           map[string]string{"message": "Region updated successfully"},
    }
    resp.SendJSON()
}
func (h *AccountBlockHandler) ApproveRegionBlock(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string  `json:"action_id"`
		Approve  bool    `json:"approve"`
		Reason   *string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "action_id is required"}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
		resp.SendJSON()
		return
	}
	if strings.TrimSpace(department) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveRegionBlock requested by user: %s (%s)", userID, fullName)
	maker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Department:  department,
	}
	err := h.service.ApproveRegionBlock(r.Context(), req.ActionID, req.Approve, req.Reason, maker)
	if err != nil {
		h.logger.Errorf("ApproveRegionBlock failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[map[string]string]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           map[string]string{"message": "Action processed successfully"},
	}
	resp.SendJSON()
}

func (h *AccountBlockHandler) GetRegionByCode(w http.ResponseWriter, r *http.Request) {
    regionCode := chi.URLParam(r, "code")
    if strings.TrimSpace(regionCode) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "region_code is required"}
        resp.SendJSON()
        return
    }
    region, err := h.service.GetRegionByCode(r.Context(), regionCode)
    if err != nil {
        h.logger.Errorf("GetRegionByCode failed: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "Region not found"}
        resp.SendJSON()
        return
    }
    resp := common.Response[map[string]any]{
        ResponseWriter: w,
        Status:         http.StatusOK,
        Data:           map[string]any{"region": region},
    }
    resp.SendJSON()
}
func (h *AccountBlockHandler) BlockDistrict(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DistrictID string `json:"district_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.DistrictID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "district_id is required"}
		resp.SendJSON()
		return
	}
	userPayload, ok := r.Context().Value(constant.ContextKey("user_payload")).(middleware.UserPayload)
	if !ok || strings.TrimSpace(userPayload.UserID) == "" || strings.TrimSpace(userPayload.FullName) == "" || strings.TrimSpace(userPayload.PhoneNumber) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
		resp.SendJSON()
		return
	}
	maker := action.CPSAction{
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
	}
	err := h.service.BlockDistrict(r.Context(), req.DistrictID, maker)
	if err != nil {
		h.logger.Errorf("BlockDistrict failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusCreated, Data: "District block action created"}
	resp.SendJSON()
}

func (h *AccountBlockHandler) GetDistrictByID(w http.ResponseWriter, r *http.Request) {
	districtID := chi.URLParam(r, "id")
	if strings.TrimSpace(districtID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "district_id is required"}
		resp.SendJSON()
		return
	}
	district, err := h.service.GetDistrictByID(r.Context(), districtID)
	if err != nil {
		h.logger.Errorf("GetDistrictByID failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "District not found"}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: district}
	resp.SendJSON()
}
func (h *AccountBlockHandler) ApproveBlockDistrict(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DistrictID string `json:"district_id"`
		Approve    bool   `json:"approve"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.DistrictID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "district_id is required"}
		resp.SendJSON()
		return
	}
	userPayload, ok := r.Context().Value(constant.ContextKey("user_payload")).(middleware.UserPayload)
	if !ok {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
		resp.SendJSON()
		return
	}
	h.logger.Infof("ApproveBlockDistrict requested by user: %s (%s)", userPayload.UserID, userPayload.FullName)
	maker := action.CPSAction{
		CheckerID:          userPayload.UserID,
		CheckerName:        userPayload.FullName,
		CheckerPhoneNumber: userPayload.PhoneNumber,
		Department:         userPayload.Department,
	}
	err := h.service.ApproveBlockDistrict(r.Context(), req.DistrictID, maker)
	if err != nil {
		h.logger.Errorf("ApproveBlockDistrict failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: "Action processed successfully"}
	resp.SendJSON()
}
func (h *AccountBlockHandler) BlockCity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CityID string `json:"city_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.CityID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "city_id is required"}
		resp.SendJSON()
		return
	}
	userPayload, ok := r.Context().Value(constant.ContextKey("user_payload")).(middleware.UserPayload)
	if !ok || strings.TrimSpace(userPayload.UserID) == "" || strings.TrimSpace(userPayload.FullName) == "" || strings.TrimSpace(userPayload.PhoneNumber) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
		resp.SendJSON()
		return
	}
	maker := action.CPSAction{
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
	}
	err := h.service.BlockCity(r.Context(), req.CityID, maker)
	if err != nil {
		h.logger.Errorf("BlockCity failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusCreated, Data: "City block action created"}
	resp.SendJSON()
}

func (h *AccountBlockHandler) GetCityByID(w http.ResponseWriter, r *http.Request) {
	cityID := chi.URLParam(r, "id")
	if strings.TrimSpace(cityID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "city_id is required"}
		resp.SendJSON()
		return
	}
	city, err := h.service.GetCityByID(r.Context(), cityID)
	if err != nil {
		h.logger.Errorf("GetCityByID failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "City not found"}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: city}
	resp.SendJSON()
}

func (h *AccountBlockHandler) ApproveBlockCity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CityID  string `json:"city_id"`
		Approve bool   `json:"approve"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.CityID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "city_id is required"}
		resp.SendJSON()
		return
	}
	userPayload, ok := r.Context().Value(constant.ContextKey("user_payload")).(middleware.UserPayload)
	if !ok {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
		resp.SendJSON()
		return
	}
	h.logger.Infof("ApproveBlockCity requested by user: %s (%s)", userPayload.UserID, userPayload.FullName)
	maker := action.CPSAction{
		CheckerID:          userPayload.UserID,
		CheckerName:        userPayload.FullName,
		CheckerPhoneNumber: userPayload.PhoneNumber,
		Department:         userPayload.Department,
	}
	err := h.service.ApproveBlockCity(r.Context(), req.CityID, maker)
	if err != nil {
		h.logger.Errorf("ApproveBlockCity failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: "Action processed successfully"}
	resp.SendJSON()
}
func (h *AccountBlockHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.UserID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "user_id is required"}
		resp.SendJSON()
		return
	}
	userPayload, ok := r.Context().Value(constant.ContextKey("user_payload")).(middleware.UserPayload)
	if !ok || strings.TrimSpace(userPayload.UserID) == "" || strings.TrimSpace(userPayload.FullName) == "" || strings.TrimSpace(userPayload.PhoneNumber) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
		resp.SendJSON()
		return
	}
	maker := action.CPSAction{
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
	}
	err := h.service.BlockUser(r.Context(), req.UserID, maker)
	if err != nil {
		h.logger.Errorf("BlockUser failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusCreated, Data: "User block action created"}
	resp.SendJSON()
}

func (h *AccountBlockHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if strings.TrimSpace(userID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "user_id is required"}
		resp.SendJSON()
		return
	}
	maker := action.CPSAction{}
	user, err := h.service.GetUserByID(r.Context(), userID, maker)
	if err != nil {
		h.logger.Errorf("GetUserByID failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "User not found"}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: user}
	resp.SendJSON()
}

func (h *AccountBlockHandler) ApproveBlockUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  string `json:"user_id"`
		Approve bool   `json:"approve"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.UserID) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "user_id is required"}
		resp.SendJSON()
		return
	}
	userPayload, ok := r.Context().Value(constant.ContextKey("user_payload")).(middleware.UserPayload)
	if !ok {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
		resp.SendJSON()
		return
	}
	h.logger.Infof("ApproveBlockUser requested by user: %s (%s)", userPayload.UserID, userPayload.FullName)
	maker := action.CPSAction{
		CheckerID:          userPayload.UserID,
		CheckerName:        userPayload.FullName,
		CheckerPhoneNumber: userPayload.PhoneNumber,
		Department:         userPayload.Department,
	}
	err := h.service.ApproveBlockUser(r.Context(), req.UserID, maker)
	if err != nil {
		h.logger.Errorf("ApproveBlockUser failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: "Action processed successfully"}
	resp.SendJSON()
}
