package encryption

import (
	"encoding/json"
	"net/http"

	encryptionDto "cbe-super-app-cps-action/internal/constants/dto/encryption"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type encryptionHandler struct {
	svc    service.EncryptionService
	logger utils.Logger
}

func InitEncryption(svc service.EncryptionService, logger utils.Logger) *encryptionHandler {
	return &encryptionHandler{
		svc:    svc,
		logger: logger,
	}
}

// Encrypt godoc
//
//	@Summary		Encrypt user password
//	@Description	Encrypt a user's password using the encryption service. Requires username and password in the request body.
//	@Tags			Encryption
//	@Accept			json
//	@Produce		json
//	@Param			request	body		encryption.EncryptionRequest										true	"Encryption request"	example({"username":"user123","password":"mypassword"})
//	@Success		200		{object}	localization.StandardResponse{data=encryption.EncryptionResponse}	"Password encrypted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}								"Bad request - invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}								"Server error"
//	@Security		BearerAuth
//	@Router			/encryption/encrypt [post]
func (enc *encryptionHandler) Encrypt(w http.ResponseWriter, r *http.Request) {
	_, span := local_util.TraceLogger(r.Context(), "handler", "encryptPassword", "handler", "encryption")
	defer span.End()
	var req encryptionDto.EncryptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		enc.logger.Errorf("[Encryption] failed to decode request: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		enc.logger.Errorf("[Encryption] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("encryption.username", req.Username))

	result, _, err := enc.svc.LocalEncryptPassword(req, "enc", "enc", "enc")
	if err != nil {
		span.RecordError(err)
		enc.logger.Errorf("[Encryption] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	enc.logger.Infof("[Encryption] password encrypted successfully for username: %s", req.Username)
	localization.SendSuccessResponse(w, localization.SuccessEncryptionGenerated, result)
}
