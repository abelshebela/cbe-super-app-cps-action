package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/wallet"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	inboundWallet "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/wallet"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

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

func createCPSUserForAuthorize(r *http.Request, actionCode string) (*model.AuthorizeCPSAction, error) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}
	return &model.AuthorizeCPSAction{
		CheckerUser: toModelUser(userContext),
		Department:  userContext.Department,
		ActionCode:  actionCode,
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
	cpsRes, err := wa.walletHandler.CreateWallet(ctx, *cpsRequest)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsRes, "Wallet created successfully")
}

func (wa *WalletAdapter) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	var updateRequest dto.UpdateWalletRequest

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		wa.logger.Errorf("invalid input", err)
		common_util.SendErrorResponse(w, common_util.InvalidReq, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r, updateRequest)
	if err != nil {
		wa.logger.Errorf("failed to create CPS user for create wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	ctx := r.Context()
	cpsAction, err := wa.walletHandler.UpdateWallet(ctx, id, *cpsReq)
	if err != nil {
		wa.logger.Errorf("failed to update wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "Wallet updated successfully")
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
	cpsAction, err := wa.walletHandler.DeleteWallet(ctx, id, *cpsReq)
	if err != nil {
		wa.logger.Errorf("failed to delete wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "Wallet deleted successfully")
}

func (wa *WalletAdapter) GetAllWallet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}

	perPage := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt <= 10 && perPageInt > 0 {
		perPage = perPageInt
	}

	search := query.Get("search")
	filter := query.Get("filter")

	filterParams := &constant.Filter{
		Page:    page,
		PerPage: perPage,
		Search:  search,
		Filters: filter,
	}

	ctx := r.Context()

	wallets, err := wa.walletHandler.GetAllWallet(ctx, filterParams)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to get all wallets", err)
		return
	}

	common_util.WriteSuccessResponse(w, wallets, "Wallets retrieved successfully")
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

func (wa *WalletAdapter) Authorize(w http.ResponseWriter, r *http.Request) {
	actionCode, ok := common_util.GetParam(r, "action_code")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'action_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForAuthorize(r, actionCode)
	if err != nil {
		wa.logger.Errorf("failed to create CPS user for create wallet", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	ctx := r.Context()
	authAction, err := wa.walletHandler.Authorize(ctx, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to delete wallet", err)
		return
	}

	common_util.WriteSuccessResponse(w, authAction, "Wallet authorized successfully")
}

func (wa *WalletAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	actionCode, ok := common_util.GetParam(r, "action_code")
	if !ok {
		wa.logger.Errorf("missing or invalid parameter 'action_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var cpsReq model.RejectCPSAction

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		wa.logger.Errorf("failed to decode wallet request %v", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)

	cpsReq.CheckerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
	cpsReq.Department = userContext.Department
	cpsReq.ActionCode = actionCode

	ctx := r.Context()
	rejectAction, err := wa.walletHandler.Reject(ctx, cpsReq)
	if err != nil {
		wa.logger.Errorf("failed to reject wallet %v", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, rejectAction, "Wallet rejected successfully")
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
	cpsAction, err := wa.walletHandler.EnableOrDisableWallet(ctx, id, model.RequestDisableWallet, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to delete wallet", err)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "Wallet disabled successfully")
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
	cpsAction, err := wa.walletHandler.EnableOrDisableWallet(ctx, id, model.RequestEnableWallet, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		wa.logger.Errorf("failed to delete wallet", err)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "Wallet enabled successfully")
}
