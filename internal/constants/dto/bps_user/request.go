package bpsuser

// type FullName struct {
// 	FirstName  string `json:"first_name"`
// 	MiddleName string `json:"middle_name"`
// 	LastName   string `json:"last_name"`
// }

type BPSUserCreateRequest struct {
	UserID      string `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	// Role        string   `json:"role"`
	BranchCode []string `json:"branch_code"`
	// HomeBranch  string   `json:"home_branch"`
	// Realm       string   `json:"realm"`
	Enabled bool `json:"enabled"`
}

type BPSUserUpdateRequest struct {
	UserID      string `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	// Role        string   `json:"role"`
	BranchCode []string `json:"branch_code"`
	// HomeBranch  string   `json:"home_branch"`
	// Realm       string   `json:"realm"`
	Enabled bool `json:"enabled"`
}
