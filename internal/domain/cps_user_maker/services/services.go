package services

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/cps_user_maker/repository"
)

type CPSUserService interface {
	CreateUserRequest(ctx context.Context, user action.CPSUser, maker action.User) error
	UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) error
	ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error
	GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error)
}

type cpsUserService struct {
	repo repository.CPSUserRepo
}

func NewCPSUserService(repo repository.CPSUserRepo) CPSUserService {
	return &cpsUserService{repo: repo}
}

func (s *cpsUserService) CreateUserRequest(ctx context.Context, user action.CPSUser, maker action.User) error {
	return s.repo.CreateUserRequest(ctx, user, maker)
}

func (s *cpsUserService) UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) error {
	return s.repo.UpdateUserRequest(ctx, updated, maker)
}

func (s *cpsUserService) ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error {
	return s.repo.ApproveUserAction(ctx, actionID, approve, reason)
}

func (s *cpsUserService) GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error) {
	return s.repo.GetPendingUserActions(ctx, actionCode)
}
func (s *cpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error) {
	return s.repo.FetchUserByUserCode(ctx, userCode)
}
