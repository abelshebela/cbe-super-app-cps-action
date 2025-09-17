package service

import (
	"net/http"

	portalcardApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/portal_card"
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
	filterParams := common_util.ExtractFilterParams(r)
	ctx := r.Context()

	cards, err := s.appService.GetAll(ctx, filterParams)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, _ := common_util.StructToMap(cards)
	common_util.BaseResponseMaker(data, w, "Portal cards fetched successfully", 200)
}
