package cpsroles

type ServiceAccessInfo struct {
	Key            string `json:"key" bson:"key"`
	AccessListName string `json:"access_list_name" bson:"access_list_name"`
}

type Channel struct {
	Name     string `json:"name"`
	MaxLimit int64  `json:"max_limit"`
}

type GlobalLimitResponse struct {
	CustomerNumber   int64     `json:"customer_number"`
	CustomerMaxLimit int64     `json:"customer_max_limit"`
	Channels         []Channel `json:"channels"`
}

type ServiceLevelLimitResponse struct {
	Name             string `json:"name"`
	USSDTranFreq     string `json:"ussd_transaction_freq"`
	SuperAppTranFreq string `json:"super_app_transaction_freq"`
	USSDMaxLimit     string `json:"ussd_max_limit"`
	SuperAppMaxLimit string `json:"super_app_max_limit"`
}

// type CPSRolesResponse struct {
// 	ID               bson.ObjectID       `json:"id" bson:"_id,omitempty"`
// 	Name             string              `json:"name,omitempty" bson:"name,omitempty"`
// 	RoleCode         string              `json:"role_code,omitempty" bson:"role_code,omitempty"`
// 	Description      string              `json:"description,omitempty" bson:"description,omitempty"`
// 	Enabled          *bool               `json:"enabled" bson:"enabled"`
// 	MakerActions     []string            `json:"maker" bson:"maker"`
// 	CheckerActions   []string            `json:"checker" bson:"checker"`
// 	AuditorActions   []string            `json:"auditor" bson:"auditor"`
// 	EnabledServices  []ServiceAccessInfo `json:"enabled_services" bson:"-"`
// 	DisabledServices []ServiceAccessInfo `json:"disabled_services" bson:"-"`
// 	GlobalLimit      Limit               `json:"global_limit" bson:"-"`
// 	IsDeleted        bool                `json:"is_deleted" bson:"is_deleted"`
// 	CreatedAt        time.Time           `json:"created_at,omitempty" bson:"created_at,omitempty"`
// 	UpdatedAt        time.Time           `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
// 	DeletedAt        time.Time           `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
// }
