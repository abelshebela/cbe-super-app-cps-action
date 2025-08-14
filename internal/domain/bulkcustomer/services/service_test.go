package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulkcustomer/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulkcustomer/repository/mocks"
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
		wantResult []entities.Branch
	}{
		{
			name:     "success",
			region:   "Addis",
			district: "Bole",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetBranch(ctx, "Addis", "Bole").
					Return([]entities.Branch{
						{BranchCode: "BR001", BranchName: "Branch 1"},
						{BranchCode: "BR002", BranchName: "Branch 2"},
					}, nil)
			},
			wantErr: false,
			wantResult: []entities.Branch{
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
					GetBranch(ctx, "Addis", "Bole").
					Return(nil, errors.New("db error"))
			},
			wantErr:    true,
			wantResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := service.GetBranch(ctx, tt.region, tt.district)
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

	branch := entities.Branch{BranchCode: "BR001"}
	maker := entities.User{UserCode: "U1"}
	tests := []struct {
		name      string
		branch    entities.Branch
		maker     entities.User
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
		wantResult []entities.Branch
	}{
		{
			name:     "success",
			region:   "Addis",
			district: "Bole",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetAllBranches(ctx, "Addis", "Bole").
					Return([]entities.Branch{
						{BranchCode: "BR003"},
						{BranchCode: "BR004"},
					}, nil)
			},
			wantErr: false,
			wantResult: []entities.Branch{
				{BranchCode: "BR003"},
				{BranchCode: "BR004"},
			},
		},
		{
			name:       "missing district",
			region:     "Addis",
			district:   "",
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
					GetAllBranches(ctx, "Addis", "Bole").
					Return(nil, errors.New("db error"))
			},
			wantErr:    true,
			wantResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := service.GetAllBranches(ctx, tt.region, tt.district)
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

	branches := []entities.Branch{
		{BranchCode: "BR001"},
		{BranchCode: "BR002"},
	}
	maker := entities.User{UserCode: "U1"}

	tests := []struct {
		name      string
		branches  []entities.Branch
		maker     entities.User
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
			branches:  []entities.Branch{},
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
			_, err := service.DisableMultipleBranches(ctx, tt.branches, tt.maker)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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

func TestBranchService_GetBranchByCode(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockBulkCustomerRepo(ctrl)
	service := NewBranchService(mockRepo)

	branch := entities.Branch{BranchCode: "BR001"}

	tests := []struct {
		name      string
		code      string
		mockSetup func()
		wantErr   bool
		want      entities.Branch
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
			name: "repo error",
			code: "BR002",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetBranchByCode(ctx, "BR002").
					Return(entities.Branch{}, errors.New("not found"))
			},
			wantErr: true,
			want:    entities.Branch{},
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
