package bank

import (
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BankAdapter struct {
	bankHandler service.BankService
	logger      utils.Logger
}

func InitBankAdapter(bankHandler bank.BankHandlerService, logger utils.Logger) inboundBank.BankAdapter {
	return &BankAdapter{
		bankHandler: bankHandler,
		logger:      logger,
	}
}

func UserContextToModel(userContext ctx_util.UserContext) model.User {
	return model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}
}

func createCPSUserForReject(r *http.Request) (*model.RejectCPSAction, error) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}

	v := &model.RejectCPSAction{
		CheckerUser: UserContextToModel(userContext),
	}

	v.Department = userContext.Department
	return v, nil
}

func createCPSUserForCreate(r *http.Request) (*model.CreateCPSAction, error) {
	userContext := ctx_util.ExtractUserContext(r)

	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}
	return &model.CreateCPSAction{
		MakerUser:  UserContextToModel(userContext),
		Department: userContext.Department,
	}, nil
}

func (b *BankAdapter) CreateOneBank(w http.ResponseWriter, r *http.Request) {
	var bankRequest bank_dto.CreateBankRequest

	file, fileHeader, err := ParseMultipartFormFile(r, "logo", 10<<20)
	if err != nil {
		b.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
		return
	}
	defer file.Close()

	bankRequest.Name = r.FormValue("name")
	bankRequest.Code = r.FormValue("code")
	bankRequest.BIC = r.FormValue("bic")
	bankRequest.Logo = fileHeader

	cpsRequest, err := createCPSUserForCreate(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	cpsRequest.ActionData = bankRequest
	_, err = b.bankHandler.CreateOneBank(r.Context(), *cpsRequest)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := map[string]interface{}{}
	localization.SendSuccessResponse(w,localization.SuccessBankCreatedRequestSent,data)
	// common_util.WriteSuccessResponse(w, data, "Bank created requrest sent successfully")
}

func (b *BankAdapter) UpdateOneBank(w http.ResponseWriter, r *http.Request) {

	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var updateRequest dto.UpdateBankRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		b.logger.Errorf("invalid input", err)
		common_util.SendErrorResponse(w, common_util.InvalidReq, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	updateRequest.ID = id
	cpsReq.ActionData = updateRequest

	cpsAction, err := b.bankHandler.UpdateOneBank(r.Context(), id, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "Bank updated request sent successfully")
}

func (b *BankAdapter) DeleteOneBank(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	cpsReq.ActionData = entity.Bank{ID: id}

	cpsAction, err := b.bankHandler.DeleteOneBank(r.Context(), id, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "Bank deleted request sent successfully.")
}

func (b *BankAdapter) GetAllBank(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	ctx := r.Context()
	banks, err := b.bankHandler.GetAllBank(ctx, filterParams)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, banks, "Banks retrieved request sent successfully")
}

func (b *BankAdapter) GetOneBank(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	bank, err := b.bankHandler.GetOneBank(r.Context(), id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, bank, "Bank retrieved successfully")
}

func (b *BankAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	actionCode, ok := common_util.GetParam(r, "action_code")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'action_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForReject(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	var rejectPayload model.RejectCPSAction
	if err := json.NewDecoder(r.Body).Decode(&rejectPayload); err != nil {
		b.logger.Errorf("failed to decode rejection payload: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidReq, 0, nil)
		return
	}

	cpsReq.ActionCode = actionCode
	cpsReq.RejectedReason = rejectPayload.RejectedReason

	rejectAction, err := b.bankHandler.Reject(r.Context(), *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, rejectAction, "Bank rejected successfully")
}

func (b *BankAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	cpsReq.ActionData = dto.UpdateBankRequest{ID: id}

	cpsAction, err := b.bankHandler.EnableOrDisableBank(r.Context(), id, model.RequestDisableBank, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "Bank disable request sent successfully")
}

func (b *BankAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	cpsReq.ActionData = dto.UpdateBankRequest{ID: id}

	_, err = b.bankHandler.EnableOrDisableBank(r.Context(), id, model.RequestEnableBank, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, map[string]interface{}{}, "Bank enable request sent successfully")
}

func (b *BankAdapter) UpdateLogo(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	if id == "" {
		common_util.SendErrorResponse(w, "BANK_ID_REQUIRED", 0, nil)
		return
	}
	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "logo", 10<<20)
	if err != nil {
		b.logger.Errorf("error parsing file: %v", err)
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	defer file.Close()

	var updateLogo dto.UpdateLogo
	updateLogo.Logo = fileHeader
	updateLogo.ID = id

	cpsRequest, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	cpsRequest.ActionData = updateLogo
	_, err = b.bankHandler.UpdateLogo(r.Context(), id, *cpsRequest)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, map[string]interface{}{}, "Bank logo updated request sent successfully.")
}
