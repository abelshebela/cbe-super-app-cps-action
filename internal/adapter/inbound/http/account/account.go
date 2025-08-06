package account

import (
	"encoding/json"
	"net/http"

	accountAppDto "cbe-super-app-member-users/internal/application/dto"
	"cbe-super-app-member-users/pkgs/utils"
)

func (h AccountAdapter) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.ContextKey("user_id")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	req := accountAppDto.CreateAccountRequest{UserID: userID}

	response, err := h.Application.CreateAccount(r.Context(), req.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "NOT_FOUND":
			status = http.StatusNotFound
		case "USER_KYC_LEVEL_ZERO":
			status = http.StatusBadRequest
		case "USER_PHONE_EXISTS":
			status = http.StatusConflict
		case "API_REQUEST_FAILED":
			status = http.StatusInternalServerError
		case "GENERAL_SIF_GENERATION_FAILED":
			status = http.StatusInternalServerError
		case "GENERAL_DB_UPDATE_FAILED":
			status = http.StatusInternalServerError
		case "GENERAL_DB_INSERT_FAILED":
			status = http.StatusInternalServerError
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"status":          "success",
		"message":         "Account created successfully",
		"customer_number": response.CustomerNumber,
		"account_number":  response.AccountNumber,
	})
}

func (h AccountAdapter) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	resp := struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: status,
		Data:   data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
