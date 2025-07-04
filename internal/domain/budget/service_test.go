package budget_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type MockLogger struct{}

func (l *MockLogger) Debug(args ...interface{})                 {}
func (l *MockLogger) Debugf(format string, args ...interface{}) {}
func (l *MockLogger) Info(args ...interface{})                  {}
func (l *MockLogger) Infof(format string, args ...interface{})  {}
func (l *MockLogger) Warn(args ...interface{})                  {}
func (l *MockLogger) Warnf(format string, args ...interface{})  {}
func (l *MockLogger) Error(args ...interface{})                 {}
func (l *MockLogger) Errorf(format string, args ...interface{}) {}
func (l *MockLogger) Fatal(args ...interface{})                 {}
func (l *MockLogger) Fatalf(format string, args ...interface{}) {}
func (l *MockLogger) Sync() error                               { return nil }

func TestBudgetService_CreateIcon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := &MockLogger{}
	service := budget.InitBudgetDomain(repo, logger)

	action := entities.CPSAction{ActionCode: "test-action"}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().CreateIconAction(gomock.Any(), action).Return(&action, nil)
		result, err := service.CreateIcon(context.Background(), action)
		assert.NoError(t, err)
		assert.Equal(t, &action, result)
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().CreateIconAction(gomock.Any(), action).Return(nil, errors.New("db error"))
		result, err := service.CreateIcon(context.Background(), action)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBudgetService_FetchIcons(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := &MockLogger{}
	service := budget.InitBudgetDomain(repo, logger)

	t.Run("success", func(t *testing.T) {
		icons := []*entities.Icon{{Icon: "icon1"}, {Icon: "icon2"}}
		repo.EXPECT().FetchIcons(gomock.Any()).Return(icons, nil)
		result, err := service.FetchIcons(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, icons, result)
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().FetchIcons(gomock.Any()).Return(nil, errors.New("fetch error"))
		result, err := service.FetchIcons(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBudgetService_UpdateIcon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := &MockLogger{}
	service := budget.InitBudgetDomain(repo, logger)

	action := entities.CPSAction{ActionCode: "update-icon"}
	id := "icon123"

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().UpdateIcon(gomock.Any(), id, action).Return(&action, nil)
		result, err := service.UpdateIcon(context.Background(), id, action)
		assert.NoError(t, err)
		assert.Equal(t, &action, result)
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().UpdateIcon(gomock.Any(), id, action).Return(nil, errors.New("update failed"))
		result, err := service.UpdateIcon(context.Background(), id, action)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBudgetService_CreateColor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := &MockLogger{}
	service := budget.InitBudgetDomain(repo, logger)

	action := entities.CPSAction{ActionCode: "color-action"}
	hex := "#FFFFFF"

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().CreateColor(gomock.Any(), hex, action).Return(&action, nil)
		result, err := service.CreateColor(context.Background(), hex, action)
		assert.NoError(t, err)
		assert.Equal(t, &action, result)
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().CreateColor(gomock.Any(), hex, action).Return(nil, errors.New("create color error"))
		result, err := service.CreateColor(context.Background(), hex, action)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBudgetService_ApproveAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	logger := &MockLogger{}
	service := budget.InitBudgetDomain(repo, logger)

	action := entities.CPSAction{ActionCode: "approve"}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().ApproveAction(gomock.Any(), action).Return(&action, nil)
		result, err := service.ApproveAction(context.Background(), action)
		assert.NoError(t, err)
		assert.Equal(t, &action, result)
	})

	t.Run("error", func(t *testing.T) {
		repo.EXPECT().ApproveAction(gomock.Any(), action).Return(nil, errors.New("approval failed"))
		result, err := service.ApproveAction(context.Background(), action)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
