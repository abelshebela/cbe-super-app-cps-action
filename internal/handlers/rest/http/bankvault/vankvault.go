package bankvault

import (
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/dto/bankvault"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/handlers/rest/http/bankvault/core"
	"cbe-super-app-cps-action/internal/service"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type handler struct {
	service service.BankVaultService
	logger  utils.Logger
}

func InitBankVaultHandler(svc service.BankVaultService, logger utils.Logger) *handler {
	return &handler{service: svc, logger: logger}
}

// CreateBankVault
//
//	@Summary		Create Bank Vault
//	@Description	Create a new bank vault with the provided information
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			request	body		bankvault.CreateBankVaultProductRequest	true	"Bank vault request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Bank vault creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/create [post]
func (h *handler) CreateBankVault(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "createBankVault", "handler", "bankVault")
	defer span.End()
	var req bankvault.CreateBankVaultProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateBankVault] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateBankVault] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	lockPeriod, err := common_utils.ParseToYears(req.LockPeriodDays)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateBankVault] lock period parse: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	product := core.ToDomainCreateBankVaultRequest(req, lockPeriod)
	id, err := h.service.CreateBankVault(ctx, product)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("bank_vault.id", id))
	h.logger.Infof("[CreateBankVault] request sent successfully with id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultCreationRequestSubmitted, nil)

}

// FindAllBankVaults
//
//	@Summary		Find All Bank Vaults
//	@Description	Find all bank vaults with the provided filters
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string									false	"saving"
//	@Param			is_deleted	query	bool									false	"true or false"
//	@Param			is_active	query		bool									false	"true or false"
//	@Param			currency	query		string									false	"ETB ,USD"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Bank vaults retrieved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/products [get]
func (h *handler) FindAllBankVaults(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "findAllBankVaults", "handler", "bankVault")
	defer span.End()
	params := common_utils.ExtractFilterParams(r)
	result, err := h.service.FindAllBankVaults(ctx, params)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[FindAllBankVaults] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("bank_vault.count", len(result.Data)))
	h.logger.Infof("[FindAllBankVaults] retrieved %d bank vaults", len(result.Data))
	localization.SendSuccessResponse(w, localization.SuccessBankVaultsRetrieved, result)
}

