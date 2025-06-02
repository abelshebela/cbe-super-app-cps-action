package users

import (
    "context"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
    "cbe-super-app-member-users/internal/port/outbound"
)

type UserService struct {
    repository outbound.Repository[entities.User]
    logger     utils.Logger
}

func NewUserService(repository outbound.Repository[entities.User], logger utils.Logger) *UserService {
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
        if err.Error() == "mongo: no documents in result" {
            s.logger.Warnf("User with ID %s not found", id)
            return nil, NewServiceError(common.DefineError.General["NOT_FOUND"])
        }
        s.logger.Errorf("Failed to fetch user with ID %s: %v", id, err)
        return nil, NewServiceError(common.DefineError.General["INTERNAL_SERVER_ERROR"])
    }

    if user == nil {
        s.logger.Warnf("User with ID %s not found", id)
        return nil, NewServiceError(common.DefineError.General["NOT_FOUND"])
    }

    activeAccounts := make([]LinkedAccountDetail, 0, len(user.LinkedAccount))
    for _, account := range user.LinkedAccount {
        if account.IsAccountActive && account.LinkedStatus && !user.IsDeleted {
            activeAccounts = append(activeAccounts, LinkedAccountDetail{
                AccountNumber:     account.AccountNumber,
                AccountBranchCode: account.AccountBranchCode,
                LinkedBranch:      account.LinkedBranch,
                IsAccountActive:   account.IsAccountActive,
                Maker:             account.MakerAndChecker.Linkers.Maker,
                Checker:           account.MakerAndChecker.Linkers.Checker,
            })
        }
    }

    response := &LinkedAccountResponse{
        UserID:         user.ID.Hex(),
        FullName:       user.FullName,
        LinkedAccounts: activeAccounts,
    }

    s.logger.Infof("Successfully fetched active linked accounts for user ID %s", id)
    return response, nil
}