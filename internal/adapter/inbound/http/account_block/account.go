package accountblock_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
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

	data, err := constant_utils.StructToMap(map[string]any{"branches": branches})
	if err != nil {
		constant_utils.SendErrorResponse(w, http.StatusInternalServerError, 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Branches retrieved successfully", http.StatusOK)
}
func (h *AccountBlockHandler) DisableSingleBranch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BranchCode string `json:"branch_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.BranchCode) == "" {
		constant_utils.BaseResponseMaker(nil, w, "branch_code is required", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {
		constant_utils.BaseResponseMaker(nil, w, "Incomplete user information", http.StatusUnauthorized)
		return
	}

	if strings.TrimSpace(department) == "" {
		constant_utils.BaseResponseMaker(nil, w, "department is required in context", http.StatusUnauthorized)
		return
	}

	branch, err := h.service.GetBranchByCode(r.Context(), req.BranchCode)
	if err != nil {
		constant_utils.BaseResponseMaker(nil, w, "Branch not found", http.StatusNotFound)
		return
	}

	maker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Department:  department,
	}

	actionCode, err := h.service.DisableSingleBranch(r.Context(), branch, maker)
	if err != nil {
		h.logger.Errorf("DisableSingleBranch failed: %v", err)
		if strings.Contains(err.Error(), "Duplicate key error") || strings.Contains(err.Error(), "already exists") {
			constant_utils.BaseResponseMaker(nil, w, "A pending disable action already exists for this branch", http.StatusConflict)
			return
		}
		constant_utils.BaseResponseMaker(nil, w, err.Error(), http.StatusInternalServerError)
		return
	}
	respData := map[string]any{
		"action_code": actionCode,
	}
	constant_utils.BaseResponseMaker(respData, w, "Branch is disabled and CPS action created", http.StatusCreated)
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

	data := map[string]any{
		"message": "Action processed successfully",
	}
	constant_utils.BaseResponseMaker(data, w, "Action processed successfully", http.StatusOK)
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

	constant_utils.BaseResponseMaker(map[string]any{"branches": branches}, w, "Branches retrieved successfully", http.StatusOK)
}
func (h *AccountBlockHandler) DisableMultipleBranches(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BranchCodes []string `json:"branch_codes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.BranchCodes) == 0 {
		constant_utils.BaseResponseMaker(nil, w, "branch_codes are required", http.StatusBadRequest)
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
		constant_utils.BaseResponseMaker(nil, w, "User info missing in context", http.StatusUnauthorized)
		return
	}

	var branches []action.Branch
	for _, code := range req.BranchCodes {
		branch, err := h.service.GetBranchByCode(r.Context(), code)
		if err != nil {
			constant_utils.BaseResponseMaker(nil, w, "Branch not found: "+code, http.StatusNotFound)
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

	actionCode, err := h.service.DisableMultipleBranches(r.Context(), branches, maker)
	if err != nil {
		h.logger.Errorf("DisableMultipleBranches failed: %v", err)
		constant_utils.BaseResponseMaker(nil, w, err.Error(), http.StatusInternalServerError)
		return
	}

	respData := map[string]any{
		"action_code": actionCode,
	}
	constant_utils.BaseResponseMaker(respData, w, "Branches are disabled and CPS action created", http.StatusCreated)
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

	data := map[string]any{
		"message": "Action processed successfully",
	}
	constant_utils.BaseResponseMaker(data, w, "Action processed successfully", http.StatusOK)
}

func (h *AccountBlockHandler) BlockRegion(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RegionCode string `json:"region_code"`
		RegionName string `json:"region_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		strings.TrimSpace(req.RegionCode) == "" || strings.TrimSpace(req.RegionName) == "" {
		constant_utils.BaseResponseMaker(nil, w, "region_code and name are required", http.StatusBadRequest)
		return
	}
	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		constant_utils.BaseResponseMaker(nil, w, "department is required in context", http.StatusUnauthorized)
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

	actionCode, err := h.service.BlockRegion(r.Context(), req.RegionCode, cpsAction)
	if err != nil {
		constant_utils.BaseResponseMaker(nil, w, err.Error(), http.StatusInternalServerError)
		return
	}
	respData := map[string]any{
		"action_code": actionCode,
	}
	constant_utils.BaseResponseMaker(respData, w, "Region block action created", http.StatusCreated)
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
	data := map[string]any{
		"message": "Region updated successfully",
	}
	constant_utils.BaseResponseMaker(data, w, "Region updated successfully", http.StatusOK)
}
func (h *AccountBlockHandler) ApproveRegionBlock(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string  `json:"action_id"`
		Approve  bool    `json:"approve"`
		Reason   *string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		h.logger.Errorf("ApproveRegionBlock: invalid request body or missing action_id, err=%v", err)
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
		h.logger.Errorf("ApproveRegionBlock: Incomplete user information in context: userID=%s, fullName=%s, phoneNumber=%s", userID, fullName, phoneNumber)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
		resp.SendJSON()
		return
	}
	if strings.TrimSpace(department) == "" {
		h.logger.Errorf("ApproveRegionBlock: department is required in context")
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveRegionBlock requested by user: %s (%s), actionID: %s, approve: %v", userID, fullName, req.ActionID, req.Approve)
	checker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Department:  department,
	}

	err := h.service.ApproveRegionBlock(r.Context(), req.ActionID, req.Approve, req.Reason, checker)
	if err != nil {
		h.logger.Errorf("ApproveRegionBlock failed for actionID=%s by user=%s: %v", req.ActionID, userID, err)
		if strings.Contains(err.Error(), "not allowed") {
			resp := common.Response[any]{ResponseWriter: w, Status: http.StatusForbidden, Data: err.Error()}
			resp.SendJSON()
			return
		}
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveRegionBlock succeeded for actionID=%s by user=%s", req.ActionID, userID)

	data := map[string]any{
		"message": "Action processed successfully",
	}
	constant_utils.BaseResponseMaker(data, w, "Action processed successfully", http.StatusOK)
}
func (h *AccountBlockHandler) GetRegionByCode(w http.ResponseWriter, r *http.Request) {
	regionCode := r.URL.Query().Get("region_code")
	if strings.TrimSpace(regionCode) == "" {
		constant_utils.BaseResponseMaker(nil, w, "region_code is required", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {
		constant_utils.BaseResponseMaker(nil, w, "Incomplete user information", http.StatusUnauthorized)
		return
	}
	if strings.TrimSpace(department) == "" {
		constant_utils.BaseResponseMaker(nil, w, "department is required in context", http.StatusUnauthorized)
		return
	}

	h.logger.Infof("GetRegionByCode requested by user: %s (%s)", userID, fullName)

	region, err := h.service.GetRegionByCode(r.Context(), regionCode)
	if err != nil {
		h.logger.Errorf("GetRegionByCode failed: %v", err)
		constant_utils.BaseResponseMaker(nil, w, "Region not found", http.StatusNotFound)
		return
	}

	data, err := constant_utils.StructToMap(region)
	if err != nil {
		constant_utils.SendErrorResponse(w, http.StatusInternalServerError, 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Region retrieved successfully", http.StatusOK)
}
func (h *AccountBlockHandler) BlockDistrict(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DistrictCode string `json:"district_code"`
		DistrictName string `json:"district_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		strings.TrimSpace(req.DistrictCode) == "" || strings.TrimSpace(req.DistrictName) == "" {
		constant_utils.BaseResponseMaker(nil, w, "district_code and name are required", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		constant_utils.BaseResponseMaker(nil, w, "department is required in context", http.StatusUnauthorized)
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

	actionCode, err := h.service.BlockDistrict(r.Context(), req.DistrictCode, cpsAction)
	if err != nil {
		constant_utils.BaseResponseMaker(nil, w, err.Error(), http.StatusInternalServerError)
		return
	}
	respData := map[string]any{
		"action_code": actionCode,
	}
	constant_utils.BaseResponseMaker(respData, w, "District block action created", http.StatusCreated)
}

func (h *AccountBlockHandler) GetDistrictByCode(w http.ResponseWriter, r *http.Request) {
	districtCode := r.URL.Query().Get("district_code")
	if strings.TrimSpace(districtCode) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "district_code is required"}
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

	h.logger.Infof("GetDistrictByCode requested by user: %s (%s)", userID, fullName)

	district, err := h.service.GetDistrictByCode(r.Context(), districtCode)
	if err != nil {
		h.logger.Errorf("GetDistrictByCode failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "District not found"}
		resp.SendJSON()
		return
	}

	data, err := constant_utils.StructToMap(district)
	if err != nil {
		constant_utils.SendErrorResponse(w, http.StatusInternalServerError, 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "District retrieved successfully", http.StatusOK)
}
func (h *AccountBlockHandler) ApproveBlockDistrict(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string  `json:"action_id"`
		Approve  bool    `json:"approve"`
		Reason   *string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		h.logger.Errorf("ApproveBlockDistrict: invalid request body or missing action_id, err=%v", err)
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
		h.logger.Errorf("ApproveBlockDistrict: Incomplete user information in context: userID=%s, fullName=%s, phoneNumber=%s", userID, fullName, phoneNumber)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
		resp.SendJSON()
		return
	}
	if strings.TrimSpace(department) == "" {
		h.logger.Errorf("ApproveBlockDistrict: department is required in context")
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveBlockDistrict requested by user: %s (%s), actionID: %s, approve: %v", userID, fullName, req.ActionID, req.Approve)
	checker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Department:  department,
	}

	err := h.service.ApproveBlockDistrict(r.Context(), req.ActionID, req.Approve, req.Reason, checker)
	if err != nil {
		h.logger.Errorf("ApproveBlockDistrict failed for actionID=%s by user=%s: %v", req.ActionID, userID, err)
		if strings.Contains(err.Error(), "not allowed") {
			resp := common.Response[any]{ResponseWriter: w, Status: http.StatusForbidden, Data: err.Error()}
			resp.SendJSON()
			return
		}
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveBlockDistrict succeeded for actionID=%s by user=%s", req.ActionID, userID)
	data := map[string]any{
		"message": "Action processed successfully",
	}
	constant_utils.BaseResponseMaker(data, w, "Action processed successfully", http.StatusOK)
}

func (h *AccountBlockHandler) GetCityByCode(w http.ResponseWriter, r *http.Request) {
	cityCode := r.URL.Query().Get("city_code")
	if strings.TrimSpace(cityCode) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "city_code is required"}
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

	h.logger.Infof("GetCityByCode requested by user: %s (%s)", userID, fullName)

	city, err := h.service.GetCityByCode(r.Context(), cityCode)
	if err != nil {
		h.logger.Errorf("GetCityByCode failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "City not found"}
		resp.SendJSON()
		return
	}

	data, err := constant_utils.StructToMap(city)
	if err != nil {
		constant_utils.SendErrorResponse(w, http.StatusInternalServerError, 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "City retrieved successfully", http.StatusOK)
}
func (h *AccountBlockHandler) BlockCity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CityCode string `json:"city_code"`
		CityName string `json:"city_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		strings.TrimSpace(req.CityCode) == "" || strings.TrimSpace(req.CityName) == "" {
		constant_utils.BaseResponseMaker(nil, w, "city_code and city_name are required", http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)
	if strings.TrimSpace(department) == "" {
		constant_utils.BaseResponseMaker(nil, w, "department is required in context", http.StatusUnauthorized)
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

	actionCode, err := h.service.BlockCity(r.Context(), req.CityCode, cpsAction)
	if err != nil {
		constant_utils.BaseResponseMaker(nil, w, err.Error(), http.StatusInternalServerError)
		return
	}
	respData := map[string]any{
		"action_code": actionCode,
	}
	constant_utils.BaseResponseMaker(respData, w, "City block action created", http.StatusCreated)
}
func (h *AccountBlockHandler) ApproveBlockCity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string  `json:"action_id"`
		Approve  bool    `json:"approve"`
		Reason   *string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		h.logger.Errorf("ApproveBlockCity: invalid request body or missing action_id, err=%v", err)
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
		h.logger.Errorf("ApproveBlockCity: Incomplete user information in context: userID=%s, fullName=%s, phoneNumber=%s", userID, fullName, phoneNumber)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
		resp.SendJSON()
		return
	}
	if strings.TrimSpace(department) == "" {
		h.logger.Errorf("ApproveBlockCity: department is required in context")
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveBlockCity requested by user: %s (%s), actionID: %s, approve: %v", userID, fullName, req.ActionID, req.Approve)
	checker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Department:  department,
	}

	err := h.service.ApproveBlockCity(r.Context(), req.ActionID, req.Approve, req.Reason, checker)
	if err != nil {
		h.logger.Errorf("ApproveBlockCity failed for actionID=%s by user=%s: %v", req.ActionID, userID, err)
		if strings.Contains(err.Error(), "not allowed") {
			resp := common.Response[any]{ResponseWriter: w, Status: http.StatusForbidden, Data: err.Error()}
			resp.SendJSON()
			return
		}
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}

	h.logger.Infof("ApproveBlockCity succeeded for actionID=%s by user=%s", req.ActionID, userID)
	data := map[string]any{
		"message": "Action processed successfully",
	}
	constant_utils.BaseResponseMaker(data, w, "Action processed successfully", http.StatusOK)
}
func (h *AccountBlockHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        PhoneNumber string `json:"phone_number"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.PhoneNumber) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "phone_number is required"}
        resp.SendJSON()
        return
    }

    userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
    fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
    phoneNumberCtx, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
    department, _ := r.Context().Value(constant.ContextKey("department")).(string)

    if strings.TrimSpace(userID) == "" ||
        strings.TrimSpace(fullName) == "" ||
        strings.TrimSpace(phoneNumberCtx) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
        resp.SendJSON()
        return
    }
    if strings.TrimSpace(department) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
        resp.SendJSON()
        return
    }

    maker := action.CPSAction{
        MakerID:          userID,
        MakerName:        fullName,
        MakerPhoneNumber: phoneNumberCtx,
        Department:       department,
    }

    user, err := h.service.GetUserByPhone(r.Context(), req.PhoneNumber, maker)
    if err != nil || strings.TrimSpace(user.UserCode) == "" {
        h.logger.Errorf("BlockUser failed to find user: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "User not found"}
        resp.SendJSON()
        return
    }

    actionCode, err := h.service.BlockUser(r.Context(), user.PhoneNumber, maker)
    if err != nil {
        h.logger.Errorf("BlockUser failed: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
        resp.SendJSON()
        return
    }
    data := map[string]any{
        "action_code": actionCode,
    }
    constant_utils.BaseResponseMaker(data, w, "User block action created", http.StatusCreated)
}
func (h *AccountBlockHandler) GetUserByPhone(w http.ResponseWriter, r *http.Request) {
    phoneNumber := r.URL.Query().Get("phone_number")
    if strings.TrimSpace(phoneNumber) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "phone_number is required"}
        resp.SendJSON()
        return
    }

    userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
    fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
    phoneNumberCtx, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
    department, _ := r.Context().Value(constant.ContextKey("department")).(string)

    if strings.TrimSpace(userID) == "" ||
        strings.TrimSpace(fullName) == "" ||
        strings.TrimSpace(phoneNumberCtx) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
        resp.SendJSON()
        return
    }
    if strings.TrimSpace(department) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
        resp.SendJSON()
        return
    }

    h.logger.Infof("GetUserByPhone requested by user: %s (%s)", userID, fullName)

    maker := action.CPSAction{
        MakerID:          userID,
        MakerName:        fullName,
        MakerPhoneNumber: phoneNumberCtx,
        Department:       department,
    }
    user, err := h.service.GetUserByPhone(r.Context(), phoneNumber, maker)
    if err != nil {
        h.logger.Errorf("GetUserByPhone failed: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "User not found"}
        resp.SendJSON()
        return
    }
    data, err := constant_utils.StructToMap(user)
    if err != nil {
        constant_utils.SendErrorResponse(w, http.StatusInternalServerError, 500, nil)
        return
    }
    constant_utils.BaseResponseMaker(data, w, "User retrieved successfully", http.StatusOK)
}
func (h *AccountBlockHandler) ApproveBlockUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ActionID string  `json:"action_id"`
        Approve  bool    `json:"approve"`
        Reason   *string `json:"reason,omitempty"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
        h.logger.Errorf("ApproveBlockUser: invalid request body or missing action_id, err=%v", err)
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
        h.logger.Errorf("ApproveBlockUser: Incomplete user information in context: userID=%s, fullName=%s, phoneNumber=%s", userID, fullName, phoneNumber)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
        resp.SendJSON()
        return
    }
    if strings.TrimSpace(department) == "" {
        h.logger.Errorf("ApproveBlockUser: department is required in context")
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "department is required in context"}
        resp.SendJSON()
        return
    }

    h.logger.Infof("ApproveBlockUser requested by user: %s (%s), actionID: %s, approve: %v", userID, fullName, req.ActionID, req.Approve)
    checker := action.User{
        UserID:      userID,
        FullName:    fullName,
        PhoneNumber: phoneNumber,
        Department:  department,
    }

    err := h.service.ApproveBlockUser(r.Context(), req.ActionID, req.Approve, req.Reason, checker)
    if err != nil {
        h.logger.Errorf("ApproveBlockUser failed for actionID=%s by user=%s: %v", req.ActionID, userID, err)
        if strings.Contains(err.Error(), "not allowed") {
            resp := common.Response[any]{ResponseWriter: w, Status: http.StatusForbidden, Data: err.Error()}
            resp.SendJSON()
            return
        }
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
        resp.SendJSON()
        return
    }

    h.logger.Infof("ApproveBlockUser succeeded for actionID=%s by user=%s", req.ActionID, userID)
    data := map[string]any{
        "message": "Action processed successfully",
    }
    constant_utils.BaseResponseMaker(data, w, "Action processed successfully", http.StatusOK)
}