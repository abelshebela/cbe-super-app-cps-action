package domain

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
// 	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
// )

// type MockPasswordRuleRepository struct {
// 	GetCurrentPasswordRuleFunc          func(ctx context.Context) (*action.PasswordRule, error)
// 	FetchPendingActionsByUniqueIDFunc   func(ctx context.Context, uniqueID string) ([]action.ActionResponse, error)
// 	CreateCpsActionFunc                 func(ctx context.Context, a action.CPSAction) (action.CPSAction, error)
// 	GetPasswordRuleUpdateActionByIDFunc func(ctx context.Context, actionID string) (*action.CPSAction, error)
// 	UpdatePasswordRuleFunc              func(ctx context.Context, rule action.PasswordRule) error
// 	UpdateCpsActionFunc                 func(ctx context.Context, a action.CPSAction) error
// }

// func (m *MockPasswordRuleRepository) GetCurrentPasswordRule(ctx context.Context) (*action.PasswordRule, error) {
// 	return m.GetCurrentPasswordRuleFunc(ctx)
// }
// func (m *MockPasswordRuleRepository) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.ActionResponse, error) {
// 	return m.FetchPendingActionsByUniqueIDFunc(ctx, uniqueID)
// }
// func (m *MockPasswordRuleRepository) CreateCpsAction(ctx context.Context, a action.CPSAction) (action.CPSAction, error) {
// 	return m.CreateCpsActionFunc(ctx, a)
// }
// func (m *MockPasswordRuleRepository) GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error) {
// 	return m.GetPasswordRuleUpdateActionByIDFunc(ctx, actionID)
// }
// func (m *MockPasswordRuleRepository) UpdatePasswordRule(ctx context.Context, rule action.PasswordRule) error {
// 	return m.UpdatePasswordRuleFunc(ctx, rule)
// }
// func (m *MockPasswordRuleRepository) UpdateCpsAction(ctx context.Context, a action.CPSAction) error {
// 	return m.UpdateCpsActionFunc(ctx, a)
// }
// func (m *MockPasswordRuleRepository) GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error) {
// 	return nil, nil
// }

// func TestPasswordRuleService_RequestPasswordRuleUpdate(t *testing.T) {
// 	repo := &MockPasswordRuleRepository{
// 		GetCurrentPasswordRuleFunc: func(ctx context.Context) (*action.PasswordRule, error) {
// 			return &action.PasswordRule{ID: "rule-1"}, nil
// 		},
// 		FetchPendingActionsByUniqueIDFunc: func(ctx context.Context, uniqueID string) ([]action.ActionResponse, error) {
// 			return nil, nil
// 		},
// 		CreateCpsActionFunc: func(ctx context.Context, a action.CPSAction) (action.CPSAction, error) {
// 			return a, nil
// 		},
// 	}
// 	service := services.NewPasswordRuleService(repo)
// 	maker := action.User{UserID: "maker-1", FullName: "Maker", PhoneNumber: "123"}
// 	rule := &action.PasswordRule{ID: "rule-1"}
// 	code, err := service.RequestPasswordRuleUpdate(context.Background(), rule, maker, "IT")
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// 	if code == "" {
// 		t.Error("expected action code, got empty string")
// 	}
// }

// func TestPasswordRuleService_RequestPasswordRuleUpdate_Errors(t *testing.T) {
// 	repo := &MockPasswordRuleRepository{GetCurrentPasswordRuleFunc: func(ctx context.Context) (*action.PasswordRule, error) {
// 		return nil, errors.New("not found")
// 	},
// 		FetchPendingActionsByUniqueIDFunc: func(ctx context.Context, uniqueID string) ([]action.ActionResponse, error) {
// 			return nil, nil
// 		},
// 		CreateCpsActionFunc: func(ctx context.Context, a action.CPSAction) (action.CPSAction, error) {
// 			return a, nil
// 		},
// 	}
// 	service := services.NewPasswordRuleService(repo)
// 	maker := action.User{UserID: "maker-1", FullName: "Maker", PhoneNumber: "123"}
// 	_, err := service.RequestPasswordRuleUpdate(context.Background(), &action.PasswordRule{ID: "rule-1"}, maker, "IT")
// 	if err == nil {
// 		t.Error("expected error for not found rule")
// 	}
// }

