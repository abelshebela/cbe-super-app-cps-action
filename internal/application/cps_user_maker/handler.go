package cpsusermaker

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/services"
)

type ApplicationService interface {
	CreateUserRequest(ctx context.Context, user action.CPSUser, maker action.User) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) (*model.CPSAction, error)
	ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error
	GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error)

	FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error)
}

type Handler struct {
	service services.CPSUserService
}

func NewApplicationHandler(service services.CPSUserService) ApplicationService {
	return &Handler{
		service: service,
	}
}
func (h *Handler) CreateUserRequest(ctx context.Context, user action.CPSUser, maker action.User) (*model.CPSAction, error) {
	return h.service.CreateUserRequest(ctx, user, maker)
}

func (h *Handler) UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) (*model.CPSAction, error) {
	return h.service.UpdateUserRequest(ctx, updated, maker)
}

func (h *Handler) ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error {
	return h.service.ApproveUserAction(ctx, actionID, approve, reason)
}

func (h *Handler) GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error) {
	return h.service.GetPendingUserActions(ctx, actionCode)
}
func (h *Handler) FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error) {
	return h.service.FetchUserByUserCode(ctx, userCode)
}
