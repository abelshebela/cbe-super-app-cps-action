package services

import (
	"context"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/repository"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSUserService interface {
	CreateUserRequest(ctx context.Context, r *http.Request, userData userDTO.CreateUserRequest) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, r *http.Request, userData userDTO.UpdateUserRequest, userCode string) (*model.CPSAction, error)
	ApproveUserAction(ctx context.Context, cpsAction model.CPSAction) error
	GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error)
}

type cpsUserService struct {
	repo repository.CPSUserRepo
}

func NewCPSUserService(repo repository.CPSUserRepo) CPSUserService {
	return &cpsUserService{repo: repo}
}

func (s *cpsUserService) CreateUserRequest(ctx context.Context, r *http.Request, userData userDTO.CreateUserRequest) (*model.CPSAction, error) {
	userPayload := ctx_util.ExtractUserContext(r)
	actionCode := utils.RandomGenerator(24)
	cpsAction := model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       actionCode,
		UniqueId:         userPayload.UserCode,
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionCreate),
		RequestAction:    string(model.RequestUser),
		CurrentAction:    userData,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return s.repo.CreateUserRequest(ctx, cpsAction)
}

func (s *cpsUserService) UpdateUserRequest(ctx context.Context, r *http.Request, userData userDTO.UpdateUserRequest, userCode string) (*model.CPSAction, error) {
	userPayload := ctx_util.ExtractUserContext(r)
	actionCode := utils.RandomGenerator(24)
	cpsAction := model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       actionCode,
		UniqueId:         userPayload.UserCode,
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionCreate),
		RequestAction:    string(model.RequestUser),
		CurrentAction:    userData,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return s.repo.UpdateUserRequest(ctx, cpsAction, userCode)
}

func (s *cpsUserService) ApproveUserAction(ctx context.Context, cpsAction model.CPSAction) error {
	return s.repo.ApproveUserAction(ctx, cpsAction)
}

func (s *cpsUserService) GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error) {
	return s.repo.GetPendingUserActions(ctx, actionCode)
}

func (s *cpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error) {
	return s.repo.FetchUserByUserCode(ctx, userCode)
}
