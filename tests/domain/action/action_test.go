package test_action

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	actionDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/mocks"
)

func TestGetServicePaginated(t *testing.T) {
	mockRepo := new(mocks.ActionRepository)
	serviceStore := actionDomain.ServiceStore{Repository: mockRepo}

	ctx := context.Background()
	limit := 10
	offset := 0

	expectedServices := []actionDomain.ServiceDetails{
		{
			ID:          stringToPointer("1"),
			ServiceCode: "SVC001",
			ServiceName: "Service One",
			Enabled:     true,
		},
		{
			ID:          stringToPointer("2"),
			ServiceCode: "SVC002",
			ServiceName: "Service Two",
			Enabled:     false,
		},
	}

	mockRepo.On("GetAllHqServicesPaginated", ctx, offset, limit).Return(expectedServices, nil)

	services, err := serviceStore.GetServicePaginated(ctx, limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedServices, services)
	mockRepo.AssertCalled(t, "GetAllHqServicesPaginated", ctx, offset, limit)
}

func TestUpdateServiceFlagRequest(t *testing.T) {
	mockRepo := new(mocks.ActionRepository)
	serviceStore := actionDomain.ServiceStore{Repository: mockRepo}

	ctx := context.Background()
	id := "service-id"
	action := true
	maker := actionDomain.User{UserID: "maker-id"}

	service := actionDomain.ServiceDetails{
		ID:          stringToPointer(id),
		ServiceCode: "SVC001",
		ServiceName: "Service One",
		Enabled:     false,
	}

	mockRepo.On("GetHqServiceById", ctx, id).Return(service, nil)
	mockRepo.On("CreateCpsAction", ctx, mock.Anything).Return(actionDomain.CPSAction{ActionCode: "CPS_1234567890"}, nil)

	actionId, err := serviceStore.UpdateServiceFlagRequest(ctx, id, action, maker)

	assert.NoError(t, err)
	assert.Equal(t, "CPS_1234567890", actionId)
	mockRepo.AssertCalled(t, "GetHqServiceById", ctx, id)
	mockRepo.AssertCalled(t, "CreateCpsAction", ctx, mock.Anything)
}

func TestUpdateServiceFlag(t *testing.T) {
	mockRepo := new(mocks.ActionRepository)
	serviceStore := actionDomain.ServiceStore{Repository: mockRepo}

	ctx := context.Background()
	actionId := "action-id"
	action := true
	checkerId := "checker-id"

	cpsAction := actionDomain.CPSAction{
		ActionCode:   actionId,
		ActionStatus: actionDomain.ActionPending,
		CurrentAction: actionDomain.CurrentAction{
			Id:     []string{"service-id"},
			Action: action,
		},
	}

	service := actionDomain.ServiceDetails{
		ID:          stringToPointer("service-id"),
		ServiceCode: "SVC001",
		ServiceName: "Service One",
		Enabled:     false,
	}

	mockRepo.On("FetchCpsActionById", ctx, actionId).Return(cpsAction, nil)
	mockRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil)
	mockRepo.On("GetHqServiceById", ctx, "service-id").Return(service, nil)
	mockRepo.On("UpdateHqService", ctx, mock.Anything).Return(nil)

	err := serviceStore.UpdateServiceFlag(ctx, actionId, action, actionDomain.User{UserID: checkerId})

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "FetchCpsActionById", ctx, actionId)
	mockRepo.AssertCalled(t, "UpdateCpsAction", ctx, mock.Anything)
	mockRepo.AssertCalled(t, "GetHqServiceById", ctx, "service-id")
	mockRepo.AssertCalled(t, "UpdateHqService", ctx, mock.Anything)
}

func TestGetAccountByAccount(t *testing.T) {
	mockRepo := new(mocks.ActionRepository)
	serviceStore := actionDomain.ServiceStore{Repository: mockRepo}

	ctx := context.Background()
	account := "1234567890"

	expectedAccounts := []actionDomain.LinkedAccount{
		{
			AccountNumber:     "1234567890",
			AccountHolderName: "John Doe",
			LinkedStatus:      true,
		},
	}

	mockRepo.On("FetchAccountsByAccountNumber", ctx, account).Return(expectedAccounts, nil)

	accounts, err := serviceStore.GetAccountByAccount(ctx, account)

	assert.NoError(t, err)
	assert.Equal(t, expectedAccounts, accounts)
	mockRepo.AssertCalled(t, "FetchAccountsByAccountNumber", ctx, account)
}

func TestRemoveCifRequest(t *testing.T) {
	mockRepo := new(mocks.ActionRepository)
	serviceStore := actionDomain.ServiceStore{Repository: mockRepo}

	ctx := context.Background()
	ids := []string{"account-id-1", "account-id-2"}
	action := false
	maker := actionDomain.User{UserID: "maker-id"}

	mockRepo.On("FetchLastCpsActionByMakerID", ctx, maker.UserID).Return(actionDomain.CPSAction{ActionStatus: actionDomain.ActionApproved}, nil)
	mockRepo.On("CreateCpsAction", ctx, mock.Anything).Return(actionDomain.CPSAction{ActionCode: "CPS_1234567890"}, nil)

	actionId, err := serviceStore.RemoveCifRequest(ctx, ids, action, maker)

	assert.NoError(t, err)
	assert.Equal(t, "CPS_1234567890", actionId)
	mockRepo.AssertCalled(t, "FetchLastCpsActionByMakerID", ctx, maker.UserID)
	mockRepo.AssertCalled(t, "CreateCpsAction", ctx, mock.Anything)
}

func TestRemoveCif(t *testing.T) {
	mockRepo := new(mocks.ActionRepository)
	serviceStore := actionDomain.ServiceStore{Repository: mockRepo}

	ctx := context.Background()
	actionId := "action-id"
	action := false
	checker := actionDomain.User{UserID: "checker-id"}

	cpsAction := actionDomain.CPSAction{
		ActionCode:   actionId,
		ActionStatus: actionDomain.ActionPending,
		CurrentAction: actionDomain.CurrentAction{
			Id:     []string{"account-id-1", "account-id-2"},
			Action: action,
		},
	}

	linkedAccounts := []actionDomain.LinkedAccount{
		{
			AccountNumber: "1234567890",
			LinkedStatus:  true,
		},
		{
			AccountNumber: "0987654321",
			LinkedStatus:  true,
		},
	}

	mockRepo.On("FetchCpsActionById", ctx, actionId).Return(cpsAction, nil)
	mockRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil)
	mockRepo.On("FetchLinkedAccountById", ctx, []string{"account-id-1", "account-id-2"}).Return(linkedAccounts, nil)
	mockRepo.On("UpdateAccount", ctx, mock.Anything).Return(actionDomain.LinkedAccount{}, nil)

	_, err := serviceStore.RemoveCif(ctx, actionId, action, "",checker)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "FetchCpsActionById", ctx, actionId)
	mockRepo.AssertCalled(t, "UpdateCpsAction", ctx, mock.Anything)
	mockRepo.AssertCalled(t, "FetchLinkedAccountById", ctx, []string{"account-id-1", "account-id-2"})
	mockRepo.AssertCalled(t, "UpdateAccount", ctx, mock.Anything)
}

func stringToPointer(s string) *string {
	return &s
}
