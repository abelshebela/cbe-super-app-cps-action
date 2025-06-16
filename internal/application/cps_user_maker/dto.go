package cpsusermaker

type CreateUserRequest struct {
    UserName           string   `json:"user_name"`
    FullName           string   `json:"full_name"`
    PhoneNumber        string   `json:"phone_number"`
    UserRole           string   `json:"user_role"`
    Department         string   `json:"department"`
    PermissionCategory []string `json:"permission_category"`
    PermissionGroups   []string `json:"permission_groups"`
}
