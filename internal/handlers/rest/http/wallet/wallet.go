package wallet

import (
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	walletInbound "cbe-super-app-cps-action/internal/constants/interfaces/wallet"
	"fmt"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/localization"
	walletcore "cbe-super-app-cps-action/internal/handlers/rest/http/wallet/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

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

	var req walletDto.WalletRequest

	file, fileHeader, err := walletcore.ParseMultipartFormFile(r, "avatar", 10<<20)
	if err != nil {
		a.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorWalletImageMissingOrInvalid, nil, nil)
		return
	}
	defer file.Close()
	req.Name = r.FormValue("name")
	req.Code = r.FormValue("code")
	req.Avatar = fileHeader

	fmt.Println("wallet req", req)
	if err := req.Validate(true); err != nil {
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