// func TestPasswordRuleService_ApproveOrRejectPasswordRuleAction(t *testing.T) {
// 	approved := false
// 	repo := &MockPasswordRuleRepository{
// 		GetPasswordRuleUpdateActionByIDFunc: func(ctx context.Context, actionID string) (*action.CPSAction, error) {
// 			return &action.CPSAction{
// 				ActionStatus:  action.ActionPending,
// 				CurrentAction: &action.PasswordRule{ID: "rule-1", MinLength: 6, MaxLength: 12, Numbers: true, CapitalLetters: true, SmallLetters: true, Characters: true},
// 			}, nil
// 		},
// 		UpdatePasswordRuleFunc: func(ctx context.Context, rule action.PasswordRule) error {
// 			approved = true
// 			return nil
// 		},
// 		UpdateCpsActionFunc: func(ctx context.Context, a action.CPSAction) error {
// 			return nil
// 		},
// 	}
// 	service := services.NewPasswordRuleService(repo)
// 	checker := action.User{UserID: "checker-1", FullName: "Checker", PhoneNumber: "456"}
// 	err := service.ApproveOrRejectPasswordRuleAction(context.Background(), "action-1", "APPROVED", checker, nil, "IT")
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// 	if !approved {
// 		t.Error("expected UpdatePasswordRule to be called")
// 	}
// }

// func TestPasswordRuleService_ApproveOrRejectPasswordRuleAction_Denied(t *testing.T) {
// 	repo := &MockPasswordRuleRepository{
// 		GetPasswordRuleUpdateActionByIDFunc: func(ctx context.Context, actionID string) (*action.CPSAction, error) {
// 			return &action.CPSAction{ActionStatus: action.ActionPending}, nil
// 		},
// 		UpdateCpsActionFunc: func(ctx context.Context, a action.CPSAction) error {
// 			return nil
// 		},
// 	}
// 	service := services.NewPasswordRuleService(repo)
// 	checker := action.User{UserID: "checker-1", FullName: "Checker", PhoneNumber: "456"}
// 	reason := "not valid"
// 	err := service.ApproveOrRejectPasswordRuleAction(context.Background(), "action-1", "DENIED", checker, &reason, "IT")
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// }

// func TestPasswordRuleService_ApproveOrRejectPasswordRuleAction_InvalidDecision(t *testing.T) {
// 	repo := &MockPasswordRuleRepository{
// 		GetPasswordRuleUpdateActionByIDFunc: func(ctx context.Context, actionID string) (*action.CPSAction, error) {
// 			return &action.CPSAction{ActionStatus: action.ActionPending}, nil
// 		},
// 	}
// 	service := services.NewPasswordRuleService(repo)
// 	checker := action.User{UserID: "checker-1", FullName: "Checker", PhoneNumber: "456"}
// 	err := service.ApproveOrRejectPasswordRuleAction(context.Background(), "action-1", "INVALID", checker, nil, "IT")
// 	if err == nil {
// 		t.Error("expected error for invalid decision")
// 	}
// }

// func TestPasswordRuleService_CheckPasswordRule(t *testing.T) {
// 	repo := &MockPasswordRuleRepository{
// 		GetCurrentPasswordRuleFunc: func(ctx context.Context) (*action.PasswordRule, error) {
// 			return &action.PasswordRule{ID: "rule-1", MinLength: 6, MaxLength: 12, Numbers: true, CapitalLetters: true, SmallLetters: true, Characters: true}, nil

// 		},
// 	}
// 	service := services.NewPasswordRuleService(repo)
// 	ok, msg := service.CheckPasswordRule(context.Background(), "Abc123!")
// 	if !ok {
// 		t.Errorf("expected password to be valid, got: %s", msg)
// 	}
// 	ok, msg = service.CheckPasswordRule(context.Background(), "abc")
// 	if ok {
// 		t.Error("expected password to be invalid")
// 	}
// }
