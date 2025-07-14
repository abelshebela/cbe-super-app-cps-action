package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/repository"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSUserService interface {
	CreateUserRequest(ctx context.Context, r *http.Request, userData userDTO.CreateUserRequest) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, r *http.Request, userData userDTO.UpdateUserRequest, userCode string) (*model.CPSAction, error)
	ApproveUserAction(ctx context.Context, r *http.Request, approved userDTO.ApproveCPSAction, actionID string) (*model.CPSAction, error)
	GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*model.CPSUser, error)
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
	userData.UserCode = "CPS_USER_" + utils.RandomGenerator(15)
	cpsAction := model.CPSAction{
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

	fmt.Println("check this one ------------------")
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
		ActionType:       string(model.ActionUpdate),
		RequestAction:    string(model.RequestUpdateUser),
		CurrentAction:    userData,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return s.repo.UpdateUserRequest(ctx, cpsAction, userCode)
}

func (s *cpsUserService) ApproveUserAction(ctx context.Context, r *http.Request, approved userDTO.ApproveCPSAction, actionID string) (*model.CPSAction, error) {
	userPayload := ctx_util.ExtractUserContext(r)
	cpsAction := model.CPSAction{
		CheckerID:          userPayload.UserID,
		CheckerName:        userPayload.FullName,
		CheckerPhoneNumber: userPayload.PhoneNumber,
		Department:         userPayload.Department,
	}

	if approved.Approved {
		cpsAction.ActionStatus = string(model.ActionApproved)
	} else if !approved.Approved {
		cpsAction.ActionStatus = string(model.ActionRejected)
		cpsAction.RejectionReason = *approved.Reason
	}
	return s.repo.ApproveUserAction(ctx, actionID, cpsAction)
}

func (s *cpsUserService) GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error) {
	data, err := s.repo.GetPendingUserActions(ctx)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("no pending action found")
	}

	return data, nil
}

func (s *cpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*model.CPSUser, error) {
	user, err := s.repo.FetchUserByUserCode(ctx, userCode)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
