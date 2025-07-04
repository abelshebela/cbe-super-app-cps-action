package entities

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/entities/type_definition"
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
