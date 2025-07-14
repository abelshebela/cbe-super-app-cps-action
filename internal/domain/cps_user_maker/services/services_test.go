package services

import (
    "context"
    "errors"
    "net/http"
    "testing"
    "time"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"

    "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
    userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/dto"
    mockrepo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/repository/mocks"
)

func TestCPSUserService_CreateUserRequest(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mockrepo.NewMockCPSUserRepo(ctrl)
    service := NewCPSUserService(mockRepo)

    req := &http.Request{}
    userData := userDTO.CreateUserRequest{ /* fill fields as needed */ }
    expected := &model.CPSAction{ActionCode: "A1", CreatedAt: time.Now()}

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().
                    CreateUserRequest(ctx, gomock.Any()).
                    Return(expected, nil)
            },
            wantErr: false,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().
                    CreateUserRequest(ctx, gomock.Any()).
                    Return(nil, errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.CreateUserRequest(ctx, req, userData)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, expected, got)
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

    req := &http.Request{}
    userData := userDTO.UpdateUserRequest{ /* fill fields as needed */ }
    userCode := "U1"
    expected := &model.CPSAction{ActionCode: "A2", CreatedAt: time.Now()}

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().
                    UpdateUserRequest(ctx, gomock.Any(), userCode).
                    Return(expected, nil)
            },
            wantErr: false,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().
                    UpdateUserRequest(ctx, gomock.Any(), userCode).
                    Return(nil, errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.UpdateUserRequest(ctx, req, userData, userCode)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, expected, got)
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

    req := &http.Request{}
    approved := userDTO.ApproveCPSAction{Approved: true, Reason: nil}
    actionID := "A1"
    expected := &model.CPSAction{ActionCode: actionID, ActionStatus: string(model.ActionApproved)}

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveUserAction(ctx, actionID, gomock.Any()).
                    Return(expected, nil)
            },
            wantErr: false,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveUserAction(ctx, actionID, gomock.Any()).
                    Return(nil, errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.ApproveUserAction(ctx, req, approved, actionID)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, expected, got)
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

    expected := []model.CPSAction{
        {ActionCode: "AC1"},
        {ActionCode: "AC2"},
    }

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
        want      []model.CPSAction
    }{
        {
            name: "success",
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetPendingUserActions(ctx).
                    Return(expected, nil)
            },
            wantErr: false,
            want:    expected,
        },
        {
            name: "repo error",
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetPendingUserActions(ctx).
                    Return(nil, errors.New("db error"))
            },
            wantErr: true,
            want:    nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.GetPendingUserActions(ctx)
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
    expected := &model.CPSUser{UserCode: userCode}

    tests := []struct {
        name      string
        mockSetup func()
        wantErr   bool
        want      *model.CPSUser
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