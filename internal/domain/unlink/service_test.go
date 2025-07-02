package unlink_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	unlinkDevice "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	mocks "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink/mocks"
)

func TestService_UnlinkDevice_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	service := unlinkDevice.NewUnlinkService(repo)

	userCode := "123"
	makerUser := "maker"
	branchCode := []string{"001"}
	homeBranch := "HB001"

	repo.
		EXPECT().
		UnlinkDevice(userCode, makerUser, branchCode, homeBranch).
		Return(nil)

	err := service.UnlinkDevice(userCode, makerUser, branchCode, homeBranch)
	assert.NoError(t, err)
}

func TestService_UnlinkDevice_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	service := unlinkDevice.NewUnlinkService(repo)

	userCode := "user1"
	makerUser := "maker1"
	branchCode := []string{"branch1"}
	homeBranch := "home1"
	expectedErr := assert.AnError

	repo.
		EXPECT().
		UnlinkDevice(userCode, makerUser, branchCode, homeBranch).
		Return(expectedErr)

	err := service.UnlinkDevice(userCode, makerUser, branchCode, homeBranch)
	assert.Equal(t, expectedErr, err)
}

func TestService_ApproveOrDecline_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	service := unlinkDevice.NewUnlinkService(repo)

	userCode := "123"
	decision := "APPROVE"
	reason := "valid reason"
	checkerUser := "checker"

	repo.
		EXPECT().
		ApproveOrDecline(userCode, decision, reason, checkerUser).
		Return(nil)

	err := service.ApproveOrDecline(userCode, decision, reason, checkerUser)
	assert.NoError(t, err)
}

func TestService_ApproveOrDecline_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	service := unlinkDevice.NewUnlinkService(repo)

	userCode := "user1"
	decision := "DECLINE"
	reason := "bad reason"
	checkerUser := "checker1"
	expectedErr := assert.AnError

	repo.
		EXPECT().
		ApproveOrDecline(userCode, decision, reason, checkerUser).
		Return(expectedErr)

	err := service.ApproveOrDecline(userCode, decision, reason, checkerUser)
	assert.Equal(t, expectedErr, err)
}
