package model

type AccountSubType struct {
	ID                 string `json:"id"                   bson:"id"`
	AccountType        string `json:"account_type"         bson:"account_type"`
	Gender             string `json:"gender"               bson:"gender"`
	AccountSubTypeName string `json:"account_sub_type_name" bson:"account_sub_type_name"`
	AccountSubTypeCode string `json:"account_sub_type_code" bson:"account_sub_type_code"`
	IsDeleted          bool   `json:"is_deleted"           bson:"is_deleted"`
	IsEnabled          int    `json:"is_enabled"           bson:"is_enabled"`
	CreatedAt          string `json:"created_at"           bson:"created_at"`
	LastModifiedAt     string `json:"last_modified_at"     bson:"last_modified_at"`
}
