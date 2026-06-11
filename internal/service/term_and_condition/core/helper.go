package term_and_condition_core

import (
	"encoding/json"

	imodel "cbe-super-app-cps-action/internal/constants/model"
)

func MapFromAction(raw interface{}) (imodel.AccountOpeningTerms, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return imodel.AccountOpeningTerms{}, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return imodel.AccountOpeningTerms{}, err
	}

	t := imodel.AccountOpeningTerms{}

	if v, ok := m["id"].(string); ok {
		t.ID = v
	}
	if v, ok := m["account_product_id"].(string); ok {
		t.AccountProductID = v
	}
	if v, ok := m["activation_time"].(string); ok {
		t.ActivationTime = v
	}
	if v, ok := m["version_label"].(string); ok {
		t.VersionLabel = v
	}
	if v, ok := m["terms_and_conditions_path"].(string); ok {
		t.TermsAndConditionsPath = v
	}
	if v, ok := m["is_enabled"].(bool); ok {
		t.IsEnabled = v
	}
	if v, ok := m["is_deleted"].(bool); ok {
		t.IsDeleted = v
	}

	return t, nil
}
