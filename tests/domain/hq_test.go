package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	hqdomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	actionmock "github.com/CBE-Super-App/cbe-super-app-cps-action/mocks/domain/action/action_manual_mock"
	hqmock "github.com/CBE-Super-App/cbe-super-app-cps-action/mocks/domain/hq"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// CustomActionRepository implements the action.Repository interface with additional methods
type CustomActionRepository struct {
	*actionmock.MockActionRepository
	FetchPendingActionsByUniqueIDFunc func(ctx context.Context, uniqueID string) ([]action.ActionResponse, error)
}

func (m *CustomActionRepository) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.ActionResponse, error) {
	if m.FetchPendingActionsByUniqueIDFunc != nil {
		return m.FetchPendingActionsByUniqueIDFunc(ctx, uniqueID)
	}
	return []action.ActionResponse{}, nil
}

func TestService_GetHQ(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockHQ        hq.HQ
		mockError     error
		expectedError bool
	}{
		{
			name:          "success",
			id:            "test-id",
			mockHQ:        hq.HQ{ID: "test-id", Name: "Test HQ"},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "empty id",
			id:            "",
			mockHQ:        hq.HQ{},
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "repository error",
			id:            "test-id",
			mockHQ:        hq.HQ{},
			mockError:     errors.New("repository error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &hqmock.MockRepository{
				GetHQByIDFunc: func(ctx context.Context, id string) (hq.HQ, error) {
					return tt.mockHQ, tt.mockError
				},
			}

			mockActionRepo := &CustomActionRepository{
				MockActionRepository: &actionmock.MockActionRepository{},
			}
			logger := sharedutils.NewLogger()

			service := hqdomain.NewService(mockRepo, mockActionRepo, logger)
			hq, err := service.GetHQ(context.Background(), tt.id)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if hq.ID != tt.mockHQ.ID {
					t.Errorf("expected HQ ID %s, got %s", tt.mockHQ.ID, hq.ID)
				}
			}
		})
	}
}

func TestService_UpdateBlockTimeRequest(t *testing.T) {
	tests := []struct {
		name          string
		request       hq.UpdateBlockTimeRequest
		mockHQ        hq.HQ
		mockError     error
		expectedError bool
	}{
		{
			name: "success",
			request: hq.UpdateBlockTimeRequest{
				ID:         "test-id",
				MakerID:    "test-maker",
				MakerName:  "Test Maker",
				MakerPhone: "1234567890",
				BlockTime:  3600, // 1 hour in seconds
			},
			mockHQ: hq.HQ{
				ID:        "test-id",
				Name:      "Test HQ",
				BlockTime: 1800, // 30 minutes in seconds
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "repository error",
			request: hq.UpdateBlockTimeRequest{
				ID:         "test-id",
				MakerID:    "test-maker",
				MakerName:  "Test Maker",
				MakerPhone: "1234567890",
				BlockTime:  3600,
			},
			mockHQ:        hq.HQ{},
			mockError:     errors.New("repository error"),
			expectedError: true,
		},
		{
			name: "empty id",
			request: hq.UpdateBlockTimeRequest{
				ID:         "",
				MakerID:    "test-maker",
				MakerName:  "Test Maker",
				MakerPhone: "1234567890",
				BlockTime:  3600,
			},
			mockHQ:        hq.HQ{},
			mockError:     nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &hqmock.MockRepository{
				GetHQByIDFunc: func(ctx context.Context, id string) (hq.HQ, error) {
					return tt.mockHQ, tt.mockError
				},
			}

			mockActionRepo := &CustomActionRepository{
				MockActionRepository: &actionmock.MockActionRepository{
					CreateCpsActionFunc: func(ctx context.Context, a action.CPSAction) (action.CPSAction, error) {
						return a, nil
					},
				},
				FetchPendingActionsByUniqueIDFunc: func(ctx context.Context, uniqueID string) ([]action.ActionResponse, error) {
					return []action.ActionResponse{}, nil
				},
			}

			logger := sharedutils.NewLogger()
			service := hqdomain.NewService(mockRepo, mockActionRepo, logger)
			actionCode, err := service.UpdateBlockTimeRequest(context.Background(), tt.request)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if actionCode == "" {
					t.Error("expected action code but got empty string")
				}
			}
		})
	}
}

