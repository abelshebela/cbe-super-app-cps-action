package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/repository"
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
	GetAllCPSUsers(ctx context.Context) ([]model.CPSUser, error)
}

type cpsUserService struct {
	repo repository.CPSUserRepo
}

func NewCPSUserService(repo repository.CPSUserRepo) CPSUserService {
	return &cpsUserService{repo: repo}
}

func (s *cpsUserService) CreateUserRequest(ctx context.Context, r *http.Request, userData userDTO.CreateUserRequest) (*model.CPSAction, error) {
	// Create the CPS action
	userPayload := ctx_util.ExtractContext(ctx)
	actionCode := utils.RandomGenerator(24)

	user := model.CPSUser{
		ID:                 bson.NewObjectID(),
		UserCode:           "CPS_USER_" + utils.RandomGenerator(15),
		FullName:           userData.FullName,
		Role:               userData.Role,
		Department:         userData.Department,
		Gender:             userData.Gender,
		PhoneNumber:        userData.PhoneNumber,
		Email:              userData.Email,
		UserName:           userData.UserName,
		Realm:              userData.Realm,
		Enabled:            userData.Enabled,
		Country:            userData.Country,
		Region:             userData.Region,
		PermissionCategory: userData.PermissionCategory,
		PermissionGroup:    userData.PermissionGroups,
	}

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
		CurrentAction:    user,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return s.repo.CreateUserRequest(ctx, cpsAction)
}

func (s *cpsUserService) UpdateUserRequest(ctx context.Context, r *http.Request, userData userDTO.UpdateUserRequest, userCode string) (*model.CPSAction, error) {
	userPayload := ctx_util.ExtractContext(ctx)
	actionCode := utils.RandomGenerator(24)
	userData.UserCode = userCode

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
	userPayload := ctx_util.ExtractContext(ctx)
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
		return nil, fmt.Errorf("NO_PENDING_ACTION_FOUND")
	}

	return data, nil
}

func (s *cpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*model.CPSUser, error) {
	user, err := s.repo.FetchUserByUserCode(ctx, userCode)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *cpsUserService) GetAllCPSUsers(ctx context.Context) ([]model.CPSUser, error) {
	users, err := s.repo.GetAllCPSUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
