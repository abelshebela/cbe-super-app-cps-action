package cpsusermaker

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bps_user"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*userDTO.BPSUser, error)
	GetAllBPSUsers(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*userDTO.BPSUser], error)
	DisableUser(ctx context.Context, userCode, makerID, phone, fullName, dept string) (*action.CPSAction, error)
	EnableUser(ctx context.Context, userCode, makerID, phone, fullName, dept string) (*action.CPSAction, error)
}

type Handler struct {
	service bps_user.Service
	logger  utils.Logger
}

func NewApplicationHandler(service bps_user.Service, logger utils.Logger) ApplicationService {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error) {
	// This method needs to be implemented based on the new service interface
	// For now, returning empty slice
	return []model.CPSAction{}, nil
}

func (h *Handler) FetchUserByUserCode(ctx context.Context, userCode string) (*userDTO.BPSUser, error) {
	user, err := h.service.GetBPSUserByUserCode(ctx, userCode)
	if err != nil {
		return nil, err
	}

	// Convert domain BPSUser to DTO BPSUser
	dtoUser := &userDTO.BPSUser{
		ID:                user.ID,
		UserCode:          user.UserCode,
		FullName:          user.FullName,
		Username:          user.Username,
		PhoneNumber:       user.PhoneNumber,
		BranchCode:        user.BranchCode,
		BranchName:        user.BranchName,
		HomeBranch:        user.HomeBranch,
		Role:              user.Role,
		Realm:             user.Realm,
	
		Enabled:           user.Enabled,
		IsDeleted:         user.IsDeleted,
	}

	return dtoUser, nil
}

func (h *Handler) GetAllBPSUsers(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*userDTO.BPSUser], error) {
	users, err := h.service.GetAllBPSUsers(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	// Convert domain BPSUsers to DTO BPSUsers
	dtoUsers := make([]*userDTO.BPSUser, 0, len(users.Data))
	for _, user := range users.Data {
		dtoUser := &userDTO.BPSUser{
			ID:                user.ID,
			UserCode:          user.UserCode,
			FullName:          user.FullName,
			Username:          user.Username,
			PhoneNumber:       user.PhoneNumber,
			BranchCode:        user.BranchCode,
			BranchName:        user.BranchName,
			HomeBranch:        user.HomeBranch,
			Role:              user.Role,
			Realm:             user.Realm,
			Enabled:           user.Enabled,
			IsDeleted:         user.IsDeleted,

		}
		dtoUsers = append(dtoUsers, dtoUser)
	}

	return &common_util.PaginatedResponse[[]*userDTO.BPSUser]{
		Data: dtoUsers,
		Meta: users.Meta,
	}, nil
}

func (h *Handler) DisableUser(ctx context.Context, userCode, makerID, phone, fullName, dept string) (*action.CPSAction, error) {
	if userCode == "" {
		h.logger.Errorf("user_code is required")
		return nil, fmt.Errorf("USER_CODE_IS_REQUIRED")
	}

	request := bps_user.DisableBPSUserRequest{
		UserCode:   userCode,
		MakerID:    makerID,
		MakerName:  fullName,
		MakerPhone: phone,
		Department: dept,
	}

	return h.service.DisableBPSUserRequest(ctx, request)
}

func (h *Handler) EnableUser(ctx context.Context, userCode, makerID, phone, fullName, dept string) (*action.CPSAction, error) {
	if userCode == "" {
		h.logger.Errorf("user_code is required")
		return nil, fmt.Errorf("USER_CODE_IS_REQUIRED")
	}

	request := bps_user.EnableBPSUserRequest{
		UserCode:   userCode,
		MakerID:    makerID,
		MakerName:  fullName,
		MakerPhone: phone,
		Department: dept,
	}

	return h.service.EnableBPSUserRequest(ctx, request)
}
