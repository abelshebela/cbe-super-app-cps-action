package cpsusermaker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/services"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	CreateUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error)
	ApproveUserAction(ctx context.Context, r *http.Request) (*model.CPSAction, error)
	GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, r *http.Request) (*model.CPSUser, error)
}

type Handler struct {
	service services.CPSUserService
	logger  utils.Logger
}

func NewApplicationHandler(service services.CPSUserService, logger utils.Logger) ApplicationService {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error) {
	var req userDTO.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind user data: %v", err)
		return nil, err
	}

	req.Normalize()
	if err := req.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		return nil, err
	}

	return h.service.CreateUserRequest(ctx, r, req)
}

func (h *Handler) UpdateUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))

	if userCode == "" {
		h.logger.Errorf("user_code is required")
		return nil, fmt.Errorf("user_code is required")
	}

	var req userDTO.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind user data: %v", err)
		return nil, err
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		return nil, err
	}

	return h.service.UpdateUserRequest(ctx, r, req, userCode)
}

func (h *Handler) ApproveUserAction(ctx context.Context, r *http.Request) (*model.CPSAction, error) {
	actionID := strings.TrimSpace(chi.URLParam(r, "action_id"))

	if actionID == "" {
		h.logger.Errorf("action_id is required")
		return nil, fmt.Errorf("action_id is required")
	}

	var req userDTO.ApproveCPSAction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind user data: %v", err)
		return nil, err
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		return nil, err
	}

	return h.service.ApproveUserAction(ctx, r, req, actionID)
}

func (h *Handler) GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error) {
	return h.service.GetPendingUserActions(ctx)
}
func (h *Handler) FetchUserByUserCode(ctx context.Context, r *http.Request) (*model.CPSUser, error) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	return h.service.FetchUserByUserCode(ctx, userCode)
}
