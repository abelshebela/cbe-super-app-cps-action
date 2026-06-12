package account_product_core

import (
	"encoding/json"

	imodel "cbe-super-app-cps-action/internal/constants/model"
)

// MapFromAction deserialises a CPS CurrentAction payload back to AccountProduct.
func MapFromAction(raw interface{}) (imodel.AccountProduct, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return imodel.AccountProduct{}, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return imodel.AccountProduct{}, err
	}

	ap := imodel.AccountProduct{}

	if v, ok := m["id"].(string); ok {
		ap.ID = v
	}
	if v, ok := m["cbs_product_code"].(string); ok {
		ap.CBSProductCode = v
	}
	if v, ok := m["product_name"].(string); ok {
		ap.ProductName = v
	}
	if v, ok := m["product_tag_line"].(string); ok {
		ap.ProductTagLine = v
	}
	if v, ok := m["product_line"].(string); ok {
		ap.ProductLine = v
	}
	if v, ok := m["account_category_id"].(string); ok {
		ap.AccountCategoryID = v
	}
	if v, ok := m["account_currency"].(string); ok {
		ap.AccountCurrency = v
	}
	if v, ok := m["minimum_opening_balance"].(float64); ok {
		ap.MinimumOpeningBalance = v
	}
	if v, ok := m["minimum_maintenance_fee"].(float64); ok {
		ap.MinimumMaintenanceFee = v
	}
	if v, ok := m["interest_fee"].(float64); ok {
		ap.InterestFee = v
	}
	if v, ok := m["faq_url"].(string); ok {
		ap.FaqURL = v
	}
	if v, ok := m["product_features"].(string); ok {
		ap.ProductFeatures = v
	}
	if v, ok := m["has_physical_card"].(bool); ok {
		ap.HasPhysicalCard = v
	}
	if v, ok := m["has_virtual_card"].(bool); ok {
		ap.HasVirtualCard = v
	}
	if v, ok := m["product_icon"].(string); ok {
		ap.ProductIcon = v
	}
	if v, ok := m["product_cover_image"].(string); ok {
		ap.ProductCoverImage = v
	}
	if v, ok := m["is_enabled"].(bool); ok {
		ap.IsEnabled = v
	}
	if v, ok := m["is_deleted"].(bool); ok {
		ap.IsDeleted = v
	}

	return ap, nil
}
