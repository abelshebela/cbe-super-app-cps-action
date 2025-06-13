package branch_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	branchapp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/branch"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bulkcustomer/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"
)

type BranchHandler struct {
	service branchapp.ApplicationService
	logger  utils.Logger
}

func NewBranchHandler(service branchapp.ApplicationService, logger utils.Logger) inbound.BranchHandler {
	return &BranchHandler{
		service: service,
		logger:  logger,
	}
}
func (h *BranchHandler) FilterSingleBranches(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Region   string `json:"region"`
		District string `json:"district"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Region) == "" || strings.TrimSpace(req.District) == "" {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "region and district are required"}
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
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: branches}
	resp.SendJSON()
}
func (h *BranchHandler) DisableSingleBranch(w http.ResponseWriter, r *http.Request) {
    var req struct {
        BranchCode string `json:"branch_code"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.BranchCode) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "branch_code is required"}
        resp.SendJSON()
        return
    }

    userID, _ := r.Context().Value(middleware.ContextKey("user_id")).(string)
    fullName, _ := r.Context().Value(middleware.ContextKey("fullname")).(string)
    phoneNumber, _ := r.Context().Value(middleware.ContextKey("phoneNumber")).(string)

    if strings.TrimSpace(userID) == "" ||
        strings.TrimSpace(fullName) == "" ||
        strings.TrimSpace(phoneNumber) == "" {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "Incomplete user information"}
        resp.SendJSON()
        return
    }

    branch, err := h.service.GetBranchByCode(r.Context(), req.BranchCode)
    if err != nil {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "Branch not found"}
        resp.SendJSON()
        return
    }

    maker := entities.User{
        UserCode:    userID,
        FullName:    fullName,
        PhoneNumber: phoneNumber,
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
    resp := common.Response[any]{ResponseWriter: w, Status: http.StatusCreated, Data: "Branch is disabled and CPS action created"}
    resp.SendJSON()
}
func (h *BranchHandler) ApproveSingleBranchDisable(w http.ResponseWriter, r *http.Request) {
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
	userPayload, ok := r.Context().Value(middleware.ContextKey("user_payload")).(middleware.UserPayload)
	if !ok {
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
		resp.SendJSON()
		return
	}
	h.logger.Infof("ApproveSingleBranchDisable requested by user: %s (%s)", userPayload.UserID, userPayload.FullName)
	err := h.service.ApproveSingleBranchDisable(r.Context(), req.ActionID, req.Approve, req.Reason)
	if err != nil {
		h.logger.Errorf("ApproveSingleBranchDisable failed: %v", err)
		resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
		resp.SendJSON()
		return
	}
	resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: "Action processed successfully"}
	resp.SendJSON()
}

func (h *BranchHandler) FilterMultipleBranches(w http.ResponseWriter, r *http.Request) {
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
    resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: branches}
    resp.SendJSON()
}
func (h *BranchHandler) DisableMultipleBranches(w http.ResponseWriter, r *http.Request) {
    var req struct {
        BranchCodes []string `json:"branch_codes"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.BranchCodes) == 0 {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: "branch_codes are required"}
        resp.SendJSON()
        return
    }
    userPayload, ok := r.Context().Value(middleware.ContextKey("user_payload")).(middleware.UserPayload)
    if !ok {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
        resp.SendJSON()
        return
    }

    var branches []entities.Branch
    for _, code := range req.BranchCodes {
        branch, err := h.service.GetBranchByCode(r.Context(), code)
        if err != nil {
            resp := common.Response[any]{ResponseWriter: w, Status: http.StatusNotFound, Data: "Branch not found: " + code}
            resp.SendJSON()
            return
        }
        branches = append(branches, branch)
    }

    maker := entities.User{
        UserCode:    userPayload.UserID,
        FullName:    userPayload.FullName,
        PhoneNumber: userPayload.PhoneNumber,
    }

       _, err := h.service.DisableMultipleBranches(r.Context(), branches, maker)
    if err != nil {
        h.logger.Errorf("DisableMultipleBranches failed: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusInternalServerError, Data: err.Error()}
        resp.SendJSON()
        return
    }
    resp := common.Response[any]{ResponseWriter: w, Status: http.StatusCreated, Data: "Branches are disabled and CPS action created"}
    resp.SendJSON()
}
func (h *BranchHandler) ApproveBulkBranchesDisable(w http.ResponseWriter, r *http.Request) {
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
    userPayload, ok := r.Context().Value(middleware.ContextKey("user_payload")).(middleware.UserPayload)
    if !ok {
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusUnauthorized, Data: "User info missing in context"}
        resp.SendJSON()
        return
    }
    h.logger.Infof("ApproveBulkBranchesDisable requested by user: %s (%s)", userPayload.UserID, userPayload.FullName)
    err := h.service.ApproveBulkBranchesDisable(r.Context(), req.ActionID, req.Approve, req.Reason)
    if err != nil {
        h.logger.Errorf("ApproveBulkBranchesDisable failed: %v", err)
        resp := common.Response[any]{ResponseWriter: w, Status: http.StatusBadRequest, Data: err.Error()}
        resp.SendJSON()
        return
    }
    resp := common.Response[any]{ResponseWriter: w, Status: http.StatusOK, Data: "Action processed successfully"}
    resp.SendJSON()
}