package bankHandler

import (
	"cbe-super-app-cps-action/internal/constants"
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	"cbe-super-app-cps-action/internal/constants/interfaces/bank"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	bank_core "cbe-super-app-cps-action/internal/handlers/rest/http/bank/core"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type paginatedBankResp types.PaginatedResponse[[]*bank_dto.BankResponse]
type bankAdapter struct {
	bankService service.BankService
	logger      utils.Logger
}

func InitBankAdapter(bankApplication service.BankService, logger utils.Logger) bank.BankHandler {
	return &bankAdapter{
		logger:      logger,
		bankService: bankApplication,
	}
}

// CreateOneBank godoc
//
//	@Summary		Create bank (maker)
//	@Description	Submit a bank create request with logo.
//	@Tags			Banks
//	@Accept			mpfd
//	@Produce		json
//	@Param			name	formData	string									true	"Bank name"	example("Commercial Bank")
//	@Param			code	formData	string									true	"Bank code"	example("CBE")
//	@Param			type	formData	string									true	"Bank type"	example("Commercial")
//	@Param			bic		formData	string									true	"Bank BIC"	example("CBETETAA")
//	@Param			logo	formData	file									true	"Bank logo (<=2MB; jpeg/png/gif/webp)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Bank create request sent"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid input or logo"
//	@Failure		409		{object}	localization.StandardResponse{data=nil}	"Duplicate"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/banks [post]
func (b *bankAdapter) CreateOneBank(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "createOneBank", "handler", "bank")
	defer span.End()
	var bankRequest bank_dto.CreateBankRequest

	file, fileHeader, err := bank_core.ParseMultipartFormFile(r, "logo", 10<<20, string(constants.CREATE), b.logger)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
		return
	}
	defer file.Close()

	bankRequest.Name = r.FormValue("name")
	bankRequest.Code = r.FormValue("code")
	bankRequest.BIC = r.FormValue("bic")
	bankRequest.Type = constants.FinancialInstitutionType(r.FormValue("type"))
	bankRequest.Logo = fileHeader

	if response_code := bank_core.ValidateBankRequest(r, &bankRequest); response_code.Code != "" {
		b.logger.Errorf("invalid input", response_code)
		span.SetAttributes(attribute.String("invalid input", response_code.Code))
		localization.SendErrorResponse(w, response_code, nil, nil)
		return
	}

	span.SetAttributes(
		attribute.String("bank.code", bankRequest.Code),
		attribute.String("bank.name", bankRequest.Name),
	)

	err = b.bankService.CreateOneBank(ctx, bankRequest)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[CreateOneBank] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	b.logger.Infof("[CreateOneBank] request sent successfully for bank code: %s", bankRequest.Code)
	localization.SendSuccessResponse(w, localization.SuccessBankCreatedRequestSent, nil)
}

// DeleteOneBank godoc
//
//	@Summary		Delete bank (maker)
//	@Description	Submit a delete request for a bank by ID.
//	@Tags			Banks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Bank ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Delete request created"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Bank not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/banks/{id} [delete]
func (b *bankAdapter) DeleteOneBank(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "deleteOneBank", "handler", "bank")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("bank.id", id))

	err := b.bankService.DeleteOneBank(ctx, id)

	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[DeleteOneBank] service error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorBankDeleteRequestFailed.Code)
		return
	}

	b.logger.Infof("[DeleteOneBank] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessDeleteRequestCreated, nil)
}

// Disable godoc
//
//	@Summary		Disable bank (checker)
//	@Description	Approve disable request for a bank.
//	@Tags			Banks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Bank ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Bank disable request created"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Bank not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/banks/{id}/disable [patch]
func (b *bankAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "disableBank", "handler", "bank")
	defer span.End()
	id := chi.URLParam(r, "id")
	span.SetAttributes(attribute.String("bank.id", id))
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := b.bankService.EnableOrDisableBank(ctx, id, false)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	b.logger.Infof("[Disable] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankDisableRequestCreated, nil)
}

// Enable godoc
//
//	@Summary		Enable bank (checker)
//	@Description	Approve enable request for a bank.
//	@Tags			Banks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Bank ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Bank enable request created"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Bank not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/banks/{id}/enable [patch]
func (b *bankAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "enableBank", "handler", "bank")
	defer span.End()
	id := chi.URLParam(r, "id")
	span.SetAttributes(attribute.String("bank.id", id))
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := b.bankService.EnableOrDisableBank(ctx, id, true)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	b.logger.Infof("[Enable] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankEnableRequestCreated, nil)
}

