package users

import "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities"

type LinkedAccountDetail struct {
    AccountNumber     string `json:"account_number"`
    AccountBranchCode string `json:"account_branch_code"`
    LinkedBranch      string `json:"linked_branch"`
    IsAccountActive   bool   `json:"is_account_active"`
    Maker             string `json:"maker"`
    Checker           string `json:"checker"`
}

type LinkedAccountResponse struct {
    UserID         string                `json:"user_id"`
    FullName       entities.FullName     `json:"full_name"`
    LinkedAccounts []LinkedAccountDetail `json:"linked_accounts"`
}