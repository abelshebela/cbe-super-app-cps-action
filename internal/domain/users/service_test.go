package users_test

import (
	"context"
	"testing"

	"cbe-super-app-member-users/internal/domain/users"
	"cbe-super-app-member-users/internal/domain/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestActiveLinkedAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger := utils.NewLogger()
	service := users.NewUserService(mockRepo, logger)

	userID := primitive.NewObjectID().Hex()
	user := &users.User{
		ID:        userID,
		FullName:  users.FullName{FirstName: "John", MiddleName: "M", LastName: "Doe"},
		IsDeleted: false,
	}
	linkedAccounts := []users.LinkedAccountDetail{
		{
			AccountNumber:     "1234567890",
			AccountBranchCode: "001",
			LinkedBranch:      "Main Branch",
			IsAccountActive:   true,
			LinkedStatus:      true,
			CurrencyCode:      "USD",
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRepo.EXPECT().FindActiveLinkedAccounts(gomock.Any(), userID).Return(linkedAccounts, nil)

		response, err := service.ActiveLinkedAccounts(context.Background(), userID)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, userID, response.UserID)
		assert.Equal(t, user.FullName, response.FullName)
		assert.Len(t, response.LinkedAccounts, 1)
		assert.Equal(t, "1234567890", response.LinkedAccounts[0].AccountNumber)
		assert.Equal(t, "USD", response.LinkedAccounts[0].CurrencyCode)
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := service.ActiveLinkedAccounts(context.Background(), "invalid-id")
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["INVALID_ID"].Code, err.(*users.ServiceError).Code)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, mongo.ErrNoDocuments)

		_, err := service.ActiveLinkedAccounts(context.Background(), userID)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["NOT_FOUND"].Code, err.(*users.ServiceError).Code)
	})

	t.Run("no linked accounts", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRepo.EXPECT().FindActiveLinkedAccounts(gomock.Any(), userID).Return([]users.LinkedAccountDetail{}, nil)

		response, err := service.ActiveLinkedAccounts(context.Background(), userID)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.LinkedAccounts, 0)
	})

	t.Run("deleted user", func(t *testing.T) {
		deletedUser := &users.User{
			ID:        userID,
			FullName:  user.FullName,
			IsDeleted: true,
		}
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(deletedUser, nil)

		response, err := service.ActiveLinkedAccounts(context.Background(), userID)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.LinkedAccounts, 0)
	})
}