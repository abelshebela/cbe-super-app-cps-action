package wallet

import (
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/wallet"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
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

func toModelUser(userContext ctx_util.UserContext) model.User {
	return model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
}

func createCPSUserForCreate(r *http.Request, actionData any) (*model.CreateCPSAction, error) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}
	return &model.CreateCPSAction{
		MakerUser:  toModelUser(userContext),
		Department: userContext.Department,
		ActionData: actionData,
	}, nil
}

func (wa *WalletAdapter) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var walletRequest dto.CreateWalletRequest

	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "avatar", 10<<20)
	if err != nil {
		wa.logger.Errorf("error parsing file: %v", err)
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	defer file.Close()

	walletRequest.Name = r.FormValue("name")
	walletRequest.Code = r.FormValue("code")
	walletRequest.Avatar = fileHeader

	cpsRequest, err := createCPSUserForCreate(r, walletRequest)
	if err != nil {
		wa.logger.Errorf("failed to create CPS user for create wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	ctx := r.Context()
	_, err = wa.walletHandler.CreateWallet(ctx, *cpsRequest)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Wallet created successfully")
}

func (wa *WalletAdapter) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	var updateRequest dto.UpdateWalletRequest

	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "avatar", 10<<20)
	if err != nil && err.Error() != common_util.ErrMissingFile {
		wa.logger.Errorf("error parsing file: %v", err)
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	if file != nil {
		defer file.Close()
	}

	updateRequest.Name = r.FormValue("name")
	updateRequest.Code = r.FormValue("code")
	updateRequest.Avatar = fileHeader

	if updateRequest.Name == "" && updateRequest.Code == "" && fileHeader == nil {
		common_util.SendErrorResponse(w, "NO_DATA_PROVIDED_FOR_UPDATE", http.StatusBadRequest, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r, updateRequest)
	if err != nil {
		wa.logger.Errorf("failed to create CPS user for create wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	_, err = wa.walletHandler.UpdateWallet(r.Context(), id, *cpsReq)
	if err != nil {
		wa.logger.Errorf("failed to update wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Wallet updated successfully")
}

func (wa *WalletAdapter) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r, entity.Wallet{
		ID: id,
	})
	if err != nil {
		wa.logger.Errorf("failed to create CPS user for create wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	ctx := r.Context()
	_, err = wa.walletHandler.DeleteWallet(ctx, id, *cpsReq)
	if err != nil {
		wa.logger.Errorf("failed to delete wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Wallet deleted successfully")
}

func (wa *WalletAdapter) GetAllWallet(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	fmt.Print(filterParams.Page, "page")
	fmt.Print(filterParams.PerPage, "per page")

	ctx := r.Context()
	wallets, err := wa.walletHandler.GetAllWallet(ctx, filterParams)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to get all wallets", err)
		return
	}

	doc, _ := common_util.StructToMap(wallets)
	common_util.BaseResponseMaker(doc, w, "Wallets retrieved successfully", 200)
}

func (wa *WalletAdapter) GetWallet(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	ctx := r.Context()

	wallet, err := wa.walletHandler.GetWallet(ctx, id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to delete wallet", err)
		return
	}

	common_util.WriteSuccessResponse(w, wallet, "Wallet retrieved successfully")
}

func (wa *WalletAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r, entity.Wallet{
		ID: id,
	})
	if err != nil {
		wa.logger.Errorf("failed to create CPS user for create wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	ctx := r.Context()
	_, err = wa.walletHandler.EnableOrDisableWallet(ctx, id, model.RequestDisableWallet, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to delete wallet", err)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Wallet disabled successfully")
}

func (wa *WalletAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r, entity.Wallet{
		ID: id,
	})
	if err != nil {
		wa.logger.Errorf("failed to create CPS user for create wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	ctx := r.Context()
	_, err = wa.walletHandler.EnableOrDisableWallet(ctx, id, model.RequestEnableWallet, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to delete wallet", err)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Wallet enabled successfully")
}
