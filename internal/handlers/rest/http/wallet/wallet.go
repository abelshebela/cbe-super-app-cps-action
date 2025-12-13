package wallet

import (
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	walletInbound "cbe-super-app-cps-action/internal/constants/interfaces/wallet"
	"cbe-super-app-cps-action/internal/constants/types"
	"errors"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/constants/localization"
	walletcore "cbe-super-app-cps-action/internal/handlers/rest/http/wallet/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type PaginatedWalletResponse types.PaginatedResponse[[]*model.Wallet]

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

// CreateWallet godoc
//
//	@Summary		Create a new wallet
//	@Description	Create a new wallet with the provided information
//	@Tags			Wallet
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name	formData	string									false	"name"
//	@Param			code	formData	string									false	"code"
//	@Param			self	formData	bool									false	"self"
//	@Param			other	formData	bool									false	"other"
//	@Param			agent	formData	bool									false	"agent"
//	@Param			avatar	formData	file									fale	"Avatar image file"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Wallet creation request sent successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/wallets [post]
func (a *walletAdapter) CreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "createWallet", "handler", "wallet")
	defer span.End()

	var req walletDto.WalletRequest
	req, err := walletcore.ParseWalletRequestFromMultipartForm(r, true)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("error fetching wallet create request data")
	}
	a.logger.Infof("this is the wallet request%+v\n", req)

	if err := req.AggregatedValidate(true); err != nil {
		span.RecordError(err)
		a.logger.Errorf("wallet request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("wallet.code", req.Code),
		attribute.String("wallet.name", req.Name),
	)

	if err := a.walletApp.CreateWallet(ctx, req); err != nil {
		span.RecordError(err)
		a.logger.Errorf("failed to create wallet: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("wallet creation request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessWalletCreationRequestSent, nil)
}

// UpdateWallet godoc
//
//	@Summary		Update a wallet
//	@Description	Update a wallet with the provided information
//	@Tags			Wallet
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string									true	"Wallet ID"
//	@Param			name	formData	string									false	"name"
//	@Param			code	formData	string									false	"code"
//	@Param			self	formData	bool									false	"self"
//	@Param			other	formData	bool									false	"other"
//	@Param			agent	formData	bool									false	"agent"
//	@Param			avatar	formData	file									fale	"Avatar image file"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Wallet update request sent successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/wallets/{id} [patch]
func (a *walletAdapter) UpdateWallet(w http.ResponseWriter, r *http.Request) {

	ctx, span := local_util.TraceLogger(r.Context(), "", "updateWallet", "handler", "wallet")
	defer span.End()
	id := chi.URLParam(r, "id")

	if id == "" {
		span.RecordError(errors.New("wallet ID is required for update"))
		a.logger.Errorf("wallet ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	req, err := walletcore.ParseWalletRequestFromMultipartForm(r, false)
	if err != nil {

		span.RecordError(err)
		a.logger.Errorf("failed to parse wallet update request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	a.logger.Infof("this is the wallet request%+v\n", req)

	// Track which fields were provided in the form
	fieldsProvided := make(map[string]bool)

	// Parse the multipart form to check which fields exist
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		a.logger.Errorf("failed to parse multipart form: %v", err)
		localization.SendBadRequestResponse(w, "Failed to parse form data")
		return
	}

	// Check for each field in the form
	if r.FormValue("name") != "" {
		fieldsProvided["name"] = true
	}
	if r.FormValue("code") != "" {
		fieldsProvided["code"] = true
	}
	if r.FormValue("self") != "" {
		fieldsProvided["self"] = true
	}
	if r.FormValue("other") != "" {
		fieldsProvided["other"] = true
	}
	if r.FormValue("agent") != "" {
		fieldsProvided["agent"] = true
	}
	if _, _, err := r.FormFile("avatar"); err == nil {
		fieldsProvided["avatar"] = true
	}

	if err := req.AggregatedValidate(false); err != nil {
		span.RecordError(err)
		a.logger.Errorf("wallet request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if req.IsEmpty() {
		span.RecordError(errors.New("no data provided for wallet update"))

		a.logger.Warnf("no data provided for wallet update, wallet ID: %s", id)
		localization.SendErrorResponse(w, localization.ErrorWalletUpdateEmptyPayload, nil, nil)
		return
	}

	// Pass fieldsProvided to the service

	span.SetAttributes(attribute.String("wallet.id", id))
	if err := a.walletApp.UpdateWallet(ctx, id, req, fieldsProvided); err != nil {
		span.RecordError(err)
		a.logger.Errorf("failed to update wallet (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("wallet update request submitted successfully, wallet ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessWalletUpdateRequestSent, nil)
}

// DeleteWallet godoc
//
//	@Summary		Delete a wallet
//	@Description	Permanently delete a wallet by ID
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Wallet ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Wallet deleted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Wallet not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/wallets/{id} [delete]
func (a *walletAdapter) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "deleteWallet", "handler", "wallet")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet ID is required for delete"))
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))

	if err := a.walletApp.DeleteWallet(ctx, id); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessWalletDeleted, nil)
}

// EnableWallet godoc
//
//	@Summary		Enable a wallet
//	@Description	Enable a wallet by ID
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Wallet ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Wallet enable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Wallet not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/wallets/{id}/enable [patch]
func (a *walletAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "enableWallet", "handler", "wallet")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet ID is required for enable"))
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))
	if err := a.walletApp.EnableOrDisableWallet(ctx, id, true); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessWalletEnableRequestSubmitted, nil)
}

// DisableWallet godoc
//
//	@Summary		Disable a wallet
//	@Description	Disable a wallet by ID
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Wallet ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Wallet disable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Wallet not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/wallets/{id}/disable [patch]
func (a *walletAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "disableWallet", "handler", "wallet")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("wallet ID is required for disable"))
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("wallet.id", id))
	if err := a.walletApp.EnableOrDisableWallet(ctx, id, false); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessWalletDisableRequestSubmitted, nil)
}

// GetWallet godoc
//
//	@Summary		Get wallet by ID
//	@Description	Retrieve a wallet's details by ID
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string												true	"Wallet ID"
//	@Success		200	{object}	localization.StandardResponse{data=model.Wallet}	"Wallet retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}				"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}				"Wallet not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/wallets/{id} [get]
func (a *walletAdapter) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getWallet", "handler", "wallet")
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

// GetWallets godoc
//
//	@Summary		List wallets
//	@Description	Retrieve wallets with pagination, filtering, and search. Filterable fields: name, code, enabled. Searchable fields: name, code.
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"Page number"		default(1)
//	@Param			per_page	query	int		false	"Items per page"	default(10)
//	@Param			name		query	string	false	"Filter by wallet name"
//	@Param			code		query	string	false	"Filter by wallet code"
//	@Param			enabled		query	bool	false	"Filter by enabled status"
//	@Param			search		query	string	false	"Search term (searches name, code)"
//	@Success		200	{object}	localization.StandardResponse{data=PaginatedWalletResponse}	"Wallets retrieved successfully"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/wallets [get]
func (a *walletAdapter) GetAllWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getAllWallets", "handler", "wallet")
	defer span.End()
	filter := local_util.ExtractFilterParams(r)
	a.logger.Infof("fetching wallets with filter: %+v", filter)

	list, err := a.walletApp.GetAllWallet(ctx, *filter)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("failed to fetch wallets: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("wallet.count", len(list.Data)))
	a.logger.Infof("wallets fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessWalletsRetrieved, list)
}
