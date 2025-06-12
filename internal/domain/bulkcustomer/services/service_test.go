package services

import (
    "context"
    "errors"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/entities"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/repository/mocks"


)

func TestBranchService_FilterSingleBranches(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockBulkCustomerRepo(ctrl)
    service := NewBranchService(mockRepo)

    tests := []struct {
        name       string
        region     string
        district   string
        mockSetup  func()
        wantErr    bool
        wantResult []string
    }{
        {
            name:     "success",
            region:   "Addis",
            district: "Bole",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FilterSingleBranches(gomock.Any(), "Addis", "Bole").
                    Return([]string{"BR001", "BR002"}, nil)
            },
            wantErr:    false,
            wantResult: []string{"BR001", "BR002"},
        },
        {
            name:     "missing region",
            region:   "",
            district: "Bole",
            mockSetup: func() {},
            wantErr:    true,
            wantResult: nil,
        },
        {
            name:     "repo error",
            region:   "Addis",
            district: "Bole",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FilterSingleBranches(gomock.Any(), "Addis", "Bole").
                    Return(nil, errors.New("db error"))
            },
            wantErr:    true,
            wantResult: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.FilterSingleBranches(ctx, tt.region, tt.district)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResult, got)
            }
        })
    }
}

func TestBranchService_DisableSingleBranch(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockBulkCustomerRepo(ctrl)
    service := NewBranchService(mockRepo)

action := &entities.CPSAction{ID: "action1"}

    tests := []struct {
        name       string
        branchCode string
        cpsData    string
        mockSetup  func()
        wantErr    bool
        wantResult *entities.CPSAction
    }{
        {
            name:       "success",
            branchCode: "BR001",
            cpsData:    "data",
            mockSetup: func() {
                mockRepo.EXPECT().
                    DisableSingleBranch(gomock.Any(), "BR001", "data").
                    Return(action, nil)
            },
            wantErr:    false,
            wantResult: action,
        },
        {
            name:       "missing branchCode",
            branchCode: "",
            cpsData:    "data",
            mockSetup:  func() {},
            wantErr:    true,
            wantResult: nil,
        },
        {
            name:       "repo error",
            branchCode: "BR001",
            cpsData:    "data",
            mockSetup: func() {
                mockRepo.EXPECT().
                    DisableSingleBranch(gomock.Any(), "BR001", "data").
                    Return(nil, errors.New("db error"))
            },
            wantErr:    true,
            wantResult: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.DisableSingleBranch(ctx, tt.branchCode, tt.cpsData)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResult, got)
            }
        })
    }
}

func TestBranchService_ApproveSingleBranchDisable(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockBulkCustomerRepo(ctrl)
    service := NewBranchService(mockRepo)

    reason := "test reason"

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
            actionID: "action1",
            approve:  true,
            reason:   &reason,
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveSingleBranchDisable(gomock.Any(), "action1", true, &reason).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name:      "missing actionID",
            actionID:  "",
            approve:   true,
            reason:    &reason,
            mockSetup: func() {},
            wantErr:   true,
        },
        {
            name:     "repo error",
            actionID: "action1",
            approve:  false,
            reason:   &reason,
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveSingleBranchDisable(gomock.Any(), "action1", false, &reason).
                    Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := service.ApproveSingleBranchDisable(ctx, tt.actionID, tt.approve, tt.reason)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestBranchService_FilterMultipleBranches(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockBulkCustomerRepo(ctrl)
    service := NewBranchService(mockRepo)

    tests := []struct {
        name       string
        region     string
        district   string
        mockSetup  func()
        wantErr    bool
        wantResult []string
    }{
        {
            name:     "success",
            region:   "Addis",
            district: "Bole",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FilterMultipleBranches(gomock.Any(), "Addis", "Bole").
                    Return([]string{"BR003", "BR004"}, nil)
            },
            wantErr:    false,
            wantResult: []string{"BR003", "BR004"},
        },
        {
            name:     "missing district",
            region:   "Addis",
            district: "",
            mockSetup: func() {},
            wantErr:    true,
            wantResult: nil,
        },
        {
            name:     "repo error",
            region:   "Addis",
            district: "Bole",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FilterMultipleBranches(gomock.Any(), "Addis", "Bole").
                    Return(nil, errors.New("db error"))
            },
            wantErr:    true,
            wantResult: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.FilterMultipleBranches(ctx, tt.region, tt.district)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResult, got)
            }
        })
    }
}

func TestBranchService_DisableMultipleBranches(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockBulkCustomerRepo(ctrl)
    service := NewBranchService(mockRepo)

    action := &entities.CPSAction{ID: "action2"}

    tests := []struct {
        name        string
        branchCodes []string
        cpsData     string
        mockSetup   func()
        wantErr     bool
        wantResult  *entities.CPSAction
    }{
        {
            name:        "success",
            branchCodes: []string{"BR001", "BR002"},
            cpsData:     "data",
            mockSetup: func() {
                mockRepo.EXPECT().
                    DisableMultipleBranches(gomock.Any(), []string{"BR001", "BR002"}, "data").
                    Return(action, nil)
            },
            wantErr:    false,
            wantResult: action,
        },
        {
            name:        "missing branchCodes",
            branchCodes: []string{},
            cpsData:     "data",
            mockSetup:   func() {},
            wantErr:     true,
            wantResult:  nil,
        },
        {
            name:        "repo error",
            branchCodes: []string{"BR001"},
            cpsData:     "data",
            mockSetup: func() {
                mockRepo.EXPECT().
                    DisableMultipleBranches(gomock.Any(), []string{"BR001"}, "data").
                    Return(nil, errors.New("db error"))
            },
            wantErr:    true,
            wantResult: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.DisableMultipleBranches(ctx, tt.branchCodes, tt.cpsData)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResult, got)
            }
        })
    }
}

func TestBranchService_ApproveBulkBranchesDisable(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockBulkCustomerRepo(ctrl)
    service := NewBranchService(mockRepo)

    reason := "bulk reason"

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
            actionID: "action2",
            approve:  true,
            reason:   &reason,
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveBulkBranchesDisable(gomock.Any(), "action2", true, &reason).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name:      "missing actionID",
            actionID:  "",
            approve:   false,
            reason:    &reason,
            mockSetup: func() {},
            wantErr:   true,
        },
        {
            name:     "repo error",
            actionID: "action2",
            approve:  false,
            reason:   &reason,
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveBulkBranchesDisable(gomock.Any(), "action2", false, &reason).
                    Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := service.ApproveBulkBranchesDisable(ctx, tt.actionID, tt.approve, tt.reason)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}