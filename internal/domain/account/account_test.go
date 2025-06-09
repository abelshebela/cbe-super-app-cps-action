package account_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"cbe-super-app-member-users/internal/domain/account"
// 	"cbe-super-app-member-users/internal/domain/users"
// 	"cbe-super-app-member-users/internal/domain/account/mocks"
// 	"cbe-super-app-member-users/internal/shared"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// )

// func TestCreateAccount(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mocks.NewMockAccountRepository(ctrl)
// 	mockAPIClient := mocks.NewMockAccountAPIClient(ctrl)
// 	logger := utils.NewLogger()
// 	service := account.NewAccountService(mockRepo, mockAPIClient, logger)

// 	userID := "6644c8e37f41d2c9c1c293fa"
// 	phoneNumber := "+251911234567"
// 	user := &account.AccountUser{
// 		ID:              userID,
// 		PhoneNumber:     phoneNumber,
// 		KYCLevel:        1,
// 		BranchCode:      "12345678",
// 		FullName:        "John Doe",
// 		RegistrationType: "SAVINGS",
// 		AndOrStatus:     true,
// 	}
// 	sif := "SIF-12345678"
// 	accountNumber := "ACC-9876543210"

// 	t.Run("success", func(t *testing.T) {
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), userID).Return(user, nil)
// 		mockAPIClient.EXPECT().LookupAccountByPhone(gomock.Any(), phoneNumber).Return(false, nil)
// 		mockAPIClient.EXPECT().GenerateSIFAndAccountNumber().Return(sif, accountNumber, nil)
// 		mockRepo.EXPECT().UpdateUserCustomerNumber(gomock.Any(), userID, sif).Return(nil)
// 		mockRepo.EXPECT().CreateLinkedAccount(gomock.Any(), gomock.Any()).DoAndReturn(
// 			func(_ context.Context, linkedAccount *account.LinkedAccount) error {
// 				assert.Equal(t, userID, linkedAccount.UserID)
// 				assert.Equal(t, sif, linkedAccount.CustomerNumber)
// 				assert.Equal(t, accountNumber, linkedAccount.AccountNumber)
// 				assert.Equal(t, user.FullName, linkedAccount.AccountHolderName)
// 				assert.Equal(t, "SAVINGS", linkedAccount.AccountType)
// 				assert.Equal(t, user.BranchCode, linkedAccount.BranchCode)
// 				assert.Equal(t, "USD", linkedAccount.CurrencyCode)
// 				assert.True(t, linkedAccount.IsMain)
// 				return nil
// 			},
// 		)

// 		result, err := service.CreateAccount(context.Background(), userID)
// 		require.NoError(t, err)
// 		require.NotNil(t, result)
// 		assert.Equal(t, sif, result.CustomerNumber)
// 		assert.Equal(t, accountNumber, result.AccountNumber)
// 	})

// 	t.Run("invalid user ID", func(t *testing.T) {
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), "invalid-id").Return(nil, users.ErrNotFound)

// 		result, err := service.CreateAccount(context.Background(), "invalid-id")
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Equal(t, common.DefineError.General["NOT_FOUND"].Code, err.(*account.ServiceError).Code)
// 	})

// 	t.Run("KYC level zero", func(t *testing.T) {
// 		zeroKYCUser := &account.AccountUser{
// 			ID:          userID,
// 			PhoneNumber: phoneNumber,
// 			KYCLevel:    0,
// 		}
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), userID).Return(zeroKYCUser, nil)

// 		result, err := service.CreateAccount(context.Background(), userID)
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Equal(t, common.DefineError.User["USER_KYC_LEVEL_ZERO"].Code, err.(*account.ServiceError).Code)
// 	})

// 	t.Run("phone number exists", func(t *testing.T) {
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), userID).Return(user, nil)
// 		mockAPIClient.EXPECT().LookupAccountByPhone(gomock.Any(), phoneNumber).Return(true, nil)

// 		result, err := service.CreateAccount(context.Background(), userID)
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Equal(t, common.DefineError.User["USER_PHONE_EXISTS"].Code, err.(*account.ServiceError).Code)
// 	})

// 	t.Run("phone lookup failed", func(t *testing.T) {
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), userID).Return(user, nil)
// 		mockAPIClient.EXPECT().LookupAccountByPhone(gomock.Any(), phoneNumber).Return(false, shared.DefineError.Account["PHONE_LOOKUP_FAILED"])

// 		result, err := service.CreateAccount(context.Background(), userID)
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Equal(t, shared.DefineError.Account["PHONE_LOOKUP_FAILED"].Code, err.(*account.ServiceError).Code)
// 	})

// 	t.Run("API request failed", func(t *testing.T) {
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), userID).Return(user, nil)
// 		mockAPIClient.EXPECT().LookupAccountByPhone(gomock.Any(), phoneNumber).Return(false, shared.DefineError.Account["API_REQUEST_FAILED"])

// 		result, err := service.CreateAccount(context.Background(), userID)
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Equal(t, shared.DefineError.Account["API_REQUEST_FAILED"].Code, err.(*account.ServiceError).Code)
// 	})

// 	t.Run("SIF generation failed", func(t *testing.T) {
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), userID).Return(user, nil)
// 		mockAPIClient.EXPECT().LookupAccountByPhone(gomock.Any(), phoneNumber).Return(false, nil)
// 		mockAPIClient.EXPECT().GenerateSIFAndAccountNumber().Return("", "", errors.New("generation failed"))

// 		result, err := service.CreateAccount(context.Background(), userID)
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Equal(t, common.DefineError.General["GENERAL_SIF_GENERATION_FAILED"].Code, err.(*account.ServiceError).Code)
// 	})

// 	t.Run("update customer number failed", func(t *testing.T) {
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), userID).Return(user, nil)
// 		mockAPIClient.EXPECT().LookupAccountByPhone(gomock.Any(), phoneNumber).Return(false, nil)
// 		mockAPIClient.EXPECT().GenerateSIFAndAccountNumber().Return(sif, accountNumber, nil)
// 		mockRepo.EXPECT().UpdateUserCustomerNumber(gomock.Any(), userID, sif).Return(mongo.ErrClientDisconnected)

// 		result, err := service.CreateAccount(context.Background(), userID)
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Equal(t, common.DefineError.General["GENERAL_DB_UPDATE_FAILED"].Code, err.(*account.ServiceError).Code)
// 	})

// 	t.Run("create linked account failed", func(t *testing.T) {
// 		mockRepo.EXPECT().FindAccountUserByID(gomock.Any(), userID).Return(user, nil)
// 		mockAPIClient.EXPECT().LookupAccountByPhone(gomock.Any(), phoneNumber).Return(false, nil)
// 		mockAPIClient.EXPECT().GenerateSIFAndAccountNumber().Return(sif, accountNumber, nil)
// 		mockRepo.EXPECT().UpdateUserCustomerNumber(gomock.Any(), userID, sif).Return(nil)
// 		mockRepo.EXPECT().CreateLinkedAccount(gomock.Any(), gomock.Any()).Return(mongo.ErrClientDisconnected)

// 		result, err := service.CreateAccount(context.Background(), userID)
// 		require.Error(t, err)
// 		assert.Nil(t, result)
// 		assert.Equal(t, common.DefineError.General["GENERAL_DB_INSERT"].Code, err.(*account.ServiceError).Code)
// 	})
// }