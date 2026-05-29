package walletDto

import (
	"mime/multipart"
)

type WalletRequest struct {
	Name           string                `form:"name" json:"name"`
	Avatar         *multipart.FileHeader `form:"avatar" json:"avatar"`
	UniqueCode     string                `form:"unique_code" json:"unique_code"`
	Self           *bool                 `form:"self" json:"self"`
	Other          *bool                 `form:"other" json:"other"`
	Agent          *bool                 `form:"agent" json:"agent"`
	SelfServiceID  string                `form:"self_service_id" json:"self_service_id"`
	OtherServiceID string                `form:"other_service_id" json:"other_service_id"`
	AgentServiceID string                `form:"agent_service_id" json:"agent_service_id"`
}
