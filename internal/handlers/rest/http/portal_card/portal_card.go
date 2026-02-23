package portalcard

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "portalCard", "portalCardAdapter", "GetAllPortalCard")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		span.AddEvent("Invalid pagination params", trace.WithAttributes(attribute.Int("page", filterParams.Page), attribute.Int("per_page", filterParams.PerPage)))
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	cards, err := s.appService.GetAll(ctx, filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[GetAllPortalCard] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Portal cards retrieved", trace.WithAttributes(attribute.Int("count", len(cards.Data))))
	log.Infof("[GetAllPortalCard] retrieved %d portal cards", len(cards.Data))
	localization.SendSuccessResponse(w, localization.SuccessPortalCardsFetched, cards)
}
