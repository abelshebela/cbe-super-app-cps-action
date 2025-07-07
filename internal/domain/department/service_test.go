package department_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    // "go.uber.org/mock/gomock"
	    "github.com/golang/mock/gomock"

    error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

    "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
    "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
    "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/mocks"
)

type dummyLogger struct{}

func (d dummyLogger) Infof(format string, args ...interface{})  {}
func (d dummyLogger) Errorf(format string, args ...interface{}) {}
func (d dummyLogger) Debugf(format string, args ...interface{}) {}
func (d dummyLogger) Warnf(format string, args ...interface{})  {}
func (d dummyLogger) Fatalf(format string, args ...interface{}) {}
func (d dummyLogger) Sync() error { return nil }

func TestCreateDepartment_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    ctx := context.Background()
    cpsRepo := mocks.NewMockCPSActionRepository(ctrl)
    deptRepo := mocks.NewMockDepartmentRepository(ctrl)
    logger := dummyLogger{}

    svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

    action := entities.CPSAction{}
    portalCards := []string{"CARD1", "CARD2"}
    deptName := "IT"

    // Expect CreateDepartment on the department repository
    deptRepo.EXPECT().
        CreateDepartment(ctx, gomock.Any()).
        Return(nil)

    err := svc.CreateDepartment(ctx, deptName, portalCards, action)
    assert.NoError(t, err)
}

func TestValidateActionRequest_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    ctx := context.Background()

    cpsRepo := mocks.NewMockCPSActionRepository(ctrl)
    deptRepo := mocks.NewMockDepartmentRepository(ctrl)
    logger := dummyLogger{}

    svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

    expectedAction := &entities.CPSAction{Department: "IT"}

    cpsRepo.EXPECT().FindByActionCode(ctx, "CPS_ABC123").Return(expectedAction, nil)

    action, err := svc.ValidateActionRequest(ctx, "CPS_ABC123", "IT")
    assert.NoError(t, err)
    assert.Equal(t, expectedAction, action)
}

func TestValidateActionRequest_NotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    ctx := context.Background()

    cpsRepo := mocks.NewMockCPSActionRepository(ctrl)
    deptRepo := mocks.NewMockDepartmentRepository(ctrl)
    logger := dummyLogger{}

    svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

    cpsRepo.EXPECT().FindByActionCode(ctx, "INVALID").Return(nil, nil)

    action, err := svc.ValidateActionRequest(ctx, "INVALID", "IT")
    assert.Nil(t, action)
    assert.EqualError(t, err, error_codes.ActionNotFound)
}

func TestValidateActionRequest_Unauthorized(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    ctx := context.Background()

    cpsRepo := mocks.NewMockCPSActionRepository(ctrl)
    deptRepo := mocks.NewMockDepartmentRepository(ctrl)
    logger := dummyLogger{}

    svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

    action := &entities.CPSAction{Department: "Finance"}
    cpsRepo.EXPECT().FindByActionCode(ctx, "CPS_XYZ").Return(action, nil)

    _, err := svc.ValidateActionRequest(ctx, "CPS_XYZ", "IT")
    assert.EqualError(t, err, error_codes.ActionNotAllowed)
}

func TestApproveActionRequest_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    ctx := context.Background()

    cpsRepo := mocks.NewMockCPSActionRepository(ctrl)
    deptRepo := mocks.NewMockDepartmentRepository(ctrl)
    logger := dummyLogger{}

    svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

    action := entities.CPSAction{}
    cpsRepo.EXPECT().ApproveActionRequest(ctx, "CPS_001", action).Return(nil)

    err := svc.ApproveActionRequest(ctx, "CPS_001", action)
    assert.NoError(t, err)
}