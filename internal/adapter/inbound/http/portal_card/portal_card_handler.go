package service

import (
	"net/http"

	portalcardApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/portal_card"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type portalCardHandler struct {
	appService portalcardApp.PortalCardApplication
	logger     utils.Logger
}

func NewportalCardHandler(domain portalcardApp.PortalCardApplication, logger utils.Logger) inbound.PortalCardBound {
	return &portalCardHandler{
		appService: domain,
		logger:     logger,
	}
}

func (s *portalCardHandler) GetAllPortalCard(w http.ResponseWriter, r *http.Request) {

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}
	if r == nil {
		common_util.SendErrorResponse(w, common_util.InvalidInput, 0, nil)
		return
	}

	data, err := s.appService.GetAll(r.Context())
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, data, "Successfuly Fetched")
}
