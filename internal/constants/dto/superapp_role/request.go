package superapproledto

type BulkAccessListByRoleRequest struct {
	AccessListIDs []string `json:"access_list_ids" validate:"required"`
}
