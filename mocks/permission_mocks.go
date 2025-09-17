package mocks

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/stretchr/testify/mock"
)

// MockPermissionDomainService is a mock implementation of the PermissionDomainService interface.
type MockPermissionDomainService struct {
	mock.Mock
}

func (m *MockPermissionDomainService) CreatePermissionGroup(oldGroupName, groupName, role string, permissionCategoryLists []string, cpsAction model.CPSAction) (model.CPSAction, error) {
	args := m.Called(oldGroupName, groupName, role, permissionCategoryLists, cpsAction)
	if args.Get(0) == nil {
		return model.CPSAction{}, args.Error(1)
	}
	return args.Get(0).(model.CPSAction), args.Error(1)
}

func (m *MockPermissionDomainService) ApprovePermissionGroup(actionCode string, action model.CPSAction) (model.CPSAction, error) {
	args := m.Called(actionCode, action)
	if args.Get(0) == nil {
		return model.CPSAction{}, args.Error(1)
	}
	return args.Get(0).(model.CPSAction), args.Error(1)
}

func (m *MockPermissionDomainService) RejectPermissionGroup(actionCode string, action model.CPSAction, reason string) (model.CPSAction, error) {
	args := m.Called(actionCode, action, reason)
	if args.Get(0) == nil {
		return model.CPSAction{}, args.Error(1)
	}
	return args.Get(0).(model.CPSAction), args.Error(1)
}

func (m *MockPermissionDomainService) UpdatePermissionGroup(groupName string, permissionCategoryLists []string) (entities.PermissionGroup, error) {
	args := m.Called(groupName, permissionCategoryLists)
	if args.Get(0) == nil {
		return entities.PermissionGroup{}, args.Error(1)
	}
	return args.Get(0).(entities.PermissionGroup), args.Error(1)
}

func (m *MockPermissionDomainService) GetPermissionGroup(groupName string) (entities.PermissionGroup, error) {
	args := m.Called(groupName)
	if args.Get(0) == nil {
		return entities.PermissionGroup{}, args.Error(1)
	}
	return args.Get(0).(entities.PermissionGroup), args.Error(1)
}

func (m *MockPermissionDomainService) ValidatePermissionCategories(ctx context.Context, categoryIDs []string) ([]string, error) {
	args := m.Called(categoryIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPermissionDomainService) GetPermissionCategories(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionCategory], error) {
	args := m.Called(ctx, filterParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common_util.PaginatedResponse[[]*entities.PermissionCategory]), args.Error(1)
}

func (m *MockPermissionDomainService) GetPermissionGroups(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionGroup], error) {
	args := m.Called(ctx, filterParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common_util.PaginatedResponse[[]*entities.PermissionGroup]), args.Error(1)
}
