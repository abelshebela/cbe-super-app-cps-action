package portalcard

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PortalCardPaginatedResponse types.PaginatedResponse[[]*model.Card]
type portalCardAdapter struct {
	appService service.PortalCardService
	logger     utils.Logger
}

func InitPortalCardAdapter(appService service.PortalCardService, logger utils.Logger) portal_card.PortalCardAdapter {
	return &portalCardAdapter{
		appService: appService,
		logger:     logger,
	}
}

// Get All Portal Cards
//
//	@Summary		Get All Portal Cards
//	@Description	Retrieves all portal cards with pagination
//	@Tags			PortalCard
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200			{object}	localization.StandardResponse{data=PortalCardPaginatedResponse}
//	@Failure		400,500		{object}	localization.StandardResponse{data=nil}
//	@Router			/portal_cards [get]
func (s *portalCardAdapter) GetAllPortalCard(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	cards, err := s.appService.GetAll(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessPortalCardsFetched, cards)
}
