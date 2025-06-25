package users_application

import (
	"context"
	"mime/multipart"

	domainUsers "cbe-super-app-member-users/internal/domain/users"
	userPort "cbe-super-app-member-users/internal/port/inbound/users"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	FetchLinkedAccounts(ctx context.Context, userID string) (*userPort.LinkedAccountResponse, error)
	GenerateEmailOTP(ctx context.Context, req userPort.OTPRequest) (string, error)
	VerifyEmailOTP(ctx context.Context, req userPort.OTPVerification) error
	UpdateProfilePicture(ctx context.Context, id string, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type UsersHandler struct {
	userService *domainUsers.UserService
	logger      utils.Logger
}

func InitUsersHandler(userService *domainUsers.UserService, logger utils.Logger, minIO config.MinioClientInterface) ApplicationService {
	return UsersHandler{
		userService: userService,
		logger:      logger,
	}
}
