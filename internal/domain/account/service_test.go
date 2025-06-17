package account_test

import (
	"context"
	"testing"

	"cbe-super-app-member-users/internal/domain/account"
	"cbe-super-app-member-users/internal/shared"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	accountPortMocks "cbe-super-app-member-users/internal/port/outbound/account/mocks"
	accountPort "cbe-super-app-member-users/internal/port/outbound/account"
)

func MockPortAccountUser() *accountPort.AccountUser {
	return &accountPort.AccountUser{
		ID:               "user-123",
		PhoneNumber:      "1234567890",
		KYCLevel:         1,
		FullName:         "John Doe",
		BranchCode:       "BR001",
		RegistrationType: "Savings",
		AndOrStatus:      true,
	}
}

func TestAccountService_CreateAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := accountPortMocks.NewMockAccountRepositoryPort(ctrl)
	mockAPIClient := accountPortMocks.NewMockAccountAPIPort(ctrl)
	logger := utils.NewLogger()

	service := account.NewAccountService(mockRepo, mockAPIClient, logger)

	ctx := context.Background()
	userID := "user-123"
	user := MockPortAccountUser()

	mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(user, nil)
	mockAPIClient.EXPECT().LookupAccountByPhone(ctx, user.PhoneNumber).Return(false, nil)
	mockRepo.EXPECT().UpdateUserCustomerNumber(ctx, userID, gomock.Any()).Return(nil)
	mockRepo.EXPECT().CreateLinkedAccount(ctx, gomock.Any()).Return(nil)

	result, err := service.CreateAccount(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.CustomerNumber)
	assert.NotEmpty(t, result.AccountNumber)
}

func TestAccountService_CreateAccount_KYCLevelZero(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := accountPortMocks.NewMockAccountRepositoryPort(ctrl)
	mockAPIClient := accountPortMocks.NewMockAccountAPIPort(ctrl)
	logger := utils.NewLogger()

	service := account.NewAccountService(mockRepo, mockAPIClient, logger)

	ctx := context.Background()
	userID := "user-123"
	user := MockPortAccountUser()
	user.KYCLevel = 0

	mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(user, nil)

	result, err := service.CreateAccount(ctx, userID)

	require.Error(t, err)
	require.Nil(t, result)
	serviceErr := err.(account.ServiceError)
	assert.Equal(t, common.DefineError.User["USER_KYC_LEVEL_ZERO"].Code, serviceErr.Code)
}

func TestAccountService_CreateAccount_PhoneExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := accountPortMocks.NewMockAccountRepositoryPort(ctrl)
	mockAPIClient := accountPortMocks.NewMockAccountAPIPort(ctrl)
	logger := utils.NewLogger()

	service := account.NewAccountService(mockRepo, mockAPIClient, logger)

	ctx := context.Background()
	userID := "user-123"
	user := MockPortAccountUser()

	mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(user, nil)
	mockAPIClient.EXPECT().LookupAccountByPhone(ctx, user.PhoneNumber).Return(true, nil)

	result, err := service.CreateAccount(ctx, userID)

	require.Error(t, err)
	require.Nil(t, result)
	serviceErr := err.(account.ServiceError)
	assert.Equal(t, common.DefineError.User["USER_PHONE_EXISTS"].Code, serviceErr.Code)
}

func TestAccountService_CreateAccount_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := accountPortMocks.NewMockAccountRepositoryPort(ctrl)
	mockAPIClient := accountPortMocks.NewMockAccountAPIPort(ctrl)
	logger := utils.NewLogger()

	service := account.NewAccountService(mockRepo, mockAPIClient, logger)

	ctx := context.Background()
	userID := "user-123"

	mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(nil, shared.ErrNotFound)

	result, err := service.CreateAccount(ctx, userID)

	require.Error(t, err)
	require.Nil(t, result)
	serviceErr := err.(account.ServiceError)
	assert.Equal(t, common.DefineError.General["NOT_FOUND"].Code, serviceErr.Code)
}

func TestAccountService_GenerateSIFAndAccountNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := accountPortMocks.NewMockAccountRepositoryPort(ctrl)
	mockAPIClient := accountPortMocks.NewMockAccountAPIPort(ctrl)
	logger := utils.NewLogger()

	service := account.NewAccountService(mockRepo, mockAPIClient, logger)

	sif, accountNumber, err := service.GenerateSIFAndAccountNumber()

	require.NoError(t, err)
	assert.NotEmpty(t, sif)
	assert.NotEmpty(t, accountNumber)
	assert.Contains(t, sif, "SIF-")
	assert.Contains(t, accountNumber, "ACC-")
	assert.Len(t, sif, 12)
	assert.Len(t, accountNumber, 14)
}