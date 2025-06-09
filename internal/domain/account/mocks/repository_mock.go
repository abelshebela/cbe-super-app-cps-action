// ./internal/domain/account/mocks/repository_mock.go
package mocks

import (
	account "cbe-super-app-member-users/internal/domain/account"
	context "context"
	reflect "reflect"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateLinkedAccount mocks base method.
func (m *MockRepository) CreateLinkedAccount(ctx context.Context, account *account.LinkedAccount) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLinkedAccount", ctx, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLinkedAccount indicates an expected call of CreateLinkedAccount.
func (mr *MockRepositoryMockRecorder) CreateLinkedAccount(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLinkedAccount", reflect.TypeOf((*MockRepository)(nil).CreateLinkedAccount), ctx, account)
}

// FindAccountUserByID mocks base method.
func (m *MockRepository) FindAccountUserByID(ctx context.Context, id string) (*account.AccountUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAccountUserByID", ctx, id)
	ret0, _ := ret[0].(*account.AccountUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAccountUserByID indicates an expected call of FindAccountUserByID.
func (mr *MockRepositoryMockRecorder) FindAccountUserByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAccountUserByID", reflect.TypeOf((*MockRepository)(nil).FindAccountUserByID), ctx, id)
}

// UpdateUserCustomerNumber mocks base method.
func (m *MockRepository) UpdateUserCustomerNumber(ctx context.Context, userID, customerNumber string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserCustomerNumber", ctx, userID, customerNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserCustomerNumber indicates an expected call of UpdateUserCustomerNumber.
func (mr *MockRepositoryMockRecorder) UpdateUserCustomerNumber(ctx, userID, customerNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserCustomerNumber", reflect.TypeOf((*MockRepository)(nil).UpdateUserCustomerNumber), ctx, userID, customerNumber)
}