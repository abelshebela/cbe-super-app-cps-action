package initiator

import (
	inbound "cbe-super-app-member-auth/internal/handlers/rest"
	userHandler "cbe-super-app-member-auth/internal/handlers/rest/http/users"
	service "cbe-super-app-member-auth/internal/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	UserHandler inbound.Users
}

func InitHandler(serviceLayer service.UserService, logger utils.Logger) Handler {
	return Handler{
		UserHandler: userHandler.Init(serviceLayer, logger),
	}
}