// GetAllBank godoc
//
//	@Summary		List banks
//	@Description	Retrieve banks with pagination and optional search.
//	@Tags			Banks
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int														false	"Page number"		default(1)	minimum(1)	example(1)
//	@Param			per_page	query		int														false	"Items per page"	default(10)	minimum(1)	maximum(100)	example(10)
//	@Param			search		query		string													false	"Search term"		example("CBE")
//	@Success		200			{object}	localization.StandardResponse{data=paginatedBankResp}	"Banks fetched"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}					"Server error"
//	@Security		BearerAuth
//	@Router			/banks [get]
func (b *bankAdapter) GetAllBank(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getAllBanks", "handler", "bank")
	defer span.End()
	filterParams := common_utils.ExtractFilterParams(r)

	banks, err := b.bankService.GetAllBank(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[GetAllBank] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("bank.count", len(banks.Data)))
	b.logger.Infof("[GetAllBank] retrieved %d banks", len(banks.Data))
	localization.SendSuccessResponse(w, localization.SuccessGetAllBanks, banks)
}

// GetOneBank godoc
//
//	@Summary		Get bank by ID
//	@Description	Retrieve a single bank by ID.
//	@Tags			Banks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string														true	"Bank ID"
//	@Success		200	{object}	localization.StandardResponse{data=bank_dto.BankResponse}	"Bank fetched"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}						"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}						"Bank not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}						"Server error"
//	@Security		BearerAuth
//	@Router			/banks/{id} [get]
func (b *bankAdapter) GetOneBank(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getOneBank", "handler", "bank")
	defer span.End()
	id := chi.URLParam(r, "id")
	span.SetAttributes(attribute.String("bank.id", id))

	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	bank, err := b.bankService.GetOneBank(ctx, id)

	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[GetOneBank] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[GetOneBank] bank retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessGetOneBank, bank)
}

// UpdateLogo godoc
//
//	@Summary		Update bank logo (maker)
//	@Description	Upload a new logo for the bank.
//	@Tags			Banks
//	@Accept			mpfd
//	@Produce		json
//	@Param			id		path		string									true	"Bank ID"
//	@Param			logo	formData	file									true	"Bank logo (<=2MB; jpeg/png/gif/webp)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Logo uploaded"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid logo"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Bank not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/banks/{id}/logo [patch]
func (b *bankAdapter) UpdateLogo(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "updateBankLogo", "handler", "bank")
	defer span.End()
	var uploadLogo bank_dto.UpdateLogo
	id := chi.URLParam(r, "id")
	span.SetAttributes(attribute.String("bank.id", id))
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	file, fileHeader, err := bank_core.ParseMultipartFormFile(r, "logo", 10<<20, string(constants.CREATE), b.logger)
	if err != nil {
		b.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
		return
	}
	defer file.Close()
	uploadLogo.Logo = fileHeader
	uploadLogo.ID = id

	err = b.bankService.UpdateLogo(ctx, id, uploadLogo)

	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[UpdateLogo] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[UpdateLogo] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessFileUploadedToMinIO, nil)
}

// UpdateOneBank godoc
//
//	@Summary		Update bank (maker)
//	@Description	Update bank name/code/bic.
//	@Tags			Banks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Bank ID"
//	@Param			request	body		bank_dto.UpdateBankRequest				true	"Update payload"	example({"name":"New Name","code":"NEW","bic":"NEWBIC"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Bank update request sent"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid input"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Bank not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/banks/{id} [patch]
func (b *bankAdapter) UpdateOneBank(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "updateOneBank", "handler", "bank")
	defer span.End()
	id := chi.URLParam(r, "id")
	span.SetAttributes(attribute.String("bank.id", id))
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	var updateRequest bank_dto.UpdateBankRequest
	file, fileHeader, err := bank_core.ParseMultipartFormFile(r, "logo", 10<<20, string(constants.Update), b.logger)
	if err != nil {
		span.RecordError(err)
		if strings.Contains(err.Error(), "file not found") {
			b.logger.Infof("No logo uploaded; skipping logo update")
		} else {
			b.logger.Errorf("error parsing file: %v", err)
			localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
			return
		}
	}
	if file == nil || fileHeader == nil {
		b.logger.Infof("No logo uploaded; skipping logo update")
	} else {
		defer file.Close()
		updateRequest.Logo = fileHeader
	}

	updateRequest.Name = r.FormValue("name")
	updateRequest.Code = r.FormValue("code")
	updateRequest.BIC = r.FormValue("bic")
	updateRequest.Type = constants.FinancialInstitutionType(r.FormValue("type"))

	if response_code := bank_core.ValidateBankRequest(r, &updateRequest); response_code.Code != "" {
		span.SetAttributes(attribute.String("invalid input", response_code.Code))
		b.logger.Errorf("invalid input", response_code)
		localization.SendErrorResponse(w, response_code, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("bank.id", id))

	if err = b.bankService.UpdateOneBank(ctx, id, updateRequest); err != nil {
		span.RecordError(err)
		b.logger.Errorf("[UpdateOneBank] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[UpdateOneBank] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankUpdatedRequestSent, nil)
}
