package accountlookup

import (
	"encoding/json"
	"net/http"

	lookup "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_lookup"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AccountlookupStore struct {
	Application lookup.UserSearchService
	logger      shared_utils.Logger
}

func NewMiniAppMerchantAdapter(app lookup.UserSearchService, logger shared_utils.Logger) *AccountlookupStore {
	return &AccountlookupStore{Application: app, logger: logger}
}

func (a *AccountlookupStore) SearchUser(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Number string `json:"number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		a.logger.Errorf("invalid JSON body: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	if reqBody.Number == "" {
		a.logger.Warnf("missing 'number' field in request body")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	res, err := a.Application.SearchUser(r.Context(), reqBody.Number)
	if err != nil {
		a.logger.Errorf("search error: %v", err)
		common_util.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	common_util.WriteSuccessResponse(w, res, "Successfully retrieved the value")
}
