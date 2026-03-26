package ussd_merchant

import (
	"cbe-super-app-cps-action/internal/constants"
	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	ussd_merchant_interface "cbe-super-app-cps-action/internal/constants/interfaces/ussd_merchant"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/ussd_merchant/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"mime/multipart"

	"cbe-super-app-cps-action/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UssdMerchantHandler struct {
	UssdMerchantService service.UssdMerchantService
	Logger              utils.Logger
}

func NewUssdMerchantHandler(ussdMerchantService service.UssdMerchantService, logger utils.Logger) ussd_merchant_interface.UssdMerchantInbound {
	return &UssdMerchantHandler{
		UssdMerchantService: ussdMerchantService,
		Logger:              logger,
	}
}

// CreateUssdMerchant godoc
//
//	@Summary		Create USSD merchant
//	@Description	Create a new USSD merchant with the provided information. Requires multipart/form-data with logo image.
//	@Tags			USSD Merchant
//	@Accept			mpfd
//	@Produce		json
//	@Param			settlement_method	formData	string	true	"Settlement method"
//	@Param			name				formData	string	true	"Merchant name"
//	@Param			phone_number		formData	string	true	"Phone number"
//	@Param			email				formData	string	true	"Email address"
//	@Param			service				formData	string	true	"Service"
//	@Param			account_number		formData	string	true	"Account number"
//	@Param			logo				formData	file	true	"Logo image (<=10MB; jpeg/png/gif/webp)"
//	@Success		200					{object}	localization.StandardResponse{data=nil}	"USSD merchant created successfully"
//	@Failure		400					{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500					{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/ussd_merchant [post]
func (u *UssdMerchantHandler) CreateUssdMerchant(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "CreateUssdMerchantHandler", "handler", "ussdMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, u.Logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var CreateDto ussd_merchant_dto.CreateUssdMerchantRequest
	file, fileHeader, err := core.ParseMultipartFormFile(r, "logo", 10<<20, string(constants.CREATE), u.Logger)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CreateUssdMerchantRequestHandler] error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
		return
	}
	defer file.Close()

	CreateDto.SettlementMethod = r.FormValue("settlement_method")
	CreateDto.Name = r.FormValue("name")
	CreateDto.PhoneNumber = local_util.FormatPhoneNumber(r.FormValue("phone_number"))
	CreateDto.Service = r.FormValue("service")
	CreateDto.Email = r.FormValue("email")
	CreateDto.AccountNumber = r.FormValue("account_number")
	CreateDto.Logo = fileHeader

	if err := CreateDto.Validate(); err != nil {
		log.Errorf("[CreateUssdMerchantHandler] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := u.UssdMerchantService.CreateUssdMerchant(ctx, CreateDto); err != nil {
		log.Errorf("[CreateUssdMerchantHandler] failed to create ussd merchant error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[CreateUssdMerchantHandler] request sent successfully: is_maker_only: %v", md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessUssdMerchantCreated, nil)
		return
	}

	log.Infof("[CreateUssdMerchantHandler] successfully created ussd merchant")
	localization.SendSuccessResponse(w, localization.SuccessUssdMerchantRequestCreated, "Ussd merchant created successfully")
}

