package account_sub_type_dto

type CreateAccountSubTypeRequest struct {
	AccountType        string `json:"account_type"`
	Gender             string `json:"gender"`
	AccountSubTypeName string `json:"account_sub_type_name"`
	AccountSubTypeCode string `json:"account_sub_type_code"`
}

type UpdateAccountSubTypeRequest struct {
	AccountType        string `json:"account_type"`
	Gender             string `json:"gender"`
	AccountSubTypeName string `json:"account_sub_type_name"`
	AccountSubTypeCode string `json:"account_sub_type_code"`
}
