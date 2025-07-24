package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cps_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/repository"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSUserService interface {
	CreateUserRequest(ctx context.Context, r *http.Request, userData userDTO.CreateUserRequest) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, r *http.Request, userData userDTO.UpdateUserRequest, userCode string) (*model.CPSAction, error)
	ApproveUserAction(ctx context.Context, r *http.Request, approved userDTO.ApproveCPSAction, actionID string) (*model.CPSAction, error)
	GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*model.CPSUser, error)
	GetAllCPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*model.CPSUser], error)
	Authorize(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	DeleteUserRequest(ctx context.Context, userCode string) (*model.CPSAction, error)
}

type cpsUserService struct {
	repo              repository.CPSUserRepo
	repoDepartment    department.DepartmentRepository
	permissionService permission.PermissionDomainService
	logger            utils.Logger
}

func NewCPSUserService(repo repository.CPSUserRepo, permissionService permission.PermissionDomainService, repoDepartment department.DepartmentRepository, logger utils.Logger) CPSUserService {
	return &cpsUserService{repo: repo, permissionService: permissionService, repoDepartment: repoDepartment, logger: logger}
}

func (s *cpsUserService) CreateUserRequest(ctx context.Context, r *http.Request, userData userDTO.CreateUserRequest) (*model.CPSAction, error) {
	// Validate PermissionCategory
	categoryIDs := make([]string, len(userData.PermissionCategory))
	for i, id := range userData.PermissionCategory {
		categoryIDs[i] = id.Hex()
	}

	_, err := s.permissionService.ValidatePermissionCategories(categoryIDs)
	if err != nil {
		return nil, err
	}

	// Validate PermissionGroups
	groupIDs := make([]string, len(userData.PermissionGroups))
	for i, id := range userData.PermissionGroups {
		groupIDs[i] = id.Hex()
	}
	_, err = s.permissionService.ValidatePermissionGroups(groupIDs)
	if err != nil {
		return nil, err
	}

	// Check if the department exists
	dept, err := s.repo.GetDepartmentByID(ctx, userData.Department.Hex())
	if err != nil {
		s.logger.Errorf("failed to check department existence: %v", err)
		return nil, err
	}
	if dept == nil {
		s.logger.Errorf("DEPARTMENT_NOT_FOUND")
		return nil, err
	}
	// Check for existing pending actions for this user
	userPayload := ctx_util.ExtractContext(ctx)
	pendingActions, err := s.repo.FetchPendingActionsByUniqueID(ctx, userPayload.UserCode)
	if err != nil && err.Error() != "NOT_FOUND" {
		s.logger.Errorf("failed to fetch pending actions for user %s: %v", userPayload.UserCode, err)
		return nil, err
	}
	if len(pendingActions) > 0 {
		s.logger.Errorf("pending action already exists for user: %s", userPayload.UserCode)
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	if userData.Role != "maker" && userData.Role != "Checker" {
		s.logger.Errorf("User role can only be either 'Maker' or 'Checker'")
		return nil, fmt.Errorf("MAKER_OR_CHECKER")
	}

	// Create the CPS action
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
		RequestAction:    string(model.RequestCpsUserCreate),
		PreviousAction:   nil,
		CurrentAction:    user,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return s.repo.CreateUserRequest(ctx, cpsAction)
}

func (s *cpsUserService) UpdateUserRequest(ctx context.Context, r *http.Request, userData userDTO.UpdateUserRequest, userCode string) (*model.CPSAction, error) {

	if len(userData.PermissionCategory) > 0 {
		categoryIDs := make([]string, len(userData.PermissionCategory))
		for i, id := range userData.PermissionCategory {
			categoryIDs[i] = id.Hex()
		}
		_, err := s.permissionService.ValidatePermissionCategories(categoryIDs)
		if err != nil {
			return nil, err
		}
	}

	if len(userData.PermissionGroups) > 0 {
		groupIDs := make([]string, len(userData.PermissionGroups))
		for i, id := range userData.PermissionGroups {
			groupIDs[i] = id.Hex()
		}
		_, err := s.permissionService.ValidatePermissionGroups(groupIDs)
		if err != nil {
			return nil, err
		}
	}

	// Check for existing pending actions for this user
	userPayload := ctx_util.ExtractContext(ctx)
	pendingActions, err := s.repo.FetchPendingActionsByUniqueID(ctx, userPayload.UserCode)
	if err != nil && err.Error() != "NOT_FOUND" {
		s.logger.Errorf("failed to fetch pending actions for user %s: %v", userPayload.UserCode, err)
		return nil, err
	}
	if len(pendingActions) > 0 {
		s.logger.Errorf("pending action already exists for user: %s", userPayload.UserCode)
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	dept, err := s.repo.GetDepartmentByID(ctx, userData.Department.Hex())
	if err != nil {
		s.logger.Errorf("failed to check department existence: %v", err)
		return nil, err
	}
	if dept == nil {
		s.logger.Errorf("DEPARTMENT_NOT_FOUND")
		return nil, err
	}

	if userData.Role != "maker" && userData.Role != "Checker" {
		s.logger.Errorf("User role can only be either 'Maker' or 'Checker'")
		return nil, fmt.Errorf("MAKER_OR_CHECKER")
	}

	userPayload = ctx_util.ExtractContext(ctx)
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
		RequestAction:    string(model.RequestCpsUserUpdate),
		PreviousAction:   nil,
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

func (s *cpsUserService) GetAllCPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*model.CPSUser], error) {
	users, err := s.repo.GetAllCPSUsers(ctx, filterParams)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *cpsUserService) Authorize(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime

	switch action.ActionType {
	case cps_constant.ActionCreate:
		return s.repo.AuthorizeCreate(ctx, action)
	case cps_constant.ActionUpdate:
		return s.repo.AuthorizeUpdate(ctx, action)
	case cps_constant.ActionDelete:
		return s.repo.AuthorizeDelete(ctx, action)

	default:
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}
}

func (s *cpsUserService) DeleteUserRequest(ctx context.Context, userCode string) (*model.CPSAction, error) {
	userPayload := ctx_util.ExtractContext(ctx)
	actionCode := utils.RandomGenerator(24)

	action := model.CPSUser{
		UserCode:  userCode,
		IsDeleted: true,
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
		ActionType:       string(model.ActionDelete),
		RequestAction:    string(model.RequestCpsUserDelete),
		PreviousAction:   nil,
		CurrentAction:    action,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return s.repo.DeleteUserRequest(ctx, userCode, cpsAction)
}
