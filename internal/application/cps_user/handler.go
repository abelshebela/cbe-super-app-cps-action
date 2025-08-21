package cpsusermaker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/services"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	CreateUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, r *http.Request) (*userDTO.CPSUserDTO, error)
	GetAllCPSUsers(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*userDTO.CPSUserDTO], error)
	DeleteUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error)
	DisableUser(ctx context.Context, r *http.Request) error
	EnableUser(ctx context.Context, r *http.Request) error
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
		if err.Error() == "the provided hex string is not a valid ObjectID" {
			return nil, fmt.Errorf("INVALID_OBJECT_ID_FORMAT")
		}
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
		return nil, fmt.Errorf("USER_CODE_IS_REQUIRED")
	}

	var req userDTO.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind user data: %v", err)
		if err.Error() == "the provided hex string is not a valid ObjectID" {
			return nil, fmt.Errorf("INVALID_OBJECT_ID_FORMAT")
		}
		return nil, err
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		return nil, err
	}

	return h.service.UpdateUserRequest(ctx, r, req, userCode)
}

func (h *Handler) FetchUserByUserCode(ctx context.Context, r *http.Request) (*userDTO.CPSUserDTO, error) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	return h.service.FetchUserByUserCode(ctx, userCode)
}

func (h *Handler) GetAllCPSUsers(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*userDTO.CPSUserDTO], error) {
	users, err := h.service.GetAllCPSUsers(ctx, filterParams)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (h *Handler) DeleteUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		h.logger.Errorf("user_code is required")
		return nil, fmt.Errorf("USER_CODE_IS_REQUIRED")
	}

	return h.service.DeleteUserRequest(ctx, userCode)
}

func (h *Handler) DisableUser(ctx context.Context, r *http.Request) error {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))

	if userCode == "" {
		h.logger.Errorf("user_code is required")
		return fmt.Errorf("USER_CODE_IS_REQUIRED")
	}

	err := h.service.DisableUser(ctx, userCode)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) EnableUser(ctx context.Context, r *http.Request) error {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))

	if userCode == "" {
		h.logger.Errorf("user_code is required")
		return fmt.Errorf("USER_CODE_IS_REQUIRED")
	}

	err := h.service.EnableUser(ctx, userCode)
	if err != nil {
		return err
	}

	return nil
}
