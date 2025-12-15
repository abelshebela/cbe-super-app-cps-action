package faydaaccount

import (
	faydaInbound "cbe-super-app-cps-action/internal/constants/interfaces/fayda"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/fayda/core"
	"net/http"

	"cbe-super-app-cps-action/internal/service"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type faydaAccountHandler struct {
	faydaService service.FaydaAccountService
	logger       utils.Logger
}

func InitFaydaHandler(faydaService service.FaydaAccountService, logger utils.Logger) faydaInbound.FaydaAccount {
	return &faydaAccountHandler{
		faydaService: faydaService,
		logger:       logger,
	}
}

// InitiateEnableFaydaAccount godoc
//
//	@Summary		Enable Fayda account
//	@Description	Initiates the process to enable a Fayda account for the given user
//	@Tags			FaydaAccount
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string						true	"User Code"
//	@Success		200			{object}	localization.ResponseCode	"Fayda account enabled successfully"
//	@Failure		400			{object}	localization.ResponseCode	"Invalid request"
//	@Failure		500			{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/fayda_account/enable/{user_code} [post]
func (f *faydaAccountHandler) InitiateEnableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "initiateEnableFaydaAccount", "handler", "fayda")
	defer span.End()

	user_code, ok := core.ExtractID(w, r, f.logger)
	if !ok {
		return
	}

	span.SetAttributes(attribute.String("fayda.user_code", user_code))

	err := f.faydaService.EnableOrDisableFayda(ctx, user_code, true)
	if err != nil {
		span.RecordError(err)
		f.logger.Errorf("[InitiateEnableFaydaAccount] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	f.logger.Infof("[InitiateEnableFaydaAccount] request sent successfully for user_code: %s", user_code)
	localization.SendSuccessResponse(w, localization.SuccessFaydaEnableActionCreated, nil)
}

// InitiateDisableFaydaAccount godoc
//
//	@Summary		Disable Fayda account
//	@Description	Initiates the process to disable a Fayda account for the given user
//	@Tags			FaydaAccount
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string						true	"User Code"
//	@Success		200			{object}	localization.ResponseCode	"Fayda account disabled successfully"
//	@Failure		400			{object}	localization.ResponseCode	"Invalid request"
//	@Failure		500			{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/fayda_account/disable/{user_code} [post]
func (f *faydaAccountHandler) InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "initiateDisableFaydaAccount", "handler", "fayda")
	defer span.End()

	user_code, ok := core.ExtractID(w, r, f.logger)
	if !ok {
		return
	}

	span.SetAttributes(attribute.String("fayda.user_code", user_code))

	err := f.faydaService.EnableOrDisableFayda(ctx, user_code, false)
	if err != nil {
		span.RecordError(err)
		f.logger.Errorf("[InitiateDisableFaydaAccount] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	f.logger.Infof("[InitiateDisableFaydaAccount] request sent successfully for user_code: %s", user_code)
	localization.SendSuccessResponse(w, localization.SuccessFaydaDisableActionCreated, nil)
}
