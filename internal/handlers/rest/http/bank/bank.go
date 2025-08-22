package bank

import (
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	"cbe-super-app-cps-action/internal/constants/interfaces/bank"
	"cbe-super-app-cps-action/internal/constants/localization"
	bank_core "cbe-super-app-cps-action/internal/handlers/rest/http/bank/core"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type bankAdapter struct {
	bankService service.BankService
	logger      utils.Logger
}

// CreateOneBank implements bank.BankAdapter.
func (b *bankAdapter) CreateOneBank(w http.ResponseWriter, r *http.Request) {
	var bankRequest bank_dto.CreateBankRequest

	file, fileHeader, err := bank_core.ParseMultipartFormFile(r, "logo", 10<<20)
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

	err = b.bankService.CreateOneBank(r.Context(), bankRequest)
	if err != nil {
		b.logger.Errorf("bank create request failed", err)
		localization.SendErrorByCodeResponse(w, err.Error())
	}
	localization.SendSuccessResponse(w, localization.SuccessBankCreatedRequestSent, nil)
}

// DeleteOneBank implements bank.BankAdapter.
func (b *bankAdapter) DeleteOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := b.bankService.DeleteOneBank(r.Context(), id)

	if err != nil {
		b.logger.Errorf("bank delete request failed", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorBankDeleteRequestFailed.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDeleteRequestCreated, nil)
}

// Disable implements bank.BankAdapter.
func (b *bankAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := b.bankService.EnableOrDisableBank(r.Context(), id, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorBankDisableRequest.Code)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDisableRequestCreated, nil)
}

// Enable implements bank.BankAdapter.
func (b *bankAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := b.bankService.EnableOrDisableBank(r.Context(), id, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorBankEnableRequestFailed.Code)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEnableRequestCreated, nil)
}

// GetAllBank implements bank.BankAdapter.
func (b *bankAdapter) GetAllBank(w http.ResponseWriter, r *http.Request) {
	filterParams := common_utils.ExtractFilterParams(r)

	banks, err := b.bankService.GetAllBank(r.Context(), filterParams)
	if err != nil {
		b.logger.Errorf("get all banks failed", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorGetAllBanksFailed.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessGetAllBanks, banks)
}

// GetOneBank implements bank.BankAdapter.
func (b *bankAdapter) GetOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	bank, err := b.bankService.GetOneBank(r.Context(), id)

	if err != nil {
		b.logger.Errorf("get bank by id failed", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorGetOneBank.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessGetOneBank, bank)
}

// UpdateLogo implements bank.BankAdapter.
func (b *bankAdapter) UpdateLogo(w http.ResponseWriter, r *http.Request) {
	var uploadLogo bank_dto.UpdateLogo
	id := chi.URLParam(r, "id")
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	file, fileHeader, err := bank_core.ParseMultipartFormFile(r, "logo", 10<<20)
	if err != nil {
		b.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
		return
	}
	defer file.Close()
	uploadLogo.Logo = fileHeader
	uploadLogo.ID = id

	err = b.bankService.UpdateLogo(r.Context(), id, uploadLogo)

	if err != nil {
		b.logger.Errorf("Error uploading bank logo", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFileUploadFailed.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessFileUploadedToMinIO, nil)
}

// UpdateOneBank implements bank.BankAdapter.
func (b *bankAdapter) UpdateOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	var updateRequest bank_dto.UpdateBankRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		b.logger.Errorf("invalid input", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	err := b.bankService.UpdateOneBank(r.Context(), id, updateRequest)
	if err != nil {
		b.logger.Errorf("bank delete request failed", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorBankUpdateFailed.Code)
	}

	localization.SendSuccessResponse(w, localization.SuccessBankUpdatedRequestSent, nil)
}

func InitBankAdapter(bankApplication service.BankService, logger utils.Logger) bank.BankAdapter {
	return &bankAdapter{
		logger:      logger,
		bankService: bankApplication,
	}
}
