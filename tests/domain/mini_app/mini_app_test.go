package test_miniapp

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	mock_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/mocks/domain/miniapp"
)

func TestCreateMiniAppAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockRepository(ctrl)
	miniAppStore := miniapp.MiniAppStore{Repository: mockRepo}

	ctx := context.Background()
	makerUser := model.User{
		UserCode:    "1234",
		FullName:    "Abebe Kebede",
		PhoneNumber: "+251911111111",
	}
	miniApp := miniapp.MiniApp{
		AppName: "Test MiniApp",
	}

	actionId := "CPS_1234567890"
	expectedAction := action.CPSAction{
		ActionCode:       actionId,
		MakerID:          makerUser.UserCode,
		MakerName:        makerUser.FullName,
		MakerPhoneNumber: makerUser.PhoneNumber,
		ActionType:       action.ActionCreate,
		RequestAction:    action.RequestCreateMiniAppMerchant,
		ActionStatus:     action.ActionPending,
		PreviosAction:    miniApp,
	}

	mockRepo.EXPECT().
		CreateMiniAppAction(ctx, gomock.Any()).
		Return(expectedAction, nil)

	// Act
	result, err := miniAppStore.CreateMiniAppAction(ctx, miniApp, makerUser)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, actionId, result)
}

func TestCheckMiniApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockRepository(ctrl)
	miniAppStore := miniapp.MiniAppStore{Repository: mockRepo}

	ctx := context.Background()
	actionId := "CPS_1234567890"
	checkerUser := model.User{
		UserCode:    "1234",
		FullName:    "Abebe Kebede",
		PhoneNumber: "+251911111111",
	}
	miniApp := miniapp.MiniApp{
		AppName: "Test MiniApp",
	}

	action := action.CPSAction{
		ActionCode:    actionId,
		ActionType:    action.ActionType(action.ActionPending),
		ActionStatus:  action.ActionPending,
		PreviosAction: miniApp,
	}

	mockRepo.EXPECT().
		GetMiniAppActionId(ctx, actionId).
		Return(action, nil)

	mockRepo.EXPECT().
		UpdateCpsAction(ctx, gomock.Any()).
		Return(nil)

	mockRepo.EXPECT().
		CreateMiniApp(ctx, miniApp).
		Return(nil)

	// Act
	err := miniAppStore.CheckMiniApp(ctx, actionId, true, checkerUser)

	// Assert
	assert.NoError(t, err)
}

func TestCheckMiniApp_Rejected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockRepository(ctrl)
	miniAppStore := miniapp.MiniAppStore{Repository: mockRepo}

	ctx := context.Background()
	actionId := "CPS_1234567890"
	checkerUser := model.User{
		UserCode:    "1234",
		FullName:    "Abebe Kebede",
		PhoneNumber: "+251911111111",
	}
	miniApp := miniapp.MiniApp{
		AppName: "Test MiniApp",
	}

	action := action.CPSAction{
		ActionCode:    actionId,
		ActionType:    action.ActionType(action.ActionPending),
		ActionStatus:  action.ActionPending,
		PreviosAction: miniApp,
	}

	mockRepo.EXPECT().
		GetMiniAppActionId(ctx, actionId).
		Return(action, nil)

	mockRepo.EXPECT().
		UpdateCpsAction(ctx, gomock.Any()).
		Return(nil)

	mockRepo.EXPECT().
		CreateMiniApp(ctx, miniApp).
		Return(nil)

	// Act
	err := miniAppStore.CheckMiniApp(ctx, actionId, false, checkerUser)

	// Assert
	assert.NoError(t, err)
}
