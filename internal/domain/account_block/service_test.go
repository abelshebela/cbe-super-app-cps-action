package account_block

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"

    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_block/mocks"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

func TestAccountService_FilterSingleBranches(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

    tests := []struct {
        name       string
        region     string
        district   string
        mockSetup  func()
        wantErr    bool
        wantResult []action.Branch
    }{
        {
            name:     "success",
            region:   "Addis",
            district: "Bole",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FilterSingleBranches(ctx, "Addis", "Bole").
                    Return([]action.Branch{
                        {BranchCode: "BR001", BranchName: "Branch 1"},
                        {BranchCode: "BR002", BranchName: "Branch 2"},
                    }, nil)
            },
            wantErr: false,
            wantResult: []action.Branch{
                {BranchCode: "BR001", BranchName: "Branch 1"},
                {BranchCode: "BR002", BranchName: "Branch 2"},
            },
        },
        {
            name:       "missing region",
            region:     "",
            district:   "Bole",
            mockSetup:  func() {},
            wantErr:    true,
            wantResult: nil,
        },
        {
            name:     "repo error",
            region:   "Addis",
            district: "Bole",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FilterSingleBranches(ctx, "Addis", "Bole").
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

func TestAccountService_DisableSingleBranch(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

    branch := action.Branch{BranchCode: "BR001"}
    maker := action.User{UserCode: "U1"}
    tests := []struct {
        name      string
        branch    action.Branch
        maker     action.User
        mockSetup func()
        wantErr   bool
    }{
        {
            name:   "success",
            branch: branch,
            maker:  maker,
            mockSetup: func() {
                mockRepo.EXPECT().
                    DisableSingleBranch(ctx, branch, maker).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name:   "repo error",
            branch: branch,
            maker:  maker,
            mockSetup: func() {
                mockRepo.EXPECT().
                    DisableSingleBranch(ctx, branch, maker).
                    Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := service.DisableSingleBranch(ctx, tt.branch, tt.maker)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestAccountService_ApproveSingleBranchDisable(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

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
                    ApproveSingleBranchDisable(ctx, "action1", true, &reason).
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
                    ApproveSingleBranchDisable(ctx, "action1", false, &reason).
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

func TestAccountService_FilterMultipleBranches(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

    tests := []struct {
        name       string
        region     string
        district   string
        mockSetup  func()
        wantErr    bool
        wantResult []action.Branch
    }{
        {
            name:     "success",
            region:   "Addis",
            district: "Bole",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FilterMultipleBranches(ctx, "Addis", "Bole").
                    Return([]action.Branch{
                        {BranchCode: "BR003"},
                        {BranchCode: "BR004"},
                    }, nil)
            },
            wantErr: false,
            wantResult: []action.Branch{
                {BranchCode: "BR003"},
                {BranchCode: "BR004"},
            },
        },
        {
            name:       "missing region",
            region:     "",
            district:   "Bole",
            mockSetup:  func() {},
            wantErr:    true,
            wantResult: nil,
        },
        {
            name:     "repo error",
            region:   "Addis",
            district: "Bole",
            mockSetup: func() {
                mockRepo.EXPECT().
                    FilterMultipleBranches(ctx, "Addis", "Bole").
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

func TestAccountService_DisableMultipleBranches(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

    branches := []action.Branch{
        {BranchCode: "BR001"},
        {BranchCode: "BR002"},
    }
    maker := action.User{UserCode: "U1"}

    tests := []struct {
        name      string
        branches  []action.Branch
        maker     action.User
        mockSetup func()
        wantErr   bool
    }{
        {
            name:     "success",
            branches: branches,
            maker:    maker,
            mockSetup: func() {
                mockRepo.EXPECT().
                    DisableMultipleBranches(ctx, branches, maker).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name:      "missing branches",
            branches:  []action.Branch{},
            maker:     maker,
            mockSetup: func() {},
            wantErr:   true,
        },
        {
            name:     "repo error",
            branches: branches,
            maker:    maker,
            mockSetup: func() {
                mockRepo.EXPECT().
                    DisableMultipleBranches(ctx, branches, maker).
                    Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := service.DisableMultipleBranches(ctx, tt.branches, tt.maker)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestAccountService_ApproveBulkBranchesDisable(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

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
                    ApproveBulkBranchesDisable(ctx, "action2", true, &reason).
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
                    ApproveBulkBranchesDisable(ctx, "action2", false, &reason).
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

func TestAccountService_GetBranchByCode(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

    branch := action.Branch{BranchCode: "BR001"}

    tests := []struct {
        name      string
        code      string
        mockSetup func()
        wantErr   bool
        want      action.Branch
    }{
        {
            name: "success",
            code: "BR001",
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetBranchByCode(ctx, "BR001").
                    Return(branch, nil)
            },
            wantErr: false,
            want:    branch,
        },
        {
            name: "missing branchCode",
            code: "",
            mockSetup: func() {},
            wantErr: true,
            want:    action.Branch{},
        },
        {
            name: "repo error",
            code: "BR002",
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetBranchByCode(ctx, "BR002").
                    Return(action.Branch{}, errors.New("not found"))
            },
            wantErr: true,
            want:    action.Branch{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            got, err := service.GetBranchByCode(ctx, tt.code)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}

func TestAccountService_BlockRegion(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

    region := action.Region{ID: "R1", RegionName: "Addis"}
    maker := action.CPSAction{Maker: action.User{UserCode: "U1"}}
    branches := []action.Branch{{BranchCode: "BR1"}, {BranchCode: "BR2"}}

    t.Run("success", func(t *testing.T) {
        mockRepo.EXPECT().UpdateRegion(ctx, gomock.Any()).Return(nil)
        mockRepo.EXPECT().FilterMultipleBranches(ctx, "Addis", "").Return(branches, nil)
        mockRepo.EXPECT().DisableMultipleBranches(ctx, gomock.Any(), maker.Maker).Return(nil)
        mockRepo.EXPECT().BlockRegion(ctx, "R1", maker).Return(nil)

        err := service.BlockRegion(ctx, region, maker)
        assert.NoError(t, err)
    })

    t.Run("missing region ID", func(t *testing.T) {
        err := service.BlockRegion(ctx, action.Region{ID: "", RegionName: ""}, maker)
        assert.Error(t, err)
    })

    t.Run("update region error", func(t *testing.T) {
        mockRepo.EXPECT().UpdateRegion(ctx, gomock.Any()).Return(errors.New("fail"))
        err := service.BlockRegion(ctx, region, maker)
        assert.Error(t, err)
    })

    t.Run("filter branches error", func(t *testing.T) {
        mockRepo.EXPECT().UpdateRegion(ctx, gomock.Any()).Return(nil)
        mockRepo.EXPECT().FilterMultipleBranches(ctx, "Addis", "").Return(nil, errors.New("fail"))
        err := service.BlockRegion(ctx, region, maker)
        assert.Error(t, err)
    })

    t.Run("disable branches error", func(t *testing.T) {
        mockRepo.EXPECT().UpdateRegion(ctx, gomock.Any()).Return(nil)
        mockRepo.EXPECT().FilterMultipleBranches(ctx, "Addis", "").Return(branches, nil)
        mockRepo.EXPECT().DisableMultipleBranches(ctx, gomock.Any(), maker.Maker).Return(errors.New("fail"))
        err := service.BlockRegion(ctx, region, maker)
        assert.Error(t, err)
    })

    t.Run("block region error", func(t *testing.T) {
        mockRepo.EXPECT().UpdateRegion(ctx, gomock.Any()).Return(nil)
        mockRepo.EXPECT().FilterMultipleBranches(ctx, "Addis", "").Return(branches, nil)
        mockRepo.EXPECT().DisableMultipleBranches(ctx, gomock.Any(), maker.Maker).Return(nil)
        mockRepo.EXPECT().BlockRegion(ctx, "R1", maker).Return(errors.New("fail"))
        err := service.BlockRegion(ctx, region, maker)
        assert.Error(t, err)
    })
}

func TestAccountService_ApproveRegionBlock(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

    reason := "region reason"

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
            actionID: "action3",
            approve:  true,
            reason:   &reason,
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveRegionBlock(ctx, "action3", true, &reason).
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
            actionID: "action3",
            approve:  false,
            reason:   &reason,
            mockSetup: func() {
                mockRepo.EXPECT().
                    ApproveRegionBlock(ctx, "action3", false, &reason).
                    Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()
            err := service.ApproveRegionBlock(ctx, tt.actionID, tt.approve, tt.reason)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestAccountService_UpdateRegion(t *testing.T) {
    ctx := context.Background()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mocks.NewMockAccountBlockRepo(ctrl)
    service := NewAccountService(mockRepo)

    region := action.Region{ID: "R1", RegionName: "Addis", UpdatedAt: time.Now()}

    t.Run("success", func(t *testing.T) {
        mockRepo.EXPECT().UpdateRegion(ctx, gomock.Any()).Return(nil)
        err := service.UpdateRegion(ctx, region)
        assert.NoError(t, err)
    })

    t.Run("missing region ID", func(t *testing.T) {
        err := service.UpdateRegion(ctx, action.Region{ID: "", RegionName: ""})
        assert.Error(t, err)
    })

    t.Run("repo error", func(t *testing.T) {
        mockRepo.EXPECT().UpdateRegion(ctx, gomock.Any()).Return(errors.New("fail"))
        err := service.UpdateRegion(ctx, region)
        assert.Error(t, err)
    })
}