func TestService_UpdateArchiveTimeRequest(t *testing.T) {
	tests := []struct {
		name          string
		request       hq.UpdateArchiveTimeRequest
		mockHQ        hq.HQ
		mockError     error
		expectedError bool
	}{
		{
			name: "success",
			request: hq.UpdateArchiveTimeRequest{
				ID:          "test-id",
				MakerID:     "test-maker",
				MakerName:   "Test Maker",
				MakerPhone:  "1234567890",
				ArchiveTime: 86400, // 24 hours in seconds
			},
			mockHQ: hq.HQ{
				ID:          "test-id",
				Name:        "Test HQ",
				ArchiveTime: 43200, // 12 hours in seconds
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "repository error",
			request: hq.UpdateArchiveTimeRequest{
				ID:          "test-id",
				MakerID:     "test-maker",
				MakerName:   "Test Maker",
				MakerPhone:  "1234567890",
				ArchiveTime: 86400,
			},
			mockHQ:        hq.HQ{},
			mockError:     errors.New("repository error"),
			expectedError: true,
		},
		{
			name: "empty id",
			request: hq.UpdateArchiveTimeRequest{
				ID:          "",
				MakerID:     "test-maker",
				MakerName:   "Test Maker",
				MakerPhone:  "1234567890",
				ArchiveTime: 86400,
			},
			mockHQ:        hq.HQ{},
			mockError:     nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &hqmock.MockRepository{
				GetHQByIDFunc: func(ctx context.Context, id string) (hq.HQ, error) {
					return tt.mockHQ, tt.mockError
				},
			}

			mockActionRepo := &CustomActionRepository{
				MockActionRepository: &actionmock.MockActionRepository{
					CreateCpsActionFunc: func(ctx context.Context, a action.CPSAction) (action.CPSAction, error) {
						return a, nil
					},
				},
				FetchPendingActionsByUniqueIDFunc: func(ctx context.Context, uniqueID string) ([]action.ActionResponse, error) {
					return []action.ActionResponse{}, nil
				},
			}

			logger := sharedutils.NewLogger()
			service := hqdomain.NewService(mockRepo, mockActionRepo, logger)
			actionCode, err := service.UpdateArchiveTimeRequest(context.Background(), tt.request)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if actionCode == "" {
					t.Error("expected action code but got empty string")
				}
			}
		})
	}
}

