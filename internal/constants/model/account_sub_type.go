package model

type AccountSubType struct {
	ID                 string `json:"id"`
	AccountType        string `json:"account_type"`
	Gender             string `json:"gender"`
	AccountSubTypeName string `json:"account_sub_type_name"`
	AccountSubTypeCode string `json:"account_sub_type_code"`
	IsDeleted          bool   `json:"is_deleted"`
	IsEnabled          int    `json:"is_enabled"`
	CreatedAt          string `json:"created_at"`
	LastModifiedAt     string `json:"last_modified_at"`
}
