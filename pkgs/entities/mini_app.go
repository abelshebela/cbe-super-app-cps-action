package entities

import (
	"time"

	"cbe-super-app-member-users/pkgs/entities/type_definition"
)

type MiniApp struct {
	ID                string
	AppName           string
	AppIcon           string
	CommisonGLAccount string
	AppType           type_definition.AppType
	MerchantID        string
	ProductCode       []type_definition.ProductCode
	Credential        []type_definition.CredentialInformation
	IsEventMiniApp    bool
	IsThreeClick      bool
	Enabled           bool
	IsDeleted         bool
	CreatedAt         time.Time
	LastModifiedAt    time.Time
	DeletedAt         time.Time
}
