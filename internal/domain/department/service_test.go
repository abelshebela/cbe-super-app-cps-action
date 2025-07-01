package department_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cbe-super-app-cps-action/internal/domain/department"
	"cbe-super-app-cps-action/internal/domain/department/entities"
	"cbe-super-app-cps-action/internal/domain/department/mock"
)

type dummyLogger struct{}

func (d dummyLogger) Infof(format string, args ...interface{})  {}
func (d dummyLogger) Errorf(format string, args ...interface{}) {}
func (d dummyLogger) Debugf(format string, args ...interface{}) {}
func (d dummyLogger) Warnf(format string, args ...interface{})  {}
func (d dummyLogger) Fatalf(format string, args ...interface{}) {}
func (d dummyLogger) Sync() error {
	return nil
}

func TestCreateDepartment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cpsRepo := mock.NewMockCPSActionRepository(ctrl)
	deptRepo := mock.NewMockDepartmentRepository(ctrl)
	logger := dummyLogger{}

	svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

	action := entities.CPSAction{}
	portalCards := []string{"CARD1", "CARD2"}
	deptName := "IT"

	cpsRepo.EXPECT().CheckRequestExists(action).Return(false, nil)
	deptRepo.EXPECT().CheckDepartmentExists(deptName).Return(false, nil)
	cpsRepo.EXPECT().CreateCPSAction(deptName, portalCards, gomock.Any()).Return(nil)

	err := svc.CreateDepartment(deptName, portalCards, action)
	assert.NoError(t, err)
}

func TestCreateDepartment_RequestExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cpsRepo := mock.NewMockCPSActionRepository(ctrl)
	deptRepo := mock.NewMockDepartmentRepository(ctrl)
	logger := dummyLogger{}

	svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

	action := entities.CPSAction{}

	cpsRepo.EXPECT().CheckRequestExists(action).Return(true, nil)

	err := svc.CreateDepartment("Finance", nil, action)
	assert.EqualError(t, err, "You have a pending request for this action")
}

func TestCreateDepartment_DepartmentExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cpsRepo := mock.NewMockCPSActionRepository(ctrl)
	deptRepo := mock.NewMockDepartmentRepository(ctrl)
	logger := dummyLogger{}

	svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

	action := entities.CPSAction{}

	cpsRepo.EXPECT().CheckRequestExists(action).Return(false, nil)
	deptRepo.EXPECT().CheckDepartmentExists("HR").Return(true, nil)

	err := svc.CreateDepartment("HR", nil, action)
	assert.EqualError(t, err, "Department already exists")
}

func TestValidateActionRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cpsRepo := mock.NewMockCPSActionRepository(ctrl)
	deptRepo := mock.NewMockDepartmentRepository(ctrl)
	logger := dummyLogger{}

	svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

	expectedAction := &entities.CPSAction{Department: "IT"}

	cpsRepo.EXPECT().FindByActionCode("CPS_ABC123").Return(expectedAction, nil)

	action, err := svc.ValidateActionRequest("CPS_ABC123", "IT")
	assert.NoError(t, err)
	assert.Equal(t, expectedAction, action)
}

func TestValidateActionRequest_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cpsRepo := mock.NewMockCPSActionRepository(ctrl)
	deptRepo := mock.NewMockDepartmentRepository(ctrl)
	logger := dummyLogger{}

	svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

	cpsRepo.EXPECT().FindByActionCode("INVALID").Return(nil, nil)

	action, err := svc.ValidateActionRequest("INVALID", "IT")
	assert.Nil(t, action)
	assert.EqualError(t, err, "action not found")
}

func TestValidateActionRequest_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cpsRepo := mock.NewMockCPSActionRepository(ctrl)
	deptRepo := mock.NewMockDepartmentRepository(ctrl)
	logger := dummyLogger{}

	svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

	action := &entities.CPSAction{Department: "Finance"}
	cpsRepo.EXPECT().FindByActionCode("CPS_XYZ").Return(action, nil)

	_, err := svc.ValidateActionRequest("CPS_XYZ", "IT")
	assert.EqualError(t, err, "You are not allowed to approve this request")
}

func TestApproveActionRequest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cpsRepo := mock.NewMockCPSActionRepository(ctrl)
	deptRepo := mock.NewMockDepartmentRepository(ctrl)
	logger := dummyLogger{}

	svc := department.InitDepartmentDomain(cpsRepo, deptRepo, logger)

	action := entities.CPSAction{}
	cpsRepo.EXPECT().ApproveActionRequest("CPS_001", action).Return(nil)

	err := svc.ApproveActionRequest("CPS_001", action)
	assert.NoError(t, err)
}