func TestService_UpdateBlockTime(t *testing.T) {
	tests := []struct {
		name          string
		request       hq.ApproveRejectRequest
		mockAction    action.CPSAction
		mockError     error
		expectedError bool
	}{
		{
			name: "success approve",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction: action.CPSAction{
				ActionCode:   "test-action",
				ActionStatus: action.ActionPending,
				CurrentAction: func() interface{} {
					type CurrentAction struct {
						HQ hq.HQ `json:"hq"`
					}
					currentAction := CurrentAction{
						HQ: hq.HQ{
							ID:        "test-id",
							Name:      "Test HQ",
							BlockTime: 3600,
						},
					}
					return currentAction
				}(),
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "success reject",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionDenied,
				RejectedReason: "Not needed",
			},
			mockAction: action.CPSAction{
				ActionCode:   "test-action",
				ActionStatus: action.ActionPending,
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "reject missing reason",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionDenied,
				RejectedReason: "",
			},
			mockAction: action.CPSAction{
				ActionCode:   "test-action",
				ActionStatus: action.ActionPending,
			},
			mockError:     nil,
			expectedError: false, // Service allows empty rejection reason
		},
		{
			name: "action not found",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction:    action.CPSAction{},
			mockError:     errors.New("action not found"),
			expectedError: true,
		},
		{
			name: "empty action code",
			request: hq.ApproveRejectRequest{
				ActionCode:     "",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction:    action.CPSAction{},
			mockError:     nil,
			expectedError: true,
		},
		{
			name: "empty checker id",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction:    action.CPSAction{},
			mockError:     nil,
			expectedError: true,
		},
		{
			name: "action not pending",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction: action.CPSAction{
				ActionCode:   "test-action",
				ActionStatus: action.ActionApproved,
			},
			mockError:     nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &hqmock.MockRepository{
				UpdateHQFunc: func(ctx context.Context, id string, hq hq.HQ) error {
					return nil
				},
			}

			mockActionRepo := &CustomActionRepository{
				MockActionRepository: &actionmock.MockActionRepository{
					FetchCpsActionByIdFunc: func(ctx context.Context, id string) (action.CPSAction, error) {
						return tt.mockAction, tt.mockError
					},
					UpdateCpsActionFunc: func(ctx context.Context, a action.CPSAction) error {
						return nil
					},
				},
			}

			logger := sharedutils.NewLogger()
			service := hqdomain.NewService(mockRepo, mockActionRepo, logger)
			err := service.UpdateBlockTime(context.Background(), tt.request)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestService_UpdateArchiveTime(t *testing.T) {
	tests := []struct {
		name          string
		request       hq.ApproveRejectRequest
		mockAction    action.CPSAction
		mockError     error
		expectedError bool
	}{
		{
			name: "success approve",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction: action.CPSAction{
				ActionCode:   "test-action",
				ActionStatus: action.ActionPending,
				CurrentAction: func() interface{} {
					type CurrentAction struct {
						HQ hq.HQ `json:"hq"`
					}
					currentAction := CurrentAction{
						HQ: hq.HQ{
							ID:          "test-id",
							Name:        "Test HQ",
							ArchiveTime: 86400,
						},
					}
					return currentAction
				}(),
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "success reject",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionDenied,
				RejectedReason: "Not needed",
			},
			mockAction: action.CPSAction{
				ActionCode:   "test-action",
				ActionStatus: action.ActionPending,
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "reject missing reason",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionDenied,
				RejectedReason: "",
			},
			mockAction: action.CPSAction{
				ActionCode:   "test-action",
				ActionStatus: action.ActionPending,
			},
			mockError:     nil,
			expectedError: false, // Service allows empty rejection reason
		},
		{
			name: "action not found",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction:    action.CPSAction{},
			mockError:     errors.New("action not found"),
			expectedError: true,
		},
		{
			name: "empty action code",
			request: hq.ApproveRejectRequest{
				ActionCode:     "",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction:    action.CPSAction{},
			mockError:     nil,
			expectedError: true,
		},
		{
			name: "empty checker id",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction:    action.CPSAction{},
			mockError:     nil,
			expectedError: true,
		},
		{
			name: "action not pending",
			request: hq.ApproveRejectRequest{
				ActionCode:     "test-action",
				CheckerID:      "test-checker",
				CheckerName:    "Test Checker",
				CheckerPhone:   "0987654321",
				Decision:       utils.DecisionApproved,
				RejectedReason: "",
			},
			mockAction: action.CPSAction{
				ActionCode:   "test-action",
				ActionStatus: action.ActionApproved,
			},
			mockError:     nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &hqmock.MockRepository{
				UpdateHQFunc: func(ctx context.Context, id string, hq hq.HQ) error {
					return nil
				},
			}

			mockActionRepo := &CustomActionRepository{
				MockActionRepository: &actionmock.MockActionRepository{
					FetchCpsActionByIdFunc: func(ctx context.Context, id string) (action.CPSAction, error) {
						return tt.mockAction, tt.mockError
					},
					UpdateCpsActionFunc: func(ctx context.Context, a action.CPSAction) error {
						return nil
					},
				},
			}

			logger := sharedutils.NewLogger()
			service := hqdomain.NewService(mockRepo, mockActionRepo, logger)
			err := service.UpdateArchiveTime(context.Background(), tt.request)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
