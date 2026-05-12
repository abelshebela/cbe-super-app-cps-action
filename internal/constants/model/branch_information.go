package model

type BranchInformation struct {
	ID                  string `json:"id" bson:"id"`
	BranchCode          string `json:"branch_code" bson:"branch_code"`
	BranchName          string `json:"branch_name" bson:"branch_name"`
	BranchAddress       string `json:"branch_address" bson:"branch_address"`
	BranchOwner         string `json:"branch_owner" bson:"branch_owner"`
	IsEnabled           bool   `json:"is_enabled" bson:"is_enabled"`
	BranchAccountNumber string `json:"branch_account_number" bson:"branch_account_number"`
}
