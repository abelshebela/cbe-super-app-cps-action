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
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSUserService interface {
	CreateUserRequest(ctx context.Context, r *http.Request, userData userDTO.CreateUserRequest) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, r *http.Request, userData userDTO.UpdateUserRequest, userCode string) (*model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*userDTO.CPSUserDTO, error)
	GetAllCPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*userDTO.CPSUserDTO], error)
	Authorize(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	DeleteUserRequest(ctx context.Context, userCode string) (*model.CPSAction, error)
	EnableUser(ctx context.Context, userCode string) error
	DisableUser(ctx context.Context, userCode string) error
}

type cpsUserService struct {
	repo              repository.CPSUserRepo
	permissionService permission.PermissionDomainService
	logger            utils.Logger
}

func NewCPSUserService(repo repository.CPSUserRepo, permissionService permission.PermissionDomainService, logger utils.Logger) CPSUserService {
	return &cpsUserService{repo: repo, permissionService: permissionService, logger: logger}
}

func (s *cpsUserService) CreateUserRequest(ctx context.Context, r *http.Request, userData userDTO.CreateUserRequest) (*model.CPSAction, error) {
	// Validate PermissionCategory
	categoryIDs := make([]string, len(userData.PermissionCategory))
	for i, id := range userData.PermissionCategory {
		categoryIDs[i] = id.Hex()
	}

	_, err := s.permissionService.ValidatePermissionCategories(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	// Validate PermissionGroups
	groupIDs := make([]string, len(userData.PermissionGroups))
	for i, id := range userData.PermissionGroups {
		groupIDs[i] = id.Hex()
	}
	_, err = s.permissionService.ValidatePermissionGroups(ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	// Validate the phone number
	phone_number := common_util.FormatPhoneNumber(userData.PhoneNumber)
	if phone_number == "" {
		return nil, fmt.Errorf("UNSUPPORTED_PHONE_NUMBER_FORMAT")
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
		PhoneNumber:        phone_number,
		Email:              userData.Email,
		UserName:           userData.UserName,
		PermissionCategory: userData.PermissionCategory,
		PermissionGroup:    userData.PermissionGroups,
	}

	userPayload := ctx_util.ExtractContext(ctx)
	cpsAction := model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       actionCode,
		UniqueId:         "",
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
		_, err := s.permissionService.ValidatePermissionCategories(ctx, categoryIDs)
		if err != nil {
			return nil, err
		}
	}

	if len(userData.PermissionGroups) > 0 {
		groupIDs := make([]string, len(userData.PermissionGroups))
		for i, id := range userData.PermissionGroups {
			groupIDs[i] = id.Hex()
		}
		_, err := s.permissionService.ValidatePermissionGroups(ctx, groupIDs)
		if err != nil {
			return nil, err
		}
	}

	userPayload := ctx_util.ExtractContext(ctx)

	if userData.Role != "" && userData.Role != "maker" && userData.Role != "checker" {
		s.logger.Errorf("User role can only be either 'maker' or 'checker'")
		return nil, fmt.Errorf("MAKER_OR_CHECKER")
	}

	// Validate the phone number
	var phone_number string
	if userData.PhoneNumber != "" {
		phone_number = common_util.FormatPhoneNumber(userData.PhoneNumber)
		if phone_number == "" {
			return nil, fmt.Errorf("UNSUPPORTED_PHONE_NUMBER_FORMAT")
		}
	}

	userPayload = ctx_util.ExtractContext(ctx)
	actionCode := utils.RandomGenerator(24)
	userData.UserCode = userCode

	// Format phone_number
	userData.PhoneNumber = phone_number

	cpsAction := model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       actionCode,
		UniqueId:         "",
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

func (s *cpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*userDTO.CPSUserDTO, error) {
	user, err := s.repo.FetchUserByUserCode(ctx, userCode)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *cpsUserService) GetAllCPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*userDTO.CPSUserDTO], error) {
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
		return s.repo.AuthorizeUserCreate(ctx, action)
	case cps_constant.ActionUpdate:
		return s.repo.AuthorizeUserUpdate(ctx, action)
	case cps_constant.ActionDelete:
		return s.repo.AuthorizeUserDelete(ctx, action)
	case cps_constant.ActionEnable:
		return s.repo.AuthorizeUserEnable(ctx, action)
	case cps_constant.ActionDisable:
		return s.repo.AuthorizeUserDisable(ctx, action)
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
		UniqueId:         userPayload.UserCode + utils.RandomGenerator(6),
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

func (s *cpsUserService) EnableUser(ctx context.Context, userCode string) error {
	userPayload := ctx_util.ExtractContext(ctx)
	data, err := s.repo.FetchPendingActionsByUniqueID(ctx, userPayload.UserCode)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}
	if data != nil {
		return common.DefineError.General["PENDING_REQUEST_EXISTS"]
	}
	actionCode := utils.RandomGenerator(24)

	action := model.CPSUser{
		UserCode: userCode,
		Enabled:  true,
	}

	cpsAction := model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       actionCode,
		UniqueId:         userPayload.UserCode + utils.RandomGenerator(6),
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionEnable),
		RequestAction:    string(model.RequestCpsUserEnable),
		PreviousAction:   nil,
		CurrentAction:    action,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return s.repo.EnableDisableUser(ctx, userCode, cpsAction, model.RequestCpsUserEnable)
}

func (s *cpsUserService) DisableUser(ctx context.Context, userCode string) error {
	userPayload := ctx_util.ExtractContext(ctx)
	actionCode := utils.RandomGenerator(24)
	data, err := s.repo.FetchPendingActionsByUniqueID(ctx, userPayload.UserCode)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}
	if data != nil {
		return common.DefineError.General["PENDING_REQUEST_EXISTS"]
	}

	action := model.CPSUser{
		UserCode: userCode,
		Enabled:  false,
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
		ActionType:       string(model.ActionDisable),
		RequestAction:    string(model.RequestCpsUserDisable),
		PreviousAction:   nil,
		CurrentAction:    action,
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}

	return s.repo.EnableDisableUser(ctx, userCode, cpsAction, model.RequestCpsUserDisable)
}
