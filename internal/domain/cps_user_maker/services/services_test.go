package services

import (
    "context"
    "errors"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"

    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
    mockrepo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/cps_user_maker/repository/mocks"
)

func TestCPSUserService_CreateUserRequest(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mockrepo.NewMockCPSUserRepo(ctrl)
    service := NewCPSUserService(mockRepo)

    user := action.CPSUser{UserCode: "U1"}
    maker := action.User{UserID: "M1"} // Use UserID, not UserCode

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().
                    CreateUserRequest(ctx, user, maker).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().
                    CreateUserRequest(ctx, user, maker).
                    Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := service.CreateUserRequest(ctx, user, maker)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestCPSUserService_UpdateUserRequest(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mockrepo.NewMockCPSUserRepo(ctrl)
    service := NewCPSUserService(mockRepo)

    updated := action.CPSUser{UserCode: "U2"}
    maker := action.User{UserID: "M2"} // Use UserID, not UserCode

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().
                    UpdateUserRequest(ctx, updated, maker).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().
                    UpdateUserRequest(ctx, updated, maker).
                    Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := service.UpdateUserRequest(ctx, updated, maker)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestCPSUserService_ApproveUserAction(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mockrepo.NewMockCPSUserRepo(ctrl)
    service := NewCPSUserService(mockRepo)

    reason := "approved"
    actionID := "A1"

    tests := []struct {
        name      string
        actionID  string
        approve   bool
        reason    *string
        mockSetup func()
        wantErr   bool
    }{
        {
            name:     "success",
            actionID: actionID,
            approve:  true,
            reason:   &reason,
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveUserAction(ctx, actionID, true, &reason).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name:     "repo error",
            actionID: actionID,
            approve:  false,
            reason:   &reason,
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveUserAction(ctx, actionID, false, &reason).
                    Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := service.ApproveUserAction(ctx, tt.actionID, tt.approve, tt.reason)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestCPSUserService_GetPendingUserActions(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mockrepo.NewMockCPSUserRepo(ctrl)
    service := NewCPSUserService(mockRepo)

    actionCode := "AC1"
    expected := []action.CPSAction{
        {ID: "1", ActionCode: "AC1"},
        {ID: "2", ActionCode: "AC1"},
    }

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
        want      []action.CPSAction
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetPendingUserActions(ctx, actionCode).
                    Return(expected, nil)
            },
            wantErr: false,
            want:    expected,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetPendingUserActions(ctx, actionCode).
                    Return(nil, errors.New("db error"))
            },
            wantErr: true,
            want:    nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.GetPendingUserActions(ctx, actionCode)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}

func TestCPSUserService_FetchUserByUserCode(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mockrepo.NewMockCPSUserRepo(ctrl)
    service := NewCPSUserService(mockRepo)

    userCode := "U123"
    expected := &action.CPSUser{UserCode: userCode}

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
        want      *action.CPSUser
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FetchUserByUserCode(ctx, userCode).
                    Return(expected, nil)
            },
            wantErr: false,
            want:    expected,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FetchUserByUserCode(ctx, userCode).
                    Return(nil, errors.New("db error"))
            },
            wantErr: true,
            want:    nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.FetchUserByUserCode(ctx, userCode)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}