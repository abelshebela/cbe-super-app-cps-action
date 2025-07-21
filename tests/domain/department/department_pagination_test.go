package department_test

import (
	"context"
	"testing"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDepartmentRepository is a mock implementation of DepartmentRepository
type MockDepartmentRepository struct {
	mock.Mock
}

func (m *MockDepartmentRepository) CheckDepartmentExists(ctx context.Context, department string) (bool, error) {
	args := m.Called(ctx, department)
	return args.Bool(0), args.Error(1)
}

func (m *MockDepartmentRepository) CreateDepartment(ctx context.Context, dept entities.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepository) UpdateDepartment(ctx context.Context, code string, department string, portalCards []string) (*entities.Department, error) {
	args := m.Called(ctx, code, department, portalCards)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Department), args.Error(1)
}

func (m *MockDepartmentRepository) GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*utils.PaginatedResponse[[]*entities.Department], error) {
	args := m.Called(ctx, filterParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.PaginatedResponse[[]*entities.Department]), args.Error(1)
}

// MockCPSActionRepository is a mock implementation of CPSActionRepository
type MockCPSActionRepository struct {
	mock.Mock
}

func (m *MockCPSActionRepository) CheckRequestExists(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error) {
	args := m.Called(ctx, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepository) CreateCPSAction(ctx context.Context, department string, portalCards []string, action entities.CPSAction) error {
	args := m.Called(ctx, department, portalCards, action)
	return args.Error(0)
}

func (m *MockCPSActionRepository) ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*entities.CPSAction, error) {
	args := m.Called(ctx, actionCode, userDept)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepository) FindByActionCode(ctx context.Context, code string) (*entities.CPSAction, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepository) UpdateActionStatus(ctx context.Context, actionCode string, status string) (*entities.CPSAction, error) {
	args := m.Called(ctx, actionCode, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepository) ApproveActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error {
	args := m.Called(ctx, actionCode, action)
	return args.Error(0)
}

func (m *MockCPSActionRepository) RejectActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error {
	args := m.Called(ctx, actionCode, action)
	return args.Error(0)
}

func (m *MockCPSActionRepository) CreateDepartmentUpdateCPSAction(ctx context.Context, req entities.CPSAction) (*entities.CPSAction, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CPSAction), args.Error(1)
}

func (m *MockCPSActionRepository) ApproveDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error {
	args := m.Called(ctx, cpsAction)
	return args.Error(0)
}

func (m *MockCPSActionRepository) RejectDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error {
	args := m.Called(ctx, cpsAction)
	return args.Error(0)
}

func TestDepartmentService_GetAllDepartments_WithPagination(t *testing.T) {
	// Setup
	mockRepo := &MockDepartmentRepository{}
	mockCPSRepo := &MockCPSActionRepository{}

	service := &department.Service{
		DepartmentRepository: mockRepo,
		CPSActionRepository:  mockCPSRepo,
	}

	ctx := context.Background()
	filterParams := &constant.Filter{
		Page:    1,
		PerPage: 10,
		Search:  "IT",
	}

	// Mock data
	mockDepartments := []*entities.Department{
		{
			DepartmentCode: "IT001",
			Department:     "Information Technology",
			PortalCards:    []string{"card1", "card2"},
			Enabled:        true,
		},
		{
			DepartmentCode: "IT002",
			Department:     "IT Support",
			PortalCards:    []string{"card3"},
			Enabled:        true,
		},
	}

	expectedResponse := &utils.PaginatedResponse[[]*entities.Department]{
		Data: mockDepartments,
		Meta: utils.PaginationMeta{
			TotalDocs:  2,
			Limit:      10,
			TotalPages: 1,
			Page:       1,
		},
	}

	// Setup mock expectations
	mockRepo.On("GetAllDepartments", ctx, filterParams).Return(expectedResponse, nil)

	// Execute
	result, err := service.GetAllDepartments(ctx, filterParams)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse.Data, result.Data)
	assert.Equal(t, expectedResponse.Meta.TotalDocs, result.Meta.TotalDocs)
	assert.Equal(t, expectedResponse.Meta.Limit, result.Meta.Limit)
	assert.Equal(t, expectedResponse.Meta.TotalPages, result.Meta.TotalPages)
	assert.Equal(t, expectedResponse.Meta.Page, result.Meta.Page)

	// Verify mock was called correctly
	mockRepo.AssertExpectations(t)
}

func TestDepartmentService_GetAllDepartments_WithoutSearch(t *testing.T) {
	// Setup
	mockRepo := &MockDepartmentRepository{}
	mockCPSRepo := &MockCPSActionRepository{}

	service := &department.Service{
		DepartmentRepository: mockRepo,
		CPSActionRepository:  mockCPSRepo,
	}

	ctx := context.Background()
	filterParams := &constant.Filter{
		Page:    2,
		PerPage: 5,
		Search:  "", // No search
	}

	// Mock data
	mockDepartments := []*entities.Department{
		{
			DepartmentCode: "HR001",
			Department:     "Human Resources",
			PortalCards:    []string{"hr_card"},
			Enabled:        true,
		},
	}

	expectedResponse := &utils.PaginatedResponse[[]*entities.Department]{
		Data: mockDepartments,
		Meta: utils.PaginationMeta{
			TotalDocs:  6,
			Limit:      5,
			TotalPages: 2,
			Page:       2,
		},
	}

	// Setup mock expectations
	mockRepo.On("GetAllDepartments", ctx, filterParams).Return(expectedResponse, nil)

	// Execute
	result, err := service.GetAllDepartments(ctx, filterParams)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse.Data, result.Data)
	assert.Equal(t, expectedResponse.Meta.TotalDocs, result.Meta.TotalDocs)
	assert.Equal(t, expectedResponse.Meta.Limit, result.Meta.Limit)
	assert.Equal(t, expectedResponse.Meta.TotalPages, result.Meta.TotalPages)
	assert.Equal(t, expectedResponse.Meta.Page, result.Meta.Page)

	// Verify mock was called correctly
	mockRepo.AssertExpectations(t)
}
