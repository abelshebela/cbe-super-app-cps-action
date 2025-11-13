package encryption

import (
	"encoding/json"
	"net/http"

	encryptionDto "cbe-super-app-cps-action/internal/constants/dto/encryption"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

func (enc *encryptionHandler) Encrypt(w http.ResponseWriter, r *http.Request) {
	var req encryptionDto.EncryptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		enc.logger.Errorf("[Encryption] validation failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	result, _, err := enc.svc.LocalEncryptPassword(req, "enc", "enc", "enc")
	if err != nil {
		enc.logger.Errorf("[Encryption] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEncryptionGenerated, result)
}
