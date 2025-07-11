package cpsusermaker

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/services"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	CreateUserRequest(ctx context.Context, r *http.Request) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) (*model.CPSAction, error)
	ApproveUserAction(ctx context.Context, cpsAction model.CPSAction) error
	GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error)
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

	if err := req.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		return nil, err
	}

	return h.service.CreateUserRequest(ctx, r, req)
}

func (h *Handler) UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) (*model.CPSAction, error) {
	return h.service.UpdateUserRequest(ctx, updated, maker)
}

func (h *Handler) ApproveUserAction(ctx context.Context, cpsAction model.CPSAction) error {
	return h.service.ApproveUserAction(ctx, cpsAction)
}

func (h *Handler) GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error) {
	return h.service.GetPendingUserActions(ctx, actionCode)
}
func (h *Handler) FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error) {
	return h.service.FetchUserByUserCode(ctx, userCode)
}
