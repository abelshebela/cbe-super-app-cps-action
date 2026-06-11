package account_product_category_dto

type CreateAPCRequest struct {
	ProductLine     string `json:"product_line"`
	CBSCategoryCode string `json:"cbs_category_code"`
	CategoryName    string `json:"category_name"`
	Description     string `json:"description"`
}

type UpdateAPCRequest struct {
	ProductLine     string `json:"product_line"`
	CBSCategoryCode string `json:"cbs_category_code"`
	CategoryName    string `json:"category_name"`
	Description     string `json:"description"`
}

type CPSAPCRequest struct {
	ID              string `json:"id"`
	ProductLine     string `json:"product_line"`
	CBSCategoryCode string `json:"cbs_category_code"`
	CategoryName    string `json:"category_name"`
	Description     string `json:"description"`
	IsEnabled       bool   `json:"is_enabled"`
	IsDeleted       bool   `json:"is_deleted"`
	CreatedAt       string `json:"created_at"`
	LastModifiedAt  string `json:"last_modified_at"`
}

type APCResponse struct {
	ID              string `json:"id"`
	ProductLine     string `json:"product_line"`
	CBSCategoryCode string `json:"cbs_category_code"`
	CategoryName    string `json:"category_name"`
	Description     string `json:"description"`
	IsEnabled       bool   `json:"is_enabled"`
	IsDeleted       bool   `json:"is_deleted"`
	CreatedAt       string `json:"created_at"`
	LastModifiedAt  string `json:"last_modified_at"`
}