// UpdateUssdMerchant godoc
//
//	@Summary		Update USSD merchant
//	@Description	Update an existing USSD merchant by ID. Logo is optional. Requires multipart/form-data.
//	@Tags			USSD Merchant
//	@Accept			mpfd
//	@Produce		json
//	@Param			id					path		string	false	"USSD Merchant ID"
//	@Param			settlement_method	formData	string	false	"Settlement method"
//	@Param			name				formData	string	false	"Merchant name"
//	@Param			phone_number		formData	string	false	"Phone number"
//	@Param			email				formData	string	false	"Email address"
//	@Param			service				formData	string	false	"Service"
//	@Param			account_number		formData	string	false	"Account number"
//	@Param			logo				formData	file	false	"Logo image (<=10MB; jpeg/png/gif/webp)"
//	@Success		200					{object}	localization.StandardResponse{data=nil}	"USSD merchant updated successfully"
//	@Failure		400					{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404					{object}	localization.StandardResponse{data=nil}	"USSD merchant not found"
//	@Failure		500					{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/ussd_merchant/{id} [patch]
func (u *UssdMerchantHandler) UpdateUssdMerchant(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateUssdMerchantHandler", "handler", "ussdMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, u.Logger)
	var fileHeader *multipart.FileHeader
	var file multipart.File
	var err error
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[UpdateUssdMerchantHandler] missing id parameter in request path")
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var req ussd_merchant_dto.UpdateUssdMerchantRequest

	// logo is optional on update; only parse/validate when the multipart file is provided
	_, _, formFileErr := r.FormFile("logo")
	if formFileErr == nil {
		file, fileHeader, err = core.ParseMultipartFormFile(r, "logo", 10<<20, string(constants.CREATE), u.Logger)
		if err != nil {
			span.RecordError(err)
			log.Errorf("[UpdateUssdMerchantHandler] error parsing file: %v", err)
			localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
			return
		}
		defer file.Close()
		req.Logo = fileHeader
	} else if formFileErr != http.ErrMissingFile {
		span.RecordError(formFileErr)
		log.Errorf("[UpdateUssdMerchantHandler] error reading form file: %v", formFileErr)
		localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
		return
	}
	// file, fileHeader, err = core.ParseMultipartFormFile(r, "logo", 10<<20, string(constants.CREATE), u.Logger)
	// if err != nil {
	// 	span.RecordError(err)
	// 	log.Errorf("[CreateUssdMerchantRequestHandler] error parsing file: %v", err)
	// 	localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
	// 	return
	// }
	// defer file.Close()

	req.SettlementMethod = r.FormValue("settlement_method")
	req.Name = r.FormValue("name")
	req.PhoneNumber = local_util.FormatPhoneNumber(r.FormValue("phone_number"))
	req.Service = r.FormValue("service")
	req.Email = r.FormValue("email")
	req.AccountNumber = r.FormValue("account_number")
	if fileHeader != nil {
		req.Logo = fileHeader
	}

	// Handle optional logo via multipart if present and validate

	if err := req.Validate(); err != nil {
		log.Errorf("[UpdateUssdMerchantHandler] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if req.PhoneNumber != "" {
		req.PhoneNumber = local_util.FormatPhoneNumber(req.PhoneNumber)
	}

	if err := u.UssdMerchantService.UpdateUssdMerchant(ctx, id, req); err != nil {
		log.Errorf("[UpdateUssdMerchantHandler] failed to update ussd merchant: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[UpdateUssdMerchantHandler] request sent successfully: is_maker_only: %v", md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessUssdMerchantUpdated, nil)
		return
	}

	log.Infof("[UpdateUssdMerchantHandler] successfully updated ussd merchant")
	localization.SendSuccessResponse(w, localization.SuccessUssdMerchantUpdateRequestCreated, nil)
}
func (u *UssdMerchantHandler) EnableUssdMerchant(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "EnableUssdMerchantHandler", "handler", "ussdMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, u.Logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[EnableUssdMerchantHandler] missing id parameter in request path")
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	if err := u.UssdMerchantService.EnableUssdMerchant(ctx, id); err != nil {
		log.Errorf("[EnableUssdMerchantHandler] failed to enable ussd merchant: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[EnableUssdMerchantHandler] request sent successfully: is_maker_only: %v", md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessUssdMerchantEnabled, nil)
		return
	}

	log.Infof("[EnableUssdMerchantHandler] successfully updated ussd merchant")
	localization.SendSuccessResponse(w, localization.SuccessUssdMerchantEnableRequestCreated, nil)
}

// DisableUssdMerchant godoc
//
//	@Summary		Disable USSD merchant
//	@Description	Disable a USSD merchant by ID
//	@Tags			USSD Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"USSD Merchant ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"USSD merchant disabled successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"USSD merchant not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/ussd_merchant/disable/{id} [patch]
func (u *UssdMerchantHandler) DisableUssdMerchant(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "DisableUssdMerchantHandler", "handler", "ussdMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, u.Logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[DisableUssdMerchantHandler] missing id parameter in request path")
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	if err := u.UssdMerchantService.DisableUssdMerchant(ctx, id); err != nil {
		log.Errorf("[DisableUssdMerchantHandler] failed to enable ussd merchant: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[DisableUssdMerchantHandler] request sent successfully: is_maker_only: %v", md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessUssdMerchantDisabled, nil)
		return
	}

	log.Infof("[DisableUssdMerchantHandler] successfully updated ussd merchant")
	localization.SendSuccessResponse(w, localization.SuccessUssdMerchantDisableRequestCreated, nil)
}
func (u *UssdMerchantHandler) GetUssdMerchant(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "GetUssdMerchantHandler", "handler", "ussdMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, u.Logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[GetUssdMerchantHandler] missing id parameter in request path")
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidID.Code)
		return
	}

	result, err := u.UssdMerchantService.GetUssdMerchantByID(ctx, id)
	if err != nil {
		log.Errorf("[GetUssdMerchantHandler] failed to fetch ussd merchant: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUssdMerchantFetched, result)
}

// GetAllUssdMerchant godoc
//
//	@Summary		Get all USSD merchants
//	@Description	Retrieve all USSD merchants with pagination and optional search
//	@Tags			USSD Merchant
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			per_page	query		int		false	"Items per page"	default(10)
//	@Param			search		query		string	false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"USSD merchants retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/ussd_merchant [get]
func (u *UssdMerchantHandler) GetAllUssdMerchant(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "GetAllUssdMerchantHandler", "handler", "ussdMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, u.Logger)

	filter := local_util.ExtractFilterParams(r)
	result, err := u.UssdMerchantService.FindAllWithPagination(ctx, filter)
	if err != nil {
		log.Errorf("[GetAllUssdMerchantHandler] failed to list ussd merchants: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUssdMerchantFetched, result)
}
