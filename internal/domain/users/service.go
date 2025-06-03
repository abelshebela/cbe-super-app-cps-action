package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UserService struct {
	repository UserRepository
	logger     utils.Logger
}

func NewUserService(repository UserRepository, logger utils.Logger) *UserService {
	return &UserService{
		repository: repository,
		logger:     logger,
	}
}

func (s *UserService) ActiveLinkedAccounts(ctx context.Context, id string) (*LinkedAccountResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("Invalid ObjectID: %v", err)
		return nil, NewServiceError(common.DefineError.General["INVALID_ID"])
	}

	user, err := s.repository.FindByID(ctx, objectID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warnf("User with ID %s not found", id)
			return nil, NewServiceError(common.DefineError.General["NOT_FOUND"])
		}
		s.logger.Errorf("Failed to fetch user with ID %s: %v", id, err)
		return nil, NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	if user == nil {
		s.logger.Warnf("User with ID %s not found", id)
		return nil, NewServiceError(common.DefineError.General["NOT_FOUND"])
	}

	if user.IsDeleted {
		s.logger.Infof("User with ID %s is deleted", id)
		return &LinkedAccountResponse{
			UserID:         user.ID,
			FullName:       user.FullName,
			LinkedAccounts: []LinkedAccountDetail{},
		}, nil
	}

	linkedAccounts, err := s.repository.FindActiveLinkedAccounts(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch linked accounts for user ID %s: %v", id, err)
		return nil, NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	response := &LinkedAccountResponse{
		UserID:         user.ID,
		FullName:       user.FullName,
		LinkedAccounts: linkedAccounts,
	}

	s.logger.Infof("Successfully fetched active linked accounts for user ID %s", id)
	return response, nil
}