// GetBankVault
//
//	@Summary		Get Bank Vault
//	@Description	Get a bank vault by the provided ID
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Bank vault ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Bank vault retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/products/{id} [get]
func (h *handler) GetBankVault(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getBankVault", "handler", "bankVault")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).WithService("bankvault_product").
			WithOperation("GetBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("bank_vault.id", id))

	result, err := h.service.GetBankVault(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[GetBankVault] bank vault retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultRetrieved, result)
}

// UpdateBankVault
//
//	@Summary		Update Bank Vault
//	@Description	Update a bank vault by the provided ID
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Bank vault ID"
//	@Param			request	body		bankvault.UpdateBankVaultProductRequest	true	"Bank vault request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Bank vault update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/products/update/{id} [patch]
func (h *handler) UpdateBankVault(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "updateBankVault", "handler", "bankVault")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("UpdateBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	var req bankvault.UpdateBankVaultProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateBankVault] decode: %v", err)
		appErr := middleware.NewValidationError("Invalid JSON input", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("UpdateBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateBankVault] validation: %v", err)
		appErr := middleware.NewValidationError(err.Error(), map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("UpdateBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	product := core.ToDomainUpdateBankVaultRequest(req, nil)

	span.SetAttributes(attribute.String("bank_vault.id", id))

	_, err = h.service.UpdateBankVault(ctx, id, product)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[UpdateBankVault] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultUpdateRequestSubmitted, nil)
}

// DeleteBankVault
//
//	@Summary		Delete Bank Vault
//	@Description	Delete a bank vault by the provided ID
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Bank vault ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Bank vault delete request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/products/delete/{id} [delete]
func (h *handler) DeleteBankVault(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "deleteBankVault", "handler", "bankVault")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("DeleteBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DeleteBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("bank_vault.id", id))

	_, err = h.service.DeleteBankVault(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DeleteBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return

	}

	h.logger.Infof("[DeleteBankVault] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultDeleteRequestSubmitted, nil)
}

// DisableBankVault
//
//	@Summary		Disable Bank Vault
//	@Description	Disable a bank vault by the provided ID
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Bank vault ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Bank vault disable request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/products/disable/{id} [patch]
func (h *handler) DisableBankVault(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "disableBankVault", "handler", "bankVault")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("DisableBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DisableBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("bank_vault.id", id))

	err = h.service.DisableBankVault(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DisableBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[DisableBankVault] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultDisableRequestSubmitted, nil)
}

// EnableBankVault
//
//	@Summary		Enable Bank Vault
//	@Description	Enable a bank vault by the provided ID
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Bank vault ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Bank vault enable request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/products/enable/{id} [patch]
func (h *handler) EnableBankVault(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "enableBankVault", "handler", "bankVault")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("EnableBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[EnableBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("bank_vault.id", id))

	err = h.service.EnableBankVault(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[EnableBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[EnableBankVault] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultEnableRequestSubmitted, nil)
}

// GetAllLockedBankVaults
//
//	@Summary		Get All Locked Bank Vaults
//	@Description	Retrieve all locked bank vaults with pagination and filters
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"Page number"
//	@Param			per_page	query	int		false	"Items per page"
//	@Param			search		query	string	false	"Search keyword"
//	@Param			filters		query	string	false	"Additional filters as JSON string"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Locked bank vaults retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/bank-vaults [get]
func (h *handler) GetAllLockedBankVaults(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getAllLockedBankVaults", "handler", "bankVault")
	defer span.End()
	params := common_utils.ExtractFilterParams(r)
	result, err := h.service.FindAllBankLockedVaultsWithPagination(ctx, params)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAllLockedBankVaults] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("bank_vault.locked.count", len(result.Data)))
	h.logger.Infof("[GetAllLockedBankVaults] retrieved %d locked bank vaults", len(result.Data))
	localization.SendSuccessResponse(w, localization.SuccessAllBankLockedVaultsRetrievedSuccessfully, result)
}

// On development
// func (h *handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
// 	transactionReference, err := common_utils.ExtractID(w, r)
// 	if transactionReference == "" {
// 		localization.SendErrorByCodeResponse(w, err.Error())
// 		return
// 	}

// 	if err != nil {
// 		h.logger.Errorf("[GetBankLockedVaults] extract ID: %v", err)
// 		localization.SendErrorByCodeResponse(w, err.Error())
// 		return
// 	}

// 	result, err := h.service.GetTransactions(r.Context(), transactionReference)
// 	if err != nil {
// 		h.logger.Errorf("[GetBankLockedVaults] service: %v", err)
// 		localization.SendErrorByCodeResponse(w, err.Error())
// 		return
// 	}
// 	h.logger.Infof("Bank Locked vaults retrieved with FT: %s", transactionReference)
// 	localization.SendSuccessResponse(w, localization.SuccessBankLockedVaultsRetrievedSuccessfully, result)
// }

// GetAllGroupVaults
//
//	@Summary		Get All Group Vaults
//	@Description	Retrieve all group vaults with pagination and filters
//	@Tags			Bank Vault
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"Page number"
//	@Param			per_page	query	int		false	"Items per page"
//	@Param			search		query	string	false	"Search keyword"
//	@Param			filters		query	string	false	"Additional filters as JSON string"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Group vaults retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault/group-vaults [get]
func (h *handler) GetAllGroupVaults(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getAllGroupVaults", "handler", "bankVault")
	defer span.End()
	params := common_utils.ExtractFilterParams(r)
	results, err := h.service.FindAllGroupVaultsWithPagination(ctx, params)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAllGroupVaults] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("bank_vault.group.count", len(results.Data)))
	h.logger.Infof("[GetAllGroupVaults] retrieved %d group vaults", len(results.Data))
	localization.SendSuccessResponse(w, localization.SuccessGroupVaultsRetrievedSuccessfully, results)
}

// On development
// func (h *handler) GetGroupVault(w http.ResponseWriter, r *http.Request) {
// 	id, err := common_utils.ExtractID(w, r)
// 	if id == "" {
// 		localization.SendErrorByCodeResponse(w, err.Error())
// 		return
// 	}

// 	if err != nil {
// 		h.logger.Errorf("[GetGroupVault] extract ID: %v", err)
// 		localization.SendErrorByCodeResponse(w, err.Error())
// 		return
// 	}

// 	result, err := h.service.GetGroupVault(r.Context(), id)
// 	if err != nil {
// 		h.logger.Errorf("[GetGroupVault] service: %v", err)
// 		localization.SendErrorByCodeResponse(w, err.Error())
// 		return
// 	}
// 	h.logger.Infof("Bank Locked vaults retrieved with ID: %s", id)
// 	localization.SendSuccessResponse(w, localization.SuccessGroupVaultRetrievedSuccessfully, result)
// }
