package wallet

import (
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	walletInbound "cbe-super-app-cps-action/internal/constants/interfaces/wallet"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/localization"
	walletcore "cbe-super-app-cps-action/internal/handlers/rest/http/wallet/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

	var req walletDto.WalletRequest
	req, err := walletcore.ParseWalletRequestFromMultipartForm(r, true)
	if err != nil {
		a.logger.Errorf("error fetching wallet create request data")
	}
	a.logger.Infof("this is the wallet request%+v\n", req)

	if err := req.AggregatedValidate(true); err != nil {
		a.logger.Errorf("wallet request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := a.walletApp.CreateWallet(r.Context(), req); err != nil {
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
	id := chi.URLParam(r, "id")

	if id == "" {
		a.logger.Errorf("wallet ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	req, err := walletcore.ParseWalletRequestFromMultipartForm(r, false)
	if err != nil {
		a.logger.Errorf("failed to parse wallet update request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	a.logger.Infof("this is the wallet request%+v\n", req)

	if err := req.AggregatedValidate(false); err != nil {
		a.logger.Errorf("wallet request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if req.IsEmpty() {
		a.logger.Warnf("no data provided for wallet update, wallet ID: %s", id)
		localization.SendErrorResponse(w, localization.ErrorWalletUpdateEmptyPayload, nil, nil)
		return
	}

	if err := a.walletApp.UpdateWallet(r.Context(), id, req); err != nil {
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
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	if err := a.walletApp.DeleteWallet(r.Context(), id); err != nil {
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
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	if err := a.walletApp.EnableOrDisableWallet(r.Context(), id, true); err != nil {
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
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	if err := a.walletApp.EnableOrDisableWallet(r.Context(), id, false); err != nil {
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
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorWalletIDRequired, nil, nil)
		return
	}

	wallet, err := a.walletApp.GetWallet(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessWalletRetrieved, wallet)
}

// GetWallets godoc
//
//	@Summary		List wallets
//	@Description	Retrieve wallets with pagination and optional search
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int															false	"Page number"		default(1)
//	@Param			per_page	query		int															false	"Items per page"	default(10)
//	@Param			search		query		string														false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=PaginatedWalletResponse}	"Wallets retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}						"Internal server error"
//	@Security		BearerAuth
//	@Router			/wallets [get]
func (a *walletAdapter) GetAllWallet(w http.ResponseWriter, r *http.Request) {
	filter := local_util.ExtractFilterParams(r)
	a.logger.Infof("fetching wallets with filter: %+v", filter)

	list, err := a.walletApp.GetAllWallet(r.Context(), *filter)
	if err != nil {
		a.logger.Errorf("failed to fetch wallets: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("wallets fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessWalletsRetrieved, list)
}
