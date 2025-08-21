package wallet

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/wallet"
	cps_entitites "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	inboundWallet "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/wallet"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type WalletAdapter struct {
	walletHandler wallet.WalletHandlerAppllication
	logger        utils.Logger
}

func InitWalletRouter(walletHandler wallet.WalletHandlerAppllication, logger utils.Logger) inboundWallet.WalletAdapter {
	return &WalletAdapter{
		walletHandler: walletHandler,
		logger:        logger,
	}
}

func createUser(r *http.Request) (*cps_entitites.User, error) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}
	return &cps_entitites.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}, nil
}

func (wa *WalletAdapter) parseRequest(r *http.Request, isCreate bool) (*WalletRequest, *cps_entitites.User, error) {
	var walletRequest WalletRequest

	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "avatar", 10<<20)
	if err != nil && (err.Error() != common_util.ErrMissingFile && !isCreate) {
		wa.logger.Errorf("error parsing file: %v", err)
		return nil, nil, err
	}
	if file != nil {
		defer file.Close()
	}

	walletRequest.Name = strings.TrimSpace(r.FormValue("name"))
	walletRequest.Code = strings.TrimSpace(r.FormValue("code"))
	walletRequest.Avatar = fileHeader

	maker, err := createUser(r)
	if err != nil {
		wa.logger.Errorf("failed to create user for create wallet", err)
		return nil, nil, err
	}

	return &walletRequest, maker, nil
}

func (wa *WalletAdapter) parseIDMaker(r *http.Request) (string, *cps_entitites.User, error) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		return "", nil, fmt.Errorf(common_util.InvalidInputParameters)
	}

	maker, err := createUser(r)
	if err != nil {
		wa.logger.Errorf("failed to create user for create wallet", err)
		return "", nil, err
	}

	return id, maker, nil
}

func (wa *WalletAdapter) CreateWallet(w http.ResponseWriter, r *http.Request) {
	walletRequest, maker, err := wa.parseRequest(r, true)
	if err != nil {
		wa.logger.Errorf("failed to parse request for create wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	if err := walletRequest.Validate(true); err != nil {
		wa.logger.Errorf("error validating wallet request: %v", err)
		common_util.SendErrorResponse(w, err, 0, nil)
		return
	}

	wallReq := ToDomainWalletRequest(walletRequest)
	err = wa.walletHandler.CreateWallet(r.Context(), *wallReq, *maker)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.BaseResponseMaker(nil, w, "Wallet creation request sent successfully", 201)
}

func (wa *WalletAdapter) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	walletRequest, maker, err := wa.parseRequest(r, true)

	if err != nil {
		wa.logger.Errorf("failed to parse request for update wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	if err := walletRequest.Validate(false); err != nil {
		wa.logger.Errorf("error validating wallet request: %v", err)
		common_util.SendErrorResponse(w, err, 0, nil)
		return
	}

	if IsEmpty(*walletRequest) {
		wa.logger.Errorf("empty wallet request")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	wallReq := ToDomainWalletRequest(walletRequest)
	err = wa.walletHandler.UpdateWallet(r.Context(), id, *wallReq, *maker)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.BaseResponseMaker(nil, w, "Wallet update request sent successfully", 200)
}

func (wa *WalletAdapter) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	id, maker, err := wa.parseIDMaker(r)
	if err != nil {
		wa.logger.Errorf("failed to parse request for delete wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	err = wa.walletHandler.DeleteWallet(r.Context(), id, *maker)
	if err != nil {
		wa.logger.Errorf("failed to delete wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.BaseResponseMaker(nil, w, "Wallet deleted successfully", 204)
}

func (wa *WalletAdapter) GetAllWallet(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	wallets, err := wa.walletHandler.GetAllWallet(r.Context(), filterParams)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to get all wallets", err)
		return
	}

	common_util.BaseResponseMaker(wallets, w, "Wallets retrieved successfully", 200)
}

func (wa *WalletAdapter) GetWallet(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	wallet, err := wa.walletHandler.GetWallet(r.Context(), id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to delete wallet", err)
		return
	}

	common_util.BaseResponseMaker(wallet, w, "Wallet retrieved successfully", 200)
}

func (wa *WalletAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	wa.handleEnableOrDisableWallet(w, r, false)
}

func (wa *WalletAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	wa.handleEnableOrDisableWallet(w, r, true)
}

func (wa *WalletAdapter) handleEnableOrDisableWallet(w http.ResponseWriter, r *http.Request, isEnable bool) {
	id, maker, err := wa.parseIDMaker(r)
	if err != nil {
		wa.logger.Errorf("failed to parse request for delete wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	err = wa.walletHandler.EnableOrDisableWallet(r.Context(), id, isEnable, *maker)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to delete wallet", err)
		return
	}

	action := "disabled"
	if isEnable {
		action = "enabled"
	}

	common_util.BaseResponseMaker(nil, w, fmt.Sprintf("Wallet %s request sent successfully", action), 200)
}
