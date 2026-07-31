package wallet

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	walletDto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/wallet"
	walletInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/wallet"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"context"
	"errors"
	"net/http"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	walletcore "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/wallet/core"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	local_model "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type PaginatedWalletResponse types.PaginatedResponse[[]*local_model.Wallet]

type walletAdapter struct {
	walletApp service.WalletService
	logger    utils.Logger
}

func InitWalletAdapter(walletApp service.WalletService, logger utils.Logger) walletInbound.WalletAdapter {
	return &walletAdapter{
		walletApp: walletApp,
		logger:    logger,
	}
}

func (a *walletAdapter) CreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createWallet", "handler", "wallet")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	var req walletDto.WalletRequest
	req, err := walletcore.ParseWalletRequestFromMultipartForm(r, true)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[WalletH][Create] parse form err")
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	log.Infof("[WalletH] req: %+v", req)

	if err := req.Validate(true); err != nil {
		span.RecordError(err)
		log.Errorf("[WalletH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("wallet.unique_code", req.UniqueCode),
		attribute.String("wallet.name", req.Name),
	)

	if err := a.walletApp.CreateWallet(ctx, req); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[WalletH][Create] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		a.logger.Infof("[WalletH][Create]  create wallet successfully")
		localization.SendSuccessResponse(w, localization.SuccessWalletCreated, nil)
	} else {
		a.logger.Infof("[WalletH][Create] request sent successfully for create wallet")
		localization.SendSuccessResponse(w, localization.SuccessWalletCreationRequestSent, nil)
	}
}

func (a *walletAdapter) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateWallet", "handler", "wallet")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	id := chi.URLParam(r, "id")
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	if id == "" {
		span.RecordError(errors.New("wallet ID is required for update"))
		log.Errorf("[WalletH][Update] id required")
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	req, err := walletcore.ParseWalletRequestFromMultipartForm(r, false)
	if err != nil {

		span.RecordError(err)
		log.Errorf("[WalletH][Update] parse form err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(false); err != nil {
		span.RecordError(err)
		log.Errorf("[WalletH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if req.IsEmpty() {
		span.RecordError(errors.New("no data provided for wallet update"))

		log.Warnf("[WalletH][Update] empty payload id: %s", id)
		localization.SendErrorResponse(w, localization.ErrorWalletUpdateEmptyPayload, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))
	if err := a.walletApp.UpdateWallet(ctx, id, req); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[WalletH][Update] svc err id: %s: %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		a.logger.Infof("[WalletH][Update]  update wallet successfully")
		localization.SendSuccessResponse(w, localization.SuccessWalletUpdated, nil)
	} else {
		a.logger.Infof("[WalletH][Update] request sent successfully for update wallet")
		localization.SendSuccessResponse(w, localization.SuccessWalletUpdateRequestSent, nil)
	}
}

func (a *walletAdapter) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "deleteWallet", "handler", "wallet")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet ID is required for delete"))
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))

	if err := a.walletApp.DeleteWallet(ctx, id); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessWalletDeleted, nil)
}

func (a *walletAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableWallet", "handler", "wallet")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet ID is required for enable"))
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))
	if err := a.walletApp.EnableOrDisableWallet(ctx, id, true); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		a.logger.Infof("[WalletH][Enable]  enable wallet successfully")
		localization.SendSuccessResponse(w, localization.SuccessWalletEnabled, nil)
	} else {
		a.logger.Infof("[WalletH][Enable] request sent successfully for enable wallet")
		localization.SendSuccessResponse(w, localization.SuccessWalletEnableRequestSubmitted, nil)
	}
}

func (a *walletAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableWallet", "handler", "wallet")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet ID is required for disable"))
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))
	if err := a.walletApp.EnableOrDisableWallet(ctx, id, false); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		a.logger.Infof("[WalletH][Disable]  disable wallet successfully")
		localization.SendSuccessResponse(w, localization.SuccessWalletDisabled, nil)
	} else {
		a.logger.Infof("[WalletH][Disable] request sent successfully for disable wallet")
		localization.SendSuccessResponse(w, localization.SuccessWalletDisableRequestSubmitted, nil)
	}
}

func (a *walletAdapter) EnableWalletService(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableWalletService", "handler", "wallet")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet service ID is required for enable"))
		localization.SendErrorResponse(w, localization.ErrorWalletServiceIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))
	if err := a.walletApp.EnableOrDisableWalletService(ctx, id, true); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		a.logger.Infof("[WalletH][Enable]  enable wallet service successfully")
		localization.SendSuccessResponse(w, localization.SuccessWalletServiceEnabled, nil)
	} else {
		a.logger.Infof("[WalletH][Enable] request sent successfully for enable wallet service")
		localization.SendSuccessResponse(w, localization.SuccessWalletServiceEnableRequestSubmitted, nil)
	}
}

func (a *walletAdapter) DisableWalletService(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableWallet", "handler", "wallet")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet service ID is required for disable"))
		localization.SendErrorResponse(w, localization.ErrorWalletServiceIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.service.id", id))
	if err := a.walletApp.EnableOrDisableWalletService(ctx, id, false); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		a.logger.Infof("[WalletH][Disable]  disable wallet service successfully")
		localization.SendSuccessResponse(w, localization.SuccessWalletServiceDisabled, nil)
	} else {
		a.logger.Infof("[WalletH][Disable] request sent successfully for disable wallet service")
		localization.SendSuccessResponse(w, localization.SuccessWalletServiceDisableRequestSubmitted, nil)
	}
}

func (a *walletAdapter) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getWallet", "handler", "wallet")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet ID is required for get"))
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))

	wallet, err := a.walletApp.GetWallet(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessWalletRetrieved, wallet)
}

func (a *walletAdapter) GetAllWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllWallets", "handler", "wallet")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	filterParams, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

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

	log.Infof("[WalletH][GetAll] filter: %+v", filter)

	list, err := a.walletApp.GetAllWallet(ctx, *filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[WalletH][GetAll] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("wallet.count", len(list.Data)))
	log.Infof("[WalletH][GetAll] ok")
	localization.SendSuccessResponse(w, localization.SuccessWalletsRetrieved, list)
}
