package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/rest"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/user"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"
)

type Handler struct {
	UserHandler rest.Users
}

func InitHandler(serviceLayer user.UsersService, logger logger.Logger) Handler{
	return Handler{
		UserHandler: user.NewUserService()
	}
}
