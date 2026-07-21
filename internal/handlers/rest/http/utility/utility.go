package utility

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	utilityInbound "cbe-super-app-cps-action/internal/constants/interfaces/utility"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"context"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type utilityAdapter struct {
	svc    service.UtilityService
	logger utils.Logger
}

func NewUtilityAdapter(svc service.UtilityService, logger utils.Logger) utilityInbound.UtilityInbound {
	return &utilityAdapter{svc: svc, logger: logger}
}

// GetByUniqueToken godoc
//
//	@Summary		Get utility action by unique token
//	@Description	Returns the opaque payload stored in the CPS action identified by unique_token.
//	@Tags			Utility
//	@Produce		json
//	@Param			unique_token	path		string									true	"Unique token"
//	@Success		200				{object}	localization.StandardResponse{data=map[string]interface{}}
//	@Failure		404				{object}	localization.StandardResponse{data=nil}
//	@Failure		500				{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/utility/{unique_token} [get]
func (a *utilityAdapter) GetByUniqueToken(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "GetByUniqueToken", "handler", "utility")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	uniqueToken := chi.URLParam(r, "unique_token")
	if uniqueToken == "" {
		log.Errorf("[UtilityH][Get] missing unique_token")
		localization.SendBadRequestResponse(w, "unique_token is required")
		return
	}

	payload, err := a.svc.GetByUniqueToken(ctx, uniqueToken, "")
	if err != nil {
		span.RecordError(err)
		log.Errorf("[UtilityH][Get] error fetching token=%s: %v", uniqueToken, err)
		w = local_util.HandlePendingResponseError(ctx, w, err)
		localization.SendInternalServerErrorResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessGetOneUtility, payload)
}
