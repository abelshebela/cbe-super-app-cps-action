package encryption

import (
	"encoding/json"
	"net/http"

	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/encryption"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type encryptionHandler struct {
	logger utils.Logger
	cfg    *config.VaultConfig
}

func NewEncryption(logger utils.Logger, config *config.VaultConfig) inbound.Encryption {
	return &encryptionHandler{
		logger: logger,
		cfg:    config,
	}
}

func (enc *encryptionHandler) Encrypt(w http.ResponseWriter, r *http.Request) {
	var req EncryptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	if err := req.Validate(); err != nil {
		local_util.SendErrorResponse(w, err.Error(), 400, nil)
		return
	}

	data := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		local_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}
	dataStr := string(dataBytes)
	encData, _, err := local_util.LocalEncryptPassword(dataStr, "enc", "enc", "enc", enc.cfg)
	if err != nil {
		local_util.SendErrorResponse(w, err, 409, nil)
		return
	}

	response := map[string]interface{}{
		"encryption": encData,
	}
	local_util.BaseResponseMaker(response, w, "Successfuly encryption generated", 200)
}
